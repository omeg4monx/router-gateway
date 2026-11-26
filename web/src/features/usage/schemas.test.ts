import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  DEFAULT_UPSERT_PRICING_VALUES,
  formatCompactTokens,
  formatDurationMs,
  formatErrorRate,
  formatUSD,
  toUpsertModelPricingRequest,
  upsertPricingSchema,
} from "./schemas";

describe("upsertPricingSchema", () => {
  test("validates model and converts non-negative finite USD rates", async () => {
    const input = {
      model: "  openai/gpt-4o  ",
      input_per_million: "2.5",
      output_per_million: "10",
      cached_per_million: "1.25",
    };

    const parsed = upsertPricingSchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data).toEqual({
        model: "openai/gpt-4o",
        input_per_million: 2.5,
        output_per_million: 10,
        cached_per_million: 1.25,
      });

      expect(toUpsertModelPricingRequest(parsed.data)).toEqual({
        input_per_million: 2.5,
        output_per_million: 10,
        cached_per_million: 1.25,
      });
    }

    const resolver = zodResolver(upsertPricingSchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
    expect(resolved.values).toEqual({
      model: "openai/gpt-4o",
      input_per_million: 2.5,
      output_per_million: 10,
      cached_per_million: 1.25,
    });
  });

  test("defaults blank cached_per_million to zero and accepts zero rates", () => {
    const parsed = upsertPricingSchema.safeParse({
      model: "local/llama-3",
      input_per_million: "0",
      output_per_million: "0",
      cached_per_million: "",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data).toEqual({
        model: "local/llama-3",
        input_per_million: 0,
        output_per_million: 0,
        cached_per_million: 0,
      });
    }
  });

  test("rejects blank, wildcard, control-char, or oversized model identifiers", () => {
    for (const badModel of ["", "   ", "*", "gpt\n4o", "gpt\t4o", "m".repeat(129)]) {
      const res = upsertPricingSchema.safeParse({
        ...DEFAULT_UPSERT_PRICING_VALUES,
        model: badModel,
        input_per_million: "1",
        output_per_million: "2",
      });
      expect(res.success).toBe(false);
    }
  });

  test("rejects missing, negative, non-numeric, or non-finite USD rates", () => {
    const badRates = ["", "  ", "-0.01", "-5", "NaN", "Infinity", "abc", "1e999"];
    for (const rate of badRates) {
      const badInput = upsertPricingSchema.safeParse({
        model: "gpt-4o",
        input_per_million: rate,
        output_per_million: "5",
        cached_per_million: "0",
      });
      expect(badInput.success).toBe(false);

      const badOutput = upsertPricingSchema.safeParse({
        model: "gpt-4o",
        input_per_million: "5",
        output_per_million: rate,
        cached_per_million: "0",
      });
      expect(badOutput.success).toBe(false);
    }

    for (const rate of ["-1", "NaN", "Infinity", "xyz", "1e999"]) {
      const badCached = upsertPricingSchema.safeParse({
        model: "gpt-4o",
        input_per_million: "2.5",
        output_per_million: "10",
        cached_per_million: rate,
      });
      expect(badCached.success).toBe(false);
    }
  });
});

describe("formatting helpers", () => {
  test("formatUSD renders zero, sub-cent precision, standard amounts, and guards non-finite inputs", () => {
    expect(formatUSD(0)).toBe("$0.00");
    expect(formatUSD(0.0042)).toBe("$0.0042");
    expect(formatUSD(0.000042)).toBe("$0.000042");
    expect(formatUSD(2.5)).toBe("$2.50");
    expect(formatUSD(1234.56)).toBe("$1,234.56");
    expect(formatUSD(Number.NaN)).toBe("$0.00");
    expect(formatUSD(Number.POSITIVE_INFINITY)).toBe("$0.00");
  });

  test("formatCompactTokens formats counts across unit scales cleanly", () => {
    expect(formatCompactTokens(0)).toBe("0");
    expect(formatCompactTokens(-5)).toBe("0");
    expect(formatCompactTokens(Number.NaN)).toBe("0");
    expect(formatCompactTokens(842)).toBe("842");
    expect(formatCompactTokens(1_250)).toBe("1.3K");
    expect(formatCompactTokens(12_000)).toBe("12K");
    expect(formatCompactTokens(1_500_000)).toBe("1.5M");
    expect(formatCompactTokens(2_400_000_000)).toBe("2.4B");
  });

  test("formatDurationMs renders ms under 1s and seconds at or above 1s", () => {
    expect(formatDurationMs(0)).toBe("0ms");
    expect(formatDurationMs(-10)).toBe("0ms");
    expect(formatDurationMs(Number.NaN)).toBe("0ms");
    expect(formatDurationMs(244.6)).toBe("245ms");
    expect(formatDurationMs(1_000)).toBe("1.00s");
    expect(formatDurationMs(1_250)).toBe("1.25s");
  });

  test("formatErrorRate converts 0..1 fraction to percentage", () => {
    expect(formatErrorRate(0)).toBe("0.0%");
    expect(formatErrorRate(0.0525)).toBe("5.3%");
    expect(formatErrorRate(1)).toBe("100.0%");
    expect(formatErrorRate(Number.NaN)).toBe("0.0%");
  });
});
