"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  TokenSaverConfig,
  TokenSaverConfigPatch,
  TokenSaverPreviewRequest,
  TokenSaverPreviewResult,
  TokenSaverStats,
} from "./schemas";

export const tokenSaverQueryKeys = {
  all: ["tokensaver"] as const,
  config: ["tokensaver", "config"] as const,
  stats: ["tokensaver", "stats"] as const,
};

export function useTokenSaverConfig() {
  return useQuery<TokenSaverConfig, Error>({
    queryKey: tokenSaverQueryKeys.config,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/token-saver/config");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUpdateTokenSaverConfigMutation() {
  const queryClient = useQueryClient();
  return useMutation<TokenSaverConfig, Error, TokenSaverConfigPatch>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/token-saver/config", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tokenSaverQueryKeys.all });
    },
  });
}

export function useTokenSaverStats() {
  return useQuery<TokenSaverStats, Error>({
    queryKey: tokenSaverQueryKeys.stats,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/token-saver/stats");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useResetTokenSaverStatsMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: async () => {
      const { error } = await apiClient.DELETE("/api/token-saver/stats");
      if (error) throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tokenSaverQueryKeys.stats });
    },
  });
}

export function useClearHeadroomCacheMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: async () => {
      const { error } = await apiClient.DELETE("/api/token-saver/cache");
      if (error) throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tokenSaverQueryKeys.stats });
    },
  });
}

export function usePreviewTokenSaverMutation() {
  return useMutation<TokenSaverPreviewResult, Error, TokenSaverPreviewRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/token-saver/preview", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}
