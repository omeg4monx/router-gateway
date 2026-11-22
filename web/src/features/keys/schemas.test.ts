import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import { createKeySchema, parseCommaSeparatedList, toCreateApiKeyRequest } from "./schemas";

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
