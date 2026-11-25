import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";
import {
  classifyTestError,
  comboTestResultKey,
  flattenCatalogResponse,
  modelTestResultKey,
  normalizeTestStatus,
  type RawCatalogGroup,
} from "./catalogUtils";
import type {
  ApiKeyAccessItem,
  CatalogModelRow,
  CatalogTestResult,
  ProviderHealthState,
} from "./types";

export const modelCatalogQueryKeys = {
  catalog: ["models", "catalog"] as const,
  healthMatrix: ["providers", "health-matrix"] as const,
  apiKeys: ["keys", "catalog-access"] as const,
};

export function useModelCatalog() {
  return useQuery<CatalogModelRow[], Error>({
    queryKey: modelCatalogQueryKeys.catalog,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/models/catalog", {
        params: { query: { includeHidden: "1" } },
      });
      if (data) {
        return flattenCatalogResponse(data.catalog as Record<string, RawCatalogGroup>);
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useProviderHealthMatrix() {
  return useQuery<Record<string, ProviderHealthState>, Error>({
    queryKey: modelCatalogQueryKeys.healthMatrix,
    refetchInterval: 30_000,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/providers/health-matrix");
      if (data) {
        const out: Record<string, ProviderHealthState> = {};
        for (const item of data.providers) {
          if (item.provider) {
            out[item.provider] = item.state;
          }
        }
        return out;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCatalogApiKeys() {
  return useQuery<ApiKeyAccessItem[], Error>({
    queryKey: modelCatalogQueryKeys.apiKeys,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/keys");
      if (data) {
        return data.keys.map((k) => {
          const allowed = Array.isArray(k.allowed_models) ? k.allowed_models : [];
          const isAll = allowed.length === 0 || allowed.includes("*");
          return {
            id: k.id,
            name: k.name,
            prefix: k.key_prefix,
            disabled: k.disabled,
            allowedModels: allowed,
            modelAccessMode: isAll ? "all" : "restricted",
          };
        });
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useUpdateKeyAllowedModelsMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, { keyId: string; allowedModels: string[] }>({
    mutationFn: async ({ keyId, allowedModels }) => {
      const { error } = await apiClient.PATCH("/api/keys/{id}", {
        params: { path: { id: keyId } },
        body: { allowed_models: allowedModels },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: modelCatalogQueryKeys.apiKeys });
    },
  });
}

export function useToggleModelVisibilityMutation() {
  const queryClient = useQueryClient();

  return useMutation<void, Error, { qualifiedIds: string[]; hide: boolean }>({
    mutationFn: async ({ qualifiedIds, hide }) => {
      for (const qualifiedId of qualifiedIds) {
        if (hide) {
          const { error } = await apiClient.PUT("/api/models/disabled", {
            body: { model: qualifiedId, reason: "Hidden from Model Catalog" },
          });
          if (error) {
            throw new Error(extractErrorMessage(error));
          }
        } else {
          const { error, response } = await apiClient.DELETE("/api/models/disabled", {
            params: { query: { model: qualifiedId } },
          });
          if (error && response.status !== 404) {
            throw new Error(extractErrorMessage(error));
          }
        }
      }
    },
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: modelCatalogQueryKeys.catalog }),
        queryClient.invalidateQueries({ queryKey: ["combos", "disabled-models"] }),
      ]);
    },
  });
}

export async function runSingleModelTest(
  providerId: string,
  modelId: string,
  signal?: AbortSignal,
): Promise<CatalogTestResult> {
  const testedAt = Date.now();
  try {
    const res = await fetch("/api/models/test", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ providerId, modelId }),
      signal,
    });
    const payload = (await res.json().catch(() => ({}))) as {
      ok?: boolean;
      status?: string;
      latencyMs?: number;
      statusCode?: number;
      error?: string;
      rateLimited?: boolean;
      isQuota?: boolean;
      isTimeout?: boolean;
    };
    if (!res.ok) {
      const errMsg = payload.error || `HTTP ${res.status}`;
      return {
        id: modelTestResultKey(providerId, modelId),
        targetType: "model",
        providerId,
        modelId,
        status: "error",
        statusCode: res.status,
        error: errMsg,
        errorClass: classifyTestError(errMsg, res.status),
        testedAt,
      };
    }
    const status = payload.isTimeout
      ? "error"
      : normalizeTestStatus(payload.status, payload.latencyMs);
    const errMsg = status === "error" ? payload.error : undefined;
    return {
      id: modelTestResultKey(providerId, modelId),
      targetType: "model",
      providerId,
      modelId,
      status,
      latencyMs: payload.latencyMs,
      statusCode: payload.statusCode,
      error: errMsg,
      errorClass:
        status === "error"
          ? classifyTestError(errMsg, payload.statusCode, {
              rateLimited: payload.rateLimited,
              isQuota: payload.isQuota,
              isTimeout: payload.isTimeout,
            })
          : undefined,
      testedAt,
    };
  } catch (err) {
    const msg = err instanceof Error ? err.message : "Network error";
    return {
      id: modelTestResultKey(providerId, modelId),
      targetType: "model",
      providerId,
      modelId,
      status: "error",
      error: msg,
      errorClass: classifyTestError(msg),
      testedAt,
    };
  }
}

export type ModelIdentityResult = components["schemas"]["ModelIdentityResult"];

const FALLBACK_IDENTITY_LIMITATION =
  "Detection combines upstream response metadata, vendor headers, and a neutral self-identification probe. Proxies that rewrite all metadata/headers and inject strict system prompts can still disguise the underlying model, so treat results as indicative evidence.";

export async function runModelIdentityProbe(
  providerId: string,
  modelId: string,
  signal?: AbortSignal,
): Promise<ModelIdentityResult> {
  try {
    const { data, error, response } = await apiClient.POST("/api/models/identify", {
      body: { providerId, modelId },
      signal,
    });
    if (data) {
      return data;
    }
    return {
      ok: false,
      claimedModel: modelId,
      verdict: "error",
      confidence: "low",
      signals: [],
      limitationNote: FALLBACK_IDENTITY_LIMITATION,
      latencyMs: 0,
      statusCode: response.status,
      error: extractErrorMessage(error),
    };
  } catch (err) {
    return {
      ok: false,
      claimedModel: modelId,
      verdict: "error",
      confidence: "low",
      signals: [],
      limitationNote: FALLBACK_IDENTITY_LIMITATION,
      latencyMs: 0,
      error: err instanceof Error ? err.message : "Network error",
    };
  }
}

export async function runSingleComboTest(
  comboName: string,
  signal?: AbortSignal,
): Promise<CatalogTestResult> {
  const testedAt = Date.now();
  try {
    const res = await fetch("/api/combos/test", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ comboName }),
      signal,
    });
    const payload = (await res.json().catch(() => ({}))) as {
      resolvedBy?: string;
      results?: Array<{
        target: string;
        status: string;
        latencyMs?: number;
        statusCode?: number;
        error?: string;
      }>;
      error?: string;
    };
    const items = Array.isArray(payload.results) ? payload.results : [];
    const okItem = items.find((i) => i.status === "ok");
    const ok = res.ok && (Boolean(payload.resolvedBy) || Boolean(okItem));
    const latencyMs = okItem?.latencyMs ?? items[0]?.latencyMs;
    const statusCode = okItem?.statusCode ?? items[0]?.statusCode ?? res.status;
    const status = ok ? normalizeTestStatus("ok", latencyMs) : "error";
    const errMsg =
      status === "error"
        ? (items.map((i) => i.error).find(Boolean) ?? payload.error ?? `HTTP ${res.status}`)
        : undefined;
    return {
      id: comboTestResultKey(comboName),
      targetType: "combo",
      comboName,
      status,
      latencyMs,
      statusCode,
      resolvedBy: payload.resolvedBy,
      error: errMsg,
      errorClass: status === "error" ? classifyTestError(errMsg, statusCode) : undefined,
      testedAt,
    };
  } catch (err) {
    const msg = err instanceof Error ? err.message : "Network error";
    return {
      id: comboTestResultKey(comboName),
      targetType: "combo",
      comboName,
      status: "error",
      error: msg,
      errorClass: classifyTestError(msg),
      testedAt,
    };
  }
}
