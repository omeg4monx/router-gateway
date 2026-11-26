import { describe, expect, test } from "bun:test";
import {
  classifyTestError,
  comboTestResultKey,
  flattenCatalogResponse,
  modelTestResultKey,
  normalizeTestStatus,
} from "./catalogUtils";

describe("models/catalogUtils", () => {
  test("flattens provider-grouped catalog into normalized rows", () => {
    const rows = flattenCatalogResponse({
      openai: {
        provider: "openai",
        models: [
          {
            id: "gpt-4o",
            name: "gpt-4o",
            type: "chat",
            subtype: "general",
            context_length: 128000,
            max_output_tokens: 16384,
            capabilities: { vision: true, streaming: true },
          },
          {
            id: "o3-mini",
            name: "o3-mini",
            type: "chat",
            subtype: "reasoning",
            hidden: true,
          },
        ],
      },
      ollama: {
        provider: "ollama",
        models: [{ id: "llama3.3", type: "chat", free: true, custom: true }],
      },
    });

    expect(rows).toHaveLength(3);
    expect(rows[0]?.qualifiedId).toBe("openai/gpt-4o");
    expect(rows[1]?.hidden).toBe(true);
    expect(rows[2]?.free).toBe(true);
  });

  test("builds stable model and combo test result keys", () => {
    expect(modelTestResultKey("openai", "gpt-4o")).toBe("model:openai:gpt-4o");
    expect(comboTestResultKey("coding-primary")).toBe("combo:coding-primary");
  });

  test("classifies test errors and normalizes slow threshold (>=5000ms)", () => {
    expect(classifyTestError("rate limit exceeded", 429)).toBe("rate-limited");
    expect(classifyTestError("insufficient balance on account", 402)).toBe("quota");
    expect(classifyTestError("request timed out", 504)).toBe("timeout");
    expect(classifyTestError("internal server error", 500)).toBe("other");

    expect(normalizeTestStatus("ok", 1200)).toBe("ok");
    expect(normalizeTestStatus("ok", 5200)).toBe("slow");
    expect(normalizeTestStatus("error", 100)).toBe("error");
  });
});
