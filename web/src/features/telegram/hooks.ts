"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  TelegramActionResponse,
  TelegramCommandRequest,
  TelegramCommandResponse,
  TelegramConfig,
  TelegramEventListResponse,
  TelegramTestRequest,
  TelegramTestResponse,
  UpdateTelegramConfigRequest,
} from "./schemas";

export const telegramQueryKeys = {
  all: ["telegram"] as const,
  config: ["telegram", "config"] as const,
  events: ["telegram", "events"] as const,
};

export function useTelegramConfig() {
  return useQuery<TelegramConfig, Error>({
    queryKey: telegramQueryKeys.config,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/telegram/config");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useTelegramEvents() {
  return useQuery<TelegramEventListResponse, Error>({
    queryKey: telegramQueryKeys.events,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/telegram/events");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    refetchInterval: 8000,
  });
}

export function useUpdateTelegramConfigMutation() {
  const queryClient = useQueryClient();
  return useMutation<TelegramConfig, Error, UpdateTelegramConfigRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.PATCH("/api/telegram/config", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.all });
    },
  });
}

export function useRotateTelegramSecretMutation() {
  const queryClient = useQueryClient();
  return useMutation<TelegramConfig, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/telegram/secret/rotate");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.all });
    },
  });
}

export function useRegisterTelegramWebhookMutation() {
  const queryClient = useQueryClient();
  return useMutation<TelegramActionResponse, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/telegram/webhook/registration");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.events });
    },
  });
}

export function useDeleteTelegramWebhookMutation() {
  const queryClient = useQueryClient();
  return useMutation<TelegramActionResponse, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.DELETE("/api/telegram/webhook/registration");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.events });
    },
  });
}

export function useSendTelegramTestMutation() {
  const queryClient = useQueryClient();
  return useMutation<TelegramTestResponse, Error, TelegramTestRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/telegram/test", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.events });
    },
  });
}

export function useSimulateTelegramCommandMutation() {
  const queryClient = useQueryClient();
  return useMutation<TelegramCommandResponse, Error, TelegramCommandRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/telegram/command", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.events });
    },
  });
}

export function useClearTelegramEventsMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, void>({
    mutationFn: async () => {
      const { error } = await apiClient.DELETE("/api/telegram/events");
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: telegramQueryKeys.events });
    },
  });
}
