import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ProviderProtocol = components["schemas"]["ProviderProtocol"];
export type ProviderAuthType = components["schemas"]["ProviderAuthType"];
export type CatalogOAuthDescriptor = components["schemas"]["CatalogOAuthDescriptor"];
export type CatalogAdapterDescriptor = components["schemas"]["CatalogAdapterDescriptor"];
export type ProviderNode = components["schemas"]["ProviderNode"];
export type ProvidersCatalogResponse = components["schemas"]["ProvidersCatalogResponse"];
export type ListProviderNodesResponse = components["schemas"]["ListProviderNodesResponse"];
export type CreateProviderNodeRequest = components["schemas"]["CreateProviderNodeRequest"];
export type TestProviderNodeRequest = components["schemas"]["TestProviderNodeRequest"];
export type ProviderConnectionSummary = components["schemas"]["ProviderConnectionSummary"];
export type ListProviderConnectionsResponse =
  components["schemas"]["ListProviderConnectionsResponse"];
export type CreateProviderConnectionRequest =
  components["schemas"]["CreateProviderConnectionRequest"];
export type UpdateProviderConnectionRequest =
  components["schemas"]["UpdateProviderConnectionRequest"];
export type TestProviderConnectionResponse =
  components["schemas"]["TestProviderConnectionResponse"];
export type OAuthStartRequest = components["schemas"]["OAuthStartRequest"];
export type OAuthStartResponse = components["schemas"]["OAuthStartResponse"];
export type OAuthExchangeRequest = components["schemas"]["OAuthExchangeRequest"];
export type OAuthStatusResponse = components["schemas"]["OAuthStatusResponse"];
export type OAuthPasteRequest = components["schemas"]["OAuthPasteRequest"];
export type OAuthPasteResponse = components["schemas"]["OAuthPasteResponse"];
export type OAuthStopResponse = components["schemas"]["OAuthStopResponse"];

export const PROVIDER_PROTOCOLS = [
  "openai-chat",
  "openai-responses",
  "anthropic",
  "gemini",
] as const;

export const OAUTH_PROVIDER_IDS = [
  "claude",
  "codex",
  "gemini-oauth",
  "qwen",
  "github-copilot",
  "kiro",
  "antigravity",
  "aws-quick",
  "freepi",
  "kimi",
  "grok-cli",
  "nous",
  "minimax-code",
  "minimax-code-global",
  "muse",
  "kilocode",
  "iflow",
  "cline",
  "clinepass",
  "gitlab",
] as const;

export const MAX_NAME_RUNES = 128;
export const MAX_ID_RUNES = 64;
export const MAX_SLICE_ITEMS = 256;
export const MAX_SLICE_ITEM_RUNES = 128;
export const MAX_HEADERS_COUNT = 64;
export const MAX_HEADER_KEY_RUNES = 128;
export const MAX_HEADER_VALUE_RUNES = 1024;
export const MAX_PRIORITY = 100_000;
export const DEFAULT_PRIORITY = 100;

export function parseCommaSeparatedList(input: string): string[] {
  const seen = new Set<string>();
  const out: string[] = [];

  for (const part of input.split(",")) {
    const trimmed = part.trim();
    if (!trimmed || seen.has(trimmed)) {
      continue;
    }
    seen.add(trimmed);
    out.push(trimmed);
  }

  return out;
}

export function formatCommaSeparatedList(items: string[]): string {
  return items.join(", ");
}

export function isValidHttpUrl(raw: string): boolean {
  const trimmed = raw.trim();
  if (!trimmed) {
    return false;
  }
  try {
    const parsed = new URL(trimmed);
    const protocol = parsed.protocol.toLowerCase();
    return (protocol === "http:" || protocol === "https:") && parsed.host.length > 0;
  } catch {
    return false;
  }
}

export type ParseHeadersResult =
  | { ok: true; headers: Record<string, string> }
  | { ok: false; error: string };

export function parseHeadersText(input: string): ParseHeadersResult {
  const headers: Record<string, string> = {};
  const lines = input.split(/\r?\n/);

  for (const rawLine of lines) {
    const line = rawLine.trim();
    if (!line) {
      continue;
    }
    const colonIdx = line.indexOf(":");
    if (colonIdx <= 0) {
      return {
        ok: false,
        error: "Each header line must use 'Key: Value' format with a non-empty header name.",
      };
    }
    const key = line.slice(0, colonIdx).trim();
    const value = line.slice(colonIdx + 1).trim();
    if (!key) {
      return { ok: false, error: "Header name must not be empty." };
    }
    if (Array.from(key).length > MAX_HEADER_KEY_RUNES) {
      return {
        ok: false,
        error: `Header name must not exceed ${MAX_HEADER_KEY_RUNES} characters.`,
      };
    }
    if (Array.from(value).length > MAX_HEADER_VALUE_RUNES) {
      return {
        ok: false,
        error: `Header value must not exceed ${MAX_HEADER_VALUE_RUNES} characters.`,
      };
    }
    headers[key] = value;
  }

  if (Object.keys(headers).length > MAX_HEADERS_COUNT) {
    return {
      ok: false,
      error: `Headers must not exceed ${MAX_HEADERS_COUNT} entries.`,
    };
  }

  return { ok: true, headers };
}

export function formatHeadersText(headers: Record<string, string> | undefined): string {
  if (!headers) {
    return "";
  }
  return Object.entries(headers)
    .map(([key, value]) => `${key}: ${value}`)
    .join("\n");
}

const commaListSchema = (label: string) =>
  z
    .string()
    .transform(parseCommaSeparatedList)
    .refine(
      (items) => items.length <= MAX_SLICE_ITEMS,
      `${label} must not exceed ${MAX_SLICE_ITEMS} items.`,
    )
    .refine(
      (items) => items.every((item) => Array.from(item).length <= MAX_SLICE_ITEM_RUNES),
      `${label} entries must not exceed ${MAX_SLICE_ITEM_RUNES} characters.`,
    );

const headersTextSchema = z.string().transform((raw, ctx) => {
  const result = parseHeadersText(raw);
  if (!result.ok) {
    ctx.addIssue({
      code: "custom",
      message: result.error,
    });
    return z.NEVER;
  }
  return result.headers;
});

const prioritySchema = z.string().transform((raw, ctx) => {
  const trimmed = raw.trim();
  if (trimmed === "") {
    return DEFAULT_PRIORITY;
  }
  if (!/^\d+$/.test(trimmed)) {
    ctx.addIssue({
      code: "custom",
      message: `Priority must be an integer between 0 and ${MAX_PRIORITY}.`,
    });
    return z.NEVER;
  }
  const parsed = Number.parseInt(trimmed, 10);
  if (!Number.isSafeInteger(parsed) || parsed < 0 || parsed > MAX_PRIORITY) {
    ctx.addIssue({
      code: "custom",
      message: `Priority must be an integer between 0 and ${MAX_PRIORITY}.`,
    });
    return z.NEVER;
  }
  return parsed;
});

const requiredNameSchema = (label: string) =>
  z
    .string()
    .trim()
    .min(1, `${label} is required.`)
    .refine(
      (value) => Array.from(value).length <= MAX_NAME_RUNES,
      `${label} must not exceed ${MAX_NAME_RUNES} characters.`,
    );

const requiredHttpUrlSchema = (label: string) =>
  z
    .string()
    .trim()
    .min(1, `${label} is required.`)
    .refine(isValidHttpUrl, `${label} must be a valid http:// or https:// URL.`);

const optionalHttpUrlSchema = (label: string) =>
  z
    .string()
    .trim()
    .refine(
      (value) => value === "" || isValidHttpUrl(value),
      `${label} must be a valid http:// or https:// URL.`,
    );

export const createConnectionSchema = z.object({
  provider_id: z.string().trim().min(1, "Provider is required."),
  name: requiredNameSchema("Connection name"),
  base_url: requiredHttpUrlSchema("Base URL"),
  api_key: z.string().trim().min(1, "API key is required."),
  priority: prioritySchema,
  enabled: z.boolean(),
  models: commaListSchema("Models"),
  headers: headersTextSchema,
});

export type CreateConnectionFormInput = z.input<typeof createConnectionSchema>;
export type CreateConnectionFormOutput = z.output<typeof createConnectionSchema>;

export const DEFAULT_CREATE_CONNECTION_VALUES: CreateConnectionFormInput = {
  provider_id: "openai",
  name: "",
  base_url: "https://api.openai.com/v1",
  api_key: "",
  priority: String(DEFAULT_PRIORITY),
  enabled: true,
  models: "gpt-4o, gpt-4o-mini, o3-mini",
  headers: "",
};

export function toCreateConnectionRequest(
  values: CreateConnectionFormOutput,
): CreateProviderConnectionRequest {
  return {
    provider_id: values.provider_id,
    name: values.name,
    base_url: values.base_url,
    auth_type: "api_key",
    api_key: values.api_key,
    priority: values.priority,
    enabled: values.enabled,
    models: values.models,
    headers: values.headers,
  };
}

export const createNodeSchema = z.object({
  id: z
    .string()
    .trim()
    .refine(
      (value) => value === "" || Array.from(value).length <= MAX_ID_RUNES,
      `Node ID must not exceed ${MAX_ID_RUNES} characters.`,
    )
    .refine(
      (value) => value === "" || /^[a-z0-9-_]+$/.test(value),
      "Node ID must contain only lowercase letters, digits, hyphens, or underscores.",
    ),
  name: requiredNameSchema("Node name"),
  protocol: z.enum(PROVIDER_PROTOCOLS),
  base_url: requiredHttpUrlSchema("Base URL"),
  models: commaListSchema("Models"),
  headers: headersTextSchema,
});

export type CreateNodeFormInput = z.input<typeof createNodeSchema>;
export type CreateNodeFormOutput = z.output<typeof createNodeSchema>;

export const DEFAULT_CREATE_NODE_VALUES: CreateNodeFormInput = {
  id: "",
  name: "",
  protocol: "openai-chat",
  base_url: "",
  models: "",
  headers: "",
};

export function toCreateNodeRequest(values: CreateNodeFormOutput): CreateProviderNodeRequest {
  return {
    ...(values.id ? { id: values.id } : {}),
    name: values.name,
    protocol: values.protocol,
    base_url: values.base_url,
    models: values.models,
    headers: values.headers,
  };
}

export const oauthStartSchema = z.object({
  provider: z.string().trim().min(1, "OAuth provider is required."),
  connection_name: z
    .string()
    .trim()
    .refine(
      (value) => Array.from(value).length <= MAX_NAME_RUNES,
      `Connection name must not exceed ${MAX_NAME_RUNES} characters.`,
    ),
  redirect_uri: optionalHttpUrlSchema("Redirect URI"),
  client_id: z.string().trim(),
  client_secret: z.string().trim(),
  base_url: optionalHttpUrlSchema("Base URL"),
  models: commaListSchema("Models"),
});

export type OAuthStartFormInput = z.input<typeof oauthStartSchema>;
export type OAuthStartFormOutput = z.output<typeof oauthStartSchema>;

export const DEFAULT_OAUTH_START_VALUES: OAuthStartFormInput = {
  provider: "claude",
  connection_name: "",
  redirect_uri: "",
  client_id: "",
  client_secret: "",
  base_url: "",
  models: "",
};

export function toOAuthStartRequest(values: OAuthStartFormOutput): OAuthStartRequest {
  return {
    ...(values.connection_name ? { connection_name: values.connection_name } : {}),
    ...(values.redirect_uri ? { redirect_uri: values.redirect_uri } : {}),
    ...(values.client_id ? { client_id: values.client_id } : {}),
    ...(values.client_secret ? { client_secret: values.client_secret } : {}),
    ...(values.base_url ? { base_url: values.base_url } : {}),
    ...(values.models.length > 0 ? { models: values.models } : {}),
  };
}

export const oauthExchangeSchema = z
  .object({
    callback_url: optionalHttpUrlSchema("Callback URL"),
    code: z.string().trim(),
    state: z.string().trim(),
    connection_name: z
      .string()
      .trim()
      .refine(
        (value) => Array.from(value).length <= MAX_NAME_RUNES,
        `Connection name must not exceed ${MAX_NAME_RUNES} characters.`,
      ),
    priority: prioritySchema,
    enabled: z.boolean(),
    headers: headersTextSchema,
  })
  .superRefine((val, ctx) => {
    const hasCallbackUrl = val.callback_url.length > 0;
    const hasCode = val.code.length > 0;
    const hasState = val.state.length > 0;
    if (hasCallbackUrl || (hasCode && hasState)) {
      return;
    }

    if (!hasCallbackUrl && !hasCode && !hasState) {
      ctx.addIssue({
        code: "custom",
        path: ["callback_url"],
        message: "Paste the callback URL or provide both authorization code and state.",
      });
      return;
    }

    if (!hasCode) {
      ctx.addIssue({
        code: "custom",
        path: ["code"],
        message: "Authorization code is required when callback URL is not provided.",
      });
    }
    if (!hasState) {
      ctx.addIssue({
        code: "custom",
        path: ["state"],
        message: "OAuth state is required when callback URL is not provided.",
      });
    }
  });

export type OAuthExchangeFormInput = z.input<typeof oauthExchangeSchema>;
export type OAuthExchangeFormOutput = z.output<typeof oauthExchangeSchema>;

export const DEFAULT_OAUTH_EXCHANGE_VALUES: OAuthExchangeFormInput = {
  callback_url: "",
  code: "",
  state: "",
  connection_name: "",
  priority: String(DEFAULT_PRIORITY),
  enabled: true,
  headers: "",
};

export function toOAuthExchangeRequest(values: OAuthExchangeFormOutput): OAuthExchangeRequest {
  return {
    ...(values.callback_url ? { callback_url: values.callback_url } : {}),
    ...(values.code ? { code: values.code } : {}),
    ...(values.state ? { state: values.state } : {}),
    ...(values.connection_name ? { connection_name: values.connection_name } : {}),
    priority: values.priority,
    enabled: values.enabled,
    headers: values.headers,
  };
}

export const updateConnectionSchema = z.object({
  name: requiredNameSchema("Connection name"),
  base_url: requiredHttpUrlSchema("Base URL"),
  priority: prioritySchema,
  enabled: z.boolean(),
  models: commaListSchema("Models"),
  headers: headersTextSchema,
  api_key: z.string().trim(),
  access_token: z.string().trim(),
  refresh_token: z.string().trim(),
});

export type UpdateConnectionFormInput = z.input<typeof updateConnectionSchema>;
export type UpdateConnectionFormOutput = z.output<typeof updateConnectionSchema>;

export const DEFAULT_UPDATE_CONNECTION_VALUES: UpdateConnectionFormInput = {
  name: "",
  base_url: "",
  priority: String(DEFAULT_PRIORITY),
  enabled: true,
  models: "",
  headers: "",
  api_key: "",
  access_token: "",
  refresh_token: "",
};

export function connectionToEditFormValues(
  conn: ProviderConnectionSummary,
): UpdateConnectionFormInput {
  return {
    name: conn.name,
    base_url: conn.base_url,
    priority: String(conn.priority),
    enabled: conn.enabled,
    models: formatCommaSeparatedList(conn.models),
    headers: formatHeadersText(conn.headers),
    api_key: "",
    access_token: "",
    refresh_token: "",
  };
}

export function toUpdateConnectionRequest(
  values: UpdateConnectionFormOutput,
  authType: ProviderAuthType,
): UpdateProviderConnectionRequest {
  const payload: UpdateProviderConnectionRequest = {
    name: values.name,
    base_url: values.base_url,
    priority: values.priority,
    enabled: values.enabled,
    models: values.models,
    headers: values.headers,
  };

  if (authType === "api_key" && values.api_key) {
    payload.api_key = values.api_key;
  }
  if (authType === "oauth") {
    if (values.access_token) {
      payload.access_token = values.access_token;
    }
    if (values.refresh_token) {
      payload.refresh_token = values.refresh_token;
    }
  }

  return payload;
}
