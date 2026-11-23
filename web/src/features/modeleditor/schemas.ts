import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ModelOverride = components["schemas"]["ModelOverride"];
export type ModelOverrideList = components["schemas"]["ModelOverrideList"];
export type CreateModelOverrideRequest = components["schemas"]["CreateModelOverrideRequest"];
export type UpdateModelOverrideRequest = components["schemas"]["UpdateModelOverrideRequest"];
export type ModelBypassRule = components["schemas"]["ModelBypassRule"];
export type ModelBypassRuleList = components["schemas"]["ModelBypassRuleList"];
export type CreateModelBypassRuleRequest = components["schemas"]["CreateModelBypassRuleRequest"];
export type UpdateModelBypassRuleRequest = components["schemas"]["UpdateModelBypassRuleRequest"];
export type ModelEditorPreviewRequest = components["schemas"]["ModelEditorPreviewRequest"];
export type ModelEditorPreviewResult = components["schemas"]["ModelEditorPreviewResult"];

export const BYPASS_ACTIONS = [
  "rewrite_model",
  "bypass_overrides",
  "bypass_tokensaver",
  "bypass_headroom",
] as const;

export function parseStripParamsList(csv: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const part of csv.split(",")) {
    const trimmed = part.trim();
    if (!trimmed || seen.has(trimmed)) {
      continue;
    }
    seen.add(trimmed);
    out.push(trimmed);
  }
  return out;
}

export function parseExtraParamsJson(
  raw: string,
): { ok: true; value: Record<string, unknown> } | { ok: false; error: string } {
  const trimmed = raw.trim();
  if (!trimmed) {
    return { ok: true, value: {} };
  }
  try {
    const parsed: unknown = JSON.parse(trimmed);
    if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
      return { ok: false, error: "Extra params must be a JSON object" };
    }
    return { ok: true, value: parsed as Record<string, unknown> };
  } catch {
    return { ok: false, error: "Extra params must be valid JSON" };
  }
}

function optionalBoundedNumber(min: number, max: number, label: string) {
  return z
    .string()
    .optional()
    .refine(
      (val) => {
        if (!val?.trim()) return true;
        const n = Number(val);
        return Number.isFinite(n) && n >= min && n <= max;
      },
      { message: `${label} must be between ${min} and ${max}` },
    );
}

export const createModelOverrideSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(128),
  model_pattern: z.string().trim().min(1, "Model pattern is required").max(128),
  target_model: z.string().trim().max(128).optional(),
  system_prompt_prefix: z.string().optional(),
  system_prompt_suffix: z.string().optional(),
  temperature: optionalBoundedNumber(0, 2, "Temperature"),
  max_tokens: optionalBoundedNumber(1, 2000000, "Max tokens"),
  top_p: optionalBoundedNumber(0, 1, "Top-P"),
  extra_params_json: z
    .string()
    .optional()
    .refine((v) => parseExtraParamsJson(v ?? "").ok, {
      message: "Extra params must be a valid JSON object",
    }),
  strip_params_csv: z.string().optional(),
  priority: optionalBoundedNumber(0, 10000, "Priority"),
});

export type CreateModelOverrideFormValues = z.infer<typeof createModelOverrideSchema>;

export const createBypassRuleSchema = z
  .object({
    name: z.string().trim().min(1, "Name is required").max(128),
    model_pattern: z.string().trim().min(1, "Model pattern is required").max(128),
    key_id: z.string().trim().max(128).optional(),
    protocol: z.string().trim().max(64).optional(),
    header_name: z.string().trim().max(128).optional(),
    header_value: z.string().trim().max(256).optional(),
    action: z.enum(BYPASS_ACTIONS),
    target_model: z.string().trim().max(128).optional(),
    priority: optionalBoundedNumber(0, 10000, "Priority"),
  })
  .refine((data) => data.action !== "rewrite_model" || Boolean(data.target_model?.trim()), {
    path: ["target_model"],
    message: "Target model is required when action is rewrite_model",
  });

export type CreateBypassRuleFormValues = z.infer<typeof createBypassRuleSchema>;
