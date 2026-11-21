import vaultKey from "../../vault/vaultKey.cjs";
import { DATA_DIR } from "../../dataDir.js";
import { makeKv } from "../helpers/kvStore.js";

const kv = makeKv("connection");

export function encryptSecret(plaintext) {
  if (typeof plaintext !== "string" || !plaintext) return plaintext;
  return vaultKey.seal(plaintext, vaultKey.getVaultKey(DATA_DIR));
}

export function decryptSecret(stored) {
  if (typeof stored !== "string" || !stored) return "";
  if (stored.split(":").length !== 3) return stored; // unencrypted / plaintext legacy
  try {
    return vaultKey.open(stored, vaultKey.getVaultKey(DATA_DIR));
  } catch {
    return "";
  }
}

const DEFAULT_CONNECTION_CONFIG = {
  telegram: {
    enabled: false,
    botToken: "",
    userId: "",
    webhookUrl: "",
    autoWebhook: true,
    botInfo: null, // { id, username, first_name }
    status: "idle", // idle | connected | error
    lastError: null,
    lastTestedAt: null,
  },
};

export async function getConnectionConfig() {
  const stored = await kv.get("config", {});
  const tgStored = stored.telegram || {};
  return {
    telegram: {
      ...DEFAULT_CONNECTION_CONFIG.telegram,
      ...tgStored,
      botToken: decryptSecret(tgStored.botToken || ""),
    },
  };
}

export async function setConnectionConfig(patch = {}) {
  const current = await getConnectionConfig();
  const next = {
    telegram: {
      ...current.telegram,
      ...(patch.telegram || {}),
    },
  };

  const toStore = {
    telegram: {
      ...next.telegram,
    },
  };

  if (typeof patch.telegram?.botToken === "string") {
    const raw = patch.telegram.botToken.trim();
    toStore.telegram.botToken = raw ? encryptSecret(raw) : "";
  } else {
    toStore.telegram.botToken = next.telegram.botToken ? encryptSecret(next.telegram.botToken) : "";
  }

  await kv.set("config", toStore);
  return getConnectionConfig();
}
