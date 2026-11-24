import { useMutation, useQuery } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  ApplyCLIToolConfigResponse,
  CLIToolConfigRequest,
  CLIToolID,
  CLIToolsCatalogResponse,
  GenerateCLIToolConfigResponse,
  RollbackCLIToolConfigResponse,
} from "./schemas";

export const cliToolsQueryKeys = {
  all: ["clitools"] as const,
  catalog: ["clitools", "catalog"] as const,
};

export function useCliToolsCatalog() {
  return useQuery<CLIToolsCatalogResponse, Error>({
    queryKey: cliToolsQueryKeys.catalog,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/cli-tools");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export interface CLIToolConfigMutationVariables {
  tool: CLIToolID;
  body: CLIToolConfigRequest;
}

export function useGenerateCliToolConfigMutation() {
  return useMutation<GenerateCLIToolConfigResponse, Error, CLIToolConfigMutationVariables>({
    mutationFn: async ({ tool, body }) => {
      const { data, error } = await apiClient.POST("/api/cli-tools/{tool}/generate", {
        params: { path: { tool } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useApplyCliToolConfigMutation() {
  return useMutation<ApplyCLIToolConfigResponse, Error, CLIToolConfigMutationVariables>({
    mutationFn: async ({ tool, body }) => {
      const { data, error } = await apiClient.POST("/api/cli-tools/{tool}/apply", {
        params: { path: { tool } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useRollbackCliToolConfigMutation() {
  return useMutation<RollbackCLIToolConfigResponse, Error, CLIToolID>({
    mutationFn: async (tool) => {
      const { data, error } = await apiClient.POST("/api/cli-tools/{tool}/rollback", {
        params: { path: { tool } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}
