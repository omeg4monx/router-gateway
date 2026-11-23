"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreateProxyPoolRequest,
  ListProxyPoolsResponse,
  ProxyBindingSummary,
  ProxyPoolSummary,
  ScanProxyCandidatesRequest,
  ScanProxyCandidatesResponse,
  TestProxyPoolRequest,
  TestProxyPoolResponse,
  UpdateProxyPoolRequest,
} from "./schemas";

export const proxyPoolsQueryKeys = {
  all: ["proxypools"] as const,
  list: ["proxypools", "list"] as const,
};

export function useProxyPools() {
  return useQuery<ListProxyPoolsResponse, Error>({
    queryKey: proxyPoolsQueryKeys.list,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/proxypools");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateProxyPoolMutation() {
  const queryClient = useQueryClient();
  return useMutation<ProxyPoolSummary, Error, CreateProxyPoolRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/proxypools", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}

export function useUpdateProxyPoolMutation() {
  const queryClient = useQueryClient();
  return useMutation<ProxyPoolSummary, Error, { id: string; body: UpdateProxyPoolRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/proxypools/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}

export function useDeleteProxyPoolMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/proxypools/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}

export function useTestProxyPoolMutation() {
  const queryClient = useQueryClient();
  return useMutation<TestProxyPoolResponse, Error, { id: string; body?: TestProxyPoolRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.POST("/api/proxypools/{id}/test", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}

export function useScanProxyCandidatesMutation() {
  const queryClient = useQueryClient();
  return useMutation<ScanProxyCandidatesResponse, Error, ScanProxyCandidatesRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/proxypools/scan", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}

export function useUpsertProxyBindingMutation() {
  const queryClient = useQueryClient();
  return useMutation<ProxyBindingSummary, Error, { providerId: string; poolId: string }>({
    mutationFn: async ({ providerId, poolId }) => {
      const { data, error } = await apiClient.PUT("/api/proxypools/bindings/{providerId}", {
        params: { path: { providerId } },
        body: { pool_id: poolId },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}

export function useDeleteProxyBindingMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (providerId) => {
      const { error } = await apiClient.DELETE("/api/proxypools/bindings/{providerId}", {
        params: { path: { providerId } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: proxyPoolsQueryKeys.all });
    },
  });
}
