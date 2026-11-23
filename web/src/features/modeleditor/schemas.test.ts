import { describe, expect, test } from "bun:test";
import {
  createBypassRuleSchema,
  createModelOverrideSchema,
  parseExtraParamsJson,
  parseStripParamsList,
} from "./schemas";

describe("modeleditor schemas", () => {
  test("validates model override creation with numeric bounds", () => {
    const parsed = createModelOverrideSchema.safeParse({
      name: "Coding Strict",
      model_pattern: "coding-*",
      target_model: "claude-sonnet-4",
      temperature: "0.2",
      max_tokens: "4096",
      top_p: "0.95",
      extra_params_json: '{"seed": 42}',
      strip_params_csv: "presence_penalty, frequency_penalty",
      priority: "10",
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects out-of-range temperature or invalid extra_params JSON", () => {
    const badTemp = createModelOverrideSchema.safeParse({
      name: "Bad Temp",
      model_pattern: "*",
      temperature: "3.5",
    });
    expect(badTemp.success).toBe(false);

    expect(parseExtraParamsJson("{bad json}").ok).toBe(false);
    expect(parseExtraParamsJson("[]").ok).toBe(false);
    expect(parseExtraParamsJson('{"seed":1}').ok).toBe(true);
    expect(parseStripParamsList(" a , b , a , ")).toEqual(["a", "b"]);
  });

  test("requires target_model when bypass action is rewrite_model", () => {
    const missingTarget = createBypassRuleSchema.safeParse({
      name: "Redirect",
      model_pattern: "legacy-*",
      action: "rewrite_model",
      target_model: "",
    });
    expect(missingTarget.success).toBe(false);

    const validBypass = createBypassRuleSchema.safeParse({
      name: "Bypass TokenSaver",
      model_pattern: "gpt-*",
      action: "bypass_tokensaver",
      header_name: "X-Raw",
      header_value: "true",
      priority: "20",
    });
    expect(validBypass.success).toBe(true);
  });
});
