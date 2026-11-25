import { Database, Play, Search, Trash2 } from "lucide-react";
import { useState } from "react";
import { Button } from "@/shared/ui/button";
import {
  useClearImporterHistoryMutation,
  useImporterHistory,
  useInspectLegacyMutation,
  useRunLegacyImportMutation,
} from "./hooks";

export function ImporterStudio() {
  const historyQuery = useImporterHistory();
  const inspectMutation = useInspectLegacyMutation();
  const runMutation = useRunLegacyImportMutation();
  const clearMutation = useClearImporterHistoryMutation();

  const [sourcePath, setSourcePath] = useState("");
  const [importSettings, setImportSettings] = useState(true);
  const [importCombos, setImportCombos] = useState(true);
  const [importKeys, setImportKeys] = useState(false);
  const [overwrite, setOverwrite] = useState(false);

  const handleInspect = (e: React.FormEvent) => {
    e.preventDefault();
    inspectMutation.mutate({ source_path: sourcePath.trim() });
  };

  const handleRun = (dryRun: boolean) => {
    runMutation.mutate({
      source_path: sourcePath.trim(),
      dry_run: dryRun,
      import_settings: importSettings,
      import_combos: importCombos,
      import_keys: importKeys,
      overwrite,
    });
  };

  const errMsg = inspectMutation.error?.message ?? runMutation.error?.message;

  return (
    <div className="space-y-6">
      <div className="border-b border-border pb-4">
        <div className="flex items-center gap-2 text-primary">
          <Database className="h-4 w-4" />
          <span className="font-mono text-xs uppercase tracking-widest">One-Way Migration</span>
        </div>
        <h1 className="mt-1 text-2xl font-semibold tracking-tight text-foreground">
          Legacy Importer
        </h1>
        <p className="text-sm text-muted-foreground">
          Inspect a legacy <code className="font-mono">data.sqlite</code> strictly read-only and
          migrate settings, fallback combos, and API keys cleanly.
        </p>
      </div>

      {errMsg && (
        <div className="rounded-md border border-destructive/50 bg-card px-4 py-2.5 text-sm text-destructive">
          {errMsg}
        </div>
      )}

      <form
        onSubmit={handleInspect}
        className="space-y-4 rounded-lg border border-border bg-card p-5"
      >
        <div className="space-y-1.5">
          <label
            htmlFor="importer-source-path"
            className="block text-xs font-medium text-muted-foreground"
          >
            Legacy SQLite File Path
          </label>
          <div className="flex flex-col gap-2 sm:flex-row">
            <input
              id="importer-source-path"
              type="text"
              required
              value={sourcePath}
              onChange={(e) => setSourcePath(e.target.value)}
              placeholder="/path/to/legacy/data.sqlite"
              className="flex-1 rounded-md border border-border bg-background px-3 py-2 font-mono text-xs text-foreground"
            />
            <Button type="submit" size="sm" variant="outline" disabled={inspectMutation.isPending}>
              <Search className="mr-1.5 h-3.5 w-3.5" />
              Inspect Read-Only
            </Button>
          </div>
        </div>

        <div className="flex flex-wrap gap-4 text-xs text-foreground">
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={importSettings}
              onChange={(e) => setImportSettings(e.target.checked)}
              className="h-4 w-4 accent-primary"
            />
            Import Settings
          </label>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={importCombos}
              onChange={(e) => setImportCombos(e.target.checked)}
              className="h-4 w-4 accent-primary"
            />
            Import Combos
          </label>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={importKeys}
              onChange={(e) => setImportKeys(e.target.checked)}
              className="h-4 w-4 accent-primary"
            />
            Mint New Gateway Keys
          </label>
          <label className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={overwrite}
              onChange={(e) => setOverwrite(e.target.checked)}
              className="h-4 w-4 accent-primary"
            />
            Overwrite Existing Combos
          </label>
        </div>

        <div className="flex flex-wrap gap-2">
          <Button
            type="button"
            size="sm"
            variant="outline"
            disabled={!sourcePath.trim() || runMutation.isPending}
            onClick={() => handleRun(true)}
          >
            Dry-Run Preview
          </Button>
          <Button
            type="button"
            size="sm"
            disabled={!sourcePath.trim() || runMutation.isPending}
            onClick={() => handleRun(false)}
          >
            <Play className="mr-1.5 h-3.5 w-3.5" />
            Run Live Import
          </Button>
        </div>
      </form>

      {inspectMutation.data && (
        <div className="rounded-lg border border-border bg-card p-5 space-y-3">
          <h2 className="text-base font-semibold text-foreground">Inspection Summary</h2>
          <div className="grid gap-3 sm:grid-cols-4">
            <div className="rounded border border-border bg-background p-3">
              <div className="text-xs text-muted-foreground">Settings</div>
              <div className="mt-1 font-mono text-xl font-semibold text-foreground">
                {inspectMutation.data.settings_count}
              </div>
            </div>
            <div className="rounded border border-border bg-background p-3">
              <div className="text-xs text-muted-foreground">Combos</div>
              <div className="mt-1 font-mono text-xl font-semibold text-foreground">
                {inspectMutation.data.combos_count}
              </div>
            </div>
            <div className="rounded border border-border bg-background p-3">
              <div className="text-xs text-muted-foreground">API Keys</div>
              <div className="mt-1 font-mono text-xl font-semibold text-foreground">
                {inspectMutation.data.keys_count}
              </div>
            </div>
            <div className="rounded border border-border bg-background p-3">
              <div className="text-xs text-muted-foreground">Legacy Providers</div>
              <div className="mt-1 font-mono text-xl font-semibold text-foreground">
                {inspectMutation.data.providers_count}
              </div>
            </div>
          </div>
          {inspectMutation.data.warnings.length > 0 && (
            <ul className="list-disc space-y-1 pl-5 text-xs text-muted-foreground">
              {inspectMutation.data.warnings.map((w) => (
                <li key={w}>{w}</li>
              ))}
            </ul>
          )}
        </div>
      )}

      {runMutation.data && (
        <div className="rounded-lg border border-primary/40 bg-card p-5 space-y-2">
          <h2 className="text-base font-semibold text-foreground">
            {runMutation.data.run.dry_run ? "Dry-Run Result" : "Import Completed"}
          </h2>
          <p className="text-xs text-foreground">{runMutation.data.run.summary}</p>
          {runMutation.data.provisioned_keys.length > 0 && (
            <div className="mt-2 space-y-1">
              <div className="text-xs font-medium text-primary">
                Newly Minted Gateway API Keys (copy once):
              </div>
              {runMutation.data.provisioned_keys.map((k) => (
                <div
                  key={k.ID}
                  className="rounded border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
                >
                  {k.Name}: {k.Secret}
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      <div className="rounded-lg border border-border bg-card p-5 space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="text-base font-semibold text-foreground">Import History</h2>
          <Button
            type="button"
            size="sm"
            variant="outline"
            disabled={clearMutation.isPending}
            onClick={() => clearMutation.mutate()}
          >
            <Trash2 className="mr-1.5 h-3.5 w-3.5" />
            Clear
          </Button>
        </div>
        {(historyQuery.data?.runs ?? []).length === 0 ? (
          <p className="text-xs text-muted-foreground">No import runs recorded yet.</p>
        ) : (
          <div className="divide-y divide-border">
            {(historyQuery.data?.runs ?? []).map((run) => (
              <div
                key={run.id}
                className="flex flex-wrap items-center justify-between gap-2 py-2.5 text-xs"
              >
                <div className="flex items-center gap-2">
                  <span className="rounded border border-border bg-background px-2 py-0.5 font-mono text-foreground">
                    {run.dry_run ? "DRY-RUN" : "LIVE"}
                  </span>
                  <span className="text-foreground">{run.summary}</span>
                </div>
                <span className="font-mono text-muted-foreground">{run.created_at}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
