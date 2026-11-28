import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { runWithConcurrency } from "./hooks";

const providersDir = resolve(import.meta.dir);
const logoSource = readFileSync(resolve(providersDir, "ProviderLogo.tsx"), "utf8");
const consoleSource = readFileSync(resolve(providersDir, "ProvidersConsole.tsx"), "utf8");
const detailSource = readFileSync(resolve(providersDir, "ProviderDetailView.tsx"), "utf8");
const customModelModalSource = readFileSync(
  resolve(providersDir, "AddCustomModelModal.tsx"),
  "utf8",
);
const modalsSource = readFileSync(resolve(providersDir, "ProvidersModals.tsx"), "utf8");
const hooksSource = readFileSync(resolve(providersDir, "hooks.ts"), "utf8");

describe("Providers bundle & architecture optimizations", () => {
  test("ProviderLogo does not import catalogUtils or providerRegistry and lazy-loads images", () => {
    expect(logoSource).not.toContain('from "./catalogUtils"');
    expect(logoSource).not.toContain('from "./providerRegistry"');
    expect(logoSource).toContain('loading="lazy"');
    expect(logoSource).toContain('decoding="async"');
  });

  test("ProvidersModals does not ship dead ProviderConnectionsDrawer", () => {
    expect(modalsSource).not.toContain("export function ProviderConnectionsDrawer(");
  });

  test("ProvidersConsole lazy-loads modals, preserves close animation latch, memoizes cards, and defers search", () => {
    expect(consoleSource).not.toContain('} from "./ProvidersModals"');
    expect(consoleSource).toContain('import("./ProvidersModals")');
    expect(consoleSource).toContain("modalsActivated");
    expect(consoleSource).toContain("memo(function ProviderCard");
    expect(consoleSource).toContain("useDeferredValue(search)");
    expect(consoleSource).toContain("catalogQuery.isPending || connectionsQuery.isPending");
  });

  test("ProviderDetailView lazy-loads modals, defers searches, bounds batch concurrency, and avoids skeleton flash", () => {
    expect(detailSource).not.toContain('} from "./ProvidersModals"');
    expect(detailSource).not.toContain('} from "./AddCustomModelModal"');
    expect(detailSource).toContain('import("./ProvidersModals")');
    expect(detailSource).toContain('import("./AddCustomModelModal")');
    expect(detailSource).toContain("modalsActivated");
    expect(detailSource).toContain("useDeferredValue(modelSearch)");
    expect(detailSource).toContain("useDeferredValue(connectionSearch)");
    expect(detailSource).toContain("catalogQuery.isPending || connectionsQuery.isPending");
    expect(detailSource).toContain("runWithConcurrency");
    expect(detailSource).toContain("copyTimerRef");
  });

  test("AddCustomModelModal uses the same Spring Pop + Blur-In + Top Accent Sweep choreography", () => {
    expect(customModelModalSource).toContain("provider-modal-root");
    expect(customModelModalSource).toContain("provider-modal-backdrop");
    expect(customModelModalSource).toContain("provider-modal-panel");
    expect(customModelModalSource).toContain("provider-modal-accent");
  });

  test("ModalShell animates spring pop + blur-in + top accent sweep with reduced-motion support", () => {
    expect(modalsSource).toContain("MODAL_EXIT_MS");
    expect(modalsSource).toContain("requestAnimationFrame");
    expect(modalsSource).toContain("provider-modal-backdrop");
    expect(modalsSource).toContain("provider-modal-panel");
    expect(modalsSource).toContain("provider-modal-accent");
    expect(modalsSource).toContain('data-state={visible ? "open" : "closed"}');
  });

  test("hooks.ts uses typed apiClient for connection models discovery instead of raw fetch", () => {
    expect(hooksSource).not.toContain("await fetch(");
    expect(hooksSource).toContain('apiClient.GET("/api/providers/connections/{id}/models"');
  });

  test("runWithConcurrency bounds in-flight workers and preserves item order", async () => {
    let active = 0;
    let peak = 0;
    const items = [1, 2, 3, 4, 5, 6, 7, 8];

    const results = await runWithConcurrency(items, 3, async (item) => {
      active += 1;
      peak = Math.max(peak, active);
      await new Promise((r) => setTimeout(r, 5));
      active -= 1;
      return item * 10;
    });

    expect(peak).toBeLessThanOrEqual(3);
    expect(peak).toBeGreaterThan(1);
    expect(results).toEqual([10, 20, 30, 40, 50, 60, 70, 80]);
    expect(await runWithConcurrency([], 4, async () => 1)).toEqual([]);
  });
});
