import { CheckCircle2, KeyRound, Search, ShieldAlert, X } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { computeUpdatedAllowedModels } from "./catalogUtils";
import { useUpdateKeyAllowedModelsMutation } from "./hooks";
import type { ApiKeyAccessItem } from "./types";

interface AssignKeysModalProps {
  open: boolean;
  onClose: () => void;
  targets: string[];
  keys: ApiKeyAccessItem[];
  keysLoading: boolean;
  keysError: string | null;
  onRetryKeys: () => void;
  isId: boolean;
}

export function AssignKeysModal({
  open,
  onClose,
  targets,
  keys,
  keysLoading,
  keysError,
  onRetryKeys,
  isId,
}: AssignKeysModalProps) {
  const updateKeyMutation = useUpdateKeyAllowedModelsMutation();

  const [mode, setMode] = useState<"add" | "remove">("add");
  const [search, setSearch] = useState("");
  const [selectedKeyIds, setSelectedKeyIds] = useState<Set<string>>(() => new Set());
  const [restrictedConfirmIds, setRestrictedConfirmIds] = useState<Set<string>>(() => new Set());
  const [applying, setApplying] = useState(false);
  const [results, setResults] = useState<Record<string, { ok: boolean; message: string }>>({});

  useEffect(() => {
    if (!open) {
      setSearch("");
      setSelectedKeyIds(new Set());
      setRestrictedConfirmIds(new Set());
      setResults({});
      setMode("add");
    }
  }, [open]);

  const filteredKeys = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return keys;
    return keys.filter(
      (k) => k.name.toLowerCase().includes(q) || k.prefix.toLowerCase().includes(q),
    );
  }, [keys, search]);

  if (!open) {
    return null;
  }

  const toggleKey = (keyId: string) => {
    setSelectedKeyIds((prev) => {
      const next = new Set(prev);
      if (next.has(keyId)) {
        next.delete(keyId);
      } else {
        next.add(keyId);
      }
      return next;
    });
  };

  const toggleRestrictedConfirm = (keyId: string) => {
    setRestrictedConfirmIds((prev) => {
      const next = new Set(prev);
      if (next.has(keyId)) {
        next.delete(keyId);
      } else {
        next.add(keyId);
      }
      return next;
    });
  };

  const selectedKeys = keys.filter((k) => selectedKeyIds.has(k.id));
  const needsRestrictedConfirm =
    mode === "add" &&
    selectedKeys.some((k) => k.modelAccessMode === "all" && !restrictedConfirmIds.has(k.id));

  const handleApply = async () => {
    if (selectedKeys.length === 0 || targets.length === 0 || applying || needsRestrictedConfirm) {
      return;
    }
    setApplying(true);
    const nextResults: Record<string, { ok: boolean; message: string }> = {};

    for (const key of selectedKeys) {
      if (mode === "remove" && key.modelAccessMode === "all") {
        nextResults[key.id] = {
          ok: true,
          message: isId ? "Tetap semua model (*)" : "No change (all models access)",
        };
        continue;
      }
      const nextAllowed = computeUpdatedAllowedModels(key, targets, mode);
      try {
        await updateKeyMutation.mutateAsync({ keyId: key.id, allowedModels: nextAllowed });
        nextResults[key.id] = {
          ok: true,
          message:
            mode === "add"
              ? isId
                ? `Diperbarui (${nextAllowed.length} aturan)`
                : `Updated (${nextAllowed.length} rules)`
              : isId
                ? `Dihapus (${nextAllowed.length} tersisa)`
                : `Removed (${nextAllowed.length} left)`,
        };
      } catch (err) {
        nextResults[key.id] = {
          ok: false,
          message: err instanceof Error ? err.message : "Update failed",
        };
      }
    }

    setResults(nextResults);
    setApplying(false);
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/55 p-4 backdrop-blur-[2px]"
      role="dialog"
      aria-modal="true"
      aria-labelledby="assign-keys-modal-title"
    >
      <div className="flex max-h-[85vh] w-full max-w-xl flex-col overflow-hidden rounded-lg border border-[var(--border-strong)] bg-[var(--bg-surface)] shadow-2xl">
        {/* Sticky Header */}
        <div className="flex items-center justify-between border-b border-[var(--border-subtle)] px-5 py-4">
          <div className="flex items-center gap-2.5">
            <KeyRound className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <div>
              <h3
                id="assign-keys-modal-title"
                className="text-sm font-semibold text-[var(--text-primary)]"
              >
                {isId ? "Tetapkan Akses Kunci API" : "Assign to API Keys"}
              </h3>
              <p className="text-xs text-[var(--text-muted)]">
                {isId
                  ? `${targets.length} item dipilih untuk diperbarui aksesnya`
                  : `${targets.length} selected ${targets.length === 1 ? "item" : "items"}`}
              </p>
            </div>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={onClose}
            aria-label={isId ? "Tutup modal" : "Close modal"}
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>

        {/* Body */}
        <div className="flex-1 space-y-4 overflow-y-auto px-5 py-4">
          {/* Selected targets preview */}
          <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3">
            <p className="mb-1.5 text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
              {isId ? "Target Dipilih" : "Selected Targets"}
            </p>
            <div className="flex flex-wrap gap-1.5">
              {targets.slice(0, 12).map((t) => (
                <Badge key={t} variant="secondary" className="font-mono text-[11px]">
                  {t}
                </Badge>
              ))}
              {targets.length > 12 ? (
                <Badge variant="outline" className="text-[11px]">
                  +{targets.length - 12} {isId ? "lainnya" : "more"}
                </Badge>
              ) : null}
            </div>
          </div>

          {/* Mode switch + Search */}
          <div className="flex flex-col gap-2.5 sm:flex-row sm:items-center sm:justify-between">
            <div className="inline-flex rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-0.5">
              <button
                type="button"
                onClick={() => setMode("add")}
                className={`rounded px-3 py-1 text-xs font-medium transition-colors ${
                  mode === "add"
                    ? "bg-[var(--accent-primary)] text-white"
                    : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
                }`}
              >
                {isId ? "Tambah Akses" : "Add Access"}
              </button>
              <button
                type="button"
                onClick={() => setMode("remove")}
                className={`rounded px-3 py-1 text-xs font-medium transition-colors ${
                  mode === "remove"
                    ? "bg-[var(--status-error)] text-white"
                    : "text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
                }`}
              >
                {isId ? "Cabut Akses" : "Remove Access"}
              </button>
            </div>

            <div className="relative flex-1 sm:max-w-xs">
              <Search
                className="pointer-events-none absolute top-1/2 left-2.5 h-3.5 w-3.5 -translate-y-1/2 text-[var(--text-muted)]"
                aria-hidden="true"
              />
              <Input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder={isId ? "Cari kunci API..." : "Search API keys..."}
                className="h-8 pl-8 text-xs"
              />
            </div>
          </div>

          {/* Keys List */}
          {keysLoading ? (
            <p className="py-6 text-center text-xs text-[var(--text-muted)]">
              {isId ? "Memuat kunci API..." : "Loading API keys..."}
            </p>
          ) : keysError ? (
            <div className="flex items-center justify-between rounded-md border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-3.5 py-2.5 text-xs text-[var(--status-error)]">
              <span>{keysError}</span>
              <Button type="button" size="xs" variant="outline" onClick={onRetryKeys}>
                {isId ? "Coba lagi" : "Retry"}
              </Button>
            </div>
          ) : filteredKeys.length === 0 ? (
            <p className="py-6 text-center text-xs text-[var(--text-muted)]">
              {isId ? "Tidak ada kunci API tersedia." : "No API keys available."}
            </p>
          ) : (
            <div className="divide-y divide-[var(--border-subtle)] rounded-md border border-[var(--border-subtle)]">
              {filteredKeys.map((key) => {
                const isSelected = selectedKeyIds.has(key.id);
                const result = results[key.id];
                const isAllMode = key.modelAccessMode === "all";

                return (
                  <div key={key.id} className="px-3.5 py-2.5">
                    <label className="flex cursor-pointer items-center justify-between gap-3">
                      <div className="flex min-w-0 items-center gap-2.5">
                        <input
                          type="checkbox"
                          checked={isSelected}
                          disabled={applying}
                          onChange={() => toggleKey(key.id)}
                          className="h-4 w-4 accent-[var(--accent-primary)]"
                        />
                        <div className="min-w-0">
                          <div className="flex flex-wrap items-center gap-1.5">
                            <span className="truncate text-xs font-semibold text-[var(--text-primary)]">
                              {key.name}
                            </span>
                            <span className="font-mono text-[11px] text-[var(--text-muted)]">
                              {key.prefix}...
                            </span>
                            {key.disabled ? (
                              <Badge variant="warning" className="text-[10px]">
                                {isId ? "Nonaktif" : "Inactive"}
                              </Badge>
                            ) : null}
                          </div>
                          <p className="mt-0.5 text-[11px] text-[var(--text-muted)]">
                            {isAllMode
                              ? isId
                                ? "Akses semua model (*)"
                                : "All models access (*)"
                              : isId
                                ? `Terbatas: ${key.allowedModels.length} model`
                                : `Restricted: ${key.allowedModels.length} models`}
                          </p>
                        </div>
                      </div>

                      {result ? (
                        <span
                          data-testid={`assign-result-${key.id}`}
                          className={`inline-flex shrink-0 items-center gap-1 text-[11px] font-medium ${
                            result.ok ? "text-[var(--status-ok)]" : "text-[var(--status-error)]"
                          }`}
                        >
                          {result.ok ? (
                            <CheckCircle2 className="h-3.5 w-3.5" aria-hidden="true" />
                          ) : (
                            <ShieldAlert className="h-3.5 w-3.5" aria-hidden="true" />
                          )}
                          <span>{result.message}</span>
                        </span>
                      ) : null}
                    </label>

                    {/* OmniRoute restricted-mode switch confirmation when adding to an 'all' key */}
                    {isSelected && isAllMode && mode === "add" ? (
                      <label
                        data-testid={`switch-restricted-${key.id}`}
                        className="mt-2 ml-6 flex cursor-pointer items-start gap-2 rounded border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-2.5 py-2 text-[11px] text-[var(--text-primary)]"
                      >
                        <input
                          type="checkbox"
                          checked={restrictedConfirmIds.has(key.id)}
                          disabled={applying}
                          onChange={() => toggleRestrictedConfirm(key.id)}
                          className="mt-0.5 h-3.5 w-3.5 accent-[var(--accent-primary)]"
                        />
                        <span>
                          {isId
                            ? "Konfirmasi ubah kunci ini dari akses semua model (*) menjadi mode terbatas hanya untuk item yang dipilih."
                            : "Confirm switching this key from all-models (*) to restricted mode with only the selected items."}
                        </span>
                      </label>
                    ) : null}
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Sticky Footer */}
        <div className="flex items-center justify-between border-t border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-5 py-3">
          <span className="text-xs text-[var(--text-muted)]">
            {isId ? `${selectedKeyIds.size} kunci dipilih` : `${selectedKeyIds.size} keys selected`}
          </span>
          <div className="flex items-center gap-2">
            <Button type="button" variant="secondary" size="sm" onClick={onClose}>
              {isId ? "Tutup" : "Close"}
            </Button>
            <Button
              type="button"
              size="sm"
              data-testid="assign-apply-btn"
              disabled={selectedKeyIds.size === 0 || applying || needsRestrictedConfirm}
              onClick={() => void handleApply()}
            >
              {applying ? (isId ? "Menerapkan..." : "Applying...") : isId ? "Terapkan" : "Apply"}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
