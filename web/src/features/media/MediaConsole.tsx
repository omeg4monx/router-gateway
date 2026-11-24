import { zodResolver } from "@hookform/resolvers/zod";
import { Globe, Image, Layers, Plus, RefreshCw, Trash2 } from "lucide-react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useCreateMediaComboMutation,
  useCreateMediaProviderMutation,
  useDeleteMediaComboMutation,
  useDeleteMediaProviderMutation,
  useMediaCombos,
  useMediaProviders,
  useUpdateMediaComboMutation,
  useUpdateMediaProviderMutation,
} from "./hooks";
import {
  type CreateMediaComboFormInput,
  type CreateMediaComboFormOutput,
  type CreateMediaProviderFormInput,
  type CreateMediaProviderFormOutput,
  createMediaComboSchema,
  createMediaProviderSchema,
  MEDIA_CAPABILITIES,
  MEDIA_DRIVERS,
  parseTargetsList,
} from "./schemas";

const GATEWAY_ENDPOINTS = [
  {
    capability: "audio",
    method: "POST",
    path: "/v1/audio/speech, /v1/audio/transcriptions, /v1/audio/translations",
  },
  {
    capability: "image",
    method: "POST",
    path: "/v1/images/generations, /v1/images/edits, /v1/images/variations",
  },
  { capability: "video", method: "POST/GET", path: "/v1/videos/generations, /v1/videos/{id}" },
  { capability: "search", method: "POST", path: "/v1/search" },
  { capability: "fetch", method: "POST", path: "/v1/web/fetch (built-in SSRF-guarded extractor)" },
];

export function MediaConsole() {
  const providersQuery = useMediaProviders();
  const combosQuery = useMediaCombos();
  const createProvider = useCreateMediaProviderMutation();
  const updateProvider = useUpdateMediaProviderMutation();
  const deleteProvider = useDeleteMediaProviderMutation();
  const createCombo = useCreateMediaComboMutation();
  const updateCombo = useUpdateMediaComboMutation();
  const deleteCombo = useDeleteMediaComboMutation();

  const providerForm = useForm<
    CreateMediaProviderFormInput,
    unknown,
    CreateMediaProviderFormOutput
  >({
    resolver: zodResolver(createMediaProviderSchema),
    defaultValues: {
      name: "",
      capability: "image",
      driver: "openai_compatible",
      base_url: "",
      api_key: "",
      default_model: "",
      priority: 100,
      enabled: true,
    },
  });

  const comboForm = useForm<CreateMediaComboFormInput, unknown, CreateMediaComboFormOutput>({
    resolver: zodResolver(createMediaComboSchema),
    defaultValues: {
      name: "",
      capability: "image",
      targets_text: "",
      enabled: true,
    },
  });

  const providers = providersQuery.data?.providers ?? [];
  const combos = combosQuery.data?.combos ?? [];

  const activeError =
    providersQuery.error ??
    combosQuery.error ??
    createProvider.error ??
    updateProvider.error ??
    deleteProvider.error ??
    createCombo.error ??
    updateCombo.error ??
    deleteCombo.error;

  const clearMutationErrors = () => {
    createProvider.reset();
    updateProvider.reset();
    deleteProvider.reset();
    createCombo.reset();
    updateCombo.reset();
    deleteCombo.reset();
  };

  const onCreateProvider = providerForm.handleSubmit((values) => {
    clearMutationErrors();
    createProvider.mutate(
      {
        name: values.name,
        capability: values.capability,
        driver: values.driver,
        base_url: values.base_url || undefined,
        api_key: values.api_key || undefined,
        default_model: values.default_model || undefined,
        priority: values.priority,
        enabled: values.enabled,
      },
      {
        onSuccess: () => {
          providerForm.reset();
        },
      },
    );
  });

  const onCreateCombo = comboForm.handleSubmit((values) => {
    clearMutationErrors();
    createCombo.mutate(
      {
        name: values.name,
        capability: values.capability,
        targets: parseTargetsList(values.targets_text),
        enabled: values.enabled,
      },
      {
        onSuccess: () => {
          comboForm.reset();
        },
      },
    );
  });

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Wave 3B — Media Gateway
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Media Providers &amp; Fallback Combos
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Configure encrypted upstream providers and ordered fallback chains for audio, images,
            videos, web search, and SSRF-guarded web fetch.
          </p>
        </div>
        <Button
          type="button"
          variant="secondary"
          size="sm"
          onClick={() => {
            void Promise.all([providersQuery.refetch(), combosQuery.refetch()]);
          }}
          disabled={providersQuery.isFetching || combosQuery.isFetching}
        >
          <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
          <span>Refresh</span>
        </Button>
      </section>

      <EmberRule />

      {activeError ? (
        <Alert variant="destructive">
          <AlertTitle>Media Gateway Error</AlertTitle>
          <AlertDescription>{activeError.message}</AlertDescription>
        </Alert>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-base">
            <Globe className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <span>Gateway Media Endpoints</span>
          </CardTitle>
          <CardDescription>
            Authenticated with standard gateway API keys (Bearer / x-api-key).
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
            {GATEWAY_ENDPOINTS.map((ep) => (
              <div
                key={ep.capability}
                className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3"
              >
                <div className="flex items-center justify-between gap-2">
                  <Badge variant="outline" className="font-mono uppercase">
                    {ep.capability}
                  </Badge>
                  <span className="font-mono text-[11px] text-[var(--text-muted)]">
                    {ep.method}
                  </span>
                </div>
                <p className="mt-2 break-all font-mono text-xs text-[var(--text-secondary)]">
                  {ep.path}
                </p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Image className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <span>Add Media Provider</span>
            </CardTitle>
            <CardDescription>
              API keys are encrypted at rest with AES-256-GCM and never returned after save.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={onCreateProvider} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="media-prov-name">Name</Label>
                  <Input
                    id="media-prov-name"
                    placeholder="openai-images"
                    {...providerForm.register("name")}
                  />
                  {providerForm.formState.errors.name ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {providerForm.formState.errors.name.message}
                    </p>
                  ) : null}
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="media-prov-cap">Capability</Label>
                  <select
                    id="media-prov-cap"
                    className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
                    {...providerForm.register("capability")}
                  >
                    {MEDIA_CAPABILITIES.map((cap) => (
                      <option key={cap} value={cap}>
                        {cap}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="media-prov-driver">Driver</Label>
                  <select
                    id="media-prov-driver"
                    className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
                    {...providerForm.register("driver")}
                  >
                    {MEDIA_DRIVERS.map((drv) => (
                      <option key={drv} value={drv}>
                        {drv}
                      </option>
                    ))}
                  </select>
                  {providerForm.formState.errors.driver ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {providerForm.formState.errors.driver.message}
                    </p>
                  ) : null}
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="media-prov-prio">Priority (0-10000)</Label>
                  <Input
                    id="media-prov-prio"
                    type="number"
                    {...providerForm.register("priority")}
                  />
                  {providerForm.formState.errors.priority ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {providerForm.formState.errors.priority.message}
                    </p>
                  ) : null}
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="media-prov-url">Base URL</Label>
                <Input
                  id="media-prov-url"
                  placeholder="https://api.openai.com"
                  {...providerForm.register("base_url")}
                />
                {providerForm.formState.errors.base_url ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {providerForm.formState.errors.base_url.message}
                  </p>
                ) : null}
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="media-prov-key">API Key (optional)</Label>
                  <Input
                    id="media-prov-key"
                    type="password"
                    placeholder="sk-..."
                    {...providerForm.register("api_key")}
                  />
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="media-prov-model">Default Model (optional)</Label>
                  <Input
                    id="media-prov-model"
                    placeholder="gpt-image-1"
                    {...providerForm.register("default_model")}
                  />
                </div>
              </div>

              <Button type="submit" size="sm" disabled={createProvider.isPending}>
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{createProvider.isPending ? "Saving..." : "Add Provider"}</span>
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Layers className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <span>Create Media Combo</span>
            </CardTitle>
            <CardDescription>
              Ordered provider failover chains triggered on upstream 401/403/429/5xx errors.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={onCreateCombo} className="space-y-4">
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="media-combo-name">Combo Name</Label>
                  <Input
                    id="media-combo-name"
                    placeholder="images-primary"
                    {...comboForm.register("name")}
                  />
                  {comboForm.formState.errors.name ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {comboForm.formState.errors.name.message}
                    </p>
                  ) : null}
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="media-combo-cap">Capability</Label>
                  <select
                    id="media-combo-cap"
                    className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
                    {...comboForm.register("capability")}
                  >
                    {MEDIA_CAPABILITIES.map((cap) => (
                      <option key={cap} value={cap}>
                        {cap}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="media-combo-targets">
                  Ordered Targets (comma or newline separated provider names or IDs)
                </Label>
                <Input
                  id="media-combo-targets"
                  placeholder="openai-images, backup-images"
                  {...comboForm.register("targets_text")}
                />
                {comboForm.formState.errors.targets_text ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {comboForm.formState.errors.targets_text.message}
                  </p>
                ) : null}
              </div>

              <Button type="submit" size="sm" disabled={createCombo.isPending}>
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{createCombo.isPending ? "Creating..." : "Create Combo"}</span>
              </Button>
            </form>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            Configured Media Providers ({providers.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          {providers.length === 0 ? (
            <p className="text-sm text-[var(--text-muted)]">No media providers configured yet.</p>
          ) : (
            <div className="space-y-2.5">
              {providers.map((p) => (
                <div
                  key={p.id}
                  className="flex flex-col gap-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="min-w-0 space-y-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-medium text-[var(--text-primary)]">{p.name}</span>
                      <Badge variant="outline" className="font-mono uppercase">
                        {p.capability}
                      </Badge>
                      <Badge variant="secondary" className="font-mono">
                        {p.driver}
                      </Badge>
                      <Badge variant={p.enabled ? "default" : "outline"}>
                        {p.enabled ? "enabled" : "disabled"}
                      </Badge>
                      {p.has_api_key ? (
                        <Badge variant="secondary" className="font-mono">
                          key vaulted
                        </Badge>
                      ) : null}
                    </div>
                    <p className="truncate font-mono text-xs text-[var(--text-muted)]">
                      {p.base_url || "(built-in local extractor)"} | priority {p.priority} | id{" "}
                      {p.id}
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="secondary"
                      size="sm"
                      onClick={() =>
                        updateProvider.mutate({ id: p.id, body: { enabled: !p.enabled } })
                      }
                      disabled={updateProvider.isPending}
                    >
                      {p.enabled ? "Disable" : "Enable"}
                    </Button>
                    <Button
                      type="button"
                      variant="destructive"
                      size="icon"
                      aria-label={`Delete provider ${p.name}`}
                      onClick={() => deleteProvider.mutate(p.id)}
                      disabled={deleteProvider.isPending}
                    >
                      <Trash2 className="h-4 w-4" aria-hidden="true" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Media Fallback Combos ({combos.length})</CardTitle>
        </CardHeader>
        <CardContent>
          {combos.length === 0 ? (
            <p className="text-sm text-[var(--text-muted)]">No media combos configured yet.</p>
          ) : (
            <div className="space-y-2.5">
              {combos.map((c) => (
                <div
                  key={c.id}
                  className="flex flex-col gap-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="min-w-0 space-y-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-medium text-[var(--text-primary)]">{c.name}</span>
                      <Badge variant="outline" className="font-mono uppercase">
                        {c.capability}
                      </Badge>
                      <Badge variant={c.enabled ? "default" : "outline"}>
                        {c.enabled ? "enabled" : "disabled"}
                      </Badge>
                    </div>
                    <p className="font-mono text-xs text-[var(--text-muted)]">
                      Chain: {c.targets.join(" -> ")}
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="secondary"
                      size="sm"
                      onClick={() =>
                        updateCombo.mutate({ id: c.id, body: { enabled: !c.enabled } })
                      }
                      disabled={updateCombo.isPending}
                    >
                      {c.enabled ? "Disable" : "Enable"}
                    </Button>
                    <Button
                      type="button"
                      variant="destructive"
                      size="icon"
                      aria-label={`Delete combo ${c.name}`}
                      onClick={() => deleteCombo.mutate(c.id)}
                      disabled={deleteCombo.isPending}
                    >
                      <Trash2 className="h-4 w-4" aria-hidden="true" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
