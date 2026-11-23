"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  PxpipeHealthReport,
  PxpipeLogListResponse,
  PxpipePreviewRequest,
  PxpipePreviewResponse,
  PxpipeStats,
  PxpipeStatus,
  UpdatePxpipeConfigRequest,
} from "./schemas";

export const pxpipeQueryKeys = {
  all: ["pxpipe"] as const,
  status: ["pxpipe", "status"] as const,
  health: ["pxpipe", "health"] as const,
  stats: ["pxpipe", "stats"] as const,
  logs: ["pxpipe", "logs"] as const,
};

export function usePxpipeStatus() {
  return useQuery<PxpipeStatus, Error>({
    queryKey: pxpipeQueryKeys.status,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/pxpipe/status");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function usePxpipeHealth() {
  return useQuery<PxpipeHealthReport, Error>({
    queryKey: pxpipeQueryKeys.health,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/pxpipe/health");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 10000,
  });
}

export function usePxpipeStats() {
  return useQuery<PxpipeStats, Error>({
    queryKey: pxpipeQueryKeys.stats,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/pxpipe/stats");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function usePxpipeLogs() {
  return useQuery<PxpipeLogListResponse, Error>({
    queryKey: pxpipeQueryKeys.logs,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/pxpipe/logs");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function useUpdatePxpipeConfigMutation() {
  const queryClient = useQueryClient();
  return useMutation<PxpipeStatus, Error, UpdatePxpipeConfigRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/pxpipe/config", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: pxpipeQueryKeys.all });
    },
  });
}

export function useStartPxpipeMutation() {
  const queryClient = useQueryClient();
  return useMutation<PxpipeStatus, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/pxpipe/start");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: pxpipeQueryKeys.all });
    },
  });
}

export function useStopPxpipeMutation() {
  const queryClient = useQueryClient();
  return useMutation<PxpipeStatus, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/pxpipe/stop");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: pxpipeQueryKeys.all });
    },
  });
}

export function useRestartPxpipeMutation() {
  const queryClient = useQueryClient();
  return useMutation<PxpipeStatus, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/pxpipe/restart");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: pxpipeQueryKeys.all });
    },
  });
}

export function useResetPxpipeStatsMutation() {
  const queryClient = useQueryClient();
  return useMutation<PxpipeStats, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/pxpipe/stats/reset");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: pxpipeQueryKeys.stats });
    },
  });
}

export function useClearPxpipeLogsMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: async () => {
      const { error } = await apiClient.DELETE("/api/pxpipe/logs");
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: pxpipeQueryKeys.logs });
    },
  });
}

export function usePreviewPxpipeMutation() {
  return useMutation<PxpipePreviewResponse, Error, PxpipePreviewRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/pxpipe/preview", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}
