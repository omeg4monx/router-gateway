"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  BackgroundJobListResponse,
  CreateNotificationRequest,
  NotificationListResponse,
  SystemNotification,
  TriggerJobResponse,
} from "./schemas";

export const notificationsQueryKeys = {
  all: ["notifications"] as const,
  list: (filters: { unreadOnly?: boolean; severity?: string }) =>
    ["notifications", "list", filters] as const,
  jobs: ["notifications", "jobs"] as const,
};

export function useNotificationsList(filters: {
  unreadOnly?: boolean;
  severity?: "info" | "success" | "warning" | "error";
}) {
  return useQuery<NotificationListResponse, Error>({
    queryKey: notificationsQueryKeys.list(filters),
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/notifications", {
        params: {
          query: {
            unread_only: filters.unreadOnly ? true : undefined,
            severity: filters.severity || undefined,
            limit: 100,
          },
        },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useNotificationsStream() {
  const queryClient = useQueryClient();
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    if (typeof EventSource === "undefined") {
      return;
    }
    const source = new EventSource("/api/notifications/stream", { withCredentials: true });
    source.onopen = () => setConnected(true);
    source.onerror = () => setConnected(false);
    source.addEventListener("notification", () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    });

    return () => {
      source.close();
      setConnected(false);
    };
  }, [queryClient]);

  return connected;
}

export function useCreateNotificationMutation() {
  const queryClient = useQueryClient();
  return useMutation<SystemNotification, Error, CreateNotificationRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/notifications", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    },
  });
}

export function useMarkNotificationReadMutation() {
  const queryClient = useQueryClient();
  return useMutation<SystemNotification, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.POST("/api/notifications/{id}/read", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    },
  });
}

export function useMarkAllNotificationsReadMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ updated: number }, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/notifications/read-all");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    },
  });
}

export function useDeleteNotificationMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ deleted: boolean }, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.DELETE("/api/notifications/{id}", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    },
  });
}

export function useClearNotificationsMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ cleared: number }, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.DELETE("/api/notifications");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    },
  });
}

export function useBackgroundJobs() {
  return useQuery<BackgroundJobListResponse, Error>({
    queryKey: notificationsQueryKeys.jobs,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/jobs");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 15_000,
  });
}

export function useTriggerJobMutation() {
  const queryClient = useQueryClient();
  return useMutation<TriggerJobResponse, Error, string>({
    mutationFn: async (name) => {
      const { data, error } = await apiClient.POST("/api/jobs/{name}/trigger", {
        params: { path: { name } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.jobs });
      void queryClient.invalidateQueries({ queryKey: notificationsQueryKeys.all });
    },
  });
}
