import { useMutation, useQuery } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  ConsoleLogsResponse,
  TranslatePayloadRequest,
  TranslatePayloadResponse,
  TranslatorProtocolsResponse,
} from "./schemas";

export const translatorQueryKeys = {
  protocols: ["translator", "protocols"] as const,
  logs: (level: string) => ["translator", "logs", level] as const,
};

export function useTranslatorProtocols() {
  return useQuery<TranslatorProtocolsResponse, Error>({
    queryKey: translatorQueryKeys.protocols,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/translator/protocols");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useTranslatePayloadMutation() {
  return useMutation<TranslatePayloadResponse, Error, TranslatePayloadRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/translator/translate", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useConsoleLogsSnapshot(level: "ALL" | "DEBUG" | "INFO" | "WARN" | "ERROR") {
  return useQuery<ConsoleLogsResponse, Error>({
    queryKey: translatorQueryKeys.logs(level),
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/translator/logs", {
        params: {
          query: {
            limit: 100,
            level: level === "ALL" ? undefined : level,
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
