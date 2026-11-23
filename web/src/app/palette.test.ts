import { describe, expect, test } from "bun:test";
import { readdirSync, readFileSync, statSync } from "node:fs";
import { join, resolve } from "node:path";

const WEB_SRC_DIR = resolve(import.meta.dir, "..");
const GLOBALS_CSS_PATH = join(WEB_SRC_DIR, "app", "globals.css");

function extractCssBlock(css: string, selectorPattern: RegExp): string {
  const match = css.match(selectorPattern);
  if (!match || match.index === undefined) {
    throw new Error(`CSS block not found for pattern: ${selectorPattern}`);
  }
  const openBrace = css.indexOf("{", match.index);
  let depth = 0;
  for (let i = openBrace; i < css.length; i += 1) {
    if (css[i] === "{") depth += 1;
    if (css[i] === "}") {
      depth -= 1;
      if (depth === 0) {
        return css.slice(openBrace + 1, i);
      }
    }
  }
  throw new Error(`Unclosed CSS block for pattern: ${selectorPattern}`);
}

function parseCustomProps(block: string): Record<string, string> {
  const props: Record<string, string> = {};
  const regex = /(--[a-z0-9-]+)\s*:\s*([^;]+);/gi;
  for (const match of block.matchAll(regex)) {
    const key = match[1]?.trim();
    const val = match[2]?.trim();
    if (key && val) {
      props[key] = val;
    }
  }
  return props;
}

function walkSourceFiles(dir: string): string[] {
  const files: string[] = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const stat = statSync(full);
    if (stat.isDirectory()) {
      files.push(...walkSourceFiles(full));
    } else if (
      (entry.endsWith(".ts") || entry.endsWith(".tsx") || entry.endsWith(".css")) &&
      !entry.endsWith(".test.ts") &&
      !entry.endsWith(".test.tsx")
    ) {
      files.push(full);
    }
  }
  return files;
}

describe("Omegamon color palette contract (main parity)", () => {
  const css = readFileSync(GLOBALS_CSS_PATH, "utf8");

  test("dark theme implements Omegamon Zwart palette from main", () => {
    const darkBlock = extractCssBlock(css, /\[data-theme="dark"\]\s*\{/);
    const props = parseCustomProps(darkBlock);

    expect(props["--bg-canvas"]?.toLowerCase()).toBe("#050608");
    expect(props["--bg-surface"]?.toLowerCase()).toBe("#101217");
    expect(props["--bg-surface-raised"]?.toLowerCase()).toBe("#171a20");
    expect(props["--bg-input"]?.toLowerCase()).toBe("#0a0c10");
    expect(props["--border-subtle"]?.toLowerCase()).toBe("#171a20");
    expect(props["--border-strong"]?.toLowerCase()).toBe("#262a33");
    expect(props["--text-primary"]?.toLowerCase()).toBe("#e6e9ee");
    expect(props["--text-secondary"]?.toLowerCase()).toBe("#a3aab6");
    expect(props["--text-muted"]?.toLowerCase()).toBe("#7e8592");

    expect(props["--color-brand-500"]?.toLowerCase()).toBe("#eb6583");
    expect(props["--color-brand-600"]?.toLowerCase()).toBe("#ef8099");
    expect(props["--accent-primary"]?.toLowerCase()).toBe("#eb6583");
    expect(props["--accent-primary-hover"]?.toLowerCase()).toBe("#ef8099");
    expect(props["--accent-text"]?.toLowerCase()).toBe("#eb6583");
    expect(props["--accent-primary-ink"]?.toLowerCase()).toBe("#050608");
    expect(props["--color-secondary"]?.toLowerCase()).toBe("#c0c6cf");
    expect(props["--color-accent"]?.toLowerCase()).toBe("#ff2d3a");
    expect(props["--color-highlight"]?.toLowerCase()).toBe("#22d3ee");
  });

  test("light theme implements Omegamon X-Antibody palette from main as :root default", () => {
    const lightBlock = extractCssBlock(css, /:root,\s*\[data-theme="light"\]\s*\{/);
    const props = parseCustomProps(lightBlock);

    expect(props["--bg-canvas"]?.toLowerCase()).toBe("#f5f6f8");
    expect(props["--bg-surface"]?.toLowerCase()).toBe("#ffffff");
    expect(props["--bg-surface-raised"]?.toLowerCase()).toBe("#eef1f5");
    expect(props["--bg-input"]?.toLowerCase()).toBe("#ffffff");
    expect(props["--border-subtle"]?.toLowerCase()).toBe("#e8ebf0");
    expect(props["--border-strong"]?.toLowerCase()).toBe("#d9dee6");
    expect(props["--text-primary"]?.toLowerCase()).toBe("#0b0d12");
    expect(props["--text-secondary"]?.toLowerCase()).toBe("#4a5260");
    expect(props["--text-muted"]?.toLowerCase()).toBe("#5f6775");

    expect(props["--color-brand-500"]?.toLowerCase()).toBe("#b5173a");
    expect(props["--color-brand-600"]?.toLowerCase()).toBe("#951330");
    expect(props["--accent-primary"]?.toLowerCase()).toBe("#b5173a");
    expect(props["--accent-primary-hover"]?.toLowerCase()).toBe("#951330");
    expect(props["--accent-text"]?.toLowerCase()).toBe("#b5173a");
    expect(props["--accent-primary-ink"]?.toLowerCase()).toBe("#ffffff");
    expect(props["--color-secondary"]?.toLowerCase()).toBe("#e8a317");
    expect(props["--color-accent"]?.toLowerCase()).toBe("#c8102e");
    expect(props["--color-highlight"]?.toLowerCase()).toBe("#00b8e6");
  });

  test("dark theme block appears after :root light default so [data-theme='dark'] wins CSS cascade", () => {
    const lightIdx = css.search(/:root,\s*\[data-theme="light"\]\s*\{/);
    const darkIdx = css.search(/\[data-theme="dark"\]\s*\{/);

    expect(lightIdx).toBeGreaterThan(-1);
    expect(darkIdx).toBeGreaterThan(-1);
    expect(darkIdx).toBeGreaterThan(lightIdx);
  });

  test("no legacy Ember Orange (#E56A4A), warm off-white (#FDFAF6), or hardcoded #0a0a0a remain in web/src", () => {
    const violations: string[] = [];
    const forbidden = /#(e56a4a|fdfaf6|0a0a0a)\b/gi;
    for (const filePath of walkSourceFiles(WEB_SRC_DIR)) {
      const content = readFileSync(filePath, "utf8");
      const matches = content.match(forbidden);
      if (matches && matches.length > 0) {
        violations.push(`${filePath}: ${[...new Set(matches)].join(", ")}`);
      }
    }
    expect(violations).toEqual([]);
  });
});
