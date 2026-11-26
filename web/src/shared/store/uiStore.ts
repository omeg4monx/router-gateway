import { create } from "zustand";

export type ThemeMode = "dark" | "light";
export type UILocale = "en" | "id";

const THEME_STORAGE_KEY = "omeg4mon-theme";
const LOCALE_STORAGE_KEY = "omeg4mon-locale";

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

function applyLocaleToDocument(locale: UILocale): void {
  if (typeof document === "undefined") {
    return;
  }
  document.documentElement.setAttribute("lang", locale);
  document.documentElement.setAttribute("data-locale", locale);
  try {
    window.localStorage.setItem(LOCALE_STORAGE_KEY, locale);
  } catch {
    // Ignore storage errors in restricted browsing contexts.
  }
}

export interface UIState {
  theme: ThemeMode;
  locale: UILocale;
  mobileNavOpen: boolean;
  sidebarCollapsed: boolean;
  hydrateTheme: () => void;
  setTheme: (theme: ThemeMode) => void;
  toggleTheme: () => void;
  hydrateLocale: () => void;
  setLocale: (locale: UILocale) => void;
  toggleLocale: () => void;
  setMobileNavOpen: (open: boolean) => void;
  closeMobileNav: () => void;
  toggleSidebarCollapsed: () => void;
}

export const useUIStore = create<UIState>((set, get) => ({
  theme: "light",
  locale: "en",
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
  hydrateLocale: () => {
    if (typeof document === "undefined") {
      return;
    }
    const attr =
      document.documentElement.getAttribute("data-locale") ??
      document.documentElement.getAttribute("lang");
    if (attr === "en" || attr === "id") {
      set({ locale: attr });
      return;
    }
    try {
      const stored = window.localStorage.getItem(LOCALE_STORAGE_KEY);
      if (stored === "en" || stored === "id") {
        applyLocaleToDocument(stored);
        set({ locale: stored });
      }
    } catch {
      // Ignore storage errors in restricted browsing contexts.
    }
  },
  setLocale: (locale) => {
    applyLocaleToDocument(locale);
    set({ locale });
  },
  toggleLocale: () => {
    const next: UILocale = get().locale === "en" ? "id" : "en";
    applyLocaleToDocument(next);
    set({ locale: next });
  },
  setMobileNavOpen: (open) => set({ mobileNavOpen: open }),
  closeMobileNav: () => set({ mobileNavOpen: false }),
  toggleSidebarCollapsed: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
}));
