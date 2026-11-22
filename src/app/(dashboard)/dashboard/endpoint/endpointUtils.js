import { parseIntervalMs } from "@/shared/utils/interval";

const DAY_MS = 24 * 60 * 60 * 1000;
const EXPIRING_SOON_DAYS = 7;

export const BULK_RESET_INTERVAL_OPTIONS = [
  { value: "never", label: "No reset" },
  { value: "1h", label: "Every 1 Hour (1h)" },
  { value: "5h", label: "Every 5 Hours (5h)" },
  { value: "1d", label: "Every 1 Day (1d)" },
  { value: "7d", label: "Every 7 Days (7d)" },
  { value: "14d", label: "Every 14 Days (14d)" },
  { value: "30d", label: "Every 30 Days (30d)" },
];

export const RESET_INTERVAL_OPTIONS = [
  ...BULK_RESET_INTERVAL_OPTIONS,
  { value: "custom", label: "Custom Interval..." },
];

export function formatTokensNumber(num) {
  if (!num || num <= 0) return "0";
  if (num >= 1_000_000_000) return (num / 1_000_000_000).toFixed(2) + "B";
  if (num >= 1_000_000) return (num / 1_000_000).toFixed(2) + "M";
  if (num >= 1_000) return (num / 1_000).toFixed(1) + "K";
  return num.toLocaleString();
}

export function maskApiKey(fullKey) {
  if (!fullKey || fullKey.length <= 10) return fullKey || "";
  return fullKey.slice(0, 6) + "•".repeat(fullKey.length - 10) + fullKey.slice(-4);
}

export function normalizeOrigin(url) {
  if (!url || typeof url !== "string") return "";
  let cleaned = url.trim().replace(/\/+$/, "");
  cleaned = cleaned.replace(/\/v1$/i, "").replace(/\/+$/, "");
  return cleaned;
}

export function getExpiryMeta(expiresAt, nowMs = Date.now()) {
  if (!expiresAt || typeof expiresAt !== "string") {
    return { status: "none", label: "" };
  }
  const ts = new Date(expiresAt).getTime();
  if (!Number.isFinite(ts)) {
    return { status: "none", label: "" };
  }
  const diffMs = ts - nowMs;
  if (diffMs <= 0) {
    return { status: "expired", label: "Expired" };
  }
  const days = Math.ceil(diffMs / DAY_MS);
  if (days <= EXPIRING_SOON_DAYS) {
    return { status: "soon", label: `Expires in ${days}d` };
  }
  const formatted = new Date(ts).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
    timeZone: "UTC",
  });
  return { status: "ok", label: `Expires ${formatted}` };
}

export function isKeyExhausted(key) {
  if (!key) return false;
  const limit = Number(key.tokenLimit) || 0;
  const used = Number(key.usedTokens) || 0;
  return limit > 0 && used >= limit;
}

export function countKeyBuckets(keys = [], nowMs = Date.now()) {
  const counts = { all: keys.length, active: 0, disabled: 0, exhausted: 0, expiring: 0 };
  for (const key of keys) {
    const disabled = key?.isActive === false;
    if (disabled) counts.disabled += 1;
    else counts.active += 1;
    if (isKeyExhausted(key)) counts.exhausted += 1;
    const expiry = getExpiryMeta(key?.expiresAt, nowMs);
    if (expiry.status === "expired" || expiry.status === "soon") {
      counts.expiring += 1;
    }
  }
  return counts;
}

export function filterApiKeys(keys = [], { query = "", statusFilter = "all", nowMs = Date.now() } = {}) {
  const q = String(query || "").trim().toLowerCase();
  return keys.filter((key) => {
    if (q) {
      const nameMatch = key?.name && String(key.name).toLowerCase().includes(q);
      const keyMatch = key?.key && String(key.key).toLowerCase().includes(q);
      if (!nameMatch && !keyMatch) return false;
    }
    if (statusFilter === "active") return key?.isActive !== false;
    if (statusFilter === "disabled") return key?.isActive === false;
    if (statusFilter === "exhausted") return isKeyExhausted(key);
    if (statusFilter === "expiring") {
      const expiry = getExpiryMeta(key?.expiresAt, nowMs);
      return expiry.status === "expired" || expiry.status === "soon";
    }
    return true;
  });
}

export function resolveDefaultSnippetModel(modelAliases, activeProviders) {
  if (modelAliases && typeof modelAliases === "object") {
    const firstAlias = Object.keys(modelAliases).find((k) => k && k.trim());
    if (firstAlias) return firstAlias;
  }
  if (Array.isArray(activeProviders)) {
    for (const conn of activeProviders) {
      if (Array.isArray(conn?.models) && conn.models.length > 0 && conn.models[0]) {
        return String(conn.models[0]);
      }
    }
  }
  return "gpt-4o";
}

export function countWhitelistEntries(ipWhitelist) {
  if (!ipWhitelist || typeof ipWhitelist !== "string") return 0;
  return ipWhitelist
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean).length;
}

export function validateCustomInterval(val) {
  const trimmed = String(val || "").trim();
  if (!trimmed) {
    return "Enter a custom reset interval such as 30m, 10h, or 3d.";
  }
  if (!/^\d+[mhd]$/.test(trimmed) || parseIntervalMs(trimmed) <= 0) {
    return "Custom interval must use a positive number ending in m, h, or d (e.g. 30m, 10h, 3d).";
  }
  return null;
}

export function toDateTimeLocalValue(raw) {
  if (!raw || typeof raw !== "string") return "";
  const trimmed = raw.trim();
  if (!trimmed) return "";
  // Preserve already-naive datetime-local strings as-is, but convert any
  // timestamp with an explicit timezone (Z or +/-HH:MM) into local wall-clock time.
  const hasTimezone = /(?:[Zz]|[+-]\d{2}:?\d{2})$/.test(trimmed);
  if (!hasTimezone && /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}/.test(trimmed)) {
    return trimmed.slice(0, 16);
  }
  const d = new Date(trimmed);
  if (Number.isNaN(d.getTime())) return "";
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}
