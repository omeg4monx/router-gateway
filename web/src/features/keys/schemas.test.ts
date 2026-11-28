import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  buildGatewayUrls,
  countKeyStatusBuckets,
  filterKeysByStatusAndSearch,
  generateKeySnippet,
  getKeyExpiryState,
  normalizeOrigin,
  resolveGatewayOrigins,
} from "./endpointUtils";
import {
  type ApiKeySummary,
  createKeySchema,
  parseCommaSeparatedList,
  toCreateApiKeyRequest,
} from "./schemas";

describe("parseCommaSeparatedList", () => {
  test("splits, trims, drops empty entries, and deduplicates preserving order", () => {
    expect(
      parseCommaSeparatedList(" gpt-4o, claude-3-7-sonnet , , gpt-4o, gemini-2.5-pro "),
    ).toEqual(["gpt-4o", "claude-3-7-sonnet", "gemini-2.5-pro"]);
  });

  test("returns an empty array for blank or comma-only input", () => {
    expect(parseCommaSeparatedList("")).toEqual([]);
    expect(parseCommaSeparatedList(" , ,  ")).toEqual([]);
  });
});

describe("createKeySchema", () => {
  test("parses comma-separated allowed_models and tags and coerces valid limits", async () => {
    const input = {
      name: "  ci-runner-key  ",
      allowed_models: "gpt-4o, claude-sonnet-4, gpt-4o",
      tags: "ci, production, ",
      rate_limit_rpm: "120",
      quota_usd: "45.5",
      expires_at: "",
    };

    const parsed = createKeySchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data).toEqual({
        name: "ci-runner-key",
        allowed_models: ["gpt-4o", "claude-sonnet-4"],
        tags: ["ci", "production"],
        rate_limit_rpm: 120,
        quota_usd: 45.5,
        expires_at: null,
      });

      const req = toCreateApiKeyRequest(parsed.data);
      expect(req).toEqual({
        name: "ci-runner-key",
        allowed_models: ["gpt-4o", "claude-sonnet-4"],
        scopes: ["gateway"],
        tags: ["ci", "production"],
        rate_limit_rpm: 120,
        quota_usd: 45.5,
        expires_at: null,
      });
    }

    const resolver = zodResolver(createKeySchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
    expect(resolved.values).toEqual({
      name: "ci-runner-key",
      allowed_models: ["gpt-4o", "claude-sonnet-4"],
      tags: ["ci", "production"],
      rate_limit_rpm: 120,
      quota_usd: 45.5,
      expires_at: null,
    });
  });

  test("defaults empty numeric fields to zero unlimited bounds", () => {
    const parsed = createKeySchema.safeParse({
      name: "local-dev",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "",
      quota_usd: "",
      expires_at: "",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data.rate_limit_rpm).toBe(0);
      expect(parsed.data.quota_usd).toBe(0);
      expect(parsed.data.allowed_models).toEqual([]);
      expect(parsed.data.tags).toEqual([]);
      expect(parsed.data.expires_at).toBeNull();
    }
  });

  test("rejects blank or whitespace-only key name", () => {
    const empty = createKeySchema.safeParse({
      name: "   ",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "0",
      quota_usd: "0",
      expires_at: "",
    });
    expect(empty.success).toBe(false);
    if (!empty.success) {
      expect(empty.error.issues[0]?.path).toEqual(["name"]);
    }
  });

  test("rejects key name longer than 128 runes", () => {
    const parsed = createKeySchema.safeParse({
      name: "k".repeat(129),
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "0",
      quota_usd: "0",
      expires_at: "",
    });
    expect(parsed.success).toBe(false);
  });

  test("rejects negative or fractional rate_limit_rpm", () => {
    const negative = createKeySchema.safeParse({
      name: "neg-rpm",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "-1",
      quota_usd: "0",
      expires_at: "",
    });
    expect(negative.success).toBe(false);

    const fractional = createKeySchema.safeParse({
      name: "float-rpm",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "12.5",
      quota_usd: "0",
      expires_at: "",
    });
    expect(fractional.success).toBe(false);
  });

  test("rejects negative or non-numeric quota_usd", () => {
    const negative = createKeySchema.safeParse({
      name: "neg-quota",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "0",
      quota_usd: "-0.01",
      expires_at: "",
    });
    expect(negative.success).toBe(false);

    const invalid = createKeySchema.safeParse({
      name: "bad-quota",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "0",
      quota_usd: "abc",
      expires_at: "",
    });
    expect(invalid.success).toBe(false);
  });

  test("normalizes valid expires_at timestamp to ISO-8601 UTC and rejects invalid dates", () => {
    const valid = createKeySchema.safeParse({
      name: "expiring-key",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "60",
      quota_usd: "10",
      expires_at: "2027-01-15T12:30:00Z",
    });
    expect(valid.success).toBe(true);
    if (valid.success) {
      expect(valid.data.expires_at).toBe("2027-01-15T12:30:00.000Z");
    }

    const invalid = createKeySchema.safeParse({
      name: "bad-expiry",
      allowed_models: "",
      tags: "",
      rate_limit_rpm: "0",
      quota_usd: "0",
      expires_at: "not-a-timestamp",
    });
    expect(invalid.success).toBe(false);
  });
});

describe("endpointUtils (main Endpoint & Keys parity)", () => {
  test("normalizes origins and builds synchronized OpenAI & Claude base URLs", () => {
    expect(normalizeOrigin("http://127.0.0.1:20128///")).toBe("http://127.0.0.1:20128");
    expect(normalizeOrigin("  https://gw.example.com/ ")).toBe("https://gw.example.com");
    expect(buildGatewayUrls("https://gw.example.com/")).toEqual({
      origin: "https://gw.example.com",
      openaiBaseUrl: "https://gw.example.com/v1",
      claudeBaseUrl: "https://gw.example.com",
    });
  });

  test("resolves server domain host separately while keeping localhost as localhost", () => {
    expect(resolveGatewayOrigins("https://router.example.com/")).toEqual({
      domainOrigin: "https://router.example.com",
      localOrigin: "http://localhost:20128",
      defaultRoute: "domain",
      isDomainHost: true,
    });

    expect(resolveGatewayOrigins("http://api.company.id:8080")).toEqual({
      domainOrigin: "http://api.company.id:8080",
      localOrigin: "http://localhost:20128",
      defaultRoute: "domain",
      isDomainHost: true,
    });

    expect(resolveGatewayOrigins("http://localhost:20128")).toEqual({
      domainOrigin: null,
      localOrigin: "http://localhost:20128",
      defaultRoute: "local",
      isDomainHost: false,
    });

    expect(resolveGatewayOrigins("http://127.0.0.1:20128")).toEqual({
      domainOrigin: null,
      localOrigin: "http://127.0.0.1:20128",
      defaultRoute: "local",
      isDomainHost: false,
    });

    expect(resolveGatewayOrigins("http://[::1]:20128")).toEqual({
      domainOrigin: null,
      localOrigin: "http://[::1]:20128",
      defaultRoute: "local",
      isDomainHost: false,
    });

    expect(resolveGatewayOrigins("")).toEqual({
      domainOrigin: null,
      localOrigin: "http://127.0.0.1:20128",
      defaultRoute: "local",
      isDomainHost: false,
    });
  });

  test("generates OpenAI and Anthropic snippets across curl, python, node, and go", () => {
    const oaiCurl = generateKeySnippet({
      format: "openai",
      lang: "curl",
      baseUrl: "http://127.0.0.1:20128/",
      apiKey: "sk-omg-secret",
      model: "gpt-4o",
    });
    expect(oaiCurl).toContain("http://127.0.0.1:20128/v1/chat/completions");
    expect(oaiCurl).toContain("Authorization: Bearer sk-omg-secret");

    const claudeGo = generateKeySnippet({
      format: "anthropic",
      lang: "go",
      baseUrl: "https://gw.example.com",
      apiKey: "sk-omg-secret",
      model: "claude-sonnet-4-5",
    });
    expect(claudeGo).toContain("https://gw.example.com/v1/messages");
    expect(claudeGo).toContain('req.Header.Set("x-api-key", "sk-omg-secret")');
  });

  test("filters and counts active, disabled, and expired keys accurately", () => {
    const now = Date.parse("2026-11-25T00:00:00Z");
    const sample: ApiKeySummary[] = [
      {
        id: "k1",
        name: "Production CLI",
        key_prefix: "sk-omg_aaa",
        allowed_models: [],
        scopes: ["gateway"],
        tags: [],
        rate_limit_rpm: 60,
        quota_usd: 10,
        disabled: false,
        expires_at: null,
        last_used_at: null,
        created_at: "2026-11-20T00:00:00Z",
        updated_at: "2026-11-20T00:00:00Z",
      },
      {
        id: "k2",
        name: "Paused Dev",
        key_prefix: "sk-omg_bbb",
        allowed_models: ["gpt-4o"],
        scopes: ["gateway"],
        tags: [],
        rate_limit_rpm: 0,
        quota_usd: 0,
        disabled: true,
        expires_at: null,
        last_used_at: null,
        created_at: "2026-11-20T00:00:00Z",
        updated_at: "2026-11-20T00:00:00Z",
      },
      {
        id: "k3",
        name: "Expired Trial",
        key_prefix: "sk-omg_ccc",
        allowed_models: [],
        scopes: ["gateway"],
        tags: [],
        rate_limit_rpm: 0,
        quota_usd: 0,
        disabled: false,
        expires_at: "2026-11-24T00:00:00Z",
        last_used_at: null,
        created_at: "2026-11-20T00:00:00Z",
        updated_at: "2026-11-20T00:00:00Z",
      },
    ];

    expect(countKeyStatusBuckets(sample, now)).toEqual({
      all: 3,
      active: 1,
      disabled: 1,
      expired: 1,
    });
    expect(filterKeysByStatusAndSearch(sample, "active", "", now).map((k) => k.id)).toEqual(["k1"]);
    expect(filterKeysByStatusAndSearch(sample, "disabled", "", now).map((k) => k.id)).toEqual([
      "k2",
    ]);
    expect(filterKeysByStatusAndSearch(sample, "expired", "", now).map((k) => k.id)).toEqual([
      "k3",
    ]);
    expect(filterKeysByStatusAndSearch(sample, "all", "sk-omg_bbb", now).map((k) => k.id)).toEqual([
      "k2",
    ]);
    expect(getKeyExpiryState("2026-11-24T00:00:00Z", now).status).toBe("expired");
    expect(getKeyExpiryState("2026-11-27T00:00:00Z", now).status).toBe("soon");
    expect(getKeyExpiryState(null, now).status).toBe("none");
  });
});

describe("KeysConsole main Endpoint & Keys parity contract", () => {
  test("keeps main Endpoint & Keys zones while excluding raw developer chrome", () => {
    const consoleSource = readFileSync(resolve(import.meta.dirname, "KeysConsole.tsx"), "utf8");
    const modalsSource = readFileSync(resolve(import.meta.dirname, "KeyModals.tsx"), "utf8");
    const gatewaySource = readFileSync(
      resolve(import.meta.dirname, "GatewayEndpointBar.tsx"),
      "utf8",
    );

    for (const src of [consoleSource, modalsSource, gatewaySource]) {
      expect(src).not.toContain("POST /api/keys");
      expect(src).not.toContain("GET /api/keys");
      expect(src).not.toContain("Gateway Credentials");
      expect(src).not.toContain("SHA-256");
      expect(src).not.toContain("key-tags");
    }

    // Main parity: Endpoint & Keys header, live health probe, route/URL console, search/filter, edit/snippet modals
    expect(consoleSource).toContain("GatewayEndpointBar");
    expect(consoleSource).toContain("KeyFormModal");
    expect(consoleSource).toContain("KeySnippetModal");
    expect(consoleSource).toContain("useUIStore");
    expect(consoleSource).toContain('role="switch"');
    expect(gatewaySource).toContain("openaiBaseUrl");
    expect(gatewaySource).toContain("claudeBaseUrl");
    expect(modalsSource).toContain('role="dialog"');
    expect(modalsSource).toContain("clearTimeout");
  });
});
