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

export function isValidProxyUrl(raw: string): boolean {
  const trimmed = raw.trim();
  if (!trimmed) {
    return false;
  }
  try {
    const parsed = new URL(trimmed);
    return ALLOWED_SCHEMES.has(parsed.protocol.toLowerCase()) && parsed.hostname.length > 0;
  } catch {
    return false;
  }
}

export function parseCandidateLines(raw: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];
  for (const line of raw.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed || seen.has(trimmed)) {
      continue;
    }
    seen.add(trimmed);
    out.push(trimmed);
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
    .refine(isValidProxyUrl, "Must be a valid http://, https://, socks5://, or socks5h:// URL"),
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
