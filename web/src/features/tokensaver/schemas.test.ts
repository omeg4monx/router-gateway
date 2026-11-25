import { describe, expect, test } from "bun:test";
import {
  HISTORY_LIMIT_PRESETS,
  matchHistoryPreset,
  TOKEN_SAVER_LEVELS,
  TOKEN_SAVER_STRATEGIES,
  tokenSaverConfigSchema,
} from "./schemas";

describe("tokensaver schemas", () => {
  test("validates valid token saver feature toggles and backend-managed bounds", () => {
    const parsed = tokenSaverConfigSchema.safeParse({
      enabled: true,
      strategy: "aggressive",
      max_history_messages: 20,
      max_tool_output_chars: 8000,
      headroom_enabled: true,
      prefix_depth: 2,
      max_cache_entries: 1024,
      cache_ttl_seconds: 1800,
      inject_cache_control: true,
      caveman_enabled: true,
      caveman_level: "full",
      ponytail_enabled: true,
      ponytail_level: "ultra",
      trim_history_enabled: true,
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects out-of-range bounds or unknown strategy/level", () => {
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
      caveman_enabled: true,
      caveman_level: "extreme",
      ponytail_enabled: false,
      ponytail_level: "full",
      trim_history_enabled: false,
    });
    expect(parsed.success).toBe(false);
  });

  test("matches xrouter-style history limit presets and exposes canonical levels", () => {
    expect(HISTORY_LIMIT_PRESETS.length).toBe(3);
    expect(matchHistoryPreset(12)).toBe("tight");
    expect(matchHistoryPreset(20)).toBe("balanced");
    expect(matchHistoryPreset(40)).toBe("deep");
    expect(matchHistoryPreset(27)).toBe("custom");
    expect(TOKEN_SAVER_STRATEGIES).toEqual(["off", "whitespace", "dedupe", "aggressive"]);
    expect(TOKEN_SAVER_LEVELS).toEqual(["lite", "full", "ultra"]);
  });
});
