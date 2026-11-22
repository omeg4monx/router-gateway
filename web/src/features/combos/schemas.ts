import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ComboStrategy = components["schemas"]["ComboStrategy"];
export type Combo = components["schemas"]["Combo"];
export type ListCombosResponse = components["schemas"]["ListCombosResponse"];
export type CreateComboRequest = components["schemas"]["CreateComboRequest"];
export type UpdateComboRequest = components["schemas"]["UpdateComboRequest"];
export type ModelAlias = components["schemas"]["ModelAlias"];
export type ListModelAliasesResponse = components["schemas"]["ListModelAliasesResponse"];
export type UpsertModelAliasRequest = components["schemas"]["UpsertModelAliasRequest"];
export type DisabledModel = components["schemas"]["DisabledModel"];
export type ListDisabledModelsResponse = components["schemas"]["ListDisabledModelsResponse"];
export type DisableModelRequest = components["schemas"]["DisableModelRequest"];
export type EnableModelRequest = components["schemas"]["EnableModelRequest"];
export type CooldownEntry = components["schemas"]["CooldownEntry"];
export type ListCooldownsResponse = components["schemas"]["ListCooldownsResponse"];

export const COMBO_STRATEGIES = ["ordered"] as const;
export const MAX_MODEL_IDENTIFIER_RUNES = 128;
export const MAX_DESCRIPTION_RUNES = 512;
export const MAX_TARGETS_COUNT = 64;

const modelIdentifierSchema = (label: string) =>
  z
    .string()
    .trim()
    .min(1, `${label} is required.`)
    .refine((value) => value !== "*", `${label} must not be wildcard '*'.`)
    .refine((value) => !/[\r\n\t]/.test(value), `${label} must not contain control whitespace.`)
    .refine(
      (value) => Array.from(value).length <= MAX_MODEL_IDENTIFIER_RUNES,
      `${label} must not exceed ${MAX_MODEL_IDENTIFIER_RUNES} characters.`,
    );

export const createComboSchema = z
  .object({
    name: modelIdentifierSchema("Combo name"),
    description: z
      .string()
      .trim()
      .refine(
        (value) => Array.from(value).length <= MAX_DESCRIPTION_RUNES,
        `Description must not exceed ${MAX_DESCRIPTION_RUNES} characters.`,
      ),
    strategy: z.enum(COMBO_STRATEGIES),
    enabled: z.boolean(),
    targets: z
      .array(modelIdentifierSchema("Target model"))
      .min(1, "Add at least one fallback target model.")
      .max(MAX_TARGETS_COUNT, `Targets must not exceed ${MAX_TARGETS_COUNT} items.`),
  })
  .superRefine((val, ctx) => {
    const seen = new Set<string>();
    val.targets.forEach((target, idx) => {
      if (target === val.name) {
        ctx.addIssue({
          code: "custom",
          path: ["targets", idx],
          message: "Combo targets must not reference the combo name itself.",
        });
      }
      if (seen.has(target)) {
        ctx.addIssue({
          code: "custom",
          path: ["targets", idx],
          message: `Duplicate target model "${target}".`,
        });
      }
      seen.add(target);
    });
  });

export type CreateComboFormInput = z.input<typeof createComboSchema>;
export type CreateComboFormOutput = z.output<typeof createComboSchema>;

export const DEFAULT_CREATE_COMBO_VALUES: CreateComboFormInput = {
  name: "",
  description: "",
  strategy: "ordered",
  enabled: true,
  targets: [],
};

export function toCreateComboRequest(values: CreateComboFormOutput): CreateComboRequest {
  return {
    name: values.name,
    description: values.description,
    strategy: values.strategy,
    enabled: values.enabled,
    targets: values.targets,
  };
}

export const upsertAliasSchema = z
  .object({
    alias: modelIdentifierSchema("Alias name"),
    target_model: modelIdentifierSchema("Target model"),
  })
  .superRefine((val, ctx) => {
    if (val.alias === val.target_model) {
      ctx.addIssue({
        code: "custom",
        path: ["target_model"],
        message: "Target model must differ from the alias name.",
      });
    }
  });

export type UpsertAliasFormInput = z.input<typeof upsertAliasSchema>;
export type UpsertAliasFormOutput = z.output<typeof upsertAliasSchema>;

export const DEFAULT_UPSERT_ALIAS_VALUES: UpsertAliasFormInput = {
  alias: "",
  target_model: "",
};

export function toUpsertAliasRequest(values: UpsertAliasFormOutput): UpsertModelAliasRequest {
  return {
    target_model: values.target_model,
  };
}

export const disableModelSchema = z.object({
  model: modelIdentifierSchema("Model identifier"),
  reason: z
    .string()
    .trim()
    .refine(
      (value) => Array.from(value).length <= MAX_DESCRIPTION_RUNES,
      `Reason must not exceed ${MAX_DESCRIPTION_RUNES} characters.`,
    ),
});

export type DisableModelFormInput = z.input<typeof disableModelSchema>;
export type DisableModelFormOutput = z.output<typeof disableModelSchema>;

export const DEFAULT_DISABLE_MODEL_VALUES: DisableModelFormInput = {
  model: "",
  reason: "",
};

export function toDisableModelRequest(values: DisableModelFormOutput): DisableModelRequest {
  return {
    model: values.model,
    reason: values.reason,
  };
}
