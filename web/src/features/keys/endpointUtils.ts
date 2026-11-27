import type { ApiKeySummary } from "./schemas";

export type GatewayRouteId = "domain" | "local" | "cloudflare" | "tailscale";
export type KeyStatusFilter = "all" | "active" | "disabled" | "expired";
export type KeySnippetFormat = "openai" | "anthropic";
export type KeySnippetLang = "curl" | "python" | "node" | "go";

export interface ResolvedGatewayOrigins {
  domainOrigin: string | null;
  localOrigin: string;
  defaultRoute: "domain" | "local";
  isDomainHost: boolean;
}

export interface KeyStatusBuckets {
  all: number;
  active: number;
  disabled: number;
  expired: number;
}

export interface KeyExpiryState {
  status: "none" | "active" | "soon" | "expired";
  label: string;
}

const DEFAULT_FALLBACK_ORIGIN = "http://127.0.0.1:20128";
const DEFAULT_LOCALHOST_ORIGIN = "http://localhost:20128";
const EXPIRING_SOON_MS = 7 * 24 * 60 * 60 * 1000;

function isLocalHostname(hostname: string): boolean {
  const clean = hostname
    .trim()
    .toLowerCase()
    .replace(/^\[|\]$/g, "");
  if (!clean) {
    return true;
  }
  if (
    clean === "localhost" ||
    clean.endsWith(".localhost") ||
    clean === "::1" ||
    clean === "0.0.0.0"
  ) {
    return true;
  }
  return /^127(?:\.\d{1,3}){3}$/.test(clean);
}

function isIpLiteral(hostname: string): boolean {
  const clean = hostname.trim().replace(/^\[|\]$/g, "");
  return /^\d{1,3}(?:\.\d{1,3}){3}$/.test(clean) || clean.includes(":");
}

export function normalizeOrigin(raw: string | undefined | null): string {
  const trimmed = (raw ?? "").trim();
  if (!trimmed) {
    return "";
  }
  return trimmed.replace(/\/+$/, "");
}

export function resolveGatewayOrigins(
  windowOrigin: string | undefined | null,
): ResolvedGatewayOrigins {
  const normalized = normalizeOrigin(windowOrigin);
  if (!normalized) {
    return {
      domainOrigin: null,
      localOrigin: DEFAULT_FALLBACK_ORIGIN,
      defaultRoute: "local",
      isDomainHost: false,
    };
  }

  try {
    const parsed = new URL(normalized);
    if (isLocalHostname(parsed.hostname)) {
      return {
        domainOrigin: null,
        localOrigin: normalized,
        defaultRoute: "local",
        isDomainHost: false,
      };
    }

    const isDomainHost = !isIpLiteral(parsed.hostname) && parsed.hostname.includes(".");
    return {
      domainOrigin: normalized,
      localOrigin: DEFAULT_LOCALHOST_ORIGIN,
      defaultRoute: "domain",
      isDomainHost,
    };
  } catch {
    return {
      domainOrigin: null,
      localOrigin: normalized,
      defaultRoute: "local",
      isDomainHost: false,
    };
  }
}

export function buildGatewayUrls(rawOrigin: string | undefined | null): {
  origin: string;
  openaiBaseUrl: string;
  claudeBaseUrl: string;
} {
  const origin = normalizeOrigin(rawOrigin) || DEFAULT_FALLBACK_ORIGIN;
  return {
    origin,
    openaiBaseUrl: `${origin}/v1`,
    claudeBaseUrl: origin,
  };
}

export function isKeyExpiredAt(expiresAt: string | null | undefined, nowMs = Date.now()): boolean {
  if (!expiresAt) {
    return false;
  }
  const parsed = Date.parse(expiresAt);
  return !Number.isNaN(parsed) && parsed <= nowMs;
}

export function getKeyExpiryState(
  expiresAt: string | null | undefined,
  nowMs = Date.now(),
  isId = false,
): KeyExpiryState {
  if (!expiresAt) {
    return { status: "none", label: "" };
  }
  const parsed = Date.parse(expiresAt);
  if (Number.isNaN(parsed)) {
    return { status: "active", label: expiresAt };
  }
  const delta = parsed - nowMs;
  if (delta <= 0) {
    return { status: "expired", label: isId ? "Kedaluwarsa" : "Expired" };
  }
  const days = Math.ceil(delta / (24 * 60 * 60 * 1000));
  if (delta <= EXPIRING_SOON_MS) {
    return {
      status: "soon",
      label: isId ? `Berakhir ${days}h lagi` : `Expires in ${days}d`,
    };
  }
  const formatted = new Date(parsed).toLocaleDateString([], {
    month: "short",
    day: "2-digit",
    year: "numeric",
  });
  return {
    status: "active",
    label: isId ? `Berakhir ${formatted}` : `Expires ${formatted}`,
  };
}

export function countKeyStatusBuckets(
  keys: readonly ApiKeySummary[],
  nowMs = Date.now(),
): KeyStatusBuckets {
  const buckets: KeyStatusBuckets = {
    all: keys.length,
    active: 0,
    disabled: 0,
    expired: 0,
  };

  for (const key of keys) {
    if (key.disabled) {
      buckets.disabled += 1;
    } else if (isKeyExpiredAt(key.expires_at, nowMs)) {
      buckets.expired += 1;
    } else {
      buckets.active += 1;
    }
  }
  return buckets;
}

export function filterKeysByStatusAndSearch(
  keys: readonly ApiKeySummary[],
  status: KeyStatusFilter,
  query: string,
  nowMs = Date.now(),
): ApiKeySummary[] {
  const needle = query.trim().toLowerCase();
  return keys.filter((key) => {
    const expired = isKeyExpiredAt(key.expires_at, nowMs);
    if (status === "active" && (key.disabled || expired)) {
      return false;
    }
    if (status === "disabled" && !key.disabled) {
      return false;
    }
    if (status === "expired" && (key.disabled || !expired)) {
      return false;
    }
    if (needle) {
      const haystack =
        `${key.name} ${key.key_prefix} ${key.allowed_models.join(" ")}`.toLowerCase();
      if (!haystack.includes(needle)) {
        return false;
      }
    }
    return true;
  });
}

export function toDatetimeLocalValue(iso: string | null | undefined): string {
  if (!iso) {
    return "";
  }
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return "";
  }
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}`;
}

export function generateKeySnippet(options: {
  format: KeySnippetFormat;
  lang: KeySnippetLang;
  baseUrl: string;
  apiKey?: string;
  model?: string;
}): string {
  const { origin, openaiBaseUrl } = buildGatewayUrls(options.baseUrl);
  const key = options.apiKey?.trim() || "<YOUR_API_KEY>";
  const model =
    options.model?.trim() || (options.format === "anthropic" ? "claude-sonnet-4-5" : "gpt-4o");

  if (options.format === "anthropic") {
    const url = `${origin}/v1/messages`;
    if (options.lang === "python") {
      return `import anthropic\n\nclient = anthropic.Anthropic(\n    base_url="${origin}",\n    api_key="${key}",\n)\nmsg = client.messages.create(\n    model="${model}",\n    max_tokens=1024,\n    messages=[{"role": "user", "content": "Hello"}],\n)\nprint(msg.content[0].text)`;
    }
    if (options.lang === "node") {
      return `import Anthropic from "@anthropic-ai/sdk";\n\nconst client = new Anthropic({\n  baseURL: "${origin}",\n  apiKey: "${key}",\n});\nconst msg = await client.messages.create({\n  model: "${model}",\n  max_tokens: 1024,\n  messages: [{ role: "user", content: "Hello" }],\n});\nconsole.log(msg.content[0].text);`;
    }
    if (options.lang === "go") {
      return `package main\n\nimport (\n\t"bytes"\n\t"fmt"\n\t"io"\n\t"net/http"\n)\n\nfunc main() {\n\tbody := []byte(\`{"model":"${model}","max_tokens":1024,"messages":[{"role":"user","content":"Hello"}]}\`)\n\treq, _ := http.NewRequest("POST", "${url}", bytes.NewBuffer(body))\n\treq.Header.Set("x-api-key", "${key}")\n\treq.Header.Set("anthropic-version", "2023-06-01")\n\treq.Header.Set("Content-Type", "application/json")\n\tresp, _ := http.DefaultClient.Do(req)\n\tdefer resp.Body.Close()\n\tb, _ := io.ReadAll(resp.Body)\n\tfmt.Println(string(b))\n}`;
    }
    return `curl -X POST "${url}" \\\n  -H "x-api-key: ${key}" \\\n  -H "anthropic-version: 2023-06-01" \\\n  -H "Content-Type: application/json" \\\n  -d '{"model":"${model}","max_tokens":1024,"messages":[{"role":"user","content":"Hello"}]}'`;
  }

  const url = `${openaiBaseUrl}/chat/completions`;
  if (options.lang === "python") {
    return `from openai import OpenAI\n\nclient = OpenAI(base_url="${openaiBaseUrl}", api_key="${key}")\nres = client.chat.completions.create(\n    model="${model}",\n    messages=[{"role": "user", "content": "Hello"}],\n)\nprint(res.choices[0].message.content)`;
  }
  if (options.lang === "node") {
    return `import OpenAI from "openai";\n\nconst client = new OpenAI({ baseURL: "${openaiBaseUrl}", apiKey: "${key}" });\nconst res = await client.chat.completions.create({\n  model: "${model}",\n  messages: [{ role: "user", content: "Hello" }],\n});\nconsole.log(res.choices[0].message.content);`;
  }
  if (options.lang === "go") {
    return `package main\n\nimport (\n\t"bytes"\n\t"fmt"\n\t"io"\n\t"net/http"\n)\n\nfunc main() {\n\tbody := []byte(\`{"model":"${model}","messages":[{"role":"user","content":"Hello"}]}\`)\n\treq, _ := http.NewRequest("POST", "${url}", bytes.NewBuffer(body))\n\treq.Header.Set("Authorization", "Bearer ${key}")\n\treq.Header.Set("Content-Type", "application/json")\n\tresp, _ := http.DefaultClient.Do(req)\n\tdefer resp.Body.Close()\n\tb, _ := io.ReadAll(resp.Body)\n\tfmt.Println(string(b))\n}`;
  }
  return `curl -X POST "${url}" \\\n  -H "Authorization: Bearer ${key}" \\\n  -H "Content-Type: application/json" \\\n  -d '{"model":"${model}","messages":[{"role":"user","content":"Hello"}]}'`;
}
