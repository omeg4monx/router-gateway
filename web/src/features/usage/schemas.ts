import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type UsageRange = "24h" | "7d" | "30d" | "all";
export type UsageLog = components["schemas"]["UsageLog"];
export type UsageLogsResponse = components["schemas"]["UsageLogsResponse"];
export type UsageSummaryTotals = components["schemas"]["UsageSummaryTotals"];
export type ModelUsageBreakdown = components["schemas"]["ModelUsageBreakdown"];
export type ProviderUsageBreakdown = components["schemas"]["ProviderUsageBreakdown"];
export type UsageSeriesBucket = components["schemas"]["UsageSeriesBucket"];
export type UsageSummaryResponse = components["schemas"]["UsageSummaryResponse"];
export type ModelPricingSource = components["schemas"]["ModelPricingSource"];
export type ModelPricing = components["schemas"]["ModelPricing"];
export type ListModelPricingResponse = components["schemas"]["ListModelPricingResponse"];
export type UpsertModelPricingRequest = components["schemas"]["UpsertModelPricingRequest"];

export const USAGE_RANGES = ["24h", "7d", "30d", "all"] as const;

export const USAGE_RANGE_LABELS: Record<UsageRange, string> = {
  "24h": "24h",
  "7d": "7d",
  "30d": "30d",
  all: "All",
};

export const USAGE_RANGE_DESCRIPTIONS: Record<UsageRange, string> = {
  "24h": "Last 24 hours",
  "7d": "Last 7 days",
  "30d": "Last 30 days",
  all: "All retained history",
};

export const MAX_MODEL_IDENTIFIER_RUNES = 128;
export const DEFAULT_LOGS_LIMIT = 50;
export const MAX_LOGS_LIMIT = 200;

const USD_NUMBER_PATTERN = /^(?:\d+(?:\.\d+)?|\.\d+)$/;

function parseNonNegativeRate(raw: string, ctx: z.RefinementCtx, label: string): number {
  const trimmed = raw.trim();
  if (!USD_NUMBER_PATTERN.test(trimmed)) {
    ctx.addIssue({
      code: "custom",
      message: `${label} must be a non-negative finite USD amount.`,
    });
    return z.NEVER;
  }
  const parsed = Number(trimmed);
  if (!Number.isFinite(parsed) || parsed < 0) {
    ctx.addIssue({
      code: "custom",
      message: `${label} must be a non-negative finite USD amount.`,
    });
    return z.NEVER;
  }
  return parsed;
}

export const upsertPricingSchema = z.object({
  model: z
    .string()
    .trim()
    .min(1, "Model identifier is required.")
    .refine((value) => value !== "*", "Model identifier must not be wildcard '*'.")
    .refine(
      (value) => !/[\r\n\t]/.test(value),
      "Model identifier must not contain control whitespace.",
    )
    .refine(
      (value) => Array.from(value).length <= MAX_MODEL_IDENTIFIER_RUNES,
      `Model identifier must not exceed ${MAX_MODEL_IDENTIFIER_RUNES} characters.`,
    ),
  input_per_million: z.string().transform((raw, ctx) => {
    const trimmed = raw.trim();
    if (trimmed === "") {
      ctx.addIssue({
        code: "custom",
        message: "Input USD / 1M tokens is required.",
      });
      return z.NEVER;
    }
    return parseNonNegativeRate(trimmed, ctx, "Input USD / 1M tokens");
  }),
  output_per_million: z.string().transform((raw, ctx) => {
    const trimmed = raw.trim();
    if (trimmed === "") {
      ctx.addIssue({
        code: "custom",
        message: "Output USD / 1M tokens is required.",
      });
      return z.NEVER;
    }
    return parseNonNegativeRate(trimmed, ctx, "Output USD / 1M tokens");
  }),
  cached_per_million: z.string().transform((raw, ctx) => {
    const trimmed = raw.trim();
    if (trimmed === "") {
      return 0;
    }
    return parseNonNegativeRate(trimmed, ctx, "Cached USD / 1M tokens");
  }),
});

export type UpsertPricingFormInput = z.input<typeof upsertPricingSchema>;
export type UpsertPricingFormOutput = z.output<typeof upsertPricingSchema>;

export const DEFAULT_UPSERT_PRICING_VALUES: UpsertPricingFormInput = {
  model: "",
  input_per_million: "",
  output_per_million: "",
  cached_per_million: "0",
};

export function toUpsertModelPricingRequest(
  values: UpsertPricingFormOutput,
): UpsertModelPricingRequest {
  return {
    input_per_million: values.input_per_million,
    output_per_million: values.output_per_million,
    cached_per_million: values.cached_per_million,
  };
}

const STANDARD_USD_FORMATTER = new Intl.NumberFormat("en-US", {
  minimumFractionDigits: 2,
  maximumFractionDigits: 2,
});

export function formatUSD(amount: number): string {
  if (!Number.isFinite(amount) || amount === 0) {
    return "$0.00";
  }
  const sign = amount < 0 ? "-" : "";
  const abs = Math.abs(amount);
  if (abs < 0.01) {
    const fourDecimals = abs.toFixed(4);
    if (Number(fourDecimals) > 0) {
      return `${sign}$${fourDecimals}`;
    }
    return `${sign}$${abs.toFixed(6)}`;
  }
  return `${sign}$${STANDARD_USD_FORMATTER.format(abs)}`;
}

export function formatCompactTokens(tokens: number): string {
  if (!Number.isFinite(tokens) || tokens <= 0) {
    return "0";
  }
  if (tokens < 1_000) {
    return Math.round(tokens).toString();
  }
  const scales: Array<[number, string]> = [
    [1_000_000_000, "B"],
    [1_000_000, "M"],
    [1_000, "K"],
  ];
  for (const [divisor, suffix] of scales) {
    if (tokens >= divisor) {
      const rounded = Math.round((tokens / divisor) * 10) / 10;
      const formatted = Number.isInteger(rounded) ? rounded.toFixed(0) : rounded.toFixed(1);
      return `${formatted}${suffix}`;
    }
  }
  return Math.round(tokens).toString();
}

export function formatDurationMs(ms: number): string {
  if (!Number.isFinite(ms) || ms <= 0) {
    return "0ms";
  }
  if (ms < 1_000) {
    return `${Math.round(ms)}ms`;
  }
  return `${(ms / 1_000).toFixed(2)}s`;
}

export function formatErrorRate(rate: number): string {
  if (!Number.isFinite(rate) || rate <= 0) {
    return "0.0%";
  }
  const clamped = Math.min(1, rate);
  return `${(clamped * 100).toFixed(1)}%`;
}
