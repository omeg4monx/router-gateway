import { describe, expect, test } from "bun:test";
import {
  buildTestFingerprint,
  CUSTOM_ENDPOINT_PRESETS,
  parseBulkCredentials,
  parseOAuthCallbackUrl,
  slugifyProviderId,
} from "./modalUtils";
import { MAX_ID_RUNES } from "./schemas";

describe("CUSTOM_ENDPOINT_PRESETS", () => {
  test("includes ollama, vllm, litellm, and lmstudio presets with valid http base URLs", () => {
    const ids = CUSTOM_ENDPOINT_PRESETS.map((p) => p.id);
    expect(ids).toEqual(["ollama", "vllm", "litellm", "lmstudio"]);
    for (const preset of CUSTOM_ENDPOINT_PRESETS) {
      expect(preset.baseUrl.startsWith("http://") || preset.baseUrl.startsWith("https://")).toBe(
        true,
      );
      expect(preset.protocol).toBe("openai-chat");
    }
  });
});

describe("slugifyProviderId", () => {
  test("normalizes display names into valid provider node IDs", () => {
    expect(slugifyProviderId("Local vLLM Cluster")).toBe("local-vllm-cluster");
    expect(slugifyProviderId("  Ollama (Home Lab) #1!  ")).toBe("ollama-home-lab-1");
    expect(slugifyProviderId("my_custom--node")).toBe("my_custom-node");
    expect(slugifyProviderId("   ")).toBe("");
  });

  test("caps slug length to MAX_ID_RUNES", () => {
    const longName = "a".repeat(MAX_ID_RUNES + 25);
    const slug = slugifyProviderId(longName);
    expect(slug.length).toBeLessThanOrEqual(MAX_ID_RUNES);
  });
});

describe("parseOAuthCallbackUrl", () => {
  test("extracts code and state from query parameters", () => {
    const parsed = parseOAuthCallbackUrl(
      "http://localhost:20128/oauth/callback?code=auth_123&state=pkce_456",
    );
    expect(parsed).toEqual({
      valid: true,
      code: "auth_123",
      state: "pkce_456",
    });
  });

  test("extracts code and state from hash fragments", () => {
    const parsed = parseOAuthCallbackUrl(
      "https://router.local/callback#code=frag_code&state=frag_state",
    );
    expect(parsed).toEqual({
      valid: true,
      code: "frag_code",
      state: "frag_state",
    });
  });

  test("returns invalid when URL is malformed or missing code/state", () => {
    expect(parseOAuthCallbackUrl("not a url").valid).toBe(false);
    expect(parseOAuthCallbackUrl("http://localhost:20128/callback?code=only").valid).toBe(false);
    expect(parseOAuthCallbackUrl("").valid).toBe(false);
  });
});

describe("parseBulkCredentials", () => {
  test("parses named and raw credential lines while counting invalid empty-key lines", () => {
    const input = ["prod-1|sk-111", "", "  sk-raw-222  ", "missing-key|", "prod-3 | sk-333"].join(
      "\n",
    );

    const result = parseBulkCredentials(input, "OpenAI");
    expect(result.entries).toEqual([
      { name: "prod-1", apiKey: "sk-111" },
      { name: "OpenAI #2", apiKey: "sk-raw-222" },
      { name: "prod-3", apiKey: "sk-333" },
    ]);
    expect(result.invalidLines).toBe(1);
  });
});

describe("buildTestFingerprint", () => {
  test("produces deterministic fingerprints and detects any connection config change", () => {
    const base = {
      providerId: "custom",
      protocol: "openai-chat",
      baseUrl: "http://127.0.0.1:8000/v1",
      apiKey: "sk-test",
      headers: "X-A: 1",
    };
    expect(buildTestFingerprint(base)).toBe(buildTestFingerprint({ ...base }));
    expect(buildTestFingerprint({ ...base, baseUrl: "http://127.0.0.1:8001/v1" })).not.toBe(
      buildTestFingerprint(base),
    );
    expect(buildTestFingerprint({ ...base, apiKey: "sk-other" })).not.toBe(
      buildTestFingerprint(base),
    );
    expect(buildTestFingerprint({ ...base, protocol: "anthropic" })).not.toBe(
      buildTestFingerprint(base),
    );
    expect(buildTestFingerprint({ ...base, headers: "X-A: 2" })).not.toBe(
      buildTestFingerprint(base),
    );
  });
});
