import { describe, expect, test } from "bun:test";
import {
  ALL_NAV_ITEMS,
  countVisibleNavItems,
  EXTENSION_GROUPS,
  findActiveNavContext,
  getShellCopy,
  isNavItemActive,
  normalizeRoute,
  resolveVisibleNavigation,
} from "./navigation";

describe("shared/ui/navigation (main LeftRail layout parity)", () => {
  test("preserves all 28 dashboard routes with unique trailing-slash hrefs and expanded accordions", () => {
    expect(ALL_NAV_ITEMS).toHaveLength(28);
    expect(EXTENSION_GROUPS).toHaveLength(2);
    expect(EXTENSION_GROUPS.every((group) => group.defaultExpanded)).toBe(true);

    const hrefs = ALL_NAV_ITEMS.map((item) => item.href);
    expect(new Set(hrefs).size).toBe(28);
    for (const href of hrefs) {
      expect(href.startsWith("/dashboard/")).toBe(true);
      expect(href.endsWith("/")).toBe(true);
    }
  });

  test("matches exact /dashboard/ root and nested sub-routes accurately", () => {
    expect(normalizeRoute("/dashboard")).toBe("/dashboard/");
    expect(isNavItemActive("/dashboard/", "/dashboard/")).toBe(true);
    expect(isNavItemActive("/dashboard/providers/", "/dashboard/")).toBe(false);
    expect(isNavItemActive("/dashboard/providers/detail/", "/dashboard/providers/")).toBe(true);

    const providerCtx = findActiveNavContext("/dashboard/providers/detail/");
    expect(providerCtx?.item.href).toBe("/dashboard/providers/");
    expect(providerCtx?.sectionLabel.en).toBe("Routing");

    const pluginCtx = findActiveNavContext("/dashboard/plugins/");
    expect(pluginCtx?.accordionId).toBe("features");
  });

  test("filters menu items across English labels, Indonesian labels, and keywords", () => {
    const all = resolveVisibleNavigation({ query: "" });
    expect(countVisibleNavItems(all)).toBe(28);

    const byEn = resolveVisibleNavigation({ query: "token saver" });
    expect(byEn.operations.map((i) => i.href)).toEqual(["/dashboard/token-saver/"]);

    const byId = resolveVisibleNavigation({ query: "pengaturan" });
    expect(byId.settings?.href).toBe("/dashboard/settings/");

    const byKeyword = resolveVisibleNavigation({ query: "bot" });
    expect(byKeyword.extensions.flatMap((g) => g.items).map((i) => i.href)).toContain(
      "/dashboard/telegram/",
    );

    const none = resolveVisibleNavigation({ query: "nonexistent-xyz-query" });
    expect(countVisibleNavItems(none)).toBe(0);
  });

  test("hides items whose backend moduleId is not registered in /api/health", () => {
    const restricted = resolveVisibleNavigation({
      query: "",
      registeredModules: ["system", "providers", "keys", "settings"],
    });

    expect(restricted.routing.map((i) => i.href)).toEqual(["/dashboard/", "/dashboard/providers/"]);
    expect(restricted.operations.map((i) => i.href)).toEqual(["/dashboard/keys/"]);
    expect(restricted.extensions).toHaveLength(0);
    expect(restricted.settings?.href).toBe("/dashboard/settings/");
  });

  test("provides localized shell copy for en and id", () => {
    const en = getShellCopy("en");
    const id = getShellCopy("id");
    expect(en.searchPlaceholder).not.toBe(id.searchPlaceholder);
    expect(en.signOutLabel).not.toBe(id.signOutLabel);
    expect(en.signOutConfirmTitle).not.toBe(id.signOutConfirmTitle);
    expect(en.signOutConfirmMessage).not.toBe(id.signOutConfirmMessage);
    expect(en.signOutCancelLabel).not.toBe(id.signOutCancelLabel);
    expect(en.signOutConfirmLabel).not.toBe(id.signOutConfirmLabel);
    expect(en.emptySearchTitle).not.toBe(id.emptySearchTitle);
    expect(en.expandAllLabel).not.toBe(id.expandAllLabel);
    expect(en.collapseAllLabel).not.toBe(id.collapseAllLabel);
  });
});
