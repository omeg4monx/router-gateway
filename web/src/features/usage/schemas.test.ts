import { describe, expect, test } from "bun:test";
import {
  computeBarPercent,
  computeSharePercent,
  computeTotalTokens,
  formatCompactTokens,
  formatDurationMs,
  formatErrorRate,
  formatUSD,
  getUsageLogStatusVariant,
  isUsageLogError,
} from "./schemas";

describe("usage data-validity helpers", () => {
  test("isUsageLogError mirrors backend error predicate (status >= 400, status 0/<200, or error_summary)", () => {
    expect(isUsageLogError({ status_code: 200, error_summary: "" })).toBe(false);
    expect(isUsageLogError({ status_code: 304, error_summary: "   " })).toBe(false);
    expect(isUsageLogError({ status_code: 200, error_summary: "stream aborted" })).toBe(true);
    expect(isUsageLogError({ status_code: 0, error_summary: "" })).toBe(true);
    expect(isUsageLogError({ status_code: 429, error_summary: "" })).toBe(true);
    expect(isUsageLogError({ status_code: 502, error_summary: "bad gateway" })).toBe(true);
  });

  test("getUsageLogStatusVariant maps status and error_summary to semantic badge variants", () => {
    expect(getUsageLogStatusVariant({ status_code: 200, error_summary: "" })).toBe("success");
    expect(getUsageLogStatusVariant({ status_code: 200, error_summary: "upstream timeout" })).toBe(
      "destructive",
    );
    expect(getUsageLogStatusVariant({ status_code: 429, error_summary: "" })).toBe("warning");
    expect(getUsageLogStatusVariant({ status_code: 400, error_summary: "bad request" })).toBe(
      "warning",
    );
    expect(getUsageLogStatusVariant({ status_code: 500, error_summary: "" })).toBe("destructive");
    expect(getUsageLogStatusVariant({ status_code: 0, error_summary: "" })).toBe("destructive");
  });

  test("computeSharePercent and computeBarPercent guard zero and non-finite totals", () => {
    expect(computeSharePercent(0, 0)).toBe(0);
    expect(computeSharePercent(5, 0)).toBe(0);
    expect(computeSharePercent(Number.NaN, 10)).toBe(0);
    expect(computeSharePercent(1, 4)).toBe(25);
    expect(computeSharePercent(15, 10)).toBe(100);

    expect(computeBarPercent(0, 100)).toBe(0);
    expect(computeBarPercent(1, 0)).toBe(0);
    expect(computeBarPercent(1, 1000)).toBe(4);
    expect(computeBarPercent(50, 100)).toBe(50);
    expect(computeBarPercent(200, 100)).toBe(100);
  });

  test("computeTotalTokens sums input, output, and cached tokens safely", () => {
    expect(computeTotalTokens({ input_tokens: 120, output_tokens: 80, cached_tokens: 50 })).toBe(
      250,
    );
    expect(
      computeTotalTokens({ input_tokens: -10, output_tokens: Number.NaN, cached_tokens: 40 }),
    ).toBe(40);
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
