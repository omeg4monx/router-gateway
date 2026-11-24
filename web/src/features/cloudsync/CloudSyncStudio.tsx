import { Cloud, Download, RefreshCw, Save, Trash2, Upload } from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "@/shared/ui/button";
import {
  type CloudSyncSnapshotBundle,
  useClearCloudSyncHistoryMutation,
  useCloudSyncConfig,
  useCloudSyncHistory,
  useExportSnapshotMutation,
  usePullCloudSyncMutation,
  usePushCloudSyncMutation,
  useRestoreSnapshotMutation,
  useUpdateCloudSyncConfigMutation,
} from "./hooks";

export function CloudSyncStudio() {
  const configQuery = useCloudSyncConfig();
  const historyQuery = useCloudSyncHistory();
  const updateMutation = useUpdateCloudSyncConfigMutation();
  const pushMutation = usePushCloudSyncMutation();
  const pullMutation = usePullCloudSyncMutation();
  const exportMutation = useExportSnapshotMutation();
  const restoreMutation = useRestoreSnapshotMutation();
  const clearHistoryMutation = useClearCloudSyncHistoryMutation();

  const [enabled, setEnabled] = useState(false);
  const [endpointURL, setEndpointURL] = useState("");
  const [authToken, setAuthToken] = useState("");
  const [intervalMin, setIntervalMin] = useState(30);
  const [conflictPolicy, setConflictPolicy] = useState<"merge" | "remote_wins" | "local_wins">(
    "merge",
  );
  const [includeSettings, setIncludeSettings] = useState(true);
  const [includeCombos, setIncludeCombos] = useState(true);
  const [snapshotText, setSnapshotText] = useState("");
  const [notice, setNotice] = useState<string | null>(null);
  const [localError, setLocalError] = useState<string | null>(null);

  useEffect(() => {
    const cfg = configQuery.data;
    if (!cfg) return;
    setEnabled(cfg.enabled);
    setEndpointURL(cfg.endpoint_url);
    setIntervalMin(cfg.sync_interval_min);
    setConflictPolicy(cfg.conflict_policy);
    setIncludeSettings(cfg.include_settings);
    setIncludeCombos(cfg.include_combos);
  }, [configQuery.data]);

  const handleSaveConfig = (e: React.FormEvent) => {
    e.preventDefault();
    setNotice(null);
    setLocalError(null);
    updateMutation.mutate(
      {
        enabled,
        endpoint_url: endpointURL.trim(),
        ...(authToken.trim() ? { auth_token: authToken.trim() } : {}),
        sync_interval_min: intervalMin,
        conflict_policy: conflictPolicy,
        include_settings: includeSettings,
        include_combos: includeCombos,
      },
      {
        onSuccess: () => {
          setAuthToken("");
          setNotice("Cloud Sync configuration saved.");
        },
      },
    );
  };

  const handleExport = () => {
    setNotice(null);
    setLocalError(null);
    exportMutation.mutate(undefined, {
      onSuccess: (bundle) => {
        setSnapshotText(JSON.stringify(bundle, null, 2));
        setNotice(`Exported snapshot (${bundle.checksum}).`);
      },
    });
  };

  const handleRestore = () => {
    setNotice(null);
    setLocalError(null);
    let parsed: CloudSyncSnapshotBundle;
    try {
      parsed = JSON.parse(snapshotText) as CloudSyncSnapshotBundle;
    } catch {
      setLocalError("Snapshot JSON is invalid.");
      return;
    }
    restoreMutation.mutate(
      { conflict_policy: conflictPolicy, snapshot: parsed },
      {
        onSuccess: (res) => setNotice(`Restored ${res.item_count} items (${res.checksum}).`),
      },
    );
  };

  const activeError =
    localError ??
    updateMutation.error?.message ??
    pushMutation.error?.message ??
    pullMutation.error?.message ??
    exportMutation.error?.message ??
    restoreMutation.error?.message;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 border-b border-border pb-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="flex items-center gap-2 text-primary">
            <Cloud className="h-4 w-4" />
            <span className="font-mono text-xs uppercase tracking-widest">Remote State Sync</span>
          </div>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight text-foreground">Cloud Sync</h1>
          <p className="text-sm text-muted-foreground">
            Synchronize settings, combos, and aliases with a remote HTTPS endpoint or portable JSON
            snapshots.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button
            type="button"
            size="sm"
            variant="outline"
            disabled={pushMutation.isPending}
            onClick={() =>
              pushMutation.mutate(undefined, {
                onSuccess: (r) => setNotice(`Pushed ${r.item_count} items (${r.checksum}).`),
              })
            }
          >
            <Upload className="mr-1.5 h-3.5 w-3.5" />
            Push Remote
          </Button>
          <Button
            type="button"
            size="sm"
            disabled={pullMutation.isPending}
            onClick={() =>
              pullMutation.mutate(undefined, {
                onSuccess: (r) => setNotice(`Pulled ${r.item_count} items (${r.checksum}).`),
              })
            }
          >
            <RefreshCw className="mr-1.5 h-3.5 w-3.5" />
            Pull Remote
          </Button>
        </div>
      </div>

      {notice && (
        <div className="rounded-md border border-primary/40 bg-card px-4 py-2.5 text-sm text-foreground">
          {notice}
        </div>
      )}
      {activeError && (
        <div className="rounded-md border border-destructive/50 bg-card px-4 py-2.5 text-sm text-destructive">
          {activeError}
        </div>
      )}

      <div className="grid gap-6 lg:grid-cols-2">
        <form
          onSubmit={handleSaveConfig}
          className="space-y-4 rounded-lg border border-border bg-card p-5"
        >
          <div className="flex items-center justify-between">
            <h2 className="text-base font-semibold text-foreground">Endpoint & Policy</h2>
            <label className="flex items-center gap-2 text-xs font-medium text-foreground">
              <input
                type="checkbox"
                checked={enabled}
                onChange={(e) => setEnabled(e.target.checked)}
                className="h-4 w-4 accent-primary"
              />
              Enabled
            </label>
          </div>

          <div className="space-y-1.5">
            <label
              htmlFor="cloudsync-endpoint"
              className="block text-xs font-medium text-muted-foreground"
            >
              Remote HTTPS Endpoint URL
            </label>
            <input
              id="cloudsync-endpoint"
              type="url"
              value={endpointURL}
              onChange={(e) => setEndpointURL(e.target.value)}
              placeholder="https://sync.example.com/v1/router-state"
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
            />
          </div>

          <div className="space-y-1.5">
            <label
              htmlFor="cloudsync-token"
              className="block text-xs font-medium text-muted-foreground"
            >
              Vaulted Bearer Token{" "}
              {configQuery.data?.has_auth_token
                ? `(stored: ${configQuery.data.auth_token_mask})`
                : "(not set)"}
            </label>
            <input
              id="cloudsync-token"
              type="password"
              value={authToken}
              onChange={(e) => setAuthToken(e.target.value)}
              placeholder="Enter new token to encrypt in vault"
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
            />
          </div>

          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5">
              <label
                htmlFor="cloudsync-policy"
                className="block text-xs font-medium text-muted-foreground"
              >
                Conflict Policy
              </label>
              <select
                id="cloudsync-policy"
                value={conflictPolicy}
                onChange={(e) =>
                  setConflictPolicy(e.target.value as "merge" | "remote_wins" | "local_wins")
                }
                className="w-full rounded-md border border-border bg-background px-3 py-1.5 text-xs text-foreground"
              >
                <option value="merge">Merge (prefer additions)</option>
                <option value="remote_wins">Remote Wins</option>
                <option value="local_wins">Local Wins</option>
              </select>
            </div>
            <div className="space-y-1.5">
              <label
                htmlFor="cloudsync-interval"
                className="block text-xs font-medium text-muted-foreground"
              >
                Interval (minutes)
              </label>
              <input
                id="cloudsync-interval"
                type="number"
                min={1}
                max={1440}
                value={intervalMin}
                onChange={(e) => setIntervalMin(Number(e.target.value) || 30)}
                className="w-full rounded-md border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
              />
            </div>
          </div>

          <div className="flex flex-wrap gap-4 pt-1 text-xs text-foreground">
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={includeSettings}
                onChange={(e) => setIncludeSettings(e.target.checked)}
                className="h-4 w-4 accent-primary"
              />
              Include Settings
            </label>
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={includeCombos}
                onChange={(e) => setIncludeCombos(e.target.checked)}
                className="h-4 w-4 accent-primary"
              />
              Include Combos & Aliases
            </label>
          </div>

          <Button type="submit" size="sm" disabled={updateMutation.isPending}>
            <Save className="mr-1.5 h-3.5 w-3.5" />
            Save Configuration
          </Button>
        </form>

        <div className="space-y-3 rounded-lg border border-border bg-card p-5">
          <div className="flex items-center justify-between">
            <h2 className="text-base font-semibold text-foreground">Portable JSON Snapshot</h2>
            <div className="flex items-center gap-2">
              <Button
                type="button"
                size="sm"
                variant="outline"
                disabled={exportMutation.isPending}
                onClick={handleExport}
              >
                <Download className="mr-1.5 h-3.5 w-3.5" />
                Export
              </Button>
              <Button
                type="button"
                size="sm"
                disabled={restoreMutation.isPending || !snapshotText.trim()}
                onClick={handleRestore}
              >
                Restore
              </Button>
            </div>
          </div>
          <textarea
            rows={12}
            value={snapshotText}
            onChange={(e) => setSnapshotText(e.target.value)}
            placeholder="Export a snapshot or paste a JSON snapshot bundle to restore..."
            className="w-full rounded-md border border-border bg-background p-3 font-mono text-xs text-foreground"
          />
        </div>
      </div>

      <div className="rounded-lg border border-border bg-card p-5 space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-semibold text-foreground">Sync History</h2>
          <Button
            type="button"
            size="sm"
            variant="outline"
            disabled={clearHistoryMutation.isPending}
            onClick={() => clearHistoryMutation.mutate()}
          >
            <Trash2 className="mr-1.5 h-3.5 w-3.5" />
            Clear
          </Button>
        </div>
        {(historyQuery.data?.events ?? []).length === 0 ? (
          <p className="text-xs text-muted-foreground">No sync events recorded yet.</p>
        ) : (
          <div className="divide-y divide-border">
            {(historyQuery.data?.events ?? []).map((ev) => (
              <div
                key={ev.id}
                className="flex flex-wrap items-center justify-between gap-2 py-2.5 text-xs"
              >
                <div className="flex items-center gap-2">
                  <span className="rounded border border-border bg-background px-2 py-0.5 font-mono uppercase text-foreground">
                    {ev.direction}
                  </span>
                  <span
                    className={
                      ev.status === "ok"
                        ? "font-medium text-primary"
                        : "font-medium text-destructive"
                    }
                  >
                    {ev.status}
                  </span>
                  <span className="text-foreground">{ev.summary}</span>
                </div>
                <span className="font-mono text-muted-foreground">{ev.created_at}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
