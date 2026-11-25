import { describe, expect, test } from "bun:test";
import { readdirSync, readFileSync, statSync } from "node:fs";
import { join, resolve } from "node:path";
import { DASHBOARD_ROUTE_PATHS } from "./App";
import { ALL_NAV_ITEMS } from "./shared/ui/navigation";

const WEB_SRC_DIR = resolve(import.meta.dir);

function walkTsFiles(dir: string): string[] {
  const files: string[] = [];
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry);
    const stat = statSync(full);
    if (stat.isDirectory()) {
      files.push(...walkTsFiles(full));
    } else if ((entry.endsWith(".ts") || entry.endsWith(".tsx")) && !entry.endsWith(".test.ts")) {
      files.push(full);
    }
  }
  return files;
}

describe("Vite SPA routing & zero-Next.js contract", () => {
  test("every dashboard nav item in navigation.ts has a matching route in App.tsx", () => {
    const registered = new Set<string>(
      DASHBOARD_ROUTE_PATHS.map((segment) => `/dashboard/${segment}/`),
    );
    registered.add("/dashboard/");

    const missing = ALL_NAV_ITEMS.map((item) => item.href).filter((href) => !registered.has(href));
    expect(missing).toEqual([]);
  });

  test("no next/* imports remain anywhere in web/src", () => {
    const violations: string[] = [];
    const nextImportPattern = /from\s+["']next(?:\/[^"']*)?["']/g;
    for (const filePath of walkTsFiles(WEB_SRC_DIR)) {
      const content = readFileSync(filePath, "utf8");
      if (nextImportPattern.test(content)) {
        violations.push(filePath);
      }
    }
    expect(violations).toEqual([]);
  });

  test("no legacy 'use client' directives remain anywhere in web/src", () => {
    const violations: string[] = [];
    const useClientPattern = /^\s*["']use client["'];?\s*$/m;
    for (const filePath of walkTsFiles(WEB_SRC_DIR)) {
      const content = readFileSync(filePath, "utf8");
      if (useClientPattern.test(content)) {
        violations.push(filePath);
      }
    }
    expect(violations).toEqual([]);
  });
});
