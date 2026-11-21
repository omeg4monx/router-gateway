import fs from "node:fs";
import path from "node:path";
import crypto from "node:crypto";
import { describe, it, expect, vi } from "vitest";
import { machineIdSync } from "node-machine-id";
import { sealCredential, unsealCredential, isSealed } from "../../src/lib/vault/credentialVault.js";
import { DATA_DIR } from "../../src/lib/dataDir.js";

function sealWith(key, plaintext) {
  const iv = crypto.randomBytes(12);
  const cipher = crypto.createCipheriv("aes-256-gcm", key, iv);
  const data = Buffer.concat([cipher.update(plaintext, "utf8"), cipher.final()]);
  return `enc:v1:${iv.toString("hex")}:${cipher.getAuthTag().toString("hex")}:${data.toString("hex")}`;
}

describe("credential vault (AES-256-GCM)", () => {
  it("seals plaintext credentials into enc:v1 format", () => {
    const secret = "sk-test-secret-key-12345";
    const sealed = sealCredential(secret);

    expect(sealed).not.toBe(secret);
    expect(isSealed(sealed)).toBe(true);
    expect(sealed.startsWith("enc:v1:")).toBe(true);
  });

  it("unseals ciphertext back to original plaintext", () => {
    const secret = "gho_sampleGithubOAuthToken123456789";
    const sealed = sealCredential(secret);
    const unsealed = unsealCredential(sealed);

    expect(unsealed).toBe(secret);
  });

  it("leaves unsealed legacy plaintext untouched (backward compatibility)", () => {
    const legacy = "plaintext-api-key-from-old-db";
    expect(unsealCredential(legacy)).toBe(legacy);
  });

  it("keys from a private random vault.key, not the machine id", () => {
    sealCredential("x");
    const keyFile = path.join(DATA_DIR, "vault.key");
    expect(fs.readFileSync(keyFile, "utf8").trim()).toMatch(/^[0-9a-f]{64}$/);
    if (process.platform !== "win32") expect(fs.statSync(keyFile).mode & 0o077).toBe(0);

    // The machine id is not key material.
    const machineKey = crypto.createHash("sha256").update(machineIdSync()).digest();
    const [, , iv, tag, data] = sealCredential("sk-new").split(":");
    const decipher = crypto.createDecipheriv("aes-256-gcm", machineKey, Buffer.from(iv, "hex"));
    decipher.setAuthTag(Buffer.from(tag, "hex"));
    expect(() => decipher.update(Buffer.from(data, "hex")) + decipher.final("utf8")).toThrow();
  });

  it("uses OMEG4MON_ROUTER_VAULT_KEY when set", () => {
    const fixture = "env-test-fixture";
    const stored = sealWith(crypto.createHash("sha256").update(fixture).digest(), "stored-credential");
    vi.stubEnv("OMEG4MON_ROUTER_VAULT_KEY", fixture);
    try {
      expect(unsealCredential(stored)).toBe("stored-credential");
    } finally {
      vi.unstubAllEnvs();
    }
  });

  it("does not open values sealed under another key", () => {
    const stored = sealWith(crypto.randomBytes(32), "sk-foreign");
    expect(unsealCredential(stored)).toBe(stored);
  });
});
