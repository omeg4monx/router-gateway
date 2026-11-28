import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreateModelOverrideRequest,
  ModelOverride,
  ModelOverrideList,
  UpdateModelOverrideRequest,
} from "./schemas";

export const modelEditorQueryKeys = {
  overrides: ["modeleditor", "overrides"] as const,
  arenaModels: ["arena", "models"] as const,
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
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.arenaModels });
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
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.arenaModels });
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
      void queryClient.invalidateQueries({ queryKey: modelEditorQueryKeys.arenaModels });
    },
  });
}

export interface QuickPromptTestInput {
  model: string;
  prompt: string;
  system_prompt?: string;
  temperature?: number;
  max_tokens?: number;
}

export interface QuickPromptTestResult {
  model: string;
  content: string;
  latencyMs: number;
  promptTokens?: number;
  completionTokens?: number;
}

export function useRunQuickPromptTestMutation() {
  return useMutation<QuickPromptTestResult, Error, QuickPromptTestInput>({
    mutationFn: async ({ model, prompt, system_prompt, temperature, max_tokens }) => {
      const { data: session, error: sessionErr } = await apiClient.POST("/api/chat/sessions", {
        body: {
          title: "Model Editor Test",
          model,
          system_prompt: system_prompt || undefined,
          temperature,
          max_tokens,
        },
      });
      if (!session) {
        throw new Error(extractErrorMessage(sessionErr));
      }

      try {
        const { data: reply, error: msgErr } = await apiClient.POST(
          "/api/chat/sessions/{id}/messages",
          {
            params: { path: { id: session.id } },
            body: { content: prompt, model },
          },
        );
        if (!reply) {
          throw new Error(extractErrorMessage(msgErr));
        }
        return {
          model: reply.session?.model || model,
          content: reply.assistant_message?.content || "",
          latencyMs: reply.assistant_message?.latency_ms ?? 0,
          promptTokens: reply.assistant_message?.prompt_tokens,
          completionTokens: reply.assistant_message?.completion_tokens,
        };
      } finally {
        await apiClient
          .DELETE("/api/chat/sessions/{id}", {
            params: { path: { id: session.id } },
          })
          .catch(() => {});
      }
    },
  });
}
