import { describe, expect, test } from "bun:test";
import { existsSync, readFileSync } from "node:fs";
import { resolve } from "node:path";
import { COMBO_STRATEGY_OPTIONS } from "@/features/combos/strategies";

const modelsDir = resolve(import.meta.dir);
const consoleSource = readFileSync(resolve(modelsDir, "ModelCatalogConsole.tsx"), "utf8");
const modalSource = readFileSync(resolve(modelsDir, "AssignKeysModal.tsx"), "utf8");
const hooksSource = readFileSync(resolve(modelsDir, "hooks.ts"), "utf8");
const pageSource = readFileSync(
  resolve(modelsDir, "../../app/(dashboard)/dashboard/models/page.tsx"),
  "utf8",
);
const combosPageSource = readFileSync(
  resolve(modelsDir, "../../app/(dashboard)/dashboard/combos/page.tsx"),
  "utf8",
);
const comboBuilderSource = readFileSync(
  resolve(modelsDir, "../combos/ComboBuilderModal.tsx"),
  "utf8",
);
const legacyConsolePath = resolve(modelsDir, "../combos/ModelsRoutingConsole.tsx");

describe("/dashboard/models & /dashboard/combos parity", () => {
  test("/dashboard/models mounts ModelCatalogConsole and /dashboard/combos mounts CombosConsole (legacy removed)", () => {
    expect(pageSource).toContain("ModelCatalogConsole");
    expect(combosPageSource).toContain("CombosConsole");
    expect(existsSync(legacyConsolePath)).toBe(false);
  });

  test("Combos supports all 22 strategies and structured model/combo-ref steps", () => {
    expect(COMBO_STRATEGY_OPTIONS).toHaveLength(22);
    expect(comboBuilderSource).toContain("combo-strategy-select");
    expect(comboBuilderSource).toContain("combo-ref");
    expect(comboBuilderSource).toContain("fallbackOnlyOnQuotaExhaustion");
    expect(comboBuilderSource).toContain("nestedComboMode");
  });

  test("ModelCatalogConsole implements OmniRoute underline tabs, visibility counts, and bulk action testids", () => {
    expect(consoleSource).toContain('id="tab-models"');
    expect(consoleSource).toContain('id="tab-combos"');
    expect(consoleSource).toContain('data-testid="catalog-visibility-counts"');
    expect(consoleSource).toContain('data-testid="assign-keys-btn"');
    expect(consoleSource).toContain('data-testid="test-selected-btn"');
    expect(consoleSource).toContain('data-testid="test-all-filtered-btn"');
    expect(consoleSource).toContain('data-testid="cancel-tests-btn"');
    expect(consoleSource).toContain('data-testid="clear-test-results-btn"');
    expect(consoleSource).toContain('data-testid="hide-filtered-btn"');
    expect(consoleSource).toContain('data-testid="show-filtered-btn"');
    expect(consoleSource).toContain("test-model-");
  });

  test("AssignKeysModal implements OmniRoute restricted-mode confirmation and per-key result badges", () => {
    expect(modalSource).toContain("switch-restricted-");
    expect(modalSource).toContain("assign-result-");
    expect(modalSource).toContain('data-testid="assign-apply-btn"');
  });

  test("hooks wire OmniRoute catalog, health-matrix, model test, and combo test endpoints", () => {
    expect(hooksSource).toContain("/api/models/catalog");
    expect(hooksSource).toContain("/api/providers/health-matrix");
    expect(hooksSource).toContain("/api/models/test");
    expect(hooksSource).toContain("/api/combos/test");
  });
});
