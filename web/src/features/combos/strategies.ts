import type { ComboStrategy } from "./schemas";

export type StrategyCategory =
  | "failover"
  | "load-balance"
  | "quota-cost"
  | "context-cache"
  | "intelligent";

export interface ComboStrategyMeta {
  value: ComboStrategy;
  category: StrategyCategory;
  labelEn: string;
  labelId: string;
  descEn: string;
  descId: string;
  usesWeights?: boolean;
  usesRouterStrategy?: boolean;
  usesJudgeModel?: boolean;
  usesHandoff?: boolean;
  usesStepPrompt?: boolean;
}

export const STRATEGY_CATEGORIES: Array<{
  id: StrategyCategory;
  labelEn: string;
  labelId: string;
}> = [
  { id: "failover", labelEn: "Failover & Order", labelId: "Failover & Urutan" },
  { id: "load-balance", labelEn: "Load Balancing", labelId: "Penyeimbang Beban" },
  { id: "quota-cost", labelEn: "Quota & Cost", labelId: "Kuota & Biaya" },
  { id: "context-cache", labelEn: "Context & Cache", labelId: "Konteks & Cache" },
  { id: "intelligent", labelEn: "Intelligent & Multi-Model", labelId: "Cerdas & Multi-Model" },
];

export const COMBO_STRATEGY_OPTIONS: readonly ComboStrategyMeta[] = [
  // Failover & Order
  {
    value: "priority",
    category: "failover",
    labelEn: "Priority",
    labelId: "Prioritas",
    descEn: "Tries targets strictly in order and fails over on error or cooldown.",
    descId: "Mencoba target berurutan dan berpindah otomatis saat error atau cooldown.",
  },
  {
    value: "fill-first",
    category: "failover",
    labelEn: "Fill First",
    labelId: "Isi Pertama",
    descEn: "Exhausts the primary target's capacity before cascading to backups.",
    descId: "Memaksimalkan kapasitas target pertama sebelum beralih ke cadangan.",
  },
  {
    value: "lkgp",
    category: "failover",
    labelEn: "Last Known Good (LKGP)",
    labelId: "Provider Sehat Terakhir (LKGP)",
    descEn: "Pins requests to the last known healthy target and fails over automatically.",
    descId: "Mengunci permintaan ke target sehat terakhir dan failover otomatis.",
  },

  // Load Balancing
  {
    value: "weighted",
    category: "load-balance",
    labelEn: "Weighted",
    labelId: "Berbobot",
    descEn: "Distributes traffic proportionally by step weight percentage with failover.",
    descId: "Membagi trafik secara proporsional berdasarkan bobot persen tiap langkah.",
    usesWeights: true,
  },
  {
    value: "round-robin",
    category: "load-balance",
    labelEn: "Round Robin",
    labelId: "Round Robin",
    descEn: "Rotates cyclically across all targets to distribute load evenly.",
    descId: "Memutar permintaan secara bergiliran ke seluruh target agar beban merata.",
  },
  {
    value: "p2c",
    category: "load-balance",
    labelEn: "Power of Two Choices (P2C)",
    labelId: "Power of Two Choices (P2C)",
    descEn: "Samples two candidates and picks the less-loaded healthy target.",
    descId: "Memilih dua kandidat acak lalu mengambil target sehat dengan beban terendah.",
  },
  {
    value: "least-used",
    category: "load-balance",
    labelEn: "Least Used",
    labelId: "Paling Jarang Dipakai",
    descEn: "Routes to the target with the fewest recent requests.",
    descId: "Mengarahkan ke target dengan jumlah penggunaan terbaru paling sedikit.",
  },
  {
    value: "random",
    category: "load-balance",
    labelEn: "Random",
    labelId: "Acak",
    descEn: "Randomly selects a target with optional session stickiness.",
    descId: "Memilih target secara acak dengan dukungan session stickiness.",
  },
  {
    value: "strict-random",
    category: "load-balance",
    labelEn: "Strict Random",
    labelId: "Acak Murni",
    descEn: "Pure uniform random selection per request with no session stickiness.",
    descId: "Pemilihan acak seragam per permintaan tanpa session stickiness.",
  },

  // Quota & Cost
  {
    value: "cost-optimized",
    category: "quota-cost",
    labelEn: "Cost Optimized",
    labelId: "Optimasi Biaya",
    descEn: "Prefers free-tier and lowest-cost models first before paid fallbacks.",
    descId: "Memprioritaskan model gratis dan termurah sebelum fallback berbayar.",
  },
  {
    value: "quota-weighted",
    category: "quota-cost",
    labelEn: "Quota Weighted",
    labelId: "Berbobot Kuota",
    descEn: "Distributes traffic proportionally based on remaining quota and weights.",
    descId: "Membagi trafik berdasarkan sisa kuota dan bobot langkah.",
    usesWeights: true,
  },
  {
    value: "headroom",
    category: "quota-cost",
    labelEn: "Headroom",
    labelId: "Headroom",
    descEn: "Selects the target with the highest available rate-limit headroom.",
    descId: "Memilih target dengan ruang rate-limit tersisa paling besar.",
  },
  {
    value: "reset-aware",
    category: "quota-cost",
    labelEn: "Reset Aware",
    labelId: "Sadar Reset Kuota",
    descEn: "Balances across accounts using session and weekly quota reset timing.",
    descId: "Menyeimbangkan akun berdasarkan waktu reset kuota sesi dan mingguan.",
  },
  {
    value: "reset-window",
    category: "quota-cost",
    labelEn: "Reset Window",
    labelId: "Jendela Reset",
    descEn: "Orders accounts by soonest upcoming quota reset window.",
    descId: "Mengurutkan akun berdasarkan jendela reset kuota terdekat.",
  },
  {
    value: "expiry-first",
    category: "quota-cost",
    labelEn: "Expiry First",
    labelId: "Kedaluwarsa Terdekat",
    descEn: "Consumes soonest-expiring credits or subscriptions first.",
    descId: "Menghabiskan kredit atau paket yang paling cepat kedaluwarsa lebih dulu.",
  },

  // Context & Cache
  {
    value: "context-optimized",
    category: "context-cache",
    labelEn: "Context Optimized",
    labelId: "Optimasi Konteks",
    descEn: "Prioritizes targets with the largest context window for long prompts.",
    descId: "Memprioritaskan target dengan jendela konteks terbesar untuk prompt panjang.",
  },
  {
    value: "cache-optimized",
    category: "context-cache",
    labelEn: "Cache Optimized",
    labelId: "Optimasi Cache",
    descEn: "Deterministic sticky routing to maximize upstream prompt cache hits.",
    descId: "Routing deterministik agar prompt cache upstream selalu hit maksimal.",
  },
  {
    value: "context-relay",
    category: "context-cache",
    labelEn: "Context Relay",
    labelId: "Estafet Konteks",
    descEn: "Relays long conversations to a handoff target once context threshold hits.",
    descId: "Mengalihkan percakapan panjang ke target handoff saat ambang konteks tercapai.",
    usesHandoff: true,
  },

  // Intelligent & Multi-Model
  {
    value: "auto",
    category: "intelligent",
    labelEn: "Auto Router",
    labelId: "Auto Router",
    descEn: "Dynamically routes requests using a configurable auto-routing policy.",
    descId: "Merutekan permintaan secara dinamis sesuai kebijakan auto-router.",
    usesRouterStrategy: true,
  },
  {
    value: "fusion",
    category: "intelligent",
    labelEn: "Fusion (Panel + Judge)",
    labelId: "Fusion (Panel + Juri)",
    descEn: "Fans out to panel models and synthesizes the final response via a judge model.",
    descId: "Mengirim paralel ke panel model lalu menyintesis jawaban lewat model juri.",
    usesJudgeModel: true,
  },
  {
    value: "pipeline",
    category: "intelligent",
    labelEn: "Pipeline",
    labelId: "Pipeline",
    descEn: "Chains steps sequentially where each step's output feeds the next.",
    descId: "Menjalankan langkah berurutan di mana output satu langkah menjadi input berikutnya.",
    usesStepPrompt: true,
  },
  {
    value: "jev",
    category: "intelligent",
    labelEn: "JEV (Judge-Evaluate-Verify)",
    labelId: "JEV (Judge-Evaluate-Verify)",
    descEn: "Multi-stage generation, evaluation, and verification routing chain.",
    descId: "Rantai generasi, evaluasi, dan verifikasi multi-tahap.",
  },
] as const;

export const AUTO_ROUTER_STRATEGIES = [
  "rules",
  "score",
  "cost",
  "eco",
  "latency",
  "fast",
  "sla-aware",
  "sla",
  "lkgp",
  "nadir",
] as const;

export const COMPRESSION_MODES = [
  "",
  "off",
  "lite",
  "standard",
  "aggressive",
  "ultra",
  "rtk",
  "stacked",
] as const;

export function getStrategyMeta(strategy: string): ComboStrategyMeta {
  return (
    COMBO_STRATEGY_OPTIONS.find((s) => s.value === strategy) ??
    (COMBO_STRATEGY_OPTIONS[0] as ComboStrategyMeta)
  );
}
