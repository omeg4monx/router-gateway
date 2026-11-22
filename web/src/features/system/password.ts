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

export function validateInitPassword(
  password: string,
  confirmPassword: string,
): InitPasswordValidation {
  const runeCount = countPasswordRunes(password);
  const byteLength = utf8ByteLength(password);
  const meetsMinLength = runeCount >= MIN_PASSWORD_RUNES;
  const withinMaxBytes = byteLength <= MAX_PASSWORD_BYTES;
  const passwordsMatch =
    confirmPassword.length > 0 && password === confirmPassword;

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
