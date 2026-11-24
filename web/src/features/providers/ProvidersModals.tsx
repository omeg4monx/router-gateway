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
  Search,
  Server,
  Settings2,
  ShieldCheck,
  Sparkles,
  Trash2,
  X,
  Zap,
} from "lucide-react";
import Link from "next/link";
import { type ReactNode, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { cn } from "@/shared/lib/utils";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import type { ProviderEntry } from "./catalogUtils";
import {
  fetchProviderOAuthStatus,
  useCreateProviderConnectionMutation,
  useCreateProviderNodeMutation,
  useExchangeProviderOAuthMutation,
  usePasteProviderOAuthMutation,
  useStartProviderOAuthMutation,
  useStopProviderOAuthMutation,
  useTestProviderNodeDraftMutation,
} from "./hooks";
import {
  buildTestFingerprint,
  CUSTOM_ENDPOINT_PRESETS,
  type CustomEndpointPreset,
  formatOAuthClientIdEnvVar,
  parseBulkCredentials,
  parseOAuthCallbackUrl,
  slugifyProviderId,
} from "./modalUtils";
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
  DEFAULT_PRIORITY,
  formatCommaSeparatedList,
  formatHeadersText,
  isValidHttpUrl,
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
  parseCommaSeparatedList,
  parseHeadersText,
  type TestProviderConnectionResponse,
  toCreateConnectionRequest,
  toCreateNodeRequest,
  toOAuthExchangeRequest,
  toOAuthStartRequest,
} from "./schemas";

const SELECT_CLASS_NAME =
  "flex h-9 w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-field)] px-3 py-1.5 font-mono text-xs text-[var(--color-text-main)] transition-colors focus:border-[var(--color-primary)] focus:outline-none";

const TEXTAREA_CLASS_NAME =
  "flex min-h-[80px] w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-field)] px-3 py-2 font-mono text-xs text-[var(--color-text-main)] placeholder:font-sans placeholder:text-[var(--color-text-muted)] transition-colors focus:border-[var(--color-primary)] focus:outline-none";

type ModalSize = "md" | "lg" | "xl";

const MODAL_SIZE_CLASS: Record<ModalSize, string> = {
  md: "max-w-lg",
  lg: "max-w-xl",
  xl: "max-w-2xl",
};

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

function DraftTestResultCard({
  result,
  localError,
  isStale,
  protocol,
  isId,
}: {
  result: TestProviderConnectionResponse | null;
  localError: string | null;
  isStale: boolean;
  protocol?: string;
  isId: boolean;
}) {
  if (localError) {
    return (
      <div className="flex items-start gap-2.5 rounded-xl border border-[var(--status-error)]/30 bg-[var(--status-error-bg)] px-3.5 py-2.5 text-xs text-[var(--status-error)]">
        <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
        <p className="leading-relaxed">{localError}</p>
      </div>
    );
  }

  if (!result) {
    return null;
  }

  return (
    <div
      className={cn(
        "space-y-1.5 rounded-xl border px-3.5 py-2.5 text-xs transition-colors",
        result.ok
          ? "border-[var(--status-ok)]/35 bg-[var(--status-ok)]/10 text-[var(--color-text-main)]"
          : "border-[var(--status-error)]/35 bg-[var(--status-error-bg)] text-[var(--status-error)]",
      )}
    >
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex items-center gap-2 font-medium">
          {result.ok ? (
            <CheckCircle2 className="h-4 w-4 shrink-0 text-[var(--status-ok)]" aria-hidden="true" />
          ) : (
            <AlertTriangle className="h-4 w-4 shrink-0" aria-hidden="true" />
          )}
          <span className="font-mono text-[11px] tabular-nums">
            {result.ok
              ? `OK ${result.status_code} · ${result.latency_ms}ms`
              : `FAIL ${result.status_code || "ERR"} · ${result.latency_ms}ms`}
          </span>
          {result.ok ? (
            <span className="text-[var(--color-text-muted)]">
              {isId
                ? `Endpoint merespons${protocol ? ` (${protocol})` : ""}`
                : `Endpoint responded${protocol ? ` (${protocol})` : ""}`}
            </span>
          ) : null}
        </div>

        {isStale ? (
          <span className="rounded-md border border-[var(--status-warn)]/40 bg-[var(--status-warn)]/15 px-1.5 py-0.5 text-[10px] font-medium text-[var(--status-warn)]">
            {isId ? "Konfigurasi berubah · tes ulang" : "Config changed · re-test"}
          </span>
        ) : null}
      </div>

      {!result.ok && result.error ? (
        <p className="break-words font-mono text-[11px] leading-relaxed opacity-90">
          {result.error}
        </p>
      ) : null}
    </div>
  );
}

function ModalShell({
  open,
  onClose,
  title,
  subtitle,
  badge,
  icon,
  children,
  footer,
  variant = "center",
  size = "lg",
}: {
  open: boolean;
  onClose: () => void;
  title: string;
  subtitle?: string;
  badge?: ReactNode;
  icon: ReactNode;
  children: ReactNode;
  footer?: ReactNode;
  variant?: "center" | "drawer";
  size?: ModalSize;
}) {
  useEffect(() => {
    if (!open) return;
    const prevOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };
    document.addEventListener("keydown", onKey);
    return () => {
      document.body.style.overflow = prevOverflow;
      document.removeEventListener("keydown", onKey);
    };
  }, [open, onClose]);

  if (!open) {
    return null;
  }

  return (
    <div
      className={cn(
        "fixed inset-0 z-50 flex",
        variant === "drawer" ? "justify-end" : "items-center justify-center p-4 sm:p-6",
      )}
    >
      <button
        type="button"
        aria-label="Close modal backdrop"
        onClick={onClose}
        className="fixed inset-0 bg-black/55 backdrop-blur-[2px]"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-label={title}
        className={cn(
          "relative z-10 flex flex-col overflow-hidden border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl",
          variant === "drawer"
            ? cn("h-full w-full border-y-0 border-r-0", MODAL_SIZE_CLASS[size])
            : cn("max-h-[90vh] w-full rounded-2xl", MODAL_SIZE_CLASS[size]),
        )}
      >
        <div className="flex shrink-0 items-start justify-between gap-3 border-b border-[var(--color-border-subtle)] px-4 py-3.5 sm:px-6 sm:py-4">
          <div className="flex min-w-0 items-center gap-3.5">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/85 p-1.5 text-[var(--color-primary)]">
              {icon}
            </div>
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2">
                <h3 className="truncate text-base font-semibold tracking-tight text-[var(--color-text-main)]">
                  {title}
                </h3>
                {badge}
              </div>
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

        <div className="custom-scrollbar flex-1 overflow-y-auto px-4 py-4 sm:px-6 sm:py-5">
          {children}
        </div>

        {footer ? (
          <div className="shrink-0 border-t border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35 px-4 py-3 sm:px-6 sm:py-3.5">
            {footer}
          </div>
        ) : null}
      </div>
    </div>
  );
}

/* ============================================================================
 * 1. Add API Key / Credential Modal (searchable picker, Single/Bulk, pre-save test)
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
  const testDraftMutation = useTestProviderNodeDraftMutation();

  const [mode, setMode] = useState<"single" | "bulk">("single");
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [pickerOpen, setPickerOpen] = useState(false);
  const [providerQuery, setProviderQuery] = useState("");
  const [bulkText, setBulkText] = useState("");
  const [bulkBusy, setBulkBusy] = useState(false);
  const [bulkError, setBulkError] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [testState, setTestState] = useState<{
    result: TestProviderConnectionResponse;
    fingerprint: string;
  } | null>(null);
  const [testValidationError, setTestValidationError] = useState<string | null>(null);

  const selectableEntries = useMemo(
    () => allEntries.filter((e) => !e.isNativeOAuth || e.hasOAuth),
    [allEntries],
  );

  const filteredProviders = useMemo(() => {
    const q = providerQuery.trim().toLowerCase();
    if (!q) {
      return selectableEntries;
    }
    return selectableEntries.filter(
      (item) =>
        item.name.toLowerCase().includes(q) ||
        item.id.toLowerCase().includes(q) ||
        item.protocol.toLowerCase().includes(q),
    );
  }, [selectableEntries, providerQuery]);

  const form = useForm<CreateConnectionFormInput, unknown, CreateConnectionFormOutput>({
    resolver: zodResolver(createConnectionSchema),
    defaultValues: DEFAULT_CREATE_CONNECTION_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const watchedProviderId = form.watch("provider_id");
  const watchedBaseUrl = form.watch("base_url");
  const watchedApiKey = form.watch("api_key");
  const watchedPriority = form.watch("priority");
  const watchedModels = form.watch("models");
  const watchedHeaders = form.watch("headers");

  const activeEntry = useMemo(
    () => allEntries.find((e) => e.id === watchedProviderId),
    [allEntries, watchedProviderId],
  );

  const currentFingerprint = useMemo(
    () =>
      buildTestFingerprint({
        providerId: activeEntry?.id ?? watchedProviderId,
        protocol: activeEntry?.protocol ?? "openai-chat",
        baseUrl: watchedBaseUrl,
        apiKey: watchedApiKey,
        headers: watchedHeaders,
      }),
    [activeEntry, watchedProviderId, watchedBaseUrl, watchedApiKey, watchedHeaders],
  );

  const isTestStale = Boolean(testState && testState.fingerprint !== currentFingerprint);

  const applyEntryDefaults = (providerId: string) => {
    const target = allEntries.find((e) => e.id === providerId);
    if (!target) return;
    form.setValue("provider_id", target.id, { shouldDirty: true });
    form.setValue("base_url", target.baseUrl || "https://api.openai.com/v1", {
      shouldDirty: false,
    });
    form.setValue("models", formatCommaSeparatedList(target.defaultModels), {
      shouldDirty: false,
    });
    form.setValue("headers", formatHeadersText(target.headers), { shouldDirty: false });
    form.setValue("name", `${target.name} Primary`, { shouldDirty: false });
    setTestState(null);
    setTestValidationError(null);
    setPickerOpen(false);
    setProviderQuery("");
  };

  useEffect(() => {
    if (!open) return;
    createConnMutation.reset();
    createNodeMutation.reset();
    testDraftMutation.reset();
    setMode("single");
    setShowApiKey(false);
    setShowAdvanced(false);
    setPickerOpen(!initialProviderId);
    setProviderQuery("");
    setBulkText("");
    setBulkError(null);
    setSubmitError(null);
    setTestState(null);
    setTestValidationError(null);

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
    testDraftMutation.reset,
  ]);

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

  const handleTestDraft = async () => {
    setTestValidationError(null);
    const values = form.getValues();
    const baseUrl = values.base_url.trim();
    if (!isValidHttpUrl(baseUrl)) {
      setTestValidationError(
        isId
          ? "Base URL harus berupa URL http:// atau https:// yang valid sebelum menguji."
          : "Base URL must be a valid http:// or https:// URL before testing.",
      );
      return;
    }
    const parsedHeaders = parseHeadersText(values.headers);
    if (!parsedHeaders.ok) {
      setTestValidationError(parsedHeaders.error);
      return;
    }
    const apiKey = values.api_key.trim();
    if (!apiKey) {
      setTestValidationError(
        isId
          ? "Masukkan API Key / kredensial terlebih dahulu untuk menguji koneksi."
          : "Enter an API key or credential before testing connectivity.",
      );
      return;
    }

    const fingerprint = buildTestFingerprint({
      providerId: activeEntry?.id ?? values.provider_id,
      protocol: activeEntry?.protocol ?? "openai-chat",
      baseUrl,
      apiKey,
      headers: values.headers,
    });

    try {
      const result = await testDraftMutation.mutateAsync({
        provider_id: activeEntry?.id ?? values.provider_id,
        protocol: activeEntry?.protocol ?? "openai-chat",
        base_url: baseUrl,
        api_key: apiKey,
        headers: parsedHeaders.headers,
      });
      setTestState({ result, fingerprint });
    } catch (err) {
      setTestValidationError(err instanceof Error ? err.message : "Connectivity test failed.");
    }
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

  const parsedBulk = useMemo(
    () => parseBulkCredentials(bulkText, activeEntry?.name ?? "Provider"),
    [bulkText, activeEntry?.name],
  );

  const onSubmitBulk = async () => {
    setBulkError(null);
    if (parsedBulk.entries.length === 0) {
      setBulkError(
        isId
          ? "Masukkan minimal satu baris kredensial yang valid."
          : "Enter at least one valid credential line.",
      );
      return;
    }
    if (!activeEntry) return;

    setBulkBusy(true);
    try {
      const baseUrl = form.getValues("base_url").trim() || activeEntry.baseUrl;
      await ensureProviderRegistered(activeEntry, baseUrl);
      let okCount = 0;
      let failCount = parsedBulk.invalidLines;
      for (const item of parsedBulk.entries) {
        try {
          await createConnMutation.mutateAsync({
            provider_id: activeEntry.id,
            name: item.name,
            api_key: item.apiKey,
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
  const modelCount = parseCommaSeparatedList(watchedModels).length;
  const isSaving = createConnMutation.isPending || createNodeMutation.isPending;

  return (
    <ModalShell
      open={open}
      onClose={onClose}
      size="lg"
      title={
        activeEntry
          ? `Setup ${activeEntry.name}`
          : isId
            ? "Tambah Koneksi API Key"
            : "Add API Key Connection"
      }
      subtitle={
        activeEntry
          ? `${activeEntry.id} · ${activeEntry.baseUrl || "Custom URL"}`
          : isId
            ? "Simpan kredensial terenkripsi untuk routing upstream."
            : "Store an encrypted upstream credential for routing."
      }
      badge={
        activeEntry ? (
          <Badge variant="outline" className="font-mono text-[10px] uppercase">
            {activeEntry.protocol}
          </Badge>
        ) : null
      }
      icon={
        activeEntry ? (
          <ProviderLogo
            providerId={activeEntry.id}
            name={activeEntry.name}
            protocol={activeEntry.protocol}
            isCustom={activeEntry.isCustom}
            size={24}
          />
        ) : (
          <KeyRound className="h-4 w-4" aria-hidden="true" />
        )
      }
      footer={
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div>
            {mode === "single" ? (
              <Button
                type="button"
                variant="secondary"
                size="sm"
                disabled={testDraftMutation.isPending || isSaving}
                onClick={() => void handleTestDraft()}
              >
                {testDraftMutation.isPending ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                ) : (
                  <Zap className="h-3.5 w-3.5" aria-hidden="true" />
                )}
                <span>
                  {testDraftMutation.isPending
                    ? isId
                      ? "Menguji..."
                      : "Testing..."
                    : isId
                      ? "Tes Koneksi"
                      : "Test Connection"}
                </span>
              </Button>
            ) : (
              <span className="text-xs text-[var(--color-text-muted)]">
                {parsedBulk.entries.length > 0
                  ? isId
                    ? `${parsedBulk.entries.length} kredensial siap diimpor`
                    : `${parsedBulk.entries.length} credentials ready`
                  : isId
                    ? "Format: nama|key atau langsung key"
                    : "Format: name|key or raw key"}
              </span>
            )}
          </div>

          <div className="flex items-center gap-2">
            <Button type="button" variant="secondary" size="sm" onClick={onClose}>
              {isId ? "Batal" : "Cancel"}
            </Button>
            {mode === "single" ? (
              <Button type="submit" form="add-apikey-form" size="sm" disabled={isSaving}>
                {isSaving
                  ? isId
                    ? "Menyimpan..."
                    : "Saving..."
                  : isId
                    ? "Simpan Koneksi"
                    : "Save Connection"}
              </Button>
            ) : (
              <Button
                type="button"
                size="sm"
                disabled={bulkBusy || parsedBulk.entries.length === 0}
                onClick={() => void onSubmitBulk()}
              >
                {bulkBusy
                  ? isId
                    ? "Mengimpor..."
                    : "Importing..."
                  : isId
                    ? `Impor ${parsedBulk.entries.length || ""} Koneksi`.trim()
                    : `Import ${parsedBulk.entries.length || ""} Keys`.trim()}
              </Button>
            )}
          </div>
        </div>
      }
    >
      <div className="space-y-4">
        {/* Selected Provider Card + Searchable Picker */}
        <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/55 p-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex min-w-0 items-center gap-2.5">
              {activeEntry ? (
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg border border-[var(--color-border)] bg-[var(--color-surface)] p-1">
                  <ProviderLogo
                    providerId={activeEntry.id}
                    name={activeEntry.name}
                    protocol={activeEntry.protocol}
                    isCustom={activeEntry.isCustom}
                    size={20}
                  />
                </div>
              ) : null}
              <div className="min-w-0">
                <p className="truncate text-xs font-semibold text-[var(--color-text-main)]">
                  {activeEntry?.name ?? (isId ? "Pilih Provider" : "Select Provider")}
                </p>
                <p className="truncate font-mono text-[10px] text-[var(--color-text-muted)]">
                  {activeEntry ? `${activeEntry.id} · ${activeEntry.protocol}` : "—"}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-1.5">
              {keyPortalUrl ? (
                <a
                  href={keyPortalUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex h-7 items-center gap-1 rounded-lg border border-[var(--color-primary)]/30 bg-[var(--color-primary)]/10 px-2.5 text-[11px] font-medium text-[var(--color-primary)] transition-colors hover:bg-[var(--color-primary)]/20"
                >
                  <span>
                    {activeEntry?.apiKeyUrl
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
              <Button
                type="button"
                variant="secondary"
                size="sm"
                className="h-7 px-2.5 text-[11px]"
                onClick={() => setPickerOpen((v) => !v)}
              >
                {pickerOpen ? (isId ? "Tutup" : "Done") : isId ? "Ganti Provider" : "Change"}
              </Button>
            </div>
          </div>

          {activeEntry && (activeEntry.noticeText || activeEntry.authHint) ? (
            <div className="mt-2.5 flex items-start gap-1.5 border-t border-[var(--color-border-subtle)] pt-2 text-[11px] leading-relaxed text-[var(--color-text-muted)]">
              <Info
                className="mt-0.5 h-3.5 w-3.5 shrink-0 text-[var(--color-primary)]"
                aria-hidden="true"
              />
              <div className="space-y-1">
                {activeEntry.noticeText ? <p>{activeEntry.noticeText}</p> : null}
                {activeEntry.authHint ? <p>{activeEntry.authHint}</p> : null}
              </div>
            </div>
          ) : null}

          {pickerOpen ? (
            <div className="mt-3 space-y-2 border-t border-[var(--color-border-subtle)] pt-3">
              <div className="relative">
                <Search
                  className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  value={providerQuery}
                  onChange={(e) => setProviderQuery(e.target.value)}
                  placeholder={
                    isId
                      ? "Cari provider (mis. openai, deepseek, groq)..."
                      : "Search provider (e.g. openai, deepseek, groq)..."
                  }
                  className="h-8 pl-8 text-xs"
                />
              </div>
              <div className="custom-scrollbar grid max-h-48 grid-cols-1 gap-1 overflow-y-auto pr-1 sm:grid-cols-2">
                {filteredProviders.map((item) => {
                  const isSelected = item.id === watchedProviderId;
                  return (
                    <button
                      key={item.id}
                      type="button"
                      onClick={() => applyEntryDefaults(item.id)}
                      className={cn(
                        "flex cursor-pointer items-center justify-between gap-2 rounded-lg border px-2.5 py-1.5 text-left transition-colors",
                        isSelected
                          ? "border-[var(--color-primary)]/50 bg-[var(--color-primary)]/10 text-[var(--color-text-main)]"
                          : "border-transparent bg-[var(--color-surface)]/80 text-[var(--color-text-muted)] hover:border-[var(--color-border)] hover:text-[var(--color-text-main)]",
                      )}
                    >
                      <div className="flex min-w-0 items-center gap-2">
                        <ProviderLogo
                          providerId={item.id}
                          name={item.name}
                          protocol={item.protocol}
                          isCustom={item.isCustom}
                          size={16}
                        />
                        <div className="min-w-0">
                          <p className="truncate text-xs font-medium text-[var(--color-text-main)]">
                            {item.name}
                          </p>
                          <p className="truncate font-mono text-[10px] text-[var(--color-text-subtle)]">
                            {item.id}
                          </p>
                        </div>
                      </div>
                      {isSelected ? (
                        <Check
                          className="h-3.5 w-3.5 shrink-0 text-[var(--color-primary)]"
                          aria-hidden="true"
                        />
                      ) : null}
                    </button>
                  );
                })}
                {filteredProviders.length === 0 ? (
                  <p className="col-span-full py-4 text-center text-xs text-[var(--color-text-muted)]">
                    {isId ? "Provider tidak ditemukan." : "No matching providers."}
                  </p>
                ) : null}
              </div>
            </div>
          ) : null}
        </div>

        {/* Mode Switcher */}
        <div className="grid grid-cols-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)] p-1 text-xs">
          <button
            type="button"
            onClick={() => setMode("single")}
            className={cn(
              "cursor-pointer rounded-lg px-3 py-1.5 font-medium transition-colors",
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
              "cursor-pointer rounded-lg px-3 py-1.5 font-medium transition-colors",
              mode === "bulk"
                ? "bg-[var(--color-surface)] text-[var(--color-text-main)] shadow-xs"
                : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
            )}
          >
            {isId ? "Impor Massal" : "Bulk Import"}
          </button>
        </div>

        {mode === "single" ? (
          <form
            id="add-apikey-form"
            className="space-y-4"
            onSubmit={form.handleSubmit(onSubmitSingle)}
            noValidate
          >
            {submitError ? (
              <Alert variant="destructive">
                <AlertTitle>{isId ? "Gagal menyimpan" : "Failed to save"}</AlertTitle>
                <AlertDescription>{submitError}</AlertDescription>
              </Alert>
            ) : null}

            <input type="hidden" {...form.register("provider_id")} />

            <div className="space-y-1.5">
              <Label htmlFor="modal-conn-name">{isId ? "Nama Koneksi" : "Connection Name"}</Label>
              <Input id="modal-conn-name" placeholder="Production Key" {...form.register("name")} />
              {form.formState.errors.name ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.name.message}
                </p>
              ) : null}
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
                    <EyeOff className="h-3.5 w-3.5" aria-hidden="true" />
                  ) : (
                    <Eye className="h-3.5 w-3.5" aria-hidden="true" />
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
                className="flex w-full cursor-pointer items-center justify-between gap-2 px-3.5 py-2.5 text-xs font-medium text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]"
              >
                <span className="flex items-center gap-2">
                  <Settings2 className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Pengaturan Lanjutan" : "Advanced Settings"}</span>
                </span>
                <span className="flex items-center gap-2">
                  <span className="rounded-md bg-[var(--color-surface)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--color-text-subtle)]">
                    P{watchedPriority || "100"} · {modelCount} {isId ? "model" : "models"}
                  </span>
                  <ChevronDown
                    className={cn("h-3.5 w-3.5 transition-transform", showAdvanced && "rotate-180")}
                    aria-hidden="true"
                  />
                </span>
              </button>

              {showAdvanced ? (
                <div className="space-y-3 border-t border-[var(--color-border-subtle)] p-3.5">
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
                      {isId ? "Daftar Model (pisahkan dengan koma)" : "Models (comma-separated)"}
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
                        : "Custom Headers (Key: Value per line)"}
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

            <DraftTestResultCard
              result={testState?.result ?? null}
              localError={testValidationError}
              isStale={isTestStale}
              protocol={activeEntry?.protocol}
              isId={isId}
            />
          </form>
        ) : (
          <div className="space-y-3.5">
            {bulkError ? (
              <Alert variant="destructive">
                <AlertDescription>{bulkError}</AlertDescription>
              </Alert>
            ) : null}

            <div className="space-y-1.5">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <Label htmlFor="modal-bulk-keys">
                  {isId
                    ? "Satu kredensial per baris (nama|key atau langsung key)"
                    : "One credential per line (name|key or raw key)"}
                </Label>
                <div className="flex items-center gap-1.5">
                  {parsedBulk.entries.length > 0 ? (
                    <Badge variant="success" className="font-mono text-[10px]">
                      {parsedBulk.entries.length} {isId ? "valid" : "valid"}
                    </Badge>
                  ) : null}
                  {parsedBulk.invalidLines > 0 ? (
                    <Badge variant="destructive" className="font-mono text-[10px]">
                      {parsedBulk.invalidLines} {isId ? "tidak valid" : "invalid"}
                    </Badge>
                  ) : null}
                </div>
              </div>
              <textarea
                id="modal-bulk-keys"
                rows={7}
                value={bulkText}
                onChange={(e) => setBulkText(e.target.value)}
                placeholder={"prod-1|sk-...\nprod-2|sk-...\nsk-auto-named-key"}
                className={TEXTAREA_CLASS_NAME}
              />
            </div>
          </div>
        )}
      </div>
    </ModalShell>
  );
}

/* ============================================================================
 * 2. OAuth Guided Wizard Modal (2-step stepper + auto callback URL parsing)
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
  const pasteOAuthMutation = usePasteProviderOAuthMutation();
  const stopOAuthMutation = useStopProviderOAuthMutation();

  const nativeOAuthEntries = useMemo(() => allEntries.filter((e) => e.isNativeOAuth), [allEntries]);

  const [activeSession, setActiveSession] = useState<{
    provider: string;
    session: OAuthStartResponse;
  } | null>(null);
  const [step, setStep] = useState<1 | 2>(1);
  const [showManualFields, setShowManualFields] = useState(false);
  const [showAdvancedOAuth, setShowAdvancedOAuth] = useState(false);
  const [callbackDetected, setCallbackDetected] = useState(false);
  const [copied, setCopied] = useState(false);

  // Amazon Quick (aws-quick) two-hop loopback + paste-back wizard state
  // (learned from xrouter OAuthModal + omniroute OAuthWaitingStep).
  const [awsQuickSession, setAwsQuickSession] = useState<OAuthStartResponse | null>(null);
  const [awsQuickPasteStep, setAwsQuickPasteStep] = useState<"discovery" | "code">("discovery");
  const [awsQuickPasteUrl, setAwsQuickPasteUrl] = useState("");
  const [awsQuickError, setAwsQuickError] = useState<string | null>(null);
  const [awsQuickBusy, setAwsQuickBusy] = useState(false);
  const awsQuickStartedRef = useRef(false);
  const awsQuickSettledRef = useRef(false);
  const awsQuickAutoStartedRef = useRef(false);

  // Device Code Flow state (github-copilot, kiro, freepi — learned from xrouter OAuthModal).
  const [deviceSession, setDeviceSession] = useState<{
    provider: string;
    session: OAuthStartResponse;
  } | null>(null);
  const [deviceError, setDeviceError] = useState<string | null>(null);
  const [copiedCode, setCopiedCode] = useState(false);
  const deviceSettledRef = useRef(false);

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

  const stopOAuthMutate = stopOAuthMutation.mutate;
  const startOAuthMutate = startOAuthMutation.mutate;
  const exchangeOAuthMutateAsync = exchangeOAuthMutation.mutateAsync;
  const pasteOAuthMutateAsync = pasteOAuthMutation.mutateAsync;
  const startOAuthReset = startOAuthMutation.reset;
  const exchangeOAuthReset = exchangeOAuthMutation.reset;
  const pasteOAuthReset = pasteOAuthMutation.reset;

  const stopAwsQuickIfOwned = useCallback(() => {
    if (awsQuickStartedRef.current) {
      awsQuickStartedRef.current = false;
      stopOAuthMutate("aws-quick");
    }
  }, [stopOAuthMutate]);

  const handleClose = useCallback(() => {
    stopAwsQuickIfOwned();
    onClose();
  }, [onClose, stopAwsQuickIfOwned]);

  const completeAwsQuickExchange = useCallback(
    async (state: string) => {
      if (awsQuickSettledRef.current) return;
      awsQuickSettledRef.current = true;
      setAwsQuickBusy(true);
      try {
        const connName = startForm.getValues("connection_name")?.trim();
        const conn = await exchangeOAuthMutateAsync({
          provider: "aws-quick",
          body: {
            state,
            priority: DEFAULT_PRIORITY,
            enabled: true,
            ...(connName ? { connection_name: connName } : {}),
          },
        });
        awsQuickStartedRef.current = false;
        onConnected(
          isId
            ? `Akun Amazon Quick "${conn.name}" berhasil dihubungkan.`
            : `Connected Amazon Quick account "${conn.name}".`,
        );
        onClose();
      } catch (err) {
        awsQuickSettledRef.current = false;
        setAwsQuickError(
          err instanceof Error ? err.message : "Failed to complete Amazon Quick OAuth",
        );
      } finally {
        setAwsQuickBusy(false);
      }
    },
    [exchangeOAuthMutateAsync, isId, onClose, onConnected, startForm],
  );

  const startAwsQuickFlow = useCallback(
    (openPopup = true) => {
      setAwsQuickPasteStep("discovery");
      setAwsQuickPasteUrl("");
      setAwsQuickError(null);
      setAwsQuickBusy(false);
      awsQuickSettledRef.current = false;
      const connName = startForm.getValues("connection_name")?.trim();
      startOAuthMutate(
        {
          provider: "aws-quick",
          body: connName ? { connection_name: connName } : {},
        },
        {
          onSuccess: (session) => {
            awsQuickStartedRef.current = true;
            setAwsQuickSession(session);
            if (openPopup && typeof window !== "undefined") {
              window.open(session.authorize_url, "oauth_popup", "width=600,height=700");
            }
          },
          onError: (err) => {
            setAwsQuickError(err.message);
          },
        },
      );
    },
    [startForm, startOAuthMutate],
  );

  const submitAwsQuickPaste = useCallback(
    async (rawPasted: string) => {
      const pasted = rawPasted.trim();
      if (!pasted) return;
      setAwsQuickBusy(true);
      setAwsQuickError(null);
      try {
        const res = await pasteOAuthMutateAsync({
          provider: "aws-quick",
          body: { url: pasted },
        });
        if (res.status === "discovery_saved") {
          setAwsQuickPasteStep("code");
          setAwsQuickPasteUrl("");
          if (res.authorize_url) {
            setAwsQuickSession((prev) =>
              prev
                ? {
                    ...prev,
                    state: res.state,
                    authorize_url: res.authorize_url ?? prev.authorize_url,
                  }
                : {
                    state: res.state,
                    authorize_url: res.authorize_url ?? "",
                    redirect_uri: res.redirect_uri ?? "",
                    expires_at: "",
                  },
            );
            if (typeof window !== "undefined") {
              window.open(res.authorize_url, "oauth_popup", "width=600,height=700");
            }
          }
          return;
        }
        if (res.status === "done") {
          await completeAwsQuickExchange(res.state);
          return;
        }
        setAwsQuickError(
          isId ? "Respons paste-back tidak dikenali." : "Unexpected paste-back response.",
        );
      } catch (err) {
        setAwsQuickError(err instanceof Error ? err.message : "Paste-back failed");
      } finally {
        setAwsQuickBusy(false);
      }
    },
    [completeAwsQuickExchange, isId, pasteOAuthMutateAsync],
  );

  const pasteAwsQuickFromClipboard = useCallback(async () => {
    setAwsQuickBusy(true);
    setAwsQuickError(null);
    try {
      const text = (await navigator.clipboard.readText()).trim();
      if (!text) return;
      setAwsQuickPasteUrl(text);
      await submitAwsQuickPaste(text);
    } catch {
      setAwsQuickError(
        isId
          ? "Akses clipboard diblokir — klik kolom di atas lalu tempel dengan Ctrl/Cmd+V."
          : "Clipboard access was blocked — click the field and paste with Ctrl/Cmd+V.",
      );
    } finally {
      setAwsQuickBusy(false);
    }
  }, [isId, submitAwsQuickPaste]);

  useEffect(() => {
    if (!open) {
      awsQuickAutoStartedRef.current = false;
      stopAwsQuickIfOwned();
      return;
    }
    startOAuthReset();
    exchangeOAuthReset();
    pasteOAuthReset();
    setActiveSession(null);
    setStep(1);
    setShowManualFields(false);
    setShowAdvancedOAuth(false);
    setCallbackDetected(false);
    setCopied(false);
    setAwsQuickSession(null);
    setAwsQuickPasteStep("discovery");
    setAwsQuickPasteUrl("");
    setAwsQuickError(null);
    setAwsQuickBusy(false);
    awsQuickSettledRef.current = false;
    setDeviceSession(null);
    setDeviceError(null);
    setCopiedCode(false);
    deviceSettledRef.current = false;

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
    startOAuthReset,
    exchangeOAuthReset,
    pasteOAuthReset,
    stopAwsQuickIfOwned,
  ]);

  const watchedProvider = startForm.watch("provider");
  const currentEntry = allEntries.find((e) => e.id === watchedProvider) ?? selectedEntry;
  const isNonNativeOAuth = Boolean(selectedEntry && !selectedEntry.isNativeOAuth);
  const isAwsQuick = watchedProvider === "aws-quick";
  const isDeviceCodeProvider =
    watchedProvider === "github-copilot" ||
    watchedProvider === "kiro" ||
    watchedProvider === "freepi";
  const oauthMeta = currentEntry?.adapter?.oauth;
  const clientIdConfigured = Boolean(oauthMeta?.client_id_configured);
  const oauthEnvVarName = formatOAuthClientIdEnvVar(watchedProvider || currentEntry?.id || "");
  const defaultRedirectPlaceholder =
    oauthMeta?.default_redirect_uri || "http://localhost:20128/api/oauth/callback";
  const awsQuickPort = useMemo(() => {
    const uri = awsQuickSession?.redirect_uri;
    if (!uri) return "<port>";
    try {
      return new URL(uri).port || "<port>";
    } catch {
      return "<port>";
    }
  }, [awsQuickSession?.redirect_uri]);

  // Auto-start Amazon Quick loopback proxy when opened directly on aws-quick
  // (matches xrouter OAuthModal zero-click startup).
  useEffect(() => {
    if (!open || !isAwsQuick || awsQuickSession || awsQuickAutoStartedRef.current) return;
    if (selectedEntry?.id !== "aws-quick") return;
    awsQuickAutoStartedRef.current = true;
    startAwsQuickFlow(false);
  }, [open, isAwsQuick, awsQuickSession, selectedEntry?.id, startAwsQuickFlow]);

  // Poll Amazon Quick loopback session status every 1.5s until settled.
  useEffect(() => {
    if (!open || !isAwsQuick || !awsQuickSession?.state) return;
    if (awsQuickSettledRef.current) return;

    let cancelled = false;
    let attempts = 0;
    let timerId: number | null = null;
    const POLL_INTERVAL_MS = 1500;
    const MAX_ATTEMPTS = 200;
    const sessionState = awsQuickSession.state;

    const tick = async () => {
      if (cancelled || awsQuickSettledRef.current) return;
      attempts += 1;
      try {
        const data = await fetchProviderOAuthStatus("aws-quick", sessionState);
        if (cancelled || awsQuickSettledRef.current) return;
        if (data.status === "done") {
          await completeAwsQuickExchange(sessionState);
          return;
        }
        if (data.status === "discovery_saved") {
          setAwsQuickPasteStep("code");
          if (data.authorize_url) {
            setAwsQuickSession((prev) =>
              prev ? { ...prev, authorize_url: data.authorize_url ?? prev.authorize_url } : prev,
            );
          }
        }
        if (data.status === "error") {
          setAwsQuickError(
            data.error ||
              (isId ? "Autentikasi Amazon Quick gagal." : "Amazon Quick authentication failed."),
          );
          return;
        }
      } catch {
        // Transient network error — keep polling.
      }
      if (attempts >= MAX_ATTEMPTS) {
        setAwsQuickError(
          isId
            ? "Waktu tunggu otorisasi habis. Silakan mulai ulang."
            : "Authorization timed out. Please try again.",
        );
        return;
      }
      timerId = window.setTimeout(tick, POLL_INTERVAL_MS);
    };

    timerId = window.setTimeout(tick, POLL_INTERVAL_MS);
    return () => {
      cancelled = true;
      if (timerId !== null) window.clearTimeout(timerId);
    };
  }, [open, isAwsQuick, awsQuickSession?.state, completeAwsQuickExchange, isId]);

  // Poll Device Code Flow status (github-copilot, kiro, freepi) until approved.
  useEffect(() => {
    if (!open || !deviceSession?.session.state) return;
    if (deviceSettledRef.current) return;

    let cancelled = false;
    let timerId: number | null = null;
    let intervalSec = Math.max(deviceSession.session.interval ?? 5, 2);
    const providerId = deviceSession.provider;
    const sessionState = deviceSession.session.state;

    const tick = async () => {
      if (cancelled || deviceSettledRef.current) return;
      try {
        const data = await fetchProviderOAuthStatus(providerId, sessionState);
        if (cancelled || deviceSettledRef.current) return;
        if (data.status === "done") {
          deviceSettledRef.current = true;
          const connName = startForm.getValues("connection_name")?.trim();
          const conn = await exchangeOAuthMutateAsync({
            provider: providerId,
            body: {
              state: sessionState,
              priority: DEFAULT_PRIORITY,
              enabled: true,
              ...(connName ? { connection_name: connName } : {}),
            },
          });
          onConnected(
            isId
              ? `Akun OAuth "${conn.name}" berhasil dihubungkan.`
              : `Connected OAuth account "${conn.name}".`,
          );
          onClose();
          return;
        }
        if (data.status === "error") {
          setDeviceError(
            data.error || (isId ? "Otorisasi perangkat gagal." : "Device authorization failed."),
          );
          return;
        }
        if (data.interval && data.interval > intervalSec) {
          intervalSec = data.interval;
        }
      } catch (err) {
        if (deviceSettledRef.current) {
          deviceSettledRef.current = false;
          setDeviceError(err instanceof Error ? err.message : "OAuth exchange failed");
          return;
        }
        // Transient poll error — keep polling.
      }
      timerId = window.setTimeout(tick, intervalSec * 1000);
    };

    timerId = window.setTimeout(tick, intervalSec * 1000);
    return () => {
      cancelled = true;
      if (timerId !== null) window.clearTimeout(timerId);
    };
  }, [open, deviceSession, exchangeOAuthMutateAsync, isId, onClose, onConnected, startForm]);

  const startDeviceCodeFlow = useCallback(() => {
    const providerId = startForm.getValues("provider");
    const connName = startForm.getValues("connection_name")?.trim();
    setDeviceError(null);
    setCopiedCode(false);
    deviceSettledRef.current = false;
    startOAuthMutate(
      {
        provider: providerId,
        body: connName ? { connection_name: connName } : {},
      },
      {
        onSuccess: (session) => {
          setDeviceSession({ provider: providerId, session });
          const target = session.authorize_url || session.verification_uri;
          if (target && typeof window !== "undefined") {
            window.open(target, "_blank", "noopener,noreferrer");
          }
        },
        onError: (err) => {
          setDeviceError(err.message);
        },
      },
    );
  }, [startForm, startOAuthMutate]);

  const handleCopyUserCode = async () => {
    const code = deviceSession?.session.user_code;
    if (!code) return;
    try {
      await navigator.clipboard.writeText(code);
      setCopiedCode(true);
      window.setTimeout(() => setCopiedCode(false), 2000);
    } catch {
      // Ignore clipboard failures
    }
  };

  const onStart = (values: OAuthStartFormOutput) => {
    if (!clientIdConfigured && !values.client_id) {
      startForm.setError("client_id", {
        type: "manual",
        message: isId
          ? `Client ID diperlukan karena ${oauthEnvVarName} belum disetel di server.`
          : `Client ID is required because ${oauthEnvVarName} is not configured on the server.`,
      });
      return;
    }
    startOAuthMutation.mutate(
      { provider: values.provider, body: toOAuthStartRequest(values) },
      {
        onSuccess: (session) => {
          setActiveSession({ provider: values.provider, session });
          setStep(2);
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

  const handleCallbackUrlChange = (rawValue: string) => {
    const parsed = parseOAuthCallbackUrl(rawValue);
    setCallbackDetected(parsed.valid);
    if (parsed.code) {
      exchangeForm.setValue("code", parsed.code, { shouldDirty: true, shouldValidate: true });
    }
    if (parsed.state) {
      exchangeForm.setValue("state", parsed.state, { shouldDirty: true, shouldValidate: true });
    }
  };

  const handleCopyUrl = async () => {
    const targetUrl = isAwsQuick
      ? awsQuickSession?.authorize_url
      : isDeviceCodeProvider
        ? (deviceSession?.session.authorize_url ?? deviceSession?.session.verification_uri)
        : activeSession?.session.authorize_url;
    if (!targetUrl) return;
    try {
      await navigator.clipboard.writeText(targetUrl);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      // Ignore clipboard failures
    }
  };

  const portalUrl = currentEntry?.signupUrl ?? currentEntry?.website;
  const hasManualErrors = Boolean(
    exchangeForm.formState.errors.code || exchangeForm.formState.errors.state,
  );

  return (
    <ModalShell
      open={open}
      onClose={handleClose}
      size="lg"
      title={
        currentEntry
          ? isId
            ? `Hubungkan ${currentEntry.name}`
            : `Connect ${currentEntry.name}`
          : isId
            ? "Hubungkan Akun OAuth"
            : "Connect OAuth Account"
      }
      subtitle={
        isAwsQuick
          ? isId
            ? "Alur loopback 2-hop Amazon Quick (registrasi OIDC otomatis + fallback paste URL)."
            : "Amazon Quick two-hop loopback flow (automatic OIDC registration + paste-back fallback)."
          : isDeviceCodeProvider
            ? isId
              ? "Alur Device Code — masukkan kode di browser, koneksi tersimpan otomatis."
              : "Device Code flow — enter the code in your browser and the connection saves automatically."
            : isId
              ? "Alur PKCE 2 langkah untuk autentikasi akun OAuth upstream."
              : "2-step PKCE flow for upstream OAuth account authentication."
      }
      badge={
        <Badge variant="outline" className="font-mono text-[10px] uppercase">
          {isAwsQuick ? "LOOPBACK OIDC" : isDeviceCodeProvider ? "DEVICE CODE" : "OAUTH PKCE"}
        </Badge>
      }
      icon={
        currentEntry ? (
          <ProviderLogo
            providerId={currentEntry.id}
            name={currentEntry.name}
            protocol={currentEntry.protocol}
            size={24}
          />
        ) : (
          <ShieldCheck className="h-4 w-4" aria-hidden="true" />
        )
      }
      footer={
        isNonNativeOAuth ? (
          <div className="flex items-center justify-end gap-2">
            <Button type="button" variant="secondary" size="sm" onClick={handleClose}>
              {isId ? "Tutup" : "Close"}
            </Button>
          </div>
        ) : isAwsQuick ? (
          <div className="flex flex-wrap items-center justify-between gap-2">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              disabled={startOAuthMutation.isPending || awsQuickBusy}
              onClick={() => startAwsQuickFlow(true)}
            >
              <RefreshCw
                className={cn("h-3.5 w-3.5", startOAuthMutation.isPending && "animate-spin")}
                aria-hidden="true"
              />
              <span>
                {awsQuickSession
                  ? isId
                    ? "Mulai Ulang Login"
                    : "Restart Sign-In"
                  : isId
                    ? "Mulai Login Amazon Quick"
                    : "Start Amazon Quick Sign-In"}
              </span>
            </Button>

            <Button type="button" variant="secondary" size="sm" onClick={handleClose}>
              {isId ? "Tutup" : "Close"}
            </Button>
          </div>
        ) : isDeviceCodeProvider ? (
          <div className="flex flex-wrap items-center justify-between gap-2">
            <Button
              type="button"
              variant={deviceSession ? "secondary" : "default"}
              size="sm"
              disabled={startOAuthMutation.isPending}
              onClick={startDeviceCodeFlow}
            >
              <RefreshCw
                className={cn("h-3.5 w-3.5", startOAuthMutation.isPending && "animate-spin")}
                aria-hidden="true"
              />
              <span>
                {deviceSession
                  ? isId
                    ? "Minta Kode Baru"
                    : "Request New Code"
                  : isId
                    ? "Mulai Device Login"
                    : "Start Device Login"}
              </span>
            </Button>

            <Button type="button" variant="secondary" size="sm" onClick={handleClose}>
              {isId ? "Tutup" : "Close"}
            </Button>
          </div>
        ) : (
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div>
              {step === 1 ? (
                <button
                  type="button"
                  onClick={() => setStep(2)}
                  className="cursor-pointer text-xs font-medium text-[var(--color-text-muted)] underline-offset-2 hover:text-[var(--color-primary)] hover:underline"
                >
                  {isId ? "Sudah punya callback URL?" : "Already have a callback URL?"}
                </button>
              ) : (
                <Button type="button" variant="secondary" size="sm" onClick={() => setStep(1)}>
                  {isId ? "← Kembali ke Langkah 1" : "← Back to Step 1"}
                </Button>
              )}
            </div>

            <div className="flex items-center gap-2">
              <Button type="button" variant="secondary" size="sm" onClick={handleClose}>
                {isId ? "Tutup" : "Close"}
              </Button>
              {step === 1 ? (
                <Button
                  type="submit"
                  form="oauth-start-form"
                  size="sm"
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
              ) : (
                <Button
                  type="submit"
                  form="oauth-exchange-form"
                  size="sm"
                  disabled={exchangeOAuthMutation.isPending}
                >
                  {exchangeOAuthMutation.isPending
                    ? isId
                      ? "Menghubungkan..."
                      : "Connecting..."
                    : isId
                      ? "Selesaikan OAuth"
                      : "Complete OAuth"}
                </Button>
              )}
            </div>
          </div>
        )
      }
    >
      <div className="space-y-4">
        {isNonNativeOAuth && selectedEntry ? (
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
                handleClose();
                onSwitchToApiKey(selectedEntry.id);
              }}
            >
              <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{isId ? "Setup via API Key / Token" : "Setup via API Key / Token"}</span>
            </Button>
          </div>
        ) : isAwsQuick ? (
          <div className="space-y-4">
            <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label htmlFor="modal-oauth-provider-quick">OAuth Provider</Label>
                <select
                  id="modal-oauth-provider-quick"
                  className={SELECT_CLASS_NAME}
                  {...startForm.register("provider", {
                    onChange: (e) => {
                      if (e.target.value !== "aws-quick") {
                        stopAwsQuickIfOwned();
                        setAwsQuickSession(null);
                        setAwsQuickError(null);
                      }
                    },
                  })}
                >
                  {nativeOAuthEntries.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.name} ({item.id})
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="modal-oauth-name-quick">
                  {isId ? "Label Koneksi (opsional)" : "Connection Label (optional)"}
                </Label>
                <Input
                  id="modal-oauth-name-quick"
                  placeholder="Amazon Quick Account"
                  {...startForm.register("connection_name")}
                />
              </div>
            </div>

            {awsQuickSession ? (
              <>
                {/* Waiting status banner (omniroute + xrouter parity) */}
                <div className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-[var(--color-primary)]/35 bg-[var(--color-surface-2)]/65 px-3.5 py-2.5">
                  <div className="flex items-center gap-2">
                    <RefreshCw
                      className="h-4 w-4 animate-spin text-[var(--color-primary)]"
                      aria-hidden="true"
                    />
                    <span className="text-xs font-medium text-[var(--color-text-main)]">
                      {awsQuickBusy
                        ? isId
                          ? "Memproses token Amazon Quick..."
                          : "Exchanging Amazon Quick tokens..."
                        : isId
                          ? "Menunggu otorisasi browser..."
                          : "Waiting for browser authorization…"}
                    </span>
                  </div>

                  <div className="flex items-center gap-1.5">
                    <Button
                      type="button"
                      size="sm"
                      variant="secondary"
                      onClick={() => void handleCopyUrl()}
                    >
                      {copied ? (
                        <>
                          <Check
                            className="h-3.5 w-3.5 text-[var(--status-ok)]"
                            aria-hidden="true"
                          />
                          <span>{isId ? "Tersalin" : "Copied"}</span>
                        </>
                      ) : (
                        <>
                          <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                          <span>{isId ? "Salin URL" : "Copy URL"}</span>
                        </>
                      )}
                    </Button>
                    <a
                      href={awsQuickSession.authorize_url}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex h-8 items-center gap-1.5 rounded-lg bg-[var(--color-primary)] px-3 text-xs font-medium text-[var(--color-on-primary)] transition-opacity hover:opacity-90"
                    >
                      <span>{isId ? "Buka Halaman Login" : "Open Sign-In Page"}</span>
                      <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                    </a>
                  </div>
                </div>

                <p className="break-all rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface)] p-2 font-mono text-[11px] text-[var(--color-text-muted)]">
                  {awsQuickSession.authorize_url}
                </p>

                {/* Two-step paste-back wizard (exact xrouter layout for remote browser / SSH) */}
                <div className="space-y-3 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35 p-3.5">
                  <div className="flex items-center gap-3">
                    <div className="h-px flex-1 bg-[var(--color-border)]" />
                    <span className="text-[10px] font-semibold uppercase tracking-wider text-[var(--color-text-muted)]">
                      {isId
                        ? "Browser di mesin lain? Tempel kedua URL"
                        : "Browser on another machine? Paste both URLs"}
                    </span>
                    <div className="h-px flex-1 bg-[var(--color-border)]" />
                  </div>

                  <div className="flex items-center gap-2 text-[11px]">
                    <span
                      className={cn(
                        "rounded-full border px-2.5 py-0.5 font-medium",
                        awsQuickPasteStep === "discovery"
                          ? "border-[var(--color-primary)] bg-[var(--color-primary)]/10 text-[var(--color-primary)]"
                          : "border-[var(--color-border)] text-[var(--color-text-muted)]",
                      )}
                    >
                      {isId ? "Langkah 1 · sign-in" : "Step 1 · sign-in"}
                    </span>
                    <span className="text-[var(--color-text-muted)]">→</span>
                    <span
                      className={cn(
                        "rounded-full border px-2.5 py-0.5 font-medium",
                        awsQuickPasteStep === "code"
                          ? "border-[var(--color-primary)] bg-[var(--color-primary)]/10 text-[var(--color-primary)]"
                          : "border-[var(--color-border)] text-[var(--color-text-muted)]",
                      )}
                    >
                      {isId ? "Langkah 2 · otorisasi" : "Step 2 · authorize"}
                    </span>
                  </div>

                  <p className="text-xs leading-relaxed text-[var(--color-text-muted)]">
                    {awsQuickPasteStep === "discovery"
                      ? isId
                        ? "Selesaikan login AWS sampai browser menampilkan error koneksi ke 127.0.0.1. Salin seluruh URL di address bar (mengandung nonce=…) lalu tempel di bawah. Halaman otorisasi akan terbuka berikutnya."
                        : "Finish the AWS sign-in until the browser shows a connection error at 127.0.0.1. Copy that whole address bar URL (it contains nonce=…) and paste it below. The authorize page opens next."
                      : isId
                        ? "Setujui akses pada halaman otorisasi. Browser akan gagal lagi di 127.0.0.1 — salin seluruh URL di address bar (mengandung code=…) lalu tempel di bawah untuk menyelesaikan."
                        : "Approve access on the authorize page. The browser then fails again at 127.0.0.1 — copy that whole address bar URL (it contains code=…) and paste it below to finish."}
                  </p>

                  <Input
                    value={awsQuickPasteUrl}
                    onChange={(e) => setAwsQuickPasteUrl(e.target.value)}
                    placeholder={
                      awsQuickPasteStep === "discovery"
                        ? `http://127.0.0.1:${awsQuickPort}/callback?nonce=…&issuerUrl=…`
                        : `http://127.0.0.1:${awsQuickPort}/?code=…&state=…`
                    }
                    className="font-mono text-xs"
                  />

                  {awsQuickError ? (
                    <Alert variant="destructive">
                      <AlertDescription>{awsQuickError}</AlertDescription>
                    </Alert>
                  ) : null}

                  {awsQuickPasteUrl &&
                  /[?&]code=/.test(awsQuickPasteUrl) &&
                  awsQuickPasteStep === "discovery" ? (
                    <p className="text-[11px] text-[var(--status-warn)]">
                      {isId
                        ? "URL ini terlihat seperti URL final (code=…). Jika login tahap 1 sudah lewat mesin ini, klik Lanjutkan untuk langsung menyelesaikan."
                        : "This looks like the final URL (code=…). If the sign-in page was reached on this machine the server may already have the tenant."}
                    </p>
                  ) : null}

                  <div className="grid grid-cols-2 gap-2">
                    <Button
                      type="button"
                      variant="secondary"
                      size="sm"
                      disabled={awsQuickBusy}
                      onClick={() => void pasteAwsQuickFromClipboard()}
                    >
                      {isId ? "Tempel dari Clipboard" : "Paste from Clipboard"}
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      disabled={awsQuickBusy || !awsQuickPasteUrl.trim()}
                      onClick={() => void submitAwsQuickPaste(awsQuickPasteUrl)}
                    >
                      {awsQuickBusy
                        ? isId
                          ? "Memproses..."
                          : "Working..."
                        : awsQuickPasteStep === "discovery"
                          ? isId
                            ? "Lanjutkan"
                            : "Continue"
                          : isId
                            ? "Hubungkan"
                            : "Connect"}
                    </Button>
                  </div>

                  <p className="text-[11px] text-[var(--color-text-muted)]">
                    {isId
                      ? "Ingin tanpa tempel manual? Jalankan dashboard di mesin ini atau buka SSH tunnel:"
                      : "Want zero pastes? Run the dashboard on this machine, or open an SSH tunnel:"}{" "}
                    <span className="font-mono text-[var(--color-text-main)]">
                      ssh -L {awsQuickPort}:127.0.0.1:{awsQuickPort} &lt;this-host&gt;
                    </span>
                  </p>
                </div>
              </>
            ) : (
              <div className="space-y-3 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 p-4 text-center">
                {awsQuickError ? (
                  <Alert variant="destructive">
                    <AlertDescription>{awsQuickError}</AlertDescription>
                  </Alert>
                ) : null}
                <p className="text-xs leading-relaxed text-[var(--color-text-muted)]">
                  {startOAuthMutation.isPending
                    ? isId
                      ? "Menyiapkan listener loopback Amazon Quick..."
                      : "Starting Amazon Quick loopback callback server..."
                    : isId
                      ? "Klik tombol di bawah untuk memulai login Amazon Quick. Client ID OIDC akan didaftarkan otomatis tanpa perlu konfigurasi manual."
                      : "Click below to start Amazon Quick sign-in. OIDC Client ID is registered dynamically — no manual Client ID required."}
                </p>
                <Button
                  type="button"
                  size="sm"
                  disabled={startOAuthMutation.isPending}
                  onClick={() => startAwsQuickFlow(true)}
                >
                  {startOAuthMutation.isPending ? (
                    <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                  ) : (
                    <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                  )}
                  <span>{isId ? "Mulai Login Amazon Quick" : "Start Amazon Quick Sign-In"}</span>
                </Button>
              </div>
            )}
          </div>
        ) : isDeviceCodeProvider ? (
          <div className="space-y-4">
            <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-2">
              <div className="space-y-1.5">
                <Label htmlFor="modal-oauth-provider-device">OAuth Provider</Label>
                <select
                  id="modal-oauth-provider-device"
                  className={SELECT_CLASS_NAME}
                  {...startForm.register("provider", {
                    onChange: () => {
                      setDeviceSession(null);
                      setDeviceError(null);
                      deviceSettledRef.current = false;
                    },
                  })}
                >
                  {nativeOAuthEntries.map((item) => (
                    <option key={item.id} value={item.id}>
                      {item.name} ({item.id})
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1.5">
                <Label htmlFor="modal-oauth-name-device">
                  {isId ? "Label Koneksi (opsional)" : "Connection Label (optional)"}
                </Label>
                <Input
                  id="modal-oauth-name-device"
                  placeholder={currentEntry ? `${currentEntry.name} Account` : "Primary Account"}
                  {...startForm.register("connection_name")}
                />
              </div>
            </div>

            {deviceError ? (
              <Alert variant="destructive">
                <AlertDescription>{deviceError}</AlertDescription>
              </Alert>
            ) : null}

            {deviceSession && deviceSession.provider === watchedProvider ? (
              <div className="space-y-3.5">
                <div className="space-y-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/55 p-3.5">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <span className="text-xs font-semibold text-[var(--color-text-main)]">
                      {isId ? "1. Buka URL Verifikasi" : "1. Open Verification URL"}
                    </span>
                    <div className="flex items-center gap-1.5">
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        onClick={() => void handleCopyUrl()}
                      >
                        {copied ? (
                          <>
                            <Check
                              className="h-3.5 w-3.5 text-[var(--status-ok)]"
                              aria-hidden="true"
                            />
                            <span>{isId ? "Tersalin" : "Copied"}</span>
                          </>
                        ) : (
                          <>
                            <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                            <span>{isId ? "Salin" : "Copy"}</span>
                          </>
                        )}
                      </Button>
                      <a
                        href={
                          deviceSession.session.authorize_url ||
                          deviceSession.session.verification_uri
                        }
                        target="_blank"
                        rel="noreferrer"
                        className="inline-flex h-8 items-center gap-1.5 rounded-lg bg-[var(--color-primary)] px-3 text-xs font-medium text-[var(--color-on-primary)] transition-opacity hover:opacity-90"
                      >
                        <span>{isId ? "Buka" : "Open"}</span>
                        <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                      </a>
                    </div>
                  </div>
                  <p className="break-all rounded-lg bg-[var(--color-surface)] p-2 font-mono text-[11px] text-[var(--color-text-muted)]">
                    {deviceSession.session.authorize_url || deviceSession.session.verification_uri}
                  </p>
                </div>

                {deviceSession.session.user_code ? (
                  <div className="rounded-xl border border-[var(--color-primary)]/35 bg-[var(--color-primary)]/10 p-4 text-center">
                    <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--color-text-muted)]">
                      {isId ? "2. Masukkan Kode Perangkat Anda" : "2. Enter Your Device Code"}
                    </p>
                    <div className="mt-2 flex items-center justify-center gap-2.5">
                      <span className="font-mono text-2xl font-bold tracking-widest text-[var(--color-primary)]">
                        {deviceSession.session.user_code}
                      </span>
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        onClick={() => void handleCopyUserCode()}
                      >
                        {copiedCode ? (
                          <Check
                            className="h-3.5 w-3.5 text-[var(--status-ok)]"
                            aria-hidden="true"
                          />
                        ) : (
                          <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                        )}
                      </Button>
                    </div>
                  </div>
                ) : null}

                <div className="flex items-center justify-center gap-2 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 py-2.5 text-xs text-[var(--color-text-muted)]">
                  <RefreshCw
                    className="h-3.5 w-3.5 animate-spin text-[var(--color-primary)]"
                    aria-hidden="true"
                  />
                  <span>
                    {isId
                      ? "Menunggu konfirmasi otorisasi di browser..."
                      : "Waiting for authorization in browser…"}
                  </span>
                </div>
              </div>
            ) : (
              <div className="space-y-3 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/45 p-4 text-center">
                <p className="text-xs leading-relaxed text-[var(--color-text-muted)]">
                  {isId
                    ? `Klik "Mulai Device Login" untuk mendapatkan kode perangkat ${currentEntry?.name ?? watchedProvider}. Setelah disetujui di browser, akun otomatis tersambung.`
                    : `Click "Start Device Login" to get a ${currentEntry?.name ?? watchedProvider} device code. Once approved in your browser, the account connects automatically.`}
                </p>
                <Button
                  type="button"
                  size="sm"
                  disabled={startOAuthMutation.isPending}
                  onClick={startDeviceCodeFlow}
                >
                  {startOAuthMutation.isPending ? (
                    <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                  ) : (
                    <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                  )}
                  <span>{isId ? "Mulai Device Login" : "Start Device Login"}</span>
                </Button>
              </div>
            )}
          </div>
        ) : (
          <>
            {/* 2-Step Visual Stepper */}
            <div className="grid grid-cols-2 gap-2">
              <button
                type="button"
                onClick={() => setStep(1)}
                className={cn(
                  "flex cursor-pointer items-center gap-2.5 rounded-xl border px-3 py-2.5 text-left transition-colors",
                  step === 1
                    ? "border-[var(--color-primary)]/50 bg-[var(--color-primary)]/10"
                    : "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 hover:border-[var(--color-border)]",
                )}
              >
                <span
                  className={cn(
                    "flex h-6 w-6 shrink-0 items-center justify-center rounded-full font-mono text-[11px] font-bold",
                    activeSession
                      ? "bg-[var(--status-ok)] text-white"
                      : step === 1
                        ? "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                        : "bg-[var(--color-surface-3)] text-[var(--color-text-muted)]",
                  )}
                >
                  {activeSession ? <Check className="h-3.5 w-3.5" aria-hidden="true" /> : "1"}
                </span>
                <div className="min-w-0">
                  <p className="truncate text-xs font-semibold text-[var(--color-text-main)]">
                    {isId ? "Buat Sesi & Login" : "Start & Sign In"}
                  </p>
                  <p className="truncate text-[10px] text-[var(--color-text-muted)]">
                    {isId ? "Generate URL PKCE" : "Generate PKCE URL"}
                  </p>
                </div>
              </button>

              <button
                type="button"
                onClick={() => setStep(2)}
                className={cn(
                  "flex cursor-pointer items-center gap-2.5 rounded-xl border px-3 py-2.5 text-left transition-colors",
                  step === 2
                    ? "border-[var(--color-primary)]/50 bg-[var(--color-primary)]/10"
                    : "border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 hover:border-[var(--color-border)]",
                )}
              >
                <span
                  className={cn(
                    "flex h-6 w-6 shrink-0 items-center justify-center rounded-full font-mono text-[11px] font-bold",
                    step === 2
                      ? "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                      : "bg-[var(--color-surface-3)] text-[var(--color-text-muted)]",
                  )}
                >
                  2
                </span>
                <div className="min-w-0">
                  <p className="truncate text-xs font-semibold text-[var(--color-text-main)]">
                    {isId ? "Verifikasi Callback" : "Verify Callback"}
                  </p>
                  <p className="truncate text-[10px] text-[var(--color-text-muted)]">
                    {isId ? "Tempel URL hasil login" : "Paste redirect URL"}
                  </p>
                </div>
              </button>
            </div>

            {step === 1 ? (
              <div className="space-y-4">
                {currentEntry && (portalUrl || currentEntry.deprecated) ? (
                  <div className="flex flex-wrap items-center justify-between gap-2 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/55 px-3.5 py-2.5">
                    <p className="flex-1 text-[11px] leading-relaxed text-[var(--color-text-muted)]">
                      {isId
                        ? "Login melalui akun pribadi Anda, lalu salin URL callback dari address bar browser setelah redirect."
                        : "Sign in with your personal account, then copy the final callback URL from your browser address bar."}
                    </p>
                    {portalUrl ? (
                      <a
                        href={portalUrl}
                        target="_blank"
                        rel="noreferrer"
                        className="inline-flex shrink-0 items-center gap-1 text-[11px] font-medium text-[var(--color-primary)] hover:underline"
                      >
                        <span>{isId ? "Portal Akun" : "Account Portal"}</span>
                        <ExternalLink className="h-3 w-3" aria-hidden="true" />
                      </a>
                    ) : null}
                  </div>
                ) : null}

                <form
                  id="oauth-start-form"
                  className="space-y-4"
                  onSubmit={startForm.handleSubmit(onStart)}
                  noValidate
                >
                  {startOAuthMutation.isError ? (
                    <Alert variant="destructive">
                      <AlertDescription>{startOAuthMutation.error.message}</AlertDescription>
                    </Alert>
                  ) : null}

                  <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-2">
                    <div className="space-y-1.5">
                      <Label htmlFor="modal-oauth-provider">OAuth Provider</Label>
                      <select
                        id="modal-oauth-provider"
                        className={SELECT_CLASS_NAME}
                        {...startForm.register("provider", {
                          onChange: () => {
                            startOAuthMutation.reset();
                            startForm.clearErrors("client_id");
                          },
                        })}
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
                        placeholder={
                          currentEntry ? `${currentEntry.name} Account` : "Primary Account"
                        }
                        {...startForm.register("connection_name")}
                      />
                    </div>
                  </div>

                  {clientIdConfigured ? (
                    <div className="flex items-center gap-2 rounded-lg border border-[var(--status-ok)]/30 bg-[var(--status-ok)]/10 px-3 py-2 text-[11px] text-[var(--color-text-main)]">
                      <CheckCircle2
                        className="h-3.5 w-3.5 shrink-0 text-[var(--status-ok)]"
                        aria-hidden="true"
                      />
                      <span>
                        {isId
                          ? "Client ID sudah dikonfigurasi di server. Anda dapat langsung membuat URL otorisasi."
                          : "OAuth Client ID is configured on the server. You can generate the authorization URL directly."}
                      </span>
                    </div>
                  ) : (
                    <div className="space-y-1.5">
                      <div className="flex items-center justify-between gap-2">
                        <Label htmlFor="modal-oauth-client-id">
                          {isId ? "OAuth Client ID *" : "OAuth Client ID *"}
                        </Label>
                        <span className="font-mono text-[10px] text-[var(--color-text-muted)]">
                          {oauthEnvVarName}
                        </span>
                      </div>
                      <Input
                        id="modal-oauth-client-id"
                        className="font-mono text-xs"
                        placeholder={
                          isId
                            ? "Masukkan OAuth Client ID aplikasi..."
                            : "Enter OAuth application Client ID..."
                        }
                        {...startForm.register("client_id")}
                      />
                      {startForm.formState.errors.client_id ? (
                        <p className="text-xs text-[var(--status-error)]">
                          {startForm.formState.errors.client_id.message}
                        </p>
                      ) : (
                        <p className="text-[11px] text-[var(--color-text-muted)]">
                          {isId
                            ? `Env var ${oauthEnvVarName} belum disetel di server. Isi Client ID di sini atau setel env var tersebut lalu restart gateway.`
                            : `${oauthEnvVarName} is not set on the server. Provide a Client ID here or set that environment variable and restart the gateway.`}
                        </p>
                      )}
                    </div>
                  )}

                  <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35">
                    <button
                      type="button"
                      onClick={() => setShowAdvancedOAuth((prev) => !prev)}
                      className="flex w-full cursor-pointer items-center justify-between px-3.5 py-2.5 text-left text-xs font-medium text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-text-main)]"
                    >
                      <span className="inline-flex items-center gap-1.5">
                        <Settings2 className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>
                          {isId
                            ? "Opsi OAuth Lanjutan (Client Secret / Redirect URI)"
                            : "Advanced OAuth Options (Client Secret / Redirect URI)"}
                        </span>
                      </span>
                      <ChevronDown
                        className={cn(
                          "h-3.5 w-3.5 transition-transform",
                          showAdvancedOAuth && "rotate-180",
                        )}
                        aria-hidden="true"
                      />
                    </button>

                    {showAdvancedOAuth ? (
                      <div className="space-y-3 border-t border-[var(--color-border-subtle)] px-3.5 py-3">
                        {clientIdConfigured ? (
                          <div className="space-y-1.5">
                            <Label htmlFor="modal-oauth-client-id-override">
                              {isId
                                ? "Override Client ID (opsional)"
                                : "Client ID Override (optional)"}
                            </Label>
                            <Input
                              id="modal-oauth-client-id-override"
                              className="font-mono text-xs"
                              placeholder={
                                isId
                                  ? "Kosongkan untuk memakai default server"
                                  : "Leave blank to use server default"
                              }
                              {...startForm.register("client_id")}
                            />
                          </div>
                        ) : null}

                        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                          <div className="space-y-1.5">
                            <Label htmlFor="modal-oauth-client-secret">
                              {isId ? "Client Secret (opsional)" : "Client Secret (optional)"}
                            </Label>
                            <Input
                              id="modal-oauth-client-secret"
                              type="password"
                              className="font-mono text-xs"
                              placeholder={
                                isId
                                  ? "Opsional untuk public PKCE client"
                                  : "Optional for public PKCE clients"
                              }
                              {...startForm.register("client_secret")}
                            />
                          </div>

                          <div className="space-y-1.5">
                            <Label htmlFor="modal-oauth-redirect-uri">
                              {isId ? "Redirect URI (opsional)" : "Redirect URI (optional)"}
                            </Label>
                            <Input
                              id="modal-oauth-redirect-uri"
                              className="font-mono text-xs"
                              placeholder={defaultRedirectPlaceholder}
                              {...startForm.register("redirect_uri")}
                            />
                            {startForm.formState.errors.redirect_uri ? (
                              <p className="text-xs text-[var(--status-error)]">
                                {startForm.formState.errors.redirect_uri.message}
                              </p>
                            ) : null}
                          </div>
                        </div>
                      </div>
                    ) : null}
                  </div>
                </form>
              </div>
            ) : (
              <div className="space-y-4">
                {activeSession ? (
                  <div className="space-y-2.5 rounded-xl border border-[var(--color-primary)]/35 bg-[var(--color-surface-2)]/65 p-3.5">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <span className="text-xs font-semibold text-[var(--color-text-main)]">
                        {isId
                          ? "1. Buka URL login & selesaikan otorisasi:"
                          : "1. Open login URL & authorize:"}
                      </span>
                      <div className="flex items-center gap-1.5">
                        <Button
                          type="button"
                          size="sm"
                          variant="secondary"
                          onClick={() => void handleCopyUrl()}
                        >
                          {copied ? (
                            <>
                              <Check
                                className="h-3.5 w-3.5 text-[var(--status-ok)]"
                                aria-hidden="true"
                              />
                              <span>{isId ? "Tersalin" : "Copied"}</span>
                            </>
                          ) : (
                            <>
                              <Copy className="h-3.5 w-3.5" aria-hidden="true" />
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
                          <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                        </a>
                      </div>
                    </div>
                    <p className="break-all rounded-lg bg-[var(--color-surface)] p-2 font-mono text-[11px] text-[var(--color-text-muted)]">
                      {activeSession.session.authorize_url}
                    </p>
                  </div>
                ) : null}

                <form
                  id="oauth-exchange-form"
                  className="space-y-3.5"
                  onSubmit={exchangeForm.handleSubmit(onExchange)}
                  noValidate
                >
                  {exchangeOAuthMutation.isError ? (
                    <Alert variant="destructive">
                      <AlertDescription>{exchangeOAuthMutation.error.message}</AlertDescription>
                    </Alert>
                  ) : null}

                  <div className="space-y-1.5">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <Label htmlFor="modal-oauth-callback">
                        {isId ? "2. Tempel Callback URL" : "2. Paste Callback URL"}
                      </Label>
                      {callbackDetected ? (
                        <Badge variant="success" className="text-[10px]">
                          <Check className="mr-1 h-3 w-3" aria-hidden="true" />
                          {isId ? "Code & State terdeteksi" : "Code & State detected"}
                        </Badge>
                      ) : null}
                    </div>
                    <Input
                      id="modal-oauth-callback"
                      placeholder="http://localhost:.../callback?code=...&state=..."
                      className="font-mono text-xs"
                      {...exchangeForm.register("callback_url", {
                        onChange: (e) => handleCallbackUrlChange(e.target.value),
                      })}
                    />
                    {exchangeForm.formState.errors.callback_url ? (
                      <p className="text-xs text-[var(--status-error)]">
                        {exchangeForm.formState.errors.callback_url.message}
                      </p>
                    ) : null}
                  </div>

                  {/* Collapsible Manual Code & State */}
                  <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35">
                    <button
                      type="button"
                      onClick={() => setShowManualFields((v) => !v)}
                      className="flex w-full cursor-pointer items-center justify-between px-3.5 py-2 text-xs font-medium text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]"
                    >
                      <span>
                        {isId
                          ? "Input manual Code & State (opsional)"
                          : "Manual Code & State (optional)"}
                      </span>
                      <ChevronDown
                        className={cn(
                          "h-3.5 w-3.5 transition-transform",
                          (showManualFields || hasManualErrors) && "rotate-180",
                        )}
                        aria-hidden="true"
                      />
                    </button>

                    {showManualFields || hasManualErrors ? (
                      <div className="grid grid-cols-1 gap-3 border-t border-[var(--color-border-subtle)] p-3.5 sm:grid-cols-2">
                        <div className="space-y-1.5">
                          <Label htmlFor="modal-oauth-code">Code</Label>
                          <Input
                            id="modal-oauth-code"
                            placeholder="Authorization code"
                            className="font-mono text-xs"
                            {...exchangeForm.register("code")}
                          />
                          {exchangeForm.formState.errors.code ? (
                            <p className="text-xs text-[var(--status-error)]">
                              {exchangeForm.formState.errors.code.message}
                            </p>
                          ) : null}
                        </div>
                        <div className="space-y-1.5">
                          <Label htmlFor="modal-oauth-state">State</Label>
                          <Input
                            id="modal-oauth-state"
                            placeholder="PKCE state"
                            className="font-mono text-xs"
                            {...exchangeForm.register("state")}
                          />
                          {exchangeForm.formState.errors.state ? (
                            <p className="text-xs text-[var(--status-error)]">
                              {exchangeForm.formState.errors.state.message}
                            </p>
                          ) : null}
                        </div>
                      </div>
                    ) : null}
                  </div>
                </form>
              </div>
            )}
          </>
        )}
      </div>
    </ModalShell>
  );
}

/* ============================================================================
 * 3. Add Custom Compatible Node Modal (Quick Presets + Pre-Save Connectivity Test)
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
  onCreated: (node: ProviderNode, meta?: { connectionCreated: boolean }) => void;
}) {
  const createNodeMutation = useCreateProviderNodeMutation();
  const createConnMutation = useCreateProviderConnectionMutation();
  const testDraftMutation = useTestProviderNodeDraftMutation();

  const [activePresetId, setActivePresetId] = useState<CustomEndpointPreset["id"] | null>(null);
  const [idManuallyEdited, setIdManuallyEdited] = useState(false);
  const [apiKey, setApiKey] = useState("");
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [testState, setTestState] = useState<{
    result: TestProviderConnectionResponse;
    fingerprint: string;
  } | null>(null);
  const [testValidationError, setTestValidationError] = useState<string | null>(null);

  const form = useForm<CreateNodeFormInput, unknown, CreateNodeFormOutput>({
    resolver: zodResolver(createNodeSchema),
    defaultValues: DEFAULT_CREATE_NODE_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const watchedProtocol = form.watch("protocol");
  const watchedBaseUrl = form.watch("base_url");
  const watchedModels = form.watch("models");
  const watchedHeaders = form.watch("headers");

  const currentFingerprint = useMemo(
    () =>
      buildTestFingerprint({
        protocol: watchedProtocol,
        baseUrl: watchedBaseUrl,
        apiKey,
        headers: watchedHeaders,
      }),
    [watchedProtocol, watchedBaseUrl, apiKey, watchedHeaders],
  );

  const isTestStale = Boolean(testState && testState.fingerprint !== currentFingerprint);

  useEffect(() => {
    if (!open) return;
    createNodeMutation.reset();
    createConnMutation.reset();
    testDraftMutation.reset();
    setActivePresetId(null);
    setIdManuallyEdited(false);
    setApiKey("");
    setShowApiKey(false);
    setShowAdvanced(false);
    setSubmitError(null);
    setTestState(null);
    setTestValidationError(null);
    form.reset({
      ...DEFAULT_CREATE_NODE_VALUES,
      protocol: initialProtocol,
    });
  }, [
    open,
    initialProtocol,
    form,
    createNodeMutation.reset,
    createConnMutation.reset,
    testDraftMutation.reset,
  ]);

  const applyPreset = (preset: CustomEndpointPreset) => {
    setActivePresetId(preset.id);
    setIdManuallyEdited(false);
    form.setValue("name", preset.displayName, { shouldDirty: true, shouldValidate: true });
    form.setValue("id", preset.id, { shouldDirty: true, shouldValidate: true });
    form.setValue("protocol", preset.protocol, { shouldDirty: true });
    form.setValue("base_url", preset.baseUrl, { shouldDirty: true, shouldValidate: true });
    form.setValue("models", preset.models, { shouldDirty: true });
    setTestValidationError(null);
  };

  const handleTestEndpoint = async () => {
    setTestValidationError(null);
    const values = form.getValues();
    const baseUrl = values.base_url.trim();
    if (!isValidHttpUrl(baseUrl)) {
      setTestValidationError(
        isId
          ? "Masukkan Base URL http:// atau https:// yang valid terlebih dahulu."
          : "Enter a valid http:// or https:// Base URL before testing.",
      );
      return;
    }

    const parsedHeaders = parseHeadersText(values.headers);
    if (!parsedHeaders.ok) {
      setTestValidationError(parsedHeaders.error);
      setShowAdvanced(true);
      return;
    }

    const fingerprint = buildTestFingerprint({
      protocol: values.protocol,
      baseUrl,
      apiKey,
      headers: values.headers,
    });

    try {
      const result = await testDraftMutation.mutateAsync({
        protocol: values.protocol,
        base_url: baseUrl,
        ...(apiKey.trim() ? { api_key: apiKey.trim() } : {}),
        headers: parsedHeaders.headers,
      });
      setTestState({ result, fingerprint });
    } catch (err) {
      setTestValidationError(err instanceof Error ? err.message : "Connectivity test failed.");
    }
  };

  const onSubmit = async (values: CreateNodeFormOutput) => {
    setSubmitError(null);
    try {
      const created = await createNodeMutation.mutateAsync(toCreateNodeRequest(values));
      const trimmedKey = apiKey.trim();
      if (trimmedKey) {
        try {
          await createConnMutation.mutateAsync({
            provider_id: created.id,
            name: `${created.name} Primary`,
            base_url: created.base_url,
            auth_type: "api_key",
            api_key: trimmedKey,
            priority: 100,
            enabled: true,
            models: created.models,
            headers: created.headers,
          });
          onCreated(created, { connectionCreated: true });
          onClose();
        } catch (connErr) {
          setSubmitError(
            isId
              ? `Endpoint "${created.name}" berhasil didaftarkan, tetapi gagal menyimpan API Key awal: ${connErr instanceof Error ? connErr.message : "error"}. Anda dapat menambahkannya lewat kartu provider.`
              : `Registered endpoint "${created.name}", but failed to save its initial API key: ${connErr instanceof Error ? connErr.message : "error"}. You can add it from the provider card.`,
          );
          onCreated(created, { connectionCreated: false });
        }
        return;
      }

      onCreated(created, { connectionCreated: false });
      onClose();
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Failed to register custom endpoint.");
    }
  };

  const modelCount = parseCommaSeparatedList(watchedModels).length;
  const parsedHeaderPreview = parseHeadersText(watchedHeaders);
  const headerCount = parsedHeaderPreview.ok ? Object.keys(parsedHeaderPreview.headers).length : 0;
  const isSaving = createNodeMutation.isPending || createConnMutation.isPending;
  const hasAdvancedErrors = Boolean(form.formState.errors.headers || form.formState.errors.models);

  return (
    <ModalShell
      open={open}
      onClose={onClose}
      size="lg"
      title={isId ? "Tambah Custom Endpoint" : "Add Custom Compatible Endpoint"}
      subtitle={
        isId
          ? "Daftarkan gateway OpenAI, Anthropic, atau Gemini kompatibel (Ollama, vLLM, LiteLLM, dll)."
          : "Register an OpenAI, Anthropic, or Gemini compatible gateway (Ollama, vLLM, LiteLLM, etc)."
      }
      badge={
        <Badge variant="outline" className="font-mono text-[10px] uppercase">
          {watchedProtocol}
        </Badge>
      }
      icon={<Server className="h-4 w-4" aria-hidden="true" />}
      footer={
        <div className="flex flex-wrap items-center justify-between gap-2">
          <Button
            type="button"
            variant="secondary"
            size="sm"
            disabled={testDraftMutation.isPending || isSaving}
            onClick={() => void handleTestEndpoint()}
          >
            {testDraftMutation.isPending ? (
              <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
            ) : (
              <Zap className="h-3.5 w-3.5" aria-hidden="true" />
            )}
            <span>
              {testDraftMutation.isPending
                ? isId
                  ? "Menguji Endpoint..."
                  : "Testing Endpoint..."
                : isId
                  ? "Tes Koneksi"
                  : "Test Connection"}
            </span>
          </Button>

          <div className="flex items-center gap-2">
            <Button type="button" variant="secondary" size="sm" onClick={onClose}>
              {isId ? "Batal" : "Cancel"}
            </Button>
            <Button type="submit" form="add-custom-node-form" size="sm" disabled={isSaving}>
              {isSaving
                ? isId
                  ? "Menyimpan..."
                  : "Saving..."
                : apiKey.trim()
                  ? isId
                    ? "Daftar & Simpan"
                    : "Register & Save"
                  : isId
                    ? "Daftar Endpoint"
                    : "Register Endpoint"}
            </Button>
          </div>
        </div>
      }
    >
      <form
        id="add-custom-node-form"
        className="space-y-4"
        onSubmit={form.handleSubmit(onSubmit)}
        noValidate
      >
        {submitError ? (
          <Alert variant="destructive">
            <AlertDescription>{submitError}</AlertDescription>
          </Alert>
        ) : null}

        {/* Quick Presets */}
        <div className="space-y-1.5">
          <span className="inline-flex items-center gap-1.5 text-[11px] font-medium text-[var(--color-text-muted)]">
            <Sparkles className="h-3 w-3 text-[var(--color-primary)]" aria-hidden="true" />
            <span>{isId ? "Preset Cepat" : "Quick Presets"}</span>
          </span>
          <div className="flex flex-wrap gap-1.5">
            {CUSTOM_ENDPOINT_PRESETS.map((preset) => {
              const isSelected = activePresetId === preset.id;
              return (
                <button
                  key={preset.id}
                  type="button"
                  onClick={() => applyPreset(preset)}
                  className={cn(
                    "cursor-pointer rounded-lg border px-2.5 py-1 text-xs font-medium transition-colors",
                    isSelected
                      ? "border-[var(--color-primary)]/50 bg-[var(--color-primary)]/12 text-[var(--color-primary)]"
                      : "border-[var(--color-border)] bg-[var(--color-surface-2)]/60 text-[var(--color-text-muted)] hover:border-[var(--color-primary)]/35 hover:text-[var(--color-text-main)]",
                  )}
                >
                  {preset.label}
                </button>
              );
            })}
          </div>
        </div>

        {/* Display Name & Provider ID */}
        <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label htmlFor="modal-node-name">{isId ? "Nama Tampilan" : "Display Name"}</Label>
            <Input
              id="modal-node-name"
              placeholder="Local vLLM Cluster"
              {...form.register("name", {
                onChange: (e) => {
                  setActivePresetId(null);
                  if (!idManuallyEdited) {
                    form.setValue("id", slugifyProviderId(e.target.value), {
                      shouldDirty: false,
                    });
                  }
                },
              })}
            />
            {form.formState.errors.name ? (
              <p className="text-xs text-[var(--status-error)]">
                {form.formState.errors.name.message}
              </p>
            ) : null}
          </div>

          <div className="space-y-1.5">
            <Label htmlFor="modal-node-id">
              {isId ? "Provider ID (otomatis)" : "Provider ID (auto)"}
            </Label>
            <Input
              id="modal-node-id"
              placeholder="local-vllm-cluster"
              className="font-mono text-xs"
              {...form.register("id", {
                onChange: (e) => {
                  setIdManuallyEdited(e.target.value.trim().length > 0);
                },
              })}
            />
            {form.formState.errors.id ? (
              <p className="text-xs text-[var(--status-error)]">
                {form.formState.errors.id.message}
              </p>
            ) : null}
          </div>
        </div>

        {/* Protocol & Base URL */}
        <div className="grid grid-cols-1 gap-3.5 sm:grid-cols-3">
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

        {/* Optional API Key (used for pre-save test + optional initial connection) */}
        <div className="space-y-1.5">
          <Label htmlFor="modal-node-apikey">
            {isId ? "API Key / Token (opsional)" : "API Key / Token (optional)"}
          </Label>
          <div className="relative">
            <Input
              id="modal-node-apikey"
              type={showApiKey ? "text" : "password"}
              autoComplete="off"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder={
                isId
                  ? "Kosongkan untuk endpoint lokal tanpa auth (mis. Ollama)"
                  : "Leave blank for local no-auth endpoints (e.g. Ollama)"
              }
              className="pr-10 font-mono text-xs"
            />
            <button
              type="button"
              onClick={() => setShowApiKey((v) => !v)}
              aria-label={showApiKey ? "Hide API key" : "Show API key"}
              className="absolute top-1/2 right-2 -translate-y-1/2 cursor-pointer rounded p-1 text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]"
            >
              {showApiKey ? (
                <EyeOff className="h-3.5 w-3.5" aria-hidden="true" />
              ) : (
                <Eye className="h-3.5 w-3.5" aria-hidden="true" />
              )}
            </button>
          </div>
          <p className="text-[11px] text-[var(--color-text-muted)]">
            {isId
              ? "Digunakan untuk Tes Koneksi dan otomatis membuat koneksi awal saat disimpan."
              : "Used for Test Connection and automatically creates the initial connection on save."}
          </p>
        </div>

        {/* Collapsible Default Models & Headers */}
        <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35">
          <button
            type="button"
            onClick={() => setShowAdvanced((v) => !v)}
            className="flex w-full cursor-pointer items-center justify-between gap-2 px-3.5 py-2.5 text-xs font-medium text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]"
          >
            <span className="flex items-center gap-2">
              <Settings2 className="h-3.5 w-3.5" aria-hidden="true" />
              <span>
                {isId
                  ? "Model Default & Custom Headers (opsional)"
                  : "Default Models & Headers (optional)"}
              </span>
            </span>
            <span className="flex items-center gap-2">
              <span className="rounded-md bg-[var(--color-surface)] px-1.5 py-0.5 font-mono text-[10px] text-[var(--color-text-subtle)]">
                {modelCount} {isId ? "model" : "models"} · {headerCount}{" "}
                {isId ? "header" : "headers"}
              </span>
              <ChevronDown
                className={cn(
                  "h-3.5 w-3.5 transition-transform",
                  (showAdvanced || hasAdvancedErrors) && "rotate-180",
                )}
                aria-hidden="true"
              />
            </span>
          </button>

          {showAdvanced || hasAdvancedErrors ? (
            <div className="space-y-3 border-t border-[var(--color-border-subtle)] p-3.5">
              <div className="space-y-1.5">
                <Label htmlFor="modal-node-models">
                  {isId
                    ? "Model Default (pisahkan dengan koma)"
                    : "Default Models (comma-separated)"}
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
                  {isId
                    ? "Headers Default (Key: Value per baris)"
                    : "Default Headers (Key: Value per line)"}
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
            </div>
          ) : null}
        </div>

        {/* Inline Pre-Save Connectivity Test Result */}
        <DraftTestResultCard
          result={testState?.result ?? null}
          localError={testValidationError}
          isStale={isTestStale}
          protocol={watchedProtocol}
          isId={isId}
        />
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
  isTestingAll = false,
  liveTestResults,
  onTestAll,
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
  isTestingAll?: boolean;
  liveTestResults: Record<string, TestProviderConnectionResponse>;
  onTestAll?: (entry: ProviderEntry) => void;
  onTestConnection: (conn: ProviderConnectionSummary) => void;
  onToggleConnection: (conn: ProviderConnectionSummary) => void;
  onDeleteConnection: (conn: ProviderConnectionSummary) => void;
  onAddConnection: (entry: ProviderEntry, mode?: "oauth" | "apikey") => void;
  onDeleteCustomNode: (node: ProviderNode) => void;
  isDeletingNode: boolean;
}) {
  const [confirmDeleteConnId, setConfirmDeleteConnId] = useState<string | null>(null);
  const [confirmDeleteNode, setConfirmDeleteNode] = useState(false);
  const [modelSearch, setModelSearch] = useState("");
  const [showAllModels, setShowAllModels] = useState(false);
  const [copiedModelId, setCopiedModelId] = useState<string | null>(null);
  const [copiedBaseUrl, setCopiedBaseUrl] = useState(false);
  const entryId = entry?.id ?? null;

  useEffect(() => {
    if (entryId !== null) {
      setConfirmDeleteConnId(null);
      setConfirmDeleteNode(false);
      setModelSearch("");
      setShowAllModels(false);
      setCopiedModelId(null);
      setCopiedBaseUrl(false);
    }
  }, [entryId]);

  const allModels = useMemo(() => {
    if (!entry) return [];
    const seen = new Set<string>();
    const list: string[] = [];
    for (const m of entry.defaultModels) {
      const trimmed = m.trim();
      if (trimmed && !seen.has(trimmed)) {
        seen.add(trimmed);
        list.push(trimmed);
      }
    }
    for (const conn of entry.connections) {
      for (const m of conn.models) {
        const trimmed = m.trim();
        if (trimmed && !seen.has(trimmed)) {
          seen.add(trimmed);
          list.push(trimmed);
        }
      }
    }
    return list;
  }, [entry]);

  const filteredModels = useMemo(() => {
    const q = modelSearch.trim().toLowerCase();
    if (!q) return allModels;
    return allModels.filter((m) => m.toLowerCase().includes(q));
  }, [allModels, modelSearch]);

  if (!entry) {
    return null;
  }

  const keyUrl = entry.apiKeyUrl ?? entry.signupUrl;
  const enabledCount = entry.connections.filter((c) => c.enabled).length;
  const authFlowBadge =
    entry.id === "aws-quick"
      ? "LOOPBACK OIDC"
      : entry.id === "github-copilot" || entry.id === "kiro" || entry.id === "freepi"
        ? "DEVICE CODE"
        : entry.isNativeOAuth
          ? "OAUTH PKCE"
          : entry.category === "free"
            ? entry.noAuth
              ? "FREE · NO AUTH"
              : "FREE TIER"
            : "API KEY";

  const primaryConnectLabel = entry.isNativeOAuth
    ? entry.id === "aws-quick"
      ? isId
        ? "Login Amazon Quick"
        : "Sign In Amazon Quick"
      : entry.id === "github-copilot" || entry.id === "kiro" || entry.id === "freepi"
        ? isId
          ? "Hubungkan Device Login"
          : "Connect via Device Code"
        : isId
          ? "Hubungkan via OAuth"
          : "Connect via OAuth"
    : isId
      ? "Tambah API Key"
      : "Add API Key";

  const visibleModels = showAllModels ? filteredModels : filteredModels.slice(0, 10);
  const modelPrefix = entry.alias || entry.id;

  const handleCopyModel = async (modelId: string) => {
    const qualified = modelId.includes("/") ? modelId : `${modelPrefix}/${modelId}`;
    try {
      await navigator.clipboard.writeText(qualified);
      setCopiedModelId(modelId);
      window.setTimeout(() => setCopiedModelId((cur) => (cur === modelId ? null : cur)), 1800);
    } catch {
      // Ignore clipboard failures
    }
  };

  const handleCopyBaseUrl = async () => {
    if (!entry.baseUrl) return;
    try {
      await navigator.clipboard.writeText(entry.baseUrl);
      setCopiedBaseUrl(true);
      window.setTimeout(() => setCopiedBaseUrl(false), 1800);
    } catch {
      // Ignore clipboard failures
    }
  };

  return (
    <ModalShell
      open={Boolean(entry)}
      onClose={onClose}
      variant="drawer"
      size="xl"
      title={entry.name}
      subtitle={`${entry.id}${entry.alias && entry.alias !== entry.id ? ` (alias: ${entry.alias})` : ""} · ${entry.protocol.toUpperCase()}`}
      badge={
        <Badge variant="outline" className="font-mono text-[10px] uppercase">
          {authFlowBadge}
        </Badge>
      }
      icon={
        <ProviderLogo
          providerId={entry.id}
          name={entry.name}
          protocol={entry.protocol}
          isCustom={entry.isCustom}
          size={24}
        />
      }
      footer={
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex flex-wrap items-center gap-2">
            <Button type="button" size="sm" onClick={() => onAddConnection(entry)}>
              <Plus className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{primaryConnectLabel}</span>
            </Button>
            {entry.isNativeOAuth ? (
              <Button
                type="button"
                size="sm"
                variant="secondary"
                onClick={() => onAddConnection(entry, "apikey")}
              >
                <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Pakai API Key / Token" : "Use API Key / Token"}</span>
              </Button>
            ) : null}
          </div>
          <Button type="button" variant="secondary" size="sm" onClick={onClose}>
            {isId ? "Tutup" : "Close"}
          </Button>
        </div>
      }
    >
      <div className="space-y-5">
        {/* 3-Column Proportional Telemetry Strip */}
        <div className="grid grid-cols-3 gap-2.5">
          <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/50 px-3 py-2.5">
            <p className="text-[10px] font-medium uppercase tracking-wider text-[var(--color-text-muted)]">
              {isId ? "Koneksi" : "Connections"}
            </p>
            <div className="mt-1 flex items-baseline gap-1.5">
              <span className="font-mono text-base font-semibold tabular-nums text-[var(--color-text-main)]">
                {entry.stats.connected}
              </span>
              <span className="font-mono text-[11px] text-[var(--color-text-muted)]">
                / {entry.stats.total} {isId ? "aktif" : "active"}
              </span>
            </div>
          </div>

          <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/50 px-3 py-2.5">
            <p className="text-[10px] font-medium uppercase tracking-wider text-[var(--color-text-muted)]">
              {isId ? "Protokol" : "Protocol"}
            </p>
            <p className="mt-1 truncate font-mono text-xs font-semibold text-[var(--color-text-main)]">
              {entry.protocol}
            </p>
          </div>

          <div className="rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/50 px-3 py-2.5">
            <p className="text-[10px] font-medium uppercase tracking-wider text-[var(--color-text-muted)]">
              {isId ? "Model" : "Models"}
            </p>
            <p className="mt-1 font-mono text-base font-semibold tabular-nums text-[var(--color-text-main)]">
              {allModels.length}
            </p>
          </div>
        </div>

        {/* Provider Overview, Endpoint & Official Links */}
        <div className="space-y-3 rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 p-3.5">
          {entry.baseUrl ? (
            <div className="flex items-center justify-between gap-2 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface)] px-3 py-2">
              <div className="min-w-0">
                <span className="block text-[10px] font-medium uppercase tracking-wider text-[var(--color-text-subtle)]">
                  Base Endpoint
                </span>
                <p className="truncate font-mono text-xs text-[var(--color-text-main)]">
                  {entry.baseUrl}
                </p>
              </div>
              <button
                type="button"
                onClick={() => void handleCopyBaseUrl()}
                title={isId ? "Salin Base URL" : "Copy Base URL"}
                className="flex h-7 w-7 shrink-0 cursor-pointer items-center justify-center rounded-md border border-[var(--color-border)] bg-[var(--color-surface-2)] text-[var(--color-text-muted)] transition-colors hover:border-[var(--color-primary)]/40 hover:text-[var(--color-primary)]"
              >
                {copiedBaseUrl ? (
                  <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" aria-hidden="true" />
                ) : (
                  <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                )}
              </button>
            </div>
          ) : null}

          {entry.serviceKinds && entry.serviceKinds.length > 0 ? (
            <div className="flex flex-wrap items-center gap-1.5">
              {entry.serviceKinds.map((kind) => (
                <span
                  key={kind}
                  className="rounded-md border border-[var(--color-border)] bg-[var(--color-surface)] px-2 py-0.5 font-mono text-[10px] uppercase text-[var(--color-text-muted)]"
                >
                  {kind}
                </span>
              ))}
            </div>
          ) : null}

          {entry.noticeText ? (
            <div className="flex items-start gap-2 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface)]/70 px-3 py-2 text-[11px] leading-relaxed text-[var(--color-text-muted)]">
              <Info className="mt-0.5 h-3.5 w-3.5 shrink-0 text-[var(--color-primary)]" />
              <p>{entry.noticeText}</p>
            </div>
          ) : null}

          {entry.website || keyUrl ? (
            <div className="flex flex-wrap items-center gap-3 pt-0.5 text-xs">
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
                        ? "Portal Akun"
                        : "Account Portal"}
                  </span>
                  <ExternalLink className="h-3 w-3" />
                </a>
              ) : null}
            </div>
          ) : null}
        </div>

        {/* Connections Section */}
        <section className="space-y-3">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              <h4 className="text-xs font-semibold uppercase tracking-wider text-[var(--color-text-main)]">
                {isId ? "Akun & Koneksi Terdaftar" : "Configured Connections"}
              </h4>
              <span className="rounded-full bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[10px] text-[var(--color-text-muted)]">
                {entry.connections.length}
              </span>
            </div>

            {enabledCount > 0 && onTestAll ? (
              <Button
                type="button"
                size="sm"
                variant="secondary"
                disabled={isTestingAll || busyConnectionId !== null}
                onClick={() => onTestAll(entry)}
              >
                {isTestingAll ? (
                  <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                ) : (
                  <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                )}
                <span>
                  {isTestingAll
                    ? isId
                      ? "Menguji..."
                      : "Testing..."
                    : isId
                      ? "Tes Semua"
                      : "Test All"}
                </span>
              </Button>
            ) : null}
          </div>

          {entry.connections.length === 0 ? (
            <div className="space-y-3 rounded-xl border border-dashed border-[var(--color-border)] bg-[var(--color-surface-2)]/25 p-5 text-center">
              <p className="text-sm font-semibold text-[var(--color-text-main)]">
                {isId ? "Belum ada akun terhubung" : "No accounts connected yet"}
              </p>
              <p className="mx-auto max-w-md text-xs leading-relaxed text-[var(--color-text-muted)]">
                {isId
                  ? `Hubungkan akun ${entry.name} pertama Anda untuk mulai merutekan request ke model-model di bawah.`
                  : `Connect your first ${entry.name} account to start routing requests to the models below.`}
              </p>
              <div className="flex flex-wrap items-center justify-center gap-2 pt-1">
                <Button type="button" size="sm" onClick={() => onAddConnection(entry)}>
                  <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{primaryConnectLabel}</span>
                </Button>
                {entry.isNativeOAuth ? (
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    onClick={() => onAddConnection(entry, "apikey")}
                  >
                    <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                    <span>{isId ? "Impor Token / API Key" : "Import Token / API Key"}</span>
                  </Button>
                ) : null}
              </div>
            </div>
          ) : (
            <div className="space-y-2.5">
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
                        <div className="flex flex-wrap items-center gap-1.5">
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
                          <Badge variant="outline" className="font-mono text-[10px]">
                            P{conn.priority}
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
                          {conn.base_url} · {isId ? "Tes terakhir:" : "Tested:"}{" "}
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
        </section>

        {/* Interactive Models Catalog Section (xrouter + omniroute parity, click to copy ID) */}
        {allModels.length > 0 ? (
          <section className="space-y-2.5">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <div className="flex items-center gap-2">
                <Layers className="h-3.5 w-3.5 text-[var(--color-primary)]" aria-hidden="true" />
                <h4 className="text-xs font-semibold uppercase tracking-wider text-[var(--color-text-main)]">
                  {isId ? "Katalog Model Tersedia" : "Available Models"}
                </h4>
                <span className="rounded-full bg-[var(--color-surface-2)] px-2 py-0.5 font-mono text-[10px] text-[var(--color-text-muted)]">
                  {allModels.length}
                </span>
              </div>
              <span className="text-[11px] text-[var(--color-text-muted)]">
                {isId ? "Klik model untuk salin ID" : "Click model to copy ID"}
              </span>
            </div>

            {allModels.length > 6 ? (
              <div className="relative">
                <Search
                  className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-[var(--color-text-muted)]"
                  aria-hidden="true"
                />
                <Input
                  value={modelSearch}
                  onChange={(e) => setModelSearch(e.target.value)}
                  placeholder={isId ? "Cari model..." : "Filter models..."}
                  className="h-8 pl-8 text-xs"
                />
              </div>
            ) : null}

            <div className="grid grid-cols-1 gap-1.5 sm:grid-cols-2">
              {visibleModels.map((modelId) => {
                const isCopied = copiedModelId === modelId;
                const qualified = modelId.includes("/") ? modelId : `${modelPrefix}/${modelId}`;
                return (
                  <button
                    key={modelId}
                    type="button"
                    onClick={() => void handleCopyModel(modelId)}
                    title={isId ? `Salin ${qualified}` : `Copy ${qualified}`}
                    className="group flex cursor-pointer items-center justify-between gap-2 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 px-2.5 py-1.5 text-left transition-colors hover:border-[var(--color-primary)]/45 hover:bg-[var(--color-surface-2)]"
                  >
                    <span className="truncate font-mono text-[11px] text-[var(--color-text-main)]">
                      {modelId}
                    </span>
                    {isCopied ? (
                      <Check
                        className="h-3.5 w-3.5 shrink-0 text-[var(--status-ok)]"
                        aria-hidden="true"
                      />
                    ) : (
                      <Copy
                        className="h-3 w-3 shrink-0 text-[var(--color-text-muted)] opacity-50 transition-opacity group-hover:opacity-100"
                        aria-hidden="true"
                      />
                    )}
                  </button>
                );
              })}
            </div>

            {filteredModels.length > 10 ? (
              <button
                type="button"
                onClick={() => setShowAllModels((prev) => !prev)}
                className="cursor-pointer text-xs font-medium text-[var(--color-primary)] hover:underline"
              >
                {showAllModels
                  ? isId
                    ? "Tampilkan lebih sedikit"
                    : "Show fewer models"
                  : isId
                    ? `Tampilkan semua (${filteredModels.length} model)`
                    : `Show all ${filteredModels.length} models`}
              </button>
            ) : null}
          </section>
        ) : null}

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
