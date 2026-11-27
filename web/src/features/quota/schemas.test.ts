import { describe, expect, test } from "bun:test";
import {
  computeQuotaOverviewStats,
  filterProviderQuotas,
  formatResetCountdown,
  getBucketProgressTone,
  getQuotaStatusVariant,
  groupProviderQuotas,
  type ProviderConnectionQuota,
  patchConnectionQuotaInResponse,
} from "./schemas";

const SAMPLE_QUOTAS: ProviderConnectionQuota[] = [
  {
    connection_id: "conn-1",
    provider_id: "tokenharbor",
    provider_name: "Token Harbor",
    fetcher_id: "tokenharbor",
    connection_name: "Harbor Office",
    auth_type: "api_key",
    enabled: true,
    status: "ok",
    plan_name: "Office Pass",
    on_plan: true,
    activated: true,
    buckets: [
      {
        id: "plan",
        label: "Office Pass",
        used_pct: 37,
        remaining_pct: 63,
        exhausted: false,
        period_started_at: "2026-11-20T09:00:00Z",
        period_ends_at: "2026-11-27T09:00:00Z",
        reset_at: "2026-11-27T09:00:00Z",
      },
      {
        id: "free_tier",
        label: "Free Tier",
        used_pct: 12,
        remaining_pct: 88,
        exhausted: false,
        period_started_at: "2026-11-20T09:00:00Z",
        period_ends_at: "2026-11-27T09:00:00Z",
        reset_at: "2026-11-27T09:00:00Z",
      },
    ],
    models: ["deepseek-v4.1-flash:free", "mimo-v2.6-flash:free"],
    fetched_at: "2026-11-27T01:00:00Z",
    cached: false,
    error: "",
  },
  {
    connection_id: "conn-2",
    provider_id: "tokenharbor",
    provider_name: "Token Harbor",
    fetcher_id: "tokenharbor",
    connection_name: "Harbor Low",
    auth_type: "api_key",
    enabled: true,
    status: "low",
    plan_name: null,
    on_plan: false,
    activated: true,
    buckets: [
      {
        id: "free_tier",
        label: "Free Tier",
        used_pct: 85,
        remaining_pct: 15,
        exhausted: false,
        period_started_at: null,
        period_ends_at: null,
        reset_at: "2026-11-28T01:00:00Z",
      },
    ],
    models: ["deepseek-v4.1-flash:free"],
    fetched_at: "2026-11-27T01:00:00Z",
    cached: true,
    error: "",
  },
  {
    connection_id: "conn-3",
    provider_id: "my-harbor",
    provider_name: "Custom Harbor",
    fetcher_id: "tokenharbor",
    connection_name: "Harbor Exhausted",
    auth_type: "api_key",
    enabled: true,
    status: "exhausted",
    plan_name: null,
    on_plan: false,
    activated: true,
    buckets: [
      {
        id: "free_tier",
        label: "Free Tier",
        used_pct: 100,
        remaining_pct: 0,
        exhausted: true,
        period_started_at: null,
        period_ends_at: null,
        reset_at: "2026-11-29T00:00:00Z",
      },
    ],
    models: [],
    fetched_at: "2026-11-27T01:00:00Z",
    cached: false,
    error: "",
  },
  {
    connection_id: "conn-4",
    provider_id: "tokenharbor",
    provider_name: "Token Harbor",
    fetcher_id: "tokenharbor",
    connection_name: "Harbor Inactive",
    auth_type: "api_key",
    enabled: true,
    status: "not_activated",
    plan_name: null,
    on_plan: false,
    activated: false,
    buckets: [
      {
        id: "free_tier",
        label: "Free Tier",
        used_pct: 0,
        remaining_pct: 100,
        exhausted: false,
        period_started_at: null,
        period_ends_at: null,
        reset_at: null,
      },
    ],
    models: [],
    fetched_at: "2026-11-27T01:00:00Z",
    cached: false,
    error: "",
  },
  {
    connection_id: "conn-5",
    provider_id: "tokenharbor",
    provider_name: "Token Harbor",
    fetcher_id: "tokenharbor",
    connection_name: "Harbor Broken",
    auth_type: "api_key",
    enabled: false,
    status: "error",
    plan_name: null,
    on_plan: false,
    activated: false,
    buckets: [],
    models: [],
    fetched_at: "2026-11-27T01:00:00Z",
    cached: false,
    error: "upstream rejected credentials (HTTP 401)",
  },
];

describe("features/quota/schemas", () => {
  test("computes overview statistics accurately across all quota statuses", () => {
    const stats = computeQuotaOverviewStats(SAMPLE_QUOTAS);
    expect(stats).toEqual({
      total: 5,
      healthy: 1,
      low: 1,
      exhausted: 1,
      notActivated: 1,
      error: 1,
    });
  });

  test("filters quotas by provider, status, and search query across name, plan, and models", () => {
    expect(
      filterProviderQuotas(SAMPLE_QUOTAS, { query: "", status: "all", providerId: "all" }),
    ).toHaveLength(5);

    expect(
      filterProviderQuotas(SAMPLE_QUOTAS, { query: "", status: "low", providerId: "all" }).map(
        (q) => q.connection_id,
      ),
    ).toEqual(["conn-2"]);

    expect(
      filterProviderQuotas(SAMPLE_QUOTAS, {
        query: "",
        status: "all",
        providerId: "my-harbor",
      }).map((q) => q.connection_id),
    ).toEqual(["conn-3"]);

    expect(
      filterProviderQuotas(SAMPLE_QUOTAS, {
        query: "office pass",
        status: "all",
        providerId: "all",
      }).map((q) => q.connection_id),
    ).toEqual(["conn-1"]);

    expect(
      filterProviderQuotas(SAMPLE_QUOTAS, {
        query: "mimo-v2.6",
        status: "all",
        providerId: "all",
      }).map((q) => q.connection_id),
    ).toEqual(["conn-1"]);
  });

  test("formats reset countdown cleanly in both English and Indonesian", () => {
    const now = new Date("2026-11-27T00:00:00Z").getTime();

    expect(formatResetCountdown("2026-11-29T05:30:00Z", "en", now)).toBe("in 2d 5h");
    expect(formatResetCountdown("2026-11-29T05:30:00Z", "id", now)).toBe("2h 5j lagi");
    expect(formatResetCountdown("2026-11-27T00:45:00Z", "en", now)).toBe("in 45m");
    expect(formatResetCountdown("2026-11-26T23:00:00Z", "en", now)).toBe("Resetting soon");
    expect(formatResetCountdown("2026-11-26T23:00:00Z", "id", now)).toBe("Segera reset");
    expect(formatResetCountdown(null, "en", now)).toBeNull();
    expect(formatResetCountdown("not-a-date", "en", now)).toBeNull();
  });

  test("maps status and bucket remaining percentages to semantic UI tones", () => {
    expect(getQuotaStatusVariant("ok")).toBe("success");
    expect(getQuotaStatusVariant("low")).toBe("warning");
    expect(getQuotaStatusVariant("exhausted")).toBe("destructive");
    expect(getQuotaStatusVariant("not_activated")).toBe("secondary");
    expect(getQuotaStatusVariant("error")).toBe("destructive");

    expect(getBucketProgressTone(80, false)).toBe("ok");
    expect(getBucketProgressTone(20, false)).toBe("warn");
    expect(getBucketProgressTone(0, true)).toBe("error");
  });

  test("groups quotas by provider with aggregated stats and average remaining percentage", () => {
    const groups = groupProviderQuotas(SAMPLE_QUOTAS);
    expect(groups).toHaveLength(2);

    const harborGroup = groups.find((group) => group.providerId === "tokenharbor");
    expect(harborGroup).toBeDefined();
    expect(harborGroup?.providerName).toBe("Token Harbor");
    expect(harborGroup?.fetcherId).toBe("tokenharbor");
    expect(harborGroup?.items).toHaveLength(4);
    expect(harborGroup?.stats).toEqual({
      total: 4,
      healthy: 1,
      low: 1,
      exhausted: 0,
      notActivated: 1,
      error: 1,
    });
    // conn-1 min remaining = 63, conn-2 = 15, conn-4 = 100 (conn-5 error has no buckets) -> avg = (63 + 15 + 100) / 3 = 59.3
    expect(harborGroup?.avgRemainingPct).toBeCloseTo(59.3, 1);

    const customGroup = groups.find((group) => group.providerId === "my-harbor");
    expect(customGroup).toBeDefined();
    expect(customGroup?.providerName).toBe("Custom Harbor");
    expect(customGroup?.items).toHaveLength(1);
    expect(customGroup?.avgRemainingPct).toBe(0);
  });

  test("patches a single connection quota immutably while preserving order", () => {
    const conn2 = SAMPLE_QUOTAS[1];
    const conn2Bucket = conn2?.buckets[0];
    if (!conn2 || !conn2Bucket) {
      throw new Error("expected sample connection and bucket");
    }

    const initial = {
      quotas: SAMPLE_QUOTAS.slice(0, 2),
      supported_providers: ["tokenharbor"],
      fetched_at: "2026-11-27T01:00:00Z",
    };

    const updatedConn2: ProviderConnectionQuota = {
      ...conn2,
      status: "exhausted",
      fetched_at: "2026-11-27T01:05:00Z",
      buckets: [
        {
          ...conn2Bucket,
          used_pct: 100,
          remaining_pct: 0,
          exhausted: true,
        },
      ],
    };

    const next = patchConnectionQuotaInResponse(initial, updatedConn2);
    expect(next.fetched_at).toBe("2026-11-27T01:05:00Z");
    expect(next.quotas.map((q) => q.connection_id)).toEqual(["conn-1", "conn-2"]);
    expect(next.quotas[1]?.status).toBe("exhausted");
    expect(next.quotas[1]?.buckets[0]?.remaining_pct).toBe(0);
    // Original object untouched
    expect(initial.quotas[1]?.status).toBe("low");
  });
});
