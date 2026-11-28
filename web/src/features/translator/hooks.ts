import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  TranslatePayloadRequest,
  TranslatePayloadResponse,
  TranslatorProtocolsResponse,
} from "./schemas";

export const translatorQueryKeys = {
  protocols: ["translator", "protocols"] as const,
  logsRoot: ["translator", "logs"] as const,
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

export function useClearConsoleLogsMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ cleared: boolean }, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.DELETE("/api/translator/logs");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: translatorQueryKeys.logsRoot });
    },
  });
}
