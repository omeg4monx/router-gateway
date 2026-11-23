import { describe, expect, test } from "bun:test";
import { parsePreviewBodyJson, tokenSaverConfigSchema } from "./schemas";

describe("tokensaver schemas", () => {
  test("validates valid token saver and headroom config bounds", () => {
    const parsed = tokenSaverConfigSchema.safeParse({
      enabled: true,
      strategy: "aggressive",
      max_history_messages: 40,
      max_tool_output_chars: 8000,
      headroom_enabled: true,
      prefix_depth: 2,
      max_cache_entries: 1024,
      cache_ttl_seconds: 1800,
      inject_cache_control: true,
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects out-of-range bounds or unknown strategy", () => {
    const parsed = tokenSaverConfigSchema.safeParse({
      enabled: true,
      strategy: "quantum",
      max_history_messages: 2,
      max_tool_output_chars: 10,
      headroom_enabled: true,
      prefix_depth: 20,
      max_cache_entries: 2,
      cache_ttl_seconds: 10,
      inject_cache_control: false,
    });
    expect(parsed.success).toBe(false);
  });

  test("validates preview JSON body object", () => {
    expect(parsePreviewBodyJson('{"messages":[]}').ok).toBe(true);
    expect(parsePreviewBodyJson("[]").ok).toBe(false);
    expect(parsePreviewBodyJson("{invalid").ok).toBe(false);
  });
});
