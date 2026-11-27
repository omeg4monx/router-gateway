import { describe, expect, test } from "bun:test";
import {
  arenaBattleFormSchema,
  chatMessageFormSchema,
  createChatSessionFormSchema,
  updateChatSessionFormSchema,
} from "./schemas";

describe("arena and chat schemas", () => {
  test("accepts valid 2-4 model battle request", () => {
    const res = arenaBattleFormSchema.safeParse({
      prompt: "Explain WAL mode in SQLite",
      system_prompt: "Be concise.",
      models: ["openai/gpt-4o", "anthropic/claude-sonnet"],
      temperature: 0.7,
      max_tokens: 512,
    });
    expect(res.success).toBe(true);
  });

  test("rejects fewer than 2 or more than 4 models", () => {
    expect(
      arenaBattleFormSchema.safeParse({
        prompt: "hi",
        system_prompt: "",
        models: ["openai/gpt-4o"],
      }).success,
    ).toBe(false);
    expect(
      arenaBattleFormSchema.safeParse({
        prompt: "hi",
        system_prompt: "",
        models: ["m1", "m2", "m3", "m4", "m5"],
      }).success,
    ).toBe(false);
  });

  test("validates chat session creation and message submission", () => {
    expect(
      createChatSessionFormSchema.safeParse({
        title: "Architecture Q&A",
        model: "openai/gpt-4o",
        system_prompt: "",
      }).success,
    ).toBe(true);
    expect(
      createChatSessionFormSchema.safeParse({
        title: "",
        model: "",
        system_prompt: "",
      }).success,
    ).toBe(false);
    expect(chatMessageFormSchema.safeParse({ content: "Hello" }).success).toBe(true);
    expect(chatMessageFormSchema.safeParse({ content: "   " }).success).toBe(false);
  });

  test("validates optional temperature and max_tokens on chat session create and update", () => {
    const created = createChatSessionFormSchema.safeParse({
      title: "Latency Check",
      model: "openai/gpt-4o",
      system_prompt: "Be brief.",
      temperature: 0.4,
      max_tokens: 1024,
    });
    expect(created.success).toBe(true);

    const invalidTemp = createChatSessionFormSchema.safeParse({
      model: "openai/gpt-4o",
      temperature: 3.5,
    });
    expect(invalidTemp.success).toBe(false);

    const updated = updateChatSessionFormSchema.safeParse({
      title: "Renamed Session",
      model: "anthropic/claude-sonnet",
      system_prompt: "Answer in bullets.",
      temperature: 0.7,
      max_tokens: 2048,
    });
    expect(updated.success).toBe(true);
  });
});
