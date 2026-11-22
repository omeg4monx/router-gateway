// The voucher claim route only serves remote callers that came through Cloudflare.
// custom-server.js strips cf-ray on direct sockets, but without the wrapper the
// header is client-controlled, so it may only count when the wrapper vouched.
import { describe, it, expect, vi, beforeEach } from "vitest";
import { NextRequest } from "next/server";

vi.mock("@/lib/auth/auditLog", () => ({ audit: vi.fn() }));
vi.mock("@/lib/localDb", () => ({
  claimVoucher: vi.fn(),
  getVoucherByCode: vi.fn(),
  getActiveBansosVoucher: vi.fn(async () => null),
}));
vi.mock("@/dashboardGuard", () => ({ isLocalRequest: vi.fn(() => false) }));
vi.mock("@/lib/security/proxyIpBlocker", () => ({ isKnownProxyIp: vi.fn(async () => false) }));
vi.mock("@/lib/security/antiAbuseEngine", () => ({
  normalizeClientIp: vi.fn((ip) => ip),
  verifyNetworkThrottle: vi.fn(() => ({ allowed: true })),
  verifyClientFingerprint: vi.fn(),
  verifySubmissionIntegrity: vi.fn(),
  mintChallenge: vi.fn(async () => ({ challenge: "c" })),
  recordFailure: vi.fn(),
  recordSuccess: vi.fn(),
  acquireConcurrencyLock: vi.fn(),
  releaseConcurrencyLock: vi.fn(),
  clearAllLockouts: vi.fn(),
}));

const { GET, POST } = await import("../../src/app/api/vouchers/claim/route.js");

const PEER_TOKEN = "peer-token-fixture";
const URL_CLAIM = "http://router.example.com/api/vouchers/claim";

function claimRequest(method, headers) {
  return new NextRequest(URL_CLAIM, { method, headers });
}

describe("voucher claim Cloudflare edge gate", () => {
  beforeEach(() => {
    process.env.NINEROUTER_PEER_TOKEN = PEER_TOKEN;
  });

  it("rejects a GET whose cf-ray the wrapper never vouched for", async () => {
    const res = await GET(claimRequest("GET", { "cf-ray": "8f00-SIN" }));

    expect(res.status).toBe(403);
  });

  it("rejects a POST whose cf-ray the wrapper never vouched for", async () => {
    const res = await POST(claimRequest("POST", { "cf-ray": "8f00-SIN" }));

    expect(res.status).toBe(403);
  });

  it("serves a GET whose cf-ray arrived through the wrapper", async () => {
    const res = await GET(claimRequest("GET", {
      "cf-ray": "8f00-SIN",
      "x-9r-real-ip": "203.0.113.9",
      "x-9r-peer-token": PEER_TOKEN,
    }));

    expect(res.status).toBe(200);
    expect((await res.json()).hasActiveBansos).toBe(false);
  });
});
