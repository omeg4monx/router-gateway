import { AlertCircle, CheckCircle2, RotateCcw, Trash2, Zap } from "lucide-react";
import { memo, useEffect, useState } from "react";
import { type UILocale, useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Skeleton } from "@/shared/ui/skeleton";
import {
  useClearHeadroomCacheMutation,
  useResetTokenSaverStatsMutation,
  useTokenSaverConfig,
  useTokenSaverStats,
  useUpdateTokenSaverConfigMutation,
} from "./hooks";
import {
  HISTORY_LIMIT_PRESETS,
  matchHistoryPreset,
  TOKEN_SAVER_LEVELS,
  type TokenSaverConfigFormValues,
  type TokenSaverConfigPatch,
  type TokenSaverLevel,
} from "./schemas";

interface ToggleSwitchProps {
  id: string;
  checked: boolean;
  disabled?: boolean;
  label: string;
  onChange: (next: boolean) => void;
}

function ToggleSwitch({ id, checked, disabled, label, onChange }: ToggleSwitchProps) {
  return (
    <button
      id={id}
      type="button"
      role="switch"
      aria-checked={checked}
      aria-label={label}
      disabled={disabled}
      onClick={() => onChange(!checked)}
      className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent-primary)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--bg-surface)] disabled:cursor-not-allowed disabled:opacity-50 ${
        checked
          ? "border-[var(--accent-primary)] bg-[var(--accent-primary)]"
          : "border-[var(--border-strong)] bg-[var(--bg-surface-raised)]"
      }`}
    >
      <span
        className={`pointer-events-none block h-4 w-4 rounded-full transition-transform ${
          checked
            ? "translate-x-5 bg-[var(--accent-primary-ink)]"
            : "translate-x-1 bg-[var(--text-muted)]"
        }`}
      />
    </button>
  );
}

const DEFAULT_CONFIG: TokenSaverConfigFormValues = {
  enabled: false,
  strategy: "whitespace",
  max_history_messages: 40,
  max_tool_output_chars: 8000,
  headroom_enabled: true,
  prefix_depth: 2,
  max_cache_entries: 1024,
  cache_ttl_seconds: 1800,
  inject_cache_control: true,
  caveman_enabled: false,
  caveman_level: "full",
  ponytail_enabled: false,
  ponytail_level: "full",
  trim_history_enabled: false,
};

const SKELETON_ROWS = ["rtk", "headroom", "caveman", "ponytail", "trim"] as const;

const CAVEMAN_COPY: Record<UILocale, Record<TokenSaverLevel, { label: string; desc: string }>> = {
  en: {
    lite: { label: "Lite", desc: "Drop filler, keep clear grammar" },
    full: { label: "Full", desc: "Drop articles, fragments OK" },
    ultra: { label: "Ultra", desc: "Telegraphic, max output compression" },
  },
  id: {
    lite: { label: "Lite", desc: "Hapus basa-basi, jaga tata bahasa jelas" },
    full: { label: "Full", desc: "Ringkas padat, fragmen kalimat diperbolehkan" },
    ultra: { label: "Ultra", desc: "Sangat telegrafis, kompresi output maksimal" },
  },
};

const PONYTAIL_COPY: Record<UILocale, Record<TokenSaverLevel, { label: string; desc: string }>> = {
  en: {
    lite: { label: "Lite", desc: "Build what was asked, note simpler option" },
    full: { label: "Full", desc: "Reuse stdlib and existing utilities first" },
    ultra: { label: "Ultra", desc: "YAGNI extremist, deletion over addition" },
  },
  id: {
    lite: { label: "Lite", desc: "Kerjakan yang diminta, sebut opsi lebih sederhana" },
    full: { label: "Full", desc: "Utamakan stdlib dan utilitas yang sudah ada" },
    ultra: { label: "Ultra", desc: "YAGNI ekstrem, utamakan penghapusan kode" },
  },
};

const HISTORY_COPY: Record<
  UILocale,
  Record<(typeof HISTORY_LIMIT_PRESETS)[number]["id"], { label: string; desc: string }>
> = {
  en: {
    tight: { label: "Tight", desc: "Keep system prompt + last 12 messages" },
    balanced: { label: "Balanced", desc: "Keep system prompt + last 20 messages" },
    deep: { label: "Deep", desc: "Keep system prompt + last 40 messages" },
  },
  id: {
    tight: { label: "Ketat", desc: "Simpan prompt sistem + 12 pesan terakhir" },
    balanced: { label: "Seimbang", desc: "Simpan prompt sistem + 20 pesan terakhir" },
    deep: { label: "Dalam", desc: "Simpan prompt sistem + 40 pesan terakhir" },
  },
};

const COPY = {
  en: {
    title: "Token Saver",
    saved: "Saved",
    saving: "Saving...",
    errorTitle: "Token Saver Error",
    retry: "Retry",
    dismiss: "Dismiss",
    tokensSavedSuffix: "tok saved",
    statusActive: "Active",
    statusOff: "Off",
    rtkTitle: "Compress tool output",
    rtkTag: "(RTK)",
    rtkDesc: "Whitespace cleanup, deduplication & oversized tool output truncation",
    headroomTitle: "Compress context",
    headroomTag: "(Headroom)",
    headroomDesc: "Sticky LRU prefix cache affinity & Anthropic ephemeral cache_control",
    cavemanTitle: "Compress LLM output",
    cavemanTag: "(Caveman)",
    cavemanDesc: "Terse-style system prompt → ~65% fewer output tokens",
    cavemanGroupLabel: "Caveman compression level",
    ponytailTitle: "Lazy senior dev",
    ponytailTag: "(Ponytail)",
    ponytailDesc: "Bias the model toward minimal code: YAGNI, reuse stdlib, deletion over addition",
    ponytailGroupLabel: "Ponytail minimalism level",
    trimTitle: "Trim chat history",
    trimDesc:
      "Keep the system prompt + last N messages; 30-50% fewer input tokens in long sessions",
    trimGroupLabel: "Chat history limit preset",
    trimCustomDesc: (count: number) => `Keep system prompt + last ${count} messages`,
    statsCompressed: "compressed",
    statsWarmCache: "warm prefixes",
    resetStats: "Reset stats",
    clearCache: "Clear cache",
    loadingLabel: "Loading Token Saver configuration",
  },
  id: {
    title: "Penghemat Token",
    saved: "Tersimpan",
    saving: "Menyimpan...",
    errorTitle: "Galat Penghemat Token",
    retry: "Coba lagi",
    dismiss: "Tutup",
    tokensSavedSuffix: "tok hemat",
    statusActive: "Aktif",
    statusOff: "Nonaktif",
    rtkTitle: "Kompresi output tool",
    rtkTag: "(RTK)",
    rtkDesc: "Pembersihan spasi, deduplikasi & pemangkasan output tool yang terlalu panjang",
    headroomTitle: "Kompresi konteks",
    headroomTag: "(Headroom)",
    headroomDesc: "Afinitas cache prefix LRU sticky & injeksi cache_control ephemeral Anthropic",
    cavemanTitle: "Kompresi output LLM",
    cavemanTag: "(Caveman)",
    cavemanDesc: "Prompt sistem bergaya ringkas → ~65% lebih sedikit token output",
    cavemanGroupLabel: "Level kompresi Caveman",
    ponytailTitle: "Lazy senior dev",
    ponytailTag: "(Ponytail)",
    ponytailDesc: "Arahkan model ke kode minimal: YAGNI, utamakan stdlib, hapus sebelum menambah",
    ponytailGroupLabel: "Level minimalisme Ponytail",
    trimTitle: "Pangkas riwayat chat",
    trimDesc: "Simpan prompt sistem + N pesan terakhir; hemat 30-50% token input pada sesi panjang",
    trimGroupLabel: "Preset batas riwayat chat",
    trimCustomDesc: (count: number) => `Simpan prompt sistem + ${count} pesan terakhir`,
    statsCompressed: "terkompresi",
    statsWarmCache: "prefix aktif",
    resetStats: "Reset statistik",
    clearCache: "Bersihkan cache",
    loadingLabel: "Memuat konfigurasi Penghemat Token",
  },
} as const;

const LEVEL_BTN_BASE =
  "rounded border px-3 py-1 text-xs font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent-primary)] focus-visible:ring-offset-1 focus-visible:ring-offset-[var(--bg-surface)] disabled:opacity-50";

/** Isolated polling leaf so 5s stats ticks never re-render the toggle controls. */
const SavedTokensBadge = memo(function SavedTokensBadge({ locale }: { locale: UILocale }) {
  const statsQuery = useTokenSaverStats();
  const saved = statsQuery.data?.compression.estimated_tokens_saved ?? 0;
  if (saved <= 0) return null;
  const formatted = saved.toLocaleString(locale === "id" ? "id-ID" : "en-US");
  return (
    <Badge variant="secondary" className="tabular-nums">
      {formatted} {COPY[locale].tokensSavedSuffix}
    </Badge>
  );
});

const HeadroomStatusBadge = memo(function HeadroomStatusBadge({
  enabled,
  locale,
}: {
  enabled: boolean;
  locale: UILocale;
}) {
  const statsQuery = useTokenSaverStats();
  const copy = COPY[locale];
  const headroom = statsQuery.data?.headroom;
  const hasTraffic = headroom ? headroom.cache_hits + headroom.cache_misses > 0 : false;
  const hitRatePct = headroom ? `${(headroom.hit_rate * 100).toFixed(0)}%` : "0%";

  return (
    <Badge variant={enabled ? "success" : "secondary"} className="tabular-nums">
      {enabled
        ? hasTraffic
          ? `${copy.statusActive} · ${hitRatePct}`
          : copy.statusActive
        : copy.statusOff}
    </Badge>
  );
});

const TokenSaverTelemetryBar = memo(function TokenSaverTelemetryBar({
  locale,
}: {
  locale: UILocale;
}) {
  const copy = COPY[locale];
  const statsQuery = useTokenSaverStats();
  const resetStats = useResetTokenSaverStatsMutation();
  const clearCache = useClearHeadroomCacheMutation();

  const stats = statsQuery.data;
  const actionError = resetStats.error?.message ?? clearCache.error?.message ?? null;
  const hasCounters =
    (stats?.compression.requests_processed ?? 0) > 0 ||
    (stats?.headroom.cache_hits ?? 0) + (stats?.headroom.cache_misses ?? 0) > 0;
  const activeEntries = stats?.headroom.active_entries ?? 0;

  if (!hasCounters && activeEntries === 0 && !actionError) {
    return null;
  }

  const numLocale = locale === "id" ? "id-ID" : "en-US";
  const compressed = (stats?.compression.requests_compressed ?? 0).toLocaleString(numLocale);
  const processed = (stats?.compression.requests_processed ?? 0).toLocaleString(numLocale);

  return (
    <div className="mt-4 space-y-2 border-t border-[var(--border-subtle)] pt-3">
      {actionError ? (
        <p role="alert" className="text-xs text-[var(--status-err)]">
          {actionError}
        </p>
      ) : null}
      <div className="flex flex-wrap items-center justify-between gap-2">
        <div className="flex flex-wrap items-center gap-3 text-xs text-[var(--text-muted)] tabular-nums">
          {hasCounters ? (
            <span>
              {compressed}/{processed} {copy.statsCompressed}
            </span>
          ) : null}
          {activeEntries > 0 ? (
            <span>
              {activeEntries.toLocaleString(numLocale)} {copy.statsWarmCache}
            </span>
          ) : null}
        </div>
        <div className="flex flex-wrap items-center gap-1.5">
          {hasCounters ? (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={resetStats.isPending}
              onClick={() => resetStats.mutate()}
              className="h-7 gap-1 px-2 text-xs text-[var(--text-secondary)]"
            >
              <RotateCcw className="h-3.5 w-3.5" />
              {copy.resetStats}
            </Button>
          ) : null}
          {activeEntries > 0 ? (
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={clearCache.isPending}
              onClick={() => clearCache.mutate()}
              className="h-7 gap-1 px-2 text-xs text-[var(--text-secondary)]"
            >
              <Trash2 className="h-3.5 w-3.5" />
              {copy.clearCache}
            </Button>
          ) : null}
        </div>
      </div>
    </div>
  );
});

export function TokenSaverConsole() {
  const locale = useUIStore((state) => state.locale);
  const copy = COPY[locale];
  const cavemanCopy = CAVEMAN_COPY[locale];
  const ponytailCopy = PONYTAIL_COPY[locale];
  const historyCopy = HISTORY_COPY[locale];

  const configQuery = useTokenSaverConfig();
  const updateConfig = useUpdateTokenSaverConfigMutation();

  const [optimisticPatch, setOptimisticPatch] = useState<TokenSaverConfigPatch | null>(null);
  const [savedNotice, setSavedNotice] = useState(false);

  useEffect(() => {
    if (!savedNotice) return;
    const timer = window.setTimeout(() => setSavedNotice(false), 2500);
    return () => window.clearTimeout(timer);
  }, [savedNotice]);

  const baseConfig: TokenSaverConfigFormValues = configQuery.data ?? DEFAULT_CONFIG;
  const cfg: TokenSaverConfigFormValues = optimisticPatch
    ? { ...baseConfig, ...optimisticPatch }
    : baseConfig;

  const applyPatch = (patch: TokenSaverConfigPatch) => {
    setSavedNotice(false);
    setOptimisticPatch((prev) => ({ ...prev, ...patch }));
    updateConfig.mutate(patch, {
      onSuccess: () => setSavedNotice(true),
      onSettled: () => setOptimisticPatch(null),
    });
  };

  const handleRtkToggle = (next: boolean) => {
    applyPatch({
      enabled: next,
      strategy: next ? "aggressive" : "off",
    });
  };

  const handleHeadroomToggle = (next: boolean) => {
    applyPatch({
      headroom_enabled: next,
      inject_cache_control: next,
    });
  };

  const isInitialLoading = configQuery.isLoading && !configQuery.data;
  const isBusy = updateConfig.isPending || isInitialLoading;
  const activeError = configQuery.error?.message ?? updateConfig.error?.message ?? null;
  const activeHistoryPreset = matchHistoryPreset(cfg.max_history_messages);

  return (
    <div className="mx-auto max-w-4xl space-y-4">
      {activeError ? (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>{copy.errorTitle}</AlertTitle>
          <AlertDescription className="flex flex-wrap items-center justify-between gap-2">
            <span>{activeError}</span>
            <div className="flex items-center gap-1.5">
              {configQuery.error ? (
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => void configQuery.refetch()}
                  className="h-7 px-2.5 text-xs"
                >
                  {copy.retry}
                </Button>
              ) : null}
              {updateConfig.error ? (
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => updateConfig.reset()}
                  className="h-7 px-2.5 text-xs"
                >
                  {copy.dismiss}
                </Button>
              ) : null}
            </div>
          </AlertDescription>
        </Alert>
      ) : null}

      <Card className="p-6">
        <div className="mb-2 flex flex-wrap items-center justify-between gap-2">
          <h1 className="flex items-center gap-2 text-lg font-semibold text-[var(--text-primary)]">
            <Zap className="h-5 w-5 text-[var(--accent-primary)]" />
            {copy.title}
          </h1>
          <div className="flex items-center gap-2">
            <SavedTokensBadge locale={locale} />
            {updateConfig.isPending ? (
              <span className="font-mono text-xs text-[var(--text-muted)]">{copy.saving}</span>
            ) : savedNotice ? (
              <span className="inline-flex items-center gap-1 text-xs font-medium text-[var(--status-ok)]">
                <CheckCircle2 className="h-3.5 w-3.5" />
                {copy.saved}
              </span>
            ) : null}
          </div>
        </div>

        {isInitialLoading ? (
          <div
            role="status"
            aria-busy="true"
            aria-label={copy.loadingLabel}
            className="divide-y divide-[var(--border-subtle)]"
          >
            <span className="sr-only">{copy.loadingLabel}</span>
            {SKELETON_ROWS.map((rowKey) => (
              <div key={rowKey} className="flex items-center justify-between gap-4 py-4">
                <div className="min-w-0 flex-1 space-y-2">
                  <Skeleton className="h-4 w-44" />
                  <Skeleton className="h-3.5 w-80 max-w-full" />
                </div>
                <Skeleton className="h-6 w-11 shrink-0 rounded-full" />
              </div>
            ))}
          </div>
        ) : (
          <>
            <div className="divide-y divide-[var(--border-subtle)]">
              {/* 1. Compress tool output (RTK) */}
              <div className="flex items-center justify-between gap-4 py-4">
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-[var(--text-primary)]">
                    {copy.rtkTitle}{" "}
                    <span className="text-xs font-normal text-[var(--accent-text)]">
                      {copy.rtkTag}
                    </span>
                  </p>
                  <p className="mt-0.5 text-sm text-[var(--text-secondary)]">{copy.rtkDesc}</p>
                </div>
                <ToggleSwitch
                  id="ts-rtk"
                  checked={cfg.enabled}
                  disabled={isBusy}
                  label={`${copy.rtkTitle} ${copy.rtkTag}`}
                  onChange={handleRtkToggle}
                />
              </div>

              {/* 2. Compress context (Headroom) */}
              <div className="flex items-center justify-between gap-4 py-4">
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="text-sm font-medium text-[var(--text-primary)]">
                      {copy.headroomTitle}{" "}
                      <span className="text-xs font-normal text-[var(--accent-text)]">
                        {copy.headroomTag}
                      </span>
                    </p>
                    <HeadroomStatusBadge enabled={cfg.headroom_enabled} locale={locale} />
                  </div>
                  <p className="mt-0.5 text-sm text-[var(--text-secondary)]">{copy.headroomDesc}</p>
                </div>
                <ToggleSwitch
                  id="ts-headroom"
                  checked={cfg.headroom_enabled}
                  disabled={isBusy}
                  label={`${copy.headroomTitle} ${copy.headroomTag}`}
                  onChange={handleHeadroomToggle}
                />
              </div>

              {/* 3. Compress LLM output (Caveman) */}
              <div className="flex items-start justify-between gap-4 py-4">
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-[var(--text-primary)]">
                    {copy.cavemanTitle}{" "}
                    <span className="text-xs font-normal text-[var(--accent-text)]">
                      {copy.cavemanTag}
                    </span>
                  </p>
                  <p className="mt-0.5 text-sm text-[var(--text-secondary)]">{copy.cavemanDesc}</p>
                  {cfg.caveman_enabled ? (
                    <div className="mt-2.5 space-y-1.5">
                      <fieldset
                        aria-label={copy.cavemanGroupLabel}
                        className="flex flex-wrap items-center gap-1.5 border-0 p-0"
                      >
                        {TOKEN_SAVER_LEVELS.map((lvl) => {
                          const isSelected = cfg.caveman_level === lvl;
                          return (
                            <button
                              key={lvl}
                              type="button"
                              aria-pressed={isSelected}
                              disabled={isBusy}
                              onClick={() => applyPatch({ caveman_level: lvl })}
                              className={`${LEVEL_BTN_BASE} ${
                                isSelected
                                  ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                                  : "border-[var(--border-strong)] bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                              }`}
                            >
                              {cavemanCopy[lvl].label}
                            </button>
                          );
                        })}
                      </fieldset>
                      <p className="text-xs text-[var(--accent-text)]">
                        {cavemanCopy[cfg.caveman_level].desc}
                      </p>
                    </div>
                  ) : null}
                </div>
                <div className="pt-0.5">
                  <ToggleSwitch
                    id="ts-caveman"
                    checked={cfg.caveman_enabled}
                    disabled={isBusy}
                    label={`${copy.cavemanTitle} ${copy.cavemanTag}`}
                    onChange={(next) => applyPatch({ caveman_enabled: next })}
                  />
                </div>
              </div>

              {/* 4. Lazy senior dev (Ponytail) */}
              <div className="flex items-start justify-between gap-4 py-4">
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-[var(--text-primary)]">
                    {copy.ponytailTitle}{" "}
                    <span className="text-xs font-normal text-[var(--accent-text)]">
                      {copy.ponytailTag}
                    </span>
                  </p>
                  <p className="mt-0.5 text-sm text-[var(--text-secondary)]">{copy.ponytailDesc}</p>
                  {cfg.ponytail_enabled ? (
                    <div className="mt-2.5 space-y-1.5">
                      <fieldset
                        aria-label={copy.ponytailGroupLabel}
                        className="flex flex-wrap items-center gap-1.5 border-0 p-0"
                      >
                        {TOKEN_SAVER_LEVELS.map((lvl) => {
                          const isSelected = cfg.ponytail_level === lvl;
                          return (
                            <button
                              key={lvl}
                              type="button"
                              aria-pressed={isSelected}
                              disabled={isBusy}
                              onClick={() => applyPatch({ ponytail_level: lvl })}
                              className={`${LEVEL_BTN_BASE} ${
                                isSelected
                                  ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                                  : "border-[var(--border-strong)] bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                              }`}
                            >
                              {ponytailCopy[lvl].label}
                            </button>
                          );
                        })}
                      </fieldset>
                      <p className="text-xs text-[var(--accent-text)]">
                        {ponytailCopy[cfg.ponytail_level].desc}
                      </p>
                    </div>
                  ) : null}
                </div>
                <div className="pt-0.5">
                  <ToggleSwitch
                    id="ts-ponytail"
                    checked={cfg.ponytail_enabled}
                    disabled={isBusy}
                    label={`${copy.ponytailTitle} ${copy.ponytailTag}`}
                    onChange={(next) => applyPatch({ ponytail_enabled: next })}
                  />
                </div>
              </div>

              {/* 5. Trim chat history */}
              <div className="flex items-start justify-between gap-4 pt-4">
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-[var(--text-primary)]">{copy.trimTitle}</p>
                  <p className="mt-0.5 text-sm text-[var(--text-secondary)]">{copy.trimDesc}</p>
                  {cfg.trim_history_enabled ? (
                    <div className="mt-2.5 space-y-1.5">
                      <fieldset
                        aria-label={copy.trimGroupLabel}
                        className="flex flex-wrap items-center gap-1.5 border-0 p-0"
                      >
                        {HISTORY_LIMIT_PRESETS.map((preset) => {
                          const isSelected = activeHistoryPreset === preset.id;
                          return (
                            <button
                              key={preset.id}
                              type="button"
                              aria-pressed={isSelected}
                              disabled={isBusy}
                              onClick={() =>
                                applyPatch({ max_history_messages: preset.maxMessages })
                              }
                              className={`${LEVEL_BTN_BASE} ${
                                isSelected
                                  ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                                  : "border-[var(--border-strong)] bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                              }`}
                            >
                              {historyCopy[preset.id].label} ({preset.maxMessages})
                            </button>
                          );
                        })}
                      </fieldset>
                      <p className="text-xs text-[var(--accent-text)]">
                        {activeHistoryPreset !== "custom"
                          ? historyCopy[activeHistoryPreset].desc
                          : copy.trimCustomDesc(cfg.max_history_messages)}
                      </p>
                    </div>
                  ) : null}
                </div>
                <div className="pt-0.5">
                  <ToggleSwitch
                    id="ts-trim-history"
                    checked={cfg.trim_history_enabled}
                    disabled={isBusy}
                    label={copy.trimTitle}
                    onChange={(next) => applyPatch({ trim_history_enabled: next })}
                  />
                </div>
              </div>
            </div>

            <TokenSaverTelemetryBar locale={locale} />
          </>
        )}
      </Card>
    </div>
  );
}
