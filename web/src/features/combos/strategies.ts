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
  howItWorksEn: string;
  howItWorksId: string;
  bestForEn: string;
  bestForId: string;
  whenToUseEn: string;
  whenToUseId: string;
  flowPreview: string;
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
    howItWorksEn:
      "Always sends the request to step #1 first. If step #1 returns an error, rate limit (429), timeout, or is currently in cooldown, the router immediately falls back to step #2, then #3, in strict list order.",
    howItWorksId:
      "Selalu mengirim request ke model urutan #1 terlebih dahulu. Jika #1 mengalami error, rate-limit (429), timeout, atau sedang cooldown, router langsung pindah ke #2, lalu #3 secara berurutan.",
    bestForEn:
      "High-reliability production apps with one flagship primary model and 1–3 backup models.",
    bestForId:
      "Aplikasi produksi yang butuh uptime tinggi dengan 1 model utama terbaik dan 1–3 model cadangan.",
    whenToUseEn:
      "Use as your default strategy whenever you want predictable quality from a primary model while staying safe from provider outages or rate limits.",
    whenToUseId:
      "Gunakan sebagai pilihan default saat kamu ingin kualitas konsisten dari model utama, tapi tetap aman jika provider utama limit atau down.",
    flowPreview: "#1 Primary → (error/429) → #2 Backup → #3 Fallback",
  },
  {
    value: "fill-first",
    category: "failover",
    labelEn: "Fill First",
    labelId: "Isi Pertama",
    descEn: "Exhausts the primary target's capacity before cascading to backups.",
    descId: "Memaksimalkan kapasitas target pertama sebelum beralih ke cadangan.",
    howItWorksEn:
      "Drains all available quota/capacity of the first target until it hits rate limits or cooldown, then moves to the next target in order and stays there until exhausted.",
    howItWorksId:
      "Menghabiskan kapasitas/kuota target pertama sampai penuh atau terkena cooldown, lalu baru berpindah dan menghabiskan target berikutnya secara berurutan.",
    bestForEn:
      "Maximizing free tiers, fixed monthly quotas, or subscription allowances before touching pay-as-you-go keys.",
    bestForId:
      "Memaksimalkan paket gratis, kuota bulanan, atau langganan OAuth sebelum menyentuh API key berbayar.",
    whenToUseEn:
      "Use when you have free/subscription accounts first in the list and paid API keys at the bottom as overflow protection.",
    whenToUseId:
      "Gunakan saat kamu menaruh akun gratis/subscription di urutan atas dan API key berbayar di bawah sebagai pengaman overflow.",
    flowPreview: "#1 (100% sampai habis) → #2 (sampai habis) → #3",
  },
  {
    value: "lkgp",
    category: "failover",
    labelEn: "Last Known Good (LKGP)",
    labelId: "Provider Sehat Terakhir (LKGP)",
    descEn: "Pins requests to the last known healthy target and fails over automatically.",
    descId: "Mengunci permintaan ke target sehat terakhir dan failover otomatis.",
    howItWorksEn:
      "Remembers which target succeeded most recently and keeps routing traffic to it. When that target fails, it fails over to the next healthy candidate and promotes the new survivor as the active target.",
    howItWorksId:
      "Mengingat target mana yang terakhir kali sukses dan terus mengarahkan request ke sana. Jika target itu gagal, router failover ke kandidat sehat berikutnya dan menjadikannya target aktif baru.",
    bestForEn:
      "Unstable upstream providers where you want to avoid repeatedly retrying a flaky #1 target on every request.",
    bestForId:
      "Provider yang sering tidak stabil agar router tidak terus-menerus membuang waktu mencoba #1 yang sedang bermasalah.",
    whenToUseEn:
      "Use when latency penalties from repeated primary retries hurt UX, and all targets in your combo have comparable quality.",
    whenToUseId:
      "Gunakan ketika retry berulang ke model #1 bikin respons lambat, dan semua model di combo punya kualitas yang setara.",
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
    howItWorksEn:
      "Selects a target probabilistically according to each step's percentage weight (e.g., 70% / 20% / 10%). If the chosen target fails, it automatically fails over to the remaining healthy steps.",
    howItWorksId:
      "Memilih target secara probabilistik sesuai bobot persen tiap langkah (mis. 70% / 20% / 10%). Jika target terpilih gagal, otomatis failover ke langkah sehat lainnya.",
    bestForEn:
      "Canary rollouts, A/B testing new models, or splitting load across accounts with unequal rate limits.",
    bestForId:
      "Canary rollout, A/B testing model baru, atau membagi beban ke akun yang punya batas RPM berbeda-beda.",
    whenToUseEn:
      "Use when you want exact traffic ratios—like sending 80% to your main model and 20% to a cheaper or experimental model.",
    whenToUseId:
      "Gunakan saat kamu ingin rasio trafik spesifik—misalnya 80% ke model stabil dan 20% untuk uji coba model baru atau yang lebih murah.",
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
    howItWorksEn:
      "Advances an atomic counter on every request: request 1 goes to #1, request 2 to #2, request 3 to #3, then wraps back to #1. Cooldown targets are skipped automatically.",
    howItWorksId:
      "Memutar giliran secara merata di setiap request: request ke-1 ke #1, ke-2 ke #2, ke-3 ke #3, lalu kembali ke #1. Target yang sedang cooldown dilewati otomatis.",
    bestForEn:
      "Multiple equal-capacity accounts or equivalent models where you want perfectly even request distribution.",
    bestForId:
      "Beberapa akun dengan kapasitas setara atau model sekelas yang ingin dibagi bebannya secara adil.",
    whenToUseEn:
      "Use when you have 2+ API keys/OAuth accounts with similar RPM limits and want to spread bursts evenly to avoid hitting 429.",
    whenToUseId:
      "Gunakan saat punya 2+ API key/OAuth dengan limit RPM mirip dan ingin meratakan lonjakan request agar tidak cepat kena 429.",
    flowPreview: "Req 1 → #1 · Req 2 → #2 · Req 3 → #3 · Req 4 → #1",
  },
  {
    value: "p2c",
    category: "load-balance",
    labelEn: "Power of Two Choices (P2C)",
    labelId: "Power of Two Choices (P2C)",
    descEn: "Samples two candidates and picks the less-loaded healthy target.",
    descId: "Memilih dua kandidat acak lalu mengambil target sehat dengan beban terendah.",
    howItWorksEn:
      "Randomly picks two healthy targets from the pool, compares their recent usage/load, and routes the request to the less-loaded one—preventing herd effects without global scanning.",
    howItWorksId:
      "Mengambil 2 kandidat sehat secara acak, membandingkan beban/penggunaan terbarunya, lalu memilih yang paling ringan—mencegah penumpukan trafik di satu titik.",
    bestForEn: "Medium-to-large pools (3+ models or accounts) under concurrent, spiky traffic.",
    bestForId: "Pool berisi 3+ model atau akun yang menerima trafik paralel dan fluktuatif.",
    whenToUseEn:
      "Use when Round Robin still causes uneven load because some requests are much slower or heavier than others.",
    whenToUseId:
      "Gunakan saat Round Robin kurang seimbang karena sebagian request jauh lebih berat/lama dibanding yang lain.",
    flowPreview: "Sample 2 Random → Compare Load → Pick Lighter Target",
  },
  {
    value: "least-used",
    category: "load-balance",
    labelEn: "Least Used",
    labelId: "Paling Jarang Dipakai",
    descEn: "Routes to the target with the fewest recent requests.",
    descId: "Mengarahkan ke target dengan jumlah penggunaan terbaru paling sedikit.",
    howItWorksEn:
      "Tracks selection counts across all healthy targets and always picks the target with the lowest usage count so far.",
    howItWorksId:
      "Mencatat jumlah pemakaian tiap target sehat dan selalu memilih target dengan angka pemakaian paling rendah.",
    bestForEn:
      "Keeping quota consumption balanced across multiple accounts even after temporary outages or manual pauses.",
    bestForId:
      "Menjaga konsumsi kuota tetap seimbang antar akun meskipun ada akun yang sempat cooldown atau baru ditambahkan.",
    whenToUseEn:
      "Use when you add/remove accounts dynamically and want newly idle accounts to catch up fairly.",
    whenToUseId:
      "Gunakan ketika kamu sering menambah/mengaktifkan akun baru dan ingin akun yang paling jarang dipakai diprioritaskan dulu.",
    flowPreview: "Scan Pool → Pick Lowest Request Count",
  },
  {
    value: "random",
    category: "load-balance",
    labelEn: "Random",
    labelId: "Acak",
    descEn: "Randomly selects a target with optional session stickiness.",
    descId: "Memilih target secara acak dengan dukungan session stickiness.",
    howItWorksEn:
      "Picks a healthy target uniformly at random while respecting session stickiness when enabled, so multi-turn conversations stay consistent.",
    howItWorksId:
      "Memilih target sehat secara acak merata sambil tetap menghormati session stickiness jika aktif, sehingga percakapan multi-turn tetap konsisten.",
    bestForEn:
      "Simple stateless load spreading across interchangeable models with conversational continuity.",
    bestForId:
      "Penyebaran beban sederhana ke model-model setara sambil tetap menjaga konsistensi sesi chat.",
    whenToUseEn:
      "Use when you want lightweight random distribution without needing strict rotation order.",
    whenToUseId: "Gunakan saat ingin distribusi acak yang ringan tanpa perlu urutan rotasi ketat.",
    flowPreview: "Random Pick (Sticky per Session) → Healthy Target",
  },
  {
    value: "strict-random",
    category: "load-balance",
    labelEn: "Strict Random",
    labelId: "Acak Murni",
    descEn: "Pure uniform random selection per request with no session stickiness.",
    descId: "Pemilihan acak seragam per permintaan tanpa session stickiness.",
    howItWorksEn:
      "Selects a healthy target purely at random on every single request, completely ignoring session affinity or previous turns.",
    howItWorksId:
      "Memilih target sehat murni secara acak di setiap request, mengabaikan session stickiness atau turn sebelumnya.",
    bestForEn:
      "Unbiased model benchmarking, synthetic data generation, or blind evaluation sampling.",
    bestForId: "Benchmark model tanpa bias, generasi data sintetis, atau sampling evaluasi acak.",
    whenToUseEn:
      "Use when every single request must be independently randomized across your target list.",
    whenToUseId:
      "Gunakan ketika setiap request harus benar-benar diacak secara independen tanpa terikat sesi.",
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
    howItWorksEn:
      "Sorts healthy targets by estimated token cost (free/OAuth tiers first, then cheapest per-token models) and only escalates to more expensive models when cheaper ones fail or hit limits.",
    howItWorksId:
      "Mengurutkan target sehat dari biaya token termurah (tier gratis/OAuth dulu, lalu model murah), dan hanya naik ke model mahal jika yang murah gagal atau limit.",
    bestForEn:
      "High-volume workloads, background batch jobs, summarization, and cost-sensitive products.",
    bestForId:
      "Workload volume tinggi, batch job, ringkasan teks, dan proyek yang ketat soal anggaran API.",
    whenToUseEn:
      "Use when keeping monthly spend minimal matters more than always using the most expensive flagship model.",
    whenToUseId:
      "Gunakan saat efisiensi biaya bulanan lebih penting daripada selalu memakai model flagship termahal.",
    flowPreview: "Free / Cheapest Model → Mid-tier → Expensive Fallback",
  },
  {
    value: "quota-weighted",
    category: "quota-cost",
    labelEn: "Quota Weighted",
    labelId: "Berbobot Kuota",
    descEn: "Distributes traffic proportionally based on remaining quota and weights.",
    descId: "Membagi trafik berdasarkan sisa kuota dan bobot langkah.",
    howItWorksEn:
      "Combines your configured step weights with each account's remaining quota signals so targets with larger remaining budgets absorb proportionally more traffic.",
    howItWorksId:
      "Menggabungkan bobot persen yang kamu atur dengan sinyal sisa kuota tiap akun, sehingga target dengan kuota tersisa lebih besar menampung trafik lebih banyak.",
    bestForEn:
      "Mixed-tier accounts (e.g., 1 Pro account + 2 Free accounts) where you want balanced exhaustion.",
    bestForId:
      "Campuran akun beda tier (mis. 1 akun Pro + 2 akun Free) agar kuota habis secara proporsional.",
    whenToUseEn:
      "Use when simple weighted splits aren't enough because accounts have different remaining quota pools.",
    whenToUseId:
      "Gunakan ketika pembagian bobot statis kurang pas karena sisa kuota masing-masing akun berbeda.",
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
    howItWorksEn:
      "Inspects live rate-limit/usage headroom across healthy targets and sends traffic to whichever target has the most breathing room before throttling.",
    howItWorksId:
      "Memeriksa ruang rate-limit yang masih longgar di semua target sehat dan mengarahkan request ke target yang paling jauh dari batas throttle.",
    bestForEn:
      "Burst-heavy agents, coding loops, and parallel tool-calling pipelines prone to HTTP 429 errors.",
    bestForId:
      "Agent dengan trafik burst tinggi, loop coding, atau parallel tool-calling yang rawan kena HTTP 429.",
    whenToUseEn:
      "Use when avoiding 429 rate-limit interruptions is your top priority across multiple accounts.",
    whenToUseId:
      "Gunakan saat menghindari error 429 adalah prioritas utama di tengah trafik paralel yang padat.",
    flowPreview: "Check Rate-Limit Headroom → Pick Most Available",
  },
  {
    value: "reset-aware",
    category: "quota-cost",
    labelEn: "Reset Aware",
    labelId: "Sadar Reset Kuota",
    descEn: "Balances across accounts using session and weekly quota reset timing.",
    descId: "Menyeimbangkan akun berdasarkan waktu reset kuota sesi dan mingguan.",
    howItWorksEn:
      "Weighs both remaining quota and upcoming session/weekly reset schedules so you burn quota that is about to refresh while preserving freshly reset accounts.",
    howItWorksId:
      "Mempertimbangkan sisa kuota sekaligus jadwal reset sesi/mingguan agar kuota yang segera di-reset dimanfaatkan dulu, sementara akun yang baru reset tetap awet.",
    bestForEn:
      "Multi-account OAuth subscriptions (Claude, Codex, Gemini, Copilot) with rolling reset timers.",
    bestForId:
      "Multi-akun OAuth (Claude, Codex, Gemini, Copilot) yang punya jadwal reset kuota berkala.",
    whenToUseEn:
      "Use when juggling several subscription accounts and you want smart pacing across reset cycles.",
    whenToUseId:
      "Gunakan saat mengelola beberapa akun langganan dan ingin pemakaian kuota diatur cerdas mengikuti siklus reset.",
    flowPreview: "Score (Remaining Quota + Reset Timer) → Optimal Account",
  },
  {
    value: "reset-window",
    category: "quota-cost",
    labelEn: "Reset Window",
    labelId: "Jendela Reset",
    descEn: "Orders accounts by soonest upcoming quota reset window.",
    descId: "Mengurutkan akun berdasarkan jendela reset kuota terdekat.",
    howItWorksEn:
      "Prioritizes whichever healthy account has the earliest upcoming quota reset window, ensuring unused quota doesn't go to waste before its reset.",
    howItWorksId:
      "Memprioritaskan akun sehat yang waktu reset kuotanya paling dekat, memastikan jatah kuota tidak hangus sebelum di-reset.",
    bestForEn: "Maximizing utilization across rotating weekly or 5-hour quota windows.",
    bestForId:
      "Memaksimalkan penggunaan kuota pada akun dengan jendela reset mingguan atau per 5 jam.",
    whenToUseEn: "Use when you want a deterministic 'use soonest-resetting account first' policy.",
    whenToUseId:
      "Gunakan jika kamu ingin aturan tegas: habiskan dulu akun yang paling cepat masuk jadwal reset.",
    flowPreview: "Sort by Soonest Reset Time → Use Earliest First",
  },
  {
    value: "expiry-first",
    category: "quota-cost",
    labelEn: "Expiry First",
    labelId: "Kedaluwarsa Terdekat",
    descEn: "Consumes soonest-expiring credits or subscriptions first.",
    descId: "Menghabiskan kredit atau paket yang paling cepat kedaluwarsa lebih dulu.",
    howItWorksEn:
      "Orders targets by token/credit expiration timestamp so promotional credits, trial keys, or expiring subscriptions are consumed before permanent pay-as-you-go keys.",
    howItWorksId:
      "Mengurutkan target berdasarkan tanggal kedaluwarsa token/kredit agar kredit promo, trial, atau paket yang mau habis terpakai lebih dulu sebelum kunci permanen.",
    bestForEn: "Promo credits, hackathon grants, trial accounts, and expiring monthly plans.",
    bestForId:
      "Kredit promo, grant gratis, akun trial, dan paket bulanan yang punya tanggal hangus.",
    whenToUseEn:
      "Use whenever you have expiring credits alongside normal billing accounts so zero credit is wasted.",
    whenToUseId:
      "Gunakan saat punya kredit berbatas waktu bersama akun reguler agar tidak ada saldo yang hangus sia-sia.",
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
    howItWorksEn:
      "Evaluates the incoming prompt size and prioritizes targets capable of handling larger context windows (e.g., 1M / 200K tokens) to prevent context-length overflow errors.",
    howItWorksId:
      "Melihat kebutuhan panjang prompt lalu memprioritaskan model dengan jendela konteks terbesar (mis. 1M / 200K token) agar tidak gagal karena context overflow.",
    bestForEn: "Whole-repository analysis, large PDF/log ingestion, and long coding sessions.",
    bestForId:
      "Analisis seluruh codebase, membaca dokumen/log besar, dan sesi coding yang sangat panjang.",
    whenToUseEn:
      "Use when your combo mixes small-context fast models with huge-context models and prompt length varies widely.",
    whenToUseId:
      "Gunakan saat combo mencampur model konteks kecil dengan model konteks jumbo dan ukuran prompt sering berubah-ubah.",
    flowPreview: "Large Prompt → Largest Context Window Model First",
  },
  {
    value: "cache-optimized",
    category: "context-cache",
    labelEn: "Cache Optimized",
    labelId: "Optimasi Cache",
    descEn: "Deterministic sticky routing to maximize upstream prompt cache hits.",
    descId: "Routing deterministik agar prompt cache upstream selalu hit maksimal.",
    howItWorksEn:
      "Deterministically sticks repeated requests with the same prefix/session to the exact same provider & account so upstream prompt caching (Anthropic/OpenAI/Gemini) stays warm.",
    howItWorksId:
      "Mengunci request berulang dengan prefix/sesi yang sama secara deterministik ke provider & akun yang sama agar prompt cache upstream selalu panas (cache hit tinggi).",
    bestForEn:
      "IDE coding agents (Cursor, Claude Code, Cline) with huge repeating system prompts and tool schemas.",
    bestForId:
      "Coding agent di IDE (Cursor, Claude Code, Cline) dengan system prompt besar yang berulang di tiap turn.",
    whenToUseEn:
      "Use when you want up to 70–90% prompt-cache savings and lower TTFT latency on multi-turn agent sessions.",
    whenToUseId:
      "Gunakan saat ingin hemat biaya token 70–90% lewat prompt cache dan mempercepat respons pada sesi agent multi-turn.",
    flowPreview: "Hash(Session/Prefix) → Pin Same Account → Max Cache Hit",
  },
  {
    value: "context-relay",
    category: "context-cache",
    labelEn: "Context Relay",
    labelId: "Estafet Konteks",
    descEn: "Relays long conversations to a handoff target once context threshold hits.",
    descId: "Mengalihkan percakapan panjang ke target handoff saat ambang konteks tercapai.",
    howItWorksEn:
      "Uses your primary fast/cheap model while conversation context stays below the handoff threshold (e.g., 80%), then seamlessly relays the thread to a larger-context handoff model.",
    howItWorksId:
      "Memakai model utama yang cepat/murah selama konteks percakapan masih di bawah ambang batas (mis. 80%), lalu otomatis estafet ke model handoff berkapasitas konteks lebih besar.",
    bestForEn:
      "Long-running chat assistants and agents that start small but occasionally grow past standard context limits.",
    bestForId:
      "Asisten chat atau agent jangka panjang yang awalnya ringan tapi kadang membengkak melewati batas konteks normal.",
    whenToUseEn:
      "Use when you want fast/cheap turns early in a session and automatic handoff to a big-context model only when needed.",
    whenToUseId:
      "Gunakan saat ingin awal sesi tetap cepat/hemat, lalu pindah otomatis ke model konteks besar hanya saat percakapan mulai penuh.",
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
    howItWorksEn:
      "Inspects request characteristics and live telemetry using your chosen policy (rules, score, cost, eco, latency, fast, sla-aware, lkgp, nadir) to pick the best candidate automatically.",
    howItWorksId:
      "Menganalisis karakteristik request dan telemetri live memakai kebijakan pilihan (rules, score, cost, eco, latency, fast, sla-aware, lkgp, nadir) untuk memilih kandidat terbaik otomatis.",
    bestForEn: "General-purpose smart endpoints that serve mixed simple and complex prompts.",
    bestForId:
      "Satu endpoint pintar serba guna yang melayani campuran prompt ringan hingga kompleks.",
    whenToUseEn:
      "Use when you want one model alias for your team/apps and let the router optimize for speed, cost, or SLA dynamically.",
    whenToUseId:
      "Gunakan saat ingin 1 nama combo saja untuk semua aplikasi, lalu biarkan router memilihkan model terbaik berdasarkan kecepatan, biaya, atau SLA.",
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
    howItWorksEn:
      "Sends your prompt in parallel to all models in the panel, collects their independent answers, then asks a Judge model to compare, verify, and synthesize one superior final answer.",
    howItWorksId:
      "Mengirim prompt secara paralel ke semua model di panel, mengumpulkan jawaban masing-masing, lalu model Juri membandingkan dan menggabungkannya jadi satu jawaban terbaik.",
    bestForEn:
      "High-stakes architecture decisions, legal/medical/security analysis, and deep research where accuracy beats cost.",
    bestForId:
      "Keputusan arsitektur krusial, audit keamanan/riset mendalam, dan tugas yang menuntut akurasi maksimal dibanding biaya.",
    whenToUseEn:
      "Use when a single model isn't reliable enough and you're willing to spend extra tokens for multi-model consensus.",
    whenToUseId:
      "Gunakan ketika 1 model saja rawan halusinasi dan kamu siap memakai token ekstra demi konsensus multi-model.",
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
    howItWorksEn:
      "Executes steps one after another like an assembly line: step #1 drafts or extracts, step #2 refines or critiques using step #1's output, and so on until the final step.",
    howItWorksId:
      "Menjalankan model secara berantai seperti lini produksi: langkah #1 membuat draf/ekstraksi, outputnya masuk ke langkah #2 untuk disempurnakan/dikritik, hingga langkah terakhir.",
    bestForEn:
      "Multi-stage workflows such as Draft → Review → Polish, or Translate → Verify → Format.",
    bestForId:
      "Alur kerja bertahap seperti Draf → Review → Polish, atau Ekstraksi → Analisis → Format JSON.",
    whenToUseEn:
      "Use when splitting a complex task across specialized models (e.g., fast drafter + strong editor) produces better results.",
    whenToUseId:
      "Gunakan saat membagi tugas rumit ke beberapa model spesialis (mis. pembuat draf cepat + editor cerdas) memberi hasil lebih rapi.",
    flowPreview: "Input → Step #1 → Output₁ → Step #2 → Output₂ → Final",
    usesStepPrompt: true,
  },
  {
    value: "jev",
    category: "intelligent",
    labelEn: "JEV (Judge-Evaluate-Verify)",
    labelId: "JEV (Judge-Evaluate-Verify)",
    descEn: "Multi-stage generation, evaluation, and verification routing chain.",
    descId: "Rantai generasi, evaluasi, dan verifikasi multi-tahap.",
    howItWorksEn:
      "Runs a structured 3-stage loop across your combo steps: generates an initial answer, evaluates its quality/completeness, and verifies facts/code before returning the response.",
    howItWorksId:
      "Menjalankan siklus 3 tahap terstruktur di langkah combo: menghasilkan jawaban awal, mengevaluasi kelengkapan/kualitasnya, lalu memverifikasi fakta/kode sebelum dikirim ke pengguna.",
    bestForEn:
      "Code generation, math/logic reasoning, and compliance-sensitive answers that need built-in self-checking.",
    bestForId:
      "Generasi kode kritis, penalaran logika/matematika, dan jawaban yang wajib melewati verifikasi mandiri.",
    whenToUseEn:
      "Use when you need stronger verification than a single pass without running full parallel Fusion on every prompt.",
    whenToUseId:
      "Gunakan saat butuh verifikasi ekstra dibanding 1 kali panggil model, tanpa harus menjalankan paralel penuh seperti Fusion.",
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
