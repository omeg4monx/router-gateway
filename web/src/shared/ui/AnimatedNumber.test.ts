import { describe, expect, test } from "bun:test";
import {
  AnimatedNumber,
  easeOutCubic,
  formatIntegerDigits,
  interpolateMetricValue,
} from "./AnimatedNumber";

describe("shared/ui/AnimatedNumber", () => {
  test("formats integer values as normal grouped digits without compact suffixes", () => {
    expect(formatIntegerDigits(0)).toBe("0");
    expect(formatIntegerDigits(999)).toBe("999");
    expect(formatIntegerDigits(12_345)).toBe((12_345).toLocaleString());
    expect(formatIntegerDigits(1_250_000)).toBe((1_250_000).toLocaleString());
    expect(formatIntegerDigits(1_250_000)).not.toMatch(/[KMB]$/);
    expect(formatIntegerDigits(42.6)).toBe("43");
    expect(formatIntegerDigits(Number.NaN)).toBe("0");
    expect(formatIntegerDigits(Number.POSITIVE_INFINITY)).toBe("0");
  });

  test("eases progress smoothly with clamped cubic ease-out", () => {
    expect(easeOutCubic(-0.5)).toBe(0);
    expect(easeOutCubic(0)).toBe(0);
    expect(easeOutCubic(0.5)).toBeCloseTo(0.875, 6);
    expect(easeOutCubic(1)).toBe(1);
    expect(easeOutCubic(1.5)).toBe(1);
  });

  test("interpolates metric values smoothly between previous and next targets", () => {
    expect(interpolateMetricValue(100, 200, 0)).toBe(100);
    expect(interpolateMetricValue(100, 200, 0.5)).toBeCloseTo(187.5, 6);
    expect(interpolateMetricValue(100, 200, 1)).toBe(200);
    expect(interpolateMetricValue(200, 100, 1)).toBe(100);
  });

  test("exports AnimatedNumber component", () => {
    expect(typeof AnimatedNumber).toBe("function");
  });
});
