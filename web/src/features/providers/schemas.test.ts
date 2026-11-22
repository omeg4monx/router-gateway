import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  connectionToEditFormValues,
  createConnectionSchema,
  createNodeSchema,
  oauthExchangeSchema,
  oauthStartSchema,
  parseCommaSeparatedList,
  parseHeadersText,
  toCreateConnectionRequest,
  toCreateNodeRequest,
  toOAuthExchangeRequest,
  toUpdateConnectionRequest,
  updateConnectionSchema,
} from "./schemas";

describe("parseCommaSeparatedList", () => {
  test("splits, trims, drops empty entries, and deduplicates preserving order", () => {
    expect(parseCommaSeparatedList(" gpt-4o, claude-sonnet-4 , , gpt-4o, gemini-2.5-pro ")).toEqual(
      ["gpt-4o", "claude-sonnet-4", "gemini-2.5-pro"],
    );
  });

  test("returns empty array for blank input", () => {
    expect(parseCommaSeparatedList("  , , ")).toEqual([]);
  });
});

describe("parseHeadersText", () => {
  test("parses Key: Value lines and preserves colons in values", () => {
    const parsed = parseHeadersText(
      "X-Org-Id: org_123\n\nX-Callback-Hint: https://example.com:8443/path\n",
    );
    expect(parsed).toEqual({
      ok: true,
      headers: {
        "X-Org-Id": "org_123",
        "X-Callback-Hint": "https://example.com:8443/path",
      },
    });
  });

  test("rejects malformed header line missing colon", () => {
    const parsed = parseHeadersText("X-Missing-Colon");
    expect(parsed.ok).toBe(false);
  });
});

describe("createConnectionSchema", () => {
  test("parses valid API-key connection with CSV models, headers, and priority", async () => {
    const input = {
      provider_id: "openai",
      name: "  Primary OpenAI  ",
      base_url: "https://api.openai.com/v1",
      api_key: "  sk-test-12345  ",
      priority: "25",
      enabled: true,
      models: "gpt-4o, gpt-4o-mini, gpt-4o",
      headers: "X-Project: proj_1\nHTTP-Referer: https://router.local",
    };

    const parsed = createConnectionSchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(toCreateConnectionRequest(parsed.data)).toEqual({
        provider_id: "openai",
        name: "Primary OpenAI",
        base_url: "https://api.openai.com/v1",
        auth_type: "api_key",
        api_key: "sk-test-12345",
        priority: 25,
        enabled: true,
        models: ["gpt-4o", "gpt-4o-mini"],
        headers: {
          "X-Project": "proj_1",
          "HTTP-Referer": "https://router.local",
        },
      });
    }

    const resolver = zodResolver(createConnectionSchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
  });

  test("rejects invalid or non-HTTP base_url", () => {
    const badScheme = createConnectionSchema.safeParse({
      provider_id: "openai",
      name: "Bad URL",
      base_url: "ftp://api.openai.com/v1",
      api_key: "sk-test",
      priority: "100",
      enabled: true,
      models: "",
      headers: "",
    });
    expect(badScheme.success).toBe(false);

    const malformed = createConnectionSchema.safeParse({
      provider_id: "openai-compatible",
      name: "Empty URL",
      base_url: "not-a-url",
      api_key: "sk-test",
      priority: "100",
      enabled: true,
      models: "",
      headers: "",
    });
    expect(malformed.success).toBe(false);
  });

  test("rejects empty API key", () => {
    const parsed = createConnectionSchema.safeParse({
      provider_id: "anthropic",
      name: "Missing Key",
      base_url: "https://api.anthropic.com",
      api_key: "   ",
      priority: "100",
      enabled: true,
      models: "",
      headers: "",
    });
    expect(parsed.success).toBe(false);
    if (!parsed.success) {
      expect(parsed.error.issues[0]?.path).toEqual(["api_key"]);
    }
  });

  test("rejects negative, fractional, or out-of-bounds priority", () => {
    for (const badPriority of ["-1", "10.5", "100001", "abc"]) {
      const parsed = createConnectionSchema.safeParse({
        provider_id: "openai",
        name: "Bad Priority",
        base_url: "https://api.openai.com/v1",
        api_key: "sk-test",
        priority: badPriority,
        enabled: true,
        models: "",
        headers: "",
      });
      expect(parsed.success).toBe(false);
    }
  });
});

describe("createNodeSchema", () => {
  test("validates custom compatible node with optional slug id and CSV models", () => {
    const parsed = createNodeSchema.safeParse({
      id: "local-vllm",
      name: "Local vLLM Cluster",
      protocol: "openai-chat",
      base_url: "http://127.0.0.1:8000/v1",
      models: "qwen-2.5-72b, llama-3.3-70b",
      headers: "X-Cluster: gpu-01",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(toCreateNodeRequest(parsed.data)).toEqual({
        id: "local-vllm",
        name: "Local vLLM Cluster",
        protocol: "openai-chat",
        base_url: "http://127.0.0.1:8000/v1",
        models: ["qwen-2.5-72b", "llama-3.3-70b"],
        headers: { "X-Cluster": "gpu-01" },
      });
    }
  });

  test("omits id when left blank and rejects invalid slug id", () => {
    const blankId = createNodeSchema.safeParse({
      id: "  ",
      name: "Auto ID Node",
      protocol: "anthropic",
      base_url: "https://proxy.internal/anthropic",
      models: "",
      headers: "",
    });
    expect(blankId.success).toBe(true);
    if (blankId.success) {
      expect(toCreateNodeRequest(blankId.data).id).toBeUndefined();
    }

    const badId = createNodeSchema.safeParse({
      id: "UpperCase/Bad",
      name: "Bad ID Node",
      protocol: "anthropic",
      base_url: "https://proxy.internal/anthropic",
      models: "",
      headers: "",
    });
    expect(badId.success).toBe(false);
  });
});

describe("oauthStartSchema", () => {
  test("validates minimal start request and rejects invalid optional URLs", () => {
    const valid = oauthStartSchema.safeParse({
      provider: "claude",
      connection_name: "Work Claude",
      redirect_uri: "",
      client_id: "",
      client_secret: "",
      base_url: "",
      models: "claude-sonnet-4-20250514",
    });
    expect(valid.success).toBe(true);

    const invalidRedirect = oauthStartSchema.safeParse({
      provider: "codex",
      connection_name: "",
      redirect_uri: "javascript:alert(1)",
      client_id: "",
      client_secret: "",
      base_url: "",
      models: "",
    });
    expect(invalidRedirect.success).toBe(false);
  });
});

describe("oauthExchangeSchema", () => {
  test("accepts callback_url alone without manual code/state", () => {
    const parsed = oauthExchangeSchema.safeParse({
      callback_url: "http://localhost:20128/api/oauth/callback?code=auth_123&state=st_456",
      code: "",
      state: "",
      connection_name: "Claude Max",
      priority: "50",
      enabled: true,
      headers: "",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(toOAuthExchangeRequest(parsed.data)).toEqual({
        callback_url: "http://localhost:20128/api/oauth/callback?code=auth_123&state=st_456",
        connection_name: "Claude Max",
        priority: 50,
        enabled: true,
        headers: {},
      });
    }
  });

  test("accepts manual code and state pair without callback_url", () => {
    const parsed = oauthExchangeSchema.safeParse({
      callback_url: "",
      code: "auth_code_abc",
      state: "oauth_state_xyz",
      connection_name: "",
      priority: "",
      enabled: true,
      headers: "",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(toOAuthExchangeRequest(parsed.data)).toEqual({
        code: "auth_code_abc",
        state: "oauth_state_xyz",
        priority: 100,
        enabled: true,
        headers: {},
      });
    }
  });

  test("rejects empty submission and incomplete code-only or state-only submission", () => {
    const empty = oauthExchangeSchema.safeParse({
      callback_url: "",
      code: "",
      state: "",
      connection_name: "",
      priority: "100",
      enabled: true,
      headers: "",
    });
    expect(empty.success).toBe(false);

    const codeOnly = oauthExchangeSchema.safeParse({
      callback_url: "",
      code: "auth_code_only",
      state: "",
      connection_name: "",
      priority: "100",
      enabled: true,
      headers: "",
    });
    expect(codeOnly.success).toBe(false);

    const stateOnly = oauthExchangeSchema.safeParse({
      callback_url: "",
      code: "",
      state: "state_only",
      connection_name: "",
      priority: "100",
      enabled: true,
      headers: "",
    });
    expect(stateOnly.success).toBe(false);
  });

  test("rejects malformed callback_url", () => {
    const parsed = oauthExchangeSchema.safeParse({
      callback_url: "not a url",
      code: "",
      state: "",
      connection_name: "",
      priority: "100",
      enabled: true,
      headers: "",
    });
    expect(parsed.success).toBe(false);
  });
});

describe("updateConnectionSchema", () => {
  test("maps connection summary to form values and omits blank credentials on patch", () => {
    const formValues = connectionToEditFormValues({
      id: "conn_1",
      provider_id: "openai",
      name: "Prod OpenAI",
      protocol: "openai-chat",
      base_url: "https://api.openai.com/v1",
      auth_type: "api_key",
      priority: 10,
      enabled: true,
      models: ["gpt-4o", "o3-mini"],
      headers: { "X-Env": "prod" },
      has_credentials: true,
      token_expires_at: null,
      last_tested_at: null,
      last_test_status: "",
      last_test_error: "",
      created_at: "2026-11-23T00:00:00Z",
      updated_at: "2026-11-23T00:00:00Z",
    });

    const parsed = updateConnectionSchema.safeParse({
      ...formValues,
      priority: "20",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      const req = toUpdateConnectionRequest(parsed.data, "api_key");
      expect(req).toEqual({
        name: "Prod OpenAI",
        base_url: "https://api.openai.com/v1",
        priority: 20,
        enabled: true,
        models: ["gpt-4o", "o3-mini"],
        headers: { "X-Env": "prod" },
      });
      expect("api_key" in req).toBe(false);
    }
  });
});
