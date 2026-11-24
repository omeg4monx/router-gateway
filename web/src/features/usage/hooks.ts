import { keepPreviousData, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  ModelPricing,
  UpsertModelPricingRequest,
  UsageLogsResponse,
  UsageRange,
  UsageSummaryResponse,
} from "./schemas";

export interface UsageLogsFilter {
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
  summary: (range: UsageRange) => ["usage", "summary", range] as const,
  logsRoot: ["usage", "logs"] as const,
  logs: (filters: UsageLogsFilter) => ["usage", "logs", filters] as const,
  pricing: ["usage", "pricing"] as const,
};

function buildLogsQueryParams(filters: UsageLogsFilter) {
  const query: {
    model?: string;
    provider?: string;
    key_id?: string;
    status?: string;
    limit?: number;
    offset?: number;
  } = {};

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

export function useUsageSummary(range: UsageRange) {
  return useQuery<UsageSummaryResponse, Error>({
    queryKey: usageQueryKeys.summary(range),
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/usage/summary", {
        params: { query: { range } },
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

export function useModelPricing() {
  return useQuery<ModelPricing[], Error>({
    queryKey: usageQueryKeys.pricing,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/usage/pricing");
      if (data) {
        return data.pricing;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export interface UpsertModelPricingVariables {
  model: string;
  body: UpsertModelPricingRequest;
}

export function useUpsertModelPricingMutation() {
  const queryClient = useQueryClient();

  return useMutation<ModelPricing, Error, UpsertModelPricingVariables>({
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
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: usageQueryKeys.pricing });
    },
  });
}

export function useDeleteModelPricingMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (model) => {
      const { error, response } = await apiClient.DELETE("/api/usage/pricing/{model}", {
        params: { path: { model } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to reset custom model pricing (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: usageQueryKeys.pricing });
    },
  });
}

export type UsageStreamStatus = "idle" | "connecting" | "live" | "error";

export function useUsageStream(enabled: boolean) {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<UsageStreamStatus>("idle");
  const [eventCount, setEventCount] = useState(0);

  useEffect(() => {
    if (!enabled || typeof EventSource === "undefined") {
      setStatus("idle");
      return;
    }

    setStatus("connecting");
    const source = new EventSource("/api/usage/stream", { withCredentials: true });
    let invalidateTimer: number | null = null;

    const scheduleInvalidation = () => {
      if (invalidateTimer !== null) {
        return;
      }
      invalidateTimer = window.setTimeout(() => {
        invalidateTimer = null;
        void queryClient.invalidateQueries({ queryKey: usageQueryKeys.summaries });
        void queryClient.invalidateQueries({ queryKey: usageQueryKeys.logsRoot });
      }, 300);
    };

    const handleOpen = () => {
      setStatus("live");
    };

    const handleUsageEvent = () => {
      setStatus("live");
      setEventCount((count) => count + 1);
      scheduleInvalidation();
    };

    const handleError = () => {
      setStatus("error");
    };

    source.addEventListener("open", handleOpen);
    source.addEventListener("usage", handleUsageEvent);
    source.addEventListener("error", handleError);

    return () => {
      if (invalidateTimer !== null) {
        window.clearTimeout(invalidateTimer);
      }
      source.removeEventListener("open", handleOpen);
      source.removeEventListener("usage", handleUsageEvent);
      source.removeEventListener("error", handleError);
      source.close();
      setStatus("idle");
    };
  }, [enabled, queryClient]);

  return { status, eventCount };
}
