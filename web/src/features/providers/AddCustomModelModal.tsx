import { Activity, AlertTriangle, Bot, CheckCircle2, Plus, RefreshCw, X } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { runSingleModelTest } from "../models/hooks";
import type { CatalogTestResult } from "../models/types";
import { mergeManualModels, normalizeManualModelInput } from "./customModelUtils";
import { useUpdateProviderConnectionMutation } from "./hooks";
import { MAX_SLICE_ITEM_RUNES, MAX_SLICE_ITEMS, type ProviderConnectionSummary } from "./schemas";

interface AddCustomModelModalProps {
  open: boolean;
  onClose: () => void;
  providerId: string;
  providerName: string;
  providerAlias?: string;
  connections: ProviderConnectionSummary[];
  isId: boolean;
  onAdded: (message: string) => void;
}

const ALL_ENABLED_TARGET = "__all_enabled__";

export function AddCustomModelModal({
  open,
  onClose,
  providerId,
  providerName,
  providerAlias,
  connections,
  isId,
  onAdded,
}: AddCustomModelModalProps) {
  const updateConnMutation = useUpdateProviderConnectionMutation();

  const [rawInput, setRawInput] = useState("");
  const [targetSelection, setTargetSelection] = useState<string>(ALL_ENABLED_TARGET);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<CatalogTestResult | null>(null);
  const testControllerRef = useRef<AbortController | null>(null);

  const enabledConnections = useMemo(
    () => connections.filter((conn) => conn.enabled),
    [connections],
  );

  const providerPrefixes = useMemo(() => {
    const prefixes = [providerId];
    if (providerAlias && providerAlias !== providerId) {
      prefixes.push(providerAlias);
    }
    return prefixes;
  }, [providerId, providerAlias]);

  const parsedModels = useMemo(
    () => normalizeManualModelInput(rawInput, providerPrefixes),
    [rawInput, providerPrefixes],
  );

  const primaryModelPrefix = providerAlias || providerId;
  const singleModelId = parsedModels.length === 1 ? parsedModels[0] : null;
  const hasPrefixStripped =
    rawInput.trim().length > 0 && parsedModels.length === 1 && parsedModels[0] !== rawInput.trim();

  const firstConnection = connections[0];
  const targetConnections = useMemo((): ProviderConnectionSummary[] => {
    if (!firstConnection) {
      return [];
    }
    if (connections.length === 1) {
      return [firstConnection];
    }
    if (targetSelection === ALL_ENABLED_TARGET) {
      return enabledConnections.length > 0 ? enabledConnections : [firstConnection];
    }
    const found = connections.find((c) => c.id === targetSelection);
    return found ? [found] : [firstConnection];
  }, [connections, enabledConnections, firstConnection, targetSelection]);

  useEffect(() => {
    if (!open) {
      testControllerRef.current?.abort();
      testControllerRef.current = null;
      return;
    }
    setRawInput("");
    setTargetSelection(
      enabledConnections.length > 0
        ? ALL_ENABLED_TARGET
        : (connections[0]?.id ?? ALL_ENABLED_TARGET),
    );
    setSubmitError(null);
    setSaving(false);
    setTesting(false);
    setTestResult(null);

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
      testControllerRef.current?.abort();
      testControllerRef.current = null;
    };
  }, [open, onClose, connections, enabledConnections.length]);

  if (!open) {
    return null;
  }

  const handleTestSingleModel = async () => {
    if (!singleModelId || enabledConnections.length === 0) {
      return;
    }
    testControllerRef.current?.abort();
    const controller = new AbortController();
    testControllerRef.current = controller;
    setTesting(true);
    setTestResult(null);
    try {
      const result = await runSingleModelTest(providerId, singleModelId, controller.signal);
      if (!controller.signal.aborted) {
        setTestResult(result);
      }
    } finally {
      if (testControllerRef.current === controller) {
        testControllerRef.current = null;
        setTesting(false);
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitError(null);

    if (targetConnections.length === 0) {
      setSubmitError(
        isId
          ? "Tambahkan minimal satu koneksi provider terlebih dahulu sebelum menambahkan model."
          : "Add at least one provider connection before adding custom models.",
      );
      return;
    }

    if (parsedModels.length === 0) {
      setSubmitError(
        isId ? "Masukkan minimal satu Model ID yang valid." : "Enter at least one valid model ID.",
      );
      return;
    }

    const tooLong = parsedModels.find((id) => Array.from(id).length > MAX_SLICE_ITEM_RUNES);
    if (tooLong) {
      setSubmitError(
        isId
          ? `Model ID "${tooLong}" melebihi batas ${MAX_SLICE_ITEM_RUNES} karakter.`
          : `Model ID "${tooLong}" exceeds ${MAX_SLICE_ITEM_RUNES} characters.`,
      );
      return;
    }

    setSaving(true);
    try {
      let totalNewCount = 0;
      for (const conn of targetConnections) {
        const { merged, added } = mergeManualModels(conn.models, parsedModels);
        if (merged.length > MAX_SLICE_ITEMS) {
          throw new Error(
            isId
              ? `Koneksi "${conn.name}" melebihi batas maksimum ${MAX_SLICE_ITEMS} model.`
              : `Connection "${conn.name}" exceeds the maximum of ${MAX_SLICE_ITEMS} models.`,
          );
        }
        if (added.length === 0) {
          continue;
        }
        await updateConnMutation.mutateAsync({
          id: conn.id,
          body: { models: merged },
        });
        totalNewCount = Math.max(totalNewCount, added.length);
      }

      if (totalNewCount === 0) {
        setSubmitError(
          isId
            ? "Model tersebut sudah terdaftar pada koneksi yang dipilih."
            : "All entered model IDs are already registered on the selected connection(s).",
        );
        return;
      }

      const label = parsedModels.length === 1 ? `"${parsedModels[0]}"` : `${totalNewCount} model`;
      onAdded(
        isId
          ? `Model manual ${label} berhasil ditambahkan dan langsung siap digunakan.`
          : `Custom model ${parsedModels.length === 1 ? `"${parsedModels[0]}"` : `(${totalNewCount})`} added and ready to route.`,
      );
      onClose();
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Failed to save custom model.");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center p-0 sm:items-center sm:p-6">
      <button
        type="button"
        tabIndex={-1}
        aria-label="Close modal backdrop"
        onClick={onClose}
        className="fixed inset-0 bg-black/55 backdrop-blur-[2px]"
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-label={
          isId ? `Tambah Model Manual — ${providerName}` : `Add Custom Model — ${providerName}`
        }
        className="relative z-10 flex max-h-[90dvh] w-full max-w-lg flex-col overflow-hidden rounded-t-2xl border border-[var(--color-border)] bg-[var(--color-surface)] shadow-2xl sm:rounded-2xl"
      >
        <div className="flex shrink-0 items-start justify-between gap-3 border-b border-[var(--color-border-subtle)] px-5 py-4">
          <div className="flex min-w-0 items-center gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-[var(--color-border)] bg-[var(--color-surface-2)] text-[var(--color-primary)]">
              <Bot className="h-5 w-5" aria-hidden="true" />
            </div>
            <div className="min-w-0">
              <h3 className="truncate text-base font-semibold tracking-tight text-[var(--color-text-main)]">
                {isId ? "Tambah Model Manual" : "Add Custom Model"}
              </h3>
              <p className="mt-0.5 truncate text-xs text-[var(--color-text-muted)]">
                {isId
                  ? `Daftarkan model kustom untuk ${providerName} agar langsung bisa dipakai di gateway & katalog.`
                  : `Register custom model ID(s) on ${providerName} for immediate gateway routing & catalog access.`}
              </p>
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            aria-label={isId ? "Tutup modal" : "Close modal"}
            className="flex h-8 w-8 shrink-0 cursor-pointer items-center justify-center rounded-lg text-[var(--color-text-muted)] transition-colors hover:bg-[var(--color-surface-2)] hover:text-[var(--color-text-main)]"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </button>
        </div>

        <form
          onSubmit={(e) => void handleSubmit(e)}
          className="flex flex-1 flex-col overflow-hidden"
        >
          <div className="space-y-4 overflow-y-auto px-5 py-4">
            <div className="space-y-1.5">
              <Label htmlFor="manual-model-id-input">
                {isId
                  ? "Model ID (bisa pisah koma untuk banyak model)"
                  : "Model ID (comma-separated for multiple)"}
              </Label>
              <div className="flex items-center gap-2">
                <Input
                  id="manual-model-id-input"
                  value={rawInput}
                  onChange={(e) => {
                    setRawInput(e.target.value);
                    setTestResult(null);
                    setSubmitError(null);
                  }}
                  placeholder={
                    isId
                      ? "contoh: gpt-5-mini, claude-opus-4-6"
                      : "e.g. gpt-5-mini, claude-opus-4-6"
                  }
                  className="font-mono text-xs"
                  autoFocus
                />
                {singleModelId && enabledConnections.length > 0 ? (
                  <Button
                    type="button"
                    size="sm"
                    variant="secondary"
                    disabled={testing || saving}
                    onClick={() => void handleTestSingleModel()}
                    className="shrink-0"
                  >
                    {testing ? (
                      <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
                    ) : (
                      <Activity className="h-3.5 w-3.5" aria-hidden="true" />
                    )}
                    <span>
                      {testing ? (isId ? "Menguji..." : "Testing...") : isId ? "Tes" : "Test"}
                    </span>
                  </Button>
                ) : null}
              </div>

              {hasPrefixStripped && singleModelId ? (
                <p className="text-[11px] text-[var(--color-text-muted)]">
                  {isId ? "Dikirim ke provider sebagai:" : "Sent to provider as:"}{" "}
                  <code className="rounded bg-[var(--color-surface-2)] px-1.5 py-0.5 font-mono text-[var(--color-text-main)]">
                    {singleModelId}
                  </code>
                </p>
              ) : parsedModels.length > 0 ? (
                <p className="text-[11px] text-[var(--color-text-muted)]">
                  {isId ? "Akan terdaftar sebagai:" : "Will be callable as:"}{" "}
                  <code className="rounded bg-[var(--color-surface-2)] px-1.5 py-0.5 font-mono text-[var(--color-text-main)]">
                    {parsedModels.length === 1
                      ? `${primaryModelPrefix}/${parsedModels[0]}`
                      : `${parsedModels.length} model (${primaryModelPrefix}/${parsedModels[0]}, ...)`}
                  </code>
                </p>
              ) : (
                <p className="text-[11px] text-[var(--color-text-muted)]">
                  {isId
                    ? `Prefix "${primaryModelPrefix}/" akan otomatis dibersihkan jika Anda menempel ID lengkap.`
                    : `Leading "${primaryModelPrefix}/" is automatically stripped if you paste a qualified ID.`}
                </p>
              )}
            </div>

            {testResult ? (
              <div
                className={cn(
                  "flex flex-wrap items-center justify-between gap-2 rounded-xl border px-3.5 py-2.5 text-xs",
                  testResult.status === "error"
                    ? "border-[var(--status-error)]/35 bg-[var(--status-error-bg)] text-[var(--status-error)]"
                    : "border-[var(--status-ok)]/35 bg-[var(--status-ok-bg)] text-[var(--color-text-main)]",
                )}
              >
                <div className="flex min-w-0 items-center gap-2">
                  {testResult.status === "error" ? (
                    <AlertTriangle className="h-4 w-4 shrink-0" aria-hidden="true" />
                  ) : (
                    <CheckCircle2
                      className="h-4 w-4 shrink-0 text-[var(--status-ok)]"
                      aria-hidden="true"
                    />
                  )}
                  <Badge
                    variant={
                      testResult.status === "ok"
                        ? "success"
                        : testResult.status === "slow"
                          ? "warning"
                          : "destructive"
                    }
                    className="tabular-nums"
                  >
                    {testResult.status === "error"
                      ? `FAIL${testResult.statusCode ? ` ${testResult.statusCode}` : ""}`
                      : `${testResult.status.toUpperCase()} · ${testResult.latencyMs ?? 0}ms`}
                  </Badge>
                  {testResult.error ? (
                    <span className="truncate font-mono text-[11px]" title={testResult.error}>
                      {testResult.error}
                    </span>
                  ) : (
                    <span className="text-[var(--color-text-muted)]">
                      {isId ? "Model merespons dengan baik" : "Model responded successfully"}
                    </span>
                  )}
                </div>
              </div>
            ) : null}

            {connections.length > 1 ? (
              <div className="space-y-1.5">
                <Label htmlFor="manual-model-target-select">
                  {isId ? "Koneksi Target" : "Target Connection"}
                </Label>
                <select
                  id="manual-model-target-select"
                  value={targetSelection}
                  onChange={(e) => setTargetSelection(e.target.value)}
                  className="flex h-9 w-full rounded-lg border border-[var(--color-border)] bg-[var(--color-field)] px-3 py-1.5 text-xs text-[var(--color-text-main)] transition-colors focus:border-[var(--color-primary)] focus:outline-none"
                >
                  {enabledConnections.length > 0 ? (
                    <option value={ALL_ENABLED_TARGET}>
                      {isId
                        ? `Semua koneksi aktif (${enabledConnections.length})`
                        : `All enabled connections (${enabledConnections.length})`}
                    </option>
                  ) : null}
                  {connections.map((conn) => (
                    <option key={conn.id} value={conn.id}>
                      {conn.name} (
                      {conn.enabled ? (isId ? "aktif" : "enabled") : isId ? "nonaktif" : "disabled"}
                      )
                    </option>
                  ))}
                </select>
              </div>
            ) : firstConnection ? (
              <div className="rounded-lg border border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/50 px-3 py-2 text-xs text-[var(--color-text-muted)]">
                {isId ? "Disimpan ke koneksi:" : "Saving to connection:"}{" "}
                <span className="font-medium text-[var(--color-text-main)]">
                  {firstConnection.name}
                </span>
              </div>
            ) : (
              <div className="flex items-start gap-2 rounded-lg border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-3 py-2.5 text-xs text-[var(--color-text-main)]">
                <AlertTriangle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                  aria-hidden="true"
                />
                <span>
                  {isId
                    ? "Provider ini belum memiliki koneksi. Tambahkan API Key atau OAuth terlebih dahulu agar model dapat disimpan dan digunakan."
                    : "This provider has no connections yet. Add an API key or OAuth connection first so custom models can be saved and routed."}
                </span>
              </div>
            )}

            {submitError ? (
              <div className="flex items-start gap-2 rounded-lg border border-[var(--status-error)]/35 bg-[var(--status-error-bg)] px-3 py-2.5 text-xs text-[var(--status-error)]">
                <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
                <span>{submitError}</span>
              </div>
            ) : null}
          </div>

          <div className="flex shrink-0 flex-col-reverse gap-2 border-t border-[var(--color-border-subtle)] bg-[var(--color-surface-2)]/35 px-5 py-3.5 sm:flex-row sm:items-center sm:justify-end">
            <Button
              type="button"
              variant="secondary"
              size="sm"
              onClick={onClose}
              disabled={saving}
              className="w-full sm:w-auto"
            >
              {isId ? "Batal" : "Cancel"}
            </Button>
            <Button
              type="submit"
              size="sm"
              disabled={saving || parsedModels.length === 0 || targetConnections.length === 0}
              className="w-full sm:w-auto"
            >
              {saving ? (
                <RefreshCw className="h-3.5 w-3.5 animate-spin" aria-hidden="true" />
              ) : (
                <Plus className="h-3.5 w-3.5" aria-hidden="true" />
              )}
              <span>
                {saving
                  ? isId
                    ? "Menyimpan..."
                    : "Saving..."
                  : parsedModels.length > 1
                    ? isId
                      ? `Tambah ${parsedModels.length} Model`
                      : `Add ${parsedModels.length} Models`
                    : isId
                      ? "Tambah Model"
                      : "Add Model"}
              </span>
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
