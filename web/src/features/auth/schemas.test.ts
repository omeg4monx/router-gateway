import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  type AuthRouteState,
  authSessionToRouteState,
  loginSchema,
  normalizeAuthPathname,
  resolvePostAuthRedirect,
} from "./schemas";

describe("loginSchema", () => {
  test("rejects empty password via safeParse and zodResolver", async () => {
    const parsed = loginSchema.safeParse({ password: "" });
    expect(parsed.success).toBe(false);

    const resolver = zodResolver(loginSchema);
    const resolved = await resolver({ password: "" }, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors.password?.message).toBe("Administrator password is required.");
  });

  test("accepts valid non-empty password without trimming or mutating characters", async () => {
    const password = "  operator-secret-\u00e9!  ";
    const parsed = loginSchema.safeParse({ password });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data.password).toBe(password);
    }

    const resolver = zodResolver(loginSchema);
    const resolved = await resolver({ password }, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
    expect(resolved.values).toEqual({ password });
  });
});

describe("normalizeAuthPathname", () => {
  test("normalizes root, missing leading slash, and trailing slash variants", () => {
    expect(normalizeAuthPathname("")).toBe("/");
    expect(normalizeAuthPathname("/")).toBe("/");
    expect(normalizeAuthPathname("/login")).toBe("/login/");
    expect(normalizeAuthPathname("/login/")).toBe("/login/");
    expect(normalizeAuthPathname("dashboard/settings")).toBe("/dashboard/settings/");
    expect(normalizeAuthPathname("/dashboard/keys/")).toBe("/dashboard/keys/");
  });
});

describe("authSessionToRouteState", () => {
  test("maps AuthMeResponse snake_case fields to AuthRouteState", () => {
    expect(
      authSessionToRouteState({
        initialized: true,
        authenticated: false,
        require_login: true,
        remote_enforced: false,
      }),
    ).toEqual({
      initialized: true,
      authenticated: false,
      requireLogin: true,
      remoteEnforced: false,
    });
  });
});

describe("resolvePostAuthRedirect (Review Focus #5)", () => {
  const pathnames = [
    "/",
    "/login",
    "/login/",
    "/dashboard",
    "/dashboard/",
    "/dashboard/settings/",
    "/init/",
  ];

  test("redirects any non-/init/ route to /init/ when not initialized, and stays on /init/", () => {
    const uninitializedStates: AuthRouteState[] = [];
    for (const authenticated of [false, true]) {
      for (const requireLogin of [false, true]) {
        for (const remoteEnforced of [false, true]) {
          uninitializedStates.push({
            initialized: false,
            authenticated,
            requireLogin,
            remoteEnforced,
          });
        }
      }
    }

    expect(uninitializedStates).toHaveLength(8);
    for (const state of uninitializedStates) {
      expect(resolvePostAuthRedirect(state, "/")).toBe("/init/");
      expect(resolvePostAuthRedirect(state, "/login/")).toBe("/init/");
      expect(resolvePostAuthRedirect(state, "/dashboard/")).toBe("/init/");
      expect(resolvePostAuthRedirect(state, "/dashboard/keys/")).toBe("/init/");
      expect(resolvePostAuthRedirect(state, "/init/")).toBeNull();
    }
  });

  test("redirects unauthenticated users to /login/ when requireLogin or remoteEnforced is active, and stays on /login/", () => {
    const loginRequiredStates: AuthRouteState[] = [
      { initialized: true, authenticated: false, requireLogin: true, remoteEnforced: false },
      { initialized: true, authenticated: false, requireLogin: false, remoteEnforced: true },
      { initialized: true, authenticated: false, requireLogin: true, remoteEnforced: true },
    ];

    for (const state of loginRequiredStates) {
      expect(resolvePostAuthRedirect(state, "/")).toBe("/login/");
      expect(resolvePostAuthRedirect(state, "/dashboard/")).toBe("/login/");
      expect(resolvePostAuthRedirect(state, "/dashboard/settings/")).toBe("/login/");
      expect(resolvePostAuthRedirect(state, "/login")).toBeNull();
      expect(resolvePostAuthRedirect(state, "/login/")).toBeNull();
    }
  });

  test("routes / and /login/ to /dashboard/ and stays on /dashboard/* when authenticated or login is not enforced", () => {
    const allowedStates: AuthRouteState[] = [
      { initialized: true, authenticated: false, requireLogin: false, remoteEnforced: false },
      { initialized: true, authenticated: true, requireLogin: false, remoteEnforced: false },
      { initialized: true, authenticated: true, requireLogin: true, remoteEnforced: false },
      { initialized: true, authenticated: true, requireLogin: false, remoteEnforced: true },
      { initialized: true, authenticated: true, requireLogin: true, remoteEnforced: true },
    ];

    for (const state of allowedStates) {
      expect(resolvePostAuthRedirect(state, "/")).toBe("/dashboard/");
      expect(resolvePostAuthRedirect(state, "/login")).toBe("/dashboard/");
      expect(resolvePostAuthRedirect(state, "/login/")).toBe("/dashboard/");
      expect(resolvePostAuthRedirect(state, "/dashboard")).toBeNull();
      expect(resolvePostAuthRedirect(state, "/dashboard/")).toBeNull();
      expect(resolvePostAuthRedirect(state, "/dashboard/keys/")).toBeNull();
      expect(resolvePostAuthRedirect(state, "/dashboard/settings/")).toBeNull();
    }
  });

  test("guarantees zero infinite redirect loops across all 16 boolean combinations", () => {
    let combinationsChecked = 0;

    for (const initialized of [false, true]) {
      for (const authenticated of [false, true]) {
        for (const requireLogin of [false, true]) {
          for (const remoteEnforced of [false, true]) {
            combinationsChecked += 1;
            const state: AuthRouteState = {
              initialized,
              authenticated,
              requireLogin,
              remoteEnforced,
            };

            for (const pathname of pathnames) {
              const target = resolvePostAuthRedirect(state, pathname);
              if (target !== null) {
                expect(resolvePostAuthRedirect(state, target)).toBeNull();
              }
            }
          }
        }
      }
    }

    expect(combinationsChecked).toBe(16);
  });
});
