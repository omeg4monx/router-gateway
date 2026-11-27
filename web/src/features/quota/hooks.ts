import {
  keepPreviousData,
  type QueryClient,
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import {
  type ProviderConnectionQuota,
  type ProviderQuotasResponse,
  patchConnectionQuotaInResponse,
} from "./schemas";

const LIVE_POLL_INTERVAL_MS = 15_000;
const PER_CONNECTION_REFRESH_COOLDOWN_MS = 2_500;
const USAGE_EVENT_DEBOUNCE_MS = 450;

export type QuotaLiveStatus = "idle" | "connecting" | "live" | "reconnecting";

export const quotaQueryKeys = {
  all: ["providers", "quotas"] as const,
  list: (providerId?: string) => ["providers", "quotas", providerId ?? "all"] as const,
};

async function fetchProviderQuotas(options?: {
  providerId?: string;
  refresh?: boolean;
}): Promise<ProviderQuotasResponse> {
  const query: { provider_id?: string; refresh?: boolean } = {};
  if (options?.providerId && options.providerId !== "all") {
    query.provider_id = options.providerId;
  }
  if (options?.refresh) {
    query.refresh = true;
  }

  const { data, error } = await apiClient.GET("/api/providers/quotas", {
    params: { query },
  });
  if (data) {
    return data;
  }
  throw new Error(extractErrorMessage(error));
}

async function fetchSingleConnectionQuota(
  connectionId: string,
  signal?: AbortSignal,
): Promise<ProviderConnectionQuota> {
  const { data, error } = await apiClient.POST("/api/providers/connections/{id}/quota/refresh", {
    params: { path: { id: connectionId } },
    signal,
  });
  if (data) {
    return data;
  }
  throw new Error(extractErrorMessage(error));
}

function upsertConnectionQuotaInCache(queryClient: QueryClient, updated: ProviderConnectionQuota) {
  queryClient.setQueriesData<ProviderQuotasResponse>({ queryKey: quotaQueryKeys.all }, (current) =>
    current ? patchConnectionQuotaInResponse(current, updated) : current,
  );
}

export function useProviderQuotas(options?: { providerId?: string; liveEnabled?: boolean }) {
  const providerId = options?.providerId;
  const liveEnabled = options?.liveEnabled ?? true;

  return useQuery<ProviderQuotasResponse, Error>({
    queryKey: quotaQueryKeys.list(providerId),
    queryFn: () => fetchProviderQuotas({ providerId }),
    staleTime: 10_000,
    placeholderData: keepPreviousData,
    refetchInterval: liveEnabled ? LIVE_POLL_INTERVAL_MS : false,
    refetchIntervalInBackground: false,
    refetchOnWindowFocus: true,
  });
}

export function useRefreshAllQuotasMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderQuotasResponse, Error, { providerId?: string } | undefined>({
    mutationFn: (vars) =>
      fetchProviderQuotas({
        providerId: vars?.providerId,
        refresh: true,
      }),
    onSuccess: (fresh, vars) => {
      queryClient.setQueryData(quotaQueryKeys.list(vars?.providerId), fresh);
      if (!vars?.providerId || vars.providerId === "all") {
        queryClient.setQueryData(quotaQueryKeys.list("all"), fresh);
      }
    },
  });
}

export function useRefreshConnectionQuotaMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderConnectionQuota, Error, string>({
    mutationFn: (connectionId) => fetchSingleConnectionQuota(connectionId),
    onSuccess: (updated) => {
      upsertConnectionQuotaInCache(queryClient, updated);
    },
  });
}

export function useQuotaLiveSync(enabled: boolean, trackedConnectionIds: readonly string[]) {
  const queryClient = useQueryClient();
  const [liveStatus, setLiveStatus] = useState<QuotaLiveStatus>("idle");
  const [syncingIds, setSyncingIds] = useState<ReadonlySet<string>>(() => new Set());

  const trackedIdsRef = useRef<Set<string>>(new Set(trackedConnectionIds));
  trackedIdsRef.current = new Set(trackedConnectionIds);

  useEffect(() => {
    if (!enabled || typeof EventSource === "undefined") {
      setLiveStatus("idle");
      setSyncingIds(new Set());
      return;
    }

    let disposed = false;
    let hadError = false;
    const controllers = new Map<string, AbortController>();
    const pendingTimers = new Map<string, number>();
    const lastRefreshedAt = new Map<string, number>();

    setLiveStatus("connecting");
    const source = new EventSource("/api/usage/stream", { withCredentials: true });

    const setConnectionSyncing = (connectionId: string, active: boolean) => {
      if (disposed) {
        return;
      }
      setSyncingIds((prev) => {
        const has = prev.has(connectionId);
        if (has === active) {
          return prev;
        }
        const next = new Set(prev);
        if (active) {
          next.add(connectionId);
        } else {
          next.delete(connectionId);
        }
        return next;
      });
    };

    const runConnectionRefresh = (connectionId: string) => {
      if (disposed || !trackedIdsRef.current.has(connectionId)) {
        return;
      }
      const existingController = controllers.get(connectionId);
      if (existingController) {
        existingController.abort();
      }
      const controller = new AbortController();
      controllers.set(connectionId, controller);
      lastRefreshedAt.set(connectionId, Date.now());
      setConnectionSyncing(connectionId, true);

      void fetchSingleConnectionQuota(connectionId, controller.signal)
        .then((updated) => {
          if (!disposed) {
            upsertConnectionQuotaInCache(queryClient, updated);
          }
        })
        .catch(() => {
          // Best-effort live sync; periodic poll & manual refresh remain active.
        })
        .finally(() => {
          if (controllers.get(connectionId) === controller) {
            controllers.delete(connectionId);
          }
          setConnectionSyncing(connectionId, false);
        });
    };

    const scheduleConnectionRefresh = (connectionId: string) => {
      if (!connectionId || !trackedIdsRef.current.has(connectionId)) {
        return;
      }
      if (pendingTimers.has(connectionId)) {
        return;
      }
      const now = Date.now();
      const last = lastRefreshedAt.get(connectionId) ?? 0;
      const elapsed = now - last;
      const delayMs =
        elapsed >= PER_CONNECTION_REFRESH_COOLDOWN_MS
          ? USAGE_EVENT_DEBOUNCE_MS
          : PER_CONNECTION_REFRESH_COOLDOWN_MS - elapsed;

      const timer = window.setTimeout(() => {
        pendingTimers.delete(connectionId);
        runConnectionRefresh(connectionId);
      }, delayMs);
      pendingTimers.set(connectionId, timer);
    };

    const handleOpen = () => {
      if (disposed) {
        return;
      }
      setLiveStatus("live");
      if (hadError) {
        hadError = false;
        void queryClient.invalidateQueries({ queryKey: quotaQueryKeys.all });
      }
    };

    const handleUsageEvent = (evt: Event) => {
      if (disposed) {
        return;
      }
      setLiveStatus("live");
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      try {
        const parsed = JSON.parse((evt as MessageEvent).data) as { connection_id?: unknown };
        if (typeof parsed.connection_id === "string" && parsed.connection_id.trim() !== "") {
          scheduleConnectionRefresh(parsed.connection_id.trim());
        }
      } catch {
        // Ignore malformed SSE frames.
      }
    };

    const handleInvalidationEvent = (evt: Event) => {
      if (disposed || !("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      try {
        const parsed = JSON.parse((evt as MessageEvent).data) as { topic?: unknown };
        if (parsed.topic === "providers.changed") {
          void queryClient.invalidateQueries({ queryKey: quotaQueryKeys.all });
        }
      } catch {
        // Ignore malformed invalidation frame.
      }
    };

    const handleError = () => {
      if (disposed) {
        return;
      }
      hadError = true;
      setLiveStatus("reconnecting");
    };

    source.addEventListener("open", handleOpen);
    source.addEventListener("usage", handleUsageEvent);
    source.addEventListener("invalidation", handleInvalidationEvent);
    source.addEventListener("error", handleError);

    return () => {
      disposed = true;
      for (const timer of pendingTimers.values()) {
        window.clearTimeout(timer);
      }
      pendingTimers.clear();
      for (const controller of controllers.values()) {
        controller.abort();
      }
      controllers.clear();
      lastRefreshedAt.clear();
      source.removeEventListener("open", handleOpen);
      source.removeEventListener("usage", handleUsageEvent);
      source.removeEventListener("invalidation", handleInvalidationEvent);
      source.removeEventListener("error", handleError);
      source.close();
      setLiveStatus("idle");
      setSyncingIds(new Set());
    };
  }, [enabled, queryClient]);

  return { liveStatus, syncingIds };
}
