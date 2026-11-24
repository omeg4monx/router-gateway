import { type QueryClient, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  CreateProviderConnectionRequest,
  CreateProviderNodeRequest,
  OAuthExchangeRequest,
  OAuthPasteRequest,
  OAuthPasteResponse,
  OAuthStartRequest,
  OAuthStartResponse,
  OAuthStatusResponse,
  OAuthStopResponse,
  ProviderConnectionSummary,
  ProviderNode,
  ProvidersCatalogResponse,
  TestProviderConnectionResponse,
  TestProviderNodeRequest,
  UpdateProviderConnectionRequest,
} from "./schemas";

export const providersQueryKeys = {
  all: ["providers"] as const,
  catalog: ["providers", "catalog"] as const,
  nodes: ["providers", "nodes"] as const,
  connections: ["providers", "connections"] as const,
  connection: (id: string) => ["providers", "connections", id] as const,
};

const MODEL_CATALOG_QUERY_KEY = ["models", "catalog"] as const;
const ARENA_MODELS_QUERY_KEY = ["arena", "models"] as const;

async function invalidateModelRoutingQueries(queryClient: QueryClient, connectionId?: string) {
  const tasks: Promise<void>[] = [
    queryClient.invalidateQueries({ queryKey: providersQueryKeys.connections }),
    queryClient.invalidateQueries({ queryKey: MODEL_CATALOG_QUERY_KEY }),
    queryClient.invalidateQueries({ queryKey: ARENA_MODELS_QUERY_KEY }),
  ];
  if (connectionId) {
    tasks.push(
      queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(connectionId) }),
    );
  }
  await Promise.all(tasks);
}

export async function fetchConnectionModels(
  connectionId: string,
): Promise<FetchConnectionModelsResponse> {
  const response = await fetch(
    `/api/providers/connections/${encodeURIComponent(connectionId)}/models`,
    {
      method: "GET",
      credentials: "include",
      headers: { Accept: "application/json" },
    },
  );
  const data = await response.json().catch(() => null);
  if (response.ok && data) {
    return data as FetchConnectionModelsResponse;
  }
  throw new Error(extractErrorMessage(data) || `Failed to load models (HTTP ${response.status})`);
}

function syncConnectionModelsInBackground(queryClient: QueryClient, connectionId: string) {
  void fetchConnectionModels(connectionId)
    .then((res) => {
      if (res.source === "upstream" && res.models.length > 0) {
        return invalidateModelRoutingQueries(queryClient, connectionId);
      }
      return undefined;
    })
    .catch(() => {
      // Best-effort background sync; manual "Load Models" remains available.
    });
}

export function useProvidersCatalog() {
  return useQuery<ProvidersCatalogResponse, Error>({
    queryKey: providersQueryKeys.catalog,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/providers/catalog");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useProviderNodes() {
  return useQuery<ProviderNode[], Error>({
    queryKey: providersQueryKeys.nodes,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/providers/nodes");
      if (data) {
        return data.nodes;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateProviderNodeMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderNode, Error, CreateProviderNodeRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/providers/nodes", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.catalog }),
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.nodes }),
      ]);
    },
  });
}

export function useDeleteProviderNodeMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error, response } = await apiClient.DELETE("/api/providers/nodes/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to delete provider node (HTTP ${response.status})`);
      }
    },
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.catalog }),
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.nodes }),
      ]);
    },
  });
}

export function useProviderConnections() {
  return useQuery<ProviderConnectionSummary[], Error>({
    queryKey: providersQueryKeys.connections,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/providers/connections");
      if (data) {
        return data.connections;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useProviderConnection(id: string) {
  const trimmedId = id.trim();

  return useQuery<ProviderConnectionSummary, Error>({
    queryKey: providersQueryKeys.connection(trimmedId),
    enabled: trimmedId.length > 0,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/providers/connections/{id}", {
        params: { path: { id: trimmedId } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateProviderConnectionMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderConnectionSummary, Error, CreateProviderConnectionRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/providers/connections", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (created) => {
      await invalidateModelRoutingQueries(queryClient, created.id);
      if (created.enabled) {
        syncConnectionModelsInBackground(queryClient, created.id);
      }
    },
  });
}

export interface UpdateProviderConnectionVariables {
  id: string;
  body: UpdateProviderConnectionRequest;
}

export function useUpdateProviderConnectionMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderConnectionSummary, Error, UpdateProviderConnectionVariables>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/providers/connections/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (updated) => {
      await invalidateModelRoutingQueries(queryClient, updated.id);
    },
  });
}

export function useDeleteProviderConnectionMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error, response } = await apiClient.DELETE("/api/providers/connections/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
      if (!response.ok) {
        throw new Error(`Failed to delete provider connection (HTTP ${response.status})`);
      }
    },
    onSuccess: async (_data, id) => {
      await invalidateModelRoutingQueries(queryClient, id);
    },
  });
}

export function useTestProviderNodeDraftMutation() {
  return useMutation<TestProviderConnectionResponse, Error, TestProviderNodeRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/providers/node-drafts/test", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export interface DiscoveredModel {
  id: string;
  name?: string;
}

export interface FetchConnectionModelsResponse {
  connection_id: string;
  provider_id: string;
  models: DiscoveredModel[];
  source: "upstream" | "fallback";
  latency_ms: number;
  warning?: string;
}

export function useFetchConnectionModelsMutation() {
  const queryClient = useQueryClient();

  return useMutation<FetchConnectionModelsResponse, Error, string>({
    mutationFn: fetchConnectionModels,
    onSuccess: async (res) => {
      if (res.source === "upstream" && res.models.length > 0) {
        await invalidateModelRoutingQueries(queryClient, res.connection_id);
      }
    },
  });
}

export function useTestProviderConnectionMutation() {
  const queryClient = useQueryClient();

  return useMutation<TestProviderConnectionResponse, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.POST("/api/providers/connections/{id}/test", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (_data, id) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connections }),
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(id) }),
      ]);
    },
  });
}

export interface StartProviderOAuthVariables {
  provider: string;
  body: OAuthStartRequest;
}

export function useStartProviderOAuthMutation() {
  return useMutation<OAuthStartResponse, Error, StartProviderOAuthVariables>({
    mutationFn: async ({ provider, body }) => {
      const { data, error } = await apiClient.POST("/api/providers/oauth/{provider}/start", {
        params: { path: { provider } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export interface ExchangeProviderOAuthVariables {
  provider: string;
  body: OAuthExchangeRequest;
}

export function useExchangeProviderOAuthMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderConnectionSummary, Error, ExchangeProviderOAuthVariables>({
    mutationFn: async ({ provider, body }) => {
      const { data, error } = await apiClient.POST("/api/providers/oauth/{provider}/exchange", {
        params: { path: { provider } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (conn) => {
      await invalidateModelRoutingQueries(queryClient, conn.id);
      if (conn.enabled) {
        syncConnectionModelsInBackground(queryClient, conn.id);
      }
    },
  });
}

export async function fetchProviderOAuthStatus(
  provider: string,
  state: string,
): Promise<OAuthStatusResponse> {
  const { data, error } = await apiClient.GET("/api/providers/oauth/{provider}/status", {
    params: { path: { provider }, query: { state } },
  });
  if (data) {
    return data;
  }
  throw new Error(extractErrorMessage(error));
}

export interface PasteProviderOAuthVariables {
  provider: string;
  body: OAuthPasteRequest;
}

export function usePasteProviderOAuthMutation() {
  return useMutation<OAuthPasteResponse, Error, PasteProviderOAuthVariables>({
    mutationFn: async ({ provider, body }) => {
      const { data, error } = await apiClient.POST("/api/providers/oauth/{provider}/paste", {
        params: { path: { provider } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useStopProviderOAuthMutation() {
  return useMutation<OAuthStopResponse, Error, string>({
    mutationFn: async (provider) => {
      const { data, error } = await apiClient.POST("/api/providers/oauth/{provider}/stop", {
        params: { path: { provider } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}
