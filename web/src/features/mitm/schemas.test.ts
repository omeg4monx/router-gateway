import { describe, expect, test } from "bun:test";
import { createMitmAliasSchema, mitmConfigSchema } from "./schemas";

describe("mitm schemas", () => {
  test("accepts valid loopback config", () => {
    const parsed = mitmConfigSchema.safeParse({
      listen_addr: "127.0.0.1:8443",
      target_gateway_url: "http://127.0.0.1:20128",
      enabled_tools: ["antigravity", "copilot"],
      auto_start: false,
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects non-loopback listen_addr and target_gateway_url", () => {
    const badListen = mitmConfigSchema.safeParse({
      listen_addr: "0.0.0.0:8443",
      target_gateway_url: "http://127.0.0.1:20128",
      enabled_tools: ["copilot"],
      auto_start: false,
    });
    expect(badListen.success).toBe(false);

    const badTarget = mitmConfigSchema.safeParse({
      listen_addr: "127.0.0.1:8443",
      target_gateway_url: "http://10.0.0.5:20128",
      enabled_tools: ["copilot"],
      auto_start: false,
    });
    expect(badTarget.success).toBe(false);
  });

  test("validates model alias creation", () => {
    const ok = createMitmAliasSchema.safeParse({
      tool: "copilot",
      source_model: "gpt-4o",
      target_model: "claude-sonnet-4",
      enabled: true,
    });
    expect(ok.success).toBe(true);

    const bad = createMitmAliasSchema.safeParse({
      tool: "copilot",
      source_model: "",
      target_model: "claude-sonnet-4",
      enabled: true,
    });
    expect(bad.success).toBe(false);
  });
});
