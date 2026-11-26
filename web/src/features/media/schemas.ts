import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type MediaCapability = components["schemas"]["MediaCapability"];
export type MediaDriver = components["schemas"]["MediaDriver"];
export type MediaProvider = components["schemas"]["MediaProvider"];
export type MediaProviderListResponse = components["schemas"]["MediaProviderListResponse"];
export type CreateMediaProviderRequest = components["schemas"]["CreateMediaProviderRequest"];
export type UpdateMediaProviderRequest = components["schemas"]["UpdateMediaProviderRequest"];
export type MediaCombo = components["schemas"]["MediaCombo"];
export type MediaComboListResponse = components["schemas"]["MediaComboListResponse"];
export type CreateMediaComboRequest = components["schemas"]["CreateMediaComboRequest"];
export type UpdateMediaComboRequest = components["schemas"]["UpdateMediaComboRequest"];

export const MEDIA_CAPABILITIES = ["audio", "image", "video", "search", "fetch"] as const;
export const MEDIA_DRIVERS = ["openai_compatible", "generic_http", "builtin_fetch"] as const;

export function parseTargetsList(raw: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const part of raw.split(/[\n,]+/)) {
    const trimmed = part.trim();
    if (!trimmed || seen.has(trimmed)) {
      continue;
    }
    seen.add(trimmed);
    out.push(trimmed);
  }
  return out;
}

export const createMediaProviderSchema = z
  .object({
    name: z.string().trim().min(1, "Name is required").max(128, "Name must be at most 128 chars"),
    capability: z.enum(MEDIA_CAPABILITIES),
    driver: z.enum(MEDIA_DRIVERS),
    base_url: z.string().trim().optional().default(""),
    api_key: z.string().trim().optional().default(""),
    default_model: z.string().trim().optional().default(""),
    priority: z.coerce.number().int().min(0).max(10000).default(100),
    enabled: z.boolean().default(true),
  })
  .superRefine((val, ctx) => {
    if (val.driver === "builtin_fetch" && val.capability !== "fetch") {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["driver"],
        message: "builtin_fetch driver is only valid for the fetch capability",
      });
    }
    if (val.driver !== "builtin_fetch") {
      if (!val.base_url) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["base_url"],
          message: "Base URL is required for remote media drivers",
        });
        return;
      }
      try {
        const u = new URL(val.base_url);
        if (u.protocol !== "http:" && u.protocol !== "https:") {
          ctx.addIssue({
            code: z.ZodIssueCode.custom,
            path: ["base_url"],
            message: "Base URL must start with http:// or https://",
          });
        }
      } catch {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          path: ["base_url"],
          message: "Enter a valid http:// or https:// URL",
        });
      }
    }
  });

export type CreateMediaProviderFormInput = z.input<typeof createMediaProviderSchema>;
export type CreateMediaProviderFormOutput = z.output<typeof createMediaProviderSchema>;

export const createMediaComboSchema = z
  .object({
    name: z.string().trim().min(1, "Combo name is required").max(128),
    capability: z.enum(MEDIA_CAPABILITIES),
    targets_text: z.string().trim(),
    enabled: z.boolean().default(true),
  })
  .superRefine((val, ctx) => {
    if (parseTargetsList(val.targets_text).length === 0) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        path: ["targets_text"],
        message: "Specify at least one provider ID or name",
      });
    }
  });

export type CreateMediaComboFormInput = z.input<typeof createMediaComboSchema>;
export type CreateMediaComboFormOutput = z.output<typeof createMediaComboSchema>;
