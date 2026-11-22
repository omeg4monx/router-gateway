import { describe, it, expect } from "vitest";
import { matchRequestProbe } from "../../src/lib/db/repos/securityLogRepo.js";

describe("matchRequestProbe: request URLs that feed the auto-jail", () => {
  it.each([
    ["/.env", "Secret file probe"],
    ["/api/.env?x=1", "Secret file probe"],
    ["/static/../../etc/passwd", "Path traversal attempt"],
    ["/files/%2e%2e/%2e%2e/win.ini", "Path traversal attempt"],
    ["/search?q=<script>alert(1)</script>", "Script injection attempt"],
    ["/search?q=%3Cscript%3Ealert(1)%3C%2Fscript%3E", "Script injection attempt"],
    ["/item?id=1 union select password from users", "SQL injection attempt"],
  ])("flags %s", (url, label) => {
    expect(matchRequestProbe(url)).toBe(label);
  });

  it.each([
    "/api/auth/reset-password",
    "/api/oauth/iflow/cookie",
    "/dashboard/providers",
    "/v1/chat/completions",
    "/login?reason=idle_timeout",
  ])("leaves the real route %s alone", (url) => {
    // Route names that mention auth material are not probes; three strikes on
    // them would jail the operator.
    expect(matchRequestProbe(url)).toBeNull();
  });

  it("gives the same answer on repeated calls (global regex state)", () => {
    expect(matchRequestProbe("/.env")).toBe("Secret file probe");
    expect(matchRequestProbe("/.env")).toBe("Secret file probe");
  });
});
