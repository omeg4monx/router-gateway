"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  Activity,
  AlertTriangle,
  Check,
  CheckCircle2,
  ChevronDown,
  Copy,
  ExternalLink,
  Eye,
  EyeOff,
  Globe,
  Info,
  KeyRound,
  Layers,
  Plus,
  Power,
  RefreshCw,
  Server,
  Settings2,
  ShieldCheck,
  Trash2,
  X,
} from "lucide-react";
import Link from "next/link";
import { type ReactNode, useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { cn } from "@/shared/lib/utils";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import type { ProviderEntry } from "./catalogUtils";
import {
  useCreateProviderConnectionMutation,
  useCreateProviderNodeMutation,
  useExchangeProviderOAuthMutation,
  useStartProviderOAuthMutation,
} from "./hooks";
import { ProviderLogo } from "./ProviderLogo";
import {
  type CreateConnectionFormInput,
  type CreateConnectionFormOutput,
  type CreateNodeFormInput,
  type CreateNodeFormOutput,
  createConnectionSchema,
  createNodeSchema,
  DEFAULT_CREATE_CONNECTION_VALUES,
  DEFAULT_CREATE_NODE_VALUES,
  DEFAULT_OAUTH_EXCHANGE_VALUES,
  DEFAULT_OAUTH_START_VALUES,
  formatCommaSeparatedList,
  formatHeadersText,
  type OAuthExchangeFormInput,
  type OAuthExchangeFormOutput,
  type OAuthStartFormInput,
  type OAuthStartFormOutput,
  type OAuthStartResponse,
  oauthExchangeSchema,
  oauthStartSchema,
  PROVIDER_PROTOCOLS,
  type ProviderConnectionSummary,
  type ProviderNode,
  type ProviderProtocol,
  type TestProviderConnectionResponse,
  toCreateConnectionRequest,
  toCreateNodeRequest,
  toOAuthExchangeRequest,
  toOAuthStartRequest,
} from "./schemas";

const SELECT_CLASS_NAME =
  "flex h-9 w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-field)] px-3 py-1.5 font-mono text-xs text-[var(--color-text-main)] transition-colors focus:border-[var(--color-primary)] focus:outline-none";

const TEXTAREA_CLASS_NAME =
  "flex min-h-[76px] w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-field)] px-3 py-2 font-mono text-xs text-[var(--color-text-main)] placeholder:font-sans placeholder:text-[var(--color-text-muted)] transition-colors focus:border-[var(--color-primary)] focus:outline-none";

function formatShortTime(iso: string | null, isId: boolean): string {
  if (!iso) {
    return isId ? "Belum pernah" : "Never";
  }
  const parsed = new Date(iso);
  if (Number.isNaN(parsed.getTime())) {
    return iso;
  }
  return parsed.toLocaleString([], {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function getCredentialFieldMeta(entry: ProviderEntry | undefined, isId: boolean) {
  if (!entry) {
    return { label: "API Key", placeholder: "sk-..." };
  }
  if (entry.id === "arena-ai") {
    return {
      label: "userToken (Session Token)",
      placeholder: "userToken / evaluation_session token",
    };
  }
  if (entry.authType === "cookie" || entry.sourceCategory === "webCookie") {
    return {
      label: isId ? "Cookie / Session Token" : "Cookie / Session Value",
      placeholder:
        entry.id === "gemini-web"
          ? "__Secure-1PSID=...; __Secure-1PSIDTS=..."
          : entry.id === "grok-web"
            ? "sso=...; sso-rw=..."
            : "Paste session cookie or token...",
    };
  }
  if (entry.id === "qoder" || entry.id === "qoder-cn") {
    return { label: "Personal Access Token (PAT)", placeholder: "pt-..." };
  }
  if (entry.id === "xai") {
    return { label: "xAI API Key", placeholder: "xai-..." };
  }
  return { label: "API Key", placeholder: "sk-..." };
}

function ModalShell({
  open,
  onClose,
  title,
  subtitle,
  icon,
  children,
  variant = "center",
}: {
  open: boolean;
  onClose: () => void;
  title: string;
  subtitle?: string;
  icon: ReactNode;
  children: ReactNode;
  variant?: "center" | "drawer";
}) {
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) {
    return null;
  }

  return (
    <div
      className={cn(
        "fixed inset-0 z-50 flex",
        variant === "drawer" ? "justify-end" : "items-center justify-center p-4",
      )}
    >
      <button
        type="button"
        aria-label="Close modal backdrop"
        onClick={onClose}
        className="fixed inset-0 bg-black/50 backdrop-blur-[2px]"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-label={title}
        className={cn(
          "relative z-10 flex flex-col border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl",
          variant === "drawer"
            ? "h-full w-full max-w-xl border-y-0 border-r-0"
            : "max-h-[90vh] w-full max-w-lg rounded-2xl",
        )}
      >
        <div className="flex items-start justify-between gap-3 border-b border-[var(--color-border-subtle)] px-5 py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/80 p-1.5 text-[var(--color-primary)]">
              {icon}
            </div>
            <div className="min-w-0">
              <h3 className="truncate text-sm font-semibold text-[var(--color-text-main)]">
                {title}
              </h3>
              {subtitle ? (
                <p className="mt-0.5 truncate text-xs text-[var(--color-text-muted)]">{subtitle}</p>
              ) : null}
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close modal"
            className="flex h-8 w-8 shrink-0 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </button>
        </div>

        <div className="custom-scrollbar flex-1 overflow-y-auto px-5 py-4">{children}</div>
      </div>
    </div>
  );
}

/* ============================================================================
 * 1. Add API Key / Credential Modal (xrouter parity: Get Key link, Single/Bulk,
 *    auto-register xrouter node in Go backend if needed)
 * ========================================================================== */
export function AddApiKeyModal({
  open,
  onClose,
  initialProviderId,
  allEntries,
  isId,
  onCreated,
}: {
  open: boolean;
  onClose: () => void;
  initialProviderId: string | null;
  allEntries: ProviderEntry[];
  isId: boolean;
  onCreated: (message: string) => void;
}) {
  const createConnMutation = useCreateProviderConnectionMutation();
  const createNodeMutation = useCreateProviderNodeMutation();

  const [mode, setMode] = useState<"single" | "bulk">("single");
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [bulkText, setBulkText] = useState("");
  const [bulkBusy, setBulkBusy] = useState(false);
  const [bulkError, setBulkError] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const selectableEntries = useMemo(
    () => allEntries.filter((e) => !e.isNativeOAuth || e.hasOAuth),
    [allEntries],
  );

  const form = useForm<CreateConnectionFormInput, unknown, CreateConnectionFormOutput>({
    resolver: zodResolver(createConnectionSchema),
    defaultValues: DEFAULT_CREATE_CONNECTION_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const watchedProviderId = form.watch("provider_id");
  const activeEntry = useMemo(
    () => allEntries.find((e) => e.id === watchedProviderId),
    [allEntries, watchedProviderId],
  );

  const applyEntryDefaults = (providerId: string) => {
    const target = allEntries.find((e) => e.id === providerId);
    if (!target) return;
    form.setValue("base_url", target.baseUrl || "https://api.openai.com/v1", {
      shouldDirty: false,
    });
    form.setValue("models", formatCommaSeparatedList(target.defaultModels), {
      shouldDirty: false,
    });
    form.setValue("headers", formatHeadersText(target.headers), { shouldDirty: false });
    if (!form.getValues("name").trim()) {
      form.setValue("name", `${target.name} Key`, { shouldDirty: false });
    }
  };

  useEffect(() => {
    if (!open) return;
    createConnMutation.reset();
    createNodeMutation.reset();
    setMode("single");
    setShowApiKey(false);
    setShowAdvanced(false);
    setBulkText("");
    setBulkError(null);
    setSubmitError(null);

    const target =
      (initialProviderId ? allEntries.find((e) => e.id === initialProviderId) : undefined) ??
      selectableEntries[0];
    if (!target) return;

    form.reset({
      ...DEFAULT_CREATE_CONNECTION_VALUES,
      provider_id: target.id,
      name: `${target.name} Primary`,
      base_url: target.baseUrl || "https://api.openai.com/v1",
      models: formatCommaSeparatedList(target.defaultModels),
      headers: formatHeadersText(target.headers),
    });
  }, [
    open,
    initialProviderId,
    allEntries,
    selectableEntries,
    form,
    createConnMutation.reset,
    createNodeMutation.reset,
  ]);

  // Ensure xrouter-catalog provider exists as a node in Go backend before creating connection
  const ensureProviderRegistered = async (entry: ProviderEntry | undefined, baseUrl: string) => {
    if (!entry || entry.isBuiltinAdapter || entry.node) {
      return;
    }
    const safeUrl = baseUrl.trim() || entry.baseUrl || "https://api.openai.com/v1";
    await createNodeMutation.mutateAsync({
      id: entry.id,
      name: entry.name,
      protocol: entry.protocol,
      base_url: safeUrl,
      models: entry.defaultModels,
      headers: entry.headers ?? {},
    });
  };

  const onSubmitSingle = async (values: CreateConnectionFormOutput) => {
    setSubmitError(null);
    try {
      await ensureProviderRegistered(activeEntry, values.base_url);
      const created = await createConnMutation.mutateAsync(toCreateConnectionRequest(values));
      onCreated(
        isId
          ? `Koneksi "${created.name}" berhasil ditambahkan.`
          : `Created connection "${created.name}".`,
      );
      onClose();
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Failed to save connection.");
    }
  };

  const onSubmitBulk = async () => {
    setBulkError(null);
    const lines = bulkText
      .split(/\r?\n/)
      .map((l) => l.trim())
      .filter(Boolean);
    if (lines.length === 0) {
      setBulkError(
        isId ? "Masukkan minimal satu baris kredensial." : "Enter at least one credential line.",
      );
      return;
    }
    if (!activeEntry) return;

    setBulkBusy(true);
    try {
      const baseUrl = form.getValues("base_url").trim() || activeEntry.baseUrl;
      await ensureProviderRegistered(activeEntry, baseUrl);
      let okCount = 0;
      let failCount = 0;
      for (let i = 0; i < lines.length; i += 1) {
        const raw = lines[i] ?? "";
        const pipeIdx = raw.indexOf("|");
        const connName =
          pipeIdx > 0 ? raw.slice(0, pipeIdx).trim() : `${activeEntry.name} #${i + 1}`;
        const apiKey = pipeIdx > 0 ? raw.slice(pipeIdx + 1).trim() : raw;
        if (!apiKey) {
          failCount += 1;
          continue;
        }
        try {
          await createConnMutation.mutateAsync({
            provider_id: activeEntry.id,
            name: connName,
            api_key: apiKey,
            base_url: baseUrl || undefined,
            models: activeEntry.defaultModels,
            priority: 100,
            enabled: true,
          });
          okCount += 1;
        } catch {
          failCount += 1;
        }
      }
      onCreated(
        isId
          ? `Impor massal selesai: ${okCount} berhasil${failCount > 0 ? `, ${failCount} gagal` : ""}.`
          : `Bulk import finished: ${okCount} added${failCount > 0 ? `, ${failCount} failed` : ""}.`,
      );
      onClose();
    } catch (err) {
      setBulkError(err instanceof Error ? err.message : "Bulk import failed.");
    } finally {
      setBulkBusy(false);
    }
  };

  const credMeta = getCredentialFieldMeta(activeEntry, isId);
  const keyPortalUrl = activeEntry?.apiKeyUrl ?? activeEntry?.signupUrl ?? activeEntry?.website;

  return (
    <ModalShell
      open={open}
      onClose={onClose}
      title={
        activeEntry
          ? isId
            ? `Setup ${activeEntry.name}`
            : `Setup ${activeEntry.name}`
          : isId
            ? "Tambah Koneksi API Key"
            : "Add API Key Connection"
      }
      subtitle={
        activeEntry
          ? `${activeEntry.id} · ${activeEntry.protocol.toUpperCase()}`
          : isId
            ? "Simpan kredensial terenkripsi untuk routing upstream."
            : "Store an encrypted upstream credential for routing."
      }
      icon={
        activeEntry ? (
          <ProviderLogo
            providerId={activeEntry.id}
            name={activeEntry.name}
            protocol={activeEntry.protocol}
            isCustom={activeEntry.isCustom}
            size={22}
          />
        ) : (
          <KeyRound className="h-4 w-4" aria-hidden="true" />
        )
      }
    >
      <div className="space-y-4">
        {/* Provider Info & Get API Key Banner (xrouter parity) */}
        {activeEntry && (keyPortalUrl || activeEntry.noticeText || activeEntry.authHint) ? (
          <div className="space-y-2 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/65 p-3">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[var(--color-text-main)]">
                <Info className="h-3.5 w-3.5 text-[var(--color-primary)]" aria-hidden="true" />
                <span>{isId ? "Panduan Kredensial" : "Credential Guide"}</span>
              </span>
              {keyPortalUrl ? (
                <a
                  href={keyPortalUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 rounded-md border border-[var(--color-primary)]/30 bg-[var(--color-primary)]/10 px-2 py-0.5 text-[11px] font-medium text-[var(--color-primary)] transition-colors hover:bg-[var(--color-primary)]/20"
                >
                  <span>
                    {activeEntry.apiKeyUrl
                      ? isId
                        ? "Ambil API Key"
                        : "Get API Key"
                      : isId
                        ? "Buka Portal"
                        : "Open Portal"}
                  </span>
                  <ExternalLink className="h-3 w-3" aria-hidden="true" />
                </a>
              ) : null}
            </div>
            {activeEntry.noticeText ? (
              <p className="text-[11px] leading-relaxed text-[var(--color-text-muted)]">
                {activeEntry.noticeText}
              </p>
            ) : null}
            {activeEntry.authHint ? (
              <p className="text-[11px] leading-relaxed text-[var(--color-text-muted)]">
                {activeEntry.authHint}
              </p>
            ) : null}
          </div>
        ) : null}

        {/* Mode Switcher: Single vs Bulk Import */}
        <div className="inline-flex rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)] p-0.5 text-xs">
          <button
            type="button"
            onClick={() => setMode("single")}
            className={cn(
              "cursor-pointer rounded-md px-3 py-1 font-medium transition-colors",
              mode === "single"
                ? "bg-[var(--color-surface)] text-[var(--color-text-main)] shadow-xs"
                : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
            )}
          >
            {isId ? "Satu Koneksi" : "Single Key"}
          </button>
          <button
            type="button"
            onClick={() => setMode("bulk")}
            className={cn(
              "cursor-pointer rounded-md px-3 py-1 font-medium transition-colors",
              mode === "bulk"
                ? "bg-[var(--color-surface)] text-[var(--color-text-main)] shadow-xs"
                : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
            )}
          >
            {isId ? "Impor Massal" : "Bulk Import"}
          </button>
        </div>

        {mode === "single" ? (
          <form className="space-y-3.5" onSubmit={form.handleSubmit(onSubmitSingle)} noValidate>
            {submitError ? (
              <Alert variant="destructive">
                <AlertTitle>{isId ? "Gagal menyimpan" : "Failed to save"}</AlertTitle>
                <AlertDescription>{submitError}</AlertDescription>
              </Alert>
            ) : null}

            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label htmlFor="modal-conn-provider">Provider</Label>
                <select
                  id="modal-conn-provider"
                  className={SELECT_CLASS_NAME}
                  {...form.register("provider_id", {
                    onChange: (e) => applyEntryDefaults(e.target.value),
                  })}
                >
                  {selectableEntries.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.name} ({item.id})
                    </option>
                  ))}
                </select>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="modal-conn-name">{isId ? "Nama Koneksi" : "Connection Name"}</Label>
                <Input
                  id="modal-conn-name"
                  placeholder="Production Key"
                  {...form.register("name")}
                />
                {form.formState.errors.name ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {form.formState.errors.name.message}
                  </p>
                ) : null}
              </div>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="modal-conn-apikey">{credMeta.label}</Label>
              <div className="relative">
                <Input
                  id="modal-conn-apikey"
                  type={showApiKey ? "text" : "password"}
                  autoComplete="off"
                  placeholder={credMeta.placeholder}
                  className="pr-10 font-mono text-xs"
                  {...form.register("api_key")}
                />
                <button
                  type="button"
                  onClick={() => setShowApiKey((v) => !v)}
                  aria-label={showApiKey ? "Hide credential" : "Show credential"}
                  className="absolute top-1/2 right-2 -translate-y-1/2 cursor-pointer rounded p-1 text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]"
                >
                  {showApiKey ? (
                    <EyeOff className="h-3.5 w-3.5" />
                  ) : (
                    <Eye className="h-3.5 w-3.5" />
                  )}
                </button>
              </div>
              {form.formState.errors.api_key ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.api_key.message}
                </p>
              ) : null}
            </div>

            {/* Collapsible Advanced Settings */}
            <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35">
              <button
                type="button"
                onClick={() => setShowAdvanced((v) => !v)}
                className="flex w-full cursor-pointer items-center justify-between px-3 py-2 text-xs font-medium text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]"
              >
                <span>
                  {isId
                    ? "Pengaturan Lanjutan (Base URL, Model, Prioritas, Headers)"
                    : "Advanced Settings (Base URL, Models, Priority, Headers)"}
                </span>
                <ChevronDown
                  className={cn("h-3.5 w-3.5 transition-transform", showAdvanced && "rotate-180")}
                />
              </button>

              {showAdvanced ? (
                <div className="space-y-3 border-t border-[var(--color-border-subtle)] p-3">
                  <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
                    <div className="space-y-1.5 sm:col-span-2">
                      <Label htmlFor="modal-conn-baseurl">Base URL</Label>
                      <Input
                        id="modal-conn-baseurl"
                        placeholder="https://api.openai.com/v1"
                        className="font-mono text-xs"
                        {...form.register("base_url")}
                      />
                    </div>
                    <div className="space-y-1.5">
                      <Label htmlFor="modal-conn-priority">{isId ? "Prioritas" : "Priority"}</Label>
                      <Input
                        id="modal-conn-priority"
                        inputMode="numeric"
                        className="font-mono text-xs"
                        {...form.register("priority")}
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="modal-conn-models">
                      {isId ? "Daftar Model (koma)" : "Models (comma-separated)"}
                    </Label>
                    <Input
                      id="modal-conn-models"
                      className="font-mono text-xs"
                      {...form.register("models")}
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="modal-conn-headers">
                      {isId
                        ? "Custom Headers (Key: Value per baris)"
                        : "Custom Headers (Key: Value)"}
                    </Label>
                    <textarea
                      id="modal-conn-headers"
                      className={TEXTAREA_CLASS_NAME}
                      placeholder="X-Title: omeg4mon-router"
                      {...form.register("headers")}
                    />
                  </div>
                </div>
              ) : null}
            </div>

            <div className="flex items-center justify-end gap-2 border-t border-[var(--color-border-subtle)] pt-3.5">
              <Button type="button" variant="secondary" size="sm" onClick={onClose}>
                {isId ? "Batal" : "Cancel"}
              </Button>
              <Button
                type="submit"
                size="sm"
                disabled={createConnMutation.isPending || createNodeMutation.isPending}
              >
                {createConnMutation.isPending || createNodeMutation.isPending
                  ? isId
                    ? "Menyimpan..."
                    : "Saving..."
                  : isId
                    ? "Simpan Koneksi"
                    : "Save Connection"}
              </Button>
            </div>
          </form>
        ) : (
          <div className="space-y-3.5">
            {bulkError ? (
              <Alert variant="destructive">
                <AlertDescription>{bulkError}</AlertDescription>
              </Alert>
            ) : null}

            <div className="space-y-1.5">
              <Label htmlFor="modal-bulk-keys">
                {isId
                  ? "Satu kredensial per baris (nama|key atau langsung key)"
                  : "One credential per line (name|key or raw key)"}
              </Label>
              <textarea
                id="modal-bulk-keys"
                rows={6}
                value={bulkText}
                onChange={(e) => setBulkText(e.target.value)}
                placeholder={"prod-1|sk-...\nprod-2|sk-...\nsk-auto-named-key"}
                className={TEXTAREA_CLASS_NAME}
              />
            </div>

            <div className="flex items-center justify-end gap-2 border-t border-[var(--color-border-subtle)] pt-3.5">
              <Button type="button" variant="secondary" size="sm" onClick={onClose}>
                {isId ? "Batal" : "Cancel"}
              </Button>
              <Button type="button" size="sm" disabled={bulkBusy} onClick={onSubmitBulk}>
                {bulkBusy
                  ? isId
                    ? "Mengimpor..."
                    : "Importing..."
                  : isId
                    ? "Impor Semua"
                    : "Import All"}
              </Button>
            </div>
          </div>
        )}
      </div>
    </ModalShell>
  );
}

/* ============================================================================
 * 2. OAuth Guided Wizard Modal (with auto callback URL parsing & Token fallback)
 * ========================================================================== */
export function OAuthWizardModal({
  open,
  onClose,
  initialProviderId,
  allEntries,
  isId,
  onConnected,
  onSwitchToApiKey,
}: {
  open: boolean;
  onClose: () => void;
  initialProviderId: string | null;
  allEntries: ProviderEntry[];
  isId: boolean;
  onConnected: (message: string) => void;
  onSwitchToApiKey: (providerId: string) => void;
}) {
  const startOAuthMutation = useStartProviderOAuthMutation();
  const exchangeOAuthMutation = useExchangeProviderOAuthMutation();

  const nativeOAuthEntries = useMemo(() => allEntries.filter((e) => e.isNativeOAuth), [allEntries]);

  const [activeSession, setActiveSession] = useState<{
    provider: string;
    session: OAuthStartResponse;
  } | null>(null);
  const [copied, setCopied] = useState(false);

  const startForm = useForm<OAuthStartFormInput, unknown, OAuthStartFormOutput>({
    resolver: zodResolver(oauthStartSchema),
    defaultValues: DEFAULT_OAUTH_START_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const exchangeForm = useForm<OAuthExchangeFormInput, unknown, OAuthExchangeFormOutput>({
    resolver: zodResolver(oauthExchangeSchema),
    defaultValues: DEFAULT_OAUTH_EXCHANGE_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const selectedEntry = useMemo(
    () => allEntries.find((e) => e.id === initialProviderId) ?? nativeOAuthEntries[0],
    [allEntries, initialProviderId, nativeOAuthEntries],
  );

  useEffect(() => {
    if (!open) return;
    startOAuthMutation.reset();
    exchangeOAuthMutation.reset();
    setActiveSession(null);
    setCopied(false);

    const pid = selectedEntry?.isNativeOAuth
      ? selectedEntry.id
      : (nativeOAuthEntries[0]?.id ?? "claude");
    startForm.reset({
      ...DEFAULT_OAUTH_START_VALUES,
      provider: pid as OAuthStartFormInput["provider"],
    });
    exchangeForm.reset(DEFAULT_OAUTH_EXCHANGE_VALUES);
  }, [
    open,
    selectedEntry,
    nativeOAuthEntries,
    startForm,
    exchangeForm,
    startOAuthMutation.reset,
    exchangeOAuthMutation.reset,
  ]);

  const watchedProvider = startForm.watch("provider");
  const currentEntry = allEntries.find((e) => e.id === watchedProvider) ?? selectedEntry;

  const onStart = (values: OAuthStartFormOutput) => {
    startOAuthMutation.mutate(
      { provider: values.provider, body: toOAuthStartRequest(values) },
      {
        onSuccess: (session) => {
          setActiveSession({ provider: values.provider, session });
          exchangeForm.setValue("state", session.state, { shouldDirty: false });
          if (values.connection_name) {
            exchangeForm.setValue("connection_name", values.connection_name, {
              shouldDirty: false,
            });
          }
        },
      },
    );
  };

  const onExchange = (values: OAuthExchangeFormOutput) => {
    const provider = activeSession?.provider ?? watchedProvider;
    exchangeOAuthMutation.mutate(
      { provider, body: toOAuthExchangeRequest(values) },
      {
        onSuccess: (conn) => {
          onConnected(
            isId
              ? `Akun OAuth "${conn.name}" berhasil dihubungkan.`
              : `Connected OAuth account "${conn.name}".`,
          );
          onClose();
        },
      },
    );
  };

  const handleCopyUrl = async () => {
    if (!activeSession) return;
    try {
      await navigator.clipboard.writeText(activeSession.session.authorize_url);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      // Ignore clipboard failures
    }
  };

  const portalUrl = currentEntry?.signupUrl ?? currentEntry?.website;

  return (
    <ModalShell
      open={open}
      onClose={onClose}
      title={
        currentEntry
          ? isId
            ? `Hubungkan ${currentEntry.name} (OAuth)`
            : `Connect ${currentEntry.name} (OAuth)`
          : isId
            ? "Hubungkan Akun OAuth"
            : "Connect OAuth Account"
      }
      subtitle={
        isId
          ? "Alur PKCE 2 langkah untuk autentikasi akun OAuth upstream."
          : "2-step PKCE flow for upstream OAuth account authentication."
      }
      icon={
        currentEntry ? (
          <ProviderLogo
            providerId={currentEntry.id}
            name={currentEntry.name}
            protocol={currentEntry.protocol}
            size={22}
          />
        ) : (
          <ShieldCheck className="h-4 w-4" aria-hidden="true" />
        )
      }
    >
      <div className="space-y-4">
        {/* Provider Notice & Portal Link */}
        {currentEntry && (portalUrl || currentEntry.deprecated) ? (
          <div className="space-y-2 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/65 p-3">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span className="inline-flex items-center gap-1.5 text-xs font-medium text-[var(--color-text-main)]">
                <AlertTriangle className="h-3.5 w-3.5 text-[var(--status-warn)]" />
                <span>{isId ? "Catatan OAuth" : "OAuth Notice"}</span>
              </span>
              {portalUrl ? (
                <a
                  href={portalUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 text-[11px] font-medium text-[var(--color-primary)] hover:underline"
                >
                  <span>{isId ? "Portal Akun" : "Account Portal"}</span>
                  <ExternalLink className="h-3 w-3" />
                </a>
              ) : null}
            </div>
            <p className="text-[11px] leading-relaxed text-[var(--color-text-muted)]">
              {isId
                ? "OAuth menggunakan sesi akun pribadi Anda. Pastikan Anda menyalin URL callback setelah login selesai."
                : "OAuth uses your personal account session. Copy the final callback URL from your browser after signing in."}
            </p>
          </div>
        ) : null}

        {/* If selected provider is from xrouter OAuth catalog but not yet Go-native OAuth */}
        {selectedEntry && !selectedEntry.isNativeOAuth ? (
          <div className="space-y-3 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/50 p-4">
            <p className="text-xs leading-relaxed text-[var(--color-text-main)]">
              {isId
                ? `Provider "${selectedEntry.name}" belum memiliki handler OAuth PKCE native di backend Go. Anda tetap dapat langsung menghubungkannya menggunakan API Key / Session Token.`
                : `"${selectedEntry.name}" uses API Key / Token import on the Go gateway. Switch to credential setup to connect it immediately.`}
            </p>
            <Button
              type="button"
              size="sm"
              onClick={() => {
                onClose();
                onSwitchToApiKey(selectedEntry.id);
              }}
            >
              <KeyRound className="h-3.5 w-3.5" />
              <span>{isId ? "Setup via API Key / Token" : "Setup via API Key / Token"}</span>
            </Button>
          </div>
        ) : (
          <>
            {/* Step 1 */}
            <form className="space-y-3" onSubmit={startForm.handleSubmit(onStart)} noValidate>
              <div className="flex items-center gap-2">
                <span className="flex h-5 w-5 items-center justify-center rounded-full bg-[var(--color-primary)] font-mono text-[11px] font-bold text-[var(--color-on-primary)]">
                  1
                </span>
                <span className="text-xs font-semibold text-[var(--color-text-main)]">
                  {isId ? "Buat Sesi Otorisasi" : "Start Authorization Session"}
                </span>
              </div>

              {startOAuthMutation.isError ? (
                <Alert variant="destructive">
                  <AlertDescription>{startOAuthMutation.error.message}</AlertDescription>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="modal-oauth-provider">OAuth Provider</Label>
                  <select
                    id="modal-oauth-provider"
                    className={SELECT_CLASS_NAME}
                    {...startForm.register("provider")}
                  >
                    {nativeOAuthEntries.map((item) => (
                      <option key={item.id} value={item.id}>
                        {item.name} ({item.id})
                      </option>
                    ))}
                  </select>
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="modal-oauth-name">
                    {isId ? "Label Koneksi (opsional)" : "Connection Label (optional)"}
                  </Label>
                  <Input
                    id="modal-oauth-name"
                    placeholder={currentEntry ? `${currentEntry.name} Account` : "Primary Account"}
                    {...startForm.register("connection_name")}
                  />
                </div>
              </div>

              <Button
                type="submit"
                size="sm"
                variant="secondary"
                disabled={startOAuthMutation.isPending}
              >
                {startOAuthMutation.isPending
                  ? isId
                    ? "Membuat URL..."
                    : "Generating URL..."
                  : isId
                    ? "Generate URL Otorisasi"
                    : "Generate Authorization URL"}
              </Button>
            </form>

            {activeSession ? (
              <div className="space-y-2.5 rounded-xl border border-[var(--color-primary)]/30 bg-[var(--color-surface-2)]/70 p-3.5">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <span className="text-xs font-semibold text-[var(--color-text-main)]">
                    {isId ? "Buka URL ini & login:" : "Open this URL & sign in:"}
                  </span>
                  <div className="flex items-center gap-1.5">
                    <Button type="button" size="sm" variant="secondary" onClick={handleCopyUrl}>
                      {copied ? (
                        <>
                          <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" />
                          <span>{isId ? "Tersalin" : "Copied"}</span>
                        </>
                      ) : (
                        <>
                          <Copy className="h-3.5 w-3.5" />
                          <span>{isId ? "Salin" : "Copy"}</span>
                        </>
                      )}
                    </Button>
                    <a
                      href={activeSession.session.authorize_url}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex h-8 items-center gap-1.5 rounded-lg bg-[var(--color-primary)] px-3 text-xs font-medium text-[var(--color-on-primary)] transition-opacity hover:opacity-90"
                    >
                      <span>{isId ? "Buka Login" : "Open Login"}</span>
                      <ExternalLink className="h-3.5 w-3.5" />
                    </a>
                  </div>
                </div>
                <p className="break-all rounded-lg bg-[var(--color-surface)] p-2 font-mono text-[11px] text-[var(--color-text-muted)]">
                  {activeSession.session.authorize_url}
                </p>
              </div>
            ) : null}

            {/* Step 2 */}
            <form
              className="space-y-3 border-t border-[var(--color-border-subtle)] pt-4"
              onSubmit={exchangeForm.handleSubmit(onExchange)}
              noValidate
            >
              <div className="flex items-center gap-2">
                <span className="flex h-5 w-5 items-center justify-center rounded-full bg-[var(--color-primary)] font-mono text-[11px] font-bold text-[var(--color-on-primary)]">
                  2
                </span>
                <span className="text-xs font-semibold text-[var(--color-text-main)]">
                  {isId ? "Tukar Callback URL / Kode" : "Exchange Callback URL / Code"}
                </span>
              </div>

              {exchangeOAuthMutation.isError ? (
                <Alert variant="destructive">
                  <AlertDescription>{exchangeOAuthMutation.error.message}</AlertDescription>
                </Alert>
              ) : null}

              <div className="space-y-1.5">
                <Label htmlFor="modal-oauth-callback">
                  {isId ? "Tempel Callback URL" : "Paste Callback URL"}
                </Label>
                <Input
                  id="modal-oauth-callback"
                  placeholder="http://localhost:.../callback?code=...&state=..."
                  className="font-mono text-xs"
                  {...exchangeForm.register("callback_url")}
                />
              </div>

              <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="modal-oauth-code">Code</Label>
                  <Input
                    id="modal-oauth-code"
                    placeholder="Authorization code"
                    className="font-mono text-xs"
                    {...exchangeForm.register("code")}
                  />
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="modal-oauth-state">State</Label>
                  <Input
                    id="modal-oauth-state"
                    placeholder="PKCE state"
                    className="font-mono text-xs"
                    {...exchangeForm.register("state")}
                  />
                </div>
              </div>

              <div className="flex items-center justify-end gap-2 pt-2">
                <Button type="button" variant="secondary" size="sm" onClick={onClose}>
                  {isId ? "Tutup" : "Close"}
                </Button>
                <Button type="submit" size="sm" disabled={exchangeOAuthMutation.isPending}>
                  {exchangeOAuthMutation.isPending
                    ? isId
                      ? "Menghubungkan..."
                      : "Connecting..."
                    : isId
                      ? "Selesaikan OAuth"
                      : "Complete OAuth"}
                </Button>
              </div>
            </form>
          </>
        )}
      </div>
    </ModalShell>
  );
}

/* ============================================================================
 * 3. Add Custom Compatible Node Modal
 * ========================================================================== */
export function AddCustomNodeModal({
  open,
  onClose,
  initialProtocol,
  isId,
  onCreated,
}: {
  open: boolean;
  onClose: () => void;
  initialProtocol: ProviderProtocol;
  isId: boolean;
  onCreated: (node: ProviderNode) => void;
}) {
  const createNodeMutation = useCreateProviderNodeMutation();

  const form = useForm<CreateNodeFormInput, unknown, CreateNodeFormOutput>({
    resolver: zodResolver(createNodeSchema),
    defaultValues: DEFAULT_CREATE_NODE_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  useEffect(() => {
    if (!open) return;
    createNodeMutation.reset();
    form.reset({
      ...DEFAULT_CREATE_NODE_VALUES,
      protocol: initialProtocol,
    });
  }, [open, initialProtocol, form, createNodeMutation.reset]);

  const onSubmit = (values: CreateNodeFormOutput) => {
    createNodeMutation.mutate(toCreateNodeRequest(values), {
      onSuccess: (created) => {
        onCreated(created);
        onClose();
      },
    });
  };

  return (
    <ModalShell
      open={open}
      onClose={onClose}
      title={isId ? "Tambah Custom Endpoint" : "Add Custom Compatible Endpoint"}
      subtitle={
        isId
          ? "Daftarkan gateway OpenAI, Anthropic, atau Gemini kompatibel (vLLM, Ollama, LiteLLM, dll)."
          : "Register an OpenAI, Anthropic, or Gemini compatible gateway (vLLM, Ollama, LiteLLM, etc)."
      }
      icon={<Server className="h-4 w-4" aria-hidden="true" />}
    >
      <form className="space-y-4" onSubmit={form.handleSubmit(onSubmit)} noValidate>
        {createNodeMutation.isError ? (
          <Alert variant="destructive">
            <AlertDescription>{createNodeMutation.error.message}</AlertDescription>
          </Alert>
        ) : null}

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label htmlFor="modal-node-id">Provider ID</Label>
            <Input
              id="modal-node-id"
              placeholder="local-vllm"
              className="font-mono text-xs"
              {...form.register("id")}
            />
            {form.formState.errors.id ? (
              <p className="text-xs text-[var(--status-error)]">
                {form.formState.errors.id.message}
              </p>
            ) : null}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="modal-node-name">{isId ? "Nama Tampilan" : "Display Name"}</Label>
            <Input
              id="modal-node-name"
              placeholder="Local vLLM Cluster"
              {...form.register("name")}
            />
            {form.formState.errors.name ? (
              <p className="text-xs text-[var(--status-error)]">
                {form.formState.errors.name.message}
              </p>
            ) : null}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div className="space-y-1.5">
            <Label htmlFor="modal-node-protocol">{isId ? "Protokol" : "Protocol"}</Label>
            <select
              id="modal-node-protocol"
              className={SELECT_CLASS_NAME}
              {...form.register("protocol")}
            >
              {PROVIDER_PROTOCOLS.map((proto) => (
                <option key={proto} value={proto}>
                  {proto}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-1.5 sm:col-span-2">
            <Label htmlFor="modal-node-url">Base URL</Label>
            <Input
              id="modal-node-url"
              placeholder="http://127.0.0.1:8000/v1"
              className="font-mono text-xs"
              {...form.register("base_url")}
            />
            {form.formState.errors.base_url ? (
              <p className="text-xs text-[var(--status-error)]">
                {form.formState.errors.base_url.message}
              </p>
            ) : null}
          </div>
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="modal-node-models">
            {isId ? "Model Default (pisahkan dengan koma)" : "Default Models (comma-separated)"}
          </Label>
          <Input
            id="modal-node-models"
            placeholder="llama-3.3-70b, qwen2.5-coder-32b"
            className="font-mono text-xs"
            {...form.register("models")}
          />
        </div>

        <div className="space-y-1.5">
          <Label htmlFor="modal-node-headers">
            {isId ? "Headers Default (opsional)" : "Default Headers (optional)"}
          </Label>
          <textarea
            id="modal-node-headers"
            className={TEXTAREA_CLASS_NAME}
            placeholder="X-Custom-Auth: token"
            {...form.register("headers")}
          />
          {form.formState.errors.headers ? (
            <p className="text-xs text-[var(--status-error)]">
              {form.formState.errors.headers.message}
            </p>
          ) : null}
        </div>

        <div className="flex items-center justify-end gap-2 border-t border-[var(--color-border-subtle)] pt-4">
          <Button type="button" variant="secondary" size="sm" onClick={onClose}>
            {isId ? "Batal" : "Cancel"}
          </Button>
          <Button type="submit" size="sm" disabled={createNodeMutation.isPending}>
            {createNodeMutation.isPending
              ? isId
                ? "Mendaftarkan..."
                : "Registering..."
              : isId
                ? "Daftar Endpoint"
                : "Register Endpoint"}
          </Button>
        </div>
      </form>
    </ModalShell>
  );
}

/* ============================================================================
 * 4. Provider Detail & Connections Slide-Over Drawer (rich xrouter metadata)
 * ========================================================================== */
export function ProviderConnectionsDrawer({
  entry,
  onClose,
  isId,
  busyConnectionId,
  liveTestResults,
  onTestConnection,
  onToggleConnection,
  onDeleteConnection,
  onAddConnection,
  onDeleteCustomNode,
  isDeletingNode,
}: {
  entry: ProviderEntry | null;
  onClose: () => void;
  isId: boolean;
  busyConnectionId: string | null;
  liveTestResults: Record<string, TestProviderConnectionResponse>;
  onTestConnection: (conn: ProviderConnectionSummary) => void;
  onToggleConnection: (conn: ProviderConnectionSummary) => void;
  onDeleteConnection: (conn: ProviderConnectionSummary) => void;
  onAddConnection: (entry: ProviderEntry) => void;
  onDeleteCustomNode: (node: ProviderNode) => void;
  isDeletingNode: boolean;
}) {
  const [confirmDeleteConnId, setConfirmDeleteConnId] = useState<string | null>(null);
  const [confirmDeleteNode, setConfirmDeleteNode] = useState(false);
  const entryId = entry?.id ?? null;

  useEffect(() => {
    if (entryId !== null) {
      setConfirmDeleteConnId(null);
      setConfirmDeleteNode(false);
    }
  }, [entryId]);

  if (!entry) {
    return null;
  }

  const keyUrl = entry.apiKeyUrl ?? entry.signupUrl;

  return (
    <ModalShell
      open={Boolean(entry)}
      onClose={onClose}
      variant="drawer"
      title={entry.name}
      subtitle={`${entry.id} · ${entry.protocol.toUpperCase()} · ${entry.baseUrl || "—"}`}
      icon={
        <ProviderLogo
          providerId={entry.id}
          name={entry.name}
          protocol={entry.protocol}
          isCustom={entry.isCustom}
          size={22}
        />
      }
    >
      <div className="space-y-4">
        {/* Provider Metadata & Links Card */}
        <div className="space-y-2.5 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/55 p-3.5">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex flex-wrap items-center gap-1.5">
              <Badge variant={entry.stats.connected > 0 ? "success" : "outline"}>
                {entry.stats.connected} {isId ? "Aktif" : "Active"}
              </Badge>
              {entry.stats.error > 0 ? (
                <Badge variant="destructive">{entry.stats.error} Error</Badge>
              ) : null}
              {entry.serviceKinds?.slice(0, 4).map((kind) => (
                <span
                  key={kind}
                  className="rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-1.5 py-0.5 font-mono text-[10px] uppercase text-[var(--color-text-muted)]"
                >
                  {kind}
                </span>
              ))}
            </div>

            <Button type="button" size="sm" onClick={() => onAddConnection(entry)}>
              <Plus className="h-3.5 w-3.5" aria-hidden="true" />
              <span>
                {entry.isNativeOAuth
                  ? isId
                    ? "Tambah Akun OAuth"
                    : "Connect OAuth"
                  : isId
                    ? "Tambah Kredensial"
                    : "Add Credential"}
              </span>
            </Button>
          </div>

          {entry.website || keyUrl ? (
            <div className="flex flex-wrap items-center gap-3 pt-1 text-xs">
              {entry.website ? (
                <a
                  href={entry.website}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 text-[var(--color-text-muted)] hover:text-[var(--color-primary)]"
                >
                  <Globe className="h-3.5 w-3.5" />
                  <span>{isId ? "Website Resmi" : "Official Website"}</span>
                  <ExternalLink className="h-3 w-3" />
                </a>
              ) : null}
              {keyUrl && keyUrl !== entry.website ? (
                <a
                  href={keyUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1 font-medium text-[var(--color-primary)] hover:underline"
                >
                  <KeyRound className="h-3.5 w-3.5" />
                  <span>
                    {entry.apiKeyUrl
                      ? isId
                        ? "Ambil API Key"
                        : "Get API Key"
                      : isId
                        ? "Daftar / Login"
                        : "Sign Up / Login"}
                  </span>
                  <ExternalLink className="h-3 w-3" />
                </a>
              ) : null}
            </div>
          ) : null}

          {entry.noticeText ? (
            <p className="text-[11px] leading-relaxed text-[var(--color-text-muted)]">
              {entry.noticeText}
            </p>
          ) : null}

          {entry.defaultModels.length > 0 ? (
            <div className="space-y-1 pt-1">
              <span className="inline-flex items-center gap-1 text-[11px] font-medium text-[var(--color-text-muted)]">
                <Layers className="h-3 w-3" />
                <span>{isId ? "Model Default" : "Default Models"}</span>
              </span>
              <div className="flex flex-wrap gap-1">
                {entry.defaultModels.slice(0, 8).map((m) => (
                  <span
                    key={m}
                    className="rounded-md bg-[var(--color-surface)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--color-text-muted)]"
                  >
                    {m}
                  </span>
                ))}
                {entry.defaultModels.length > 8 ? (
                  <span className="px-1 font-mono text-[10px] text-[var(--color-text-subtle)]">
                    +{entry.defaultModels.length - 8}
                  </span>
                ) : null}
              </div>
            </div>
          ) : null}
        </div>

        {/* Connections List */}
        {entry.connections.length === 0 ? (
          <div className="rounded-xl border border-dashed border-[var(--color-border)] p-6 text-center">
            <p className="text-sm font-medium text-[var(--color-text-main)]">
              {isId ? "Belum ada koneksi untuk provider ini" : "No connections configured yet"}
            </p>
            <p className="mt-1 text-xs text-[var(--color-text-muted)]">
              {isId
                ? "Tambahkan kredensial pertama agar model dari provider ini dapat digunakan."
                : "Add your first credential to start routing requests through this provider."}
            </p>
          </div>
        ) : (
          <div className="space-y-3">
            {entry.connections.map((conn) => {
              const live = liveTestResults[conn.id];
              const isBusy = busyConnectionId === conn.id;
              const isConfirming = confirmDeleteConnId === conn.id;

              return (
                <div
                  key={conn.id}
                  className="rounded-xl border border-[var(--color-border)] bg-[var(--color-surface)] p-3.5 transition-colors hover:border-[var(--color-primary)]/35"
                >
                  <div className="flex items-start justify-between gap-2">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="truncate text-sm font-semibold text-[var(--color-text-main)]">
                          {conn.name}
                        </span>
                        <Badge variant={conn.enabled ? "success" : "outline"}>
                          {conn.enabled
                            ? isId
                              ? "Aktif"
                              : "Enabled"
                            : isId
                              ? "Nonaktif"
                              : "Paused"}
                        </Badge>
                        {live ? (
                          <Badge
                            variant={live.ok ? "success" : "destructive"}
                            className="tabular-nums"
                          >
                            {live.ok
                              ? `OK ${live.status_code} · ${live.latency_ms}ms`
                              : `FAIL ${live.status_code || "ERR"} · ${live.latency_ms}ms`}
                          </Badge>
                        ) : conn.last_test_status === "ok" ? (
                          <Badge variant="success">
                            <CheckCircle2 className="mr-1 h-3 w-3" />
                            Healthy
                          </Badge>
                        ) : conn.last_test_status === "error" ? (
                          <Badge variant="destructive">Error</Badge>
                        ) : null}
                      </div>

                      <p className="mt-1 truncate font-mono text-[11px] text-[var(--color-text-muted)]">
                        {conn.base_url} · P{conn.priority} · {isId ? "Tes terakhir:" : "Tested:"}{" "}
                        {formatShortTime(conn.last_tested_at, isId)}
                      </p>
                    </div>
                  </div>

                  {conn.last_test_error ? (
                    <p className="mt-2 rounded-lg bg-[var(--status-error-bg)] px-2.5 py-1.5 font-mono text-[11px] text-[var(--status-error)]">
                      {conn.last_test_error}
                    </p>
                  ) : null}

                  <div className="mt-3 flex flex-wrap items-center justify-between gap-2 border-t border-[var(--color-border-subtle)] pt-2.5">
                    <div className="flex flex-wrap items-center gap-1.5">
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        disabled={busyConnectionId !== null}
                        onClick={() => onTestConnection(conn)}
                      >
                        {isBusy ? (
                          <RefreshCw className="h-3.5 w-3.5 animate-spin" />
                        ) : (
                          <Activity className="h-3.5 w-3.5" />
                        )}
                        <span>{isId ? "Tes" : "Test"}</span>
                      </Button>

                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        disabled={busyConnectionId !== null}
                        onClick={() => onToggleConnection(conn)}
                      >
                        <Power className="h-3.5 w-3.5" />
                        <span>
                          {conn.enabled
                            ? isId
                              ? "Nonaktifkan"
                              : "Disable"
                            : isId
                              ? "Aktifkan"
                              : "Enable"}
                        </span>
                      </Button>

                      <Link
                        href={`/dashboard/providers/detail/?id=${encodeURIComponent(conn.id)}`}
                        onClick={onClose}
                        className="inline-flex h-8 items-center gap-1.5 rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)] px-2.5 text-xs font-medium text-[var(--color-text-main)] transition-colors hover:border-[var(--color-primary)]/40 hover:text-[var(--color-primary)]"
                      >
                        <Settings2 className="h-3.5 w-3.5" />
                        <span>{isId ? "Edit Detail" : "Edit Detail"}</span>
                      </Link>
                    </div>

                    {isConfirming ? (
                      <div className="flex items-center gap-1.5">
                        <Button
                          type="button"
                          size="sm"
                          variant="destructive"
                          disabled={busyConnectionId !== null}
                          onClick={() => {
                            onDeleteConnection(conn);
                            setConfirmDeleteConnId(null);
                          }}
                        >
                          {isId ? "Konfirmasi" : "Confirm"}
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          variant="secondary"
                          onClick={() => setConfirmDeleteConnId(null)}
                        >
                          {isId ? "Batal" : "Cancel"}
                        </Button>
                      </div>
                    ) : (
                      <button
                        type="button"
                        disabled={busyConnectionId !== null}
                        onClick={() => setConfirmDeleteConnId(conn.id)}
                        aria-label={`Delete ${conn.name}`}
                        className="flex h-8 w-8 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--status-error-bg)] hover:text-[var(--status-error)] disabled:opacity-40"
                      >
                        <Trash2 className="h-3.5 w-3.5" />
                      </button>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* Delete Custom Node Option */}
        {entry.node && entry.isCustom ? (
          <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/50 p-3.5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div>
                <p className="text-xs font-semibold text-[var(--color-text-main)]">
                  {isId ? "Hapus Custom Endpoint" : "Remove Custom Endpoint"}
                </p>
                <p className="text-[11px] text-[var(--color-text-muted)]">
                  {isId
                    ? "Menghapus definisi node ini dari katalog custom provider."
                    : "Removes this custom node definition from your provider catalog."}
                </p>
              </div>
              {confirmDeleteNode ? (
                <div className="flex items-center gap-1.5">
                  <Button
                    type="button"
                    size="sm"
                    variant="destructive"
                    disabled={isDeletingNode}
                    onClick={() => {
                      if (entry.node) {
                        onDeleteCustomNode(entry.node);
                        onClose();
                      }
                    }}
                  >
                    {isId ? "Ya, Hapus" : "Yes, Remove"}
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    onClick={() => setConfirmDeleteNode(false)}
                  >
                    {isId ? "Batal" : "Cancel"}
                  </Button>
                </div>
              ) : (
                <Button
                  type="button"
                  size="sm"
                  variant="secondary"
                  onClick={() => setConfirmDeleteNode(true)}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                  <span>{isId ? "Hapus Node" : "Delete Node"}</span>
                </Button>
              )}
            </div>
          </div>
        ) : null}
      </div>
    </ModalShell>
  );
}
