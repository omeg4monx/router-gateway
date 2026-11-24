import { MAX_ID_RUNES, type ProviderProtocol } from "./schemas";

export interface CustomEndpointPreset {
  id: "ollama" | "vllm" | "litellm" | "lmstudio";
  label: string;
  displayName: string;
  protocol: ProviderProtocol;
  baseUrl: string;
  models: string;
}

export const CUSTOM_ENDPOINT_PRESETS: ReadonlyArray<CustomEndpointPreset> = [
  {
    id: "ollama",
    label: "Ollama",
    displayName: "Local Ollama",
    protocol: "openai-chat",
    baseUrl: "http://127.0.0.1:11434/v1",
    models: "llama3.2, qwen2.5-coder",
  },
  {
    id: "vllm",
    label: "vLLM",
    displayName: "Local vLLM",
    protocol: "openai-chat",
    baseUrl: "http://127.0.0.1:8000/v1",
    models: "llama-3.3-70b, qwen2.5-coder-32b",
  },
  {
    id: "litellm",
    label: "LiteLLM",
    displayName: "LiteLLM Proxy",
    protocol: "openai-chat",
    baseUrl: "http://127.0.0.1:4000/v1",
    models: "gpt-4o, claude-3-5-sonnet",
  },
  {
    id: "lmstudio",
    label: "LM Studio",
    displayName: "LM Studio",
    protocol: "openai-chat",
    baseUrl: "http://127.0.0.1:1234/v1",
    models: "local-model",
  },
];

export function slugifyProviderId(input: string): string {
  const normalized = input
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9-_]+/g, "-")
    .replace(/-{2,}/g, "-")
    .replace(/^-+|-+$/g, "");

  if (normalized.length <= MAX_ID_RUNES) {
    return normalized;
  }
  return normalized.slice(0, MAX_ID_RUNES).replace(/-+$/g, "");
}

export interface ParsedOAuthCallback {
  valid: boolean;
  code: string;
  state: string;
}

export function parseOAuthCallbackUrl(rawUrl: string): ParsedOAuthCallback {
  const trimmed = rawUrl.trim();
  if (!trimmed) {
    return { valid: false, code: "", state: "" };
  }
  try {
    const parsed = new URL(trimmed);
    let code = parsed.searchParams.get("code")?.trim() ?? "";
    let state = parsed.searchParams.get("state")?.trim() ?? "";

    if ((!code || !state) && parsed.hash.length > 1) {
      const hashParams = new URLSearchParams(parsed.hash.slice(1));
      code = code || (hashParams.get("code")?.trim() ?? "");
      state = state || (hashParams.get("state")?.trim() ?? "");
    }

    if (code && state) {
      return { valid: true, code, state };
    }
    return { valid: false, code, state };
  } catch {
    return { valid: false, code: "", state: "" };
  }
}

export interface BulkCredentialEntry {
  name: string;
  apiKey: string;
}

export interface ParsedBulkCredentials {
  entries: BulkCredentialEntry[];
  invalidLines: number;
}

export function parseBulkCredentials(rawText: string, providerName: string): ParsedBulkCredentials {
  const lines = rawText
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean);

  const entries: BulkCredentialEntry[] = [];
  let invalidLines = 0;

  for (const [idx, raw] of lines.entries()) {
    const pipeIdx = raw.indexOf("|");
    if (pipeIdx >= 0) {
      const name = raw.slice(0, pipeIdx).trim() || `${providerName} #${idx + 1}`;
      const apiKey = raw.slice(pipeIdx + 1).trim();
      if (!apiKey) {
        invalidLines += 1;
        continue;
      }
      entries.push({ name, apiKey });
      continue;
    }

    entries.push({
      name: `${providerName} #${idx + 1}`,
      apiKey: raw,
    });
  }

  return { entries, invalidLines };
}

export interface DraftTestFingerprintInput {
  providerId?: string;
  protocol: string;
  baseUrl: string;
  apiKey?: string;
  headers?: string;
}

export function buildTestFingerprint(input: DraftTestFingerprintInput): string {
  return JSON.stringify([
    (input.providerId ?? "").trim(),
    input.protocol.trim(),
    input.baseUrl.trim(),
    (input.apiKey ?? "").trim(),
    (input.headers ?? "").trim(),
  ]);
}
