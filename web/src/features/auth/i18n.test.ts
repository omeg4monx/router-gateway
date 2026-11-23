import { describe, expect, test } from "bun:test";
import { describeAuthPosture, formatAuthError, formatUptime, getLoginCopy } from "./i18n";

describe("auth/i18n", () => {
  test("provides complete localized copy for en and id", () => {
    const en = getLoginCopy("en");
    const id = getLoginCopy("id");

    expect(en.signInHeading).not.toBe(id.signInHeading);
    expect(en.passwordLabel).not.toBe(id.passwordLabel);
    expect(en.submitLabel).not.toBe(id.submitLabel);
    expect(en.helpTitle).not.toBe(id.helpTitle);
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

  test("describes all auth postures with semantic variants in en and id", () => {
    expect(describeAuthPosture({ isPending: true, isError: false }, "en").badgeVariant).toBe(
      "secondary",
    );
    expect(describeAuthPosture({ isPending: false, isError: true }, "en").badgeVariant).toBe(
      "destructive",
    );
    expect(
      describeAuthPosture(
        {
          isPending: false,
          isError: false,
          session: {
            authenticated: false,
            subject: "",
            require_login: true,
            remote_enforced: true,
            initialized: true,
          },
        },
        "id",
      ),
    ).toMatchObject({
      badgeVariant: "default",
    });
    expect(
      describeAuthPosture(
        {
          isPending: false,
          isError: false,
          session: {
            authenticated: false,
            subject: "",
            require_login: false,
            remote_enforced: false,
            initialized: true,
          },
        },
        "en",
      ).badgeVariant,
    ).toBe("success");
    expect(
      describeAuthPosture(
        {
          isPending: false,
          isError: false,
          session: {
            authenticated: false,
            subject: "",
            require_login: true,
            remote_enforced: false,
            initialized: false,
          },
        },
        "en",
      ).badgeVariant,
    ).toBe("warning");
  });

  test("formats uptime deterministically across edge cases in en and id", () => {
    expect(formatUptime(undefined, "en")).toBe("—");
    expect(formatUptime(-5, "en")).toBe("< 1m");
    expect(formatUptime(45, "id")).toBe("< 1 mnt");
    expect(formatUptime(3660, "en")).toBe("1h 1m");
    expect(formatUptime(3660, "id")).toBe("1j 1m");
    expect(formatUptime(90061, "en")).toBe("1d 1h");
    expect(formatUptime(90061, "id")).toBe("1h 1j");
  });
});
