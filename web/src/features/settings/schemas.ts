import { z } from "zod";
import {
  countPasswordRunes,
  MAX_PASSWORD_BYTES,
  MIN_PASSWORD_RUNES,
  utf8ByteLength,
} from "@/features/system/password";
import type { components } from "@/shared/api/schema";

export type SettingItem = components["schemas"]["SettingItem"];
export type SettingValue = components["schemas"]["SettingValue"];
export type SettingHistoryItem = components["schemas"]["SettingHistoryItem"];
export type PatchSettingsRequest = components["schemas"]["PatchSettingsRequest"];
export type RollbackSettingRequest = components["schemas"]["RollbackSettingRequest"];

export const UI_THEME_OPTIONS = ["dark", "light", "system"] as const;
export const UI_LOCALE_OPTIONS = ["en", "id"] as const;

export const settingsFormSchema = z.object({
  ui_theme: z.enum(UI_THEME_OPTIONS),
  ui_locale: z.enum(UI_LOCALE_OPTIONS),
  payload_logging_enabled: z.boolean(),
  payload_retention_days: z
    .number({ error: "Payload retention must be a whole number." })
    .int("Payload retention must be a whole number.")
    .min(1, "Payload retention must be between 1 and 90 days.")
    .max(90, "Payload retention must be between 1 and 90 days."),
  metadata_retention_days: z
    .number({ error: "Metadata retention must be a whole number." })
    .int("Metadata retention must be a whole number.")
    .min(1, "Metadata retention must be between 1 and 3650 days.")
    .max(3650, "Metadata retention must be between 1 and 3650 days."),
  default_timeout_seconds: z
    .number({ error: "Default timeout must be a whole number." })
    .int("Default timeout must be a whole number.")
    .min(5, "Default timeout must be between 5 and 600 seconds.")
    .max(600, "Default timeout must be between 5 and 600 seconds."),
  allow_response_chain_fallback: z.boolean(),
  reason: z
    .string()
    .refine((v) => Array.from(v).length <= 512, "Change reason must not exceed 512 characters."),
});

export type SettingsFormValues = z.infer<typeof settingsFormSchema>;

export const DEFAULT_SETTINGS_FORM_VALUES: SettingsFormValues = {
  ui_theme: "light",
  ui_locale: "en",
  payload_logging_enabled: false,
  payload_retention_days: 7,
  metadata_retention_days: 90,
  default_timeout_seconds: 120,
  allow_response_chain_fallback: true,
  reason: "",
};

const SETTING_FIELD_SPECS: ReadonlyArray<{
  settingKey: string;
  field: Exclude<keyof SettingsFormValues, "reason">;
}> = [
  { settingKey: "ui.theme", field: "ui_theme" },
  { settingKey: "ui.locale", field: "ui_locale" },
  { settingKey: "usage.payload_logging_enabled", field: "payload_logging_enabled" },
  { settingKey: "usage.payload_retention_days", field: "payload_retention_days" },
  { settingKey: "usage.metadata_retention_days", field: "metadata_retention_days" },
  { settingKey: "gateway.default_timeout_seconds", field: "default_timeout_seconds" },
  {
    settingKey: "gateway.allow_response_chain_fallback",
    field: "allow_response_chain_fallback",
  },
];

function isUiTheme(value: unknown): value is SettingsFormValues["ui_theme"] {
  return value === "dark" || value === "light" || value === "system";
}

function isUiLocale(value: unknown): value is SettingsFormValues["ui_locale"] {
  return value === "en" || value === "id";
}

export function settingsItemsToFormValues(items: SettingItem[]): SettingsFormValues {
  const byKey = new Map(items.map((item) => [item.key, item.value]));

  const rawTheme = byKey.get("ui.theme");
  const rawLocale = byKey.get("ui.locale");
  const rawPayloadLogging = byKey.get("usage.payload_logging_enabled");
  const rawPayloadRetention = byKey.get("usage.payload_retention_days");
  const rawMetadataRetention = byKey.get("usage.metadata_retention_days");
  const rawTimeout = byKey.get("gateway.default_timeout_seconds");
  const rawFallback = byKey.get("gateway.allow_response_chain_fallback");

  return {
    ui_theme: isUiTheme(rawTheme) ? rawTheme : DEFAULT_SETTINGS_FORM_VALUES.ui_theme,
    ui_locale: isUiLocale(rawLocale) ? rawLocale : DEFAULT_SETTINGS_FORM_VALUES.ui_locale,
    payload_logging_enabled:
      typeof rawPayloadLogging === "boolean"
        ? rawPayloadLogging
        : DEFAULT_SETTINGS_FORM_VALUES.payload_logging_enabled,
    payload_retention_days:
      typeof rawPayloadRetention === "number" && Number.isInteger(rawPayloadRetention)
        ? rawPayloadRetention
        : DEFAULT_SETTINGS_FORM_VALUES.payload_retention_days,
    metadata_retention_days:
      typeof rawMetadataRetention === "number" && Number.isInteger(rawMetadataRetention)
        ? rawMetadataRetention
        : DEFAULT_SETTINGS_FORM_VALUES.metadata_retention_days,
    default_timeout_seconds:
      typeof rawTimeout === "number" && Number.isInteger(rawTimeout)
        ? rawTimeout
        : DEFAULT_SETTINGS_FORM_VALUES.default_timeout_seconds,
    allow_response_chain_fallback:
      typeof rawFallback === "boolean"
        ? rawFallback
        : DEFAULT_SETTINGS_FORM_VALUES.allow_response_chain_fallback,
    reason: "",
  };
}

export function buildSettingsPatchRequest(
  values: SettingsFormValues,
  currentItems: SettingItem[],
): PatchSettingsRequest | null {
  const baseline = settingsItemsToFormValues(currentItems);
  const updates: Record<string, SettingValue> = {};

  for (const { settingKey, field } of SETTING_FIELD_SPECS) {
    if (values[field] !== baseline[field]) {
      updates[settingKey] = values[field];
    }
  }

  if (Object.keys(updates).length === 0) {
    return null;
  }

  const trimmedReason = values.reason.trim();
  return trimmedReason.length > 0 ? { updates, reason: trimmedReason } : { updates };
}

export const changePasswordSchema = z
  .object({
    currentPassword: z.string().min(1, "Current administrator password is required."),
    newPassword: z
      .string()
      .refine(
        (value) => countPasswordRunes(value) >= MIN_PASSWORD_RUNES,
        `New password must be at least ${MIN_PASSWORD_RUNES} characters.`,
      )
      .refine(
        (value) => utf8ByteLength(value) <= MAX_PASSWORD_BYTES,
        `New password must not exceed ${MAX_PASSWORD_BYTES} UTF-8 bytes.`,
      ),
    confirmPassword: z.string(),
  })
  .superRefine((data, ctx) => {
    if (data.newPassword !== data.confirmPassword) {
      ctx.addIssue({
        code: "custom",
        path: ["confirmPassword"],
        message: "New password confirmation does not match.",
      });
    }
  });

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;
