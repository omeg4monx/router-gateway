import type { UILocale } from "@/shared/store/uiStore";
import type { PluginState, UpdatePluginRequest } from "./hooks";

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
