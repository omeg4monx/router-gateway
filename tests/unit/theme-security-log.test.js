import { describe, it, expect } from "vitest";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const root = resolve(__dirname, "../..");
const read = (p) => readFileSync(resolve(root, p), "utf8");

describe("theme: the site-wide mode is a password-session setting", () => {
  it("defaults the site to glass in settings", () => {
    expect(read("src/lib/db/repos/settingsRepo.js")).toMatch(/theme:\s*"glass"/);
  });

  it("GET /api/theme answers one value and never leaks the rest", () => {
    const route = read("src/app/api/theme/route.js");
    expect(route).toMatch(/json\(\s*\{\s*theme:/);
    // GET never reads anything but the one value
    expect(route).not.toMatch(/const settings = await getSettings\(\);[\s\S]{0,120}json\(settings/);
  });

  it("POST /api/theme refuses a key-signed session", () => {
    const route = read("src/app/api/theme/route.js");
    expect(route).toMatch(/session\.role === "apikey"/);
    expect(route).toContain("403");
  });

  it("PATCH /api/settings refuses to write theme for a key-signed session", () => {
    const settings = read("src/app/api/settings/route.js");
    expect(settings).toContain("ADMIN_ONLY_SETTING_KEYS");
    expect(settings).toContain("apikey");
    expect(settings).toContain("403");
  });

  it("ThemeModal shows a lock to a key session instead of clickable modes", () => {
    const modal = read("src/shared/components/ThemeModal.js");
    expect(modal).toContain("canChangeTheme");
    expect(modal).toContain("disabled={!canChangeTheme}");
    expect(modal).toContain("Administrator setting");
  });

  it("HeaderMenu hides the Theme row from a key session", () => {
    const menu = read("src/shared/components/HeaderMenu.js");
    expect(menu).toContain('role !== "apikey"');
    expect(menu).toContain("canChangeTheme && (");
  });

  it("the theme is loaded from the server for every visitor", () => {
    const store = read("src/store/themeStore.js");
    expect(store).toContain("setServerTheme");
    expect(store).toContain("/api/theme");
    // no per-browser persistence as the source of truth anymore
    expect(store).not.toContain('zustand/middleware');
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