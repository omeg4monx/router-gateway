import type { UsageLog, UsageLogsResponse } from "./schemas";

export interface StreamLogsFilter {
  model?: string;
  provider?: string;
  key_id?: string;
  status?: string;
  limit?: number;
  offset?: number;
}

function matchesLogFilter(log: UsageLog, filter?: StreamLogsFilter): boolean {
  if (!filter) {
    return true;
  }
  if (typeof filter.offset === "number" && filter.offset > 0) {
    return false;
  }
  const model = filter.model?.trim();
  if (model && log.model !== model) {
    return false;
  }
  const provider = filter.provider?.trim();
  if (provider && log.provider_id !== provider) {
    return false;
  }
  const keyId = filter.key_id?.trim();
  if (keyId && log.key_id !== keyId) {
    return false;
  }
  const status = filter.status?.trim().toLowerCase();
  if (status === "ok" || status === "success") {
    const isOk = log.status_code >= 200 && log.status_code < 400 && !log.error_summary.trim();
    if (!isOk) return false;
  } else if (status === "error") {
    const isErr =
      log.status_code < 200 || log.status_code >= 400 || Boolean(log.error_summary.trim());
    if (!isErr) return false;
  }
  return true;
}

export function prependUsageLogToResponse(
  current: UsageLogsResponse,
  incoming: UsageLog,
  filter?: StreamLogsFilter,
): UsageLogsResponse {
  if (current.offset > 0 || !matchesLogFilter(incoming, filter)) {
    return current;
  }
  if (current.logs.some((item) => item.id === incoming.id)) {
    return current;
  }
  const maxItems =
    typeof filter?.limit === "number" && filter.limit > 0
      ? filter.limit
      : current.limit > 0
        ? current.limit
        : 50;
  return {
    ...current,
    logs: [incoming, ...current.logs].slice(0, maxItems),
    total: current.total + 1,
  };
}
