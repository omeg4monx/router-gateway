import { describe, it, expect } from "vitest";
import fs from "node:fs";
import path from "node:path";
import { parseIntervalMs } from "@/shared/utils/interval";

const ROOT = path.resolve(process.cwd());
const readRel = (rel) => fs.readFileSync(path.join(ROOT, rel), "utf8");
import {
  RESET_INTERVAL_OPTIONS,
  BULK_RESET_INTERVAL_OPTIONS,
  formatTokensNumber,
  maskApiKey,
  normalizeOrigin,
  getExpiryMeta,
  isKeyExhausted,
  filterApiKeys,
  countKeyBuckets,
  resolveDefaultSnippetModel,
  countWhitelistEntries,
  validateCustomInterval,
} from "@/app/(dashboard)/dashboard/endpoint/endpointUtils";
import { generateEndpointSnippet } from "@/app/(dashboard)/dashboard/endpoint/endpointSnippets";

const NOW = new Date("2026-11-22T12:00:00.000Z").getTime();

describe("endpointUtils", () => {
  it("exposes valid parseIntervalMs options for BulkEditKeys and Create/Edit forms", () => {
    expect(BULK_RESET_INTERVAL_OPTIONS.length).toBeGreaterThanOrEqual(6);
    for (const opt of BULK_RESET_INTERVAL_OPTIONS) {
      if (opt.value === "never") continue;
      expect(parseIntervalMs(opt.value), `invalid bulk interval "${opt.value}"`).toBeGreaterThan(0);
    }
    const resetValues = RESET_INTERVAL_OPTIONS.map((o) => o.value);
    expect(resetValues).toContain("never");
    expect(resetValues).toContain("custom");
    for (const val of resetValues) {
      if (val === "never" || val === "custom") continue;
      expect(parseIntervalMs(val), `invalid reset interval "${val}"`).toBeGreaterThan(0);
    }
  });

  it("formats token numbers cleanly", () => {
    expect(formatTokensNumber(0)).toBe("0");
    expect(formatTokensNumber(-5)).toBe("0");
    expect(formatTokensNumber(850)).toBe("850");
    expect(formatTokensNumber(12_500)).toBe("12.5K");
    expect(formatTokensNumber(3_250_000)).toBe("3.25M");
    expect(formatTokensNumber(2_500_000_000)).toBe("2.50B");
  });

  it("masks API keys safely while preserving prefix and suffix", () => {
    expect(maskApiKey("")).toBe("");
    expect(maskApiKey("short")).toBe("short");
    const masked = maskApiKey("sk-omeg4mon-secret-1234");
    expect(masked.startsWith("sk-ome")).toBe(true);
    expect(masked.endsWith("1234")).toBe(true);
    expect(masked).toContain("•");
    expect(masked.length).toBe("sk-omeg4mon-secret-1234".length);
  });

  it("normalizes origin URLs by stripping trailing slashes and /v1", () => {
    expect(normalizeOrigin("http://localhost:20128/v1/")).toBe("http://localhost:20128");
    expect(normalizeOrigin("http://localhost:20128/v1")).toBe("http://localhost:20128");
    expect(normalizeOrigin("https://demo.trycloudflare.com///")).toBe("https://demo.trycloudflare.com");
    expect(normalizeOrigin("")).toBe("");
  });

  it("computes expiry metadata safely for empty, invalid, expired, soon, and future dates", () => {
    expect(getExpiryMeta(null, NOW)).toEqual({ status: "none", label: "" });
    expect(getExpiryMeta("", NOW)).toEqual({ status: "none", label: "" });
    expect(getExpiryMeta("not-a-date", NOW)).toEqual({ status: "none", label: "" });

    const expired = getExpiryMeta("2026-11-20T00:00:00.000Z", NOW);
    expect(expired).toEqual({ status: "expired", label: "Expired" });

    const soon = getExpiryMeta("2026-11-25T12:00:00.000Z", NOW);
    expect(soon.status).toBe("soon");
    expect(soon.label).toContain("3d");
    expect(soon.label).not.toMatch(/NaN|Invalid/i);

    const later = getExpiryMeta("2027-01-15T12:00:00.000Z", NOW);
    expect(later.status).toBe("ok");
    expect(later.label).toMatch(/^Expires /);
    expect(later.label).not.toMatch(/NaN|Invalid/i);
  });

  it("partitions and filters API keys across all status buckets and search query", () => {
    const sampleKeys = [
      { id: "1", name: "Prod Active", key: "sk-prod-1111", isActive: true, tokenLimit: 1000, usedTokens: 200 },
      { id: "2", name: "Disabled Key", key: "sk-dis-2222", isActive: false, tokenLimit: 0, usedTokens: 0 },
      { id: "3", name: "Exhausted Key", key: "sk-exh-3333", isActive: true, tokenLimit: 500, usedTokens: 500 },
      { id: "4", name: "Expiring Soon", key: "sk-soon-4444", isActive: true, tokenLimit: 0, usedTokens: 10, expiresAt: "2026-11-24T12:00:00.000Z" },
      { id: "5", name: "Already Expired", key: "sk-past-5555", isActive: true, tokenLimit: 0, usedTokens: 0, expiresAt: "2026-11-01T00:00:00.000Z" },
    ];

    expect(isKeyExhausted(sampleKeys[0])).toBe(false);
    expect(isKeyExhausted(sampleKeys[2])).toBe(true);

    expect(countKeyBuckets(sampleKeys, NOW)).toEqual({
      all: 5,
      active: 4,
      disabled: 1,
      exhausted: 1,
      expiring: 2,
    });

    expect(filterApiKeys(sampleKeys, { statusFilter: "disabled", nowMs: NOW }).map((k) => k.id)).toEqual(["2"]);
    expect(filterApiKeys(sampleKeys, { statusFilter: "exhausted", nowMs: NOW }).map((k) => k.id)).toEqual(["3"]);
    expect(filterApiKeys(sampleKeys, { statusFilter: "expiring", nowMs: NOW }).map((k) => k.id)).toEqual(["4", "5"]);
    expect(filterApiKeys(sampleKeys, { statusFilter: "all", query: "soon", nowMs: NOW }).map((k) => k.id)).toEqual(["4"]);
    expect(filterApiKeys(sampleKeys, { statusFilter: "active", query: "sk-exh", nowMs: NOW }).map((k) => k.id)).toEqual(["3"]);
  });

  it("resolves default snippet model from aliases, providers, or safe fallback", () => {
    expect(resolveDefaultSnippetModel({ "claude-sonnet": "anthropic/claude-sonnet-4" }, [])).toBe("claude-sonnet");
    expect(resolveDefaultSnippetModel({}, [{ provider: "openai", models: ["gpt-4.1"] }])).toBe("gpt-4.1");
    expect(resolveDefaultSnippetModel({}, [])).toBe("gpt-4o");
  });

  it("counts IP whitelist entries and validates custom reset intervals", () => {
    expect(countWhitelistEntries("")).toBe(0);
    expect(countWhitelistEntries(" 192.168.1.1 , 10.0.0.0/24 , ")).toBe(2);

    expect(validateCustomInterval("10h")).toBeNull();
    expect(validateCustomInterval("3d")).toBeNull();
    expect(validateCustomInterval("30m")).toBeNull();
    expect(validateCustomInterval("")).toBeTruthy();
    expect(validateCustomInterval("daily")).toBeTruthy();
    expect(validateCustomInterval("0h")).toBeTruthy();
    expect(validateCustomInterval("abc")).toBeTruthy();
  });
});

describe("endpointSnippets", () => {
  it("generates OpenAI snippets across curl, python, node, and go without double /v1", () => {
    for (const lang of ["curl", "python", "node", "go"]) {
      const code = generateEndpointSnippet({
        format: "openai",
        lang,
        baseUrl: "http://localhost:20128/v1/",
        apiKey: "sk-test-key",
        model: "gpt-4o",
      });
      expect(code).toContain("sk-test-key");
      expect(code).toContain("gpt-4o");
      expect(code).not.toContain("/v1/v1");
      expect(code).not.toContain("20128//v1");
    }
    const curl = generateEndpointSnippet({
      format: "openai",
      lang: "curl",
      baseUrl: "http://localhost:20128/v1/",
      apiKey: "",
      model: "",
    });
    expect(curl).toContain("http://localhost:20128/v1/chat/completions");
    expect(curl).toContain("Authorization: Bearer <YOUR_API_KEY>");
    expect(curl).toContain('"model":"gpt-4o"');
  });

  it("generates Anthropic snippets targeting /v1/messages with x-api-key across all languages", () => {
    for (const lang of ["curl", "python", "node", "go"]) {
      const code = generateEndpointSnippet({
        format: "anthropic",
        lang,
        baseUrl: "https://demo.trycloudflare.com/v1/",
        apiKey: "sk-claude-key",
        model: "claude-sonnet-4-5",
      });
      expect(code).toContain("sk-claude-key");
      expect(code).toContain("claude-sonnet-4-5");
      expect(code).not.toContain("/v1/v1");
    }
    const curl = generateEndpointSnippet({
      format: "anthropic",
      lang: "curl",
      baseUrl: "https://demo.trycloudflare.com/v1/",
      apiKey: "sk-claude-key",
      model: "claude-sonnet-4-5",
    });
    expect(curl).toContain("https://demo.trycloudflare.com/v1/messages");
    expect(curl).toContain("x-api-key: sk-claude-key");
    expect(curl).toContain("anthropic-version: 2023-06-01");
  });
});

describe("BulkEditKeys and KeyShareModal UX contracts", () => {
  it("uses valid parseIntervalMs reset intervals and sticky Modal footer in BulkEditKeys", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/BulkEditKeys.js");
    expect(src).toContain("BULK_RESET_INTERVAL_OPTIONS");
    for (const broken of ['"hourly"', '"daily"', '"weekly"', '"monthly"']) {
      expect(src).not.toContain(broken);
    }
    expect(src).toContain("footer={");
    expect(src).toContain("accent-[var(--color-primary)]");
    expect(src).toContain("selectedCount");
  });

  it("uses sticky Modal footer and soft rounded-xl list surface in KeyShareModal", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/KeyShareModal.js");
    expect(src).toContain("footer={");
    expect(src).toContain("rounded-xl");
  });
});

describe("KeyFormModal and KeySnippetModal UX contracts", () => {
  it("implements a 3-tab KeyFormModal with sticky footer, inline validation, English copy, and permission lock", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/KeyFormModal.js");
    expect(src).toContain("SegmentedControl");
    expect(src).toContain('"general"');
    expect(src).toContain('"limits"');
    expect(src).toContain('"permissions"');
    expect(src).toContain("footer={");
    expect(src).toContain("validateCustomInterval");
    expect(src).toContain('role="alert"');
    expect(src).toContain("isApiKeyUser");
    expect(src).not.toMatch(/\balert\s*\(/);
    expect(src).not.toContain("Maksimal");
    expect(src).not.toContain("Saat ini");
  });

  it("implements KeySnippetModal with OpenAI/Anthropic switcher, copy confirmation, and sticky footer", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/KeySnippetModal.js");
    expect(src).toContain("generateEndpointSnippet");
    expect(src).toContain('"openai"');
    expect(src).toContain('"anthropic"');
    expect(src).toContain("useCopyToClipboard");
    expect(src).toContain("Copied!");
    expect(src).toContain("footer={");
  });
});

describe("GatewayConsole UX contracts", () => {
  it("renders inline StatusAlert for tunnel and tailscale statuses and 1-click auth gate recovery", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/GatewayConsole.js");
    expect(src).toContain("StatusAlert");
    expect(src).toContain("tunnel.status");
    expect(src).toContain("tailscale.status");
    expect(src).toContain("Enable Auth & Continue");
    expect(src).toContain("onEnableAuthAndTunnel");
  });

  it("synchronizes OpenAI/Claude URLs with safe Local fallback, smooth accordion snippets, and admin gating", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/GatewayConsole.js");
    expect(src).toContain("localOrigin");
    expect(src).toContain("Showing Local URL");
    expect(src).toContain("accordion-grid");
    expect(src).toContain("generateEndpointSnippet");
    expect(src).toContain('"openai"');
    expect(src).toContain('"anthropic"');
    expect(src).toContain("!isApiKeyUser");
    expect(src).toContain("Require API Key Authentication");
  });
});

describe("ApiKeyCard UX contracts", () => {
  it("renders expiry chip, admin-gated checkbox, own-key guard, calm quota bar without pulse/ping or duplicate Tokens chip", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/ApiKeyCard.js");
    expect(src).toContain("getExpiryMeta");
    expect(src).toContain("canBulkSelect");
    expect(src).toContain("isOwnKey");
    expect(src).not.toContain("animate-pulse");
    expect(src).not.toContain("animate-ping");
    expect(src).not.toContain("Tokens:");
  });

  it("isolates destructive Delete into an accessible overflow menu and keeps 32px+ quick actions", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/components/ApiKeyCard.js");
    expect(src).toContain("more_vert");
    expect(src).toContain("Escape");
    expect(src).toContain("onOpenSnippet");
    expect(src).toContain("size-8");
    expect(src).toContain("Delete Key");
  });
});

describe("EndpointPageClient integration and bug-fix contracts", () => {
  it("removes loadSettings ReferenceError, eliminates alert(), and uses real /api/health probe", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/EndpointPageClient.js");
    expect(src).not.toContain("setRequireLogin");
    expect(src).not.toContain("setHasPassword");
    expect(src).not.toMatch(/\balert\s*\(/);
    expect(src).toContain("useNotificationStore");
    expect(src).not.toContain("ONLINE · 200 OK");
    expect(src).toContain("/api/health");
    expect(src).toContain("performance.now()");
    expect(src).not.toContain("Klik tombol di bawah ini");
  });

  it("wires GatewayConsole, ApiKeyCard, KeyFormModal, KeySnippetModal, admin Select All, status filters, and Copy cURL on Created Key", () => {
    const src = readRel("src/app/(dashboard)/dashboard/endpoint/EndpointPageClient.js");
    expect(src).toContain("<GatewayConsole");
    expect(src).toContain("<ApiKeyCard");
    expect(src).toContain("<KeyFormModal");
    expect(src).toContain("<KeySnippetModal");
    expect(src).toContain("<BulkEditKeys");
    expect(src).toContain("selectedKeyIds");
    expect(src).toContain("toggleSelectAllKeys");
    expect(src).toContain("statusFilter");
    expect(src).not.toContain("shown only once");
    expect(src).toContain("Copy cURL");
  });
});
