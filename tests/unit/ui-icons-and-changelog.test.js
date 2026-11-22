import { describe, it, expect } from "vitest";
import { readFileSync, statSync, existsSync } from "node:fs";
import { resolve } from "node:path";
import { getProviderIconSrc } from "../../src/shared/utils/providerIcon.js";

const root = resolve(__dirname, "../..");

describe("provider icon URLs", () => {
  it.each([
    ["minimax-code", "/providers/minimax-cn.png"],
    ["minimax-code-global", "/providers/minimax.png"],
  ])("serves an existing MiniMax asset for %s", (id, expected) => {
    const pathname = new URL(getProviderIconSrc(id), "http://localhost").pathname;
    expect(pathname).toBe(expected);
    expect(statSync(resolve(root, `public${pathname}`)).isFile()).toBe(true);
  });

  it.each(["ovh", "dothoosh"])("uses the text fallback rather than requesting a missing asset for %s", (id) => {
    expect(getProviderIconSrc(id)).toBeNull();
  });
});