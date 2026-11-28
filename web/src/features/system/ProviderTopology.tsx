import {
  Background,
  BackgroundVariant,
  BaseEdge,
  Controls,
  type EdgeProps,
  getBezierPath,
  Handle,
  type NodeProps,
  Position,
  ReactFlow,
  type ReactFlowInstance,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { useCallback, useEffect, useMemo, useRef } from "react";
import { ProviderLogo } from "@/features/providers/ProviderLogo";
import { cn } from "@/shared/lib/utils";
import {
  buildTopologyLayout,
  type ProviderNodeData,
  type RouterNodeData,
  type TopologyEdgeData,
  type TopologyEdgeItem,
  type TopologyEdgeState,
  type TopologyNode,
  type TopologyProviderInput,
} from "./topologyLayout";

function prefersReducedMotion(): boolean {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return false;
  }
  return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

const FLOW_PULSES = [
  { id: "pulse-0", r: 3.2, dur: "1.45s", begin: "0s", fill: "var(--color-highlight)" },
  { id: "pulse-1", r: 2.4, dur: "1.45s", begin: "0.48s", fill: "var(--status-ok)" },
  { id: "pulse-2", r: 2.0, dur: "1.45s", begin: "0.96s", fill: "var(--color-text-main)" },
] as const;

const FIT_OPTIONS = { padding: 0.18, duration: 280, maxZoom: 1.12 };

function RouterHubNode({ data }: NodeProps<TopologyNode & { data: RouterNodeData }>) {
  const powering = data.activeCount > 0;

  return (
    <div
      className={cn(
        "relative z-[1] flex h-[60px] w-[186px] items-center gap-2.5 rounded-2xl border px-3.5 transition-all duration-500 ease-out",
        powering
          ? "topology-router-core border-[var(--color-highlight)]/85 bg-[var(--color-surface)]"
          : "border-[var(--color-primary)]/70 bg-[var(--color-surface)] shadow-[0_8px_24px_rgba(0,0,0,0.08)]",
      )}
    >
      <Handle
        type="source"
        position={Position.Top}
        id="top"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />
      <Handle
        type="source"
        position={Position.Bottom}
        id="bottom"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />
      <Handle
        type="source"
        position={Position.Left}
        id="left"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />
      <Handle
        type="source"
        position={Position.Right}
        id="right"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />

      <div
        className={cn(
          "flex h-9 w-9 shrink-0 items-center justify-center rounded-xl border transition-colors duration-500",
          powering
            ? "border-[var(--color-highlight)]/45 bg-[var(--color-highlight)]/12"
            : "border-[var(--color-primary)]/25 bg-[var(--color-primary)]/10",
        )}
      >
        <img
          src="/icon.svg"
          alt="OMEG4MON ROUTER"
          className={cn(
            "h-5 w-5 shrink-0 transition-transform duration-500",
            powering && "scale-105",
          )}
        />
      </div>

      <div className="min-w-0 flex-1">
        <div
          className={cn(
            "truncate font-mono text-[12px] font-bold tracking-wider transition-colors duration-300",
            powering ? "text-[var(--color-highlight)]" : "text-[var(--color-text-main)]",
          )}
        >
          OMEG4MON
        </div>
        <div className="mt-0.5 truncate font-mono text-[10px] tracking-wide text-[var(--color-text-muted)]">
          Gateway · {data.providerCount}
        </div>
      </div>

      {powering ? (
        <span className="topology-router-badge inline-flex h-5 min-w-5 shrink-0 items-center justify-center rounded-full bg-[var(--color-highlight)] px-1.5 font-mono text-[10px] font-bold text-[var(--color-on-primary)]">
          {data.activeCount}
        </span>
      ) : (
        <span
          className="h-2 w-2 shrink-0 rounded-full bg-[var(--status-ok)]/85"
          aria-hidden="true"
        />
      )}
    </div>
  );
}

function ProviderLeafNode({ data }: NodeProps<TopologyNode & { data: ProviderNodeData }>) {
  const isActive = data.state === "active";
  const isError = data.state === "error";
  const isLast = data.state === "last";

  return (
    <div
      className={cn(
        "flex min-h-[60px] w-[196px] items-center gap-2.5 rounded-xl border bg-[var(--color-surface)] px-3 py-2 transition-all duration-500 ease-out",
        isActive
          ? "border-[var(--color-highlight)]/85 shadow-[0_0_20px_color-mix(in_srgb,var(--color-highlight)_26%,transparent)]"
          : isError
            ? "border-[var(--status-error)]/80 shadow-[0_0_16px_color-mix(in_srgb,var(--status-error)_22%,transparent)]"
            : isLast
              ? "border-[var(--status-warn)]/75 shadow-[0_0_14px_color-mix(in_srgb,var(--status-warn)_18%,transparent)]"
              : "border-[var(--color-border)]/90 shadow-[0_2px_10px_rgba(0,0,0,0.04)] hover:border-[var(--color-primary)]/45",
      )}
    >
      <Handle
        type="target"
        position={Position.Top}
        id="top"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />
      <Handle
        type="target"
        position={Position.Bottom}
        id="bottom"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />
      <Handle
        type="target"
        position={Position.Left}
        id="left"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />
      <Handle
        type="target"
        position={Position.Right}
        id="right"
        className="!h-0 !w-0 !border-0 !bg-transparent"
      />

      <div
        className={cn(
          "flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border transition-colors duration-300",
          isActive
            ? "border-[var(--color-highlight)]/35 bg-[var(--color-highlight)]/10"
            : "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/75",
        )}
      >
        <ProviderLogo
          providerId={data.providerId}
          name={data.label}
          protocol={data.protocol}
          isCustom={data.isCustom}
          size={20}
        />
      </div>

      <div className="min-w-0 flex-1">
        <div className="flex items-center justify-between gap-1.5">
          <span
            className={cn(
              "truncate text-xs font-semibold transition-colors duration-300",
              isActive
                ? "text-[var(--color-highlight)]"
                : isError
                  ? "text-[var(--status-error)]"
                  : "text-[var(--color-text-main)]",
            )}
            title={data.label}
          >
            {data.label}
          </span>
          {isActive ? (
            <span className="relative flex h-2 w-2 shrink-0" aria-hidden="true">
              <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-[var(--color-highlight)] opacity-70 motion-reduce:animate-none" />
              <span className="relative inline-flex h-2 w-2 rounded-full bg-[var(--color-highlight)]" />
            </span>
          ) : isError ? (
            <span
              className="inline-flex h-2 w-2 shrink-0 rounded-full bg-[var(--status-error)]"
              aria-hidden="true"
            />
          ) : isLast ? (
            <span
              className="inline-flex h-2 w-2 shrink-0 rounded-full bg-[var(--status-warn)]"
              aria-hidden="true"
            />
          ) : null}
        </div>
        <div className="mt-0.5 flex items-center gap-1.5 font-mono text-[10px] tabular-nums text-[var(--color-text-muted)]">
          <span className="font-medium text-[var(--color-text-muted)]">
            {data.requests.toLocaleString()} req
          </span>
          {data.connectionCount > 0 ? (
            <>
              <span className="opacity-50">·</span>
              <span>
                {data.activeConnections}/{data.connectionCount} on
              </span>
            </>
          ) : null}
        </div>
      </div>
    </div>
  );
}

function EnergyTopologyEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
}: EdgeProps<TopologyEdgeItem & { data?: TopologyEdgeData }>) {
  const [edgePath] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
    curvature: 0.34,
  });

  const state: TopologyEdgeState = data?.state ?? "idle";

  if (state === "idle") {
    return (
      <BaseEdge
        id={id}
        path={edgePath}
        style={{ stroke: "var(--color-border)", strokeWidth: 1.25, opacity: 0.52 }}
      />
    );
  }

  if (state === "last") {
    return (
      <g>
        <BaseEdge
          id={id}
          path={edgePath}
          style={{ stroke: "var(--status-warn)", strokeWidth: 1.6, opacity: 0.32 }}
        />
        <path
          d={edgePath}
          fill="none"
          stroke="var(--status-warn)"
          strokeWidth={1.8}
          strokeOpacity={0.78}
          strokeLinecap="round"
          className="topology-edge-last"
        />
      </g>
    );
  }

  if (state === "error") {
    return (
      <g>
        <path
          d={edgePath}
          fill="none"
          stroke="var(--status-error)"
          strokeWidth={6}
          strokeOpacity={0.14}
          strokeLinecap="round"
        />
        <BaseEdge
          id={id}
          path={edgePath}
          style={{ stroke: "var(--status-error)", strokeWidth: 2, opacity: 0.88 }}
        />
      </g>
    );
  }

  const filterId = `topology-glow-${id.replace(/[^a-zA-Z0-9_-]/g, "")}`;

  return (
    <g>
      <defs>
        <filter id={filterId} x="-20%" y="-20%" width="140%" height="140%">
          <feGaussianBlur stdDeviation="2.6" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      {/* Soft ambient fiber bloom */}
      <path
        d={edgePath}
        fill="none"
        stroke="var(--color-highlight)"
        strokeWidth={7}
        strokeOpacity={0.2}
        strokeLinecap="round"
        filter={`url(#${filterId})`}
        className="topology-edge-bloom"
      />
      {/* Base continuous fiber */}
      <BaseEdge
        id={id}
        path={edgePath}
        style={{ stroke: "var(--color-highlight)", strokeWidth: 2, opacity: 0.55 }}
      />
      {/* Silky flowing stream dash */}
      <path
        d={edgePath}
        fill="none"
        stroke="var(--status-ok)"
        strokeWidth={2.2}
        strokeOpacity={0.9}
        strokeLinecap="round"
        className="topology-edge-flow"
      />

      {/* Smooth comet light pulses fading in & out cleanly along path */}
      {!prefersReducedMotion() &&
        FLOW_PULSES.map((pulse) => (
          <circle
            key={`${id}-${pulse.id}`}
            r={pulse.r}
            fill={pulse.fill}
            filter={`url(#${filterId})`}
            opacity={0}
          >
            <animate
              attributeName="opacity"
              values="0;0.95;0.95;0"
              keyTimes="0;0.14;0.84;1"
              dur={pulse.dur}
              begin={pulse.begin}
              repeatCount="indefinite"
            />
            <animateMotion
              dur={pulse.dur}
              begin={pulse.begin}
              repeatCount="indefinite"
              calcMode="spline"
              keyTimes="0;1"
              keySplines="0.4 0 0.2 1"
              path={edgePath}
            />
          </circle>
        ))}
    </g>
  );
}

const NODE_TYPES = {
  router: RouterHubNode,
  provider: ProviderLeafNode,
};

const EDGE_TYPES = {
  topology: EnergyTopologyEdge,
};

export interface ProviderTopologyProps {
  providers: readonly TopologyProviderInput[];
  activeProviders: ReadonlySet<string>;
  lastProvider: string;
  errorProvider: string;
  isId: boolean;
}

export default function ProviderTopology({
  providers,
  activeProviders,
  lastProvider,
  errorProvider,
  isId,
}: ProviderTopologyProps) {
  const activeKey = useMemo(() => Array.from(activeProviders).sort().join(","), [activeProviders]);

  const activeSet = useMemo(() => new Set(activeKey ? activeKey.split(",") : []), [activeKey]);
  const lastSet = useMemo(
    () => new Set(lastProvider ? [lastProvider.toLowerCase()] : []),
    [lastProvider],
  );
  const errorSet = useMemo(
    () => new Set(errorProvider ? [errorProvider.toLowerCase()] : []),
    [errorProvider],
  );

  const { nodes, edges } = useMemo(
    () => buildTopologyLayout(providers, activeSet, lastSet, errorSet),
    [providers, activeSet, lastSet, errorSet],
  );

  const providersKey = useMemo(
    () =>
      providers
        .map((p) => p.id)
        .sort()
        .join(","),
    [providers],
  );

  const rfInstanceRef = useRef<ReactFlowInstance<TopologyNode, TopologyEdgeItem> | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const fitTimerRef = useRef<number | null>(null);

  const scheduleFitView = useCallback(() => {
    if (fitTimerRef.current !== null) {
      window.clearTimeout(fitTimerRef.current);
    }
    fitTimerRef.current = window.setTimeout(() => {
      fitTimerRef.current = null;
      void rfInstanceRef.current?.fitView(FIT_OPTIONS);
    }, 60);
  }, []);

  const handleInit = useCallback(
    (instance: ReactFlowInstance<TopologyNode, TopologyEdgeItem>) => {
      rfInstanceRef.current = instance;
      scheduleFitView();
    },
    [scheduleFitView],
  );

  useEffect(() => {
    const el = containerRef.current;
    if (!el || typeof ResizeObserver === "undefined") {
      return;
    }
    const ro = new ResizeObserver(() => {
      scheduleFitView();
    });
    ro.observe(el);
    return () => {
      ro.disconnect();
    };
  }, [scheduleFitView]);

  useEffect(() => {
    if (nodes.length > 0) {
      scheduleFitView();
    }
    return () => {
      if (fitTimerRef.current !== null) {
        window.clearTimeout(fitTimerRef.current);
        fitTimerRef.current = null;
      }
    };
  }, [nodes.length, scheduleFitView]);

  if (providers.length === 0) {
    return (
      <div className="flex h-[360px] w-full flex-col items-center justify-center gap-2 px-6 text-center text-xs text-[var(--color-text-muted)] sm:h-[410px]">
        <span>
          {isId
            ? "Belum ada provider terhubung untuk ditampilkan di topologi."
            : "No providers connected yet to visualize in the topology map."}
        </span>
      </div>
    );
  }

  return (
    <div
      ref={containerRef}
      className="relative h-[360px] w-full min-w-0 sm:h-[410px]"
      style={{
        background:
          "radial-gradient(ellipse 68% 60% at 50% 50%, color-mix(in srgb, var(--color-primary) 5%, transparent), transparent 72%)",
      }}
    >
      <ReactFlow
        key={providersKey}
        nodes={nodes}
        edges={edges}
        nodeTypes={NODE_TYPES}
        edgeTypes={EDGE_TYPES}
        fitView
        fitViewOptions={FIT_OPTIONS}
        minZoom={0.35}
        maxZoom={1.6}
        onInit={handleInit}
        proOptions={{ hideAttribution: true }}
        panOnDrag
        zoomOnScroll
        zoomOnPinch
        zoomOnDoubleClick
        preventScrolling={false}
        nodesDraggable={false}
        nodesConnectable={false}
        elementsSelectable={false}
      >
        <Background
          variant={BackgroundVariant.Dots}
          gap={22}
          size={1}
          color="var(--color-border)"
        />
        <Controls showInteractive={false} className="topology-flow-controls" />
      </ReactFlow>
    </div>
  );
}
