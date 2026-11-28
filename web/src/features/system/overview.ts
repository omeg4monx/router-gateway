export function formatClockUptime(totalSeconds: number): string {
  const safe = Math.max(0, Math.floor(totalSeconds));
  const days = Math.floor(safe / 86_400);
  const hours = Math.floor((safe % 86_400) / 3_600);
  const minutes = Math.floor((safe % 3_600) / 60);
  const seconds = safe % 60;

  const pad = (n: number) => String(n).padStart(2, "0");
  if (days > 0) {
    return `${days}d ${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`;
  }
  return `${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`;
}

export function formatActivityOriginLabel(
  keyId: string,
  keyNameMap: ReadonlyMap<string, string>,
  isId: boolean,
): string {
  if (keyId === "admin-internal") {
    return isId ? "Studio Admin" : "Admin Studio";
  }
  if (keyId === "admin-probe") {
    return isId ? "Uji Model" : "Model Probe";
  }
  return keyNameMap.get(keyId) ?? keyId;
}
