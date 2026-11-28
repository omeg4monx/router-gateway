import type { ComboStrategy } from "./schemas";

export interface StrategyGuideEntry {
  howItWorksEn: string;
  howItWorksId: string;
  bestForEn: string;
  bestForId: string;
  whenToUseEn: string;
  whenToUseId: string;
}

export const STRATEGY_GUIDE_CONTENT: Record<ComboStrategy, StrategyGuideEntry> = {
  priority: {
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
  },
  "fill-first": {
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
  },
  lkgp: {
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
  },
  weighted: {
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
  },
  "round-robin": {
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
  },
  p2c: {
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
  },
  "least-used": {
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
  },
  random: {
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
  },
  "strict-random": {
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
  },
  "cost-optimized": {
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
  },
  "quota-weighted": {
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
  },
  headroom: {
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
  },
  "reset-aware": {
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
  },
  "reset-window": {
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
  },
  "expiry-first": {
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
  },
  "context-optimized": {
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
  },
  "cache-optimized": {
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
  },
  "context-relay": {
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
  },
  auto: {
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
  },
  fusion: {
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
  },
  pipeline: {
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
  },
  jev: {
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
  },
};
