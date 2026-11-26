import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";
import type {
  ApiKeySecretResponse,
  ApiKeySummary,
  CreateApiKeyRequest,
  UpdateApiKeyRequest,
} from "./schemas";

export type ListApiKeysResponse = components["schemas"]["ListApiKeysResponse"];

export const keysQueryKeys = {
  all: ["keys"] as const,
};

export function useApiKeys() {
  return useQuery<ApiKeySummary[], Error>({
    queryKey: keysQueryKeys.all,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/keys");
      if (data) {
        return data.keys;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateApiKeyMutation() {
  const queryClient = useQueryClient();

  return useMutation<ApiKeySecretResponse, Error, CreateApiKeyRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/keys", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: keysQueryKeys.all });
    },
  });
}

export interface UpdateApiKeyVariables {
  id: string;
  body: UpdateApiKeyRequest;
}

export function useUpdateApiKeyMutation() {
  const queryClient = useQueryClient();

  return useMutation<ApiKeySummary, Error, UpdateApiKeyVariables>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/keys/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: keysQueryKeys.all });
    },
  });
}

export function useRotateApiKeyMutation() {
  const queryClient = useQueryClient();

  return useMutation<ApiKeySecretResponse, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.POST("/api/keys/{id}/rotate", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: keysQueryKeys.all });
    },
  });
}

export function useDeleteApiKeyMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error, response } = await apiClient.DELETE("/api/keys/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to delete API key (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: keysQueryKeys.all });
    },
  });
}
