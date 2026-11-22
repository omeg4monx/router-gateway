import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  buildSettingsPatchRequest,
  changePasswordSchema,
  DEFAULT_SETTINGS_FORM_VALUES,
  settingsFormSchema,
  settingsItemsToFormValues,
} from "./schemas";

describe("settingsFormSchema", () => {
  test("accepts valid boundary values for all typed settings", async () => {
    const minValid = {
      ui_theme: "dark" as const,
      ui_locale: "en" as const,
      payload_logging_enabled: false,
      payload_retention_days: 1,
      metadata_retention_days: 1,
      default_timeout_seconds: 5,
      allow_response_chain_fallback: true,
      reason: "Minimum bounds check",
    };
    expect(settingsFormSchema.safeParse(minValid).success).toBe(true);

    const maxValid = {
      ui_theme: "light" as const,
      ui_locale: "id" as const,
      payload_logging_enabled: true,
      payload_retention_days: 90,
      metadata_retention_days: 3650,
      default_timeout_seconds: 600,
      allow_response_chain_fallback: false,
      reason: "",
    };
    const resolver = zodResolver(settingsFormSchema);
    const resolved = await resolver(maxValid, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
    expect(resolved.values).toEqual(maxValid);
  });

  test("enforces payload_retention_days in 1..90", () => {
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        payload_retention_days: 0,
      }).success,
    ).toBe(false);
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        payload_retention_days: 91,
      }).success,
    ).toBe(false);
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        payload_retention_days: 7.5,
      }).success,
    ).toBe(false);
  });

  test("enforces metadata_retention_days in 1..3650", () => {
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        metadata_retention_days: 0,
      }).success,
    ).toBe(false);
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        metadata_retention_days: 3651,
      }).success,
    ).toBe(false);
  });

  test("enforces default_timeout_seconds in 5..600", () => {
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        default_timeout_seconds: 4,
      }).success,
    ).toBe(false);
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        default_timeout_seconds: 601,
      }).success,
    ).toBe(false);
  });

  test("enforces reason max length of 512 runes", () => {
    // 512 surrogate-pair runes = 1024 UTF-16 code units, valid under 512-rune limit
    const maxRunes = "🚀".repeat(512);
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        reason: maxRunes,
      }).success,
    ).toBe(true);

    // 513 runes exceeds limit
    const overRunes = "🚀".repeat(513);
    expect(
      settingsFormSchema.safeParse({
        ...DEFAULT_SETTINGS_FORM_VALUES,
        reason: overRunes,
      }).success,
    ).toBe(false);
  });

  test("hydrates form values from API settings items and builds minimal patch payload", () => {
    const apiItems = [
      {
        key: "ui.theme",
        type: "enum" as const,
        value: "light",
        default_value: "dark",
        revision: 2,
        updated_at: "2026-11-23T00:00:00Z",
        description: "Dashboard color theme mode.",
        options: ["dark", "light", "system"],
      },
      {
        key: "usage.payload_retention_days",
        type: "integer" as const,
        value: 14,
        default_value: 7,
        revision: 1,
        updated_at: "2026-11-23T00:00:00Z",
        description: "Retention days.",
        min: 1,
        max: 90,
      },
    ];

    const hydrated = settingsItemsToFormValues(apiItems);
    expect(hydrated.ui_theme).toBe("light");
    expect(hydrated.payload_retention_days).toBe(14);
    expect(hydrated.default_timeout_seconds).toBe(120);

    expect(buildSettingsPatchRequest(hydrated, apiItems)).toBeNull();

    const patch = buildSettingsPatchRequest(
      {
        ...hydrated,
        ui_theme: "dark",
        default_timeout_seconds: 180,
        reason: "  Switch back to dark  ",
      },
      apiItems,
    );
    expect(patch).toEqual({
      updates: {
        "ui.theme": "dark",
        "gateway.default_timeout_seconds": 180,
      },
      reason: "Switch back to dark",
    });
  });
});

describe("changePasswordSchema", () => {
  test("accepts valid password change within 8..72 rune/UTF-8 byte parity with Go", async () => {
    const valid = {
      currentPassword: "old-password-123",
      newPassword: "p".repeat(72),
      confirmPassword: "p".repeat(72),
    };
    const resolver = zodResolver(changePasswordSchema);
    const resolved = await resolver(valid, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
  });

  test("rejects empty current password", () => {
    const res = changePasswordSchema.safeParse({
      currentPassword: "",
      newPassword: "new-valid-password",
      confirmPassword: "new-valid-password",
    });
    expect(res.success).toBe(false);
  });

  test("rejects new password under 8 runes (including multi-byte characters)", () => {
    // 7 2-byte runes = 14 UTF-8 bytes, but < 8 runes
    const sevenRunes = "\u00e9".repeat(7);
    const res = changePasswordSchema.safeParse({
      currentPassword: "current-password",
      newPassword: sevenRunes,
      confirmPassword: sevenRunes,
    });
    expect(res.success).toBe(false);

    // 8 2-byte runes = 16 UTF-8 bytes -> valid
    const eightRunes = "\u00e9".repeat(8);
    const ok = changePasswordSchema.safeParse({
      currentPassword: "current-password",
      newPassword: eightRunes,
      confirmPassword: eightRunes,
    });
    expect(ok.success).toBe(true);
  });

  test("rejects new password exceeding 72 UTF-8 bytes even when rune count is <= 72", () => {
    // 37 2-byte runes = 74 UTF-8 bytes > 72 bcrypt limit
    const overByteLimit = "\u00e9".repeat(37);
    const res = changePasswordSchema.safeParse({
      currentPassword: "current-password",
      newPassword: overByteLimit,
      confirmPassword: overByteLimit,
    });
    expect(res.success).toBe(false);

    // 36 2-byte runes = 72 UTF-8 bytes -> exact boundary, valid
    const exactByteLimit = "\u00e9".repeat(36);
    const ok = changePasswordSchema.safeParse({
      currentPassword: "current-password",
      newPassword: exactByteLimit,
      confirmPassword: exactByteLimit,
    });
    expect(ok.success).toBe(true);
  });

  test("rejects mismatched password confirmation", () => {
    const res = changePasswordSchema.safeParse({
      currentPassword: "current-password",
      newPassword: "new-valid-password",
      confirmPassword: "different-password",
    });
    expect(res.success).toBe(false);
    if (!res.success) {
      expect(res.error.issues.some((issue) => issue.path.includes("confirmPassword"))).toBe(true);
    }
  });
});
