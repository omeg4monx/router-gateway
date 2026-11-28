import { isUsageLogError, type UsageLog, type UsageLogsResponse, type UsageRange } from "./schemas";

export interface StreamLogsFilter {
  range?: UsageRange;
  q?: string;
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
  const search = filter.q?.trim().toLowerCase();
  if (search) {
    const matchesSearch =
      log.model.toLowerCase().includes(search) ||
      log.provider_id.toLowerCase().includes(search) ||
      log.request_id.toLowerCase().includes(search) ||
      log.error_summary.toLowerCase().includes(search);
    if (!matchesSearch) {
      return false;
    }
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
    if (isUsageLogError(log)) return false;
  } else if (status === "error") {
    if (!isUsageLogError(log)) return false;
  } else if (status && /^\d{3}$/.test(status)) {
    const code = Number(status);
    if (code >= 100 && code <= 599 && log.status_code !== code) {
      return false;
    }
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
