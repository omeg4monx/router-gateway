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
