"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import { THEME_CONFIG } from "@/shared/constants/config";

export const THEME_IDS = ["dark", "light"];
const DEFAULT_THEME = "dark";
const DEFAULT_ACCENT = "omegamon";
const DEFAULT_CUSTOM_COLOR = "#0860A8";

// Theme-native primaries; must match --color-brand-500 in globals.css (:root / .dark).
const OMEGAMON_PRIMARY = { light: "#0860A8", dark: "#0EA8CF" };

const splitSwatch = (light, dark) => `linear-gradient(135deg, ${light} 50%, ${dark} 50%)`;
const preset = (id, label, light, dark, filter) => ({ id, label, light, dark, filter, swatch: splitSwatch(light[0], dark[0]) });

// "omegamon" sets no inline overrides, so the CSS palette of the active mode drives.
// The rest are optional custom palettes, each with a [primary, hover] pair per mode tuned like
// the Omegamon ones: primary text >= 4.5:1 on /15 tints, on-primary (white / obsidian) >= 6:1.
export const ACCENT_PALETTES = [
  { id: "omegamon", label: "Omegamon Default", swatch: splitSwatch(OMEGAMON_PRIMARY.light, OMEGAMON_PRIMARY.dark) },
  preset("ember", "Ember Orange", ["#a93618", "#8a2c14"], ["#e67051", "#ea866c"], "none"),
  preset("blue", "Electric Blue", ["#1451d7", "#1145b6"], ["#6591f1", "#81a5f4"], "hue-rotate(195deg) saturate(1.8)"),
  preset("emerald", "Emerald Green", ["#046d4c", "#034b34"], ["#06a875", "#07c68a"], "hue-rotate(115deg) saturate(1.5)"),
  preset("purple", "Violet Purple", ["#7028eb", "#5e14db"], ["#a87cf3", "#bb98f6"], "hue-rotate(240deg) saturate(1.7)"),
  preset("rose", "Ruby Rose", ["#b5173a", "#951330"], ["#eb6583", "#ef8099"], "hue-rotate(330deg) saturate(1.6)"),
  preset("cyan", "Cyber Cyan", ["#06677f", "#044b5d"], ["#099fc3", "#0ab7e0"], "hue-rotate(170deg) saturate(1.8)"),
];

const ON_PRIMARY = { light: "#ffffff", dark: "#050608" };

// Unknown/legacy values ("ember" mode, the old "orange" default accent) fall back to the defaults.
const normalizeTheme = (theme) => (THEME_IDS.includes(theme) ? theme : DEFAULT_THEME);
const normalizeAccent = (accent) =>
  accent === "custom" || ACCENT_PALETTES.some((p) => p.id === accent) ? accent : DEFAULT_ACCENT;

export function darkenColor(hex, percent = 15) {
  let num = parseInt(String(hex).replace("#", ""), 16);
  if (isNaN(num)) return hex;
  let r = Math.max(0, (num >> 16) - Math.round(255 * (percent / 100)));
  let g = Math.max(0, ((num >> 8) & 0x00FF) - Math.round(255 * (percent / 100)));
  let b = Math.max(0, (num & 0x0000FF) - Math.round(255 * (percent / 100)));
  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
}

export function hexToHueRotate(hex) {
  const rgb = parseInt(String(hex).replace("#", ""), 16);
  if (isNaN(rgb)) return "none";
  const r = (rgb >> 16) / 255;
  const g = ((rgb >> 8) & 255) / 255;
  const b = (rgb & 255) / 255;
  const max = Math.max(r, g, b), min = Math.min(r, g, b);
  let h = 0;
  if (max !== min) {
    const d = max - min;
    switch (max) {
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h = Math.round(h * 60);
  }
  const delta = (h - 12 + 360) % 360;
  return `hue-rotate(${delta}deg) saturate(1.4)`;
}

const useThemeStore = create(
  persist(
    (set, get) => ({
      theme: DEFAULT_THEME,
      accent: DEFAULT_ACCENT,
      customColor: DEFAULT_CUSTOM_COLOR,

      setTheme: (theme) => {
        const id = normalizeTheme(theme);
        set({ theme: id });
        applyTheme(id, get().accent, get().customColor);
        syncThemeToServer(id, get().accent, get().customColor);
      },

      setAccent: (accent) => {
        const valid = normalizeAccent(accent);
        set({ accent: valid });
        applyTheme(get().theme, valid, get().customColor);
        syncThemeToServer(get().theme, valid, get().customColor);
      },

      setCustomColor: (color) => {
        if (!color || typeof color !== "string") return;
        const normalized = color.startsWith("#") ? color : `#${color}`;
        set({ accent: "custom", customColor: normalized });
        applyTheme(get().theme, "custom", normalized);
        syncThemeToServer(get().theme, "custom", normalized);
      },

      toggleTheme: () => {
        const next = get().theme === "light" ? "dark" : "light";
        set({ theme: next });
        applyTheme(next, get().accent, get().customColor);
        syncThemeToServer(next, get().accent, get().customColor);
      },

      initTheme: async () => {
        // Apply current client store immediately; normalizing migrates legacy persisted values.
        set({ theme: normalizeTheme(get().theme), accent: normalizeAccent(get().accent) });
        applyTheme(get().theme, get().accent, get().customColor);

        // Sync with global theme from server so public portals and new devices inherit instance theme
        try {
          const res = await fetch("/api/theme", { cache: "no-store" });
          if (res.ok) {
            const data = await res.json();
            if (data?.ok) {
              const hasClientLocal = typeof window !== "undefined" && localStorage.getItem(THEME_CONFIG.storageKey);
              if (!hasClientLocal || data.accent === "custom" || data.theme) {
                set({
                  theme: normalizeTheme(data.theme || get().theme),
                  accent: normalizeAccent(data.accent || get().accent),
                  customColor: data.customColor || get().customColor,
                });
                applyTheme(get().theme, get().accent, get().customColor);
              }
            }
          }
        } catch {}
      },
    }),
    {
      name: THEME_CONFIG.storageKey,
    },
  ),
);

// Sync theme to server for global public portals
function syncThemeToServer(theme, accent, customColor) {
  if (typeof window === "undefined") return;
  fetch("/api/theme", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ theme, accent, customColor }),
  }).catch(() => {});
}

// Pairs of [primary, hover] custom properties written inline by non-default accents.
const ACCENT_PROPS = ["--color-brand-500", "--color-brand-600", "--color-primary", "--color-primary-hover", "--brand-500", "--brand-600"];
const setAccentProps = (root, primary, hover, onPrimary) => {
  ACCENT_PROPS.forEach((prop, i) => root.style.setProperty(prop, i % 2 ? hover : primary));
  root.style.setProperty("--color-on-primary", onPrimary);
};

// Relative luminance (WCAG); picks whichever of white/obsidian reads better on a user-chosen fill.
function luminance(hex) {
  const n = parseInt(String(hex).replace("#", ""), 16);
  const lin = (v) => (v <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4);
  return 0.2126 * lin(((n >> 16) & 255) / 255) + 0.7152 * lin(((n >> 8) & 255) / 255) + 0.0722 * lin((n & 255) / 255);
}
const onColorFor = (hex) => {
  const fill = luminance(hex) + 0.05;
  return 1.05 / fill >= fill / (luminance(ON_PRIMARY.dark) + 0.05) ? ON_PRIMARY.light : ON_PRIMARY.dark;
};

// Apply theme & accent palette to document
function applyTheme(theme, accent = DEFAULT_ACCENT, customColor = DEFAULT_CUSTOM_COLOR) {
  if (typeof window === "undefined") return;

  const root = document.documentElement;
  const mode = normalizeTheme(theme);
  root.classList.toggle("dark", mode === "dark");

  let primary;
  if (accent === "custom" && customColor) {
    primary = customColor;
    setAccentProps(root, customColor, darkenColor(customColor, 15), onColorFor(customColor));
    root.style.setProperty("--logo-filter", hexToHueRotate(customColor));
  } else {
    const palette = ACCENT_PALETTES.find((p) => p.id === accent) || ACCENT_PALETTES[0];
    if (palette.light) {
      const [main, hover] = palette[mode];
      primary = main;
      setAccentProps(root, main, hover, ON_PRIMARY[mode]);
      root.style.setProperty("--logo-filter", palette.filter || "none");
    } else {
      // Theme-native palette: drop inline overrides so :root / .dark tokens apply.
      primary = OMEGAMON_PRIMARY[mode];
      [...ACCENT_PROPS, "--color-on-primary"].forEach((prop) => root.style.removeProperty(prop));
      root.style.setProperty("--logo-filter", hexToHueRotate(primary));
    }
  }

  const metaThemeColor = document.querySelector('meta[name="theme-color"]');
  if (metaThemeColor) metaThemeColor.setAttribute("content", primary);
}

export default useThemeStore;
