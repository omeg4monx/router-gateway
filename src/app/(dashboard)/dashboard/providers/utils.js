import { getErrorCode, getRelativeTime } from "@/shared/utils";
import { OAUTH_PROVIDERS, APIKEY_PROVIDERS } from "@/shared/constants/config";
import {
  FREE_PROVIDERS,
  FREE_TIER_PROVIDERS,
  WEB_COOKIE_PROVIDERS,
  OPENAI_COMPATIBLE_PREFIX,
  ANTHROPIC_COMPATIBLE_PREFIX,
  getProviderAlias,
} from "@/shared/constants/providers";

export const STATUS_FILTER_OPTIONS = [
  { value: "all", label: "All" },
  { value: "active", label: "Active" },
  { value: "inactive", label: "Inactive" },
  { value: "none", label: "No connection" },
];

export const PROVIDER_STATUS_PILLS = [
  { value: "all", label: "All" },
  { value: "active", label: "Active" },
  { value: "error", label: "Errors" },
  { value: "inactive", label: "Inactive" },
  { value: "none", label: "Not Setup" },
];

export const CATEGORY_TABS = [
  { value: "all", label: "All Providers", icon: "dns" },
  { value: "custom", label: "Custom", icon: "terminal" },
  { value: "oauth", label: "OAuth", icon: "verified_user" },
  { value: "apikey", label: "API Keys", icon: "key" },
  { value: "free", label: "Free & Web", icon: "public" },
];

// noAuth providers (e.g. free proxies) are always usable even though they
// never have a stored connection record, so they never fall into "none".
export function getConnectionStatus(stats, isNoAuth = false) {
  if (isNoAuth) return "active";
  if (!stats || stats.total === 0) return "none";
  return stats.allDisabled ? "inactive" : "active";
}

export function matchesStatusFilter(statusFilter, stats, isNoAuth = false) {
  if (statusFilter === "all") return true;
  if (statusFilter === "error") {
    return Boolean(!isNoAuth && stats && stats.error > 0 && !stats.allDisabled);
  }
  return getConnectionStatus(stats, isNoAuth) === statusFilter;
}

export function getConnectionErrorTag(connection) {
  if (!connection) return null;

  const explicitType = connection.lastErrorType;
  if (explicitType === "runtime_error") return "RUNTIME";
  if (
    explicitType === "upstream_auth_error" ||
    explicitType === "auth_missing" ||
    explicitType === "token_refresh_failed" ||
    explicitType === "token_expired"
  ) {
    return "AUTH";
  }
  if (explicitType === "upstream_rate_limited") return "429";
  if (explicitType === "upstream_unavailable") return "5XX";
  if (explicitType === "network_error") return "NET";

  const numericCode = Number(connection.errorCode);
  if (Number.isFinite(numericCode) && numericCode >= 400) {
    return String(numericCode);
  }

  const fromMessage = getErrorCode(connection.lastError);
  if (fromMessage === "401" || fromMessage === "403") return "AUTH";
  if (fromMessage && fromMessage !== "ERR") return fromMessage;

  const msg = (connection.lastError || "").toLowerCase();
  if (
    msg.includes("runtime") ||
    msg.includes("not runnable") ||
    msg.includes("not installed")
  ) {
    return "RUNTIME";
  }
  if (
    msg.includes("invalid api key") ||
    msg.includes("token invalid") ||
    msg.includes("revoked") ||
    msg.includes("unauthorized")
  ) {
    return "AUTH";
  }

  return "ERR";
}

export function getProviderStats(connections, providerId, authType, nowMs = Date.now()) {
  const authTypes = Array.isArray(authType) ? authType : [authType];
  const providerConnections = (connections || []).filter(
    (c) => c.provider === providerId && authTypes.includes(c.authType),
  );

  const getEffectiveStatus = (conn) => {
    const isCooldown = Object.entries(conn).some(
      ([k, v]) =>
        k.startsWith("modelLock_") && v && new Date(v).getTime() > nowMs,
    );
    return conn.testStatus === "unavailable" && !isCooldown
      ? "active"
      : conn.testStatus;
  };

  const enabledConns = providerConnections.filter((c) => c.isActive !== false);

  const connected = enabledConns.filter((c) => {
    const status = getEffectiveStatus(c);
    return status === "active" || status === "success";
  }).length;

  const errorConns = enabledConns.filter((c) => {
    const status = getEffectiveStatus(c);
    return (
      status === "error" || status === "expired" || status === "unavailable"
    );
  });

  const error = errorConns.length;
  const total = providerConnections.length;
  const enabled = enabledConns.length;
  const allDisabled = total > 0 && enabled === 0;

  const latestError = [...errorConns].sort(
    (a, b) => new Date(b.lastErrorAt || 0).getTime() - new Date(a.lastErrorAt || 0).getTime(),
  )[0];
  const errorCode = latestError ? getConnectionErrorTag(latestError) : null;
  const errorTime = latestError?.lastErrorAt
    ? getRelativeTime(latestError.lastErrorAt)
    : null;
  const latestErrorAtMs = latestError?.lastErrorAt
    ? new Date(latestError.lastErrorAt).getTime() || 0
    : 0;

  return {
    connected,
    error,
    total,
    enabled,
    errorCode,
    errorTime,
    latestErrorAtMs,
    allDisabled,
  };
}

// Dual-auth providers (oauth + apikey) store API keys as authType "apikey"
// (and sometimes "api_key"). Card stats must count both so totals match detail.
export function resolveProviderAuthTypes(info, key) {
  if (key === "kiro") return ["oauth", "apikey", "api_key"];
  const modes = info?.authModes;
  if (!Array.isArray(modes)) {
    return key in FREE_TIER_PROVIDERS || key in APIKEY_PROVIDERS
      ? ["oauth", "apikey", "api_key"]
      : "oauth";
  }
  if (!modes.includes("apikey")) return "oauth";
  return ["oauth", "apikey", "api_key"];
}

export function isMoonshotCompatible(provider, providerId, isCompatible) {
  return Boolean(
    isCompatible &&
      (provider?.brand === "moonshot" ||
        providerId?.startsWith("moonshot-") ||
        /moonshot|kimi/i.test(provider?.name || "")),
  );
}

export function getProviderTypeLabel(provider, providerId, authType) {
  const isCompatible =
    authType === "compatible"
      ? !providerId?.startsWith(ANTHROPIC_COMPATIBLE_PREFIX)
      : providerId?.startsWith(OPENAI_COMPATIBLE_PREFIX);
  const isAnthropicCompatible =
    providerId?.startsWith(ANTHROPIC_COMPATIBLE_PREFIX) ||
    provider?.type === "anthropic-compatible";

  if (isAnthropicCompatible) return "Anthropic Compatible";
  if (isMoonshotCompatible(provider, providerId, isCompatible)) {
    return "MoonshotAI Compatible";
  }
  if (isCompatible) {
    return provider?.apiType === "responses"
      ? "OpenAI Responses"
      : "OpenAI Compatible";
  }
  if (authType === "oauth") return "OAuth 2.0";
  if (authType === "free") return "Free / Public";
  if (authType === "cookie") return "Session Cookie";
  return "API Key";
}

export function matchesProviderSearch(query, meta) {
  const q = (query || "").trim().toLowerCase();
  if (!q) return true;
  if (!meta) return false;
  const fields = [
    meta.name,
    meta.id,
    meta.alias,
    meta.prefix,
    meta.baseUrl,
    meta.typeLabel,
    meta.category,
  ];
  return fields.some(
    (value) => typeof value === "string" && value.toLowerCase().includes(q),
  );
}

export function buildProviderEntries(providerNodes = [], connections = [], nowMs = Date.now()) {
  const entries = [];

  for (const node of providerNodes || []) {
    if (node.type !== "openai-compatible" && node.type !== "anthropic-compatible") {
      continue;
    }
    const isAnthropic = node.type === "anthropic-compatible";
    const isOpenAI = node.type === "openai-compatible";
    const isMoonshot =
      isOpenAI &&
      (node.brand === "moonshot" || /moonshot|kimi/i.test(node.name || ""));
    const info = {
      id: node.id,
      name:
        node.name ||
        (isAnthropic ? "Anthropic Compatible" : "OpenAI Compatible"),
      brand: node.brand,
      color: isAnthropic ? "#D97757" : isMoonshot ? "#6366F1" : "#10A37F",
      textIcon: isAnthropic ? "AC" : "OC",
      apiType: node.apiType,
      logo: node.logo,
      type: node.type,
      prefix: node.prefix,
      baseUrl: node.baseUrl,
    };
    entries.push({
      id: node.id,
      info,
      category: "custom",
      section: "custom",
      authTypes: "apikey",
      cardAuthType: "compatible",
      typeLabel: getProviderTypeLabel(info, node.id, "compatible"),
      alias: node.prefix || "",
      prefix: node.prefix || "",
      baseUrl: node.baseUrl || "",
      stats: getProviderStats(connections, node.id, "apikey", nowMs),
      isNoAuth: false,
      isCompatible: isOpenAI,
      isAnthropicCompatible: isAnthropic,
    });
  }

  for (const [key, info] of Object.entries(OAUTH_PROVIDERS)) {
    if (info.hidden) continue;
    const authTypes = resolveProviderAuthTypes(info, key);
    entries.push({
      id: key,
      info,
      category: "oauth",
      section: "oauth",
      authTypes,
      cardAuthType: "oauth",
      typeLabel: getProviderTypeLabel(info, key, "oauth"),
      alias: info.alias || getProviderAlias(key) || "",
      prefix: "",
      baseUrl: "",
      stats: getProviderStats(connections, key, authTypes, nowMs),
      isNoAuth: Boolean(info.noAuth),
      isCompatible: false,
      isAnthropicCompatible: false,
    });
  }

  for (const [key, info] of Object.entries(FREE_PROVIDERS)) {
    if (info.hidden) continue;
    const authTypes = resolveProviderAuthTypes(info, key);
    entries.push({
      id: key,
      info,
      category: "free",
      section: "free",
      authTypes,
      cardAuthType: "free",
      typeLabel: getProviderTypeLabel(info, key, "free"),
      alias: info.alias || getProviderAlias(key) || "",
      prefix: "",
      baseUrl: "",
      stats: getProviderStats(connections, key, authTypes, nowMs),
      isNoAuth: Boolean(info.noAuth),
      isCompatible: false,
      isAnthropicCompatible: false,
    });
  }

  for (const [key, info] of Object.entries(FREE_TIER_PROVIDERS)) {
    if (info.hidden || !(info.serviceKinds ?? ["llm"]).includes("llm")) continue;
    const authTypes = resolveProviderAuthTypes(info, key);
    const cardAuthType = Array.isArray(authTypes)
      ? (authTypes[0] ?? "apikey")
      : authTypes;
    entries.push({
      id: key,
      info,
      category: "free",
      section: "freeTier",
      authTypes,
      cardAuthType,
      typeLabel: getProviderTypeLabel(info, key, cardAuthType),
      alias: info.alias || getProviderAlias(key) || "",
      prefix: "",
      baseUrl: "",
      stats: getProviderStats(connections, key, authTypes, nowMs),
      isNoAuth: Boolean(info.noAuth),
      isCompatible: false,
      isAnthropicCompatible: false,
    });
  }

  for (const [key, info] of Object.entries(WEB_COOKIE_PROVIDERS)) {
    if (info.hidden) continue;
    const authTypes = ["cookie", "apikey"];
    entries.push({
      id: key,
      info,
      category: "free",
      section: "webCookie",
      authTypes,
      cardAuthType: "cookie",
      typeLabel: getProviderTypeLabel(info, key, "cookie"),
      alias: info.alias || getProviderAlias(key) || "",
      prefix: "",
      baseUrl: "",
      stats: getProviderStats(connections, key, authTypes, nowMs),
      isNoAuth: Boolean(info.noAuth),
      isCompatible: false,
      isAnthropicCompatible: false,
    });
  }

  for (const [key, info] of Object.entries(APIKEY_PROVIDERS)) {
    if (info.hidden || !(info.serviceKinds ?? ["llm"]).includes("llm")) continue;
    entries.push({
      id: key,
      info,
      category: "apikey",
      section: "apikey",
      authTypes: "apikey",
      cardAuthType: "apikey",
      typeLabel: getProviderTypeLabel(info, key, "apikey"),
      alias: info.alias || getProviderAlias(key) || "",
      prefix: "",
      baseUrl: "",
      stats: getProviderStats(connections, key, "apikey", nowMs),
      isNoAuth: Boolean(info.noAuth),
      isCompatible: false,
      isAnthropicCompatible: false,
    });
  }

  return entries;
}

function getConfiguredRank(entry) {
  const s = entry.stats || {};
  if (s.error > 0 && !s.allDisabled) return 0; // Needs attention first
  if (s.connected > 0 && !s.allDisabled) return 1; // Active connected second
  if (!s.allDisabled) return 2; // Untested configured or 0-conn custom
  return 3; // Disabled last
}

function compareByPriorityAndName(a, b) {
  const pa = a.info?.priority ?? 999;
  const pb = b.info?.priority ?? 999;
  if (pa !== pb) return pa - pb;
  return (a.info?.name || "").localeCompare(b.info?.name || "");
}

export function sortConfiguredEntries(entries) {
  return [...entries].sort((a, b) => {
    const ra = getConfiguredRank(a);
    const rb = getConfiguredRank(b);
    if (ra !== rb) return ra - rb;
    if (ra === 0) {
      const errDiff =
        (b.stats?.latestErrorAtMs || 0) - (a.stats?.latestErrorAtMs || 0);
      if (errDiff !== 0) return errDiff;
    }
    if (ra === 1) return compareByPriorityAndName(a, b);
    return (a.info?.name || "").localeCompare(b.info?.name || "");
  });
}

export function sortCatalogSection(entries, section) {
  return [...entries].sort((a, b) => {
    if (section === "free") {
      const noAuthDiff = (b.isNoAuth ? 1 : 0) - (a.isNoAuth ? 1 : 0);
      if (noAuthDiff !== 0) return noAuthDiff;
    }
    return compareByPriorityAndName(a, b);
  });
}

export function partitionConfiguredAndCatalog(entries) {
  const configuredRaw = [];
  const catalogOauth = [];
  const catalogFree = [];
  const catalogApikey = [];

  for (const entry of entries || []) {
    const isConfigured = (entry.stats?.total || 0) > 0 || entry.category === "custom";
    if (isConfigured) {
      configuredRaw.push(entry);
    } else if (entry.category === "oauth") {
      catalogOauth.push(entry);
    } else if (entry.category === "free") {
      catalogFree.push(entry);
    } else if (entry.category === "apikey") {
      catalogApikey.push(entry);
    }
  }

  return {
    configured: sortConfiguredEntries(configuredRaw),
    catalog: {
      oauth: sortCatalogSection(catalogOauth, "oauth"),
      free: sortCatalogSection(catalogFree, "free"),
      apikey: sortCatalogSection(catalogApikey, "apikey"),
    },
  };
}

export function sortFilteredEntries(entries) {
  return [...(entries || [])].sort((a, b) => {
    const getRank = (e) => {
      const s = e.stats || {};
      if (!e.isNoAuth && s.error > 0 && !s.allDisabled) return 0;
      if (s.connected > 0 && !s.allDisabled) return 1;
      if (s.total > 0 && !s.allDisabled) return 2;
      if (e.isNoAuth) return 3;
      if (s.total > 0 && s.allDisabled) return 4;
      if (e.category === "custom") return 5;
      return 6;
    };
    const ra = getRank(a);
    const rb = getRank(b);
    if (ra !== rb) return ra - rb;
    if (ra === 0) {
      const errDiff =
        (b.stats?.latestErrorAtMs || 0) - (a.stats?.latestErrorAtMs || 0);
      if (errDiff !== 0) return errDiff;
    }
    return compareByPriorityAndName(a, b);
  });
}

export function computeFleetTelemetry(entries = [], connections = []) {
  let activeProviders = 0;
  let configuredActiveProviders = 0;
  let noAuthProviders = 0;
  let attentionProviders = 0;
  let disabledProviders = 0;

  for (const entry of entries) {
    if (matchesStatusFilter("active", entry.stats, entry.isNoAuth)) {
      activeProviders++;
      if (entry.isNoAuth && (entry.stats?.total || 0) === 0) {
        noAuthProviders++;
      } else {
        configuredActiveProviders++;
      }
    }
    if (matchesStatusFilter("error", entry.stats, entry.isNoAuth)) {
      attentionProviders++;
    }
    if (matchesStatusFilter("inactive", entry.stats, entry.isNoAuth)) {
      disabledProviders++;
    }
  }

  const totalConnections = (connections || []).length;
  const enabledConnections = (connections || []).filter(
    (c) => c.isActive !== false,
  ).length;

  return {
    activeProviders,
    configuredActiveProviders,
    noAuthProviders,
    enabledConnections,
    totalConnections,
    attentionProviders,
    disabledProviders,
  };
}

export function countStatusBuckets(entries = []) {
  const counts = { all: entries.length, active: 0, error: 0, inactive: 0, none: 0 };
  for (const entry of entries) {
    if (matchesStatusFilter("active", entry.stats, entry.isNoAuth)) counts.active++;
    if (matchesStatusFilter("error", entry.stats, entry.isNoAuth)) counts.error++;
    if (matchesStatusFilter("inactive", entry.stats, entry.isNoAuth)) counts.inactive++;
    if (matchesStatusFilter("none", entry.stats, entry.isNoAuth)) counts.none++;
  }
  return counts;
}

export function countCategoryBuckets(entries = []) {
  const counts = { all: entries.length, custom: 0, oauth: 0, apikey: 0, free: 0 };
  for (const entry of entries) {
    if (entry.category in counts) counts[entry.category]++;
  }
  return counts;
}

function getTestBatchAuthGroup(providerId, connection) {
  if (connection?.authType) {
    if (connection.authType === "oauth") {
      if (FREE_PROVIDERS[providerId]) return "free";
      return "oauth";
    }
    return connection.authType;
  }
  if (FREE_PROVIDERS[providerId]) return "free";
  if (OAUTH_PROVIDERS[providerId]) return "oauth";
  if (APIKEY_PROVIDERS[providerId]) return "apikey";
  if (
    typeof providerId === "string" &&
    (providerId.startsWith(OPENAI_COMPATIBLE_PREFIX) ||
      providerId.startsWith(ANTHROPIC_COMPATIBLE_PREFIX))
  ) {
    return "compatible";
  }
  return "apikey";
}

export function countActiveConnectionsByMode(connections = [], mode = "all") {
  const active = (connections || []).filter((c) => c.isActive !== false);
  if (mode === "all") return active.length;
  if (mode === "compatible") {
    return active.filter(
      (c) =>
        typeof c.provider === "string" &&
        (c.provider.startsWith(OPENAI_COMPATIBLE_PREFIX) ||
          c.provider.startsWith(ANTHROPIC_COMPATIBLE_PREFIX)),
    ).length;
  }
  return active.filter((c) => getTestBatchAuthGroup(c.provider, c) === mode)
    .length;
}
