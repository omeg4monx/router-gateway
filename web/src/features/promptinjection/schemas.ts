import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type PromptGuardConfig = components["schemas"]["PromptGuardConfig"];
export type PromptGuardConfigPatch = components["schemas"]["PromptGuardConfigPatch"];
export type PromptGuardRule = components["schemas"]["PromptGuardRule"];
export type PromptGuardRuleList = components["schemas"]["PromptGuardRuleList"];
export type CreatePromptGuardRuleRequest = components["schemas"]["CreatePromptGuardRuleRequest"];
export type UpdatePromptGuardRuleRequest = components["schemas"]["UpdatePromptGuardRuleRequest"];
export type PromptGuardTestRequest = components["schemas"]["PromptGuardTestRequest"];
export type PromptGuardTestResult = components["schemas"]["PromptGuardTestResult"];
export type PromptGuardEvent = components["schemas"]["PromptGuardEvent"];
export type PromptGuardEventList = components["schemas"]["PromptGuardEventList"];

export const PROMPT_GUARD_MODES = ["block", "sanitize", "audit"] as const;
export const PROMPT_GUARD_SENSITIVITIES = ["low", "medium", "high"] as const;
export const PROMPT_GUARD_MATCH_TYPES = ["substring", "regex"] as const;
export const PROMPT_GUARD_STATUS_CODES = [400, 403, 422] as const;

export const promptGuardConfigSchema = z.object({
  enabled: z.boolean(),
  mode: z.enum(PROMPT_GUARD_MODES),
  sensitivity: z.enum(PROMPT_GUARD_SENSITIVITIES),
  block_status_code: z.union([z.literal(400), z.literal(403), z.literal(422)]),
});

export type PromptGuardConfigFormValues = z.infer<typeof promptGuardConfigSchema>;

export const createPromptGuardRuleSchema = z.object({
  name: z.string().trim().min(1, "Rule name is required").max(128, "Max 128 characters"),
  pattern: z.string().min(1, "Pattern is required").max(512, "Max 512 bytes"),
  match_type: z.enum(PROMPT_GUARD_MATCH_TYPES),
  action: z.enum(PROMPT_GUARD_MODES),
  severity: z.enum(PROMPT_GUARD_SENSITIVITIES),
  description: z.string().trim().max(256).optional(),
});

export type CreatePromptGuardRuleFormValues = z.infer<typeof createPromptGuardRuleSchema>;

export const promptGuardTestSchema = z.object({
  text: z.string().trim().min(1, "Enter sample prompt text to inspect"),
  mode: z.enum(PROMPT_GUARD_MODES).optional(),
  sensitivity: z.enum(PROMPT_GUARD_SENSITIVITIES).optional(),
});

export type PromptGuardTestFormValues = z.infer<typeof promptGuardTestSchema>;
