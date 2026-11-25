export type ModelIntelligenceTier = "flagship" | "reasoning" | "balanced" | "fast";
export type ModelSortField = "name" | "intelligence";
export type ModelSortDirection = "asc" | "desc";

export interface ModelIntelligenceProfile {
  score: number;
  tier: ModelIntelligenceTier;
  source: "curated" | "heuristic";
}

export const MODEL_TIER_BADGE: Record<ModelIntelligenceTier, { label: string; className: string }> =
  {
    flagship: {
      label: "Flagship",
      className: "bg-[var(--color-primary)]/15 text-[var(--color-primary)]",
    },
    reasoning: {
      label: "Reasoning",
      className: "bg-[var(--status-warn-bg)] text-[var(--status-warn)]",
    },
    balanced: {
      label: "Balanced",
      className: "bg-[var(--color-surface-2)] text-[var(--color-text-secondary)]",
    },
    fast: {
      label: "Fast",
      className: "bg-[var(--status-ok-bg)] text-[var(--status-ok)]",
    },
  };

interface CuratedRule {
  pattern: RegExp;
  score: number;
  tier: ModelIntelligenceTier;
}

// Ordered from most specific to broader families so specialized variants
// (e.g. mini, flash-lite, thinking, pro) match before their base family.
const CURATED_MODEL_RULES: readonly CuratedRule[] = [
  // Anthropic Claude
  { pattern: /^claude-(?:opus-4[.-]6|4[.-]6-opus)/, score: 98, tier: "flagship" },
  { pattern: /^claude-(?:opus-4[.-]5|4[.-]5-opus)/, score: 97, tier: "flagship" },
  { pattern: /^claude-(?:opus-4|4[.-]?[01]?-opus)/, score: 95, tier: "flagship" },
  { pattern: /^claude-(?:sonnet-4[.-]6|4[.-]6-sonnet)/, score: 92, tier: "balanced" },
  { pattern: /^claude-(?:sonnet-4[.-]5|4[.-]5-sonnet)/, score: 90, tier: "balanced" },
  { pattern: /^claude-(?:sonnet-4|4-sonnet)/, score: 88, tier: "balanced" },
  { pattern: /^claude-3[.-]7-sonnet/, score: 86, tier: "balanced" },
  { pattern: /^claude-3[.-]5-sonnet/, score: 84, tier: "balanced" },
  { pattern: /^claude-3-opus/, score: 83, tier: "flagship" },
  { pattern: /^claude-(?:haiku-4[.-]5|4[.-]5-haiku)/, score: 77, tier: "fast" },
  { pattern: /^claude-3[.-]5-haiku/, score: 72, tier: "fast" },
  { pattern: /^claude-3-haiku/, score: 65, tier: "fast" },

  // OpenAI Reasoning & GPT
  { pattern: /^o3-pro\b/, score: 97, tier: "reasoning" },
  { pattern: /^o3-mini\b/, score: 87, tier: "reasoning" },
  { pattern: /^o3\b/, score: 94, tier: "reasoning" },
  { pattern: /^o4-mini\b/, score: 90, tier: "reasoning" },
  { pattern: /^o1-pro\b/, score: 93, tier: "reasoning" },
  { pattern: /^o1-mini\b/, score: 84, tier: "reasoning" },
  { pattern: /^o1\b/, score: 90, tier: "reasoning" },
  { pattern: /^gpt-5[.-]?\d*-codex/, score: 95, tier: "reasoning" },
  { pattern: /^gpt-5[.-]?\d*-pro/, score: 98, tier: "flagship" },
  { pattern: /^gpt-5[.-]?\d*-mini/, score: 81, tier: "fast" },
  { pattern: /^gpt-5[.-]?\d*-nano/, score: 69, tier: "fast" },
  { pattern: /^gpt-5[.-][2-9]/, score: 97, tier: "flagship" },
  { pattern: /^gpt-5\b/, score: 95, tier: "flagship" },
  { pattern: /^gpt-4[.-]5/, score: 89, tier: "flagship" },
  { pattern: /^gpt-4[.-]1-nano/, score: 65, tier: "fast" },
  { pattern: /^gpt-4[.-]1-mini/, score: 76, tier: "fast" },
  { pattern: /^gpt-4[.-]1\b/, score: 87, tier: "balanced" },
  { pattern: /^gpt-4o-mini/, score: 72, tier: "fast" },
  { pattern: /^gpt-4o\b/, score: 84, tier: "balanced" },
  { pattern: /^gpt-4-turbo/, score: 80, tier: "balanced" },
  { pattern: /^gpt-4\b/, score: 78, tier: "balanced" },
  { pattern: /^gpt-3[.-]5/, score: 58, tier: "fast" },

  // Google Gemini
  { pattern: /^gemini-3[.-]?\d*-pro/, score: 97, tier: "flagship" },
  { pattern: /^gemini-3[.-]?\d*-flash-lite/, score: 74, tier: "fast" },
  { pattern: /^gemini-3[.-]?\d*-flash/, score: 83, tier: "fast" },
  { pattern: /^gemini-2[.-]5-pro/, score: 94, tier: "flagship" },
  { pattern: /^gemini-2[.-]5-flash-thinking/, score: 86, tier: "reasoning" },
  { pattern: /^gemini-2[.-]5-flash-lite/, score: 69, tier: "fast" },
  { pattern: /^gemini-2[.-]5-flash/, score: 79, tier: "fast" },
  { pattern: /^gemini-2[.-]0-flash-thinking/, score: 83, tier: "reasoning" },
  { pattern: /^gemini-2[.-]0-pro/, score: 84, tier: "balanced" },
  { pattern: /^gemini-2[.-]0-flash/, score: 73, tier: "fast" },
  { pattern: /^gemini-1[.-]5-pro/, score: 79, tier: "balanced" },
  { pattern: /^gemini-1[.-]5-flash/, score: 68, tier: "fast" },

  // DeepSeek
  { pattern: /^deepseek-r1/, score: 92, tier: "reasoning" },
  { pattern: /^deepseek-reasoner/, score: 92, tier: "reasoning" },
  { pattern: /^deepseek-v3[.-]2/, score: 91, tier: "flagship" },
  { pattern: /^deepseek-v3/, score: 89, tier: "flagship" },
  { pattern: /^deepseek-chat/, score: 87, tier: "balanced" },

  // xAI Grok
  { pattern: /^grok-4/, score: 95, tier: "flagship" },
  { pattern: /^grok-3-mini/, score: 79, tier: "fast" },
  { pattern: /^grok-3/, score: 90, tier: "flagship" },
  { pattern: /^grok-2/, score: 80, tier: "balanced" },

  // Qwen
  { pattern: /^qwq\b/, score: 88, tier: "reasoning" },
  { pattern: /^qwen3-max/, score: 92, tier: "flagship" },
  { pattern: /^qwen3-coder/, score: 87, tier: "balanced" },
  { pattern: /^qwen3-235b/, score: 89, tier: "flagship" },
  { pattern: /^qwen-max/, score: 89, tier: "flagship" },
  { pattern: /^qwen-plus/, score: 81, tier: "balanced" },
  { pattern: /^qwen-turbo/, score: 70, tier: "fast" },
  { pattern: /^qwen-?2[.-]5-72b/, score: 82, tier: "balanced" },

  // Meta Llama & Mistral
  { pattern: /^llama-4-maverick/, score: 86, tier: "balanced" },
  { pattern: /^llama-4-scout/, score: 77, tier: "fast" },
  { pattern: /^llama-3[.-]1-405b/, score: 85, tier: "flagship" },
  { pattern: /^llama-3[.-][13]-70b/, score: 81, tier: "balanced" },
  { pattern: /^llama-3[.-]\d-8b/, score: 64, tier: "fast" },
  { pattern: /^mistral-large/, score: 83, tier: "balanced" },
  { pattern: /^mistral-medium/, score: 77, tier: "balanced" },
  { pattern: /^mistral-small|ministral/, score: 68, tier: "fast" },
  { pattern: /^codestral/, score: 80, tier: "balanced" },
  { pattern: /^kimi-k2/, score: 90, tier: "flagship" },
  { pattern: /^glm-4[.-]5/, score: 86, tier: "balanced" },
];

function normalizeModelIdentifier(raw: string): string {
  const lower = raw.trim().toLowerCase();
  const withoutTag = lower.replace(/:(?:free|beta|exacto|extended|thinking)$/i, "");
  const segments = withoutTag.split("/");
  const base = segments[segments.length - 1] ?? withoutTag;
  return base.replace(/-(?:20\d{6}|20\d{2}-\d{2}-\d{2}|latest|preview)$/i, "");
}

function estimateHeuristicProfile(
  normalizedId: string,
  modelName?: string,
): ModelIntelligenceProfile {
  const combined = `${normalizedId} ${(modelName ?? "").toLowerCase()}`;
  let score = 64;
  let explicitTier: ModelIntelligenceTier | null = null;

  if (/\b(?:thinking|reasoner|reasoning|r1|qwq|o1|o3|o4)\b/.test(combined)) {
    score += 22;
    explicitTier = "reasoning";
  } else if (/\b(?:opus|ultra|max|flagship|frontier|pro|large)\b/.test(combined)) {
    score += 18;
    explicitTier = "flagship";
  } else if (/\b(?:sonnet|plus|medium|coder|instruct)\b/.test(combined)) {
    score += 10;
    explicitTier = "balanced";
  }

  if (/\b(?:nano|tiny|micro)\b/.test(combined)) {
    score -= 18;
    explicitTier = "fast";
  } else if (/\b(?:lite|mini|flash|haiku|instant|turbo|fast|small|air)\b/.test(combined)) {
    score -= 10;
    if (explicitTier !== "reasoning") {
      explicitTier = "fast";
    }
  }

  const paramMatch = combined.match(/\b(\d+(?:\.\d+)?)b\b/);
  if (paramMatch) {
    const billions = Number.parseFloat(paramMatch[1] ?? "0");
    if (billions >= 200) {
      score += 20;
    } else if (billions >= 65) {
      score += 14;
    } else if (billions >= 24) {
      score += 7;
    } else if (billions > 0 && billions <= 9) {
      score -= 10;
      if (explicitTier !== "reasoning") {
        explicitTier = "fast";
      }
    }
  }

  const versionMatch = normalizedId.match(/(?:^|[-_ ])([2-9])(?:\.(\d))?(?:$|[-_ ])/);
  if (versionMatch) {
    const major = Number.parseInt(versionMatch[1] ?? "0", 10);
    const minor = Number.parseInt(versionMatch[2] ?? "0", 10);
    score += Math.min(6, Math.max(0, major - 2) + (minor >= 5 ? 1 : 0));
  }

  const clamped = Math.max(10, Math.min(99, Math.round(score)));
  const tier: ModelIntelligenceTier =
    explicitTier ?? (clamped >= 86 ? "flagship" : clamped >= 72 ? "balanced" : "fast");

  return {
    score: clamped,
    tier,
    source: "heuristic",
  };
}

export function getModelIntelligenceProfile(
  modelId: string,
  modelName?: string,
): ModelIntelligenceProfile {
  const normalized = normalizeModelIdentifier(modelId);
  for (const rule of CURATED_MODEL_RULES) {
    if (rule.pattern.test(normalized)) {
      return {
        score: rule.score,
        tier: rule.tier,
        source: "curated",
      };
    }
  }
  return estimateHeuristicProfile(normalized, modelName);
}

export function sortProviderModels<T extends { id: string; name?: string }>(
  models: readonly T[],
  field: ModelSortField,
  direction: ModelSortDirection,
): T[] {
  return [...models].sort((a, b) => {
    if (field === "intelligence") {
      const scoreA = getModelIntelligenceProfile(a.id, a.name).score;
      const scoreB = getModelIntelligenceProfile(b.id, b.name).score;
      if (scoreA !== scoreB) {
        return direction === "desc" ? scoreB - scoreA : scoreA - scoreB;
      }
    }
    const nameCmp = a.id.localeCompare(b.id, undefined, { numeric: true, sensitivity: "base" });
    return direction === "asc" ? nameCmp : -nameCmp;
  });
}
