"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreateMitmAliasRequest,
  GenerateMitmCARequest,
  MitmAlias,
  MitmAliasListResponse,
  MitmCAInfo,
  MitmEventListResponse,
  MitmStatus,
  UpdateMitmConfigRequest,
} from "./schemas";

export const mitmQueryKeys = {
  all: ["mitm"] as const,
  status: ["mitm", "status"] as const,
  aliases: ["mitm", "aliases"] as const,
  events: ["mitm", "events"] as const,
};

export function useMitmStatus() {
  return useQuery<MitmStatus, Error>({
    queryKey: mitmQueryKeys.status,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/mitm/status");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function useMitmAliases() {
  return useQuery<MitmAliasListResponse, Error>({
    queryKey: mitmQueryKeys.aliases,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/mitm/aliases");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useMitmEvents() {
  return useQuery<MitmEventListResponse, Error>({
    queryKey: mitmQueryKeys.events,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/mitm/events");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 5000,
  });
}

export function useUpdateMitmConfigMutation() {
  const queryClient = useQueryClient();
  return useMutation<MitmStatus, Error, UpdateMitmConfigRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/mitm/config", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.all });
    },
  });
}

export function useStartMitmMutation() {
  const queryClient = useQueryClient();
  return useMutation<MitmStatus, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/mitm/start");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.all });
    },
  });
}

export function useStopMitmMutation() {
  const queryClient = useQueryClient();
  return useMutation<MitmStatus, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/mitm/stop");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.all });
    },
  });
}

export function useGenerateMitmCAMutation() {
  const queryClient = useQueryClient();
  return useMutation<MitmCAInfo, Error, GenerateMitmCARequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/mitm/ca/generate", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.status });
    },
  });
}

export function useCreateMitmAliasMutation() {
  const queryClient = useQueryClient();
  return useMutation<MitmAlias, Error, CreateMitmAliasRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/mitm/aliases", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.aliases });
    },
  });
}

export function useDeleteMitmAliasMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/mitm/aliases/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.aliases });
    },
  });
}

export function useClearMitmEventsMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: async () => {
      const { error } = await apiClient.DELETE("/api/mitm/events");
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mitmQueryKeys.events });
    },
  });
}
