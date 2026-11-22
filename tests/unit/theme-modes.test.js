/**
 * Theme switching stays on the fork's three modes with no orphan classes.
 *
 * Ember and Dark keep `.dark` on <html> (components render from dark-tuned
 * tokens); `.ember` adds the ember overrides; Light clears both. An unknown
 * persisted value must degrade to the ember default rather than a class
 * combination the CSS never defines, which would silently render a broken theme.
 *
 * The suite runs in the node environment, so window/document/localStorage are
 * stubbed before the store module loads — zustand-persist reads localStorage at
 * import time. fetch is stubbed so the server sync stays offline.
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";

let classes;
let store;

function stubDom() {
  classes = new Set();
  const classList = {
    add: (c) => classes.add(c),
    remove: (c) => classes.delete(c),
    toggle: (c, on) => (on ? classes.add(c) : classes.delete(c)),
    contains: (c) => classes.has(c),
  };
  const mem = {};
  vi.stubGlobal("window", globalThis);
  vi.stubGlobal("document", {
    documentElement: { classList, style: { setProperty: () => {} } },
    querySelector: () => null,
  });
  vi.stubGlobal("localStorage", {
    getItem: (k) => (k in mem ? mem[k] : null),
    setItem: (k, v) => { mem[k] = String(v); },
    removeItem: (k) => { delete mem[k]; },
  });
  vi.stubGlobal("fetch", vi.fn(async () => ({ ok: false })));
}

describe("theme mode contract", () => {
  beforeEach(async () => {
    vi.resetModules();
    stubDom();
    const mod = await import("../../src/store/themeStore.js");
    store = mod.default;
    globalThis.__THEME_IDS__ = mod.THEME_IDS;
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("only knows Ember, Dark and Light", () => {
    expect(globalThis.__THEME_IDS__).toEqual(["ember", "dark", "light"]);
  });

  it("applies ember as .dark.ember, dark as plain .dark, light as neither", () => {
    const setTheme = store.getState().setTheme;

    setTheme("ember");
    expect(classes.has("dark")).toBe(true);
    expect(classes.has("ember")).toBe(true);

    setTheme("dark");
    expect(classes.has("dark")).toBe(true);
    expect(classes.has("ember")).toBe(false);

    setTheme("light");
    expect(classes.has("dark")).toBe(false);
    expect(classes.has("ember")).toBe(false);
  });

  it("degrades an unknown persisted value to the ember default", () => {
    const setTheme = store.getState().setTheme;

    setTheme("glass");
    expect(store.getState().theme).toBe("ember");
    expect(classes.has("ember")).toBe(true);
    expect(classes.has("dark")).toBe(true);
  });

  it("cycles ember → dark → light → ember starting from the ember default", () => {
    const { toggleTheme } = store.getState();

    expect(store.getState().theme).toBe("ember");

    toggleTheme();
    expect(store.getState().theme).toBe("dark");
    expect(classes.has("ember")).toBe(false);

    toggleTheme();
    expect(store.getState().theme).toBe("light");
    expect(classes.has("dark")).toBe(false);

    toggleTheme();
    expect(store.getState().theme).toBe("ember");
    expect(classes.has("ember")).toBe(true);
  });

  it("re-applies the persisted mode on init (page reload path)", async () => {
    store.getState().setTheme("ember");

    // Simulate reload: classes wiped, store hydrated from persistence.
    classes.clear();
    await store.getState().initTheme();

    expect(classes.has("dark")).toBe(true);
    expect(classes.has("ember")).toBe(true);
  });
});
