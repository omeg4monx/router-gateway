import { describe, expect, test } from "bun:test";
import { createSecurityAuditSchema, parseAuditMetadataJson } from "./schemas";

describe("security schemas", () => {
  test("validates manual audit event creation with optional metadata JSON", () => {
    const parsed = createSecurityAuditSchema.safeParse({
      category: "keys",
      action: "key.rotated",
      severity: "warning",
      resource_id: "key_123",
      summary: "Rotated production API key",
      metadata_text: '{"reason":"scheduled"}',
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parseAuditMetadataJson(parsed.data.metadata_text)).toEqual({ reason: "scheduled" });
    }
  });

  test("rejects missing category, action, summary, or invalid metadata JSON", () => {
    expect(
      createSecurityAuditSchema.safeParse({
        category: "",
        action: "test",
        severity: "info",
        resource_id: "",
        summary: "ok",
        metadata_text: "{}",
      }).success,
    ).toBe(false);

    expect(
      createSecurityAuditSchema.safeParse({
        category: "security",
        action: "manual.check",
        severity: "critical",
        resource_id: "",
        summary: "Check failed",
        metadata_text: "[1, 2, 3]",
      }).success,
    ).toBe(false);

    expect(
      createSecurityAuditSchema.safeParse({
        category: "security",
        action: "manual.check",
        severity: "invalid",
        resource_id: "",
        summary: "Check failed",
        metadata_text: "",
      }).success,
    ).toBe(false);
  });
});
