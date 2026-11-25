import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type TokenSaverConfig = components["schemas"]["TokenSaverConfig"];
export type TokenSaverConfigPatch = components["schemas"]["TokenSaverConfigPatch"];
export type TokenSaverStats = components["schemas"]["TokenSaverStats"];

export const TOKEN_SAVER_STRATEGIES = ["off", "whitespace", "dedupe", "aggressive"] as const;
export type TokenSaverStrategy = (typeof TOKEN_SAVER_STRATEGIES)[number];

export const TOKEN_SAVER_LEVELS = ["lite", "full", "ultra"] as const;
export type TokenSaverLevel = (typeof TOKEN_SAVER_LEVELS)[number];

export const HISTORY_LIMIT_PRESETS = [
  { id: "tight", maxMessages: 12 },
  { id: "balanced", maxMessages: 20 },
  { id: "deep", maxMessages: 40 },
] as const;

export type HistoryPresetId = (typeof HISTORY_LIMIT_PRESETS)[number]["id"] | "custom";

export function matchHistoryPreset(maxMessages: number): HistoryPresetId {
  const matched = HISTORY_LIMIT_PRESETS.find((preset) => preset.maxMessages === maxMessages);
  return matched ? matched.id : "custom";
}

export const tokenSaverConfigSchema = z.object({
  enabled: z.boolean(),
  strategy: z.enum(TOKEN_SAVER_STRATEGIES),
  max_history_messages: z.number().int().min(4).max(500),
  max_tool_output_chars: z.number().int().min(32).max(200000),
  headroom_enabled: z.boolean(),
  prefix_depth: z.number().int().min(1).max(10),
  max_cache_entries: z.number().int().min(8).max(10000),
  cache_ttl_seconds: z.number().int().min(60).max(86400),
  inject_cache_control: z.boolean(),
  caveman_enabled: z.boolean(),
  caveman_level: z.enum(TOKEN_SAVER_LEVELS),
  ponytail_enabled: z.boolean(),
  ponytail_level: z.enum(TOKEN_SAVER_LEVELS),
  trim_history_enabled: z.boolean(),
});

export type TokenSaverConfigFormValues = z.infer<typeof tokenSaverConfigSchema>;
