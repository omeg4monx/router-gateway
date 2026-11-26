# Post-Login Dashboard Layout & Categorized Searchable Sidebar Design

- **Date:** 2026-11-24
- **Target:** `web/src/shared/ui/AppShell.tsx` + `web/src/shared/ui/navigation.ts`
- **Status:** Approved (Approach A)

## Intent & Why
Replace the flat 29-item sidebar and bare topbar in `AppShell` with a modern, proportional, friendly, bilingual (`en` | `id`) control-plane shell featuring logical collapsible categories, instant cross-language menu search (`Ctrl+K` / `/`), breadcrumb context, and live gateway status.

## Approaches Considered
1. **Approach A — Categorized Collapsible Sidebar + Inline Instant Search & `Ctrl+K` (Selected):**
   - Single-column sidebar (`w-[272px]` expanded / `w-[72px]` icon-rail collapsed) with inline search input at top, 7 logical collapsible categories (6 dashboard + 1 account/setup), and live gateway status footer.
   - Trade-off: Keeps horizontal space efficient and 1-click navigation across categories while preventing vertical overload via accordion groups + instant search.
2. **Approach B — Dual-Rail Icon Navigation + Sub-Panel (Rejected):**
   - Wider footprint (~310px) and requires 2 clicks to switch across categories.

## Proportionality & Standardization Contract
- **Design Tokens:** 100% semantic CSS variables from `web/src/app/globals.css` (light default Omegamon X-Antibody + dark Omegamon Zwart).
- **Dimensions & Rhythm:**
  - Sidebar width: `w-[272px]` expanded, `w-[72px]` collapsed (`transition-[width] duration-200`).
  - Topbar height: `h-16` sticky translucent (`bg-[var(--bg-surface)]/90 backdrop-blur-md`).
  - Header controls (`Sidebar toggle`, `Quick search`, `LocaleToggle`, `ThemeToggle`, `Sign Out`): uniform `h-9` height with `gap-2`.
  - Sidebar search input: `h-9 rounded-lg`.
  - Nav items: uniform `h-9 px-2.5 rounded-lg text-[13px] font-medium` with `h-4 w-4` icons and subtle left active accent indicator.
  - Main content container: `mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 lg:px-8`.

## Navigation Taxonomy (29/29 Routes Preserved)
1. **Core & Monitoring** (`3`): System Overview, Usage & Telemetry, Notifications & Jobs
2. **Routing & Providers** (`5`): Providers, Models & Routing, Model Editor, Translator, Token Saver
3. **Network & Gateway** (`5`): Proxy Pools, Media Gateway, Tunnel, MITM Proxy, PXPIPE
4. **Access & Security** (`5`): API Keys, Teams, Vouchers, Prompt Guard, Security
5. **Workspace & Studios** (`5`): Basic Chat, Model Arena, Telegram Bot, Plugins, Skills
6. **System & Tools** (`4`): Cloud Sync, Legacy Importer, CLI Tools, Settings
7. **Account & Setup** (`2`): Sign In, Initialization

## Search & Keyboard UX
- Instant filtering across English labels, Indonesian labels, category titles, and curated keywords.
- Keyboard shortcuts: `Ctrl+K` / `Cmd+K` or `/` (when not typing in another input/textarea) focuses sidebar search (expanding sidebar if collapsed, or opening mobile drawer on `< lg`).
- `Escape` clears active search query first, or closes mobile drawer.
- When searching: all matching groups expand automatically with match count indicator; non-matching query shows friendly empty state + reset button.
