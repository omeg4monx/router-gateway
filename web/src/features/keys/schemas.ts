import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ApiKeySummary = components["schemas"]["ApiKeySummary"];
export type CreateApiKeyRequest = components["schemas"]["CreateApiKeyRequest"];
export type UpdateApiKeyRequest = components["schemas"]["UpdateApiKeyRequest"];
export type ApiKeySecretResponse = components["schemas"]["ApiKeySecretResponse"];

export const MAX_KEY_NAME_RUNES = 128;
export const MAX_KEY_SLICE_ITEMS = 256;
export const MAX_KEY_SLICE_ITEM_RUNES = 128;

export function parseCommaSeparatedList(input: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];

  for (const part of input.split(",")) {
    const trimmed = part.trim();
    if (!trimmed || seen.has(trimmed)) {
      continue;
    }
    seen.add(trimmed);
    out.push(trimmed);
  }

  return out;
}

const commaListSchema = (label: string) =>
  z
    .string()
    .transform(parseCommaSeparatedList)
    .refine(
      (items) => items.length <= MAX_KEY_SLICE_ITEMS,
      `${label} must not exceed ${MAX_KEY_SLICE_ITEMS} items.`,
    )
    .refine(
      (items) => items.every((item) => Array.from(item).length <= MAX_KEY_SLICE_ITEM_RUNES),
      `${label} entries must not exceed ${MAX_KEY_SLICE_ITEM_RUNES} characters.`,
    );

export const createKeySchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, "Key name is required.")
    .refine(
      (value) => Array.from(value).length <= MAX_KEY_NAME_RUNES,
      `Key name must not exceed ${MAX_KEY_NAME_RUNES} characters.`,
    ),
  allowed_models: commaListSchema("Allowed models"),
  tags: commaListSchema("Tags"),
  rate_limit_rpm: z.string().transform((raw, ctx) => {
    const trimmed = raw.trim();
    if (trimmed === "") {
      return 0;
    }
    if (!/^\d+$/.test(trimmed)) {
      ctx.addIssue({
        code: "custom",
        message: "Rate limit RPM must be a non-negative whole number.",
      });
      return z.NEVER;
    }
    const parsed = Number.parseInt(trimmed, 10);
    if (!Number.isSafeInteger(parsed) || parsed < 0) {
      ctx.addIssue({
        code: "custom",
        message: "Rate limit RPM must be a non-negative whole number.",
      });
      return z.NEVER;
    }
    return parsed;
  }),
  quota_usd: z.string().transform((raw, ctx) => {
    const trimmed = raw.trim();
    if (trimmed === "") {
      return 0;
    }
    if (!/^\d+(\.\d+)?$/.test(trimmed)) {
      ctx.addIssue({
        code: "custom",
        message: "Quota USD must be a non-negative number.",
      });
      return z.NEVER;
    }
    const parsed = Number(trimmed);
    if (!Number.isFinite(parsed) || parsed < 0) {
      ctx.addIssue({
        code: "custom",
        message: "Quota USD must be a non-negative number.",
      });
      return z.NEVER;
    }
    return parsed;
  }),
  expires_at: z.string().transform((raw, ctx) => {
    const trimmed = raw.trim();
    if (trimmed === "") {
      return null;
    }
    const parsedDate = new Date(trimmed);
    if (Number.isNaN(parsedDate.getTime())) {
      ctx.addIssue({
        code: "custom",
        message: "Expiration must be a valid date and time.",
      });
      return z.NEVER;
    }
    return parsedDate.toISOString();
  }),
});

export type CreateKeyFormInput = z.input<typeof createKeySchema>;
export type CreateKeyFormOutput = z.output<typeof createKeySchema>;

export const DEFAULT_CREATE_KEY_VALUES: CreateKeyFormInput = {
  name: "",
  allowed_models: "",
  tags: "",
  rate_limit_rpm: "0",
  quota_usd: "0",
  expires_at: "",
};

export function toCreateApiKeyRequest(values: CreateKeyFormOutput): CreateApiKeyRequest {
  return {
    name: values.name,
    allowed_models: values.allowed_models,
    scopes: ["gateway"],
    tags: values.tags,
    rate_limit_rpm: values.rate_limit_rpm,
    quota_usd: values.quota_usd,
    expires_at: values.expires_at,
  };
}
