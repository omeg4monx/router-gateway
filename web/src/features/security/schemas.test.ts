import { describe, expect, test } from "bun:test";
import {
  AUDIT_PAGE_LIMIT,
  computeAuditPagination,
  createSecurityAuditSchema,
  parseAuditMetadataJson,
} from "./schemas";

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

  test("computes audit table pagination bounds cleanly across empty, middle, and final pages", () => {
    expect(AUDIT_PAGE_LIMIT).toBe(10);

    expect(
      computeAuditPagination({ total: 0, pageSize: AUDIT_PAGE_LIMIT, offset: 0, itemCount: 0 }),
    ).toEqual({
      currentPage: 1,
      totalPages: 1,
      showingFrom: 0,
      showingTo: 0,
      canPrev: false,
      canNext: false,
    });

    expect(
      computeAuditPagination({ total: 25, pageSize: AUDIT_PAGE_LIMIT, offset: 0, itemCount: 10 }),
    ).toEqual({
      currentPage: 1,
      totalPages: 3,
      showingFrom: 1,
      showingTo: 10,
      canPrev: false,
      canNext: true,
    });

    expect(
      computeAuditPagination({ total: 25, pageSize: AUDIT_PAGE_LIMIT, offset: 10, itemCount: 10 }),
    ).toEqual({
      currentPage: 2,
      totalPages: 3,
      showingFrom: 11,
      showingTo: 20,
      canPrev: true,
      canNext: true,
    });

    expect(
      computeAuditPagination({ total: 25, pageSize: AUDIT_PAGE_LIMIT, offset: 20, itemCount: 5 }),
    ).toEqual({
      currentPage: 3,
      totalPages: 3,
      showingFrom: 21,
      showingTo: 25,
      canPrev: true,
      canNext: false,
    });
  });
});
