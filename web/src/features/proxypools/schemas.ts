import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ProxyPoolScheme = components["schemas"]["ProxyPoolScheme"];
export type ProxyPoolStatus = components["schemas"]["ProxyPoolStatus"];
export type ProxyPoolSummary = components["schemas"]["ProxyPoolSummary"];
export type ProxyBindingSummary = components["schemas"]["ProxyBindingSummary"];
export type ListProxyPoolsResponse = components["schemas"]["ListProxyPoolsResponse"];
export type CreateProxyPoolRequest = components["schemas"]["CreateProxyPoolRequest"];
export type UpdateProxyPoolRequest = components["schemas"]["UpdateProxyPoolRequest"];
export type TestProxyPoolRequest = components["schemas"]["TestProxyPoolRequest"];
export type TestProxyPoolResponse = components["schemas"]["TestProxyPoolResponse"];
export type ScanProxyCandidatesRequest = components["schemas"]["ScanProxyCandidatesRequest"];
export type ScanProxyCandidateResult = components["schemas"]["ScanProxyCandidateResult"];
export type ScanProxyCandidatesResponse = components["schemas"]["ScanProxyCandidatesResponse"];
export type UpsertProxyBindingRequest = components["schemas"]["UpsertProxyBindingRequest"];

const ALLOWED_SCHEMES = new Set(["http:", "https:", "socks5:", "socks5h:"]);

/**
 * Normalizes Webshare / ProxyScrape shorthand (`host:port` or `host:port:user:pass`)
 * into a standard `http://` proxy URL while leaving explicit scheme URLs untouched.
 */
export function normalizeProxyInput(raw: string): string {
  const trimmed = raw.trim();
  if (!trimmed || trimmed.includes("://")) {
    return trimmed;
  }
  const parts = trimmed.split(":");
  if (parts.length === 2 && parts[0] && /^\d+$/.test(parts[1] ?? "")) {
    return `http://${parts[0]}:${parts[1]}`;
  }
  if (parts.length === 4 && parts[0] && /^\d+$/.test(parts[1] ?? "") && parts[2]) {
    const [host, port, user, pass] = parts;
    return `http://${encodeURIComponent(user)}:${encodeURIComponent(pass ?? "")}@${host}:${port}`;
  }
  return trimmed;
}

export function isValidProxyUrl(raw: string): boolean {
  const normalized = normalizeProxyInput(raw);
  if (!normalized) {
    return false;
  }
  try {
    const parsed = new URL(normalized);
    return ALLOWED_SCHEMES.has(parsed.protocol.toLowerCase()) && parsed.hostname.length > 0;
  } catch {
    return false;
  }
}

export function parseCandidateLines(raw: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const line of raw.split(/\r?\n/)) {
    const normalized = normalizeProxyInput(line);
    if (!normalized || seen.has(normalized)) {
      continue;
    }
    seen.add(normalized);
    out.push(normalized);
  }
  return out;
}

export const createProxyPoolSchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, "Pool name is required")
    .max(128, "Pool name must be 128 characters or fewer"),
  url: z
    .string()
    .trim()
    .transform(normalizeProxyInput)
    .refine(
      isValidProxyUrl,
      "Must be a valid http://, https://, socks5://, socks5h://, or host:port[:user:pass] proxy",
    ),
  enabled: z.boolean().default(true),
});

export type CreateProxyPoolFormInput = z.input<typeof createProxyPoolSchema>;
export type CreateProxyPoolFormValues = z.output<typeof createProxyPoolSchema>;

export const scanCandidatesSchema = z.object({
  candidatesText: z
    .string()
    .refine((v) => parseCandidateLines(v).length > 0, "Enter at least one candidate proxy URL")
    .refine((v) => parseCandidateLines(v).length <= 50, "Maximum 50 candidate URLs per scan"),
  autoAdd: z.boolean().default(false),
  namePrefix: z.string().trim().max(64).optional(),
  timeoutMs: z.coerce.number().int().min(200).max(15000).default(3000),
});

export type ScanCandidatesFormInput = z.input<typeof scanCandidatesSchema>;
export type ScanCandidatesFormValues = z.output<typeof scanCandidatesSchema>;
