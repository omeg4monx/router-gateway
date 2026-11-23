import { describe, expect, test } from "bun:test";
import { tunnelConfigSchema } from "./schemas";

describe("tunnel schemas", () => {
  test("accepts valid loopback cloudflare and tailscale tunnel configs", () => {
    const cf = tunnelConfigSchema.safeParse({
      provider: "cloudflare",
      local_url: "http://127.0.0.1:20128",
      tailscale_mode: "funnel",
      auto_start: false,
    });
    expect(cf.success).toBe(true);

    const ts = tunnelConfigSchema.safeParse({
      provider: "tailscale",
      local_url: "http://localhost:20128",
      tailscale_mode: "serve",
      auto_start: true,
    });
    expect(ts.success).toBe(true);
  });

  test("rejects non-loopback local_url and unsupported provider", () => {
    const nonLoopback = tunnelConfigSchema.safeParse({
      provider: "cloudflare",
      local_url: "http://10.0.0.5:20128",
      tailscale_mode: "funnel",
      auto_start: false,
    });
    expect(nonLoopback.success).toBe(false);

    const badProv = tunnelConfigSchema.safeParse({
      provider: "ngrok",
      local_url: "http://127.0.0.1:20128",
      tailscale_mode: "funnel",
      auto_start: false,
    });
    expect(badProv.success).toBe(false);
  });
});
