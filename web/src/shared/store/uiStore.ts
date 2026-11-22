"use client";

import { create } from "zustand";

export type ThemeMode = "dark" | "light";

const THEME_STORAGE_KEY = "omeg4mon-theme";

function applyThemeToDocument(theme: ThemeMode): void {
  if (typeof document === "undefined") {
    return;
  }
  document.documentElement.setAttribute("data-theme", theme);
  try {
    window.localStorage.setItem(THEME_STORAGE_KEY, theme);
  } catch {
    // Ignore storage errors in restricted browsing contexts.
  }
}

export interface UIState {
  theme: ThemeMode;
  mobileNavOpen: boolean;
  sidebarCollapsed: boolean;
  hydrateTheme: () => void;
  setTheme: (theme: ThemeMode) => void;
  toggleTheme: () => void;
  setMobileNavOpen: (open: boolean) => void;
  closeMobileNav: () => void;
  toggleSidebarCollapsed: () => void;
}

export const useUIStore = create<UIState>((set, get) => ({
  theme: "dark",
  mobileNavOpen: false,
  sidebarCollapsed: false,
  hydrateTheme: () => {
    if (typeof document === "undefined") {
      return;
    }
    const attr = document.documentElement.getAttribute("data-theme");
    if (attr === "light" || attr === "dark") {
      set({ theme: attr });
    }
  },
  setTheme: (theme) => {
    applyThemeToDocument(theme);
    set({ theme });
  },
  toggleTheme: () => {
    const next: ThemeMode = get().theme === "dark" ? "light" : "dark";
    applyThemeToDocument(next);
    set({ theme: next });
  },
  setMobileNavOpen: (open) => set({ mobileNavOpen: open }),
  closeMobileNav: () => set({ mobileNavOpen: false }),
  toggleSidebarCollapsed: () =>
    set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
}));
