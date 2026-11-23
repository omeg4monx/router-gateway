import { describe, expect, test } from "bun:test";
import { claimVoucherSchema, createVoucherSchema } from "./schemas";

describe("vouchers schemas", () => {
  test("validates voucher creation and claim inputs", () => {
    expect(
      createVoucherSchema.safeParse({
        code: "PROMO2026",
        name: "Promo 2026",
        description: "Hackathon voucher",
        quota_usd: 5,
        rate_limit_rpm: 30,
        allowed_models_text: "gpt-4o, claude-3-7-sonnet",
        team_id: "",
        max_redemptions: 10,
      }).success,
    ).toBe(true);

    expect(
      claimVoucherSchema.safeParse({
        code: "PROMO2026",
        claimer_id: "alice",
        key_name: "alice-key",
      }).success,
    ).toBe(true);
  });

  test("rejects invalid max_redemptions or blank claim code", () => {
    expect(
      createVoucherSchema.safeParse({
        code: "",
        name: "Invalid",
        description: "",
        quota_usd: 0,
        rate_limit_rpm: 0,
        allowed_models_text: "",
        team_id: "",
        max_redemptions: 0,
      }).success,
    ).toBe(false);

    expect(
      claimVoucherSchema.safeParse({
        code: "  ",
        claimer_id: "",
        key_name: "",
      }).success,
    ).toBe(false);
  });
});
