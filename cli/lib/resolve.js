import fs from "node:fs";
import path from "node:path";

const PLATFORM_MAP = {
  linux: "linux",
  darwin: "darwin",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

/**
 * Returns the canonical release binary filename for a given Node platform & arch.
 * @param {string} platform
 * @param {string} arch
 * @returns {string}
 */
export function binaryFileName(platform, arch) {
  const goos = PLATFORM_MAP[platform];
  if (!goos) {
    throw new Error(`Unsupported platform: ${platform}`);
  }
  const goarch = ARCH_MAP[arch];
  if (!goarch) {
    throw new Error(`Unsupported architecture: ${arch}`);
  }
  const ext = goos === "windows" ? ".exe" : "";
  return `omeg4mon-router-${goos}-${goarch}${ext}`;
}

function isExecutableFile(candidatePath) {
  try {
    const stat = fs.statSync(candidatePath);
    if (!stat.isFile()) {
      return false;
    }
    if (process.platform === "win32") {
      return true;
    }
    fs.accessSync(candidatePath, fs.constants.X_OK);
    return true;
  } catch {
    return false;
  }
}

/**
 * Resolves the omeg4mon-router single binary path using deterministic precedence:
 *   1. OMEG4MON_ROUTER_BIN environment override
 *   2. <pkgRoot>/vendor/<platform-binary>
 *   3. <pkgRoot>/../bin/omeg4mon-router[.exe] (local dev build)
 *   4. <pkgRoot>/../dist/<platform-binary> (cross-build output)
 *   5. omeg4mon-router on PATH
 */
export function resolveBinaryPath({
  env = process.env,
  platform = process.platform,
  arch = process.arch,
  pkgRoot,
}) {
  const envOverride = env?.OMEG4MON_ROUTER_BIN?.trim();
  if (envOverride) {
    if (!isExecutableFile(envOverride)) {
      throw new Error(`OMEG4MON_ROUTER_BIN is not an executable file: ${envOverride}`);
    }
    return envOverride;
  }

  const binName = binaryFileName(platform, arch);
  const localDevName = platform === "win32" ? "omeg4mon-router.exe" : "omeg4mon-router";
  const candidates = [
    path.join(pkgRoot, "vendor", binName),
    path.resolve(pkgRoot, "..", "bin", localDevName),
    path.resolve(pkgRoot, "..", "dist", binName),
  ];

  for (const candidate of candidates) {
    if (isExecutableFile(candidate)) {
      return candidate;
    }
  }

  const pathEnv = env?.PATH || "";
  for (const dir of pathEnv.split(path.delimiter)) {
    if (!dir) continue;
    const candidate = path.join(dir, localDevName);
    if (isExecutableFile(candidate)) {
      return candidate;
    }
  }

  throw new Error(
    `Could not locate omeg4mon-router binary (${binName}). Run 'make build' or set OMEG4MON_ROUTER_BIN.`,
  );
}
