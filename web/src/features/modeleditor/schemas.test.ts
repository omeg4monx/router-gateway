import { describe, expect, test } from "bun:test";
import { BUILTIN_PROMPT_TEMPLATES } from "./builtinPromptTemplates";
import { filterPromptTemplates } from "./PromptTemplateCombobox";
import {
  classifyPromptTestOutcome,
  estimatePromptTokens,
  getParameterGuideEntry,
  MODEL_EDITOR_COPY,
  PARAMETER_GUIDE_ENTRIES,
  QUICK_TEST_PRESETS,
} from "./parameterGuide";
import {
  buildOverridePatchPayload,
  createModelOverrideSchema,
  overrideToFormValues,
  parseExtraParamsJson,
  parseStripParamsList,
} from "./schemas";
import { buildTargetModelOptions, filterTargetModelOptions } from "./TargetModelCombobox";

describe("modeleditor schemas", () => {
  test("validates model override creation with numeric bounds", () => {
    const parsed = createModelOverrideSchema.safeParse({
      name: "Coding Strict",
      target_model: "claude-sonnet-4",
      system_prompt: "Answer directly without preamble.",
      system_prompt_position: "prefix",
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
      target_model: "gpt-4o",
      system_prompt_position: "prefix",
      temperature: "3.5",
    });
    expect(badTemp.success).toBe(false);

    expect(parseExtraParamsJson("{bad json}").ok).toBe(false);
    expect(parseExtraParamsJson("[]").ok).toBe(false);
    expect(parseExtraParamsJson('{"seed":1}').ok).toBe(true);
    expect(parseStripParamsList(" a , b , a , ")).toEqual(["a", "b"]);
  });

  test("provides complete bilingual parameter guide entries", () => {
    expect(PARAMETER_GUIDE_ENTRIES.length).toBe(7);
    const disclaimerEntry = getParameterGuideEntry("disclaimer");
    expect(disclaimerEntry.id).toBe("disclaimer");
    expect(disclaimerEntry.purpose.id).toContain("30–50%");
    expect(disclaimerEntry.purpose.id).toContain("50–80%");
    expect(disclaimerEntry.purpose.id).toContain("0%");
    expect(disclaimerEntry.whenDecreased.effect.id).toContain("thinking");
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

    expect(QUICK_TEST_PRESETS.map((p) => p.id)).toEqual([
      "account_farming",
      "sql_injection",
      "illegal_story",
      "lock_picking",
    ]);
    for (const preset of QUICK_TEST_PRESETS) {
      expect(preset.prompt.id.length).toBeGreaterThan(10);
      expect(preset.prompt.en.length).toBeGreaterThan(10);
    }
  });

  test("builds deduplicated target model options and filters by model or provider", () => {
    const options = buildTargetModelOptions([
      {
        providerId: "anthropic",
        provider: "Anthropic",
        id: "claude-sonnet-4",
        qualifiedId: "anthropic/claude-sonnet-4",
        name: "Claude Sonnet 4",
        type: "chat",
      },
      {
        providerId: "openai",
        provider: "OpenAI",
        id: "gpt-4o",
        qualifiedId: "openai/gpt-4o",
        name: "GPT-4o",
        type: "chat",
      },
      {
        providerId: "openai",
        provider: "OpenAI",
        id: "gpt-4o",
        qualifiedId: "openai/gpt-4o",
        name: "Duplicate",
        type: "chat",
      },
      {
        providerId: "openai",
        provider: "OpenAI",
        id: "hidden-model",
        qualifiedId: "openai/hidden-model",
        name: "Hidden",
        type: "chat",
        hidden: true,
      },
    ]);

    expect(options.map((o) => o.id)).toEqual(["claude-sonnet-4", "gpt-4o"]);
    expect(filterTargetModelOptions(options, "anthropic").map((o) => o.id)).toEqual([
      "claude-sonnet-4",
    ]);
    expect(filterTargetModelOptions(options, "4o").map((o) => o.id)).toEqual(["gpt-4o"]);
  });

  test("classifies prompt test outcomes into answered, refusal, empty, and error", () => {
    expect(classifyPromptTestOutcome({ error: "HTTP 502 upstream timeout" }, "id").verdict).toBe(
      "error",
    );

    expect(classifyPromptTestOutcome({ content: "   " }, "id").verdict).toBe("empty");

    const enRefusal = classifyPromptTestOutcome(
      { content: "I'm sorry, but I cannot assist with SQL injection payloads." },
      "en",
    );
    expect(enRefusal.verdict).toBe("refusal");
    expect(enRefusal.badgeLabel).toBe("Refusal");

    const idRefusal = classifyPromptTestOutcome(
      { content: "Maaf, saya tidak dapat membantu membuat skrip farming akun." },
      "id",
    );
    expect(idRefusal.verdict).toBe("refusal");

    const answered = classifyPromptTestOutcome(
      {
        content: "Single-pin picking applies light torque with a tension wrench...",
        completionTokens: 42,
      },
      "id",
    );
    expect(answered.verdict).toBe("answered");
    expect(answered.badgeLabel).toBe("Menjawab");
    expect(answered.summary).toContain("42 tok");

    expect(estimatePromptTokens("")).toBe(0);
    expect(estimatePromptTokens("abcd")).toBe(1);
    expect(estimatePromptTokens("Answer concisely.")).toBe(5);

    expect(BUILTIN_PROMPT_TEMPLATES.map((t) => t.name)).toEqual([
      "GLM 5.3",
      "DeepSeek V2",
      "DeepSeek 4.1",
      "Gemini 4 Neon",
      "Gemini 4 Hydrogen",
      "DeepSeek 4.1 & 1017",
    ]);
    for (const tpl of BUILTIN_PROMPT_TEMPLATES) {
      expect(tpl.content.length).toBeGreaterThan(100);
    }
    expect(filterPromptTemplates(BUILTIN_PROMPT_TEMPLATES, "hydrogen").map((t) => t.name)).toEqual([
      "Gemini 4 Hydrogen",
    ]);
    expect(filterPromptTemplates(BUILTIN_PROMPT_TEMPLATES, "neon").map((t) => t.name)).toEqual([
      "Gemini 4 Neon",
    ]);
  });

  test("maps ModelOverride to edit form values and builds PATCH payload with clear flags", () => {
    const existing = {
      id: "mov_123",
      name: "claude-dual",
      model_pattern: "claude-dual",
      target_model: "claude-sonnet-4",
      system_prompt_prefix: "",
      system_prompt_suffix: "Anchor instruction",
      temperature: 0.4,
      max_tokens: 4096,
      extra_params: { _prompt_mode: "dual", seed: 7 },
      strip_params: ["presence_penalty", "top_k"],
      priority: 100,
      enabled: true,
      created_at: "2026-01-01T00:00:00Z",
      updated_at: "2026-01-01T00:00:00Z",
    };

    const formValues = overrideToFormValues(existing);
    expect(formValues.name).toBe("claude-dual");
    expect(formValues.target_model).toBe("claude-sonnet-4");
    expect(formValues.system_prompt).toBe("Anchor instruction");
    expect(formValues.system_prompt_position).toBe("dual");
    expect(formValues.temperature).toBe("0.4");
    expect(formValues.max_tokens).toBe("4096");
    expect(formValues.strip_params_csv).toBe("presence_penalty, top_k");
    expect(JSON.parse(formValues.extra_params_json ?? "{}")).toEqual({ seed: 7 });

    // Edit: switch to prefix, clear temperature and max_tokens.
    const patch = buildOverridePatchPayload(
      {
        ...formValues,
        system_prompt_position: "prefix",
        temperature: "",
        max_tokens: "",
        extra_params_json: '{"seed": 99}',
      },
      existing,
    );
    expect(patch.ok).toBe(true);
    if (!patch.ok) return;
    expect(patch.value.system_prompt_prefix).toBe("Anchor instruction");
    expect(patch.value.system_prompt_suffix).toBe("");
    expect(patch.value.clear_temperature).toBe(true);
    expect(patch.value.clear_max_tokens).toBe(true);
    expect(patch.value.temperature).toBeUndefined();
    expect(patch.value.max_tokens).toBeUndefined();
    expect(patch.value.extra_params).toEqual({ seed: 99 });
  });
});
