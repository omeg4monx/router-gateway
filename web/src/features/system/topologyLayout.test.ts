import { describe, expect, test } from "bun:test";
import { buildTopologyLayout, groupTopologyProviders } from "./topologyLayout";

describe("features/system/topologyLayout", () => {
  test("groupTopologyProviders deduplicates connections by provider_id and attaches request stats", () => {
    const grouped = groupTopologyProviders(
      [
        { provider_id: "openai", protocol: "openai-chat", enabled: true },
        { provider_id: "openai", protocol: "openai-chat", enabled: false },
        { provider_id: "omeg4mon", protocol: "openai-chat", enabled: true },
      ],
      [
        { provider_id: "openai", requests: 14 },
        { provider_id: "omeg4mon", requests: 5 },
        { provider_id: "anthropic", requests: 2 },
      ],
      new Map([["omeg4mon", "Omeg4mon Custom"]]),
    );

    expect(grouped).toHaveLength(3);
    const openai = grouped.find((p) => p.id === "openai");
    expect(openai?.connectionCount).toBe(2);
    expect(openai?.activeConnections).toBe(1);
    expect(openai?.requests).toBe(14);

    const custom = grouped.find((p) => p.id === "omeg4mon");
    expect(custom?.label).toBe("Omeg4mon Custom");
    expect(custom?.requests).toBe(5);

    const historical = grouped.find((p) => p.id === "anthropic");
    expect(historical?.connectionCount).toBe(0);
    expect(historical?.requests).toBe(2);
  });

  test("buildTopologyLayout returns router-only node when no providers exist", () => {
    const layout = buildTopologyLayout([], new Set(), new Set(), new Set());
    expect(layout.nodes).toHaveLength(1);
    expect(layout.nodes[0]?.id).toBe("router");
    expect(layout.edges).toHaveLength(0);
  });

  test("buildTopologyLayout distributes provider nodes radially with active, last, and error edge states", () => {
    const providers = [
      {
        id: "openai",
        label: "OpenAI",
        protocol: "openai-chat",
        isCustom: false,
        connectionCount: 1,
        activeConnections: 1,
        requests: 10,
      },
      {
        id: "anthropic",
        label: "Anthropic",
        protocol: "anthropic",
        isCustom: false,
        connectionCount: 1,
        activeConnections: 1,
        requests: 4,
      },
      {
        id: "omeg4mon",
        label: "Omeg4mon",
        protocol: "openai-chat",
        isCustom: true,
        connectionCount: 1,
        activeConnections: 1,
        requests: 3,
      },
      {
        id: "gemini",
        label: "Gemini",
        protocol: "gemini",
        isCustom: false,
        connectionCount: 1,
        activeConnections: 1,
        requests: 1,
      },
    ];

    const layout = buildTopologyLayout(
      providers,
      new Set(["openai"]),
      new Set(["anthropic"]),
      new Set(["omeg4mon"]),
    );

    // 1 center router node + 4 provider nodes
    expect(layout.nodes).toHaveLength(5);
    expect(layout.edges).toHaveLength(4);

    const routerNode = layout.nodes.find((n) => n.id === "router");
    expect(routerNode?.data.activeCount).toBe(1);

    const openaiEdge = layout.edges.find((e) => e.target === "provider-openai");
    expect(openaiEdge?.data?.state).toBe("active");

    const anthropicEdge = layout.edges.find((e) => e.target === "provider-anthropic");
    expect(anthropicEdge?.data?.state).toBe("last");

    const customEdge = layout.edges.find((e) => e.target === "provider-omeg4mon");
    expect(customEdge?.data?.state).toBe("error");

    const geminiEdge = layout.edges.find((e) => e.target === "provider-gemini");
    expect(geminiEdge?.data?.state).toBe("idle");

    for (const node of layout.nodes) {
      expect(Number.isFinite(node.position.x)).toBe(true);
      expect(Number.isFinite(node.position.y)).toBe(true);
    }
  });

  test("buildTopologyLayout uses balanced horizontal wing layouts for 1 and 2 providers", () => {
    const p1 = {
      id: "omeg4mon",
      label: "Omeg4mon",
      protocol: "openai-chat",
      isCustom: true,
      connectionCount: 1,
      activeConnections: 1,
      requests: 8,
    };
    const p2 = {
      id: "boz",
      label: "Boz",
      protocol: "openai-chat",
      isCustom: true,
      connectionCount: 1,
      activeConnections: 1,
      requests: 3,
    };

    const singleLayout = buildTopologyLayout([p1], new Set(), new Set(), new Set());
    expect(singleLayout.edges[0]?.sourceHandle).toBe("right");
    expect(singleLayout.edges[0]?.targetHandle).toBe("left");

    const pairLayout = buildTopologyLayout([p1, p2], new Set(), new Set(), new Set());
    const handles = pairLayout.edges.map((e) => e.sourceHandle).sort();
    expect(handles).toEqual(["left", "right"]);
  });

  test("buildTopologyLayout keeps 3-provider layout visually centered and compact", () => {
    const trio = [
      {
        id: "omeg4mon",
        label: "Omeg4mon",
        protocol: "openai-chat",
        isCustom: true,
        connectionCount: 1,
        activeConnections: 1,
        requests: 9,
      },
      {
        id: "boz",
        label: "Boz",
        protocol: "openai-chat",
        isCustom: true,
        connectionCount: 1,
        activeConnections: 1,
        requests: 4,
      },
      {
        id: "routeva",
        label: "Routeva",
        protocol: "openai-chat",
        isCustom: true,
        connectionCount: 1,
        activeConnections: 1,
        requests: 2,
      },
    ];
    const layout = buildTopologyLayout(trio, new Set(), new Set(), new Set());
    const providerNodes = layout.nodes.filter((n) => n.type === "provider");
    expect(providerNodes).toHaveLength(3);
    const xs = providerNodes.map((n) => n.position.x);
    const spanX = Math.max(...xs) - Math.min(...xs);
    expect(spanX).toBeGreaterThan(260);
    expect(spanX).toBeLessThan(520);
  });
});
