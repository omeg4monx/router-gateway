/**
 * Theme switching stays on two modes with no orphan classes.
 *
 * Dark (Omegamon Zwart) puts `.dark` on <html>; Light (Omegamon X-Antibody)
 * clears it. Unknown or legacy persisted values (the removed "ember" mode, the
 * old "orange" default accent) must degrade to the defaults rather than a class
 * or inline-override combination the CSS no longer defines.
 *
 * The suite runs in the node environment, so window/document/localStorage are
 * stubbed before the store module loads — zustand-persist reads localStorage at
 * import time. fetch is stubbed so the server sync stays offline.
 */

import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";

let classes;
let inlineProps;
let mem;
let store;

function stubDom() {
  classes = new Set();
  inlineProps = new Map();
  const classList = {
    add: (c) => classes.add(c),
    remove: (c) => classes.delete(c),
    toggle: (c, on) => (on ? classes.add(c) : classes.delete(c)),
    contains: (c) => classes.has(c),
  };
  vi.stubGlobal("window", globalThis);
  vi.stubGlobal("document", {
    documentElement: {
      classList,
      style: {
        setProperty: (k, v) => inlineProps.set(k, v),
        removeProperty: (k) => inlineProps.delete(k),
      },
    },
    querySelector: () => null,
  });
  vi.stubGlobal("localStorage", {
    getItem: (k) => (k in mem ? mem[k] : null),
    setItem: (k, v) => { mem[k] = String(v); },
    removeItem: (k) => { delete mem[k]; },
  });
  vi.stubGlobal("fetch", vi.fn(async () => ({ ok: false })));
}

async function loadStore() {
  vi.resetModules();
  stubDom();
  const mod = await import("../../src/store/themeStore.js");
  store = mod.default;
  return mod;
}

describe("theme mode contract", () => {
  beforeEach(() => {
    mem = {};
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("only knows Dark and Light", async () => {
    const mod = await loadStore();
    expect(mod.THEME_IDS).toEqual(["dark", "light"]);
  });

  it("applies dark as .dark and light as no class", async () => {
    await loadStore();
    const { setTheme } = store.getState();

    setTheme("dark");
    expect(classes.has("dark")).toBe(true);

    setTheme("light");
    expect(classes.has("dark")).toBe(false);
  });

  it("degrades an unknown or removed mode to the dark default", async () => {
    await loadStore();
    store.getState().setTheme("ember");
    expect(store.getState().theme).toBe("dark");
    expect(classes.has("dark")).toBe(true);
    expect(classes.has("ember")).toBe(false);
  });

  it("toggles dark ↔ light starting from the dark default", async () => {
    await loadStore();
    const { toggleTheme } = store.getState();
    expect(store.getState().theme).toBe("dark");

    toggleTheme();
    expect(store.getState().theme).toBe("light");
    expect(classes.has("dark")).toBe(false);

    toggleTheme();
    expect(store.getState().theme).toBe("dark");
    expect(classes.has("dark")).toBe(true);
  });

  it("default Omegamon accent leaves the CSS palette in charge (no inline primary)", async () => {
    await loadStore();
    store.getState().setAccent("blue");
    expect(inlineProps.has("--color-primary")).toBe(true);

    store.getState().setAccent("omegamon");
    expect(inlineProps.has("--color-primary")).toBe(false);
    expect(inlineProps.has("--color-brand-500")).toBe(false);
    expect(inlineProps.has("--color-on-primary")).toBe(false);
  });

  it("accent presets swap to a per-mode variant with matching on-primary text", async () => {
    const { ACCENT_PALETTES } = await loadStore();
    const blue = ACCENT_PALETTES.find((p) => p.id === "blue");
    store.getState().setAccent("blue");

    store.getState().setTheme("light");
    expect(inlineProps.get("--color-primary")).toBe(blue.light[0]);
    expect(inlineProps.get("--color-on-primary")).toBe("#ffffff");

    store.getState().setTheme("dark");
    expect(inlineProps.get("--color-primary")).toBe(blue.dark[0]);
    expect(inlineProps.get("--color-on-primary")).toBe("#050608");
  });

  it("custom colors get whichever on-primary text reads better", async () => {
    await loadStore();
    store.getState().setCustomColor("#facc15"); // bright yellow
    expect(inlineProps.get("--color-on-primary")).toBe("#050608");
    store.getState().setCustomColor("#1e3a8a"); // deep navy
    expect(inlineProps.get("--color-on-primary")).toBe("#ffffff");
  });

  it("migrates a persisted legacy ember/orange state on init", async () => {
    const { THEME_CONFIG } = await import("../../src/shared/constants/config");
    mem[THEME_CONFIG.storageKey] = JSON.stringify({
      state: { theme: "ember", accent: "orange", customColor: "#E56A4A" },
      version: 0,
    });
    await loadStore();
    expect(store.getState().theme).toBe("ember"); // hydrated legacy value, pre-migration

    await store.getState().initTheme();

    expect(store.getState().theme).toBe("dark");
    expect(store.getState().accent).toBe("omegamon");
    expect(classes.has("dark")).toBe(true);
    expect(inlineProps.has("--color-primary")).toBe(false);
  });

  it("re-applies the persisted mode on init (page reload path)", async () => {
    await loadStore();
    store.getState().setTheme("light");
    store.getState().setTheme("dark");

    // Simulate reload: classes wiped, store hydrated from persistence.
    classes.clear();
    await store.getState().initTheme();

    expect(classes.has("dark")).toBe(true);
  });
});
