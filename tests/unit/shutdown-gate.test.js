import { describe, it, expect, vi, beforeEach } from "vitest";
import fs from "node:fs";
import path from "node:path";
import { FULL_PERMISSIONS, canOpenPage } from "@/lib/auth/permissionPaths";

// Shutdown and update are administrator actions. The routes themselves do not read
// the session: the auth guard (src/dashboardGuard.js, run by custom-server.js ahead
// of Next) lists both in ALWAYS_PROTECTED, so the gate is asserted there.
const mocks = vi.hoisted(() => ({
  next: Symbol("next"),
  getSettings: vi.fn(),
  verifyDashboardAuthToken: vi.fn(),
  getDashboardAuthSession: vi.fn(),
}));

vi.mock("next/server", () => ({
  NextResponse: {
    next: vi.fn(() => mocks.next),
    json: vi.fn((body, init) => ({ status: init?.status || 200, body })),
    redirect: vi.fn((url) => ({ status: 307, url })),
  },
}));
vi.mock("@/lib/localDb", () => ({ getSettings: mocks.getSettings, validateApiKey: vi.fn() }));
vi.mock("@/shared/utils/machineId", () => ({ getConsistentMachineId: vi.fn(async () => "cli-token") }));
vi.mock("@/lib/auth/dashboardSession", () => ({
  verifyDashboardAuthToken: mocks.verifyDashboardAuthToken,
  getDashboardAuthSession: mocks.getDashboardAuthSession,
}));

const { proxy } = await import("../../src/dashboardGuard.js");

const ADMIN_ROUTES = ["/api/version/shutdown", "/api/version/update"];

function post(pathname, token) {
  return {
    method: "POST",
    nextUrl: { pathname, searchParams: new URLSearchParams() },
    headers: new Headers(),
    cookies: { get: (name) => (name === "auth_token" && token ? { name, value: token } : undefined) },
    url: `http://localhost${pathname}`,
  };
}

describe("shutdown and update are administrator actions", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Login switched off must not open these routes.
    mocks.getSettings.mockResolvedValue({ requireLogin: false });
    mocks.verifyDashboardAuthToken.mockResolvedValue(false);
    mocks.getDashboardAuthSession.mockResolvedValue(null);
  });

  it.each(ADMIN_ROUTES)("%s refuses a key session, even one holding every permission", async (route) => {
    mocks.verifyDashboardAuthToken.mockResolvedValue(true);
    mocks.getDashboardAuthSession.mockResolvedValue({ role: "apikey", permissions: FULL_PERMISSIONS });

    expect((await proxy(post(route, "key-session"))).status).toBe(403);
  });

  it.each(ADMIN_ROUTES)("%s refuses a caller with no valid session", async (route) => {
    expect((await proxy(post(route))).status).toBe(401);
    expect((await proxy(post(route, "forged"))).status).toBe(401);
  });

  it.each(ADMIN_ROUTES)("%s lets a password session through", async (route) => {
    mocks.verifyDashboardAuthToken.mockResolvedValue(true);
    mocks.getDashboardAuthSession.mockResolvedValue({ role: "admin" });

    expect(await proxy(post(route, "password-session"))).toBe(mocks.next);
  });

  // custom-server.js answers a guard that throws with 503 (auth-fail-closed.test.js).
  it("a session that cannot be read is refused, not assumed", async () => {
    mocks.getDashboardAuthSession.mockRejectedValue(new Error("session store unreadable"));

    await expect(proxy(post("/api/version/shutdown", "any"))).rejects.toThrow();
  });

  it("a key session never opens the profile page that carries the red button", () => {
    expect(canOpenPage("/dashboard/profile", FULL_PERMISSIONS)).toBe(false);
  });

  // The menu is gated where it is mounted, not inside HeaderMenu.js.
  it("the dashboard shell never mounts the header menu (and its Shutdown item) for a key session", () => {
    const read = (rel) => fs.readFileSync(path.join(process.cwd(), rel), "utf8");
    expect(read("src/shared/components/LeftRail.js")).toMatch(/!isApiKeyUser \? \(\s*<HeaderMenu/);
    expect(read("src/shared/components/MobileNav.js")).not.toContain("<HeaderMenu");
  });
});
