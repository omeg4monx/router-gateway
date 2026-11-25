import {
  COMBO_STRATEGY_OPTIONS,
  type ComboStrategyMeta,
  STRATEGY_CATEGORIES,
  type StrategyCategory,
} from "@/features/combos/strategies";
import type { components } from "@/shared/api/schema";

export type ProviderAccountStrategy = components["schemas"]["ProviderAccountStrategy"];
export type ProviderStrategyState = components["schemas"]["ProviderStrategyState"];
export type ListProviderStrategiesResponse =
  components["schemas"]["ListProviderStrategiesResponse"];

export const DEFAULT_PROVIDER_STRATEGY: ProviderAccountStrategy = "fill-first";

export const PROVIDER_ACCOUNT_STRATEGY_IDS: readonly ProviderAccountStrategy[] = [
  "fill-first",
  "priority",
  "round-robin",
  "p2c",
  "least-used",
  "random",
  "strict-random",
  "lkgp",
  "cache-optimized",
  "headroom",
  "reset-aware",
  "reset-window",
  "expiry-first",
] as const;

export const QUICK_PROVIDER_STRATEGIES: readonly ProviderAccountStrategy[] = [
  "fill-first",
  "round-robin",
  "p2c",
  "least-used",
] as const;

const PROVIDER_STRATEGY_SET = new Set<string>(PROVIDER_ACCOUNT_STRATEGY_IDS);

export interface ProviderStrategyMeta extends Omit<ComboStrategyMeta, "value"> {
  value: ProviderAccountStrategy;
}

export const PROVIDER_STRATEGY_OPTIONS: readonly ProviderStrategyMeta[] =
  PROVIDER_ACCOUNT_STRATEGY_IDS.map((id) => {
    const meta = COMBO_STRATEGY_OPTIONS.find((opt) => opt.value === id);
    if (!meta) {
      throw new Error(`Missing strategy metadata for provider strategy "${id}"`);
    }
    return {
      ...meta,
      value: id,
    };
  });

export interface ProviderStrategyGroup {
  id: StrategyCategory;
  labelEn: string;
  labelId: string;
  options: ProviderStrategyMeta[];
}

export const PROVIDER_STRATEGY_GROUPS: readonly ProviderStrategyGroup[] = STRATEGY_CATEGORIES.map(
  (category) => ({
    ...category,
    options: PROVIDER_STRATEGY_OPTIONS.filter((opt) => opt.category === category.id),
  }),
).filter((group) => group.options.length > 0);

export function isProviderAccountStrategy(value: string): value is ProviderAccountStrategy {
  return PROVIDER_STRATEGY_SET.has(value);
}

export function getProviderStrategyMeta(strategy: string | undefined | null): ProviderStrategyMeta {
  const matched = strategy
    ? PROVIDER_STRATEGY_OPTIONS.find((opt) => opt.value === strategy)
    : undefined;
  return matched ?? (PROVIDER_STRATEGY_OPTIONS[0] as ProviderStrategyMeta);
}
