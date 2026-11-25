import { describe, expect, test } from "bun:test";
import { formatAuthError, getLoginCopy } from "./i18n";

const EXPECTED_MINIMAL_KEYS = [
  "signInHeading",
  "pendingLabel",
  "unreachableTitle",
  "retryLabel",
  "retryingLabel",
  "uninitializedTitle",
  "uninitializedBody",
  "uninitializedCta",
  "redirectTitle",
  "redirectCta",
  "passwordLabel",
  "passwordPlaceholder",
  "showPasswordLabel",
  "hidePasswordLabel",
  "capsLockWarning",
  "submitLabel",
  "submittingLabel",
  "errorTitle",
  "technicalDetailLabel",
].sort();

describe("auth/i18n", () => {
  test("provides only essential localized login copy for en and id", () => {
    const en = getLoginCopy("en");
    const id = getLoginCopy("id");

    expect(Object.keys(en).sort()).toEqual(EXPECTED_MINIMAL_KEYS);
    expect(Object.keys(id).sort()).toEqual(EXPECTED_MINIMAL_KEYS);
    expect(en.signInHeading).not.toBe(id.signInHeading);
    expect(en.passwordLabel).not.toBe(id.passwordLabel);
    expect(en.submitLabel).not.toBe(id.submitLabel);
  });

  test("maps invalid admin password error to friendly localized copy with technical detail", () => {
    const enErr = formatAuthError("invalid admin password", "en");
    const idErr = formatAuthError("invalid admin password", "id");

    expect(enErr).not.toBeNull();
    expect(enErr?.message.toLowerCase()).toContain("password");
    expect(enErr?.detail).toBe("invalid admin password");

    expect(idErr).not.toBeNull();
    expect(idErr?.message.toLowerCase()).toContain("kata sandi");
    expect(idErr?.detail).toBe("invalid admin password");
  });

  test("maps required password validation error to localized copy without duplicating detail", () => {
    const idErr = formatAuthError("Administrator password is required.", "id");
    expect(idErr?.message.toLowerCase()).toContain("wajib");
    expect(idErr?.detail).toBeUndefined();
  });

  test("preserves unrecognized backend error messages cleanly", () => {
    const err = formatAuthError("upstream timeout after 5s", "en");
    expect(err?.message).toBe("upstream timeout after 5s");
    expect(err?.detail).toBeUndefined();
    expect(formatAuthError(null, "en")).toBeNull();
    expect(formatAuthError("   ", "id")).toBeNull();
  });
});
