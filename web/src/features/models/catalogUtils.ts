import type { Combo } from "@/features/combos/schemas";
import type {
  ApiKeyAccessItem,
  CatalogComboRow,
  CatalogModelRow,
  CatalogTestResult,
  ComboCatalogFilters,
  ComboSortField,
  ModelCatalogFilters,
  ModelSortField,
  ProviderHealthState,
  SortDirection,
  TestErrorClass,
  TestResultStatus,
} from "./types";

export const TEST_RESULTS_STORAGE_KEY = "router.catalogTestResults.v1";
export const MAX_CACHED_TEST_RESULTS = 2000;
export const SLOW_LATENCY_MS = 5000;

export const DEFAULT_MODEL_FILTERS: ModelCatalogFilters = {
  query: "",
  providerId: "all",
  type: "all",
  subtype: "all",
  capability: "all",
  pricing: "all",
  providerHealth: "all",
  testResult: "all",
  visibility: "all",
  minContextLength: undefined,
  minMaxOutputTokens: undefined,
};

export const DEFAULT_COMBO_FILTERS: ComboCatalogFilters = {
  query: "",
  strategy: "all",
  status: "all",
  testResult: "all",
  minMembers: undefined,
  maxMembers: undefined,
};

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

export function loadCachedTestResults(): Record<string, CatalogTestResult> {
  if (typeof window === "undefined" || !window.localStorage) {
    return {};
  }
  try {
    const raw = window.localStorage.getItem(TEST_RESULTS_STORAGE_KEY);
    if (!raw) {
      return {};
    }
    const parsed: unknown = JSON.parse(raw);
    if (!parsed || typeof parsed !== "object" || Array.isArray(parsed)) {
      return {};
    }
    const out: Record<string, CatalogTestResult> = {};
    for (const [k, v] of Object.entries(parsed)) {
      if (
        v &&
        typeof v === "object" &&
        !Array.isArray(v) &&
        "status" in v &&
        typeof (v as CatalogTestResult).status === "string" &&
        "testedAt" in v &&
        typeof (v as CatalogTestResult).testedAt === "number" &&
        Number.isFinite((v as CatalogTestResult).testedAt)
      ) {
        out[k] = v as CatalogTestResult;
      }
    }
    return out;
  } catch {
    return {};
  }
}

export function saveCachedTestResults(
  incoming: CatalogTestResult[],
): Record<string, CatalogTestResult> {
  const existing = loadCachedTestResults();
  for (const item of incoming) {
    existing[item.id] = item;
  }
  const entries = Object.entries(existing);
  let trimmed = existing;
  if (entries.length > MAX_CACHED_TEST_RESULTS) {
    entries.sort(([, a], [, b]) => (b.testedAt || 0) - (a.testedAt || 0));
    trimmed = {};
    for (let i = 0; i < MAX_CACHED_TEST_RESULTS; i++) {
      const entry = entries[i];
      if (entry) {
        trimmed[entry[0]] = entry[1];
      }
    }
  }
  if (typeof window !== "undefined" && window.localStorage) {
    try {
      window.localStorage.setItem(TEST_RESULTS_STORAGE_KEY, JSON.stringify(trimmed));
    } catch {
      // Ignore storage quota errors.
    }
  }
  return trimmed;
}

export function clearCachedTestResults(): void {
  if (typeof window !== "undefined" && window.localStorage) {
    try {
      window.localStorage.removeItem(TEST_RESULTS_STORAGE_KEY);
    } catch {
      // Ignore storage errors.
    }
  }
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

export function normalizeComboRows(combos: Combo[] | undefined | null): CatalogComboRow[] {
  if (!Array.isArray(combos)) {
    return [];
  }
  return combos.map((combo) => {
    const models =
      Array.isArray(combo.steps) && combo.steps.length > 0
        ? combo.steps.map((step) => ({
            model: step.kind === "combo-ref" ? (step.comboName ?? "") : (step.model ?? ""),
            kind: step.kind === "combo-ref" ? ("combo-ref" as const) : ("model" as const),
            label: step.label,
          }))
        : (combo.targets ?? []).map((target) => ({
            model: target,
            kind: target.includes("/") ? ("model" as const) : ("combo-ref" as const),
          }));
    return {
      id: combo.id,
      name: combo.name,
      displayName: combo.name,
      strategy: combo.strategy || "priority",
      description: combo.description || undefined,
      models,
      memberCount: models.length,
      status: combo.enabled ? "active" : "paused",
      enabled: combo.enabled,
    };
  });
}

export function computeVisibilityCounts(models: CatalogModelRow[]): {
  visible: number;
  hidden: number;
  total: number;
} {
  const hidden = models.filter((m) => m.hidden === true).length;
  return {
    visible: models.length - hidden,
    hidden,
    total: models.length,
  };
}

export function getEnabledCapabilities(model: CatalogModelRow): string[] {
  if (!model.capabilities) {
    return [];
  }
  return Object.entries(model.capabilities)
    .filter(([, enabled]) => enabled === true)
    .map(([cap]) => cap.toLowerCase())
    .sort();
}

export function filterModels(
  models: CatalogModelRow[],
  filters: ModelCatalogFilters,
  healthMatrix: Record<string, ProviderHealthState>,
  testResults: Record<string, CatalogTestResult>,
): CatalogModelRow[] {
  const q = filters.query.trim().toLowerCase();

  return models.filter((m) => {
    if (q) {
      const hay =
        `${m.id} ${m.name} ${m.providerId} ${m.provider} ${m.type} ${m.subtype ?? ""}`.toLowerCase();
      if (!hay.includes(q)) {
        return false;
      }
    }
    if (filters.providerId !== "all" && m.providerId !== filters.providerId) {
      return false;
    }
    if (filters.type !== "all" && m.type !== filters.type) {
      return false;
    }
    if (filters.subtype !== "all" && (m.subtype ?? "") !== filters.subtype) {
      return false;
    }
    if (filters.capability !== "all") {
      const caps = getEnabledCapabilities(m);
      if (!caps.includes(filters.capability.toLowerCase())) {
        return false;
      }
    }
    if (filters.pricing === "free" && !m.free) {
      return false;
    }
    if (filters.pricing === "paid" && m.free) {
      return false;
    }
    if (filters.visibility === "visible" && m.hidden === true) {
      return false;
    }
    if (filters.visibility === "hidden" && m.hidden !== true) {
      return false;
    }
    if (filters.providerHealth !== "all") {
      const state = healthMatrix[m.providerId] ?? "untested";
      if (state !== filters.providerHealth) {
        return false;
      }
    }
    if (filters.testResult !== "all") {
      const res = testResults[modelTestResultKey(m.providerId, m.id)];
      if (filters.testResult === "untested") {
        if (res) {
          return false;
        }
      } else if (!res || res.status !== filters.testResult) {
        return false;
      }
    }
    if (
      typeof filters.minContextLength === "number" &&
      Number.isFinite(filters.minContextLength) &&
      (m.context_length ?? 0) < filters.minContextLength
    ) {
      return false;
    }
    if (
      typeof filters.minMaxOutputTokens === "number" &&
      Number.isFinite(filters.minMaxOutputTokens) &&
      (m.max_output_tokens ?? 0) < filters.minMaxOutputTokens
    ) {
      return false;
    }
    return true;
  });
}

function testResultSortRank(res: CatalogTestResult | undefined): number {
  if (!res) return 3;
  if (res.status === "ok") return 0;
  if (res.status === "slow") return 1;
  return 2;
}

export function sortModels(
  models: CatalogModelRow[],
  field: ModelSortField,
  direction: SortDirection,
  testResults: Record<string, CatalogTestResult>,
): CatalogModelRow[] {
  const factor = direction === "asc" ? 1 : -1;
  return [...models].sort((a, b) => {
    let cmp = 0;
    switch (field) {
      case "id":
        cmp = a.id.localeCompare(b.id);
        break;
      case "provider":
        cmp = a.provider.localeCompare(b.provider) || a.id.localeCompare(b.id);
        break;
      case "type":
        cmp = a.type.localeCompare(b.type) || a.id.localeCompare(b.id);
        break;
      case "context_length":
        cmp = (a.context_length ?? 0) - (b.context_length ?? 0);
        break;
      case "max_output_tokens":
        cmp = (a.max_output_tokens ?? 0) - (b.max_output_tokens ?? 0);
        break;
      case "testResult": {
        const ra = testResults[modelTestResultKey(a.providerId, a.id)];
        const rb = testResults[modelTestResultKey(b.providerId, b.id)];
        cmp = testResultSortRank(ra) - testResultSortRank(rb);
        if (cmp === 0) {
          cmp =
            (ra?.latencyMs ?? Number.MAX_SAFE_INTEGER) - (rb?.latencyMs ?? Number.MAX_SAFE_INTEGER);
        }
        break;
      }
    }
    return cmp * factor;
  });
}

export function filterCombos(
  combos: CatalogComboRow[],
  filters: ComboCatalogFilters,
  testResults: Record<string, CatalogTestResult>,
): CatalogComboRow[] {
  const q = filters.query.trim().toLowerCase();
  return combos.filter((c) => {
    if (q) {
      const membersStr = c.models.map((m) => m.model).join(" ");
      const hay = `${c.name} ${c.description ?? ""} ${c.strategy} ${membersStr}`.toLowerCase();
      if (!hay.includes(q)) {
        return false;
      }
    }
    if (filters.strategy !== "all" && c.strategy !== filters.strategy) {
      return false;
    }
    if (filters.status !== "all" && c.status !== filters.status) {
      return false;
    }
    if (
      typeof filters.minMembers === "number" &&
      Number.isFinite(filters.minMembers) &&
      c.memberCount < filters.minMembers
    ) {
      return false;
    }
    if (
      typeof filters.maxMembers === "number" &&
      Number.isFinite(filters.maxMembers) &&
      c.memberCount > filters.maxMembers
    ) {
      return false;
    }
    if (filters.testResult !== "all") {
      const res = testResults[comboTestResultKey(c.name)];
      if (filters.testResult === "untested") {
        if (res) {
          return false;
        }
      } else if (!res || res.status !== filters.testResult) {
        return false;
      }
    }
    return true;
  });
}

export function sortCombos(
  combos: CatalogComboRow[],
  field: ComboSortField,
  direction: SortDirection,
  testResults: Record<string, CatalogTestResult>,
): CatalogComboRow[] {
  const factor = direction === "asc" ? 1 : -1;
  return [...combos].sort((a, b) => {
    let cmp = 0;
    switch (field) {
      case "name":
        cmp = a.name.localeCompare(b.name);
        break;
      case "strategy":
        cmp = a.strategy.localeCompare(b.strategy) || a.name.localeCompare(b.name);
        break;
      case "status":
        cmp = a.status.localeCompare(b.status) || a.name.localeCompare(b.name);
        break;
      case "memberCount":
        cmp = a.memberCount - b.memberCount;
        break;
      case "testResult": {
        const ra = testResults[comboTestResultKey(a.name)];
        const rb = testResults[comboTestResultKey(b.name)];
        cmp = testResultSortRank(ra) - testResultSortRank(rb);
        break;
      }
    }
    return cmp * factor;
  });
}

export function patternMatchesIdentifier(pattern: string, identifier: string): boolean {
  const p = pattern.trim();
  const id = identifier.trim();
  if (!p || !id) return false;
  if (p === "*" || p === id) return true;
  if (p.endsWith("/*")) {
    const prefix = p.slice(0, -1);
    return id.startsWith(prefix);
  }
  return false;
}

export function keyAllowsIdentifier(
  key: ApiKeyAccessItem,
  qualifiedId: string,
  bareId?: string,
): "all" | "direct" | "pattern" | "none" {
  if (key.modelAccessMode === "all") {
    return "all";
  }
  for (const rule of key.allowedModels) {
    if (rule === qualifiedId || (bareId && rule === bareId)) {
      return "direct";
    }
  }
  for (const rule of key.allowedModels) {
    if (
      patternMatchesIdentifier(rule, qualifiedId) ||
      (bareId && patternMatchesIdentifier(rule, bareId))
    ) {
      return "pattern";
    }
  }
  return "none";
}

export function countKeyAccess(
  keys: ApiKeyAccessItem[],
  qualifiedId: string,
  bareId?: string,
): { allowed: number; total: number } {
  let allowed = 0;
  for (const key of keys) {
    if (keyAllowsIdentifier(key, qualifiedId, bareId) !== "none") {
      allowed++;
    }
  }
  return { allowed, total: keys.length };
}

export function computeUpdatedAllowedModels(
  key: ApiKeyAccessItem,
  targets: string[],
  mode: "add" | "remove",
): string[] {
  const cleanTargets = targets.map((t) => t.trim()).filter(Boolean);
  if (mode === "add") {
    // Switching from wildcard/all or adding to restricted list
    const base = key.modelAccessMode === "all" ? [] : key.allowedModels.filter((m) => m !== "*");
    const merged = new Set(base);
    for (const t of cleanTargets) {
      merged.add(t);
    }
    return [...merged];
  }
  // Remove mode
  if (key.modelAccessMode === "all") {
    return key.allowedModels;
  }
  const removeSet = new Set(cleanTargets);
  return key.allowedModels.filter((m) => !removeSet.has(m));
}

export function formatTokenCount(count: number | undefined): string {
  if (typeof count !== "number" || count <= 0) {
    return "—";
  }
  if (count >= 1_000_000) {
    const val = count / 1_000_000;
    return `${Number.isInteger(val) ? val : val.toFixed(1)}M`;
  }
  if (count >= 1_000) {
    const val = Math.round(count / 1_000);
    return `${val}K`;
  }
  return String(count);
}

export function formatRelativeTime(timestampMs: number, nowMs: number, isId: boolean): string {
  const diffSec = Math.max(0, Math.floor((nowMs - timestampMs) / 1000));
  if (diffSec < 60) {
    return isId ? "Baru saja" : "Just now";
  }
  const mins = Math.floor(diffSec / 60);
  if (mins < 60) {
    return isId ? `${mins} mnt lalu` : `${mins}m ago`;
  }
  const hours = Math.floor(mins / 60);
  return isId ? `${hours} jam lalu` : `${hours}h ago`;
}
