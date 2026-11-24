import {
  PROVIDER_REGISTRY,
  PROVIDER_REGISTRY_MAP,
  type ProviderAuthMode,
  type ProviderModelSpec,
  type ProviderRegistrySpec,
  type ProviderSourceCategory,
} from "./providerRegistry";
import type {
  CatalogAdapterDescriptor,
  ProviderConnectionSummary,
  ProviderNode,
  ProviderProtocol,
} from "./schemas";

export type { ProviderModelSpec };

export type ProviderCategory = "oauth" | "apikey" | "free" | "custom";
export type ProviderStatusFilter = "all" | "active" | "error" | "inactive" | "unconfigured";
export type ProviderCategoryTab = "all" | ProviderCategory;
export type ProviderSortOrder = "default" | "asc" | "desc";

export interface ProviderEntryStats {
  total: number;
  connected: number;
  error: number;
  disabled: number;
  allDisabled: boolean;
  isHealthy: boolean;
  latestErrorMessage: string | null;
  latestTestedAt: string | null;
}

export interface ProviderEntry {
  id: string;
  name: string;
  alias?: string;
  aliases?: string[];
  category: ProviderCategory;
  sourceCategory: ProviderSourceCategory | "custom";
  authType: ProviderAuthMode;
  protocol: ProviderProtocol;
  baseUrl: string;
  defaultModels: string[];
  models: ProviderModelSpec[];
  headers?: Record<string, string>;
  website?: string;
  apiKeyUrl?: string;
  signupUrl?: string;
  noticeText?: string;
  authHint?: string;
  noAuth?: boolean;
  hasOAuth?: boolean;
  deprecated?: boolean;
  serviceKinds?: string[];
  isCustom: boolean;
  isBuiltinAdapter: boolean;
  isNativeOAuth: boolean;
  adapter?: CatalogAdapterDescriptor;
  node?: ProviderNode;
  registrySpec?: ProviderRegistrySpec;
  connections: ProviderConnectionSummary[];
  stats: ProviderEntryStats;
}

export interface FleetTelemetry {
  activeProviders: number;
  totalCatalog: number;
  enabledConnections: number;
  totalConnections: number;
  attentionProviders: number;
  disabledProviders: number;
  oauthConnections: number;
  apiKeyConnections: number;
}

// Map Go backend adapter IDs to their matching registry metadata ID when names differ.
const ADAPTER_TO_REGISTRY_ID: Record<string, string> = {
  "gemini-oauth": "gemini-cli",
  "github-copilot": "github",
  "openai-responses": "openai",
  "aws-quick": "aws-quick",
};

export function getProviderMonogram(name: string, id: string): string {
  const source = (name || id || "PR").trim();
  const words = source.split(/[\s\-_/]+/).filter(Boolean);
  const first = words[0];
  const second = words[1];
  if (first && second) {
    return `${first[0] ?? ""}${second[0] ?? ""}`.toUpperCase();
  }
  return source.slice(0, 2).toUpperCase();
}

function computeEntryStats(connections: ProviderConnectionSummary[]): ProviderEntryStats {
  let connected = 0;
  let error = 0;
  let disabled = 0;
  let latestErrorMessage: string | null = null;
  let latestTestedAt: string | null = null;

  for (const conn of connections) {
    if (conn.last_tested_at && (!latestTestedAt || conn.last_tested_at > latestTestedAt)) {
      latestTestedAt = conn.last_tested_at;
    }
    if (!conn.enabled) {
      disabled += 1;
      continue;
    }
    if (conn.last_test_status === "error") {
      error += 1;
      if (conn.last_test_error) {
        latestErrorMessage = conn.last_test_error;
      }
    } else {
      connected += 1;
    }
  }

  const total = connections.length;
  const allDisabled = total > 0 && disabled === total;
  const isHealthy = !allDisabled && connected > 0 && error === 0;

  return {
    total,
    connected,
    error,
    disabled,
    allDisabled,
    isHealthy,
    latestErrorMessage,
    latestTestedAt,
  };
}

export function buildProviderEntries(
  adapters: CatalogAdapterDescriptor[],
  customNodes: ProviderNode[],
  connections: ProviderConnectionSummary[],
): ProviderEntry[] {
  const byProviderId = new Map<string, ProviderConnectionSummary[]>();
  for (const conn of connections) {
    const list = byProviderId.get(conn.provider_id) ?? [];
    list.push(conn);
    byProviderId.set(conn.provider_id, list);
  }

  const customNodeMap = new Map<string, ProviderNode>(customNodes.map((n) => [n.id, n]));
  const seen = new Set<string>();
  const registryCovered = new Set<string>();
  const entries: ProviderEntry[] = [];

  // 1. Built-in Go backend adapters (enriched with registry metadata)
  for (const adapter of adapters) {
    seen.add(adapter.id);
    const regId = ADAPTER_TO_REGISTRY_ID[adapter.id] ?? adapter.id;
    const regSpec = PROVIDER_REGISTRY_MAP.get(regId);
    if (regSpec) {
      registryCovered.add(regSpec.id);
    }
    const conns = byProviderId.get(adapter.id) ?? [];
    const isCompatAdapter =
      adapter.id === "openai-compatible" || adapter.id === "anthropic-compatible";
    const category: ProviderCategory = isCompatAdapter
      ? "custom"
      : adapter.auth_type === "oauth"
        ? "oauth"
        : regSpec?.category === "free"
          ? "free"
          : "apikey";

    const resolvedDefaultModels =
      regSpec && regSpec.defaultModels.length > 0 ? regSpec.defaultModels : adapter.default_models;
    const resolvedModels: ProviderModelSpec[] =
      regSpec?.models && regSpec.models.length > 0
        ? [...regSpec.models]
        : resolvedDefaultModels.map((id) => ({ id }));

    entries.push({
      id: adapter.id,
      name: regSpec?.name ?? adapter.display_name,
      alias: regSpec?.alias,
      aliases: regSpec?.aliases,
      category,
      sourceCategory: isCompatAdapter ? "custom" : (regSpec?.sourceCategory ?? "apikey"),
      authType: adapter.auth_type === "oauth" ? "oauth" : (regSpec?.authType ?? "api_key"),
      protocol: adapter.protocol,
      baseUrl: adapter.default_base_url || regSpec?.baseUrl || "",
      defaultModels: resolvedDefaultModels,
      models: resolvedModels,
      website: regSpec?.website,
      apiKeyUrl: regSpec?.apiKeyUrl,
      signupUrl: regSpec?.signupUrl,
      noticeText: regSpec?.noticeText,
      authHint: regSpec?.authHint,
      noAuth: regSpec?.noAuth,
      hasOAuth: adapter.auth_type === "oauth" || regSpec?.hasOAuth,
      deprecated: regSpec?.deprecated,
      serviceKinds: regSpec?.serviceKinds,
      isCustom: isCompatAdapter,
      isBuiltinAdapter: true,
      isNativeOAuth: adapter.auth_type === "oauth",
      adapter,
      registrySpec: regSpec,
      connections: conns,
      stats: computeEntryStats(conns),
    });
  }

  // 2. Built-in provider registry (all remaining providers not already covered by Go adapters)
  for (const regSpec of PROVIDER_REGISTRY) {
    if (seen.has(regSpec.id) || registryCovered.has(regSpec.id)) {
      continue;
    }
    seen.add(regSpec.id);
    const conns = byProviderId.get(regSpec.id) ?? [];
    const registeredNode = customNodeMap.get(regSpec.id);
    const resolvedRegModels =
      registeredNode && registeredNode.models.length > 0
        ? registeredNode.models
        : regSpec.defaultModels;
    const regMetaById = new Map((regSpec.models ?? []).map((m) => [m.id, m]));
    entries.push({
      id: regSpec.id,
      name: regSpec.name,
      alias: regSpec.alias,
      aliases: regSpec.aliases,
      category: regSpec.category,
      sourceCategory: regSpec.sourceCategory,
      authType: regSpec.authType,
      protocol: registeredNode?.protocol ?? regSpec.protocol,
      baseUrl: registeredNode?.base_url || regSpec.baseUrl,
      defaultModels: resolvedRegModels,
      models: resolvedRegModels.map((id) => regMetaById.get(id) ?? { id }),
      headers: registeredNode?.headers,
      website: regSpec.website,
      apiKeyUrl: regSpec.apiKeyUrl,
      signupUrl: regSpec.signupUrl,
      noticeText: regSpec.noticeText,
      authHint: regSpec.authHint,
      noAuth: regSpec.noAuth,
      hasOAuth: regSpec.hasOAuth,
      deprecated: regSpec.deprecated,
      serviceKinds: regSpec.serviceKinds,
      isCustom: false,
      isBuiltinAdapter: false,
      isNativeOAuth: false,
      node: registeredNode,
      registrySpec: regSpec,
      connections: conns,
      stats: computeEntryStats(conns),
    });
  }

  // 3. User-registered custom endpoints (not in built-in registry or Go adapters)
  for (const node of customNodes) {
    if (seen.has(node.id)) {
      continue;
    }
    seen.add(node.id);
    const conns = byProviderId.get(node.id) ?? [];
    entries.push({
      id: node.id,
      name: node.name,
      category: "custom",
      sourceCategory: "custom",
      authType: "api_key",
      protocol: node.protocol,
      baseUrl: node.base_url,
      defaultModels: node.models,
      models: node.models.map((id) => ({ id })),
      headers: node.headers,
      isCustom: true,
      isBuiltinAdapter: false,
      isNativeOAuth: false,
      node,
      connections: conns,
      stats: computeEntryStats(conns),
    });
  }

  // 4. Orphan connections (safety fallback)
  for (const [providerId, conns] of byProviderId.entries()) {
    if (seen.has(providerId)) {
      continue;
    }
    const first = conns[0];
    entries.push({
      id: providerId,
      name: providerId,
      category: first?.auth_type === "oauth" ? "oauth" : "custom",
      sourceCategory: "custom",
      authType: first?.auth_type ?? "api_key",
      protocol: first?.protocol ?? "openai-chat",
      baseUrl: first?.base_url ?? "",
      defaultModels: first?.models ?? [],
      models: (first?.models ?? []).map((id) => ({ id })),
      isCustom: first?.auth_type !== "oauth",
      isBuiltinAdapter: false,
      isNativeOAuth: false,
      connections: conns,
      stats: computeEntryStats(conns),
    });
  }

  return entries;
}

export function computeFleetTelemetry(
  entries: ProviderEntry[],
  connections: ProviderConnectionSummary[],
): FleetTelemetry {
  let activeProviders = 0;
  let attentionProviders = 0;
  let disabledProviders = 0;

  for (const entry of entries) {
    if (entry.stats.total === 0) {
      continue;
    }
    if (entry.stats.allDisabled) {
      disabledProviders += 1;
    } else if (entry.stats.error > 0) {
      attentionProviders += 1;
      if (entry.stats.connected > 0) {
        activeProviders += 1;
      }
    } else if (entry.stats.connected > 0) {
      activeProviders += 1;
    }
  }

  let enabledConnections = 0;
  let oauthConnections = 0;
  let apiKeyConnections = 0;
  for (const conn of connections) {
    if (conn.enabled) {
      enabledConnections += 1;
    }
    if (conn.auth_type === "oauth") {
      oauthConnections += 1;
    } else {
      apiKeyConnections += 1;
    }
  }

  return {
    activeProviders,
    totalCatalog: entries.length,
    enabledConnections,
    totalConnections: connections.length,
    attentionProviders,
    disabledProviders,
    oauthConnections,
    apiKeyConnections,
  };
}

export function matchesStatusFilter(entry: ProviderEntry, filter: ProviderStatusFilter): boolean {
  switch (filter) {
    case "all":
      return true;
    case "active":
      return entry.stats.total > 0 && !entry.stats.allDisabled && entry.stats.connected > 0;
    case "error":
      return entry.stats.total > 0 && !entry.stats.allDisabled && entry.stats.error > 0;
    case "inactive":
      return entry.stats.allDisabled;
    case "unconfigured":
      return entry.stats.total === 0;
  }
}

export function matchesProviderSearch(entry: ProviderEntry, query: string): boolean {
  const needle = query.trim().toLowerCase();
  if (!needle) {
    return true;
  }
  if (
    entry.id.toLowerCase().includes(needle) ||
    entry.name.toLowerCase().includes(needle) ||
    entry.protocol.toLowerCase().includes(needle) ||
    entry.baseUrl.toLowerCase().includes(needle) ||
    entry.category.toLowerCase().includes(needle) ||
    entry.sourceCategory.toLowerCase().includes(needle) ||
    (entry.alias?.toLowerCase().includes(needle) ?? false)
  ) {
    return true;
  }
  if (entry.aliases?.some((a) => a.toLowerCase().includes(needle))) {
    return true;
  }
  if (entry.serviceKinds?.some((k) => k.toLowerCase().includes(needle))) {
    return true;
  }
  if (entry.defaultModels.some((m) => m.toLowerCase().includes(needle))) {
    return true;
  }
  return entry.connections.some(
    (c) =>
      c.name.toLowerCase().includes(needle) ||
      c.models.some((m) => m.toLowerCase().includes(needle)),
  );
}

export function countStatusBuckets(entries: ProviderEntry[]): Record<ProviderStatusFilter, number> {
  const counts: Record<ProviderStatusFilter, number> = {
    all: entries.length,
    active: 0,
    error: 0,
    inactive: 0,
    unconfigured: 0,
  };
  for (const entry of entries) {
    if (matchesStatusFilter(entry, "active")) counts.active += 1;
    if (matchesStatusFilter(entry, "error")) counts.error += 1;
    if (matchesStatusFilter(entry, "inactive")) counts.inactive += 1;
    if (matchesStatusFilter(entry, "unconfigured")) counts.unconfigured += 1;
  }
  return counts;
}

export function countCategoryBuckets(
  entries: ProviderEntry[],
): Record<ProviderCategoryTab, number> {
  const counts: Record<ProviderCategoryTab, number> = {
    all: entries.length,
    oauth: 0,
    apikey: 0,
    free: 0,
    custom: 0,
  };
  for (const entry of entries) {
    counts[entry.category] += 1;
  }
  return counts;
}

function getEntrySortRank(entry: ProviderEntry): number {
  if (entry.stats.total > 0 && !entry.stats.allDisabled && entry.stats.error > 0) return 0;
  if (entry.stats.total > 0 && !entry.stats.allDisabled && entry.stats.connected > 0) return 1;
  if (entry.stats.allDisabled) return 2;
  if (entry.isBuiltinAdapter) return 3;
  return 4;
}

export function sortProviderEntries(
  entries: ProviderEntry[],
  sortOrder: ProviderSortOrder = "default",
): ProviderEntry[] {
  if (sortOrder === "asc") {
    return [...entries].sort((a, b) =>
      a.name.localeCompare(b.name, undefined, { sensitivity: "base" }),
    );
  }
  if (sortOrder === "desc") {
    return [...entries].sort((a, b) =>
      b.name.localeCompare(a.name, undefined, { sensitivity: "base" }),
    );
  }
  return [...entries].sort((a, b) => {
    const rankDiff = getEntrySortRank(a) - getEntrySortRank(b);
    if (rankDiff !== 0) {
      return rankDiff;
    }
    const prioA = a.registrySpec?.priority ?? 500;
    const prioB = b.registrySpec?.priority ?? 500;
    if (prioA !== prioB) {
      return prioA - prioB;
    }
    return a.name.localeCompare(b.name);
  });
}
