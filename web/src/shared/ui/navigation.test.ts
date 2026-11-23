import { describe, expect, test } from "bun:test";
import {
  ALL_NAV_ITEMS,
  filterNavGroups,
  findActiveNavContext,
  getShellCopy,
  isNavItemActive,
  NAV_GROUPS,
  normalizeRoute,
} from "./navigation";

describe("shared/ui/navigation", () => {
  test("preserves all 27 dashboard routes across 6 operational categories with unique trailing-slash hrefs and expanded defaults", () => {
    expect(NAV_GROUPS).toHaveLength(6);
    expect(ALL_NAV_ITEMS).toHaveLength(27);
    expect(NAV_GROUPS.every((group) => group.defaultExpanded)).toBe(true);

    const hrefs = ALL_NAV_ITEMS.map((item) => item.href);
    expect(new Set(hrefs).size).toBe(27);
    expect(hrefs).not.toContain("/login/");
    expect(hrefs).not.toContain("/init/");
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

    const ctx = findActiveNavContext("/dashboard/providers/detail/");
    expect(ctx?.item.href).toBe("/dashboard/providers/");
    expect(ctx?.group.id).toBe("routing");
  });

  test("filters menu items across English labels, Indonesian labels, and keywords", () => {
    expect(filterNavGroups("", "en")).toHaveLength(6);

    const byEn = filterNavGroups("token saver", "en");
    expect(byEn.flatMap((g) => g.items).map((i) => i.href)).toEqual(["/dashboard/token-saver/"]);

    const byId = filterNavGroups("pengaturan", "id");
    expect(byId.flatMap((g) => g.items).map((i) => i.href)).toContain("/dashboard/settings/");

    const byKeyword = filterNavGroups("bot", "en");
    expect(byKeyword.flatMap((g) => g.items).map((i) => i.href)).toContain("/dashboard/telegram/");

    const none = filterNavGroups("nonexistent-xyz-query", "en");
    expect(none).toHaveLength(0);
  });

  test("provides localized shell copy for en and id", () => {
    const en = getShellCopy("en");
    const id = getShellCopy("id");
    expect(en.searchPlaceholder).not.toBe(id.searchPlaceholder);
    expect(en.signOutLabel).not.toBe(id.signOutLabel);
    expect(en.emptySearchTitle).not.toBe(id.emptySearchTitle);
    expect(en.expandAllLabel).not.toBe(id.expandAllLabel);
    expect(en.collapseAllLabel).not.toBe(id.collapseAllLabel);
  });
});
