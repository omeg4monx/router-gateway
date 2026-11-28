import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  Check,
  CheckCircle2,
  ChevronDown,
  Copy,
  ExternalLink,
  Eye,
  EyeOff,
  Info,
  KeyRound,
  RefreshCw,
  Save,
  Search,
  Server,
  Settings2,
  ShieldCheck,
  Sparkles,
  X,
  Zap,
} from "lucide-react";
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
  useUpdateProviderConnectionMutation,
} from "./hooks";
import {
  buildTestFingerprint,
  CUSTOM_ENDPOINT_PRESETS,
  type CustomEndpointPreset,
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
  connectionToEditFormValues,
  createConnectionSchema,
  createNodeSchema,
  DEFAULT_CREATE_CONNECTION_VALUES,
  DEFAULT_CREATE_NODE_VALUES,
  DEFAULT_OAUTH_EXCHANGE_VALUES,
  DEFAULT_OAUTH_START_VALUES,
  DEFAULT_PRIORITY,
  DEFAULT_UPDATE_CONNECTION_VALUES,
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
  toUpdateConnectionRequest,
  type UpdateConnectionFormInput,
  type UpdateConnectionFormOutput,
  updateConnectionSchema,
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

const MODAL_EXIT_MS = 200;

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
  const [rendered, setRendered] = useState(open);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    if (open) {
      setRendered(true);
      let frame2 = 0;
      const frame1 = window.requestAnimationFrame(() => {
        frame2 = window.requestAnimationFrame(() => {
          setVisible(true);
        });
      });
      return () => {
        window.cancelAnimationFrame(frame1);
        if (frame2) {
          window.cancelAnimationFrame(frame2);
        }
      };
    }
    setVisible(false);
    const timer = window.setTimeout(() => {
      setRendered(false);
    }, MODAL_EXIT_MS);
    return () => window.clearTimeout(timer);
  }, [open]);

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

  if (!rendered) {
    return null;
  }

  return (
    <div
      data-state={visible ? "open" : "closed"}
      className={cn(
        "provider-modal-root fixed inset-0 z-50 flex",
        variant === "drawer"
          ? "justify-end"
          : "items-end justify-center p-0 sm:items-center sm:p-6",
        !visible && "pointer-events-none",
      )}
    >
      <button
        type="button"
        aria-label="Close modal backdrop"
        onClick={onClose}
        className="provider-modal-backdrop fixed inset-0 bg-black/55"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-label={title}
        className={cn(
          "provider-modal-panel relative z-10 flex flex-col overflow-hidden border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl",
          variant === "drawer"
            ? cn(
                "provider-modal-panel-drawer h-[100dvh] w-full border-y-0 border-r-0",
                MODAL_SIZE_CLASS[size],
              )
            : cn(
                "max-h-[92dvh] w-full rounded-t-2xl sm:max-h-[90dvh] sm:rounded-2xl",
                MODAL_SIZE_CLASS[size],
              ),
        )}
      >
        <span
          aria-hidden="true"
          className="provider-modal-accent pointer-events-none absolute inset-x-0 top-0 z-20 h-[2px]"
        />
        <div className="flex shrink-0 items-start justify-between gap-3 border-b border-[var(--color-border-subtle)] px-4 py-3.5 sm:px-6 sm:py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/85 p-1.5 text-[var(--color-primary)]">
              {icon}
            </div>
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-1.5 sm:gap-2">
                <h3 className="truncate text-base font-semibold tracking-tight text-[var(--color-text-main)]">
                  {title}
                </h3>
                {badge}
              </div>
              {subtitle ? (
                <p className="mt-0.5 line-clamp-2 text-xs text-[var(--color-text-muted)] sm:truncate">
                  {subtitle}
                </p>
              ) : null}
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label="Close modal"
            className="flex h-9 w-9 shrink-0 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)] sm:h-8 sm:w-8"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </button>
        </div>

        <div className="custom-scrollbar flex-1 overflow-y-auto px-4 py-4 sm:px-6 sm:py-5">
          {children}
        </div>

        {footer ? (
          <div className="shrink-0 border-t border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35 px-4 pt-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] sm:px-6 sm:py-3.5">
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
  lockProvider = false,
}: {
  open: boolean;
  onClose: () => void;
  initialProviderId: string | null;
  allEntries: ProviderEntry[];
  isId: boolean;
  onCreated: (message: string) => void;
  lockProvider?: boolean;
}) {
  const createConnMutation = useCreateProviderConnectionMutation();
  const createNodeMutation = useCreateProviderNodeMutation();
  const testDraftMutation = useTestProviderNodeDraftMutation();

  const [mode, setMode] = useState<"single" | "bulk">("single");
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(true);
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
    setShowAdvanced(true);
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
          ? lockProvider
            ? isId
              ? `Tambah ${credMeta.label} ${activeEntry.name}`
              : `Add ${activeEntry.name} ${credMeta.label}`
            : `Setup ${activeEntry.name}`
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
        <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
          <div className="w-full sm:w-auto">
            {mode === "single" ? (
              <Button
                type="button"
                variant="secondary"
                size="sm"
                className="w-full sm:w-auto"
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
              <span className="block text-xs text-[var(--color-text-muted)]">
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

          <div className="grid w-full grid-cols-2 gap-2 sm:flex sm:w-auto sm:items-center">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
              onClick={onClose}
            >
              {isId ? "Batal" : "Cancel"}
            </Button>
            {mode === "single" ? (
              <Button
                type="submit"
                form="add-apikey-form"
                size="sm"
                className="w-full sm:w-auto"
                disabled={isSaving}
              >
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
                className="w-full sm:w-auto"
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
        {/* Selected Provider Card + Searchable Picker (hidden when locked to a provider page) */}
        {!lockProvider ? (
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
        ) : null}

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
              <div className="flex items-center gap-2">
                <div className="relative flex-1">
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
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  className="shrink-0"
                  disabled={!watchedApiKey?.trim() || testDraftMutation.isPending || isSaving}
                  onClick={() => void handleTestDraft()}
                >
                  {testDraftMutation.isPending
                    ? isId
                      ? "Memeriksa..."
                      : "Checking..."
                    : isId
                      ? "Cek"
                      : "Check"}
                </Button>
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
                    P{watchedPriority || String(DEFAULT_PRIORITY)} · {modelCount}{" "}
                    {isId ? "model" : "models"}
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
  lockProvider = false,
}: {
  open: boolean;
  onClose: () => void;
  initialProviderId: string | null;
  allEntries: ProviderEntry[];
  isId: boolean;
  onConnected: (message: string) => void;
  onSwitchToApiKey: (providerId: string) => void;
  lockProvider?: boolean;
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
  const [callbackDetected, setCallbackDetected] = useState(false);
  const [copied, setCopied] = useState(false);
  const [outcome, setOutcome] = useState<
    { status: "success"; accountName: string } | { status: "error"; message: string } | null
  >(null);
  const pkceAutoStartedRef = useRef<string | null>(null);
  const deviceAutoStartedRef = useRef<string | null>(null);
  const callbackHandledRef = useRef<string | null>(null);

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
  const modalInitKeyRef = useRef<string | null>(null);

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
        setOutcome({ status: "success", accountName: conn.name });
      } catch (err) {
        awsQuickSettledRef.current = false;
        const msg = err instanceof Error ? err.message : "Failed to complete Amazon Quick OAuth";
        setAwsQuickError(msg);
        setOutcome({ status: "error", message: msg });
      } finally {
        setAwsQuickBusy(false);
      }
    },
    [exchangeOAuthMutateAsync, isId, onConnected, startForm],
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
      modalInitKeyRef.current = null;
      awsQuickAutoStartedRef.current = false;
      pkceAutoStartedRef.current = null;
      deviceAutoStartedRef.current = null;
      stopAwsQuickIfOwned();
      return;
    }
    const initKey = initialProviderId ?? "";
    if (modalInitKeyRef.current === initKey) {
      return;
    }
    modalInitKeyRef.current = initKey;

    startOAuthReset();
    exchangeOAuthReset();
    pasteOAuthReset();
    setActiveSession(null);
    setCallbackDetected(false);
    setCopied(false);
    setOutcome(null);
    pkceAutoStartedRef.current = null;
    deviceAutoStartedRef.current = null;
    callbackHandledRef.current = null;
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
    initialProviderId,
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
  const oauthMeta = currentEntry?.adapter?.oauth;
  const isNonNativeOAuth = Boolean(selectedEntry && !selectedEntry.isNativeOAuth);
  const isAwsQuick = oauthMeta?.flow_type === "loopback_proxy" || watchedProvider === "aws-quick";
  const isDeviceCodeProvider =
    oauthMeta?.flow_type === "device_code" ||
    watchedProvider === "github-copilot" ||
    watchedProvider === "kiro" ||
    watchedProvider === "freepi" ||
    watchedProvider === "codebuddy-cn" ||
    watchedProvider === "codebuddy-intl";
  const awsQuickPort = useMemo(() => {
    const uri = awsQuickSession?.redirect_uri;
    if (!uri) return "<port>";
    try {
      return new URL(uri).port || "<port>";
    } catch {
      return "<port>";
    }
  }, [awsQuickSession?.redirect_uri]);

  // Auto-start Amazon Quick loopback proxy and open sign-in popup immediately on open
  // (matches xrouter OAuthModal zero-click startup).
  useEffect(() => {
    if (!open || !isAwsQuick || awsQuickSession || awsQuickAutoStartedRef.current || outcome) {
      return;
    }
    awsQuickAutoStartedRef.current = true;
    startAwsQuickFlow(true);
  }, [open, isAwsQuick, awsQuickSession, outcome, startAwsQuickFlow]);

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
          setDeviceSession(null);
          setOutcome({ status: "success", accountName: conn.name });
          handleClose();
          return;
        }
        if (data.status === "error") {
          const msg =
            data.error || (isId ? "Otorisasi perangkat gagal." : "Device authorization failed.");
          setDeviceError(msg);
          setOutcome({ status: "error", message: msg });
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
  }, [open, deviceSession, exchangeOAuthMutateAsync, handleClose, isId, onConnected, startForm]);

  const startDeviceCodeFlow = useCallback(() => {
    const providerId = startForm.getValues("provider");
    const connName = startForm.getValues("connection_name")?.trim();
    setDeviceError(null);
    setOutcome(null);
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
          setOutcome({ status: "error", message: err.message });
        },
      },
    );
  }, [startForm, startOAuthMutate]);

  // Auto-start Device Code flow on open (OmniRoute parity)
  useEffect(() => {
    if (!open || !isDeviceCodeProvider || isNonNativeOAuth || outcome) return;
    if (deviceAutoStartedRef.current === watchedProvider) return;
    deviceAutoStartedRef.current = watchedProvider;
    startDeviceCodeFlow();
  }, [open, isDeviceCodeProvider, isNonNativeOAuth, outcome, watchedProvider, startDeviceCodeFlow]);

  // Auto-start PKCE / AuthCode flow immediately on open (OmniRoute / xrouter zero-config parity)
  useEffect(() => {
    if (!open || isAwsQuick || isDeviceCodeProvider || isNonNativeOAuth || outcome) {
      return;
    }
    if (pkceAutoStartedRef.current === watchedProvider) return;
    pkceAutoStartedRef.current = watchedProvider;
    const dynamicRedirectUri = (() => {
      if (typeof window === "undefined" || watchedProvider === "codex") {
        return undefined;
      }
      const { protocol, hostname, port, origin } = window.location;
      if (
        (watchedProvider === "antigravity" || watchedProvider === "gemini-oauth") &&
        (hostname === "localhost" || hostname === "127.0.0.1")
      ) {
        const portPart = port ? `:${port}` : "";
        return `${protocol}//localhost${portPart}/api/providers/oauth/callback`;
      }
      return `${origin}/api/providers/oauth/callback`;
    })();
    startOAuthMutate(
      {
        provider: watchedProvider,
        body: dynamicRedirectUri ? { redirect_uri: dynamicRedirectUri } : {},
      },
      {
        onSuccess: (session) => {
          setActiveSession({ provider: watchedProvider, session });
          exchangeForm.setValue("state", session.state, { shouldDirty: false });
          if (session.authorize_url && typeof window !== "undefined") {
            window.open(session.authorize_url, "oauth_popup", "width=600,height=700");
          }
        },
        onError: (err) => {
          setOutcome({ status: "error", message: err.message });
        },
      },
    );
  }, [
    open,
    isAwsQuick,
    isDeviceCodeProvider,
    isNonNativeOAuth,
    outcome,
    watchedProvider,
    startOAuthMutate,
    exchangeForm,
  ]);

  // Listen for automatic popup/tab OAuth callback via postMessage, BroadcastChannel, and localStorage.
  useEffect(() => {
    if (!open || !activeSession || outcome) return;
    const expectedState = activeSession.session.state;
    const targetProvider = activeSession.provider;

    const processCallbackPayload = (raw: unknown) => {
      if (!raw || typeof raw !== "object") return;
      const payload = raw as {
        code?: string;
        state?: string;
        error?: string;
        errorDescription?: string;
        callbackUrl?: string;
      };
      if (payload.state && expectedState && payload.state !== expectedState) {
        return;
      }
      const dedupKey = `${expectedState}:${payload.code ?? payload.error ?? ""}`;
      if (callbackHandledRef.current === dedupKey) return;

      if (payload.error) {
        callbackHandledRef.current = dedupKey;
        setOutcome({
          status: "error",
          message: payload.errorDescription || payload.error,
        });
        return;
      }

      if (!payload.code && !payload.callbackUrl) return;
      callbackHandledRef.current = dedupKey;
      const connName = startForm.getValues("connection_name")?.trim();
      void exchangeOAuthMutateAsync({
        provider: targetProvider,
        body: {
          ...(payload.callbackUrl ? { callback_url: payload.callbackUrl } : {}),
          ...(payload.code ? { code: payload.code } : {}),
          state: payload.state || expectedState,
          priority: DEFAULT_PRIORITY,
          enabled: true,
          ...(connName ? { connection_name: connName } : {}),
        },
      })
        .then((conn) => {
          onConnected(
            isId
              ? `Akun OAuth "${conn.name}" berhasil dihubungkan.`
              : `Connected OAuth account "${conn.name}".`,
          );
          setOutcome({ status: "success", accountName: conn.name });
        })
        .catch((err) => {
          callbackHandledRef.current = null;
          setOutcome({
            status: "error",
            message: err instanceof Error ? err.message : "OAuth exchange failed",
          });
        });
    };

    const handleMessage = (event: MessageEvent) => {
      const isLocal =
        event.origin.includes("localhost") ||
        event.origin.includes("127.0.0.1") ||
        event.origin === window.location.origin;
      if (!isLocal) return;
      if (event.data?.type === "oauth_callback") {
        processCallbackPayload(event.data.data);
      }
    };
    window.addEventListener("message", handleMessage);

    let channel: BroadcastChannel | null = null;
    try {
      channel = new BroadcastChannel("oauth_callback");
      channel.onmessage = (event) => processCallbackPayload(event.data);
    } catch {
      // BroadcastChannel unavailable in some environments.
    }

    const handleStorage = (event: StorageEvent) => {
      if (event.key === "oauth_callback" && event.newValue) {
        try {
          processCallbackPayload(JSON.parse(event.newValue));
          window.localStorage.removeItem("oauth_callback");
        } catch {
          // Ignore malformed storage payload.
        }
      }
    };
    window.addEventListener("storage", handleStorage);

    return () => {
      window.removeEventListener("message", handleMessage);
      window.removeEventListener("storage", handleStorage);
      channel?.close();
    };
  }, [open, activeSession, outcome, exchangeOAuthMutateAsync, isId, onConnected, startForm]);

  // Poll fixed-port loopback callback status (e.g. Codex port 1455 or any strategy with fixed_port)
  // for zero copy-paste auto-completion, mirroring xrouter's OAuthModal proxy lifecycle.
  useEffect(() => {
    const hasLoopbackProxy =
      Boolean(activeSession?.session.fixed_port) || activeSession?.provider === "codex";
    if (!open || !hasLoopbackProxy || !activeSession?.session.state || outcome) {
      return;
    }
    let cancelled = false;
    let timerId: number | null = null;
    const loopbackProvider = activeSession.provider;
    const sessionState = activeSession.session.state;

    const tick = async () => {
      if (cancelled || callbackHandledRef.current === `${loopbackProvider}:${sessionState}`) return;
      try {
        const data = await fetchProviderOAuthStatus(loopbackProvider, sessionState);
        if (cancelled) return;
        if (data.status === "done") {
          callbackHandledRef.current = `${loopbackProvider}:${sessionState}`;
          const connName = startForm.getValues("connection_name")?.trim();
          const conn = await exchangeOAuthMutateAsync({
            provider: loopbackProvider,
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
          setDeviceSession(null);
          setOutcome({ status: "success", accountName: conn.name });
          handleClose();
          return;
        }
        if (data.status === "error") {
          setOutcome({
            status: "error",
            message:
              data.error || (isId ? "Otorisasi OAuth gagal." : "OAuth authorization failed."),
          });
          return;
        }
      } catch {
        // Transient poll error — keep polling.
      }
      timerId = window.setTimeout(tick, 1500);
    };

    timerId = window.setTimeout(tick, 1500);
    return () => {
      cancelled = true;
      if (timerId !== null) window.clearTimeout(timerId);
    };
  }, [
    open,
    activeSession,
    outcome,
    exchangeOAuthMutateAsync,
    handleClose,
    isId,
    onConnected,
    startForm,
  ]);

  const pasteCallbackFromClipboard = async () => {
    try {
      const text = (await navigator.clipboard.readText()).trim();
      if (!text) return;
      exchangeForm.setValue("callback_url", text, { shouldDirty: true, shouldValidate: true });
      handleCallbackUrlChange(text);
    } catch {
      // Clipboard blocked — user can paste manually with Ctrl/Cmd+V.
    }
  };

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
          setOutcome({ status: "success", accountName: conn.name });
        },
        onError: (err) => {
          setOutcome({ status: "error", message: err.message });
        },
      },
    );
  };

  const handleRetryOAuth = () => {
    setOutcome(null);
    setAwsQuickError(null);
    setDeviceError(null);
    startOAuthReset();
    exchangeOAuthReset();
    if (isAwsQuick) {
      startAwsQuickFlow(true);
      return;
    }
    if (isDeviceCodeProvider) {
      deviceAutoStartedRef.current = null;
      startDeviceCodeFlow();
      return;
    }
    pkceAutoStartedRef.current = null;
    setActiveSession(null);
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
        lockProvider
          ? undefined
          : isAwsQuick
            ? isId
              ? "Alur loopback 2-hop Amazon Quick."
              : "Amazon Quick loopback flow."
            : isDeviceCodeProvider
              ? isId
                ? "Masukkan kode di browser, koneksi tersimpan otomatis."
                : "Enter the code in your browser — connects automatically."
              : isId
                ? "Autentikasi akun OAuth upstream."
                : "Authenticate your upstream OAuth account."
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
        outcome ? undefined : isNonNativeOAuth ? (
          <div className="flex items-center justify-end gap-2">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
              onClick={handleClose}
            >
              {isId ? "Tutup" : "Close"}
            </Button>
          </div>
        ) : isAwsQuick ? (
          <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
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

            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
              onClick={handleClose}
            >
              {isId ? "Tutup" : "Close"}
            </Button>
          </div>
        ) : isDeviceCodeProvider ? (
          <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
            <Button
              type="button"
              variant={deviceSession ? "secondary" : "default"}
              size="sm"
              className="w-full sm:w-auto"
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

            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
              onClick={handleClose}
            >
              {isId ? "Tutup" : "Close"}
            </Button>
          </div>
        ) : (
          <div className="grid w-full grid-cols-2 gap-2 sm:flex sm:items-center sm:justify-end">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
              onClick={handleClose}
            >
              {isId ? "Batal" : "Cancel"}
            </Button>
            <Button
              type="submit"
              form="oauth-exchange-form"
              size="sm"
              className="w-full sm:w-auto"
              disabled={!activeSession || exchangeOAuthMutation.isPending}
            >
              {exchangeOAuthMutation.isPending
                ? isId
                  ? "Menghubungkan..."
                  : "Connecting..."
                : isId
                  ? "Hubungkan"
                  : "Connect"}
            </Button>
          </div>
        )
      }
    >
      <div className="space-y-4">
        {outcome?.status === "success" ? (
          <div className="py-6 text-center">
            <div className="mx-auto mb-4 flex size-16 items-center justify-center rounded-full bg-[var(--status-ok-bg)]">
              <CheckCircle2 className="h-8 w-8 text-[var(--status-ok)]" aria-hidden="true" />
            </div>
            <h3 className="mb-2 text-lg font-semibold text-[var(--color-text-main)]">
              {isId ? "Berhasil Terhubung!" : "Connected Successfully!"}
            </h3>
            <p className="mb-5 text-sm text-[var(--color-text-muted)]">
              {isId
                ? `Akun ${currentEntry?.name ?? watchedProvider} (${outcome.accountName}) berhasil dihubungkan.`
                : `Your ${currentEntry?.name ?? watchedProvider} account (${outcome.accountName}) has been connected.`}
            </p>
            <Button type="button" size="sm" className="w-full" onClick={handleClose}>
              {isId ? "Selesai" : "Done"}
            </Button>
          </div>
        ) : outcome?.status === "error" ? (
          <div className="py-6 text-center">
            <div className="mx-auto mb-4 flex size-16 items-center justify-center rounded-full bg-[var(--status-error-bg)]">
              <AlertTriangle className="h-8 w-8 text-[var(--status-error)]" aria-hidden="true" />
            </div>
            <h3 className="mb-2 text-lg font-semibold text-[var(--color-text-main)]">
              {isId ? "Koneksi Gagal" : "Connection Failed"}
            </h3>
            <p className="mb-5 break-words text-sm text-[var(--status-error)]">{outcome.message}</p>
            <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
              <Button
                type="button"
                variant="secondary"
                size="sm"
                className="w-full"
                onClick={handleRetryOAuth}
              >
                {isId ? "Coba Lagi" : "Try Again"}
              </Button>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="w-full"
                onClick={handleClose}
              >
                {isId ? "Batal" : "Cancel"}
              </Button>
            </div>
          </div>
        ) : isNonNativeOAuth && selectedEntry ? (
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
            {!lockProvider ? (
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
            ) : null}

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

                  <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
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

                  <p className="text-[11px] leading-relaxed text-[var(--color-text-muted)]">
                    {isId
                      ? "Ingin tanpa tempel manual? Jalankan dashboard di mesin ini atau buka SSH tunnel:"
                      : "Want zero pastes? Run the dashboard on this machine, or open an SSH tunnel:"}{" "}
                    <span className="break-all font-mono text-[var(--color-text-main)]">
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
            {!lockProvider ? (
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
            ) : null}

            {deviceError ? (
              <Alert variant="destructive">
                <AlertDescription>{deviceError}</AlertDescription>
              </Alert>
            ) : null}

            {deviceSession && deviceSession.provider === watchedProvider ? (
              <div className="space-y-3.5">
                {deviceSession.session.user_code ? (
                  <div className="space-y-3 rounded-xl border border-[var(--color-primary)]/35 bg-[var(--color-primary)]/10 p-4 text-center">
                    <p className="text-[11px] font-medium uppercase tracking-wider text-[var(--color-text-muted)]">
                      {isId ? "Kode Perangkat Anda" : "Your Device Code"}
                    </p>
                    <button
                      type="button"
                      onClick={() => void handleCopyUserCode()}
                      title={isId ? "Klik untuk menyalin kode" : "Click to copy code"}
                      className="mx-auto inline-flex cursor-pointer items-center gap-2.5 rounded-xl border border-[var(--color-primary)]/25 bg-[var(--color-surface)] px-4 py-2 transition-colors hover:border-[var(--color-primary)]"
                    >
                      <span className="break-all font-mono text-xl font-bold tracking-widest text-[var(--color-primary)] sm:text-2xl">
                        {deviceSession.session.user_code}
                      </span>
                      {copiedCode ? (
                        <Check
                          className="h-4 w-4 shrink-0 text-[var(--status-ok)]"
                          aria-hidden="true"
                        />
                      ) : (
                        <Copy
                          className="h-4 w-4 shrink-0 text-[var(--color-text-muted)]"
                          aria-hidden="true"
                        />
                      )}
                    </button>
                    <div className="flex flex-col gap-2 pt-1 sm:flex-row sm:justify-center">
                      <Button
                        type="button"
                        size="sm"
                        className="w-full sm:w-auto"
                        onClick={() => {
                          void handleCopyUserCode();
                          const target =
                            deviceSession.session.authorize_url ||
                            deviceSession.session.verification_uri;
                          if (target && typeof window !== "undefined") {
                            window.open(target, "_blank", "noopener,noreferrer");
                          }
                        }}
                      >
                        <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                        <span>
                          {copiedCode
                            ? isId
                              ? "Tersalin! Membuka Login..."
                              : "Copied! Opening Login..."
                            : isId
                              ? "Salin Kode & Buka"
                              : "Copy Code & Open"}
                        </span>
                      </Button>
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        className="w-full sm:w-auto"
                        onClick={() => void handleCopyUrl()}
                      >
                        {copied ? (
                          <Check
                            className="h-3.5 w-3.5 text-[var(--status-ok)]"
                            aria-hidden="true"
                          />
                        ) : (
                          <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                        )}
                        <span>{isId ? "Salin URL" : "Copy URL"}</span>
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="space-y-2 rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)]/55 p-3.5">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <span className="text-xs font-semibold text-[var(--color-text-main)]">
                        {isId ? "Buka URL Verifikasi" : "Open Verification URL"}
                      </span>
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
                )}

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
        ) : !activeSession ? (
          <div className="flex flex-col items-center justify-center gap-3 py-8 text-center">
            <RefreshCw
              className="h-6 w-6 animate-spin text-[var(--color-primary)]"
              aria-hidden="true"
            />
            <p className="text-sm font-medium text-[var(--color-text-main)]">
              {isId ? "Membuka halaman login…" : "Opening sign-in window…"}
            </p>
            <p className="text-xs text-[var(--color-text-muted)]">
              {isId
                ? `Menyiapkan sesi OAuth untuk ${currentEntry?.name ?? watchedProvider}.`
                : `Preparing OAuth session for ${currentEntry?.name ?? watchedProvider}.`}
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {activeSession ? (
              <>
                {/* OmniRoute Option A: Waiting for popup authorization */}
                <div className="flex items-center gap-2 rounded-lg border border-[var(--color-border)] bg-[var(--color-surface-2)]/55 px-3 py-2.5">
                  <RefreshCw
                    className="h-4 w-4 shrink-0 animate-spin text-[var(--color-primary)]"
                    aria-hidden="true"
                  />
                  <span className="text-xs font-medium text-[var(--color-text-main)]">
                    {isId
                      ? "Menunggu otorisasi popup browser…"
                      : "Waiting for popup authorization…"}
                  </span>
                </div>

                {/* OmniRoute Divider */}
                <div className="flex items-center gap-3">
                  <div className="h-px flex-1 bg-[var(--color-border)]" />
                  <span className="text-[10px] font-semibold uppercase tracking-wider text-[var(--color-text-muted)]">
                    {isId ? "Atau tempel URL callback manual" : "Or paste callback URL manually"}
                  </span>
                  <div className="h-px flex-1 bg-[var(--color-border)]" />
                </div>

                {/* OmniRoute Step 1: Compact Open / Copy URL */}
                <div className="space-y-2.5 rounded-xl border border-[var(--color-primary)]/35 bg-[var(--color-surface-2)]/65 p-3.5">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <span className="text-xs font-semibold text-[var(--color-text-main)]">
                      {isId
                        ? "Langkah 1: Buka halaman login di browser"
                        : "Step 1: Open sign-in page in your browser"}
                    </span>
                    <div className="flex w-full items-center gap-1.5 sm:w-auto">
                      <Button
                        type="button"
                        size="sm"
                        variant="secondary"
                        className="flex-1 sm:flex-none"
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
                      <Button
                        type="button"
                        size="sm"
                        className="flex-1 sm:flex-none"
                        onClick={() => {
                          window.open(
                            activeSession.session.authorize_url,
                            "oauth_popup",
                            "width=600,height=700",
                          );
                        }}
                      >
                        <span>{isId ? "Buka Login" : "Open Login"}</span>
                        <ExternalLink className="h-3.5 w-3.5" aria-hidden="true" />
                      </Button>
                    </div>
                  </div>
                  <Input
                    value={activeSession.session.authorize_url}
                    readOnly
                    aria-label="OAuth Authorization URL"
                    className="h-8 bg-[var(--color-surface)] font-mono text-[11px] text-[var(--color-text-muted)]"
                  />
                </div>
              </>
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
                    {isId ? "Langkah 2: Tempel Callback URL" : "Step 2: Paste Callback URL"}
                  </Label>
                  {callbackDetected ? (
                    <Badge variant="success" className="text-[10px]">
                      <Check className="mr-1 h-3 w-3" aria-hidden="true" />
                      {isId ? "Code & State terdeteksi" : "Code & State detected"}
                    </Badge>
                  ) : null}
                </div>
                <div className="flex items-center gap-2">
                  <Input
                    id="modal-oauth-callback"
                    placeholder="http://localhost:.../callback?code=...&state=..."
                    className="flex-1 font-mono text-xs"
                    {...exchangeForm.register("callback_url", {
                      onChange: (e) => handleCallbackUrlChange(e.target.value),
                    })}
                  />
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    className="shrink-0"
                    onClick={() => void pasteCallbackFromClipboard()}
                  >
                    {isId ? "Tempel" : "Paste"}
                  </Button>
                </div>
                {exchangeForm.formState.errors.callback_url ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {exchangeForm.formState.errors.callback_url.message}
                  </p>
                ) : null}
              </div>
            </form>
          </div>
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
        <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:items-center sm:justify-between">
          <Button
            type="button"
            variant="secondary"
            size="sm"
            className="w-full sm:w-auto"
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

          <div className="grid w-full grid-cols-2 gap-2 sm:flex sm:w-auto sm:items-center">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
              onClick={onClose}
            >
              {isId ? "Batal" : "Cancel"}
            </Button>
            <Button
              type="submit"
              form="add-custom-node-form"
              size="sm"
              className="w-full sm:w-auto"
              disabled={isSaving}
            >
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
 * 5. Edit Connection Modal (in-place edit & credential rotation on Provider Detail)
 * ========================================================================== */
export function EditConnectionModal({
  open,
  connection,
  isId,
  onClose,
  onUpdated,
}: {
  open: boolean;
  connection: ProviderConnectionSummary | null;
  isId: boolean;
  onClose: () => void;
  onUpdated?: (message: string) => void;
}) {
  const updateMutation = useUpdateProviderConnectionMutation();
  const [showApiKey, setShowApiKey] = useState(false);
  const [showAccessToken, setShowAccessToken] = useState(false);
  const [showRefreshToken, setShowRefreshToken] = useState(false);
  const [secretFieldsUnlocked, setSecretFieldsUnlocked] = useState(false);

  const editForm = useForm<UpdateConnectionFormInput, unknown, UpdateConnectionFormOutput>({
    resolver: zodResolver(updateConnectionSchema),
    defaultValues: DEFAULT_UPDATE_CONNECTION_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const {
    reset,
    formState: { errors },
  } = editForm;
  const resetMutation = updateMutation.reset;

  useEffect(() => {
    if (!open || !connection) {
      return;
    }
    reset(connectionToEditFormValues(connection));
    resetMutation();
    setShowApiKey(false);
    setShowAccessToken(false);
    setShowRefreshToken(false);
    setSecretFieldsUnlocked(false);
  }, [open, connection, reset, resetMutation]);

  const onSubmit = (values: UpdateConnectionFormOutput) => {
    if (!connection) {
      return;
    }
    updateMutation.reset();
    updateMutation.mutate(
      {
        id: connection.id,
        body: toUpdateConnectionRequest(values, connection.auth_type),
      },
      {
        onSuccess: (updated) => {
          onUpdated?.(
            isId
              ? `Perubahan pada koneksi "${updated.name}" disimpan.`
              : `Saved changes to "${updated.name}".`,
          );
          onClose();
        },
      },
    );
  };

  const formId = "edit-connection-modal-form";

  return (
    <ModalShell
      open={open && connection !== null}
      onClose={onClose}
      size="lg"
      title={
        connection
          ? isId
            ? `Edit Koneksi — ${connection.name}`
            : `Edit Connection — ${connection.name}`
          : isId
            ? "Edit Koneksi"
            : "Edit Connection"
      }
      subtitle={
        isId
          ? "Atur prioritas routing, daftar model, header upstream, atau rotasi kredensial tanpa meninggalkan halaman."
          : "Configure routing priority, allowed models, upstream headers, or rotate vaulted credentials in place."
      }
      badge={
        connection ? (
          <Badge variant={connection.auth_type === "oauth" ? "default" : "secondary"}>
            {connection.auth_type === "oauth" ? "OAuth" : "API Key"}
          </Badge>
        ) : undefined
      }
      icon={<Settings2 className="h-5 w-5" aria-hidden="true" />}
      footer={
        <div className="flex flex-col-reverse gap-2 sm:flex-row sm:items-center sm:justify-end">
          <Button
            type="button"
            variant="secondary"
            size="sm"
            className="w-full sm:w-auto"
            onClick={onClose}
            disabled={updateMutation.isPending}
          >
            {isId ? "Batal" : "Cancel"}
          </Button>
          <Button
            type="submit"
            form={formId}
            size="sm"
            className="w-full sm:w-auto"
            disabled={updateMutation.isPending}
          >
            <Save className="h-3.5 w-3.5" aria-hidden="true" />
            <span>
              {updateMutation.isPending
                ? isId
                  ? "Menyimpan..."
                  : "Saving..."
                : isId
                  ? "Simpan Perubahan"
                  : "Save Changes"}
            </span>
          </Button>
        </div>
      }
    >
      {connection ? (
        <form
          id={formId}
          autoComplete="off"
          onSubmit={editForm.handleSubmit(onSubmit)}
          noValidate
          className="space-y-4"
        >
          {updateMutation.error ? (
            <Alert variant="destructive">
              <AlertTriangle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>
                  {isId ? "Gagal memperbarui koneksi" : "Failed to update connection"}
                </AlertTitle>
                <AlertDescription className="font-mono">
                  {updateMutation.error.message}
                </AlertDescription>
              </div>
            </Alert>
          ) : null}

          <div className="grid grid-cols-1 gap-3.5 md:grid-cols-2">
            <div>
              <Label htmlFor="modal-edit-name">{isId ? "Nama Koneksi" : "Connection Name"}</Label>
              <Input
                id="modal-edit-name"
                type="text"
                placeholder={isId ? "contoh: OpenAI Utama" : "e.g. OpenAI Primary"}
                aria-invalid={Boolean(errors.name)}
                className="mt-1.5 placeholder:text-[var(--color-text-muted)]/50"
                {...editForm.register("name")}
              />
              {errors.name ? (
                <p className="mt-1 text-xs text-[var(--status-error)]">{errors.name.message}</p>
              ) : null}
            </div>

            <div>
              <Label htmlFor="modal-edit-priority">
                {isId
                  ? "Prioritas Routing (lebih kecil lebih dulu)"
                  : "Routing Priority (lower runs first)"}
              </Label>
              <Input
                id="modal-edit-priority"
                type="number"
                min={0}
                max={100000}
                step={1}
                placeholder={String(DEFAULT_PRIORITY)}
                aria-invalid={Boolean(errors.priority)}
                className="mt-1.5 tabular-nums placeholder:text-[var(--color-text-muted)]/50"
                {...editForm.register("priority")}
              />
              {errors.priority ? (
                <p className="mt-1 text-xs text-[var(--status-error)]">{errors.priority.message}</p>
              ) : null}
            </div>

            <div className="md:col-span-2">
              <Label htmlFor="modal-edit-base-url">Upstream Base URL</Label>
              <Input
                id="modal-edit-base-url"
                type="url"
                placeholder="https://api.openai.com/v1"
                aria-invalid={Boolean(errors.base_url)}
                className="mt-1.5 placeholder:text-[var(--color-text-muted)]/50"
                {...editForm.register("base_url")}
              />
              {errors.base_url ? (
                <p className="mt-1 text-xs text-[var(--status-error)]">{errors.base_url.message}</p>
              ) : null}
            </div>

            <div className="md:col-span-2">
              <Label htmlFor="modal-edit-models">
                {isId
                  ? "Model yang Diizinkan (pisahkan dengan koma)"
                  : "Allowed / Advertised Models (comma-separated)"}
              </Label>
              <Input
                id="modal-edit-models"
                type="text"
                placeholder="gpt-4o, gpt-4o-mini, claude-sonnet-4-5"
                aria-invalid={Boolean(errors.models)}
                className="mt-1.5 placeholder:text-[var(--color-text-muted)]/50"
                {...editForm.register("models")}
              />
              {errors.models ? (
                <p className="mt-1 text-xs text-[var(--status-error)]">{errors.models.message}</p>
              ) : null}
            </div>

            <div className="md:col-span-2">
              <Label htmlFor="modal-edit-headers">
                {isId
                  ? "Header Upstream Custom (satu Key: Value per baris)"
                  : "Custom Upstream Headers (one Key: Value per line)"}
              </Label>
              <textarea
                id="modal-edit-headers"
                placeholder="X-Custom-Header: value"
                aria-invalid={Boolean(errors.headers)}
                className={`${TEXTAREA_CLASS_NAME} mt-1.5 placeholder:text-[var(--color-text-muted)]/50`}
                {...editForm.register("headers")}
              />
              {errors.headers ? (
                <p className="mt-1 text-xs text-[var(--status-error)]">{errors.headers.message}</p>
              ) : null}
            </div>
          </div>

          {connection.auth_type === "api_key" ? (
            <div className="rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 p-3.5">
              <Label htmlFor="modal-edit-api-key">
                {isId ? "Rotasi API Key (opsional)" : "Rotate API Key (optional)"}
              </Label>
              <div className="relative mt-1.5">
                <Input
                  id="modal-edit-api-key"
                  type={showApiKey ? "text" : "password"}
                  autoComplete="off"
                  autoCorrect="off"
                  autoCapitalize="off"
                  spellCheck={false}
                  data-lpignore="true"
                  data-1p-ignore="true"
                  readOnly={!secretFieldsUnlocked}
                  onPointerDown={() => setSecretFieldsUnlocked(true)}
                  onFocus={() => setSecretFieldsUnlocked(true)}
                  placeholder={
                    isId
                      ? "Kosongkan untuk mempertahankan API key terenkripsi saat ini"
                      : "Leave blank to keep existing vaulted API key"
                  }
                  className="pr-10 placeholder:text-[var(--color-text-muted)]/50 [&:-webkit-autofill]:shadow-[inset_0_0_0_1000px_var(--bg-input)] [&:-webkit-autofill]:[-webkit-text-fill-color:var(--text-primary)]"
                  {...editForm.register("api_key")}
                />
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  onClick={() => setShowApiKey((prev) => !prev)}
                  aria-label={showApiKey ? "Hide API key" : "Show API key"}
                  className="absolute top-1.5 right-1.5"
                >
                  {showApiKey ? (
                    <EyeOff className="h-4 w-4" aria-hidden="true" />
                  ) : (
                    <Eye className="h-4 w-4" aria-hidden="true" />
                  )}
                </Button>
              </div>
            </div>
          ) : (
            <div className="space-y-3 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/40 p-3.5">
              <span className="text-xs font-medium text-[var(--color-text-main)]">
                {isId ? "Rotasi Token OAuth (opsional)" : "OAuth Token Rotation (optional)"}
              </span>
              <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                <div>
                  <Label htmlFor="modal-edit-access-token">
                    {isId ? "Access Token Baru" : "New Access Token"}
                  </Label>
                  <div className="relative mt-1.5">
                    <Input
                      id="modal-edit-access-token"
                      type={showAccessToken ? "text" : "password"}
                      autoComplete="off"
                      autoCorrect="off"
                      autoCapitalize="off"
                      spellCheck={false}
                      data-lpignore="true"
                      data-1p-ignore="true"
                      readOnly={!secretFieldsUnlocked}
                      onPointerDown={() => setSecretFieldsUnlocked(true)}
                      onFocus={() => setSecretFieldsUnlocked(true)}
                      placeholder={isId ? "Kosongkan jika tetap" : "Leave blank to keep existing"}
                      className="pr-10 placeholder:text-[var(--color-text-muted)]/50 [&:-webkit-autofill]:shadow-[inset_0_0_0_1000px_var(--bg-input)] [&:-webkit-autofill]:[-webkit-text-fill-color:var(--text-primary)]"
                      {...editForm.register("access_token")}
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      onClick={() => setShowAccessToken((prev) => !prev)}
                      aria-label={showAccessToken ? "Hide access token" : "Show access token"}
                      className="absolute top-1.5 right-1.5"
                    >
                      {showAccessToken ? (
                        <EyeOff className="h-4 w-4" aria-hidden="true" />
                      ) : (
                        <Eye className="h-4 w-4" aria-hidden="true" />
                      )}
                    </Button>
                  </div>
                </div>

                <div>
                  <Label htmlFor="modal-edit-refresh-token">
                    {isId ? "Refresh Token Baru" : "New Refresh Token"}
                  </Label>
                  <div className="relative mt-1.5">
                    <Input
                      id="modal-edit-refresh-token"
                      type={showRefreshToken ? "text" : "password"}
                      autoComplete="off"
                      autoCorrect="off"
                      autoCapitalize="off"
                      spellCheck={false}
                      data-lpignore="true"
                      data-1p-ignore="true"
                      readOnly={!secretFieldsUnlocked}
                      onPointerDown={() => setSecretFieldsUnlocked(true)}
                      onFocus={() => setSecretFieldsUnlocked(true)}
                      placeholder={isId ? "Kosongkan jika tetap" : "Leave blank to keep existing"}
                      className="pr-10 placeholder:text-[var(--color-text-muted)]/50 [&:-webkit-autofill]:shadow-[inset_0_0_0_1000px_var(--bg-input)] [&:-webkit-autofill]:[-webkit-text-fill-color:var(--text-primary)]"
                      {...editForm.register("refresh_token")}
                    />
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon-sm"
                      onClick={() => setShowRefreshToken((prev) => !prev)}
                      aria-label={showRefreshToken ? "Hide refresh token" : "Show refresh token"}
                      className="absolute top-1.5 right-1.5"
                    >
                      {showRefreshToken ? (
                        <EyeOff className="h-4 w-4" aria-hidden="true" />
                      ) : (
                        <Eye className="h-4 w-4" aria-hidden="true" />
                      )}
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          )}

          <label
            htmlFor="modal-edit-enabled"
            className="flex cursor-pointer items-center justify-between rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/30 px-3.5 py-2.5"
          >
            <span className="text-xs font-medium text-[var(--color-text-main)]">
              {isId
                ? "Aktifkan koneksi untuk routing gateway"
                : "Connection enabled for gateway routing"}
            </span>
            <input
              id="modal-edit-enabled"
              type="checkbox"
              className="h-4 w-4 accent-[var(--color-primary)]"
              {...editForm.register("enabled")}
            />
          </label>
        </form>
      ) : null}
    </ModalShell>
  );
}
