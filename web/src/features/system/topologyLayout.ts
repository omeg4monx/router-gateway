import type { Edge, Node } from "@xyflow/react";

export type TopologyEdgeState = "active" | "last" | "error" | "idle";

export interface TopologyProviderInput {
  id: string;
  label: string;
  protocol: string;
  isCustom: boolean;
  connectionCount: number;
  activeConnections: number;
  requests: number;
}

export interface RouterNodeData extends Record<string, unknown> {
  activeCount: number;
  providerCount: number;
}

export interface ProviderNodeData extends Record<string, unknown> {
  providerId: string;
  label: string;
  protocol: string;
  isCustom: boolean;
  connectionCount: number;
  activeConnections: number;
  requests: number;
  state: TopologyEdgeState;
}

export interface TopologyEdgeData extends Record<string, unknown> {
  state: TopologyEdgeState;
}

export type TopologyNode = Node<RouterNodeData, "router"> | Node<ProviderNodeData, "provider">;

export type TopologyEdgeItem = Edge<TopologyEdgeData, "topology">;

interface SlotPlacement {
  cx: number;
  cy: number;
  sourceHandle: "top" | "bottom" | "left" | "right";
  targetHandle: "top" | "bottom" | "left" | "right";
}

export function groupTopologyProviders(
  connections: ReadonlyArray<{
    provider_id: string;
    protocol: string;
    enabled: boolean;
  }>,
  usageByProvider: ReadonlyArray<{
    provider_id: string;
    requests: number;
  }>,
  displayNameMap: ReadonlyMap<string, string> = new Map(),
  customNodeIds: ReadonlySet<string> = new Set(),
): TopologyProviderInput[] {
  const map = new Map<string, TopologyProviderInput>();

  for (const conn of connections) {
    const id = conn.provider_id.trim();
    if (!id) continue;
    const existing = map.get(id);
    if (existing) {
      existing.connectionCount += 1;
      if (conn.enabled) existing.activeConnections += 1;
    } else {
      map.set(id, {
        id,
        label: displayNameMap.get(id) ?? id,
        protocol: conn.protocol,
        isCustom: customNodeIds.has(id),
        connectionCount: 1,
        activeConnections: conn.enabled ? 1 : 0,
        requests: 0,
      });
    }
  }

  for (const item of usageByProvider) {
    const id = item.provider_id.trim();
    if (!id) continue;
    const existing = map.get(id);
    if (existing) {
      existing.requests += item.requests;
    } else {
      map.set(id, {
        id,
        label: displayNameMap.get(id) ?? id,
        protocol: "openai-chat",
        isCustom: customNodeIds.has(id),
        connectionCount: 0,
        activeConnections: 0,
        requests: item.requests,
      });
    }
  }

  return Array.from(map.values()).sort((a, b) => {
    if (b.requests !== a.requests) return b.requests - a.requests;
    return a.label.localeCompare(b.label);
  });
}

function computeProviderSlots(
  count: number,
  nodeW: number,
): {
  routerCenter: { x: number; y: number };
  slots: SlotPlacement[];
} {
  if (count === 1) {
    return {
      routerCenter: { x: -145, y: 0 },
      slots: [{ cx: 145, cy: 0, sourceHandle: "right", targetHandle: "left" }],
    };
  }

  if (count === 2) {
    return {
      routerCenter: { x: 0, y: 0 },
      slots: [
        { cx: -235, cy: 0, sourceHandle: "left", targetHandle: "right" },
        { cx: 235, cy: 0, sourceHandle: "right", targetHandle: "left" },
      ],
    };
  }

  if (count === 3) {
    return {
      routerCenter: { x: 0, y: 10 },
      slots: [
        { cx: 0, cy: -120, sourceHandle: "top", targetHandle: "bottom" },
        { cx: 215, cy: 96, sourceHandle: "right", targetHandle: "left" },
        { cx: -215, cy: 96, sourceHandle: "left", targetHandle: "right" },
      ],
    };
  }

  if (count === 4) {
    return {
      routerCenter: { x: 0, y: 0 },
      slots: [
        { cx: -230, cy: -94, sourceHandle: "left", targetHandle: "right" },
        { cx: 230, cy: -94, sourceHandle: "right", targetHandle: "left" },
        { cx: 230, cy: 94, sourceHandle: "right", targetHandle: "left" },
        { cx: -230, cy: 94, sourceHandle: "left", targetHandle: "right" },
      ],
    };
  }

  const nodeGap = 24;
  const minRx = ((nodeW + nodeGap) * count) / (2 * Math.PI);
  const rx = Math.max(255, minRx);
  const ry = Math.max(150, rx * 0.56);
  const slots: SlotPlacement[] = [];

  for (let i = 0; i < count; i += 1) {
    const angle = -Math.PI / 2 + (2 * Math.PI * i) / count;
    const cx = rx * Math.cos(angle);
    const cy = ry * Math.sin(angle);

    let sourceHandle: SlotPlacement["sourceHandle"];
    let targetHandle: SlotPlacement["targetHandle"];
    if (
      Math.abs(angle + Math.PI / 2) < Math.PI / 4 ||
      Math.abs(angle - (3 * Math.PI) / 2) < Math.PI / 4
    ) {
      sourceHandle = "top";
      targetHandle = "bottom";
    } else if (Math.abs(angle - Math.PI / 2) < Math.PI / 4) {
      sourceHandle = "bottom";
      targetHandle = "top";
    } else if (cx > 0) {
      sourceHandle = "right";
      targetHandle = "left";
    } else {
      sourceHandle = "left";
      targetHandle = "right";
    }

    slots.push({ cx, cy, sourceHandle, targetHandle });
  }

  return { routerCenter: { x: 0, y: 0 }, slots };
}

export function buildTopologyLayout(
  providers: readonly TopologyProviderInput[],
  activeSet: ReadonlySet<string>,
  lastSet: ReadonlySet<string>,
  errorSet: ReadonlySet<string>,
): { nodes: TopologyNode[]; edges: TopologyEdgeItem[] } {
  const nodeW = 196;
  const nodeH = 60;
  const routerW = 186;
  const routerH = 60;

  const count = providers.length;
  if (count === 0) {
    return {
      nodes: [
        {
          id: "router",
          type: "router",
          position: { x: -routerW / 2, y: -routerH / 2 },
          data: { activeCount: activeSet.size, providerCount: 0 },
          draggable: false,
        },
      ],
      edges: [],
    };
  }

  const { routerCenter, slots } = computeProviderSlots(count, nodeW);
  const routerNode: Node<RouterNodeData, "router"> = {
    id: "router",
    type: "router",
    position: {
      x: Number((routerCenter.x - routerW / 2).toFixed(1)),
      y: Number((routerCenter.y - routerH / 2).toFixed(1)),
    },
    data: { activeCount: activeSet.size, providerCount: count },
    draggable: false,
  };

  const nodes: TopologyNode[] = [routerNode];
  const edges: TopologyEdgeItem[] = [];

  for (let i = 0; i < count; i += 1) {
    const provider = providers[i];
    const slot = slots[i];
    if (!provider || !slot) continue;

    const key = provider.id.toLowerCase();
    let state: TopologyEdgeState = "idle";
    if (activeSet.has(key)) {
      state = "active";
    } else if (errorSet.has(key)) {
      state = "error";
    } else if (lastSet.has(key)) {
      state = "last";
    }

    const nodeId = `provider-${provider.id}`;
    nodes.push({
      id: nodeId,
      type: "provider",
      position: {
        x: Number((slot.cx - nodeW / 2).toFixed(1)),
        y: Number((slot.cy - nodeH / 2).toFixed(1)),
      },
      data: {
        providerId: provider.id,
        label: provider.label,
        protocol: provider.protocol,
        isCustom: provider.isCustom,
        connectionCount: provider.connectionCount,
        activeConnections: provider.activeConnections,
        requests: provider.requests,
        state,
      },
      draggable: false,
    });

    edges.push({
      id: `e-${nodeId}`,
      type: "topology",
      source: "router",
      sourceHandle: slot.sourceHandle,
      target: nodeId,
      targetHandle: slot.targetHandle,
      animated: false,
      data: { state },
    });
  }

  return { nodes, edges };
}
