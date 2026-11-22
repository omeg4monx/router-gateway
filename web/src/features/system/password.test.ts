import { describe, expect, test } from "bun:test";
import {
  countPasswordRunes,
  MAX_PASSWORD_BYTES,
  MIN_PASSWORD_RUNES,
  utf8ByteLength,
  validateInitPassword,
} from "./password";

describe("init password validation (Go utf8.RuneCountInString parity)", () => {
  test("counts astral Unicode surrogate pairs as single runes, not UTF-16 code units", () => {
    // 4 mathematical bold characters (U+1D400): 8 UTF-16 code units, 4 Unicode runes, 16 UTF-8 bytes.
    const fourRunes = "\u{1D400}\u{1D401}\u{1D402}\u{1D403}";
    expect(fourRunes.length).toBe(8);
    expect(countPasswordRunes(fourRunes)).toBe(4);
    expect(utf8ByteLength(fourRunes)).toBe(16);

    const result = validateInitPassword(fourRunes, fourRunes);
    expect(result.meetsMinLength).toBe(false);
    expect(result.error).toBe(
      `Password must be at least ${MIN_PASSWORD_RUNES} characters.`,
    );
  });

  test("accepts 8 Unicode runes within 72 UTF-8 bytes", () => {
    const eightRunes = "\u{1D400}".repeat(8); // 8 runes, 32 UTF-8 bytes
    const result = validateInitPassword(eightRunes, eightRunes);
    expect(result.runeCount).toBe(8);
    expect(result.byteLength).toBe(32);
    expect(result.meetsMinLength).toBe(true);
    expect(result.withinMaxBytes).toBe(true);
    expect(result.passwordsMatch).toBe(true);
    expect(result.error).toBeNull();
  });

  test("rejects multi-byte passwords exceeding 72 UTF-8 bytes even when rune count <= 72", () => {
    // 25 3-byte runes = 75 UTF-8 bytes (> 72 bcrypt limit)
    const seventyFiveBytes = "\u20AC".repeat(25);
    const result = validateInitPassword(seventyFiveBytes, seventyFiveBytes);
    expect(result.runeCount).toBe(25);
    expect(result.byteLength).toBe(75);
    expect(result.withinMaxBytes).toBe(false);
    expect(result.error).toBe(
      `Password must not exceed ${MAX_PASSWORD_BYTES} UTF-8 bytes.`,
    );
  });

  test("rejects mismatched confirmation", () => {
    const result = validateInitPassword("correct-horse", "different-horse");
    expect(result.meetsMinLength).toBe(true);
    expect(result.withinMaxBytes).toBe(true);
    expect(result.passwordsMatch).toBe(false);
    expect(result.error).toBe("Password confirmation does not match.");
  });
});
