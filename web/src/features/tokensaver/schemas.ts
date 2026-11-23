import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type TokenSaverConfig = components["schemas"]["TokenSaverConfig"];
export type TokenSaverConfigPatch = components["schemas"]["TokenSaverConfigPatch"];
export type TokenSaverStats = components["schemas"]["TokenSaverStats"];
export type TokenSaverPreviewRequest = components["schemas"]["TokenSaverPreviewRequest"];
export type TokenSaverPreviewResult = components["schemas"]["TokenSaverPreviewResult"];

export const TOKEN_SAVER_STRATEGIES = ["off", "whitespace", "dedupe", "aggressive"] as const;

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
});

export type TokenSaverConfigFormValues = z.infer<typeof tokenSaverConfigSchema>;

export function parsePreviewBodyJson(
  raw: string,
): { ok: true; value: Record<string, unknown> } | { ok: false; error: string } {
  const trimmed = raw.trim();
  if (!trimmed) {
    return { ok: false, error: "JSON payload is required" };
  }
  try {
    const parsed: unknown = JSON.parse(trimmed);
    if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
      return { ok: false, error: "Payload must be a JSON object" };
    }
    return { ok: true, value: parsed as Record<string, unknown> };
  } catch {
    return { ok: false, error: "Payload must be valid JSON" };
  }
}
