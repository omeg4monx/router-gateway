import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type SecurityPostureReport = components["schemas"]["SecurityPostureReport"];
export type SecurityAuditEvent = components["schemas"]["SecurityAuditEvent"];
export type CreateSecurityAuditRequest = components["schemas"]["CreateSecurityAuditRequest"];
export type SecurityAuditListResponse = components["schemas"]["SecurityAuditListResponse"];
export type ClearSecurityAuditResponse = components["schemas"]["ClearSecurityAuditResponse"];

export const AUDIT_SEVERITIES = ["info", "warning", "critical"] as const;

export function parseAuditMetadataJson(raw: string): Record<string, unknown> {
  const trimmed = raw.trim();
  if (trimmed === "") {
    return {};
  }
  const parsed: unknown = JSON.parse(trimmed);
  if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
    throw new Error("Metadata must be a JSON object");
  }
  return parsed as Record<string, unknown>;
}

export const createSecurityAuditSchema = z.object({
  category: z.string().trim().min(1, "Category is required").max(64),
  action: z.string().trim().min(1, "Action is required").max(128),
  severity: z.enum(AUDIT_SEVERITIES),
  resource_id: z.string().trim().max(128),
  summary: z.string().trim().min(1, "Summary is required").max(500),
  metadata_text: z.string().refine(
    (val) => {
      try {
        parseAuditMetadataJson(val);
        return true;
      } catch {
        return false;
      }
    },
    { message: "Metadata must be a valid JSON object" },
  ),
});

export type CreateSecurityAuditFormValues = z.infer<typeof createSecurityAuditSchema>;
