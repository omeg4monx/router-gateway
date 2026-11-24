import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertCircle,
  Check,
  ChevronDown,
  Code2,
  Copy,
  KeyRound,
  Terminal,
  Waypoints,
  X,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { cn } from "@/shared/lib/utils";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  generateKeySnippet,
  type KeySnippetFormat,
  type KeySnippetLang,
  toDatetimeLocalValue,
} from "./endpointUtils";
import {
  type ApiKeySummary,
  type CreateKeyFormInput,
  type CreateKeyFormOutput,
  createKeySchema,
  DEFAULT_CREATE_KEY_VALUES,
} from "./schemas";

interface KeyFormModalProps {
  isOpen: boolean;
  isId: boolean;
  mode: "create" | "edit";
  initialKey: ApiKeySummary | null;
  isPending: boolean;
  errorMessage: string | null;
  onClose: () => void;
  onSubmit: (values: CreateKeyFormOutput, editingKeyId?: string) => void;
}

export function KeyFormModal({
  isOpen,
  isId,
  mode,
  initialKey,
  isPending,
  errorMessage,
  onClose,
  onSubmit,
}: KeyFormModalProps) {
  const [showAdvancedLimits, setShowAdvancedLimits] = useState(false);
  const nameInputRef = useRef<HTMLInputElement | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<CreateKeyFormInput, unknown, CreateKeyFormOutput>({
    resolver: zodResolver(createKeySchema),
    defaultValues: DEFAULT_CREATE_KEY_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  useEffect(() => {
    if (!isOpen) {
      return;
    }
    if (mode === "edit" && initialKey) {
      reset({
        name: initialKey.name,
        allowed_models: initialKey.allowed_models.join(", "),
        tags: initialKey.tags.join(", "),
        rate_limit_rpm: String(initialKey.rate_limit_rpm ?? 0),
        quota_usd: String(initialKey.quota_usd ?? 0),
        expires_at: toDatetimeLocalValue(initialKey.expires_at),
      });
      setShowAdvancedLimits(
        initialKey.allowed_models.length > 0 ||
          initialKey.rate_limit_rpm > 0 ||
          initialKey.quota_usd > 0 ||
          Boolean(initialKey.expires_at),
      );
    } else {
      reset(DEFAULT_CREATE_KEY_VALUES);
      setShowAdvancedLimits(false);
    }

    const focusTimer = window.setTimeout(() => nameInputRef.current?.focus(), 0);
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape" && !isPending) {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.clearTimeout(focusTimer);
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen, mode, initialKey, reset, isPending, onClose]);

  if (!isOpen) {
    return null;
  }

  const nameField = register("name");

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        tabIndex={-1}
        aria-label={isId ? "Tutup modal" : "Close dialog"}
        onClick={() => {
          if (!isPending) onClose();
        }}
        className="fixed inset-0 bg-black/55 backdrop-blur-[2px]"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="key-form-dialog-title"
        className="relative z-10 w-full max-w-md rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl"
      >
        <div className="flex items-center justify-between border-b border-[var(--color-border-subtle)] px-5 py-4">
          <div className="flex items-center gap-2">
            <KeyRound className="h-4 w-4 text-[var(--color-primary)]" aria-hidden="true" />
            <h3
              id="key-form-dialog-title"
              className="text-sm font-semibold text-[var(--color-text-main)]"
            >
              {mode === "edit"
                ? isId
                  ? "Edit Kunci API"
                  : "Edit API Key"
                : isId
                  ? "Buat Kunci API Baru"
                  : "Create API Key"}
            </h3>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            disabled={isPending}
            aria-label={isId ? "Tutup modal" : "Close dialog"}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        <form
          onSubmit={handleSubmit((values) => onSubmit(values, initialKey?.id))}
          noValidate
          className="space-y-4 p-5"
        >
          {errorMessage ? (
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>
                  {mode === "edit"
                    ? isId
                      ? "Gagal memperbarui kunci"
                      : "Failed to update API key"
                    : isId
                      ? "Gagal membuat kunci"
                      : "Failed to create API key"}
                </AlertTitle>
                <AlertDescription className="font-mono">{errorMessage}</AlertDescription>
              </div>
            </Alert>
          ) : null}

          <div>
            <Label htmlFor="key-name">{isId ? "Nama Kunci" : "Key Name"}</Label>
            <Input
              id="key-name"
              type="text"
              placeholder={
                isId ? "contoh: cli-lokal atau backend-prod" : "e.g. local-cli or prod-backend"
              }
              aria-invalid={Boolean(errors.name)}
              className="mt-1.5"
              {...nameField}
              ref={(node) => {
                nameField.ref(node);
                nameInputRef.current = node;
              }}
            />
            {errors.name ? (
              <p className="mt-1 text-xs text-[var(--status-error)]">{errors.name.message}</p>
            ) : null}
          </div>

          <div className="border-t border-[var(--color-border-subtle)] pt-3">
            <button
              type="button"
              onClick={() => setShowAdvancedLimits((prev) => !prev)}
              aria-expanded={showAdvancedLimits}
              className="flex w-full cursor-pointer items-center justify-between text-xs font-medium text-[var(--color-text-muted)] transition-colors hover:text-[var(--color-text-main)]"
            >
              <span>
                {isId ? "Batasan & akses model (opsional)" : "Limits & model access (optional)"}
              </span>
              <ChevronDown
                className={cn(
                  "h-3.5 w-3.5 transition-transform",
                  showAdvancedLimits && "rotate-180",
                )}
                aria-hidden="true"
              />
            </button>

            {showAdvancedLimits ? (
              <div className="mt-3 space-y-3">
                <div>
                  <Label htmlFor="key-allowed-models">
                    {isId ? "Model yang Diizinkan" : "Allowed Models"}
                  </Label>
                  <Input
                    id="key-allowed-models"
                    type="text"
                    placeholder={
                      isId
                        ? "Kosongkan untuk semua model (pisahkan koma)"
                        : "Leave blank for all models (comma-separated)"
                    }
                    aria-invalid={Boolean(errors.allowed_models)}
                    className="mt-1.5"
                    {...register("allowed_models")}
                  />
                  {errors.allowed_models ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {errors.allowed_models.message}
                    </p>
                  ) : null}
                </div>

                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                  <div>
                    <Label htmlFor="key-rpm">{isId ? "Batas RPM" : "Rate Limit (RPM)"}</Label>
                    <Input
                      id="key-rpm"
                      type="number"
                      min={0}
                      step={1}
                      placeholder="0"
                      aria-invalid={Boolean(errors.rate_limit_rpm)}
                      className="mt-1.5 tabular-nums"
                      {...register("rate_limit_rpm")}
                    />
                    {errors.rate_limit_rpm ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.rate_limit_rpm.message}
                      </p>
                    ) : null}
                  </div>

                  <div>
                    <Label htmlFor="key-quota">{isId ? "Kuota (USD)" : "Quota (USD)"}</Label>
                    <Input
                      id="key-quota"
                      type="number"
                      min={0}
                      step="0.01"
                      placeholder="0"
                      aria-invalid={Boolean(errors.quota_usd)}
                      className="mt-1.5 tabular-nums"
                      {...register("quota_usd")}
                    />
                    {errors.quota_usd ? (
                      <p className="mt-1 text-xs text-[var(--status-error)]">
                        {errors.quota_usd.message}
                      </p>
                    ) : null}
                  </div>
                </div>

                <div>
                  <Label htmlFor="key-expires">{isId ? "Tanggal Kedaluwarsa" : "Expires At"}</Label>
                  <Input
                    id="key-expires"
                    type="datetime-local"
                    aria-invalid={Boolean(errors.expires_at)}
                    className="mt-1.5"
                    {...register("expires_at")}
                  />
                  {errors.expires_at ? (
                    <p className="mt-1 text-xs text-[var(--status-error)]">
                      {errors.expires_at.message}
                    </p>
                  ) : null}
                </div>
              </div>
            ) : null}
          </div>

          <div className="flex items-center justify-end gap-2 pt-2">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              onClick={onClose}
              disabled={isPending}
            >
              {isId ? "Batal" : "Cancel"}
            </Button>
            <Button type="submit" size="sm" disabled={isPending}>
              <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
              <span>
                {isPending
                  ? isId
                    ? "Menyimpan..."
                    : "Saving..."
                  : mode === "edit"
                    ? isId
                      ? "Simpan Perubahan"
                      : "Save Changes"
                    : isId
                      ? "Buat Key"
                      : "Create Key"}
              </span>
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}

interface KeyReadyModalProps {
  isId: boolean;
  secretState: {
    action: "created" | "rotated";
    keyName: string;
    secret: string;
  } | null;
  activeOrigin: string;
  defaultModel: string;
  onClose: () => void;
}

export function KeyReadyModal({
  isId,
  secretState,
  activeOrigin,
  defaultModel,
  onClose,
}: KeyReadyModalProps) {
  const [copiedTarget, setCopiedTarget] = useState<"key" | "curl" | null>(null);
  const timerRef = useRef<number | null>(null);

  useEffect(() => {
    return () => {
      if (timerRef.current !== null) {
        window.clearTimeout(timerRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (!secretState) {
      return;
    }
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [secretState, onClose]);

  if (!secretState) {
    return null;
  }

  const copyText = async (text: string, kind: "key" | "curl") => {
    try {
      await navigator.clipboard.writeText(text);
      setCopiedTarget(kind);
      if (timerRef.current !== null) {
        window.clearTimeout(timerRef.current);
      }
      timerRef.current = window.setTimeout(() => {
        timerRef.current = null;
        setCopiedTarget(null);
      }, 2200);
    } catch {
      // Fallback: user can select text directly
    }
  };

  const curlCmd = generateKeySnippet({
    format: "openai",
    lang: "curl",
    baseUrl: activeOrigin,
    apiKey: secretState.secret,
    model: defaultModel,
  });

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        tabIndex={-1}
        aria-label={isId ? "Tutup modal" : "Close dialog"}
        onClick={onClose}
        className="fixed inset-0 bg-black/55 backdrop-blur-[2px]"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="key-ready-dialog-title"
        className="relative z-10 w-full max-w-lg rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl"
      >
        <div className="flex items-center justify-between border-b border-[var(--color-border-subtle)] px-5 py-4">
          <h3
            id="key-ready-dialog-title"
            className="text-sm font-semibold text-[var(--color-text-main)]"
          >
            {secretState.action === "rotated"
              ? isId
                ? "Kunci API Dirotasi"
                : "API Key Rotated"
              : isId
                ? "Kunci API Siap Dipakai"
                : "API Key Ready"}
          </h3>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            aria-label={isId ? "Tutup" : "Close"}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        <div className="space-y-4 p-5">
          <div className="rounded-xl border border-[var(--color-success-border)] bg-[var(--color-success-bg)] p-3.5 text-xs text-[var(--color-success-fg)]">
            <p className="font-semibold">
              {isId
                ? `"${secretState.keyName}" siap digunakan.`
                : `"${secretState.keyName}" is ready to use.`}
            </p>
            <p className="mt-0.5 opacity-90">
              {isId
                ? "Salin secret sekarang — hanya ditampilkan sekali dan tidak dapat diambil kembali dari server."
                : "Copy this secret now — it is only shown once and cannot be recovered after closing."}
            </p>
          </div>

          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <input
              type="text"
              readOnly
              value={secretState.secret}
              aria-label={isId ? "Secret kunci API" : "API key secret"}
              className="min-w-0 flex-1 rounded-lg border border-[var(--color-border)] bg-[var(--color-field)] px-3 py-2 font-mono text-xs text-[var(--color-text-main)] select-all outline-none"
            />
            <Button
              type="button"
              variant="secondary"
              size="sm"
              className="shrink-0"
              onClick={() => void copyText(secretState.secret, "key")}
            >
              {copiedTarget === "key" ? (
                <>
                  <Check className="h-3.5 w-3.5 text-[var(--color-success)]" aria-hidden="true" />
                  <span>{isId ? "Tersalin" : "Copied!"}</span>
                </>
              ) : (
                <>
                  <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{isId ? "Salin Key" : "Copy Key"}</span>
                </>
              )}
            </Button>
          </div>

          <Button
            type="button"
            variant="outline"
            size="sm"
            className="w-full"
            onClick={() => void copyText(curlCmd, "curl")}
          >
            {copiedTarget === "curl" ? (
              <>
                <Check className="h-3.5 w-3.5 text-[var(--color-success)]" aria-hidden="true" />
                <span>{isId ? "Perintah cURL Tersalin!" : "cURL Command Copied!"}</span>
              </>
            ) : (
              <>
                <Terminal className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Salin Perintah Uji cURL" : "Copy Ready-to-Run cURL"}</span>
              </>
            )}
          </Button>
        </div>

        <div className="flex justify-end border-t border-[var(--color-border-subtle)] px-5 py-3">
          <Button type="button" size="sm" onClick={onClose}>
            {isId ? "Selesai" : "Done"}
          </Button>
        </div>
      </div>
    </div>
  );
}

interface KeySnippetModalProps {
  isId: boolean;
  apiKey: ApiKeySummary | null;
  revealedSecret?: string;
  activeOrigin: string;
  activeRouteLabel: string;
  defaultModel: string;
  onClose: () => void;
}

const MODAL_LANGS: ReadonlyArray<{ id: KeySnippetLang; label: string }> = [
  { id: "curl", label: "cURL" },
  { id: "python", label: "Python" },
  { id: "node", label: "Node.js" },
  { id: "go", label: "Go" },
];

export function KeySnippetModal({
  isId,
  apiKey,
  revealedSecret,
  activeOrigin,
  activeRouteLabel,
  defaultModel,
  onClose,
}: KeySnippetModalProps) {
  const [format, setFormat] = useState<KeySnippetFormat>("openai");
  const [lang, setLang] = useState<KeySnippetLang>("curl");
  const [copied, setCopied] = useState(false);
  const timerRef = useRef<number | null>(null);

  useEffect(() => {
    return () => {
      if (timerRef.current !== null) {
        window.clearTimeout(timerRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (!apiKey) {
      return;
    }
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [apiKey, onClose]);

  const code = useMemo(() => {
    if (!apiKey) return "";
    const resolvedModel = apiKey.allowed_models[0] || defaultModel;
    return generateKeySnippet({
      format,
      lang,
      baseUrl: activeOrigin,
      apiKey: revealedSecret || `${apiKey.key_prefix}...`,
      model: resolvedModel,
    });
  }, [apiKey, revealedSecret, format, lang, activeOrigin, defaultModel]);

  if (!apiKey) {
    return null;
  }

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      if (timerRef.current !== null) {
        window.clearTimeout(timerRef.current);
      }
      timerRef.current = window.setTimeout(() => {
        timerRef.current = null;
        setCopied(false);
      }, 2200);
    } catch {
      // ignore
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        tabIndex={-1}
        aria-label={isId ? "Tutup modal" : "Close dialog"}
        onClick={onClose}
        className="fixed inset-0 bg-black/55 backdrop-blur-[2px]"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="key-snippet-dialog-title"
        className="relative z-10 w-full max-w-xl rounded-2xl border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl"
      >
        <div className="flex items-center justify-between border-b border-[var(--color-border-subtle)] px-5 py-4">
          <div className="flex min-w-0 items-center gap-2">
            <Code2 className="h-4 w-4 shrink-0 text-[var(--color-primary)]" aria-hidden="true" />
            <h3
              id="key-snippet-dialog-title"
              className="truncate text-sm font-semibold text-[var(--color-text-main)]"
            >
              {isId ? "Cuplikan Kode" : "Code Snippet"} — {apiKey.name}
            </h3>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            aria-label={isId ? "Tutup" : "Close"}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        <div className="space-y-3.5 p-5">
          <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="inline-flex rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] p-0.5">
              {(
                [
                  { id: "openai", label: "OpenAI (/v1)" },
                  { id: "anthropic", label: "Anthropic (/v1/messages)" },
                ] as const
              ).map((item) => (
                <button
                  key={item.id}
                  type="button"
                  onClick={() => setFormat(item.id)}
                  className={cn(
                    "cursor-pointer rounded-md px-2.5 py-1 text-[11px] font-medium transition-colors",
                    format === item.id
                      ? "bg-[var(--color-primary)] text-[var(--color-on-primary)]"
                      : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                  )}
                >
                  {item.label}
                </button>
              ))}
            </div>

            <span className="inline-flex items-center gap-1.5 rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] px-2.5 py-1 font-mono text-[11px] text-[var(--color-text-muted)]">
              <Waypoints className="h-3.5 w-3.5 text-[var(--color-primary)]" aria-hidden="true" />
              <span>
                {isId ? "Rute" : "Route"}: {activeRouteLabel}
              </span>
            </span>
          </div>

          <div className="inline-flex rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)] p-0.5">
            {MODAL_LANGS.map((item) => (
              <button
                key={item.id}
                type="button"
                onClick={() => setLang(item.id)}
                className={cn(
                  "cursor-pointer rounded-md px-2.5 py-1 font-mono text-[11px] font-medium transition-colors",
                  lang === item.id
                    ? "bg-[var(--color-surface)] text-[var(--color-text-main)] shadow-xs"
                    : "text-[var(--color-text-muted)] hover:text-[var(--color-text-main)]",
                )}
              >
                {item.label}
              </button>
            ))}
          </div>

          <pre className="custom-scrollbar max-h-72 overflow-x-auto whitespace-pre-wrap break-all rounded-xl border border-[var(--color-border-subtle)] bg-[var(--color-bg-alt)] p-4 font-mono text-xs leading-relaxed text-[var(--color-text-main)]">
            {code}
          </pre>

          {!revealedSecret ? (
            <p className="text-[11px] text-[var(--color-text-subtle)]">
              {isId
                ? "Ganti placeholder prefix dengan secret lengkap yang Anda simpan saat pembuatan kunci."
                : "Replace the masked prefix placeholder with your full saved API key secret."}
            </p>
          ) : null}
        </div>

        <div className="flex items-center justify-end gap-2 border-t border-[var(--color-border-subtle)] px-5 py-3">
          <Button type="button" variant="ghost" size="sm" onClick={onClose}>
            {isId ? "Tutup" : "Close"}
          </Button>
          <Button type="button" size="sm" onClick={() => void handleCopy()}>
            {copied ? (
              <>
                <Check className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Tersalin!" : "Copied!"}</span>
              </>
            ) : (
              <>
                <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                <span>{isId ? "Salin ke Clipboard" : "Copy to Clipboard"}</span>
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}
