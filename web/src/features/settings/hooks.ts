"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { authKeys } from "@/features/auth/hooks";
import { systemKeys } from "@/features/system/hooks";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";
import type {
  PatchSettingsRequest,
  RollbackSettingRequest,
  SettingHistoryItem,
  SettingItem,
} from "./schemas";

export type ChangePasswordRequest = components["schemas"]["ChangePasswordRequest"];
export type ChangePasswordResponse = components["schemas"]["ChangePasswordResponse"];
export type UpdateAuthPolicyRequest = components["schemas"]["UpdateAuthPolicyRequest"];
export type AuthPolicyResponse = components["schemas"]["AuthPolicyResponse"];

export const settingsQueryKeys = {
  all: ["settings"] as const,
  history: (limit: number) => ["settings", "history", limit] as const,
  historyRoot: ["settings", "history"] as const,
};

export function useSettings() {
  return useQuery<SettingItem[], Error>({
    queryKey: settingsQueryKeys.all,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/settings");
      if (data) {
        return data.settings;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useSettingsHistory(limit = 50) {
  return useQuery<SettingHistoryItem[], Error>({
    queryKey: settingsQueryKeys.history(limit),
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/settings/history", {
        params: { query: { limit } },
      });
      if (data) {
        return data.items;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function usePatchSettingsMutation() {
  const queryClient = useQueryClient();

  return useMutation<SettingItem[], Error, PatchSettingsRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.PATCH("/api/settings", {
        body: payload,
      });
      if (data) {
        return data.settings;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (updatedSettings) => {
      queryClient.setQueryData<SettingItem[]>(settingsQueryKeys.all, updatedSettings);
      await queryClient.invalidateQueries({ queryKey: settingsQueryKeys.historyRoot });
    },
  });
}

export function useRollbackSettingMutation() {
  const queryClient = useQueryClient();

  return useMutation<SettingItem, Error, RollbackSettingRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/settings/rollback", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: settingsQueryKeys.all }),
        queryClient.invalidateQueries({ queryKey: settingsQueryKeys.historyRoot }),
      ]);
    },
  });
}

export function useChangePasswordMutation() {
  const queryClient = useQueryClient();

  return useMutation<ChangePasswordResponse, Error, ChangePasswordRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.PUT("/api/auth/password", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: authKeys.session });
    },
  });
}

export function useUpdateAuthPolicyMutation() {
  const queryClient = useQueryClient();

  return useMutation<AuthPolicyResponse, Error, UpdateAuthPolicyRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.PUT("/api/auth/policy", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: authKeys.session }),
        queryClient.invalidateQueries({ queryKey: systemKeys.initStatus }),
      ]);
    },
  });
}
