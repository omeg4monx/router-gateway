import { beforeEach, describe, expect, test } from "bun:test";
import { useUIStore } from "./uiStore";

describe("useUIStore (Zustand global UI store)", () => {
  beforeEach(() => {
    useUIStore.setState({
      theme: "dark",
      mobileNavOpen: false,
      sidebarCollapsed: false,
    });
  });

  test("toggles theme between dark and light deterministically", () => {
    expect(useUIStore.getState().theme).toBe("dark");
    useUIStore.getState().toggleTheme();
    expect(useUIStore.getState().theme).toBe("light");
    useUIStore.getState().toggleTheme();
    expect(useUIStore.getState().theme).toBe("dark");
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
