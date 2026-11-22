import { describe, it, expect, vi } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const root = resolve(__dirname, "../..");
const read = (p) => readFileSync(resolve(root, p), "utf8");

describe("theme: the site-wide mode is a password-session setting", () => {
  it("defaults the site to ember when the owner never picked one", () => {
    const route = read("src/app/api/theme/route.js");
    expect(route).toMatch(/theme:\s*"ember"/);
  });

  it("GET /api/theme answers the theme record only and never leaks settings", () => {
    const route = read("src/app/api/theme/route.js");
    expect(route).toContain('getMeta("global_theme_settings"');
    expect(route).not.toContain("getSettings");
  });

  it("POST /api/theme refuses a key-signed session and keeps the stored theme", async () => {
    const setMeta = vi.fn();
    vi.resetModules();
    vi.doMock("@/lib/auth/dashboardPermissions", () => ({
      getSessionContext: async () => ({ session: { role: "apikey" } }),
    }));
    vi.doMock("@/lib/db/helpers/metaStore", () => ({ getMeta: vi.fn(), setMeta }));
    try {
      const { POST } = await import("@/app/api/theme/route.js");
      const res = await POST(new Request("http://x/api/theme", { method: "POST", body: JSON.stringify({ theme: "dark" }) }));
      expect(res.status).toBe(403);
      expect(setMeta).not.toHaveBeenCalled();
    } finally {
      vi.doUnmock("@/lib/auth/dashboardPermissions");
      vi.doUnmock("@/lib/db/helpers/metaStore");
    }
  });

  it("PATCH /api/settings is closed to a key-signed session", async () => {
    const { requiredPermissionsForApiPath } = await import("@/lib/auth/permissionPaths.js");
    // null = no permission opens it, so dashboardGuard answers 403 to any API key session
    expect(requiredPermissionsForApiPath("/api/settings", "PATCH")).toBeNull();
  });

  it("the dashboard shell hides theme controls from a key session", () => {
    const rail = read("src/shared/components/LeftRail.js");
    expect(rail).toMatch(/!isApiKeyUser \? \(\s*<HeaderMenu/);
    const mobile = read("src/shared/components/MobileNav.js");
    expect(mobile).toMatch(/!isApiKeyUser && \(\s*<button[\s\S]{0,200}setThemeModalOpen\(true\)/);
  });

  it("the theme is loaded from the server for every visitor", () => {
    const store = read("src/store/themeStore.js");
    expect(store).toMatch(/initTheme:[\s\S]*fetch\("\/api\/theme"/);
  });
});

describe("security log: sign-ins, refusals and probes are recorded and flagged", () => {
  it("the repo classifies a default-password sign-in as critical", () => {
    const repo = read("src/lib/db/repos/securityLogRepo.js");
    expect(repo).toContain("login_default_password");
    expect(repo).toContain('event.severity === "critical"');
    // probe patterns, not just honest failures
    expect(repo).toContain("SQL injection attempt");
    expect(repo).toContain("Path traversal attempt");
    expect(repo).toContain("Secret file probe");
  });

  it("the login route records every branch", () => {
    const login = read("src/app/api/auth/login/route.js");
    for (const type of [
      "login_success",
      "login_failed",
      "login_locked",
      "apikey_login_success",
      "apikey_login_failed",
      "login_default_password",
    ]) {
      expect(login).toContain(type);
    }
    // the raw secret never reaches the log
    expect(login).not.toMatch(/recordSecurityEvent\([\s\S]{0,200}keyStr(?!.*slice)/);
  });

  it("the security-events table lives in the schema with a version bump", () => {
    const schema = read("src/lib/db/schema.js");
    expect(schema).toContain("securityEvents");
    expect(schema).toContain("accessEvents");
    expect(schema).not.toContain("SCHEMA_VERSION = 11");
  });

  it("custom-server feeds probe-looking responses to the security log", () => {
    const server = read("custom-server.js");
    expect(server).toContain("onRequestFinished");
    expect(server).toContain('res.once("finish"');
    expect(server).toContain("matchRequestProbe");
    expect(server).toContain("recordSecurityEvent");
  });
});