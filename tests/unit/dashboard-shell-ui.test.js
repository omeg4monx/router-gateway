import { describe, it, expect } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const root = resolve(__dirname, "../..");
const read = (rel) => readFileSync(resolve(root, rel), "utf8");

describe("dashboard shell CSS primitives", () => {
  const css = read("src/app/globals.css");

  it("defines Tactile Omegamon Studio motion and surface primitives", () => {
    expect(css).toContain("--ease-spring:");
    expect(css).toContain(".shell-rail");
    expect(css).toContain(".shell-popover");
    expect(css).toContain(".accordion-grid");
    expect(css).toContain('.accordion-grid[data-open="true"]');
    expect(css).toContain(".accordion-inner");
    expect(css).toContain(".popover-in");
    expect(css).toContain(".flyout-in");
    expect(css).toContain(".backdrop-fade-in");
    expect(css).toContain(".rail-icon");
  });

  it("guards shell animations under prefers-reduced-motion: reduce", () => {
    expect(css).toMatch(/@media\s*\(prefers-reduced-motion:\s*reduce\)/);
  });
});

describe("LeftRail tactile sidebar", () => {
  const rail = read("src/shared/components/LeftRail.js");

  it("renders floating tactile rail with spring width transition and unclipped outer aside", () => {
    expect(rail).toContain("shell-rail");
    expect(rail).toContain("rounded-2xl");
    expect(rail).toContain("transition-[width]");
    expect(rail).toContain("w-[228px]");
    expect(rail).toContain("w-[60px]");
    expect(rail).toContain("overflow-visible");
    expect(rail).not.toMatch(/open \? "w-\[228px\] overflow-hidden"/);
  });

  it("uses animated accordion-grid with inert and aria-hidden for collapsed inline groups", () => {
    expect(rail).toContain("accordion-grid");
    expect(rail).toContain("accordion-inner");
    expect(rail).toMatch(/data-open=\{expanded \? "true" : "false"\}/);
    expect(rail).toContain("aria-hidden={!expanded}");
    expect(rail).toMatch(/inert=\{!expanded \|\| undefined\}/);
  });

  it("renders a compact floating tooltip that clears on scroll, open toggle, and route change", () => {
    expect(rail).toContain('role="tooltip"');
    expect(rail).toContain("shell-popover");
    expect(rail).toContain("onScroll");
    expect(rail).toContain("setTooltip(null)");
    expect(rail).toMatch(/title=\{open \? label : undefined\}/);
  });

  it("uses filled active icons and tactile compact flyouts", () => {
    expect(rail).toContain("rail-icon");
    expect(rail).toContain("fill-1");
    expect(rail).toContain("flyout-in");
    expect(rail).toContain("left-[78px]");
  });

  it("mounts rail-docked NotificationBell and preserves FlyoutLink and HeaderMenu admin gate", () => {
    expect(rail).toContain('placement="right-bottom"');
    expect(rail).toContain("NotificationBell");
    expect(rail).toContain("export function FlyoutLink");
    expect(rail).toMatch(/!isApiKeyUser \? \(\s*<HeaderMenu/);
  });
});

describe("NotificationBell placement and tactile popover", () => {
  const bell = read("src/shared/components/NotificationBell.js");

  it("supports right-bottom placement while preserving default bottom-right", () => {
    expect(bell).toMatch(/placement\s*=\s*"bottom-right"/);
    expect(bell).toContain('"right-bottom"');
    expect(bell).toContain("left-full");
    expect(bell).toContain("bottom-0");
    expect(bell).toContain("shell-popover");
    expect(bell).toContain("popover-in");
    expect(bell).toContain("if (!isAdmin) return null");
    expect(bell).toMatch(/unread/i);
  });
});

describe("HeaderMenu tactile popover and keyboard handling", () => {
  const menu = read("src/shared/components/HeaderMenu.js");

  it("uses shell-popover, Escape dismiss, destructive divider, and uniform item padding", () => {
    expect(menu).toContain("shell-popover");
    expect(menu).toContain("popover-in");
    expect(menu).toMatch(/e\.key\s*===\s*"Escape"/);
    expect(menu).toContain("border-t border-border-subtle/80");
    const paddingMatches = menu.match(/px-2\.5 py-2/g) || [];
    expect(paddingMatches.length).toBeGreaterThanOrEqual(2);
  });
});
