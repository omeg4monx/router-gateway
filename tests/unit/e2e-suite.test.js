import { describe, it, expect } from "vitest";
import { exportDb } from "../../src/lib/db/index.js";
import { applyFingerprintTools } from "../../open-sse/utils/opencodeFingerprint.js";

describe("E2E Verification Suite", () => {
  it("Auto-Backup: granular scope excludes usage history when requested", async () => {
    const compactExport = await exportDb({ usage: false });
    expect(compactExport.settings).toBeDefined();
    expect(compactExport.usageHistory).toBeUndefined();
    expect(compactExport.usageDaily).toBeUndefined();

    const fullExport = await exportDb(null);
    expect(fullExport.settings).toBeDefined();
    expect(fullExport.usageHistory).toBeDefined();
  });

  it("opencode-cc Hardening: tool surface capping trims tools over 38K characters", () => {
    const fakeTools = Array.from({ length: 40 }, (_, i) => ({
      type: "function",
      function: {
        name: `custom_tool_${i}`,
        description: "A very long detailed description ".repeat(40),
        parameters: { type: "object", properties: { arg: { type: "string" } } }
      }
    }));
    const body = {
      tools: fakeTools,
      messages: [{ role: "user", content: "hello" }]
    };
    expect(JSON.stringify(body.tools).length).toBeGreaterThan(38000);

    applyFingerprintTools(body, false);
    expect(JSON.stringify(body.tools).length).toBeLessThanOrEqual(38000);
  });
});
