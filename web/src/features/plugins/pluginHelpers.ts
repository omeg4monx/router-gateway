import type { UILocale } from "@/shared/store/uiStore";
import type { PluginState, UpdatePluginRequest } from "./hooks";

export interface PluginDraft {
  modelsText: string;
  level: string;
}

export function parseModelAllowlist(raw: string): string[] {
  return raw
    .split(",")
    .map((entry) => entry.trim())
    .filter(Boolean);
}

export function formatModelAllowlist(models?: readonly string[]): string {
  return (models ?? []).join(", ");
}

export function getPluginLevel(plugin: PluginState): string {
  const config = (plugin.config ?? {}) as Record<string, unknown>;
  return typeof config.level === "string" ? config.level : "full";
}

export function buildPluginDraft(plugin: PluginState): PluginDraft {
  return {
    modelsText: formatModelAllowlist(plugin.models),
    level: getPluginLevel(plugin),
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
  if (plugin.key === "antiSlop" && draft.level !== getPluginLevel(plugin)) {
    return true;
  }
  return false;
}

export function buildPluginUpdatePayload(
  plugin: PluginState,
  draft: PluginDraft,
  overrideEnabled?: boolean,
): UpdatePluginRequest {
  const config: Record<string, unknown> = { ...(plugin.config ?? {}) };
  if (plugin.key === "antiSlop") {
    config.level = draft.level;
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

  if (plugin.key === "antiSlop") {
    return `${modelsLabel} · ${getPluginLevel(plugin)}`;
  }
  return modelsLabel;
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
