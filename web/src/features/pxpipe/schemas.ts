import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type PxpipeConfig = components["schemas"]["PxpipeConfig"];
export type PxpipeStatus = components["schemas"]["PxpipeStatus"];
export type UpdatePxpipeConfigRequest = components["schemas"]["UpdatePxpipeConfigRequest"];
export type PxpipeHealthCheck = components["schemas"]["PxpipeHealthCheck"];
export type PxpipeHealthReport = components["schemas"]["PxpipeHealthReport"];
export type PxpipeStats = components["schemas"]["PxpipeStats"];
export type PxpipeLogEntry = components["schemas"]["PxpipeLogEntry"];
export type PxpipeLogListResponse = components["schemas"]["PxpipeLogListResponse"];
export type PxpipePreviewRequest = components["schemas"]["PxpipePreviewRequest"];
export type PxpipePreviewResponse = components["schemas"]["PxpipePreviewResponse"];

export const pxpipeConfigSchema = z.object({
  enabled: z.boolean().default(true),
  auto_start: z.boolean().default(true),
  strip_null_fields: z.boolean().default(true),
  collapse_whitespace: z.boolean().default(true),
  dedupe_system_prompts: z.boolean().default(true),
  max_history_messages: z.coerce
    .number()
    .int("Must be an integer")
    .min(0, "Minimum is 0 (unlimited)")
    .max(500, "Maximum is 500"),
});

export type PxpipeConfigFormInput = z.input<typeof pxpipeConfigSchema>;
export type PxpipeConfigFormOutput = z.output<typeof pxpipeConfigSchema>;

export const pxpipePreviewSchema = z.object({
  payloadJson: z
    .string()
    .trim()
    .min(2, "JSON object payload is required")
    .refine((raw) => {
      try {
        const parsed: unknown = JSON.parse(raw);
        return typeof parsed === "object" && parsed !== null && !Array.isArray(parsed);
      } catch {
        return false;
      }
    }, "Must be a valid JSON object"),
});

export type PxpipePreviewFormInput = z.input<typeof pxpipePreviewSchema>;
export type PxpipePreviewFormOutput = z.output<typeof pxpipePreviewSchema>;
