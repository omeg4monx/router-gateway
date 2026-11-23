# `/dashboard` System Overview UI/UX Redesign (Approach A)

- **Date:** 2026-11-24
- **Target:** `web/src/features/system/SystemOverview.tsx` + `web/src/features/system/overview.ts`
- **Status:** Approved (Approach A)

## Intent & Why
Replace the lopsided, developer-raw `/dashboard` overview with a balanced, proportional, modern, bilingual (`en` | `id`) command-center overview that pairs live runtime health with operational pulse (providers, API keys, 24h traffic) and quick actions.

## Key Proportion & UX Fixes
1. **Equal-Height 4-Column KPI Row (`grid-cols-1 sm:grid-cols-2 xl:grid-cols-4`, `p-5 rounded-xl`):**
   - Gateway & SQLite Health, Process Uptime, Providers & API Keys, and 24h Traffic & Tokens (eliminating the 27-module comma wall inside Card 4).
2. **Quick Actions Row (`grid-cols-2 sm:grid-cols-4 gap-3`):**
   - Direct shortcuts to Providers, API Keys, Models & Routing, and Basic Chat.
3. **Balanced 7/5 Bottom Grid (`lg:grid-cols-12`):**
   - Left (`lg:col-span-7`): Interactive 2-column compact **Registered Kernel Modules** catalog with live filter input (`h-8`), bounded scroll height (`max-h-[420px]`), bilingual descriptions, and direct links to each module's studio page.
   - Right (`lg:col-span-5`): **Build & Security Posture** card with semantic Badges (`Active` / `Enforced` / `Open`), one-click copy for Git Commit, and quick Settings link.
4. **100% Bilingual EN / ID** synced with `useUIStore((s) => s.locale)` and semantic CSS tokens from `globals.css`.
