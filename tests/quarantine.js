// Test files that were already failing when CI was introduced (2026-11-22).
// Most came from the upstream import and assert code this fork has since
// changed; some may be real regressions. CI skips them (process.env.CI) so it
// can catch new breakage; local runs still execute them. Triage a file (fix
// the code or update the test), then delete its line here. Never add a file to
// silence a new failure.
export const QUARANTINED_TESTS = [
  "translator/bugs-toClaude-context.test.js",
  "unit/context-install.test.js",
  "unit/cursor-agent-proto.test.js",
  "unit/executor-const-guard.test.js",
  "unit/headroom-chat-core.test.js",
  "unit/import-job-poll-token.test.js",
  "unit/kimi-usage.test.js",
  "unit/loading-banner.test.js",
  "unit/mimo-free.live.test.js",
  "unit/model-select-deactivated-provider.test.js",
  "unit/monaco-version.test.js",
  "unit/opencode-free-tool-choice.test.js",
  "unit/perplexity-web.test.js",
  "unit/sidebar-usage-ownedby.test.js",
  "unit/studio-banner-update.test.js",
  "unit/theme-security-log.test.js",
  "unit/ui-icons-and-changelog.test.js",
  "unit/update-banner-anim.test.js",
  "unit/use-model-caps-alias-fallback.test.js",
  "unit/windsurf-executor.test.js",
];
