"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  StartTunnelRequest,
  TunnelEventListResponse,
  TunnelStatus,
  UpdateTunnelConfigRequest,
} from "./schemas";

export const tunnelQueryKeys = {
  all: ["tunnel"] as const,
  status: ["tunnel", "status"] as const,
  events: ["tunnel", "events"] as const,
};

export function useTunnelStatus() {
  return useQuery<TunnelStatus, Error>({
    queryKey: tunnelQueryKeys.status,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/tunnel/status");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function useTunnelEvents() {
  return useQuery<TunnelEventListResponse, Error>({
    queryKey: tunnelQueryKeys.events,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/tunnel/events");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function useUpdateTunnelConfigMutation() {
  const queryClient = useQueryClient();
  return useMutation<TunnelStatus, Error, UpdateTunnelConfigRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/tunnel/config", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tunnelQueryKeys.all });
    },
  });
}

export function useStartTunnelMutation() {
  const queryClient = useQueryClient();
  return useMutation<TunnelStatus, Error, StartTunnelRequest | undefined>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/tunnel/start", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tunnelQueryKeys.all });
    },
  });
}

export function useStopTunnelMutation() {
  const queryClient = useQueryClient();
  return useMutation<TunnelStatus, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/tunnel/stop");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tunnelQueryKeys.all });
    },
  });
}
