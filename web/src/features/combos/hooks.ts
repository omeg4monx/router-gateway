"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  Combo,
  CooldownEntry,
  CreateComboRequest,
  DisabledModel,
  DisableModelRequest,
  ModelAlias,
  UpdateComboRequest,
  UpsertModelAliasRequest,
} from "./schemas";

export const combosQueryKeys = {
  all: ["combos"] as const,
  list: ["combos", "list"] as const,
  aliases: ["combos", "aliases"] as const,
  disabledModels: ["combos", "disabled-models"] as const,
  cooldowns: ["combos", "cooldowns"] as const,
};

export function useCombos() {
  return useQuery<Combo[], Error>({
    queryKey: combosQueryKeys.list,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/combos");
      if (data) {
        return data.combos;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateComboMutation() {
  const queryClient = useQueryClient();

  return useMutation<Combo, Error, CreateComboRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/combos", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.list });
    },
  });
}

export interface UpdateComboVariables {
  id: string;
  body: UpdateComboRequest;
}

export function useUpdateComboMutation() {
  const queryClient = useQueryClient();

  return useMutation<Combo, Error, UpdateComboVariables>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/combos/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.list });
    },
  });
}

export function useDeleteComboMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error, response } = await apiClient.DELETE("/api/combos/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to delete combo (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.list });
    },
  });
}

export function useModelAliases() {
  return useQuery<ModelAlias[], Error>({
    queryKey: combosQueryKeys.aliases,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/aliases");
      if (data) {
        return data.aliases;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export interface UpsertModelAliasVariables {
  alias: string;
  body: UpsertModelAliasRequest;
}

export function useUpsertModelAliasMutation() {
  const queryClient = useQueryClient();

  return useMutation<ModelAlias, Error, UpsertModelAliasVariables>({
    mutationFn: async ({ alias, body }) => {
      const { data, error } = await apiClient.PUT("/api/aliases/{alias}", {
        params: { path: { alias } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.aliases });
    },
  });
}

export function useDeleteModelAliasMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (alias) => {
      const { error, response } = await apiClient.DELETE("/api/aliases/{alias}", {
        params: { path: { alias } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to delete model alias (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.aliases });
    },
  });
}

export function useDisabledModels() {
  return useQuery<DisabledModel[], Error>({
    queryKey: combosQueryKeys.disabledModels,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/models/disabled");
      if (data) {
        return data.disabled_models;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useDisableModelMutation() {
  const queryClient = useQueryClient();

  return useMutation<DisabledModel, Error, DisableModelRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.PUT("/api/models/disabled", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.disabledModels });
    },
  });
}

export function useEnableModelMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (model) => {
      const { error, response } = await apiClient.DELETE("/api/models/disabled", {
        params: { query: { model } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to re-enable model (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.disabledModels });
    },
  });
}

export function useModelCooldowns() {
  return useQuery<CooldownEntry[], Error>({
    queryKey: combosQueryKeys.cooldowns,
    refetchInterval: 10_000,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/models/cooldowns");
      if (data) {
        return data.cooldowns;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useClearModelCooldownMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (connectionId) => {
      const { error, response } = await apiClient.DELETE("/api/models/cooldowns/{connectionId}", {
        params: { path: { connectionId } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to clear cooldown (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: combosQueryKeys.cooldowns });
    },
  });
}
