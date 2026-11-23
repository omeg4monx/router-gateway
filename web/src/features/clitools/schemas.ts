import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type CLIToolID = components["schemas"]["CLIToolID"];
export type CLIToolFormat = components["schemas"]["CLIToolFormat"];
export type CLIToolDescriptor = components["schemas"]["CLIToolDescriptor"];
export type CLIToolsCatalogResponse = components["schemas"]["CLIToolsCatalogResponse"];
export type CLIToolConfigRequest = components["schemas"]["CLIToolConfigRequest"];
export type GenerateCLIToolConfigResponse = components["schemas"]["GenerateCLIToolConfigResponse"];
export type ApplyCLIToolConfigResponse = components["schemas"]["ApplyCLIToolConfigResponse"];
export type RollbackCLIToolConfigResponse = components["schemas"]["RollbackCLIToolConfigResponse"];

export const CLI_TOOL_IDS: readonly CLIToolID[] = [
  "claude-code",
  "codex-cli",
  "opencode",
  "cline",
  "roo-code",
  "droid",
] as const;

export const MAX_CLI_INPUT_RUNES = 512;

function hasControlChar(value: string): boolean {
  for (let i = 0; i < value.length; i++) {
    const code = value.charCodeAt(i);
    if (code < 0x20 || code === 0x7f) {
      return true;
    }
  }
  return false;
}

const optionalCliFieldSchema = (label: string) =>
  z
    .string()
    .optional()
    .default("")
    .transform((value) => value.trim())
    .refine((value) => !hasControlChar(value), `${label} must not contain control characters.`)
    .refine(
      (value) => Array.from(value).length <= MAX_CLI_INPUT_RUNES,
      `${label} must not exceed ${MAX_CLI_INPUT_RUNES} characters.`,
    );

export const generateCliConfigSchema = z.object({
  base_url: z
    .string()
    .trim()
    .min(1, "Gateway base URL is required.")
    .refine(
      (value) => !hasControlChar(value),
      "Gateway base URL must not contain control characters.",
    )
    .refine(
      (value) => Array.from(value).length <= MAX_CLI_INPUT_RUNES,
      `Gateway base URL must not exceed ${MAX_CLI_INPUT_RUNES} characters.`,
    )
    .transform((value, ctx) => {
      try {
        const parsed = new URL(value);
        if ((parsed.protocol !== "http:" && parsed.protocol !== "https:") || !parsed.hostname) {
          ctx.addIssue({
            code: "custom",
            message: "Gateway base URL must be a valid http:// or https:// URL.",
          });
          return z.NEVER;
        }
        return value.replace(/\/+$/, "");
      } catch {
        ctx.addIssue({
          code: "custom",
          message: "Gateway base URL must be a valid http:// or https:// URL.",
        });
        return z.NEVER;
      }
    }),
  api_key: optionalCliFieldSchema("API key"),
  model: optionalCliFieldSchema("Default model"),
});

export type GenerateCliConfigFormInput = z.input<typeof generateCliConfigSchema>;
export type GenerateCliConfigFormOutput = z.output<typeof generateCliConfigSchema>;

export const DEFAULT_GENERATE_CLI_CONFIG_VALUES: GenerateCliConfigFormInput = {
  base_url: "",
  api_key: "",
  model: "",
};

export function toCLIToolConfigRequest(values: GenerateCliConfigFormOutput): CLIToolConfigRequest {
  const request: CLIToolConfigRequest = {
    base_url: values.base_url,
  };
  if (values.api_key) {
    request.api_key = values.api_key;
  }
  if (values.model) {
    request.model = values.model;
  }
  return request;
}
