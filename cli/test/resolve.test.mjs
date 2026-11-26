import { describe, expect, test } from "bun:test";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import { binaryFileName, resolveBinaryPath } from "../lib/resolve.js";

describe("cli binary resolver", () => {
  test("maps supported platform and arch pairs to release binary names", () => {
    expect(binaryFileName("linux", "x64")).toBe("omeg4mon-router-linux-amd64");
    expect(binaryFileName("linux", "arm64")).toBe("omeg4mon-router-linux-arm64");
    expect(binaryFileName("darwin", "x64")).toBe("omeg4mon-router-darwin-amd64");
    expect(binaryFileName("darwin", "arm64")).toBe("omeg4mon-router-darwin-arm64");
    expect(binaryFileName("win32", "x64")).toBe("omeg4mon-router-windows-amd64.exe");
  });

  test("throws on unsupported platform or architecture", () => {
    expect(() => binaryFileName("freebsd", "x64")).toThrow(/Unsupported platform/);
    expect(() => binaryFileName("linux", "ia32")).toThrow(/Unsupported architecture/);
  });

  test("honors OMEG4MON_ROUTER_BIN when present and executable", () => {
    const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "om-cli-test-"));
    const fakeBin = path.join(tmpDir, "omeg4mon-router-custom");
    fs.writeFileSync(fakeBin, "#!/bin/sh\nexit 0\n", { mode: 0o755 });
    try {
      const resolved = resolveBinaryPath({
        env: { OMEG4MON_ROUTER_BIN: fakeBin },
        platform: "linux",
        arch: "x64",
        pkgRoot: tmpDir,
      });
      expect(resolved).toBe(fakeBin);
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });

  test("finds platform binary inside vendor directory before fallback paths", () => {
    const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "om-cli-vendor-"));
    const vendorDir = path.join(tmpDir, "vendor");
    fs.mkdirSync(vendorDir, { recursive: true });
    const vendorBin = path.join(vendorDir, "omeg4mon-router-linux-amd64");
    fs.writeFileSync(vendorBin, "#!/bin/sh\nexit 0\n", { mode: 0o755 });
    try {
      const resolved = resolveBinaryPath({
        env: {},
        platform: "linux",
        arch: "x64",
        pkgRoot: tmpDir,
      });
      expect(resolved).toBe(vendorBin);
    } finally {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    }
  });
});
