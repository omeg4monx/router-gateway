import type { CatalogModelRow, TestErrorClass, TestResultStatus } from "./types";

export const SLOW_LATENCY_MS = 5000;

export function modelTestResultKey(providerId: string, modelId: string): string {
  return `model:${providerId}:${modelId}`;
}

export function comboTestResultKey(comboName: string): string {
  return `combo:${comboName}`;
}

export function classifyTestError(
  error?: string,
  statusCode?: number,
  flags?: { rateLimited?: boolean; isQuota?: boolean; isTimeout?: boolean },
): TestErrorClass | undefined {
  const msg = error ?? "";
  if (flags?.rateLimited || statusCode === 429 || /rate.?limit|429/i.test(msg)) {
    return "rate-limited";
  }
  if (flags?.isQuota || /quota|insufficient balance|credit|billing|exhausted/i.test(msg)) {
    return "quota";
  }
  if (
    flags?.isTimeout ||
    statusCode === 408 ||
    statusCode === 504 ||
    /timeout|timed out|abort/i.test(msg)
  ) {
    return "timeout";
  }
  if (msg || statusCode) {
    return "other";
  }
  return undefined;
}

export function normalizeTestStatus(
  status: string | undefined,
  latencyMs?: number,
): TestResultStatus {
  if (status === "slow") {
    return "slow";
  }
  if (status === "ok") {
    return typeof latencyMs === "number" && latencyMs >= SLOW_LATENCY_MS ? "slow" : "ok";
  }
  return "error";
}

export interface RawCatalogGroup {
  provider: string;
  models: Array<{
    id: string;
    name?: string;
    type?: string;
    subtype?: string;
    custom?: boolean;
    free?: boolean;
    capabilities?: Record<string, boolean>;
    context_length?: number;
    max_output_tokens?: number;
    input_modalities?: string[];
    output_modalities?: string[];
    supported_endpoints?: string[];
    hidden?: boolean;
  }>;
}

export function flattenCatalogResponse(
  catalog: Record<string, RawCatalogGroup> | undefined | null,
): CatalogModelRow[] {
  if (!catalog || typeof catalog !== "object") {
    return [];
  }
  const rows: CatalogModelRow[] = [];
  for (const [providerId, group] of Object.entries(catalog)) {
    if (!group || !Array.isArray(group.models)) {
      continue;
    }
    const providerLabel =
      typeof group.provider === "string" && group.provider ? group.provider : providerId;
    for (const item of group.models) {
      if (!item || typeof item.id !== "string" || !item.id.trim()) {
        continue;
      }
      rows.push({
        providerId,
        provider: providerLabel,
        id: item.id,
        qualifiedId: `${providerId}/${item.id}`,
        name: typeof item.name === "string" && item.name.trim() ? item.name : item.id,
        type: typeof item.type === "string" && item.type ? item.type : "unknown",
        subtype: item.subtype,
        custom: item.custom === true ? true : undefined,
        free: item.free === true ? true : undefined,
        capabilities: item.capabilities,
        context_length: typeof item.context_length === "number" ? item.context_length : undefined,
        max_output_tokens:
          typeof item.max_output_tokens === "number" ? item.max_output_tokens : undefined,
        input_modalities: item.input_modalities,
        output_modalities: item.output_modalities,
        supported_endpoints: item.supported_endpoints,
        hidden: item.hidden === true ? true : undefined,
      });
    }
  }
  return rows;
}
