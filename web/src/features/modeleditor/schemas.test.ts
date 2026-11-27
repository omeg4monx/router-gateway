import { describe, expect, test } from "bun:test";
import {
  getParameterGuideEntry,
  MODEL_EDITOR_COPY,
  PARAMETER_GUIDE_ENTRIES,
} from "./parameterGuide";
import { createModelOverrideSchema, parseExtraParamsJson, parseStripParamsList } from "./schemas";

describe("modeleditor schemas", () => {
  test("validates model override creation with numeric bounds", () => {
    const parsed = createModelOverrideSchema.safeParse({
      name: "Coding Strict",
      model_pattern: "coding-*",
      target_model: "claude-sonnet-4",
      temperature: "0.2",
      max_tokens: "4096",
      extra_params_json: '{"seed": 42}',
      strip_params_csv: "presence_penalty, frequency_penalty",
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

  test("provides complete bilingual parameter guide entries", () => {
    expect(PARAMETER_GUIDE_ENTRIES.length).toBeGreaterThanOrEqual(7);
    const tempEntry = getParameterGuideEntry("temperature");
    expect(tempEntry.id).toBe("temperature");
    expect(tempEntry.defaultValue.id.length).toBeGreaterThan(0);
    expect(tempEntry.whenIncreased.effect.id.length).toBeGreaterThan(0);
    expect(tempEntry.whenDecreased.effect.id.length).toBeGreaterThan(0);

    const extraEntry = getParameterGuideEntry("extra_params");
    expect(extraEntry.id).toBe("extra_params");
    expect(extraEntry.example).toContain("seed");

    expect(MODEL_EDITOR_COPY.id.guideTitle).toBeTruthy();
    expect(MODEL_EDITOR_COPY.en.guideTitle).toBeTruthy();
  });
});
