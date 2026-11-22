# Infrastructure foundation design

## Intent and scope

Keep `main` untouched at `3be429e` for this work. Use it as a reference for existing user workflows and as a benchmark comparison, not as a specification for feature or implementation parity. Build the first rewrite stage on `infra/foundation`; integrate it into `migration/next-go`, not `main`. The foundation gives the later Next.js TypeScript and Go stages a reproducible dependency install, independent PR checks and a safe local startup recipe. Its lack of product behavior changes is a stage boundary, not a constraint on the rewrite.

## Rewrite direction beyond this stage

- **Gateway:** Go owns LLM request handling and backend business logic. Compare gateway overhead, time to first token, p95 latency under concurrent streaming, and RSS using the same hardware, workload, and controlled upstream. Set numerical targets after establishing a baseline.
- **Dashboard:** build an interactive, accessible Next.js/TypeScript UI. Prefer useful operational information over visual churn: provider health, fallback reasons, errors, and usage should be easy to find and understand. Validate key flows at desktop and mobile sizes, including keyboard and failure states.
- **Maintainability:** choose modern tooling for a concrete benefit: fewer overlapping responsibilities, repeatable builds, easier debugging, or simpler deployment. Existing behavior can be redesigned or removed. For intentional changes to public APIs or persisted data, document the impact and provide a versioning or migration path; do not copy bugs or internal designs for parity.

## Current state

- `.nvmrc` pins Node `22.23.4`; `go.mod` specifies Go `1.26.0`. CI installs root and test dependencies without tracked npm locks and uses npm 11 for tests because of a recorded npm 10 resolver failure. The repo tracks a Bun lock for its existing optional Bun workflow.
- `.github/workflows/test.yml` runs Vitest from the repository root; it does not build Next.js or check Go. The default Docker Compose deployment is Node; Go has its own binary and experimental Dockerfile.
- `start:all` starts both processes and forwards the inherited environment. Both implementations use a `db/data.sqlite` path under `DATA_DIR`; Go also has a fallback that may select another local checkout's data if the variable is unset.
- Docker publishing checks out its Dockerfile from the workflow revision and its build context from the release tag. Making that Dockerfile unconditionally require an npm lock would break republishing old tags with no lock.

## Alternatives and decision

1. **Keep the current stack, make source builds repeatable and test both runtimes** (chosen). Reuses the existing npm build, Vitest, Go modules and GitHub Actions without changing release packaging.
2. **Replace the tooling or switch Docker publishing to a new build contract now.** More levers to change and verify, including old-tag republishing, before any product benefit. Leave Bun, Docker images and release workflows as-is in this stage; reproducibility here means npm dependency resolution and the checked source build, not byte-identical images or parity with the production image.

## Design and acceptance

- Keep Node `22.23.4` and Go `1.26.0`; use npm `11.20.0` for root and test installs in the supported local and CI path. Track only the root and `tests/` npm lockfiles in addition to existing tracked locks. Use clean `npm ci` installs; package and lock mismatch must fail rather than silently resolve. Keep existing direct dependency ranges and the optional Bun workflow.
- Extend the existing test workflow into separate Node and Go jobs. Retain `main` pushes, all PRs, and manual dispatch; add pushes to `migration/next-go`. Use checkout without persisted credentials and `contents: read`; no workflow secrets, publishing, deployments or write permissions. Pin the actions used in this workflow to verified upstream commit SHAs. Node must run the existing root-scoped Vitest command and `npm run build`; Go must verify modules, run all Go tests and build the binary with `CGO_ENABLED=0` and read-only module resolution. Both jobs must work on a fresh Linux checkout.
- Document two local terminal commands under one development recipe: start existing Node dev server on `127.0.0.1:20128` and run the existing Go entry point on `127.0.0.1:8080`; set explicit separate disposable `DATA_DIR` paths and a non-default disposable password in both processes. Check `/api/health` independently. Do not pass production secrets or mount existing runtime data. Starting Go does not redirect any Node routes.
- No changes to API contracts, production runtime flags, auth, SQLite schema, provider routing, deployment defaults, publish jobs, or application source. Do not run builds against a real database. Revert infra commits in reverse order on the integration branch if needed; no user data or release needs rollback.

## Handoff

The builder may choose the stable GitHub Actions commit corresponding to each named action version, the disposable directory names used by the recipe, and formatting consistent with the current workflow. No implementation or rollout is authorized by this design alone.
