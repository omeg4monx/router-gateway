import { describe, expect, test } from "bun:test";
import {
  DEFAULT_PROVIDER_STRATEGY,
  getProviderStrategyMeta,
  isProviderAccountStrategy,
  PROVIDER_ACCOUNT_STRATEGY_IDS,
  PROVIDER_STRATEGY_GROUPS,
  PROVIDER_STRATEGY_OPTIONS,
} from "./providerStrategy";

describe("providerStrategy", () => {
  test("uses fill-first as the active default strategy", () => {
    expect(DEFAULT_PROVIDER_STRATEGY).toBe("fill-first");
    expect(getProviderStrategyMeta(undefined).value).toBe("fill-first");
    expect(getProviderStrategyMeta("unknown").value).toBe("fill-first");
  });

  test("exposes all 13 account-level strategies with complete bilingual metadata", () => {
    expect(PROVIDER_ACCOUNT_STRATEGY_IDS.length).toBe(13);
    expect(PROVIDER_STRATEGY_OPTIONS.length).toBe(13);
    for (const option of PROVIDER_STRATEGY_OPTIONS) {
      expect(option.labelEn.length).toBeGreaterThan(0);
      expect(option.labelId.length).toBeGreaterThan(0);
      expect(option.descEn.length).toBeGreaterThan(0);
      expect(option.descId.length).toBeGreaterThan(0);
    }
  });

  test("excludes combo-only multi-model strategies", () => {
    for (const comboOnly of ["fusion", "pipeline", "jev", "context-relay", "auto", "weighted"]) {
      expect(isProviderAccountStrategy(comboOnly)).toBe(false);
    }
  });

  test("groups every provider strategy into non-empty categories", () => {
    const groupedCount = PROVIDER_STRATEGY_GROUPS.reduce((sum, g) => sum + g.options.length, 0);
    expect(groupedCount).toBe(PROVIDER_STRATEGY_OPTIONS.length);
  });
});
