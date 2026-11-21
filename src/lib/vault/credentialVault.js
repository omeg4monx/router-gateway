import crypto from "node:crypto";
import vaultKey from "./vaultKey.cjs";
import { DATA_DIR } from "../dataDir.js";

const VAULT_SALT = "xrouter-credential-vault-salt-2026";

let customMasterKey = null;

export function setVaultMasterKey(key) {
  if (typeof key === "string" && key.trim()) {
    customMasterKey = crypto.createHash("sha256").update(key.trim()).digest();
  }
}

function primaryKey() {
  return customMasterKey || vaultKey.getVaultKey(DATA_DIR);
}

/**
 * Encrypt plaintext credential string using AES-256-GCM.
 * Output format: "enc:v1:<iv_hex>:<tag_hex>:<ciphertext_hex>"
 */
export function sealCredential(plaintext) {
  if (typeof plaintext !== "string" || !plaintext) return plaintext;
  if (plaintext.startsWith("enc:v1:")) return plaintext; // Already sealed

  return `enc:v1:${vaultKey.seal(plaintext, primaryKey())}`;
}

/**
 * Decrypt credential sealed with sealCredential.
 * Backward compatible with plaintext legacy values.
 */
export function unsealCredential(stored) {
  if (typeof stored !== "string" || !stored) return stored;
  if (!stored.startsWith("enc:v1:")) return stored; // Plaintext legacy, return as-is

  const parts = stored.split(":");
  if (parts.length !== 5) return stored;

  try {
    // Legacy machine-id key: values sealed before vault.key existed stay readable
    // and are re-sealed with the new key on their next write.
    return vaultKey.open(parts.slice(2).join(":"), [primaryKey(), vaultKey.legacyMachineKey(VAULT_SALT)]);
  } catch (err) {
    console.error("[CredentialVault] Failed to unseal credential:", err.message);
    return stored;
  }
}

/**
 * Check if a credential value is currently sealed with AES-256-GCM.
 */
export function isSealed(value) {
  return typeof value === "string" && value.startsWith("enc:v1:");
}
