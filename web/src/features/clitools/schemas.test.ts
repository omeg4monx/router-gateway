import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  CLI_TOOL_IDS,
  DEFAULT_GENERATE_CLI_CONFIG_VALUES,
  generateCliConfigSchema,
  MAX_CLI_INPUT_RUNES,
  toCLIToolConfigRequest,
} from "./schemas";

describe("generateCliConfigSchema", () => {
  test("validates http:// and https:// base_url, trims optional fields, and works with zodResolver", async () => {
    const input = {
      base_url: "  http://127.0.0.1:20128/  ",
      api_key: "  om_test_secret_123  ",
      model: "  anthropic/claude-sonnet-4  ",
    };

    const parsed = generateCliConfigSchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data).toEqual({
        base_url: "http://127.0.0.1:20128",
        api_key: "om_test_secret_123",
        model: "anthropic/claude-sonnet-4",
      });
      expect(toCLIToolConfigRequest(parsed.data)).toEqual({
        base_url: "http://127.0.0.1:20128",
        api_key: "om_test_secret_123",
        model: "anthropic/claude-sonnet-4",
      });
    }

    const resolver = zodResolver(generateCliConfigSchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
    expect(resolved.values).toEqual({
      base_url: "http://127.0.0.1:20128",
      api_key: "om_test_secret_123",
      model: "anthropic/claude-sonnet-4",
    });
  });

  test("accepts https:// URLs with path prefixes and omits blank optional api_key/model", () => {
    const parsed = generateCliConfigSchema.safeParse({
      base_url: "https://gateway.example.com/router/",
      api_key: "   ",
      model: "",
    });
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data).toEqual({
        base_url: "https://gateway.example.com/router",
        api_key: "",
        model: "",
      });
      expect(toCLIToolConfigRequest(parsed.data)).toEqual({
        base_url: "https://gateway.example.com/router",
      });
    }

    const minimal = generateCliConfigSchema.safeParse({
      base_url: "https://router.local:8443",
    });
    expect(minimal.success).toBe(true);
    if (minimal.success) {
      expect(toCLIToolConfigRequest(minimal.data)).toEqual({
        base_url: "https://router.local:8443",
      });
    }
  });

  test("rejects missing, malformed, non-HTTP(S), control-char, or oversized base_url", () => {
    const badUrls = [
      "",
      "   ",
      "not-a-url",
      "127.0.0.1:20128",
      "ftp://127.0.0.1:20128",
      "javascript:alert(1)",
      "http://:8080",
      "https://",
      "http://127.0.0.1:20128\nfoo",
      `https://example.com/${"a".repeat(MAX_CLI_INPUT_RUNES)}`,
    ];

    for (const badUrl of badUrls) {
      const res = generateCliConfigSchema.safeParse({
        ...DEFAULT_GENERATE_CLI_CONFIG_VALUES,
        base_url: badUrl,
      });
      expect(res.success).toBe(false);
    }
  });

  test("rejects control characters or oversized values in optional api_key and model", () => {
    const oversized = "x".repeat(MAX_CLI_INPUT_RUNES + 1);

    for (const bad of ["abc\ndef", "abc\tdef", oversized]) {
      expect(
        generateCliConfigSchema.safeParse({
          base_url: "http://127.0.0.1:20128",
          api_key: bad,
          model: "",
        }).success,
      ).toBe(false);

      expect(
        generateCliConfigSchema.safeParse({
          base_url: "http://127.0.0.1:20128",
          api_key: "",
          model: bad,
        }).success,
      ).toBe(false);
    }
  });
});

describe("CLI_TOOL_IDS", () => {
  test("contains all six Wave 1 supported CLI tools", () => {
    expect(CLI_TOOL_IDS).toEqual([
      "claude-code",
      "codex-cli",
      "opencode",
      "cline",
      "roo-code",
      "droid",
    ]);
  });
});
