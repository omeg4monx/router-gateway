import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";

export type CloudSyncConfig = components["schemas"]["CloudSyncConfig"];
export type UpdateCloudSyncConfigRequest = components["schemas"]["UpdateCloudSyncConfigRequest"];
export type CloudSyncActionResponse = components["schemas"]["CloudSyncActionResponse"];
export type CloudSyncSnapshotBundle = components["schemas"]["CloudSyncSnapshotBundle"];
export type RestoreCloudSyncSnapshotRequest =
  components["schemas"]["RestoreCloudSyncSnapshotRequest"];
export type CloudSyncHistoryListResponse = components["schemas"]["CloudSyncHistoryListResponse"];

export const cloudSyncQueryKeys = {
  all: ["cloudsync"] as const,
  config: ["cloudsync", "config"] as const,
  history: ["cloudsync", "history"] as const,
};

export function useCloudSyncConfig() {
  return useQuery<CloudSyncConfig, Error>({
    queryKey: cloudSyncQueryKeys.config,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/cloud-sync/config");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCloudSyncHistory() {
  return useQuery<CloudSyncHistoryListResponse, Error>({
    queryKey: cloudSyncQueryKeys.history,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/cloud-sync/history");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUpdateCloudSyncConfigMutation() {
  const qc = useQueryClient();
  return useMutation<CloudSyncConfig, Error, UpdateCloudSyncConfigRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/cloud-sync/config", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: cloudSyncQueryKeys.all });
    },
  });
}

export function usePushCloudSyncMutation() {
  const qc = useQueryClient();
  return useMutation<CloudSyncActionResponse, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/cloud-sync/push");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: cloudSyncQueryKeys.all });
    },
  });
}

export function usePullCloudSyncMutation() {
  const qc = useQueryClient();
  return useMutation<CloudSyncActionResponse, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/cloud-sync/pull");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: cloudSyncQueryKeys.all });
    },
  });
}

export function useExportSnapshotMutation() {
  const qc = useQueryClient();
  return useMutation<CloudSyncSnapshotBundle, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.GET("/api/cloud-sync/snapshot");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: cloudSyncQueryKeys.history });
    },
  });
}

export function useRestoreSnapshotMutation() {
  const qc = useQueryClient();
  return useMutation<CloudSyncActionResponse, Error, RestoreCloudSyncSnapshotRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/cloud-sync/restore", { body });
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: cloudSyncQueryKeys.all });
    },
  });
}

export function useClearCloudSyncHistoryMutation() {
  const qc = useQueryClient();
  return useMutation<{ deleted: number }, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.DELETE("/api/cloud-sync/history");
      if (data) return data;
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void qc.invalidateQueries({ queryKey: cloudSyncQueryKeys.history });
    },
  });
}
