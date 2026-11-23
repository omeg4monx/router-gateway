# Login Page UI/UX Redesign & Bilingual (`en`/`id`) Implementation Plan

> **For agentic workers:** Execute task-by-task using TDD (`test-driven-development`) and verify with `verification-before-completion`.

**Goal:** Redesign `/login/` with a modern, friendly, informative, strictly proportional and standardized control-plane UI/UX supporting both English (`en`) and Bahasa Indonesia (`id`).

**Architecture:** Extend the existing Zustand `useUIStore` with persisted `locale: "en" | "id"` + anti-FOUC bootstrap in `layout.tsx` and a reusable `LocaleToggle` control. Add a typed bilingual dictionary & pure presentation helpers (`web/src/features/auth/i18n.ts`), then upgrade `LoginForm.tsx` and `web/src/app/login/page.tsx` using existing `/api/health`, `/api/version`, and `/api/auth/me` hooks with strict spacing/typography/sizing standardization across dark & light themes.

**Tech Stack:** Next.js 16 (static export), React 19, TypeScript, Tailwind CSS v4, Zustand, React Hook Form + Zod, TanStack Query, Lucide React, Bun Test, Biome.

**Spec:** `docs/superpowers/specs/2026-11-24-login-ui-ux-design.md`

## Global Constraints
- Preserve all existing OpenAPI contracts and `loginSchema` default validation behavior (`schemas.test.ts`).
- Use only semantic CSS variables from `web/src/app/globals.css` (zero hardcoded hex colors in TSX) for full dark/light parity.
- All interactive controls must follow standardized heights (`h-9` utility bar controls, `h-11` primary login inputs & CTA), consistent border radii (`rounded-xl` main card, `rounded-lg` tiles, `rounded-md` inputs/buttons), and 4px/8px spacing rhythm.
- Static export compatible (`NEXT_OUTPUT_EXPORT=true next build`).

## Review Focus
- Invalid or missing `localStorage` locale values must fall back safely to `"en"` without throwing in restricted storage contexts.
- Unrecognized backend error messages must still display cleanly without blanking out or crashing.
- Zero/negative/fractional uptime values from `/api/health` must format deterministically in both `en` and `id`.
- Switching locale mid-typing on `/login` must preserve the entered password and update all visible labels, posture descriptions, and errors immediately.
- Mobile viewports (`< lg`) must display the sign-in card above the fold before the gateway telemetry panel.

---

### Task 1: Persisted `locale` (`en` | `id`) in `useUIStore`, `layout.tsx` Bootstrap, & `LocaleToggle`
- **Files:** `web/src/shared/store/uiStore.ts`, `web/src/shared/store/uiStore.test.ts`, `web/src/app/layout.tsx`, `web/src/shared/ui/LocaleToggle.tsx`
- **Interfaces:**
  - `export type UILocale = "en" | "id"`
  - `useUIStore`: `locale`, `hydrateLocale()`, `setLocale(locale)`, `toggleLocale()`
  - `LocaleToggle`: segmented `EN | ID` control matching `ThemeToggle` height (`h-9`)

### Task 2: Auth Bilingual Dictionary & Pure UX Helpers (`web/src/features/auth/i18n.ts`)
- **Files:** `web/src/features/auth/i18n.ts`, `web/src/features/auth/i18n.test.ts`
- **Interfaces:**
  - `getLoginCopy(locale: UILocale): LoginCopy`
  - `formatAuthError(rawMessage: string | null | undefined, locale: UILocale): { title: string; message: string; detail?: string } | null`
  - `describeAuthPosture(input: { isPending: boolean; isError: boolean; session?: AuthMeResponse }, locale: UILocale): { code: string; badgeVariant: "default" | "secondary" | "success" | "warning" | "destructive"; summary: string }`
  - `formatUptime(seconds: number | null | undefined, locale: UILocale): string`

### Task 3: Upgraded `LoginForm.tsx`
- **Files:** `web/src/features/auth/LoginForm.tsx`
- **Scope:** Localized labels/placeholders/buttons, friendly error mapping with optional technical detail, real-time Caps Lock detection (`getModifierState("CapsLock")`), standardized `h-11` input & CTA with loading spinner (`Loader2`), accessible show/hide password toggle, and collapsible "Need help signing in?" guidance.

### Task 4: Upgraded `web/src/app/login/page.tsx`
- **Files:** `web/src/app/login/page.tsx`
- **Scope:** Proportional 5/7 desktop split + mobile-first ordering (form above the fold), subtle ambient Ruby Rose / cyan lighting + hairline grid backdrop, unified `h-9` utility bar (`LocaleToggle` + `ThemeToggle`), live gateway telemetry panel (`/api/health` + `/api/version` + `/api/auth/me` posture card), skeleton loading state, actionable offline recovery state, welcoming uninitialized `/init/` state, and smooth redirect state.

### Task 5: Full Verification & Graphify Update
- **Commands:** `cd web && bun test && bun run typecheck && bun run check && bun run build`, then `graphify update .`.
