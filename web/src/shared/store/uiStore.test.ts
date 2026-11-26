import { beforeEach, describe, expect, test } from "bun:test";
import { useUIStore } from "./uiStore";

describe("useUIStore (Zustand global UI store)", () => {
  test("defaults theme to light on initial store state", () => {
    expect(useUIStore.getInitialState().theme).toBe("light");
  });

  beforeEach(() => {
    useUIStore.setState({
      theme: "light",
      locale: "en",
      mobileNavOpen: false,
      sidebarCollapsed: false,
    });
  });

  test("toggles theme between light and dark deterministically", () => {
    expect(useUIStore.getState().theme).toBe("light");
    useUIStore.getState().toggleTheme();
    expect(useUIStore.getState().theme).toBe("dark");
    useUIStore.getState().toggleTheme();
    expect(useUIStore.getState().theme).toBe("light");
  });

  test("sets and toggles locale between en and id deterministically", () => {
    expect(useUIStore.getState().locale).toBe("en");
    useUIStore.getState().toggleLocale();
    expect(useUIStore.getState().locale).toBe("id");
    useUIStore.getState().setLocale("en");
    expect(useUIStore.getState().locale).toBe("en");
    useUIStore.getState().setLocale("id");
    expect(useUIStore.getState().locale).toBe("id");
  });

  test("manages mobile navigation drawer open/close state", () => {
    expect(useUIStore.getState().mobileNavOpen).toBe(false);
    useUIStore.getState().setMobileNavOpen(true);
    expect(useUIStore.getState().mobileNavOpen).toBe(true);
    useUIStore.getState().closeMobileNav();
    expect(useUIStore.getState().mobileNavOpen).toBe(false);
  });

  test("toggles sidebar collapsed state", () => {
    expect(useUIStore.getState().sidebarCollapsed).toBe(false);
    useUIStore.getState().toggleSidebarCollapsed();
    expect(useUIStore.getState().sidebarCollapsed).toBe(true);
    useUIStore.getState().toggleSidebarCollapsed();
    expect(useUIStore.getState().sidebarCollapsed).toBe(false);
  });
});
