import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  COMBO_STRATEGIES,
  createComboSchema,
  disableModelSchema,
  extractTargetsFieldErrors,
  MAX_MODEL_IDENTIFIER_RUNES,
  toCreateComboRequest,
  toDisableModelRequest,
  toUpsertAliasRequest,
  upsertAliasSchema,
} from "./schemas";

describe("createComboSchema", () => {
  test("supports all 22 OmniRoute combo strategies", () => {
    expect(COMBO_STRATEGIES).toHaveLength(22);
    for (const strategy of COMBO_STRATEGIES) {
      const parsed = createComboSchema.safeParse({
        name: `combo-${strategy}`,
        description: "",
        strategy,
        enabled: true,
        targets: ["openai/gpt-4o"],
      });
      expect(parsed.success).toBe(true);
    }
  });

  test("parses valid priority combo and trims name, description, and targets", async () => {
    const input = {
      name: "  coding-fallback  ",
      description: "  Primary GPT-4o with Claude Sonnet fallback  ",
      strategy: "priority" as const,
      enabled: true,
      targets: ["  openai/gpt-4o ", "anthropic/claude-sonnet-4-20250514"],
    };

    const parsed = createComboSchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(toCreateComboRequest(parsed.data)).toEqual({
        name: "coding-fallback",
        description: "Primary GPT-4o with Claude Sonnet fallback",
        strategy: "priority",
        enabled: true,
        targets: ["openai/gpt-4o", "anthropic/claude-sonnet-4-20250514"],
      });
    }

    const resolver = zodResolver(createComboSchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
  });

  test("rejects empty combo name or wildcard '*'", () => {
    const emptyName = createComboSchema.safeParse({
      name: "   ",
      description: "",
      strategy: "priority",
      enabled: true,
      targets: ["openai/gpt-4o"],
    });
    expect(emptyName.success).toBe(false);

    const wildcardName = createComboSchema.safeParse({
      name: "*",
      description: "",
      strategy: "priority",
      enabled: true,
      targets: ["openai/gpt-4o"],
    });
    expect(wildcardName.success).toBe(false);
  });

  test("requires at least 1 unique trimmed target and rejects duplicates or self-reference", () => {
    const noTargets = createComboSchema.safeParse({
      name: "empty-chain",
      description: "",
      strategy: "priority",
      enabled: true,
      targets: [],
    });
    expect(noTargets.success).toBe(false);

    const blankTarget = createComboSchema.safeParse({
      name: "blank-target",
      description: "",
      strategy: "priority",
      enabled: true,
      targets: ["   "],
    });
    expect(blankTarget.success).toBe(false);

    const duplicateTargets = createComboSchema.safeParse({
      name: "dup-chain",
      description: "",
      strategy: "priority",
      enabled: true,
      targets: ["openai/gpt-4o", " openai/gpt-4o "],
    });
    expect(duplicateTargets.success).toBe(false);

    const selfReferential = createComboSchema.safeParse({
      name: "self-combo",
      description: "",
      strategy: "priority",
      enabled: true,
      targets: ["openai/gpt-4o", " self-combo "],
    });
    expect(selfReferential.success).toBe(false);
  });

  test("extractTargetsFieldErrors surfaces per-element zodResolver errors and array-level errors", async () => {
    const resolver = zodResolver(createComboSchema);

    const selfRefResolved = await resolver(
      {
        name: "self-combo",
        description: "",
        strategy: "priority",
        enabled: true,
        targets: ["openai/gpt-4o", "self-combo"],
      },
      undefined,
      { fields: {}, shouldUseNativeValidation: false },
    );
    const selfRefExtracted = extractTargetsFieldErrors(selfRefResolved.errors.targets);
    expect(selfRefExtracted.itemMessages[1]).toContain("must not reference the combo name");
    expect(selfRefExtracted.summaryMessage).toContain("must not reference the combo name");

    const tooLongTarget = "x".repeat(MAX_MODEL_IDENTIFIER_RUNES + 1);
    const tooLongResolved = await resolver(
      {
        name: "long-target-combo",
        description: "",
        strategy: "priority",
        enabled: true,
        targets: [tooLongTarget],
      },
      undefined,
      { fields: {}, shouldUseNativeValidation: false },
    );
    const tooLongExtracted = extractTargetsFieldErrors(tooLongResolved.errors.targets);
    expect(tooLongExtracted.itemMessages[0]).toContain(`${MAX_MODEL_IDENTIFIER_RUNES} characters`);
    expect(tooLongExtracted.summaryMessage).toContain(`${MAX_MODEL_IDENTIFIER_RUNES} characters`);

    const emptyResolved = await resolver(
      {
        name: "empty-combo",
        description: "",
        strategy: "priority",
        enabled: true,
        targets: [],
      },
      undefined,
      { fields: {}, shouldUseNativeValidation: false },
    );
    const emptyExtracted = extractTargetsFieldErrors(emptyResolved.errors.targets);
    expect(emptyExtracted.itemMessages).toEqual({});
    expect(emptyExtracted.summaryMessage).toBe("Add at least one fallback target model.");
  });
});

describe("upsertAliasSchema", () => {
  test("parses valid alias mapping and trims whitespace", async () => {
    const input = {
      alias: "  fast-chat  ",
      target_model: "  openai/gpt-4o-mini  ",
    };

    const parsed = upsertAliasSchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(parsed.data.alias).toBe("fast-chat");
      expect(toUpsertAliasRequest(parsed.data)).toEqual({
        target_model: "openai/gpt-4o-mini",
      });
    }

    const resolver = zodResolver(upsertAliasSchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
  });

  test("rejects self-referential alias === target_model after trimming", () => {
    const parsed = upsertAliasSchema.safeParse({
      alias: " gpt-4o ",
      target_model: "gpt-4o",
    });
    expect(parsed.success).toBe(false);
    if (!parsed.success) {
      expect(parsed.error.issues.some((issue) => issue.path.includes("target_model"))).toBe(true);
    }
  });

  test("rejects empty or wildcard alias/target_model", () => {
    expect(
      upsertAliasSchema.safeParse({
        alias: "",
        target_model: "openai/gpt-4o",
      }).success,
    ).toBe(false);

    expect(
      upsertAliasSchema.safeParse({
        alias: "smart",
        target_model: "*",
      }).success,
    ).toBe(false);
  });
});

describe("disableModelSchema", () => {
  test("parses model identifier and optional reason", async () => {
    const input = {
      model: "  openai/gpt-4-turbo-preview  ",
      reason: "  Deprecated upstream; route to gpt-4o instead  ",
    };

    const parsed = disableModelSchema.safeParse(input);
    expect(parsed.success).toBe(true);
    if (parsed.success) {
      expect(toDisableModelRequest(parsed.data)).toEqual({
        model: "openai/gpt-4-turbo-preview",
        reason: "Deprecated upstream; route to gpt-4o instead",
      });
    }

    const resolver = zodResolver(disableModelSchema);
    const resolved = await resolver(input, undefined, {
      fields: {},
      shouldUseNativeValidation: false,
    });
    expect(resolved.errors).toEqual({});
  });

  test("allows empty reason and rejects empty or wildcard model", () => {
    const emptyReason = disableModelSchema.safeParse({
      model: "claude-3-opus",
      reason: "   ",
    });
    expect(emptyReason.success).toBe(true);
    if (emptyReason.success) {
      expect(toDisableModelRequest(emptyReason.data)).toEqual({
        model: "claude-3-opus",
        reason: "",
      });
    }

    expect(
      disableModelSchema.safeParse({
        model: "   ",
        reason: "invalid",
      }).success,
    ).toBe(false);

    expect(
      disableModelSchema.safeParse({
        model: "*",
        reason: "wildcard",
      }).success,
    ).toBe(false);
  });
});
