import { z } from "zod";

export const MIN_PASSWORD_RUNES = 8;
export const MAX_PASSWORD_BYTES = 72;

export interface InitPasswordValidation {
  runeCount: number;
  byteLength: number;
  meetsMinLength: boolean;
  withinMaxBytes: boolean;
  passwordsMatch: boolean;
  error: string | null;
}

export function countPasswordRunes(value: string): number {
  return Array.from(value).length;
}

export function utf8ByteLength(value: string): number {
  return new TextEncoder().encode(value).length;
}

export const initPasswordSchema = z
  .object({
    password: z
      .string()
      .refine(
        (value) => countPasswordRunes(value) >= MIN_PASSWORD_RUNES,
        `Password must be at least ${MIN_PASSWORD_RUNES} characters.`,
      )
      .refine(
        (value) => utf8ByteLength(value) <= MAX_PASSWORD_BYTES,
        `Password must not exceed ${MAX_PASSWORD_BYTES} UTF-8 bytes.`,
      ),
    confirmPassword: z.string(),
  })
  .superRefine((data, ctx) => {
    if (data.password !== data.confirmPassword) {
      ctx.addIssue({
        code: "custom",
        path: ["confirmPassword"],
        message: "Password confirmation does not match.",
      });
    }
  });

export type InitPasswordValues = z.infer<typeof initPasswordSchema>;

export function validateInitPassword(
  password: string,
  confirmPassword: string,
): InitPasswordValidation {
  const runeCount = countPasswordRunes(password);
  const byteLength = utf8ByteLength(password);
  const meetsMinLength = runeCount >= MIN_PASSWORD_RUNES;
  const withinMaxBytes = byteLength <= MAX_PASSWORD_BYTES;
  const passwordsMatch = confirmPassword.length > 0 && password === confirmPassword;

  let error: string | null = null;
  if (!meetsMinLength) {
    error = `Password must be at least ${MIN_PASSWORD_RUNES} characters.`;
  } else if (!withinMaxBytes) {
    error = `Password must not exceed ${MAX_PASSWORD_BYTES} UTF-8 bytes.`;
  } else if (password !== confirmPassword) {
    error = "Password confirmation does not match.";
  }

  return {
    runeCount,
    byteLength,
    meetsMinLength,
    withinMaxBytes,
    passwordsMatch,
    error,
  };
}
