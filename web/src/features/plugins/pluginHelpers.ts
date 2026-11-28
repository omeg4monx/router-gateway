import type { CatalogModelRow } from "@/features/models/types";
import type { ProvidersCatalogResponse } from "@/features/providers/schemas";
import type { UILocale } from "@/shared/store/uiStore";
import type { PluginState, UpdatePluginRequest } from "./hooks";

export interface PluginCatalogModelOption {
  id: string;
  providerId: string;
  providerLabel: string;
}

export function buildPluginCatalogModelOptions(
  modelRows: readonly CatalogModelRow[] | undefined,
  providersCatalog: ProvidersCatalogResponse | undefined,
): PluginCatalogModelOption[] {
  const displayNames = new Map<string, string>();
  for (const adapter of providersCatalog?.adapters ?? []) {
    if (adapter.id) {
      displayNames.set(adapter.id, adapter.display_name || adapter.id);
    }
  }
  for (const node of providersCatalog?.custom_nodes ?? []) {
    if (node.id) {
      displayNames.set(node.id, node.name || node.id);
    }
  }

  const seen = new Set<string>();
  const out: PluginCatalogModelOption[] = [];

  const push = (providerId: string, modelId: string, fallbackLabel?: string) => {
    const pid = providerId.trim();
    const mid = modelId.trim();
    if (!mid) return;
    const normalizedModel = pid && mid.startsWith(`${pid}/`) ? mid.slice(pid.length + 1) : mid;
    if (!normalizedModel) return;
    const qualifiedId = pid ? `${pid}/${normalizedModel}` : normalizedModel;
    if (seen.has(qualifiedId)) return;
    seen.add(qualifiedId);
    out.push({
      id: qualifiedId,
      providerId: pid,
      providerLabel: displayNames.get(pid) || fallbackLabel?.trim() || pid,
    });
  };

  for (const row of modelRows ?? []) {
    push(row.providerId, row.id, row.provider);
  }
  for (const adapter of providersCatalog?.adapters ?? []) {
    for (const modelId of adapter.default_models ?? []) {
      push(adapter.id, modelId, adapter.display_name);
    }
  }
  for (const node of providersCatalog?.custom_nodes ?? []) {
    for (const modelId of node.models ?? []) {
      push(node.id, modelId, node.name);
    }
  }

  return out;
}

export function filterPluginCatalogModelOptions(
  options: readonly PluginCatalogModelOption[],
  query: string,
): PluginCatalogModelOption[] {
  const trimmed = query.trim().toLowerCase();
  if (!trimmed) return [...options];
  const compact = trimmed.replace(/\s+/g, "");
  return options.filter((opt) => {
    const idLower = opt.id.toLowerCase();
    const pidLower = opt.providerId.toLowerCase();
    const labelLower = opt.providerLabel.toLowerCase();
    const labelCompact = labelLower.replace(/\s+/g, "");
    return (
      idLower.includes(trimmed) ||
      pidLower.includes(trimmed) ||
      labelLower.includes(trimmed) ||
      idLower.includes(compact) ||
      pidLower.includes(compact) ||
      labelCompact.includes(compact)
    );
  });
}

export interface PluginDraft {
  modelsText: string;
  level: string;
  effort: string;
  budgetTokens: string;
  maxTokens: string;
  forceJson: boolean;
  contextWindow: string;
  minTurns: string;
  parallelToolCalls: boolean;
}

export function parseModelAllowlist(raw: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const part of raw.split(",")) {
    const trimmed = part.trim();
    if (!trimmed || seen.has(trimmed)) continue;
    seen.add(trimmed);
    out.push(trimmed);
  }
  return out;
}

export function formatModelAllowlist(models?: readonly string[]): string {
  return (models ?? []).join(", ");
}

export function toggleModelInAllowlist(modelsText: string, modelId: string): string {
  const target = modelId.trim();
  if (!target) return modelsText;
  const current = parseModelAllowlist(modelsText);
  const exists = current.includes(target);
  const next = exists ? current.filter((item) => item !== target) : [...current, target];
  return formatModelAllowlist(next);
}

function readConfigNumber(config: Record<string, unknown>, key: string, fallback: number): number {
  const raw = config[key];
  if (typeof raw === "number" && Number.isFinite(raw)) return Math.trunc(raw);
  if (typeof raw === "string") {
    const parsed = Number.parseInt(raw, 10);
    if (Number.isFinite(parsed)) return parsed;
  }
  return fallback;
}

function parsePositiveInt(raw: string, fallback: number, min = 0): number {
  const parsed = Number.parseInt(raw.trim(), 10);
  if (!Number.isFinite(parsed) || parsed < min) return fallback;
  return parsed;
}

export function normalizeAntiSlopLevel(raw: unknown): string {
  if (typeof raw !== "string") return "full";
  switch (raw.trim().toLowerCase()) {
    case "lite":
    case "light":
      return "lite";
    case "standard":
      return "standard";
    case "ultra":
      return "ultra";
    default:
      return "full";
  }
}

export function getPluginLevel(plugin: PluginState): string {
  const config = (plugin.config ?? {}) as Record<string, unknown>;
  return normalizeAntiSlopLevel(config.level);
}

export function buildPluginDraft(plugin: PluginState): PluginDraft {
  const config = (plugin.config ?? {}) as Record<string, unknown>;
  const rawEffort = typeof config.effort === "string" ? config.effort.trim().toLowerCase() : "high";
  const effort =
    rawEffort === "low" || rawEffort === "medium" || rawEffort === "high" ? rawEffort : "high";

  return {
    modelsText: formatModelAllowlist(plugin.models),
    level: getPluginLevel(plugin),
    effort,
    budgetTokens: String(readConfigNumber(config, "budget_tokens", 10240)),
    maxTokens: String(readConfigNumber(config, "max_tokens", 0)),
    forceJson: Boolean(config.force),
    contextWindow: String(readConfigNumber(config, "context_window", 128000)),
    minTurns: String(readConfigNumber(config, "min_turns", 4)),
    parallelToolCalls:
      typeof config.parallel_tool_calls === "boolean" ? config.parallel_tool_calls : true,
  };
}

export function isPluginDraftDirty(plugin: PluginState, draft: PluginDraft): boolean {
  const currentModels = plugin.models ?? [];
  const nextModels = parseModelAllowlist(draft.modelsText);
  if (
    currentModels.length !== nextModels.length ||
    currentModels.some((model, index) => model !== nextModels[index])
  ) {
    return true;
  }

  const baseline = buildPluginDraft(plugin);
  switch (plugin.key) {
    case "antiSlop":
      return draft.level !== baseline.level;
    case "thinkDeeper":
      return draft.effort !== baseline.effort || draft.budgetTokens !== baseline.budgetTokens;
    case "speedMode":
      return draft.maxTokens !== baseline.maxTokens;
    case "jsonGuard":
      return draft.forceJson !== baseline.forceJson;
    case "contextSqueezer":
      return draft.contextWindow !== baseline.contextWindow || draft.minTurns !== baseline.minTurns;
    case "openaiToolBridge":
      return draft.parallelToolCalls !== baseline.parallelToolCalls;
    default:
      return false;
  }
}

export function buildToggleEnabledPayload(
  plugin: PluginState,
  nextEnabled: boolean,
): UpdatePluginRequest {
  return {
    enabled: nextEnabled,
    models: [...(plugin.models ?? [])],
    config: { ...(plugin.config ?? {}) },
  };
}

export function buildPluginUpdatePayload(
  plugin: PluginState,
  draft: PluginDraft,
  overrideEnabled?: boolean,
): UpdatePluginRequest {
  const config: Record<string, unknown> = { ...(plugin.config ?? {}) };
  switch (plugin.key) {
    case "antiSlop":
      config.level = normalizeAntiSlopLevel(draft.level);
      break;
    case "thinkDeeper":
      config.effort = draft.effort;
      config.budget_tokens = parsePositiveInt(draft.budgetTokens, 10240, 1024);
      break;
    case "speedMode":
      config.max_tokens = parsePositiveInt(draft.maxTokens, 0, 0);
      break;
    case "jsonGuard":
      config.force = draft.forceJson;
      break;
    case "contextSqueezer":
      config.context_window = parsePositiveInt(draft.contextWindow, 128000, 1024);
      config.min_turns = parsePositiveInt(draft.minTurns, 4, 1);
      break;
    case "openaiToolBridge":
      config.parallel_tool_calls = draft.parallelToolCalls;
      break;
  }
  return {
    enabled: overrideEnabled ?? plugin.enabled,
    models: parseModelAllowlist(draft.modelsText),
    config,
  };
}

export function summarizePluginScope(plugin: PluginState, locale: UILocale): string {
  const modelCount = plugin.models?.length ?? 0;
  const modelsLabel =
    modelCount === 0
      ? locale === "id"
        ? "Semua model"
        : "All models"
      : locale === "id"
        ? `${modelCount} model`
        : modelCount === 1
          ? "1 model"
          : `${modelCount} models`;

  const config = (plugin.config ?? {}) as Record<string, unknown>;
  switch (plugin.key) {
    case "antiSlop":
      return `${modelsLabel} · ${getPluginLevel(plugin)}`;
    case "thinkDeeper": {
      const effort = typeof config.effort === "string" && config.effort ? config.effort : "high";
      return `${modelsLabel} · ${effort}`;
    }
    case "speedMode": {
      const mt = readConfigNumber(config, "max_tokens", 0);
      return mt > 0 ? `${modelsLabel} · max ${mt}` : modelsLabel;
    }
    case "jsonGuard":
      return config.force ? `${modelsLabel} · force` : modelsLabel;
    case "contextSqueezer": {
      const win = readConfigNumber(config, "context_window", 128000);
      const compactWin = win >= 1000 ? `${Math.round(win / 1000)}k` : String(win);
      return `${modelsLabel} · ${compactWin} ctx`;
    }
    default:
      return modelsLabel;
  }
}

export function parsePreviewPayload(
  raw: string,
): { ok: true; payload: Record<string, unknown> } | { ok: false } {
  try {
    const parsed: unknown = JSON.parse(raw);
    if (parsed !== null && typeof parsed === "object" && !Array.isArray(parsed)) {
      return { ok: true, payload: parsed as Record<string, unknown> };
    }
    return { ok: false };
  } catch {
    return { ok: false };
  }
}
