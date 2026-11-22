import { describe, expect, test } from "bun:test";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  countPasswordRunes,
  initPasswordSchema,
  utf8ByteLength,
  validateInitPassword,
} from "./password";

describe("init password validation (Go utf8.RuneCountInString parity + Zod schema)", () => {
  test("counts astral Unicode surrogate pairs as single runes, not UTF-16 code units", () => {
    const fourRunes = "\u{1F511}\u{1F511}\u{1F511}\u{1F511}";
    expect(fourRunes.length).toBe(8);
    expect(countPasswordRunes(fourRunes)).toBe(4);
    expect(utf8ByteLength(fourRunes)).toBe(16);

    const result = validateInitPassword(fourRunes, fourRunes);
    expect(result.meetsMinLength).toBe(false);
    expect(result.error).toBe("Password must be at least 8 characters.");

    const parsed = initPasswordSchema.safeParse({
      password: fourRunes,
      confirmPassword: fourRunes,
    });
    expect(parsed.success).toBe(false);
  });

  test("accepts 8 Unicode runes within 72 UTF-8 bytes in both helper and zodResolver", async () => {
    const eightRunes = "\u00e9\u00e9\u00e9\u00e9\u00e9\u00e9\u00e9\u00e9";
    expect(countPasswordRunes(eightRunes)).toBe(8);
    expect(utf8ByteLength(eightRunes)).toBe(16);

    const result = validateInitPassword(eightRunes, eightRunes);
    expect(result.meetsMinLength).toBe(true);
    expect(result.withinMaxBytes).toBe(true);
    expect(result.passwordsMatch).toBe(true);
    expect(result.error).toBeNull();

    const resolver = zodResolver(initPasswordSchema);
    const resolved = await resolver(
      { password: eightRunes, confirmPassword: eightRunes },
      undefined,
      { fields: {}, shouldUseNativeValidation: false },
    );
    expect(resolved.errors).toEqual({});
    expect(resolved.values).toEqual({
      password: eightRunes,
      confirmPassword: eightRunes,
    });
  });

  test("rejects multi-byte passwords exceeding 72 UTF-8 bytes via zodResolver", async () => {
    const tooManyBytes = "\u00e9".repeat(37); // 37 runes, 74 UTF-8 bytes
    expect(countPasswordRunes(tooManyBytes)).toBe(37);
    expect(utf8ByteLength(tooManyBytes)).toBe(74);

    const resolver = zodResolver(initPasswordSchema);
    const resolved = await resolver(
      { password: tooManyBytes, confirmPassword: tooManyBytes },
      undefined,
      { fields: {}, shouldUseNativeValidation: false },
    );
    expect(resolved.errors.password?.message).toBe(
      "Password must not exceed 72 UTF-8 bytes.",
    );
  });

  test("rejects mismatched confirmation via zodResolver on confirmPassword field", async () => {
    const resolver = zodResolver(initPasswordSchema);
    const resolved = await resolver(
      { password: "valid-password", confirmPassword: "different-password" },
      undefined,
      { fields: {}, shouldUseNativeValidation: false },
    );
    expect(resolved.errors.confirmPassword?.message).toBe(
      "Password confirmation does not match.",
    );
  });
});
