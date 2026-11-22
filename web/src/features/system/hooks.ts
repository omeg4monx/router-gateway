"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  extractErrorMessage,
  type HealthResponse,
  type InitCompleteResponse,
  type InitRequest,
  type InitStatusResponse,
  type VersionResponse,
} from "@/shared/api/client";

export const systemKeys = {
  health: ["system", "health"] as const,
  version: ["system", "version"] as const,
  initStatus: ["system", "init"] as const,
};

function isHealthResponse(value: unknown): value is HealthResponse {
  if (typeof value !== "object" || value === null) {
    return false;
  }
  const candidate = value as Record<string, unknown>;
  return (
    (candidate.status === "ok" || candidate.status === "degraded") &&
    (candidate.db === "ok" || candidate.db === "error") &&
    Array.isArray(candidate.modules) &&
    typeof candidate.uptimeSeconds === "number"
  );
}

export function useSystemHealth() {
  return useQuery<HealthResponse, Error>({
    queryKey: systemKeys.health,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/health");
      if (data) {
        return data;
      }
      if (isHealthResponse(error)) {
        return error;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 30_000,
  });
}

export function useSystemVersion() {
  return useQuery<VersionResponse, Error>({
    queryKey: systemKeys.version,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/version");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    staleTime: 60_000,
  });
}

export function useInitStatus() {
  return useQuery<InitStatusResponse, Error>({
    queryKey: systemKeys.initStatus,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/init");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useInitializeSystem() {
  const queryClient = useQueryClient();

  return useMutation<InitCompleteResponse, Error, InitRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/init", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: systemKeys.initStatus });
    },
  });
}
