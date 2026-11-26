import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";

export type InspectLegacyRequest = components["schemas"]["InspectLegacyRequest"];
export type InspectLegacyResponse = components["schemas"]["InspectLegacyResponse"];
export type RunLegacyImportRequest = components["schemas"]["RunLegacyImportRequest"];
export type RunLegacyImportResponse = components["schemas"]["RunLegacyImportResponse"];
export type ImporterHistoryListResponse = components["schemas"]["ImporterHistoryListResponse"];

export const importerQueryKeys = {
  all: ["importer"] as const,
  history: ["importer", "history"] as const,
};

export function useImporterHistory() {
  return useQuery<ImporterHistoryListResponse, Error>({
    queryKey: importerQueryKeys.history,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/importer/history");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useInspectLegacyMutation() {
  return useMutation<InspectLegacyResponse, Error, InspectLegacyRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/importer/inspect", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useRunLegacyImportMutation() {
  const qc = useQueryClient();
  return useMutation<RunLegacyImportResponse, Error, RunLegacyImportRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/importer/run", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: importerQueryKeys.history });
    },
  });
}

export function useClearImporterHistoryMutation() {
  const qc = useQueryClient();
  return useMutation<{ deleted: number }, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.DELETE("/api/importer/history");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: importerQueryKeys.history });
    },
  });
}
