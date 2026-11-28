import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const providersDir = resolve(import.meta.dir);
const consoleSource = readFileSync(resolve(providersDir, "ProvidersConsole.tsx"), "utf8");
const modalsSource = readFileSync(resolve(providersDir, "ProvidersModals.tsx"), "utf8");
const detailSource = readFileSync(resolve(providersDir, "ProviderDetailView.tsx"), "utf8");

describe("ProvidersConsole mobile-first layout contracts", () => {
  test("loading skeleton header stacks on mobile without overflowing fixed widths", () => {
    expect(consoleSource).not.toContain(
      '<div className="flex items-center justify-between gap-4">',
    );
    expect(consoleSource).toContain("max-w-full");
  });

  test("header action controls use full-width mobile grid and non-clipping dropdown anchor", () => {
    expect(consoleSource).toContain("grid w-full grid-cols-2 gap-2 sm:flex sm:w-auto");
    expect(consoleSource).toContain("absolute left-0 right-0 top-full z-40 mt-1.5 w-full");
    expect(consoleSource).toContain("sm:left-auto sm:right-0 sm:w-64");
  });

  test("KPI cards stack metric and sublabel on mobile before aligning inline on sm+", () => {
    expect(consoleSource).toContain("flex flex-col gap-0.5 sm:flex-row sm:items-baseline sm:gap-2");
  });

  test("scrollable filter fieldset and category tablist constrain min-w-0 to prevent horizontal page overflow", () => {
    expect(consoleSource).toContain("min-w-0 w-full");
    expect(consoleSource).toContain(
      "custom-scrollbar flex min-w-0 flex-1 items-center gap-1.5 overflow-x-auto",
    );
  });

  test("ProviderCard quick test button provides mobile-sized touch target", () => {
    expect(consoleSource).toContain(
      "h-9 w-9 cursor-pointer items-center justify-center rounded-lg border transition-colors sm:h-7 sm:w-7",
    );
  });
});

describe("ProvidersModals mobile-first layout contracts", () => {
  test("ModalShell uses bottom-sheet presentation and dynamic viewport height on mobile", () => {
    expect(modalsSource).toContain("items-end justify-center p-0 sm:items-center sm:p-6");
    expect(modalsSource).toContain(
      "max-h-[92dvh] w-full rounded-t-2xl sm:max-h-[90dvh] sm:rounded-2xl",
    );
    expect(modalsSource).toContain("h-[100dvh]");
    expect(modalsSource).toContain("env(safe-area-inset-bottom)");
  });

  test("OAuthWizardModal AWS Quick paste buttons stack on mobile", () => {
    expect(modalsSource).toContain("grid grid-cols-1 gap-2 sm:grid-cols-2");
  });

  test("modal footers stack vertically on mobile before switching to row on sm+", () => {
    expect(modalsSource).toContain(
      "flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between",
    );
  });
});

describe("ProviderDetailView mobile-first layout contracts", () => {
  test("CardHeader blocks stack vertically on mobile to prevent PATCH URL and button overflow", () => {
    expect(detailSource).toContain(
      "flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between",
    );
    expect(detailSource).toContain("break-all font-mono text-xs text-[var(--text-muted)]");
  });

  test("form action footer stacks full-width buttons on mobile", () => {
    expect(detailSource).toContain(
      "flex flex-col-reverse gap-2 pt-1 sm:flex-row sm:items-center sm:justify-end",
    );
  });
});
