import type { ProviderConnectionSummary } from "@/features/providers/schemas";
import type { ComboStep } from "./schemas";

export interface ProviderAccountGroup {
  providerId: string;
  displayName: string;
  accounts: ProviderConnectionSummary[];
  hasMultipleAccounts: boolean;
}

export function stripProviderPrefix(providerId: string, rawModel: string): string {
  const trimmed = rawModel.trim();
  const prefix = `${providerId}/`;
  if (trimmed.startsWith(prefix)) {
    return trimmed.slice(prefix.length);
  }
  return trimmed;
}

export function qualifyProviderModel(providerId: string, rawModel: string): string {
  const trimmed = rawModel.trim();
  if (!trimmed) return "";
  if (trimmed.includes("/")) return trimmed;
  return `${providerId}/${trimmed}`;
}

export function buildProviderAccountGroups(
  connections: readonly ProviderConnectionSummary[],
  providerLabels: Record<string, string> = {},
): ProviderAccountGroup[] {
  const byProvider = new Map<string, ProviderConnectionSummary[]>();

  for (const conn of connections) {
    if (!conn.enabled) continue;
    const list = byProvider.get(conn.provider_id) ?? [];
    list.push(conn);
    byProvider.set(conn.provider_id, list);
  }

  const groups: ProviderAccountGroup[] = [];
  for (const [providerId, accounts] of byProvider.entries()) {
    const sortedAccounts = [...accounts].sort((a, b) => {
      if (a.priority !== b.priority) return a.priority - b.priority;
      return a.name.localeCompare(b.name);
    });
    groups.push({
      providerId,
      displayName: providerLabels[providerId] ?? providerId,
      accounts: sortedAccounts,
      hasMultipleAccounts: sortedAccounts.length > 1,
    });
  }

  return groups.sort((a, b) => a.displayName.localeCompare(b.displayName));
}

export function getModelsForSelection(
  group: ProviderAccountGroup | undefined,
  connectionId: string,
  catalogModelsByProvider: Record<string, readonly string[]> = {},
): string[] {
  if (!group) return [];
  const sourceAccounts = connectionId
    ? group.accounts.filter((a) => a.id === connectionId)
    : group.accounts;

  const seen = new Set<string>();
  const models: string[] = [];

  const appendModel = (raw: string) => {
    const clean = stripProviderPrefix(group.providerId, raw);
    if (!clean || seen.has(clean)) return;
    seen.add(clean);
    models.push(clean);
  };

  for (const acc of sourceAccounts) {
    for (const raw of acc.models) {
      appendModel(raw);
    }
  }

  // Merge provider catalog models in dynamic mode, or as fallback when a pinned
  // connection has no explicit per-account model list.
  const shouldIncludeCatalog = !connectionId || models.length === 0;
  if (shouldIncludeCatalog) {
    const catalogList = catalogModelsByProvider[group.providerId] ?? [];
    for (const raw of catalogList) {
      appendModel(raw);
    }
  }

  return models.sort((a, b) => a.localeCompare(b));
}

export function isDuplicateModelStep(
  steps: readonly ComboStep[],
  qualifiedModel: string,
  connectionId: string,
): boolean {
  const normalizedConn = connectionId.trim();
  return steps.some(
    (step) =>
      step.kind === "model" &&
      step.model === qualifiedModel &&
      (step.connectionId ?? "").trim() === normalizedConn,
  );
}

export function formatStepAccountLabel(
  connectionId: string | undefined,
  connectionsById: ReadonlyMap<string, ProviderConnectionSummary>,
  isId: boolean,
): string {
  const trimmed = (connectionId ?? "").trim();
  if (!trimmed) {
    return isId ? "Dynamic (Semua Akun)" : "Dynamic (All Accounts)";
  }
  const conn = connectionsById.get(trimmed);
  return conn ? conn.name : trimmed;
}

export function formatAuthTypeBadge(authType: string, isId: boolean): string {
  switch (authType) {
    case "oauth":
      return "OAuth";
    case "api_key":
      return "API Key";
    case "none":
      return isId ? "Tanpa Auth" : "No Auth";
    default:
      return authType;
  }
}
