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
  flowPreview: string;
  usesWeights?: boolean;
  usesRouterStrategy?: boolean;
  usesJudgeModel?: boolean;
  usesHandoff?: boolean;
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
    flowPreview: "#1 Primary → (error/429) → #2 Backup → #3 Fallback",
  },
  {
    value: "fill-first",
    category: "failover",
    labelEn: "Fill First",
    labelId: "Isi Pertama",
    descEn: "Exhausts the primary target's capacity before cascading to backups.",
    descId: "Memaksimalkan kapasitas target pertama sebelum beralih ke cadangan.",
    flowPreview: "#1 (100% until full) → #2 (until full) → #3",
  },
  {
    value: "lkgp",
    category: "failover",
    labelEn: "Last Known Good (LKGP)",
    labelId: "Provider Sehat Terakhir (LKGP)",
    descEn: "Pins requests to the last known healthy target and fails over automatically.",
    descId: "Mengunci permintaan ke target sehat terakhir dan failover otomatis.",
    flowPreview: "Active Healthy Target ↺ (error) → Promote Next Healthy",
  },

  // Load Balancing
  {
    value: "weighted",
    category: "load-balance",
    labelEn: "Weighted",
    labelId: "Berbobot",
    descEn: "Distributes traffic proportionally by step weight percentage with failover.",
    descId: "Membagi trafik secara proporsional berdasarkan bobot persen tiap langkah.",
    flowPreview: "70% Model A | 20% Model B | 10% Model C (+ failover)",
    usesWeights: true,
  },
  {
    value: "round-robin",
    category: "load-balance",
    labelEn: "Round Robin",
    labelId: "Round Robin",
    descEn: "Rotates cyclically across all targets to distribute load evenly.",
    descId: "Memutar permintaan secara bergiliran ke seluruh target agar beban merata.",
    flowPreview: "Req 1 → #1 · Req 2 → #2 · Req 3 → #3 · Req 4 → #1",
  },
  {
    value: "p2c",
    category: "load-balance",
    labelEn: "Power of Two Choices (P2C)",
    labelId: "Power of Two Choices (P2C)",
    descEn: "Samples two candidates and picks the less-loaded healthy target.",
    descId: "Memilih dua kandidat acak lalu mengambil target sehat dengan beban terendah.",
    flowPreview: "Sample 2 Random → Compare Load → Pick Lighter Target",
  },
  {
    value: "least-used",
    category: "load-balance",
    labelEn: "Least Used",
    labelId: "Paling Jarang Dipakai",
    descEn: "Routes to the target with the fewest recent requests.",
    descId: "Mengarahkan ke target dengan jumlah penggunaan terbaru paling sedikit.",
    flowPreview: "Scan Pool → Pick Lowest Request Count",
  },
  {
    value: "random",
    category: "load-balance",
    labelEn: "Random",
    labelId: "Acak",
    descEn: "Randomly selects a target with optional session stickiness.",
    descId: "Memilih target secara acak dengan dukungan session stickiness.",
    flowPreview: "Random Pick (Sticky per Session) → Healthy Target",
  },
  {
    value: "strict-random",
    category: "load-balance",
    labelEn: "Strict Random",
    labelId: "Acak Murni",
    descEn: "Pure uniform random selection per request with no session stickiness.",
    descId: "Pemilihan acak seragam per permintaan tanpa session stickiness.",
    flowPreview: "Every Request → Pure Uniform Random(Pool)",
  },

  // Quota & Cost
  {
    value: "cost-optimized",
    category: "quota-cost",
    labelEn: "Cost Optimized",
    labelId: "Optimasi Biaya",
    descEn: "Prefers free-tier and lowest-cost models first before paid fallbacks.",
    descId: "Memprioritaskan model gratis dan termurah sebelum fallback berbayar.",
    flowPreview: "Free / Cheapest Model → Mid-tier → Expensive Fallback",
  },
  {
    value: "quota-weighted",
    category: "quota-cost",
    labelEn: "Quota Weighted",
    labelId: "Berbobot Kuota",
    descEn: "Distributes traffic proportionally based on remaining quota and weights.",
    descId: "Membagi trafik berdasarkan sisa kuota dan bobot langkah.",
    flowPreview: "Weight % × Remaining Quota → Proportional Pick",
    usesWeights: true,
  },
  {
    value: "headroom",
    category: "quota-cost",
    labelEn: "Headroom",
    labelId: "Headroom",
    descEn: "Selects the target with the highest available rate-limit headroom.",
    descId: "Memilih target dengan ruang rate-limit tersisa paling besar.",
    flowPreview: "Check Rate-Limit Headroom → Pick Most Available",
  },
  {
    value: "reset-aware",
    category: "quota-cost",
    labelEn: "Reset Aware",
    labelId: "Sadar Reset Kuota",
    descEn: "Balances across accounts using session and weekly quota reset timing.",
    descId: "Menyeimbangkan akun berdasarkan waktu reset kuota sesi dan mingguan.",
    flowPreview: "Score (Remaining Quota + Reset Timer) → Optimal Account",
  },
  {
    value: "reset-window",
    category: "quota-cost",
    labelEn: "Reset Window",
    labelId: "Jendela Reset",
    descEn: "Orders accounts by soonest upcoming quota reset window.",
    descId: "Mengurutkan akun berdasarkan jendela reset kuota terdekat.",
    flowPreview: "Sort by Soonest Reset Time → Use Earliest First",
  },
  {
    value: "expiry-first",
    category: "quota-cost",
    labelEn: "Expiry First",
    labelId: "Kedaluwarsa Terdekat",
    descEn: "Consumes soonest-expiring credits or subscriptions first.",
    descId: "Menghabiskan kredit atau paket yang paling cepat kedaluwarsa lebih dulu.",
    flowPreview: "Soonest Expiry Account → Later Expiry → Permanent Key",
  },

  // Context & Cache
  {
    value: "context-optimized",
    category: "context-cache",
    labelEn: "Context Optimized",
    labelId: "Optimasi Konteks",
    descEn: "Prioritizes targets with the largest context window for long prompts.",
    descId: "Memprioritaskan target dengan jendela konteks terbesar untuk prompt panjang.",
    flowPreview: "Large Prompt → Largest Context Window Model First",
  },
  {
    value: "cache-optimized",
    category: "context-cache",
    labelEn: "Cache Optimized",
    labelId: "Optimasi Cache",
    descEn: "Deterministic sticky routing to maximize upstream prompt cache hits.",
    descId: "Routing deterministik agar prompt cache upstream selalu hit maksimal.",
    flowPreview: "Hash(Session/Prefix) → Pin Same Account → Max Cache Hit",
  },
  {
    value: "context-relay",
    category: "context-cache",
    labelEn: "Context Relay",
    labelId: "Estafet Konteks",
    descEn: "Relays long conversations to a handoff target once context threshold hits.",
    descId: "Mengalihkan percakapan panjang ke target handoff saat ambang konteks tercapai.",
    flowPreview: "Context < Threshold → Primary · ≥ Threshold → Handoff Model",
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
    flowPreview: "Analyze Prompt + Telemetry (Policy) → Best Candidate",
    usesRouterStrategy: true,
  },
  {
    value: "fusion",
    category: "intelligent",
    labelEn: "Fusion (Panel + Judge)",
    labelId: "Fusion (Panel + Juri)",
    descEn: "Fans out to panel models and synthesizes the final response via a judge model.",
    descId: "Mengirim paralel ke panel model lalu menyintesis jawaban lewat model juri.",
    flowPreview: "Prompt → [Model A ∥ Model B ∥ Model C] → Judge → Final",
    usesJudgeModel: true,
  },
  {
    value: "pipeline",
    category: "intelligent",
    labelEn: "Pipeline",
    labelId: "Pipeline",
    descEn: "Chains steps sequentially where each step's output feeds the next.",
    descId: "Menjalankan langkah berurutan di mana output satu langkah menjadi input berikutnya.",
    flowPreview: "Input → Step #1 → Output₁ → Step #2 → Output₂ → Final",
  },
  {
    value: "jev",
    category: "intelligent",
    labelEn: "JEV (Judge-Evaluate-Verify)",
    labelId: "JEV (Judge-Evaluate-Verify)",
    descEn: "Multi-stage generation, evaluation, and verification routing chain.",
    descId: "Rantai generasi, evaluasi, dan verifikasi multi-tahap.",
    flowPreview: "Generate (#1) → Evaluate (#2) → Verify (#3) → Final",
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
