# Post-Login Dashboard Layout & Categorized Sidebar Implementation Plan

**Goal:** Redesign `AppShell` with categorized collapsible sidebar navigation, instant bilingual search (`Ctrl+K` / `/`), proportional modern topbar with breadcrumbs & `LocaleToggle`, and live gateway status footer.

**Spec:** `docs/superpowers/specs/2026-11-24-dashboard-shell-navigation-design.md`

## Tasks
1. **Task 1 (TDD): Navigation Taxonomy, Bilingual Copy & Search Engine (`web/src/shared/ui/navigation.ts`)**
   - Write failing tests in `web/src/shared/ui/navigation.test.ts` verifying all 29 routes, uniqueness, nested active matching, EN/ID labels, and cross-language keyword search.
   - Implement `NAV_GROUPS`, `getShellCopy(locale)`, `normalizeRoute`, `isNavItemActive`, `findActiveNavContext`, and `filterNavGroups`.
2. **Task 2: Modern Proportional `AppShell.tsx`**
   - Upgrade `web/src/shared/ui/AppShell.tsx` with:
     - Inline search bar + `Ctrl+K` / `/` shortcut + clear button + empty state.
     - Collapsible category groups with item counts, auto-expanding the active route's group and search matches.
     - Collapsed `72px` icon rail with group dividers, tooltips, and search-expand trigger.
     - Sticky translucent topbar (`h-16`) with category/page breadcrumb, uniform `h-9` utility controls (`Quick Search`, `LocaleToggle`, `ThemeToggle`, `Sign Out`).
     - Live gateway status footer (`/api/health` + `/api/version`).
     - Translucent mobile drawer backdrop.
3. **Task 3: Full Verification & Graphify Update**
   - Run `cd web && bun test && bun run typecheck && bun run check && bun run build` and `graphify update .`.
