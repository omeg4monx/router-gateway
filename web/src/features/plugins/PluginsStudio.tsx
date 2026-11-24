import { Play, Puzzle, Save } from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "@/shared/ui/button";
import {
  type PluginState,
  usePluginsList,
  usePreviewPluginsMutation,
  useUpdatePluginMutation,
} from "./hooks";

interface DraftPlugin {
  enabled: boolean;
  modelsText: string;
  level: string;
}

const SAMPLE_PAYLOAD = JSON.stringify(
  {
    model: "gpt-4o",
    messages: [{ role: "user", content: "Summarize the router architecture in JSON." }],
  },
  null,
  2,
);

export function PluginsStudio() {
  const pluginsQuery = usePluginsList();
  const updateMutation = useUpdatePluginMutation();
  const previewMutation = usePreviewPluginsMutation();

  const [drafts, setDrafts] = useState<Record<string, DraftPlugin>>({});
  const [previewModel, setPreviewModel] = useState("gpt-4o");
  const [previewPayload, setPreviewPayload] = useState(SAMPLE_PAYLOAD);
  const [previewError, setPreviewError] = useState<string | null>(null);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);

  useEffect(() => {
    if (!pluginsQuery.data?.plugins) return;
    const next: Record<string, DraftPlugin> = {};
    for (const p of pluginsQuery.data.plugins) {
      const cfg = (p.config ?? {}) as Record<string, unknown>;
      next[p.key] = {
        enabled: p.enabled,
        modelsText: (p.models ?? []).join(", "),
        level: typeof cfg.level === "string" ? cfg.level : "full",
      };
    }
    setDrafts(next);
  }, [pluginsQuery.data]);

  const handleSavePlugin = (plugin: PluginState) => {
    setStatusMessage(null);
    const draft = drafts[plugin.key] ?? {
      enabled: plugin.enabled,
      modelsText: (plugin.models ?? []).join(", "),
      level: "full",
    };
    const models = draft.modelsText
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);
    const config: Record<string, unknown> = { ...(plugin.config ?? {}) };
    if (plugin.key === "antiSlop") {
      config.level = draft.level;
    }

    updateMutation.mutate(
      {
        key: plugin.key,
        body: {
          enabled: draft.enabled,
          models,
          config,
        },
      },
      {
        onSuccess: () => setStatusMessage(`Saved configuration for ${plugin.name}.`),
      },
    );
  };

  const handlePreview = () => {
    setPreviewError(null);
    let parsed: Record<string, unknown>;
    try {
      parsed = JSON.parse(previewPayload) as Record<string, unknown>;
    } catch {
      setPreviewError("Preview payload must be valid JSON.");
      return;
    }
    previewMutation.mutate({
      model: previewModel.trim() || "gpt-4o",
      payload: parsed,
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-2 border-b border-border pb-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <div className="flex items-center gap-2 text-primary">
            <Puzzle className="h-4 w-4" />
            <span className="font-mono text-xs uppercase tracking-widest">
              Gateway Interceptors
            </span>
          </div>
          <h1 className="mt-1 text-2xl font-semibold tracking-tight text-foreground">
            Request Plugins
          </h1>
          <p className="text-sm text-muted-foreground">
            Configure model-scoped request enhancers executed inside the Go gateway pipeline.
          </p>
        </div>
      </div>

      {statusMessage && (
        <div className="rounded-md border border-primary/40 bg-card px-4 py-2.5 text-sm text-foreground">
          {statusMessage}
        </div>
      )}

      {updateMutation.error && (
        <div className="rounded-md border border-destructive/50 bg-card px-4 py-2.5 text-sm text-destructive">
          {updateMutation.error.message}
        </div>
      )}

      {pluginsQuery.isLoading ? (
        <div className="rounded-lg border border-border bg-card p-6 text-sm text-muted-foreground">
          Loading plugins catalog...
        </div>
      ) : pluginsQuery.error ? (
        <div className="rounded-lg border border-destructive/50 bg-card p-6 text-sm text-destructive">
          {pluginsQuery.error.message}
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2">
          {(pluginsQuery.data?.plugins ?? []).map((plugin) => {
            const draft = drafts[plugin.key] ?? {
              enabled: plugin.enabled,
              modelsText: (plugin.models ?? []).join(", "),
              level: "full",
            };
            return (
              <div
                key={plugin.key}
                className="flex flex-col justify-between rounded-lg border border-border bg-card p-5"
              >
                <div className="space-y-3">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <div className="flex items-center gap-2">
                        <h2 className="text-base font-semibold text-foreground">{plugin.name}</h2>
                        <span className="rounded border border-border bg-background px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
                          {plugin.category}
                        </span>
                      </div>
                      <p className="mt-1 text-xs text-muted-foreground">{plugin.description}</p>
                    </div>
                    <label className="flex items-center gap-2 text-xs font-medium text-foreground">
                      <input
                        type="checkbox"
                        checked={draft.enabled}
                        onChange={(e) =>
                          setDrafts((prev) => ({
                            ...prev,
                            [plugin.key]: { ...draft, enabled: e.target.checked },
                          }))
                        }
                        className="h-4 w-4 accent-primary"
                      />
                      Enabled
                    </label>
                  </div>

                  <div className="space-y-2">
                    <label
                      htmlFor={`plugin-models-${plugin.key}`}
                      className="block text-xs font-medium text-muted-foreground"
                    >
                      Model Allowlist (comma-separated, blank = all models)
                    </label>
                    <input
                      id={`plugin-models-${plugin.key}`}
                      type="text"
                      value={draft.modelsText}
                      onChange={(e) =>
                        setDrafts((prev) => ({
                          ...prev,
                          [plugin.key]: { ...draft, modelsText: e.target.value },
                        }))
                      }
                      placeholder="gpt-4o, claude-*"
                      className="w-full rounded-md border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
                    />
                  </div>

                  {plugin.key === "antiSlop" && (
                    <div className="space-y-1.5">
                      <label
                        htmlFor={`plugin-level-${plugin.key}`}
                        className="block text-xs font-medium text-muted-foreground"
                      >
                        Strictness Level
                      </label>
                      <select
                        id={`plugin-level-${plugin.key}`}
                        value={draft.level}
                        onChange={(e) =>
                          setDrafts((prev) => ({
                            ...prev,
                            [plugin.key]: { ...draft, level: e.target.value },
                          }))
                        }
                        className="w-full rounded-md border border-border bg-background px-3 py-1.5 text-xs text-foreground"
                      >
                        <option value="light">Light</option>
                        <option value="standard">Standard</option>
                        <option value="full">Full</option>
                      </select>
                    </div>
                  )}
                </div>

                <div className="mt-4 flex items-center justify-between border-t border-border pt-3">
                  <span className="font-mono text-[11px] text-muted-foreground">{plugin.key}</span>
                  <Button
                    type="button"
                    size="sm"
                    disabled={updateMutation.isPending}
                    onClick={() => handleSavePlugin(plugin)}
                  >
                    <Save className="mr-1.5 h-3.5 w-3.5" />
                    Save
                  </Button>
                </div>
              </div>
            );
          })}
        </div>
      )}

      <div className="rounded-lg border border-border bg-card p-5 space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-base font-semibold text-foreground">Dry-Run Interceptor Preview</h2>
            <p className="text-xs text-muted-foreground">
              Simulate active plugins against a sample OpenAI/Anthropic request payload.
            </p>
          </div>
          <Button
            type="button"
            size="sm"
            onClick={handlePreview}
            disabled={previewMutation.isPending}
          >
            <Play className="mr-1.5 h-3.5 w-3.5" />
            Run Preview
          </Button>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <label
              htmlFor="plugins-preview-model"
              className="block text-xs font-medium text-muted-foreground"
            >
              Target Model
            </label>
            <input
              id="plugins-preview-model"
              type="text"
              value={previewModel}
              onChange={(e) => setPreviewModel(e.target.value)}
              className="w-full rounded-md border border-border bg-background px-3 py-1.5 font-mono text-xs text-foreground"
            />
            <label
              htmlFor="plugins-preview-payload"
              className="block text-xs font-medium text-muted-foreground"
            >
              Input JSON Payload
            </label>
            <textarea
              id="plugins-preview-payload"
              rows={8}
              value={previewPayload}
              onChange={(e) => setPreviewPayload(e.target.value)}
              className="w-full rounded-md border border-border bg-background p-3 font-mono text-xs text-foreground"
            />
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-muted-foreground">Transformed Output</span>
              {previewMutation.data && (
                <span className="font-mono text-xs text-primary">
                  Applied:{" "}
                  {previewMutation.data.applied_plugins.length > 0
                    ? previewMutation.data.applied_plugins.join(", ")
                    : "none"}
                </span>
              )}
            </div>
            {(previewError || previewMutation.error) && (
              <div className="rounded border border-destructive/50 bg-background p-2.5 text-xs text-destructive">
                {previewError ?? previewMutation.error?.message}
              </div>
            )}
            <pre className="h-[234px] overflow-auto rounded-md border border-border bg-background p-3 font-mono text-xs text-foreground">
              {previewMutation.data
                ? JSON.stringify(previewMutation.data.payload, null, 2)
                : "Run a preview to inspect the transformed request body."}
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}
