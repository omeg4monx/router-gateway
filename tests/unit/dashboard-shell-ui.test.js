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

  it("applies shell-rail backdrop-filter on ::before so position:fixed flyouts use viewport coordinates", () => {
    const railRule = css.match(/\.shell-rail\s*\{([^}]*)\}/)?.[1] || "";
    const darkRailRule = css.match(/\.dark\s+\.shell-rail\s*\{([^}]*)\}/)?.[1] || "";
    const railBeforeRule = css.match(/\.shell-rail::before\s*\{([^}]*)\}/)?.[1] || "";
    const darkRailBeforeRule = css.match(/\.dark\s+\.shell-rail::before\s*\{([^}]*)\}/)?.[1] || "";

    expect(railRule).not.toContain("backdrop-filter");
    expect(darkRailRule).not.toContain("backdrop-filter");
    expect(railBeforeRule).toContain("backdrop-filter");
    expect(railBeforeRule).toContain("pointer-events: none");
    expect(railBeforeRule).toContain("z-index: -1");
    expect(darkRailBeforeRule).toContain("background-color");
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

describe("DashboardLayout framed workspace shell", () => {
  const layout = read("src/shared/components/layouts/DashboardLayout.js");

  it("frames the desktop workspace column and upgrades toasts to shell-popover cards", () => {
    expect(layout).toContain("lg:rounded-2xl");
    expect(layout).toContain("lg:my-2.5");
    expect(layout).toContain("lg:mr-2.5");
    expect(layout).toContain("lg:overflow-hidden");
    expect(layout).toContain("shell-popover");
    expect(layout).toContain("rounded-xl");
    expect(layout).toContain("showBell={false}");
  });

  it("keeps <main> and the desktop column wrapper free of modal-trapping stacking contexts", () => {
    const mainMatch = layout.match(/<main\s+className="([^"]*)"/);
    expect(mainMatch).not.toBeNull();
    const mainClasses = mainMatch[1];
    expect(mainClasses).not.toContain("isolate");
    expect(mainClasses).not.toContain("transform");
    expect(mainClasses).not.toContain("filter");
    expect(mainClasses).not.toContain("backdrop-blur");

    const colMatch = layout.match(/\{\/\*\s*Content column\s*\*\/\}\s*<div\s+className="([^"]*)"/);
    expect(colMatch).not.toBeNull();
    const colClasses = colMatch[1];
    expect(colClasses).not.toContain("backdrop-blur");
    expect(colClasses).not.toContain("transform");
    expect(colClasses).not.toContain("isolate");
  });

  it("preserves TaskDock mounted once in the content column", () => {
    const dockMatches = layout.match(/<TaskDock\s*\/>/g) || [];
    expect(dockMatches).toHaveLength(1);
  });
});

describe("PageHeading frosted header and icon chip", () => {
  const heading = read("src/shared/components/PageHeading.js");

  it("uses frosted backdrop-blur-md header, bottom ember-rule, and rounded icon chip", () => {
    expect(heading).toContain("backdrop-blur-md");
    expect(heading).toContain("bg-surface/65");
    expect(heading).toContain("ember-rule");
    expect(heading).toContain("rounded-xl");
    expect(heading).toContain("bg-primary/10");
    expect(heading).toContain("size-9");
  });

  it("preserves NotificationBell with optional showBell prop defaulting to true", () => {
    expect(heading).toContain("NotificationBell");
    expect(heading).toMatch(/showBell\s*=\s*true/);
    expect(heading).toContain("showBell ?");
  });
});

describe("StatusBar brand consistency and interactive controls", () => {
  const bar = read("src/shared/components/StatusBar.js");

  it("replaces legacy local_fire_department icon with LogoMark", () => {
    expect(bar).not.toContain("local_fire_department");
    expect(bar).toContain("LogoMark");
    expect(bar).toContain('className="size-3.5"');
  });

  it("exposes interactive health re-check button with latency measurement", () => {
    expect(bar).toContain('type="button"');
    expect(bar).toContain("performance.now()");
    expect(bar).toMatch(/latencyMs/);
    expect(bar).toMatch(/ms</);
    expect(bar).toContain("animate-spin");
  });

  it("exposes click-to-copy host button with guarded clipboard call and timer cleanup", () => {
    expect(bar).toContain("navigator.clipboard");
    expect(bar).toContain("writeText");
    expect(bar).toContain("try {");
    expect(bar).toContain("catch");
    expect(bar).toContain('"Copied"');
    expect(bar).toContain("clearTimeout");
  });
});
