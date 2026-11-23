"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreateMediaComboRequest,
  CreateMediaProviderRequest,
  MediaCombo,
  MediaComboListResponse,
  MediaProvider,
  MediaProviderListResponse,
  UpdateMediaComboRequest,
  UpdateMediaProviderRequest,
} from "./schemas";

export const mediaQueryKeys = {
  all: ["media"] as const,
  providers: ["media", "providers"] as const,
  combos: ["media", "combos"] as const,
};

export function useMediaProviders() {
  return useQuery<MediaProviderListResponse, Error>({
    queryKey: mediaQueryKeys.providers,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/media/providers");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateMediaProviderMutation() {
  const queryClient = useQueryClient();
  return useMutation<MediaProvider, Error, CreateMediaProviderRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/media/providers", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mediaQueryKeys.providers });
    },
  });
}

export function useUpdateMediaProviderMutation() {
  const queryClient = useQueryClient();
  return useMutation<MediaProvider, Error, { id: string; body: UpdateMediaProviderRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/media/providers/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mediaQueryKeys.providers });
    },
  });
}

export function useDeleteMediaProviderMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ deleted: boolean }, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.DELETE("/api/media/providers/{id}", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mediaQueryKeys.providers });
    },
  });
}

export function useMediaCombos() {
  return useQuery<MediaComboListResponse, Error>({
    queryKey: mediaQueryKeys.combos,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/media/combos");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateMediaComboMutation() {
  const queryClient = useQueryClient();
  return useMutation<MediaCombo, Error, CreateMediaComboRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/media/combos", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mediaQueryKeys.combos });
    },
  });
}

export function useUpdateMediaComboMutation() {
  const queryClient = useQueryClient();
  return useMutation<MediaCombo, Error, { id: string; body: UpdateMediaComboRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/media/combos/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mediaQueryKeys.combos });
    },
  });
}

export function useDeleteMediaComboMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ deleted: boolean }, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.DELETE("/api/media/combos/{id}", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: mediaQueryKeys.combos });
    },
  });
}
