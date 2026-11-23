"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreatePromptGuardRuleRequest,
  PromptGuardConfig,
  PromptGuardConfigPatch,
  PromptGuardEventList,
  PromptGuardRule,
  PromptGuardRuleList,
  PromptGuardTestRequest,
  PromptGuardTestResult,
  UpdatePromptGuardRuleRequest,
} from "./schemas";

export const promptGuardQueryKeys = {
  all: ["promptinjection"] as const,
  config: ["promptinjection", "config"] as const,
  rules: ["promptinjection", "rules"] as const,
  events: ["promptinjection", "events"] as const,
};

export function usePromptGuardConfig() {
  return useQuery<PromptGuardConfig, Error>({
    queryKey: promptGuardQueryKeys.config,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/prompt-guard/config");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUpdatePromptGuardConfigMutation() {
  const queryClient = useQueryClient();
  return useMutation<PromptGuardConfig, Error, PromptGuardConfigPatch>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/prompt-guard/config", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: promptGuardQueryKeys.config });
    },
  });
}

export function usePromptGuardRules() {
  return useQuery<PromptGuardRuleList, Error>({
    queryKey: promptGuardQueryKeys.rules,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/prompt-guard/rules");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreatePromptGuardRuleMutation() {
  const queryClient = useQueryClient();
  return useMutation<PromptGuardRule, Error, CreatePromptGuardRuleRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/prompt-guard/rules", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: promptGuardQueryKeys.rules });
    },
  });
}

export function useUpdatePromptGuardRuleMutation() {
  const queryClient = useQueryClient();
  return useMutation<PromptGuardRule, Error, { id: string; body: UpdatePromptGuardRuleRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/prompt-guard/rules/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: promptGuardQueryKeys.rules });
    },
  });
}

export function useDeletePromptGuardRuleMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/prompt-guard/rules/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: promptGuardQueryKeys.rules });
    },
  });
}

export function useTestPromptGuardMutation() {
  return useMutation<PromptGuardTestResult, Error, PromptGuardTestRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/prompt-guard/test", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function usePromptGuardEvents() {
  return useQuery<PromptGuardEventList, Error>({
    queryKey: promptGuardQueryKeys.events,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/prompt-guard/events", {
        params: { query: { limit: 50 } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useClearPromptGuardEventsMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: async () => {
      const { error } = await apiClient.DELETE("/api/prompt-guard/events");
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: promptGuardQueryKeys.events });
    },
  });
}
