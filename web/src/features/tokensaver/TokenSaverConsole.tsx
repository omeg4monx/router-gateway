import { CheckCircle2, Zap } from "lucide-react";
import { useEffect, useState } from "react";
import { type UILocale, useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Card } from "@/shared/ui/card";
import {
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
      className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border transition-colors disabled:cursor-not-allowed disabled:opacity-50 ${
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
  strategy: "aggressive",
  max_history_messages: 20,
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
    ponytailTitle: "Lazy senior dev",
    ponytailTag: "(Ponytail)",
    ponytailDesc: "Bias the model toward minimal code: YAGNI, reuse stdlib, deletion over addition",
    trimTitle: "Trim chat history",
    trimDesc:
      "Keep the system prompt + last N messages; 30-50% fewer input tokens in long sessions",
  },
  id: {
    title: "Penghemat Token",
    saved: "Tersimpan",
    saving: "Menyimpan...",
    errorTitle: "Galat Penghemat Token",
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
    ponytailTitle: "Lazy senior dev",
    ponytailTag: "(Ponytail)",
    ponytailDesc: "Arahkan model ke kode minimal: YAGNI, utamakan stdlib, hapus sebelum menambah",
    trimTitle: "Pangkas riwayat chat",
    trimDesc: "Simpan prompt sistem + N pesan terakhir; hemat 30-50% token input pada sesi panjang",
  },
} as const;

export function TokenSaverConsole() {
  const locale = useUIStore((state) => state.locale);
  const copy = COPY[locale];
  const cavemanCopy = CAVEMAN_COPY[locale];
  const ponytailCopy = PONYTAIL_COPY[locale];
  const historyCopy = HISTORY_COPY[locale];

  const configQuery = useTokenSaverConfig();
  const statsQuery = useTokenSaverStats();
  const updateConfig = useUpdateTokenSaverConfigMutation();

  const [cfg, setCfg] = useState<TokenSaverConfigFormValues>(DEFAULT_CONFIG);
  const [savedNotice, setSavedNotice] = useState(false);

  useEffect(() => {
    if (configQuery.data) {
      setCfg({
        enabled: configQuery.data.enabled,
        strategy: configQuery.data.strategy,
        max_history_messages: configQuery.data.max_history_messages,
        max_tool_output_chars: configQuery.data.max_tool_output_chars,
        headroom_enabled: configQuery.data.headroom_enabled,
        prefix_depth: configQuery.data.prefix_depth,
        max_cache_entries: configQuery.data.max_cache_entries,
        cache_ttl_seconds: configQuery.data.cache_ttl_seconds,
        inject_cache_control: configQuery.data.inject_cache_control,
        caveman_enabled: configQuery.data.caveman_enabled,
        caveman_level: configQuery.data.caveman_level,
        ponytail_enabled: configQuery.data.ponytail_enabled,
        ponytail_level: configQuery.data.ponytail_level,
        trim_history_enabled: configQuery.data.trim_history_enabled,
      });
    }
  }, [configQuery.data]);

  const applyPatch = (patch: TokenSaverConfigPatch) => {
    setSavedNotice(false);
    setCfg((prev) => ({ ...prev, ...patch }));
    updateConfig.mutate(patch, {
      onSuccess: () => setSavedNotice(true),
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

  const isBusy = updateConfig.isPending || configQuery.isLoading;
  const activeError =
    configQuery.error?.message ?? statsQuery.error?.message ?? updateConfig.error?.message ?? null;

  const stats = statsQuery.data;
  const hitRatePct = stats ? `${(stats.headroom.hit_rate * 100).toFixed(0)}%` : "0%";
  const activeHistoryPreset = matchHistoryPreset(cfg.max_history_messages);

  return (
    <div className="mx-auto max-w-4xl space-y-4">
      {activeError ? (
        <Alert variant="destructive">
          <Zap className="h-4 w-4" />
          <AlertTitle>{copy.errorTitle}</AlertTitle>
          <AlertDescription>{activeError}</AlertDescription>
        </Alert>
      ) : null}

      <Card className="p-6">
        <div className="mb-2 flex flex-wrap items-center justify-between gap-2">
          <h1 className="flex items-center gap-2 text-lg font-semibold text-[var(--text-primary)]">
            <Zap className="h-5 w-5 text-[var(--accent-primary)]" />
            {copy.title}
          </h1>
          <div className="flex items-center gap-2">
            {stats && stats.compression.estimated_tokens_saved > 0 ? (
              <Badge variant="secondary" className="tabular-nums">
                {stats.compression.estimated_tokens_saved} tok saved
              </Badge>
            ) : null}
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

        <div className="divide-y divide-[var(--border-subtle)]">
          {/* 1. Compress tool output (RTK) */}
          <div className="flex items-center justify-between gap-4 py-4">
            <div className="min-w-0 flex-1">
              <p className="text-sm font-medium text-[var(--text-primary)]">
                {copy.rtkTitle}{" "}
                <span className="text-xs font-normal text-[var(--accent-text)]">{copy.rtkTag}</span>
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
                <Badge variant={cfg.headroom_enabled ? "success" : "secondary"}>
                  {cfg.headroom_enabled
                    ? stats && stats.headroom.cache_hits + stats.headroom.cache_misses > 0
                      ? `${copy.statusActive} · ${hitRatePct}`
                      : copy.statusActive
                    : copy.statusOff}
                </Badge>
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
                  <div className="flex flex-wrap items-center gap-1.5">
                    {TOKEN_SAVER_LEVELS.map((lvl) => (
                      <button
                        key={lvl}
                        type="button"
                        disabled={isBusy}
                        onClick={() => applyPatch({ caveman_level: lvl })}
                        className={`rounded border px-3 py-1 text-xs font-medium transition-colors disabled:opacity-50 ${
                          cfg.caveman_level === lvl
                            ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                            : "border-[var(--border-strong)] bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                        }`}
                      >
                        {cavemanCopy[lvl].label}
                      </button>
                    ))}
                  </div>
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
                  <div className="flex flex-wrap items-center gap-1.5">
                    {TOKEN_SAVER_LEVELS.map((lvl) => (
                      <button
                        key={lvl}
                        type="button"
                        disabled={isBusy}
                        onClick={() => applyPatch({ ponytail_level: lvl })}
                        className={`rounded border px-3 py-1 text-xs font-medium transition-colors disabled:opacity-50 ${
                          cfg.ponytail_level === lvl
                            ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                            : "border-[var(--border-strong)] bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                        }`}
                      >
                        {ponytailCopy[lvl].label}
                      </button>
                    ))}
                  </div>
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
                  <div className="flex flex-wrap items-center gap-1.5">
                    {HISTORY_LIMIT_PRESETS.map((preset) => (
                      <button
                        key={preset.id}
                        type="button"
                        disabled={isBusy}
                        onClick={() => applyPatch({ max_history_messages: preset.maxMessages })}
                        className={`rounded border px-3 py-1 text-xs font-medium transition-colors disabled:opacity-50 ${
                          activeHistoryPreset === preset.id
                            ? "border-[var(--accent-primary)] bg-[var(--accent-primary)] text-[var(--accent-primary-ink)]"
                            : "border-[var(--border-strong)] bg-transparent text-[var(--text-secondary)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]"
                        }`}
                      >
                        {historyCopy[preset.id].label} ({preset.maxMessages})
                      </button>
                    ))}
                  </div>
                  <p className="text-xs text-[var(--accent-text)]">
                    {activeHistoryPreset !== "custom"
                      ? historyCopy[activeHistoryPreset].desc
                      : `Keep system prompt + last ${cfg.max_history_messages} messages`}
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
      </Card>
    </div>
  );
}
