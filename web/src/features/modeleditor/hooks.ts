"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreateModelBypassRuleRequest,
  CreateModelOverrideRequest,
  ModelBypassRule,
  ModelBypassRuleList,
  ModelEditorPreviewRequest,
  ModelEditorPreviewResult,
  ModelOverride,
  ModelOverrideList,
  UpdateModelBypassRuleRequest,
  UpdateModelOverrideRequest,
} from "./schemas";

export const modelEditorQueryKeys = {
  all: ["modeleditor"] as const,
  overrides: ["modeleditor", "overrides"] as const,
  bypass: ["modeleditor", "bypass"] as const,
};

export function useModelOverrides() {
  return useQuery<ModelOverrideList, Error>({
    queryKey: modelEditorQueryKeys.overrides,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/model-editor/overrides");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateModelOverrideMutation() {
  const queryClient = useQueryClient();
  return useMutation<ModelOverride, Error, CreateModelOverrideRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/model-editor/overrides", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.overrides });
    },
  });
}

export function useUpdateModelOverrideMutation() {
  const queryClient = useQueryClient();
  return useMutation<ModelOverride, Error, { id: string; body: UpdateModelOverrideRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/model-editor/overrides/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.overrides });
    },
  });
}

export function useDeleteModelOverrideMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/model-editor/overrides/{id}", {
        params: { path: { id } },
      });
      if (error) throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.overrides });
    },
  });
}

export function useModelBypassRules() {
  return useQuery<ModelBypassRuleList, Error>({
    queryKey: modelEditorQueryKeys.bypass,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/model-editor/bypass");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateBypassRuleMutation() {
  const queryClient = useQueryClient();
  return useMutation<ModelBypassRule, Error, CreateModelBypassRuleRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/model-editor/bypass", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.bypass });
    },
  });
}

export function useUpdateBypassRuleMutation() {
  const queryClient = useQueryClient();
  return useMutation<ModelBypassRule, Error, { id: string; body: UpdateModelBypassRuleRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/model-editor/bypass/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.bypass });
    },
  });
}

export function useDeleteBypassRuleMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/model-editor/bypass/{id}", {
        params: { path: { id } },
      });
      if (error) throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.bypass });
    },
  });
}

export function usePreviewModelEditorMutation() {
  return useMutation<ModelEditorPreviewResult, Error, ModelEditorPreviewRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/model-editor/preview", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}
