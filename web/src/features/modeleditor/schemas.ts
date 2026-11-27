import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ModelOverride = components["schemas"]["ModelOverride"];
export type ModelOverrideList = components["schemas"]["ModelOverrideList"];
export type CreateModelOverrideRequest = components["schemas"]["CreateModelOverrideRequest"];
export type UpdateModelOverrideRequest = components["schemas"]["UpdateModelOverrideRequest"];
export type ModelEditorPreviewRequest = components["schemas"]["ModelEditorPreviewRequest"];
export type ModelEditorPreviewResult = components["schemas"]["ModelEditorPreviewResult"];

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
  temperature: optionalBoundedNumber(0, 2, "Temperature"),
  max_tokens: optionalBoundedNumber(1, 2000000, "Max tokens"),
  extra_params_json: z
    .string()
    .optional()
    .refine((v) => parseExtraParamsJson(v ?? "").ok, {
      message: "Extra params must be a valid JSON object",
    }),
  strip_params_csv: z.string().optional(),
});

export type CreateModelOverrideFormValues = z.infer<typeof createModelOverrideSchema>;
