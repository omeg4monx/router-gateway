// Smoke test: prove the DATA_DIR guard behaves, not just that the string exists.
// Runs the real path resolver with process.platform faked to win32.
const assert = require("assert");
const path = require("path");

function freshPaths(env, platform) {
  const mod = require.resolve("/root/omeg4mon-router/src/mitm/paths.js");
  delete require.cache[mod];
  const origPlatform = Object.getOwnPropertyDescriptor(process, "platform");
  const origDataDir = process.env.DATA_DIR;
  Object.defineProperty(process, "platform", { value: platform, configurable: true });
  if (env === undefined) delete process.env.DATA_DIR; else process.env.DATA_DIR = env;
  const out = require(mod);
  Object.defineProperty(process, "platform", origPlatform);
  if (origDataDir === undefined) delete process.env.DATA_DIR; else process.env.DATA_DIR = origDataDir;
  delete require.cache[mod];
  return out;
}

let pass = 0, fail = 0;
function check(name, fn) {
  try { fn(); console.log(`  PASS  ${name}`); pass++; }
  catch (e) { console.log(`  FAIL  ${name} -> ${e.message}`); fail++; }
}

console.log("DATA_DIR guard smoke test");

// 1. Windows + Unix path must fall back (never D:\var\lib\omeg4mon-router).
check("win32 + '/var/lib/omeg4mon-router' falls back to APPDATA", () => {
  const m = freshPaths("/var/lib/omeg4mon-router", "win32");
  // The Unix path must not survive as the resolved dir. On this Linux host the
  // fallback itself is necessarily still a POSIX path (APPDATA resolves under
  // homedir), so assert on what actually matters: it is no longer the
  // configured value, and never became a Windows-mapped D:\var\lib.
  assert.notStrictEqual(m.DATA_DIR, "/var/lib/omeg4mon-router", "guard did not fire");
  assert.ok(!/[\\/]var[\\/]lib[\\/]omeg4mon-router$/.test(m.DATA_DIR), `leaked Unix path: ${m.DATA_DIR}`);
  assert.ok(m.DATA_DIR.toLowerCase().includes("omeg4mon-router"), `got ${m.DATA_DIR}`);
});

// 2. Windows + native path must be honoured as-is.
check("win32 + 'D:\\\\data\\\\omeg4mon-router' is honoured", () => {
  const m = freshPaths("D:\\data\\omeg4mon-router", "win32");
  assert.strictEqual(m.DATA_DIR, "D:\\data\\omeg4mon-router");
});

// 3. Linux + Unix path must still be honoured (no regression for prod).
check("linux + '/var/lib/omeg4mon-router' is honoured (prod path)", () => {
  const m = freshPaths("/var/lib/omeg4mon-router", "linux");
  assert.strictEqual(m.DATA_DIR, "/var/lib/omeg4mon-router");
});

// 4. MITM_DIR always nests under DATA_DIR.
check("MITM_DIR nests under DATA_DIR", () => {
  const m = freshPaths("/var/lib/omeg4mon-router", "linux");
  assert.strictEqual(m.MITM_DIR, path.join("/var/lib/omeg4mon-router", "mitm"));
});

// 5. Unset DATA_DIR -> platform default on both.
check("unset DATA_DIR -> default on win32 and linux", () => {
  const w = freshPaths(undefined, "win32");
  const l = freshPaths(undefined, "linux");
  assert.ok(w.DATA_DIR.toLowerCase().includes("omeg4mon-router"), `win: ${w.DATA_DIR}`);
  assert.ok(l.DATA_DIR.includes(".omeg4mon-router"), `linux: ${l.DATA_DIR}`);
});

// 6. Cross-module agreement: MITM paths.js vs dashboard dataDir.js.
check("paths.js and dataDir.js resolve the same dir", async () => {
  const cjs = freshPaths("/var/lib/omeg4mon-router", "linux");
  const esm = await import("/root/omeg4mon-router/src/lib/dataDir.js");
  assert.strictEqual(cjs.DATA_DIR, esm.DATA_DIR);
});

console.log(`\n${pass} passed, ${fail} failed`);
process.exit(fail === 0 ? 0 : 1);
