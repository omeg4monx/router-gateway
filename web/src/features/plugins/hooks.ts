import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";

export type PluginState = components["schemas"]["PluginState"];
export type PluginListResponse = components["schemas"]["PluginListResponse"];
export type UpdatePluginRequest = components["schemas"]["UpdatePluginRequest"];
export type PreviewPluginsRequest = components["schemas"]["PreviewPluginsRequest"];
export type PreviewPluginsResponse = components["schemas"]["PreviewPluginsResponse"];

export const pluginsQueryKeys = {
  all: ["plugins"] as const,
  list: ["plugins", "list"] as const,
};

export function usePluginsList() {
  return useQuery<PluginListResponse, Error>({
    queryKey: pluginsQueryKeys.list,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/plugins");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUpdatePluginMutation() {
  const qc = useQueryClient();
  return useMutation<PluginState, Error, { key: string; body: UpdatePluginRequest }>({
    mutationFn: async ({ key, body }) => {
      const { data, error } = await apiClient.PUT("/api/plugins/{key}", {
        params: { path: { key } },
        body,
      });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: pluginsQueryKeys.list });
    },
  });
}

export function usePreviewPluginsMutation() {
  return useMutation<PreviewPluginsResponse, Error, PreviewPluginsRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/plugins/preview", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}
