import { describe, it, expect, vi, beforeEach } from "vitest";

const mocks = vi.hoisted(() => ({
  getSettings: vi.fn(),
  updateSettings: vi.fn(async (body) => ({ ...body })),
}));

vi.mock("@/lib/localDb", () => ({ getSettings: mocks.getSettings, updateSettings: mocks.updateSettings }));
vi.mock("@/lib/db/repos/settingsHistoryRepo.js", () => ({ snapshotSettings: vi.fn() }));
vi.mock("@/lib/network/outboundProxy", () => ({ applyOutboundProxyEnv: vi.fn() }));
vi.mock("open-sse/services/combo.js", () => ({ resetComboRotation: vi.fn() }));
vi.mock("@/lib/auth/sessionVersions", () => ({ bumpVersion: vi.fn(), sessionSubject: vi.fn(), getVersion: vi.fn() }));
vi.mock("@/lib/auth/auditLog", () => ({ audit: vi.fn() }));

const { PATCH } = await import("../../src/app/api/settings/route.js");
const headroom = await import("../../src/app/api/headroom/proxy/[...path]/route.js");

function patch(body) {
  return PATCH(new Request("http://localhost/api/settings", { method: "PATCH", body: JSON.stringify(body) }));
}

describe("first password set", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.unstubAllEnvs();
    mocks.getSettings.mockResolvedValue({}); // no stored hash yet
  });

  it("rejects without the current (initial) password", async () => {
    const res = await patch({ newPassword: "n3w-secret" });
    expect(res.status).toBe(400);
    expect(mocks.updateSettings).not.toHaveBeenCalled();
  });

  it("rejects the built-in default when INITIAL_PASSWORD is configured", async () => {
    vi.stubEnv("INITIAL_PASSWORD", "custom-initial");
    const res = await patch({ newPassword: "n3w-secret", currentPassword: "omeg4mon123" });
    expect(res.status).toBe(401);
  });

  it("accepts the configured INITIAL_PASSWORD", async () => {
    vi.stubEnv("INITIAL_PASSWORD", "custom-initial");
    const res = await patch({ newPassword: "n3w-secret", currentPassword: "custom-initial" });
    expect(res.status).toBe(200);
    expect(mocks.updateSettings.mock.calls[0][0].password).toMatch(/^\$2/);
  });
});

describe("headroom proxy", () => {
  it("never forwards dashboard credentials or custom-server peer headers", async () => {
    mocks.getSettings.mockResolvedValue({ headroomUrl: "http://127.0.0.1:8787" });
    const fetchSpy = vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response("ok"));
    await headroom.GET(
      new Request("http://localhost/api/headroom/proxy/stats", {
        headers: {
          cookie: "auth_token=jwt",
          authorization: "Bearer sk-x",
          "x-9r-peer-token": "peer-secret",
          "x-9r-real-ip": "203.0.113.7",
          accept: "application/json",
        },
      }),
      { params: Promise.resolve({ path: ["stats"] }) }
    );
    const sent = fetchSpy.mock.calls[0][1].headers;
    expect(sent.get("cookie")).toBeNull();
    expect(sent.get("authorization")).toBeNull();
    expect(sent.get("x-9r-peer-token")).toBeNull();
    expect(sent.get("x-9r-real-ip")).toBeNull();
    expect(sent.get("accept")).toBe("application/json");
    fetchSpy.mockRestore();
  });
});
