import { z } from "zod";
import type { components } from "@/shared/api/schema";

export const loginSchema = z.object({
  password: z.string().min(1, "Administrator password is required."),
});

export type LoginFormValues = z.infer<typeof loginSchema>;

export interface AuthRouteState {
  initialized: boolean;
  authenticated: boolean;
  requireLogin: boolean;
  remoteEnforced: boolean;
}

export type AuthRedirectTarget = "/init/" | "/login/" | "/dashboard/";

type AuthSessionPosture = Pick<
  components["schemas"]["AuthMeResponse"],
  "initialized" | "authenticated" | "require_login" | "remote_enforced"
>;

export function normalizeAuthPathname(pathname: string): string {
  const trimmed = pathname.trim();
  if (!trimmed || trimmed === "/") {
    return "/";
  }
  const withLeadingSlash = trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
  return withLeadingSlash.endsWith("/") ? withLeadingSlash : `${withLeadingSlash}/`;
}

export function authSessionToRouteState(session: AuthSessionPosture): AuthRouteState {
  return {
    initialized: session.initialized,
    authenticated: session.authenticated,
    requireLogin: session.require_login,
    remoteEnforced: session.remote_enforced,
  };
}

export function resolvePostAuthRedirect(
  state: AuthRouteState,
  currentPathname: string,
): AuthRedirectTarget | null {
  const normalized = normalizeAuthPathname(currentPathname);

  if (!state.initialized) {
    return normalized === "/init/" ? null : "/init/";
  }

  const loginRequired = state.requireLogin || state.remoteEnforced;
  if (loginRequired && !state.authenticated) {
    return normalized === "/login/" ? null : "/login/";
  }

  if (normalized === "/" || normalized === "/login/") {
    return "/dashboard/";
  }

  return null;
}
