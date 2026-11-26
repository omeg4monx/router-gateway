import { describe, expect, test } from "bun:test";
import { pxpipeConfigSchema, pxpipePreviewSchema } from "./schemas";

describe("pxpipe schemas", () => {
  test("accepts valid pipeline config", () => {
    const parsed = pxpipeConfigSchema.safeParse({
      enabled: true,
      auto_start: true,
      strip_null_fields: true,
      collapse_whitespace: true,
      dedupe_system_prompts: true,
      max_history_messages: 20,
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects negative or excessive max_history_messages", () => {
    const neg = pxpipeConfigSchema.safeParse({
      enabled: true,
      auto_start: true,
      strip_null_fields: true,
      collapse_whitespace: true,
      dedupe_system_prompts: true,
      max_history_messages: -1,
    });
    expect(neg.success).toBe(false);

    const high = pxpipeConfigSchema.safeParse({
      enabled: true,
      auto_start: true,
      strip_null_fields: true,
      collapse_whitespace: true,
      dedupe_system_prompts: true,
      max_history_messages: 999,
    });
    expect(high.success).toBe(false);
  });

  test("validates JSON preview payload", () => {
    const ok = pxpipePreviewSchema.safeParse({
      payloadJson: '{"model":"gpt-4o","messages":[{"role":"user","content":"hi"}]}',
    });
    expect(ok.success).toBe(true);

    const bad = pxpipePreviewSchema.safeParse({
      payloadJson: "{not-json",
    });
    expect(bad.success).toBe(false);
  });
});
