// Single source of the at-rest encryption key. CommonJS so the MITM server
// (plain node, CJS) and the Next app (ESM) share one implementation.
const crypto = require("crypto");
const fs = require("fs");
const path = require("path");

const ALGO = "aes-256-gcm";
const keyCache = new Map();

// Random 32-byte key in <dataDir>/vault.key (0600), or XROUTER_VAULT_KEY.
// Unlike the old machine-id key, copying the DB alone is not enough to decrypt.
function getVaultKey(dataDir) {
  const fromEnv = process.env.XROUTER_VAULT_KEY?.trim();
  if (fromEnv) return crypto.createHash("sha256").update(fromEnv).digest();

  const file = path.join(dataDir, "vault.key");
  if (keyCache.has(file)) return keyCache.get(file);

  let hex;
  try {
    hex = fs.readFileSync(file, "utf8").trim();
  } catch (e) {
    if (e.code !== "ENOENT") throw e;
    fs.mkdirSync(dataDir, { recursive: true });
    hex = crypto.randomBytes(32).toString("hex");
    try {
      fs.writeFileSync(file, hex, { mode: 0o600, flag: "wx" });
    } catch (w) {
      if (w.code !== "EEXIST") throw w;
      hex = fs.readFileSync(file, "utf8").trim(); // another process created it first
    }
  }
  if (!/^[0-9a-f]{64}$/.test(hex)) throw new Error(`[vault] ${file} is not a 64-char hex key`);
  const key = Buffer.from(hex, "hex");
  keyCache.set(file, key);
  return key;
}

// Pre-vault.key scheme, kept only to read values written before the switch.
function legacyMachineKey(salt) {
  try {
    const { machineIdSync } = require("node-machine-id");
    return crypto.createHash("sha256").update(machineIdSync() + salt).digest();
  } catch {
    return crypto.createHash("sha256").update(salt).digest();
  }
}

// "<iv>:<tag>:<ciphertext>" hex, the format every caller already stores.
function seal(plaintext, key) {
  const iv = crypto.randomBytes(12);
  const cipher = crypto.createCipheriv(ALGO, key, iv);
  const data = Buffer.concat([cipher.update(plaintext, "utf8"), cipher.final()]);
  return `${iv.toString("hex")}:${cipher.getAuthTag().toString("hex")}:${data.toString("hex")}`;
}

// Tries each key in order; throws if none authenticates the value.
function open(stored, keys) {
  const [ivHex, tagHex, dataHex] = stored.split(":");
  let lastErr;
  for (const key of keys) {
    try {
      const decipher = crypto.createDecipheriv(ALGO, key, Buffer.from(ivHex, "hex"));
      decipher.setAuthTag(Buffer.from(tagHex, "hex"));
      return decipher.update(Buffer.from(dataHex, "hex"), undefined, "utf8") + decipher.final("utf8");
    } catch (e) {
      lastErr = e;
    }
  }
  throw lastErr;
}

module.exports = { getVaultKey, legacyMachineKey, seal, open };
