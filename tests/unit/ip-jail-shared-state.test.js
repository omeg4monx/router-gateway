import { describe, it, expect, vi, beforeEach } from "vitest";

// custom-server.js loads ipJail as native ESM while the Next bundle carries its
// own copy, in the same process. Two module instances must act as one jail.
const IP = "198.51.100.23";

async function freshJail() {
  vi.resetModules();
  return import("@/lib/security/ipJail.js");
}

beforeEach(async () => {
  (await freshJail()).unjailIp(IP);
});

describe("ipJail state is shared across module instances", () => {
  it("a jail set by one instance is enforced and lifted by another", async () => {
    const a = await freshJail();
    const b = await freshJail();
    expect(a).not.toBe(b);

    a.jailIp(IP, "test");
    expect(b.isIpJailed(IP)).toBe(true);

    b.unjailIp(IP);
    expect(a.isIpJailed(IP)).toBe(false);
  });

  it("login failures counted in different instances add up to one jail", async () => {
    const a = await freshJail();
    const b = await freshJail();
    for (let i = 0; i < 5; i++) (i % 2 ? a : b).recordLoginFailure(IP);
    expect(a.isIpJailed(IP)).toBe(true);
  });
});

describe("recording a jail does not count as a probe strike", () => {
  it("two probes plus the jail's own critical event stay below the threshold", async () => {
    const jail = await freshJail();
    const { recordSecurityEvent } = await import("@/lib/db/repos/securityLogRepo.js");

    jail.recordProbeViolation(IP, "Secret file probe");
    jail.recordProbeViolation(IP, "Secret file probe");
    // What jailIp writes for any jail; its reason can mention "password".
    await recordSecurityEvent({
      type: "ip.jailed",
      severity: "critical",
      ip: IP,
      path: "/api/security/jail",
      status: 403,
      detail: JSON.stringify({ reason: "Failed dashboard password attempt" }),
    });

    expect(jail.isIpJailed(IP)).toBe(false);
  });

  it("an explicit request probe still strikes", async () => {
    const jail = await freshJail();
    const { recordSecurityEvent } = await import("@/lib/db/repos/securityLogRepo.js");
    for (let i = 0; i < 3; i++) {
      await recordSecurityEvent({ type: "request_probe", severity: "critical", ip: IP, path: "/.env", status: 404 });
    }
    expect(jail.isIpJailed(IP)).toBe(true);
  });
});
