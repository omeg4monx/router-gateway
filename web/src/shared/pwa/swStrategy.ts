export type FetchRouteKind = "ignore" | "network-only" | "navigate" | "asset";

export const NETWORK_ONLY_PREFIXES = [
  "/api",
  "/v1",
  "/v1beta",
  "/codex",
  "/responses",
  "/models",
] as const;

export function isNetworkOnlyPathname(pathname: string): boolean {
  if (pathname === "/sw.js") {
    return true;
  }
  return NETWORK_ONLY_PREFIXES.some(
    (prefix) => pathname === prefix || pathname.startsWith(`${prefix}/`),
  );
}

export function classifyFetchRequest(
  url: URL,
  method: string,
  mode: string,
  origin: string,
): FetchRouteKind {
  if (url.origin !== origin) {
    return "ignore";
  }
  if (method.toUpperCase() !== "GET") {
    return "network-only";
  }
  if (isNetworkOnlyPathname(url.pathname)) {
    return "network-only";
  }
  if (mode === "navigate") {
    return "navigate";
  }
  return "asset";
}
