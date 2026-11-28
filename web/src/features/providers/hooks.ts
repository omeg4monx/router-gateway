import { type QueryClient, useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { ProviderAccountStrategy, ProviderStrategyState } from "./providerStrategy";
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
  strategies: ["providers", "strategies"] as const,
  strategy: (providerId: string) => ["providers", "strategies", providerId] as const,
};

const MODEL_CATALOG_QUERY_KEY = ["models", "catalog"] as const;
const ARENA_MODELS_QUERY_KEY = ["arena", "models"] as const;

export async function invalidateModelRoutingQueries(
  queryClient: QueryClient,
  connectionId?: string,
) {
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

export async function runWithConcurrency<T, R>(
  items: readonly T[],
  limit: number,
  worker: (item: T, index: number) => Promise<R>,
): Promise<R[]> {
  if (items.length === 0) {
    return [];
  }
  const concurrency = Math.max(1, Math.min(Math.floor(limit) || 1, items.length));
  const results = new Array<R>(items.length);
  let nextIndex = 0;

  const runWorker = async () => {
    while (nextIndex < items.length) {
      const current = nextIndex;
      nextIndex += 1;
      results[current] = await worker(items[current] as T, current);
    }
  };

  await Promise.all(Array.from({ length: concurrency }, () => runWorker()));
  return results;
}

export async function fetchConnectionModels(
  connectionId: string,
): Promise<FetchConnectionModelsResponse> {
  const { data, error } = await apiClient.GET("/api/providers/connections/{id}/models", {
    params: { path: { id: connectionId } },
  });
  if (data) {
    return data;
  }
  throw new Error(extractErrorMessage(error));
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
    staleTime: 5 * 60_000,
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

export async function updateProviderConnection({
  id,
  body,
}: UpdateProviderConnectionVariables): Promise<ProviderConnectionSummary> {
  const { data, error } = await apiClient.PATCH("/api/providers/connections/{id}", {
    params: { path: { id } },
    body,
  });
  if (data) {
    return data;
  }
  throw new Error(extractErrorMessage(error));
}

export function useUpdateProviderConnectionMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderConnectionSummary, Error, UpdateProviderConnectionVariables>({
    mutationFn: updateProviderConnection,
    onSuccess: async (updated) => {
      await invalidateModelRoutingQueries(queryClient, updated.id);
    },
  });
}

export interface BatchToggleConnectionsVariables {
  connectionIds: string[];
  enabled: boolean;
}

export function useBatchToggleProviderConnectionsMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderConnectionSummary[], Error, BatchToggleConnectionsVariables>({
    mutationFn: async ({ connectionIds, enabled }) =>
      runWithConcurrency(connectionIds, 4, (id) =>
        updateProviderConnection({ id, body: { enabled } }),
      ),
    onSuccess: async (updated) => {
      await invalidateModelRoutingQueries(queryClient);
      await Promise.all(
        updated.map((conn) =>
          queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(conn.id) }),
        ),
      );
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
  is_free?: boolean;
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

export async function testProviderConnection(id: string): Promise<TestProviderConnectionResponse> {
  const { data, error } = await apiClient.POST("/api/providers/connections/{id}/test", {
    params: { path: { id } },
  });
  if (data) {
    return data;
  }
  throw new Error(extractErrorMessage(error));
}

export function useTestProviderConnectionMutation() {
  const queryClient = useQueryClient();

  return useMutation<TestProviderConnectionResponse, Error, string>({
    mutationFn: testProviderConnection,
    onSuccess: async (_data, id) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connections }),
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(id) }),
      ]);
    },
  });
}

export interface BatchTestConnectionResult {
  connectionId: string;
  connectionName: string;
  providerId: string;
  ok: boolean;
  statusCode: number;
  latencyMs: number;
  error?: string;
}

export function useBatchTestProviderConnectionsMutation() {
  const queryClient = useQueryClient();

  return useMutation<BatchTestConnectionResult[], Error, ProviderConnectionSummary[]>({
    mutationFn: async (connections) =>
      runWithConcurrency(connections, 4, async (conn) => {
        try {
          const res = await testProviderConnection(conn.id);
          return {
            connectionId: conn.id,
            connectionName: conn.name,
            providerId: conn.provider_id,
            ok: res.ok,
            statusCode: res.status_code,
            latencyMs: res.latency_ms,
            error: res.error || undefined,
          };
        } catch (err) {
          return {
            connectionId: conn.id,
            connectionName: conn.name,
            providerId: conn.provider_id,
            ok: false,
            statusCode: 0,
            latencyMs: 0,
            error: err instanceof Error ? err.message : "Test failed",
          };
        }
      }),
    onSuccess: async (_results, connections) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.connections }),
        ...connections.map((conn) =>
          queryClient.invalidateQueries({ queryKey: providersQueryKeys.connection(conn.id) }),
        ),
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

export function useProviderStrategy(providerId: string) {
  const trimmedId = providerId.trim();

  return useQuery<ProviderStrategyState, Error>({
    queryKey: providersQueryKeys.strategy(trimmedId),
    enabled: trimmedId.length > 0,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/providers/{providerId}/strategy", {
        params: { path: { providerId: trimmedId } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export interface UpdateProviderStrategyVariables {
  providerId: string;
  strategy: ProviderAccountStrategy;
}

export function useUpdateProviderStrategyMutation() {
  const queryClient = useQueryClient();

  return useMutation<ProviderStrategyState, Error, UpdateProviderStrategyVariables>({
    mutationFn: async ({ providerId, strategy }) => {
      const { data, error } = await apiClient.PUT("/api/providers/{providerId}/strategy", {
        params: { path: { providerId } },
        body: { strategy },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (updated, { providerId }) => {
      queryClient.setQueryData(providersQueryKeys.strategy(providerId), updated);
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.strategy(providerId) }),
        queryClient.invalidateQueries({ queryKey: providersQueryKeys.strategies }),
      ]);
    },
  });
}
