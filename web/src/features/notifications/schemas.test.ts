import { describe, expect, test } from "bun:test";
import { createNotificationSchema, formatIntervalMs } from "./schemas";

describe("notifications schemas", () => {
  test("validates notification creation payload", () => {
    const parsed = createNotificationSchema.safeParse({
      title: "Maintenance Window",
      message: "Upstream provider rotation scheduled at 02:00 UTC",
      severity: "warning",
      category: "operations",
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects empty title/message or unsupported severity", () => {
    expect(
      createNotificationSchema.safeParse({
        title: "   ",
        message: "Body",
        severity: "info",
        category: "system",
      }).success,
    ).toBe(false);

    expect(
      createNotificationSchema.safeParse({
        title: "Alert",
        message: "Body",
        severity: "critical",
        category: "system",
      }).success,
    ).toBe(false);
  });

  test("formats background job intervals cleanly", () => {
    expect(formatIntervalMs(500)).toBe("500ms");
    expect(formatIntervalMs(30_000)).toBe("30s");
    expect(formatIntervalMs(300_000)).toBe("5m");
    expect(formatIntervalMs(3_600_000)).toBe("1h");
    expect(formatIntervalMs(21_600_000)).toBe("6h");
  });
});
