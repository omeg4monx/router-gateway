// Automatic-backup config + status (channel: telegram | github). Secrets are
// encrypted at rest with the vault.key (see src/lib/vault/vaultKey.cjs). Lives in
// its own KV scope so tokens never travel through the settings blob; it also
// round-trips inside DB backups (see exportDb/importDb) so a restored instance
// keeps its schedule.
import vaultKey from "../../vault/vaultKey.cjs";
import { DATA_DIR } from "../../dataDir.js";
import { makeKv } from "../helpers/kvStore.js";

const kv = makeKv("autoBackup");

const LEGACY_SALT = "xrouter-backup-token-salt";

export function encryptSecret(plaintext) {
  if (typeof plaintext !== "string" || !plaintext) return plaintext;
  return vaultKey.seal(plaintext, vaultKey.getVaultKey(DATA_DIR));
}

export function decryptSecret(stored) {
  if (typeof stored !== "string" || !stored) return "";
  const [ivHex, tagHex, dataHex] = stored.split(":");
  if (!ivHex || !tagHex || !dataHex) return stored; // pre-encryption plaintext value
  try {
    return vaultKey.open(stored, [vaultKey.getVaultKey(DATA_DIR), vaultKey.legacyMachineKey(LEGACY_SALT)]);
  } catch {
    return ""; // unknown key / tampered value
  }
}

const DEFAULT_CONFIG = {
  enabled: false,
  channel: "telegram", // telegram | github
  intervalHours: 24,
  tgBotToken: "",
  tgChatId: "",
  ghToken: "",
  ghRepo: "", // owner/repo
  includeUsage: false, // Default false to prevent exceeding Telegram 45MB/50MB limit
  sections: ["settings", "providers", "apiKeys", "combos", "customModels"],
};

export async function getAutoBackupConfig() {
  const stored = { ...DEFAULT_CONFIG, ...(await kv.get("config", {})) };
  return {
    ...stored,
    tgBotToken: decryptSecret(stored.tgBotToken),
    ghToken: decryptSecret(stored.ghToken),
  };
}

// Patches may carry plaintext secrets; they are encrypted before storage.
// An empty-string secret means "keep the stored one" and must be omitted.
export async function setAutoBackupConfig(patch) {
  const current = { ...(await kv.get("config", {})) };
  const next = { ...current };
  for (const [key, value] of Object.entries(patch)) {
    if (key === "tgBotToken" || key === "ghToken") {
      if (typeof value === "string" && value) next[key] = encryptSecret(value.trim());
      continue;
    }
    next[key] = value;
  }
  await kv.set("config", next);
  return getAutoBackupConfig();
}

export async function getAutoBackupStatus() {
  return kv.get("status", null);
}

export async function setAutoBackupStatus(status) {
  await kv.set("status", status);
}
