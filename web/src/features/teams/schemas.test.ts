import { describe, expect, test } from "bun:test";
import { createTeamSchema } from "./schemas";

describe("teams schemas", () => {
  test("validates valid team creation with numeric budget and RPM", () => {
    const parsed = createTeamSchema.safeParse({
      name: "Platform Eng",
      description: "Shared quota for platform team",
      budget_usd: 50,
      rate_limit_rpm: 120,
      enabled: true,
    });
    expect(parsed.success).toBe(true);
  });

  test("rejects blank name, negative budget, or negative RPM", () => {
    expect(
      createTeamSchema.safeParse({
        name: "   ",
        description: "",
        budget_usd: 10,
        rate_limit_rpm: 0,
        enabled: true,
      }).success,
    ).toBe(false);

    expect(
      createTeamSchema.safeParse({
        name: "Bad Budget",
        description: "",
        budget_usd: -1,
        rate_limit_rpm: 0,
        enabled: true,
      }).success,
    ).toBe(false);
  });
});
