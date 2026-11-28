import { z } from "zod";
import type { components, paths } from "@/shared/api/schema";

export type ModelOverride = components["schemas"]["ModelOverride"];
export type ModelOverrideList = components["schemas"]["ModelOverrideList"];
export type CreateModelOverrideRequest = components["schemas"]["CreateModelOverrideRequest"];
export type UpdateModelOverrideRequest = components["schemas"]["UpdateModelOverrideRequest"];
export type ModelEditorPreviewRequest =
  paths["/api/model-editor/preview"]["post"]["requestBody"]["content"]["application/json"];
export type ModelEditorPreviewResult =
  paths["/api/model-editor/preview"]["post"]["responses"]["200"]["content"]["application/json"];

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
  target_model: z.string().trim().min(1, "Model is required").max(128),
  system_prompt: z.string().optional(),
  system_prompt_position: z.enum(["suffix", "dual", "prefix"]),
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

export function overrideToFormValues(override: ModelOverride): CreateModelOverrideFormValues {
  const extraRaw = override.extra_params ?? {};
  const isDual = extraRaw._prompt_mode === "dual";
  const cleanedExtra: Record<string, unknown> = {};
  for (const [key, value] of Object.entries(extraRaw)) {
    if (key === "_prompt_mode") continue;
    cleanedExtra[key] = value;
  }

  let systemPrompt = "";
  let position: "suffix" | "dual" | "prefix" = "suffix";
  if (isDual) {
    position = "dual";
    systemPrompt = override.system_prompt_suffix || override.system_prompt_prefix || "";
  } else if (override.system_prompt_prefix && !override.system_prompt_suffix) {
    position = "prefix";
    systemPrompt = override.system_prompt_prefix;
  } else {
    position = "suffix";
    systemPrompt = override.system_prompt_suffix || override.system_prompt_prefix || "";
  }

  return {
    name: override.name,
    target_model: override.target_model || "",
    system_prompt: systemPrompt,
    system_prompt_position: position,
    temperature: override.temperature !== undefined ? String(override.temperature) : "",
    max_tokens: override.max_tokens !== undefined ? String(override.max_tokens) : "",
    extra_params_json: Object.keys(cleanedExtra).length > 0 ? JSON.stringify(cleanedExtra) : "{}",
    strip_params_csv: (override.strip_params ?? []).join(", "),
  };
}

export function buildOverridePatchPayload(
  values: CreateModelOverrideFormValues,
  existing: ModelOverride,
): { ok: true; value: UpdateModelOverrideRequest } | { ok: false; error: string } {
  const extra = parseExtraParamsJson(values.extra_params_json ?? "");
  if (!extra.ok) {
    return extra;
  }

  const trimmedPrompt = values.system_prompt?.trim() ?? "";
  const isDual = values.system_prompt_position === "dual" && Boolean(trimmedPrompt);
  const extraParams = isDual ? { ...extra.value, _prompt_mode: "dual" } : extra.value;

  const tempTrimmed = values.temperature?.trim() ?? "";
  const maxTrimmed = values.max_tokens?.trim() ?? "";

  const payload: UpdateModelOverrideRequest = {
    name: values.name.trim(),
    model_pattern:
      existing.model_pattern === existing.name ? values.name.trim() : existing.model_pattern,
    target_model: values.target_model.trim(),
    system_prompt_prefix: values.system_prompt_position === "prefix" ? trimmedPrompt : "",
    system_prompt_suffix: values.system_prompt_position === "suffix" || isDual ? trimmedPrompt : "",
    extra_params: extraParams,
    strip_params: parseStripParamsList(values.strip_params_csv ?? ""),
  };

  if (tempTrimmed) {
    payload.temperature = Number(tempTrimmed);
  } else if (existing.temperature !== undefined) {
    payload.clear_temperature = true;
  }

  if (maxTrimmed) {
    payload.max_tokens = Number(maxTrimmed);
  } else if (existing.max_tokens !== undefined) {
    payload.clear_max_tokens = true;
  }

  return { ok: true, value: payload };
}
