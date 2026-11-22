# Web TypeScript migration design

`main` is a reference, not a parity target. Base `web/typescript` on `migration/next-go` at `1b1182655f`; propose integration there only. This stage establishes a typed, reproducible Next.js web foundation and proves it on the gateway status surface. It does not rewrite the dashboard in one PR or port JavaScript backend code that Go will replace.

## Decisions

- **Bun is mandatory** for installing and running the app and test scripts in local development and PR CI, with one Bun lock per existing package root and no npm fallback. Pin Bun `1.4.5`. Run Next.js, the custom server, and Vitest on **Node.js `26.13.1` LTS**, the latest verified stable v26 patch. Bun is the toolchain entry point, not the Next server runtime: forcing `bun --bun` would bypass the requested Node runtime and change SQLite adapter selection. At implementation start, recheck official stable release metadata and update these concrete pins together only if newer stable patches are available in the selected lines (Node v26 and Bun v1.4).
- **Incremental strict TypeScript.** Keep `allowJs: true`, `checkJs: false`, and `strict: true`. Next.js owns generated route declarations; preserve existing server/client boundaries and the `@/*` and `open-sse` aliases. Align ambient Node and React types with Node 26 and React 19. No new global `any` shim, TypeScript suppression, or disabled Next build validation.
- **Pilot:** type `/dashboard/endpoint`'s server page, browser health helpers and status alert. Preserve auth, key handling, tunnel actions and endpoint paths; demonstrate usable status and copy behavior in the existing UI. Reuse current semantic tokens and layout. Later typed slices should cover the dashboard shell and features, followed by remaining public UI and shared client code; use route inventory to track completion. Do not rename `src/app/api/**`, `open-sse/**`, `src/sse/**` or `custom-server.js` to TypeScript.
- **Boundaries:** Go checks and product routes are unchanged. Docker publishing and global CLI packaging are separate legacy release paths and not part of the Bun cutover here. Root Docker still needs a separate Bun migration and old-tag republish check before a Bun-only release is claimed; it must not be presented as migrated by this stage.
- **Standalone verification:** run the source-root custom-server entry point under Node v26. It launches the built standalone Next.js server while retaining access to the source guard and existing asset staging. This proves supported source startup, not that the traced directory is independently distributable.

## Alternatives

A bulk rename would conflate frontend typing with backend migration and obscure failures. Running all JavaScript on Bun would add a runtime and SQLite compatibility migration. Choose mandatory Bun tooling with Node runtime and one strict feature slice; each later slice can then build on the same verified workflow.

## Acceptance

Clean frozen Bun installs and an out-of-sync fixture test for both package roots; explicit Node 26 execution in dev, production standalone and Vitest; strict type checking and production build in PR CI; existing root-scoped Vitest green without live provider credentials. Pilot browser QA on a disposable loopback instance covers keyboard operation, status text and links, and desktop/mobile layouts. No release, push, or merge of `web/typescript` as part of planning.
