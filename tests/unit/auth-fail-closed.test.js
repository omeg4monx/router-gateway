// Boots a real HTTP server through custom-server.js and checks its auth guard fails
// CLOSED: a request whose guard cannot load or throws is answered 503 and never
// reaches the Next handler. Origin: a crafted Cookie header once made the guard
// throw and the wrapper then handed the request to Next with no auth at all.
import { describe, it, expect, beforeAll, afterAll, afterEach, vi } from "vitest";
import { createRequire } from "node:module";
import http from "node:http";

const require = createRequire(import.meta.url);
const fs = require("node:fs");
const nodeModule = require("node:module");

// Not a valid percent-encoding: decodeURIComponent throws on it.
const BAD_COOKIE = "auth_token=%E0%A4%A";

let server;
let baseUrl;
let nextCalls;
let nextHandler;

beforeAll(async () => {
  // custom-server.js applies the checkout's .env on load, which can point DATA_DIR
  // at a real data dir or BASE_URL at a live port. Keep the test's environment.
  const env = { ...process.env };
  require("../../custom-server.js");
  const peerToken = process.env.NINEROUTER_PEER_TOKEN;
  for (const key of Object.keys(process.env)) if (!(key in env)) delete process.env[key];
  Object.assign(process.env, env, { NINEROUTER_PEER_TOKEN: peerToken });

  server = http.createServer((req, res) => {
    nextCalls++;
    return nextHandler(req, res);
  });
  await new Promise((resolve) => server.listen(0, "127.0.0.1", resolve));
  baseUrl = `http://127.0.0.1:${server.address().port}`;
});

afterAll(async () => {
  await new Promise((resolve) => server.close(resolve));
});

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});

async function get(pathname, headers = {}, handler = (req, res) => res.end("next")) {
  nextCalls = 0;
  nextHandler = handler;
  const res = await fetch(baseUrl + pathname, { headers, redirect: "manual" });
  await res.arrayBuffer();
  return { status: res.status, nextCalls };
}

describe("custom-server auth guard fails closed", () => {
  // Runs first: once the guard has loaded it stays cached.
  it("answers 503 while the guard cannot load, then loads it on the next request", async () => {
    const existsSync = fs.existsSync;
    vi.spyOn(fs, "existsSync").mockImplementation(
      (p) => !String(p).endsWith("dashboardGuard.js") && existsSync(p),
    );
    expect(await get("/api/settings/database")).toEqual({ status: 503, nextCalls: 0 });

    vi.restoreAllMocks();
    expect(await get("/api/settings/database")).toEqual({ status: 401, nextCalls: 0 });
  });

  it("loads the guard once, not once per request", async () => {
    const register = vi.spyOn(nodeModule, "register");
    expect(await get("/claim")).toEqual({ status: 200, nextCalls: 1 });
    expect(await get("/claim")).toEqual({ status: 200, nextCalls: 1 });
    expect(register).not.toHaveBeenCalled();
  });

  it("answers 503 instead of passing through when the guard step throws", async () => {
    vi.stubGlobal("Headers", class {
      constructor() {
        throw new Error("boom");
      }
    });
    expect(await get("/api/settings/database")).toEqual({ status: 503, nextCalls: 0 });
  });

  it("answers 503 and never re-forwards when the wrapper chain rejects", async () => {
    const crash = () => {
      throw new Error("handler crashed");
    };
    expect(await get("/claim", {}, crash)).toEqual({ status: 503, nextCalls: 1 });
  });

  it("ignores a cookie that does not decode instead of throwing", async () => {
    // A public page still opens...
    expect(await get("/claim", { cookie: BAD_COOKIE })).toEqual({ status: 200, nextCalls: 1 });
    // ...and a protected route sees no session rather than a broken one.
    expect(await get("/api/settings/database", { cookie: BAD_COOKIE })).toEqual({ status: 401, nextCalls: 0 });
  });
});
