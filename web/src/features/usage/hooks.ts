import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  ListModelPricingResponse,
  ModelPricing,
  UpsertModelPricingRequest,
  UsageLog,
  UsageLogsResponse,
  UsageRange,
  UsageSummaryResponse,
} from "./schemas";
import { prependUsageLogToResponse } from "./streamCache";

export interface UsageLogsFilter {
  range?: UsageRange;
  q?: string;
  model?: string;
  provider?: string;
  key_id?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

export const usageQueryKeys = {
  all: ["usage"] as const,
  summaries: ["usage", "summary"] as const,
  summary: (range: UsageRange, includeSeries = true) =>
    includeSeries
      ? (["usage", "summary", range] as const)
      : (["usage", "summary", range, "no-series"] as const),
  logsRoot: ["usage", "logs"] as const,
  logs: (filters: UsageLogsFilter) => ["usage", "logs", filters] as const,
  pricing: ["usage", "pricing"] as const,
};

function buildLogsQueryParams(filters: UsageLogsFilter) {
  const query: {
    range?: UsageRange;
    q?: string;
    model?: string;
    provider?: string;
    key_id?: string;
    status?: string;
    limit?: number;
    offset?: number;
  } = {};

  if (filters.range) {
    query.range = filters.range;
  }
  if (filters.q?.trim()) {
    query.q = filters.q.trim();
  }
  if (filters.model?.trim()) {
    query.model = filters.model.trim();
  }
  if (filters.provider?.trim()) {
    query.provider = filters.provider.trim();
  }
  if (filters.key_id?.trim()) {
    query.key_id = filters.key_id.trim();
  }
  if (filters.status?.trim()) {
    query.status = filters.status.trim();
  }
  if (typeof filters.limit === "number") {
    query.limit = filters.limit;
  }
  if (typeof filters.offset === "number") {
    query.offset = filters.offset;
  }
  return query;
}

export function useUsageSummary(range: UsageRange, options?: { includeSeries?: boolean }) {
  const includeSeries = options?.includeSeries ?? true;
  return useQuery<UsageSummaryResponse, Error>({
    queryKey: usageQueryKeys.summary(range, includeSeries),
    placeholderData: keepPreviousData,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/usage/summary", {
        params: {
          query: includeSeries ? { range } : { range, include_series: false },
        },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUsageLogs(filters: UsageLogsFilter) {
  return useQuery<UsageLogsResponse, Error>({
    queryKey: usageQueryKeys.logs(filters),
    placeholderData: keepPreviousData,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/usage/logs", {
        params: { query: buildLogsQueryParams(filters) },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export type UsageStreamStatus = "idle" | "connecting" | "live" | "error";

const INVALIDATION_TOPIC_QUERY_KEYS: Record<string, readonly (readonly string[])[]> = {
  "providers.changed": [["providers"], ["models", "catalog"], ["arena", "models"]],
  "keys.changed": [["keys"]],
  "combos.changed": [["combos"], ["models", "catalog"]],
  "proxypools.changed": [["proxypools"]],
};

export function useUsageStream(enabled: boolean) {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<UsageStreamStatus>("idle");
  const [eventCount, setEventCount] = useState(0);
  const [lastEvent, setLastEvent] = useState<UsageLog | null>(null);

  useEffect(() => {
    if (!enabled || typeof EventSource === "undefined") {
      setStatus("idle");
      return;
    }

    setStatus("connecting");
    const source = new EventSource("/api/usage/stream", { withCredentials: true });
    let summaryTimer: number | null = null;
    let hadError = false;

    const scheduleSummaryInvalidation = () => {
      if (summaryTimer !== null) {
        return;
      }
      summaryTimer = window.setTimeout(() => {
        summaryTimer = null;
        void queryClient.invalidateQueries({ queryKey: usageQueryKeys.summaries });
      }, 350);
    };

    const handleOpen = () => {
      setStatus("live");
      if (hadError) {
        hadError = false;
        void queryClient.invalidateQueries({ queryKey: usageQueryKeys.summaries });
        void queryClient.invalidateQueries({ queryKey: usageQueryKeys.logsRoot });
      }
    };

    const handleUsageEvent = (evt: Event) => {
      setStatus("live");
      setEventCount((count) => count + 1);
      let patchedCache = false;
      if ("data" in evt && typeof (evt as MessageEvent).data === "string") {
        try {
          const parsed = JSON.parse((evt as MessageEvent).data) as Partial<UsageLog>;
          if (parsed && typeof parsed.provider_id === "string") {
            const logEntry: UsageLog = {
              id: typeof parsed.id === "number" ? parsed.id : Date.now(),
              request_id: typeof parsed.request_id === "string" ? parsed.request_id : "",
              key_id: typeof parsed.key_id === "string" ? parsed.key_id : "",
              model: typeof parsed.model === "string" ? parsed.model : "",
              provider_id: parsed.provider_id,
              connection_id: typeof parsed.connection_id === "string" ? parsed.connection_id : "",
              protocol: typeof parsed.protocol === "string" ? parsed.protocol : "",
              input_tokens: typeof parsed.input_tokens === "number" ? parsed.input_tokens : 0,
              output_tokens: typeof parsed.output_tokens === "number" ? parsed.output_tokens : 0,
              cached_tokens: typeof parsed.cached_tokens === "number" ? parsed.cached_tokens : 0,
              cost_usd: typeof parsed.cost_usd === "number" ? parsed.cost_usd : 0,
              latency_ms: typeof parsed.latency_ms === "number" ? parsed.latency_ms : 0,
              ttft_ms: typeof parsed.ttft_ms === "number" ? parsed.ttft_ms : 0,
              status_code: typeof parsed.status_code === "number" ? parsed.status_code : 200,
              error_summary: typeof parsed.error_summary === "string" ? parsed.error_summary : "",
              stream: Boolean(parsed.stream),
              created_at:
                typeof parsed.created_at === "string"
                  ? parsed.created_at
                  : new Date().toISOString(),
            };
            setLastEvent(logEntry);
            // Optimistically prepend into matching usage-log queries for 0ms UI updates.
            const cachedQueries = queryClient.getQueriesData<UsageLogsResponse>({
              queryKey: usageQueryKeys.logsRoot,
            });
            for (const [queryKey, data] of cachedQueries) {
              if (!data) continue;
              const filters =
                Array.isArray(queryKey) && queryKey.length >= 3 && typeof queryKey[2] === "object"
                  ? (queryKey[2] as UsageLogsFilter)
                  : undefined;
              const next = prependUsageLogToResponse(data, logEntry, filters);
              if (next !== data) {
                queryClient.setQueryData(queryKey, next);
              }
            }
            patchedCache = true;
          }
        } catch {
          // Fallback to full invalidation if payload is malformed.
        }
      }
      if (!patchedCache) {
        void queryClient.invalidateQueries({ queryKey: usageQueryKeys.logsRoot });
      }
      scheduleSummaryInvalidation();
    };

    const handleInvalidationEvent = (evt: Event) => {
      if (!("data" in evt) || typeof (evt as MessageEvent).data !== "string") {
        return;
      }
      try {
        const parsed = JSON.parse((evt as MessageEvent).data) as { topic?: string };
        const keys = parsed.topic ? INVALIDATION_TOPIC_QUERY_KEYS[parsed.topic] : undefined;
        if (keys) {
          for (const key of keys) {
            void queryClient.invalidateQueries({ queryKey: key });
          }
        }
      } catch {
        // Ignore malformed invalidation frame.
      }
    };

    const handleError = () => {
      hadError = true;
      setStatus("error");
    };

    source.addEventListener("open", handleOpen);
    source.addEventListener("usage", handleUsageEvent);
    source.addEventListener("invalidation", handleInvalidationEvent);
    source.addEventListener("error", handleError);

    return () => {
      if (summaryTimer !== null) {
        window.clearTimeout(summaryTimer);
      }
      source.removeEventListener("open", handleOpen);
      source.removeEventListener("usage", handleUsageEvent);
      source.removeEventListener("invalidation", handleInvalidationEvent);
      source.removeEventListener("error", handleError);
      source.close();
      setStatus("idle");
    };
  }, [enabled, queryClient]);

  return { status, eventCount, lastEvent };
}

export function useModelPricing(enabled: boolean) {
  return useQuery<ListModelPricingResponse, Error>({
    queryKey: usageQueryKeys.pricing,
    enabled,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/usage/pricing");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUpsertModelPricingMutation() {
  const queryClient = useQueryClient();
  return useMutation<ModelPricing, Error, { model: string; body: UpsertModelPricingRequest }>({
    mutationFn: async ({ model, body }) => {
      const { data, error } = await apiClient.PUT("/api/usage/pricing/{model}", {
        params: { path: { model } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: usageQueryKeys.pricing });
    },
  });
}

export function useDeleteModelPricingMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (model) => {
      const { error } = await apiClient.DELETE("/api/usage/pricing/{model}", {
        params: { path: { model } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: usageQueryKeys.pricing });
    },
  });
}
