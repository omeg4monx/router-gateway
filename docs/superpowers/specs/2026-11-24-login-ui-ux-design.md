# Login Page UI/UX Redesign & Bilingual Localization (`en` / `id`)

- **Date:** 2026-11-24
- **Route:** `/login/` (`web/src/app/login/page.tsx`)
- **Status:** Approved (Approach A)

## Intent & Why
Upgrade `/login/` from a flat, developer-raw split screen into a modern, friendly, informative, and strictly proportional/standardized operator entry surface supporting both English (`en`) and Bahasa Indonesia (`id`), without changing any backend API contracts.

## Approaches Considered
1. **Approach A — Modern Control-Plane Split (Selected):**
   - Desktop (`lg:grid-cols-12`, 5/7 split): Left brand + live gateway telemetry panel; right focused sign-in canvas with utility bar (`EN | ID` + `ThemeToggle`).
   - Mobile (`< lg`): Compact top bar, sign-in card rendered first above the fold, gateway telemetry cards stacked cleanly below.
   - Trade-off: Slightly richer left panel than a bare centered card, balanced by strict spacing, proportion, and visual hierarchy rules.
2. **Approach B — Centered Spotlight Card (Rejected):**
   - Single centered card with minimal footer badges. Rejected because it leaves insufficient room to surface live gateway health, auth posture explanations, and operator guidance cleanly.

## Proportionality & Standardization Contract
- **Design Tokens:** Strictly use semantic CSS variables from `web/src/app/globals.css` (Omegamon Zwart dark & Omegamon X-Antibody light) so both themes preserve >= 4.5:1 text contrast and >= 6:1 primary button contrast.
- **Spacing & Rhythm:** Standardized 4px/8px grid (`gap-2`, `gap-3`, `gap-4`, `p-6 sm:p-8` for primary card, `p-4` for telemetry tiles, `space-y-5` for form rhythm).
- **Geometry & Elevation:** `rounded-xl` for main cards, `rounded-lg` for nested telemetry/status tiles, `rounded-md` for inputs and buttons (matching existing UI primitives). Subtle radial ambient lighting + hairline architectural grid on canvas.
- **Typography Scale:** Page title `text-3xl font-semibold tracking-tight`, card heading `text-xl font-semibold tracking-tight`, section eyebrows `font-mono text-[11px] uppercase tracking-[0.14em]`, body `text-sm leading-relaxed`, metadata `font-mono text-xs tabular-nums`.
- **Control Sizing:** All header utility controls (`LocaleToggle`, `ThemeToggle`) share uniform `h-9` height; primary form inputs and submit button share `h-10`–`h-11` proportional sizing; icons standardized (`h-4 w-4` inline, `h-3.5 w-3.5` compact controls).

## Functional & UX Scope
1. **Bilingual Localization (`en` | `id`):**
   - Extend `useUIStore` (`web/src/shared/store/uiStore.ts`) with `locale: "en" | "id"`, `hydrateLocale`, `setLocale`, `toggleLocale` persisted under `omeg4mon-locale` and synced to `document.documentElement.lang`.
   - Anti-FOUC bootstrap in `web/src/app/layout.tsx` reading `omeg4mon-locale`.
   - Reusable `LocaleToggle` component (`web/src/shared/ui/LocaleToggle.tsx`).
   - Typed login/auth dictionary & pure helpers (`web/src/features/auth/i18n.ts`): localized strings, friendly API error mapping, uptime formatting, and auth posture descriptors.
2. **Live Gateway Telemetry (Zero Backend Changes):**
   - Consume existing `useSystemHealth()` (`/api/health`), `useSystemVersion()` (`/api/version`), and `useAuthSession()` (`/api/auth/me`).
   - Live status pill (`Online` / `Degraded` / `Offline`), database status (`SQLite OK` / `Error`), formatted uptime, and build version.
   - Human-friendly Auth Posture tile explaining `REMOTE ENFORCED`, `ENFORCED`, `LOOPBACK OPEN`, and `UNINITIALIZED`.
3. **Friendly Form & State Polish:**
   - Skeleton loading card while `/api/auth/me` is pending.
   - Actionable offline/unreachable state with troubleshooting checklist and Retry button.
   - Welcoming uninitialized state directing operators to `/init/`.
   - Smooth redirecting state for active sessions / loopback open mode.
   - Enhanced `LoginForm`: localized validation & friendly error copy (with subtle technical detail), real-time Caps Lock detection indicator, password visibility toggle, submit spinner (`Loader2`), and collapsible "Need help signing in?" guidance.

## Verification
- Unit tests (`bun test`) covering `uiStore` locale persistence, `auth/i18n` translations, error mapping, uptime formatting, and auth posture helpers.
- Typecheck (`bun run typecheck`), lint/format (`bun run check`), and production static export build (`bun run build`).
