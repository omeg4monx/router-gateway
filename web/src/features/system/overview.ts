import type { UILocale } from "@/shared/store/uiStore";

export interface ModuleMetadata {
  href: string | null;
  description: string;
}

interface RawModuleEntry {
  href: string | null;
  descriptions: Record<UILocale, string>;
}

const MODULE_REGISTRY: Record<string, RawModuleEntry> = {
  system: {
    href: "/dashboard/",
    descriptions: {
      en: "Health probes, build metadata, initialization & lifecycle",
      id: "Probe kesehatan, metadata build, inisialisasi & siklus hidup",
    },
  },
  auth: {
    href: "/dashboard/settings/",
    descriptions: {
      en: "Admin session cookies, password rotation & remote-bind policy",
      id: "Cookie sesi admin, rotasi kata sandi & kebijakan akses jarak jauh",
    },
  },
  settings: {
    href: "/dashboard/settings/",
    descriptions: {
      en: "Typed runtime settings, revision history & atomic rollback",
      id: "Pengaturan runtime terstruktur, riwayat revisi & rollback atomik",
    },
  },
  keys: {
    href: "/dashboard/keys/",
    descriptions: {
      en: "Gateway API keys, SHA-256 verification, model scopes & limits",
      id: "Kunci API gateway, verifikasi SHA-256, cakupan model & batas kuota",
    },
  },
  providers: {
    href: "/dashboard/providers/",
    descriptions: {
      en: "Upstream LLM adapters, OAuth PKCE & AES-256-GCM vault",
      id: "Adapter LLM upstream, OAuth PKCE & brankas kredensial AES-256-GCM",
    },
  },
  combos: {
    href: "/dashboard/combos/",
    descriptions: {
      en: "Ordered fallback chains, model aliases & cooldown routing",
      id: "Rantai failover berurutan, alias model & routing cooldown",
    },
  },
  gateway: {
    href: "/dashboard/combos/",
    descriptions: {
      en: "OpenAI, Anthropic, Responses & Gemini SSE translation proxy",
      id: "Proxy translasi SSE untuk OpenAI, Anthropic, Responses & Gemini",
    },
  },
  usage: {
    href: "/dashboard/usage/",
    descriptions: {
      en: "Async token accounting, model pricing & request telemetry",
      id: "Pencatatan token asinkron, tarif model & telemetri permintaan",
    },
  },
  proxypools: {
    href: "/dashboard/proxy-pools/",
    descriptions: {
      en: "Encrypted HTTP/SOCKS5 outbound proxy pools & routing",
      id: "Pool proxy keluar HTTP/SOCKS5 terenkripsi & perutean provider",
    },
  },
  translator: {
    href: "/dashboard/translator/",
    descriptions: {
      en: "Interactive protocol translator workbench & live SSE console",
      id: "Workbench penerjemah protokol interaktif & konsol SSE langsung",
    },
  },
  promptinjection: {
    href: "/dashboard/prompt-guard/",
    descriptions: {
      en: "Prompt injection guard, RE2 detection rules & audit trail",
      id: "Pelindung prompt injection, aturan deteksi RE2 & jejak audit",
    },
  },
  modeleditor: {
    href: "/dashboard/model-editor/",
    descriptions: {
      en: "Custom model parameter overrides & conditional bypass rules",
      id: "Override parameter model kustom & aturan bypass kondisional",
    },
  },
  tokensaver: {
    href: "/dashboard/token-saver/",
    descriptions: {
      en: "Token compression & Headroom LRU prefix cache affinity",
      id: "Kompresi token & afinitas cache prefiks Headroom LRU",
    },
  },
  security: {
    href: "/dashboard/security/",
    descriptions: {
      en: "Live security posture report & redacted control-plane audit",
      id: "Laporan postur keamanan live & log audit control-plane",
    },
  },
  notifications: {
    href: "/dashboard/notifications/",
    descriptions: {
      en: "System notifications inbox, live SSE alerts & background jobs",
      id: "Kotak notifikasi sistem, peringatan SSE & kontrol tugas latar",
    },
  },
  media: {
    href: "/dashboard/media/",
    descriptions: {
      en: "Audio, image, video, search & SSRF-guarded web fetch gateway",
      id: "Gateway audio, gambar, video, pencarian & web fetch aman SSRF",
    },
  },
  tunnel: {
    href: "/dashboard/tunnel/",
    descriptions: {
      en: "Cloudflare quick tunnel & Tailscale serve/funnel manager",
      id: "Manajer Cloudflare quick tunnel & Tailscale serve/funnel",
    },
  },
  mitm: {
    href: "/dashboard/mitm/",
    descriptions: {
      en: "Loopback TLS MITM proxy, ECDSA Root CA & IDE host aliases",
      id: "Proxy MITM TLS loopback, ECDSA Root CA & alias host IDE",
    },
  },
  pxpipe: {
    href: "/dashboard/pxpipe/",
    descriptions: {
      en: "In-process request transform pipeline & savings telemetry",
      id: "Pipeline transformasi permintaan in-process & telemetri hemat",
    },
  },
  telegram: {
    href: "/dashboard/telegram/",
    descriptions: {
      en: "Vaulted Telegram Bot controller, webhook & eventbus alerts",
      id: "Kontroler Bot Telegram terenkripsi, webhook & peringatan event",
    },
  },
  arena: {
    href: "/dashboard/arena/",
    descriptions: {
      en: "Multi-model comparison arena, leaderboard & Basic Chat",
      id: "Arena perbandingan multi-model, papan peringkat & obrolan dasar",
    },
  },
  plugins: {
    href: "/dashboard/plugins/",
    descriptions: {
      en: "Extensible runtime plugin registry & lifecycle controls",
      id: "Registri plugin runtime ekstensibel & kontrol siklus hidup",
    },
  },
  importer: {
    href: "/dashboard/importer/",
    descriptions: {
      en: "Legacy database & configuration migration importer",
      id: "Alat impor migrasi database & konfigurasi versi lama",
    },
  },
};

export function getModuleMetadata(moduleId: string, locale: UILocale): ModuleMetadata {
  const entry = MODULE_REGISTRY[moduleId];
  if (!entry) {
    return {
      href: null,
      description:
        locale === "id"
          ? "Modul aktif pada control plane kernel"
          : "Active kernel control plane module",
    };
  }
  return {
    href: entry.href,
    description: entry.descriptions[locale],
  };
}

export function filterModules(
  modules: readonly string[],
  query: string,
  locale: UILocale,
): string[] {
  const trimmed = query.trim().toLowerCase();
  if (!trimmed) {
    return [...modules];
  }
  const terms = trimmed.split(/\s+/).filter(Boolean);
  return modules.filter((mod) => {
    const entry = MODULE_REGISTRY[mod];
    const haystack = entry
      ? `${mod} ${entry.descriptions.en} ${entry.descriptions.id}`.toLowerCase()
      : `${mod} ${getModuleMetadata(mod, locale).description}`.toLowerCase();
    return terms.every((term) => haystack.includes(term));
  });
}

export function formatClockUptime(totalSeconds: number): string {
  const safe = Math.max(0, Math.floor(totalSeconds));
  const days = Math.floor(safe / 86_400);
  const hours = Math.floor((safe % 86_400) / 3_600);
  const minutes = Math.floor((safe % 3_600) / 60);
  const seconds = safe % 60;

  const pad = (n: number) => String(n).padStart(2, "0");
  if (days > 0) {
    return `${days}d ${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`;
  }
  return `${pad(hours)}h ${pad(minutes)}m ${pad(seconds)}s`;
}

export function formatActivityOriginLabel(
  keyId: string,
  keyNameMap: ReadonlyMap<string, string>,
  isId: boolean,
): string {
  if (keyId === "admin-internal") {
    return isId ? "Studio Admin" : "Admin Studio";
  }
  if (keyId === "admin-probe") {
    return isId ? "Uji Model" : "Model Probe";
  }
  return keyNameMap.get(keyId) ?? keyId;
}

export type SnippetFormat = "openai" | "anthropic";
export type SnippetLang = "curl" | "python" | "node";

export function generateEndpointSnippet(options: {
  format: SnippetFormat;
  lang: SnippetLang;
  baseUrl: string;
  apiKey?: string;
  model?: string;
}): string {
  const origin = options.baseUrl.replace(/\/+$/, "") || "http://127.0.0.1:20128";
  const key = options.apiKey?.trim() || "sk-your-api-key";
  const model =
    options.model?.trim() || (options.format === "anthropic" ? "claude-sonnet-4-5" : "gpt-4o");

  if (options.format === "anthropic") {
    if (options.lang === "python") {
      return `from anthropic import Anthropic\n\nclient = Anthropic(\n    base_url="${origin}",\n    api_key="${key}",\n)\n\nmsg = client.messages.create(\n    model="${model}",\n    max_tokens=256,\n    messages=[{"role": "user", "content": "Hello!"}],\n)\nprint(msg.content[0].text)`;
    }
    if (options.lang === "node") {
      return `import Anthropic from "@anthropic-ai/sdk";\n\nconst client = new Anthropic({\n  baseURL: "${origin}",\n  apiKey: "${key}",\n});\n\nconst msg = await client.messages.create({\n  model: "${model}",\n  max_tokens: 256,\n  messages: [{ role: "user", content: "Hello!" }],\n});\nconsole.log(msg.content);`;
    }
    return `curl ${origin}/v1/messages \\\n  -H "Content-Type: application/json" \\\n  -H "x-api-key: ${key}" \\\n  -H "anthropic-version: 2023-06-01" \\\n  -d '{"model":"${model}","max_tokens":256,"messages":[{"role":"user","content":"Hello!"}]}'`;
  }

  if (options.lang === "python") {
    return `from openai import OpenAI\n\nclient = OpenAI(\n    base_url="${origin}/v1",\n    api_key="${key}",\n)\n\nres = client.chat.completions.create(\n    model="${model}",\n    messages=[{"role": "user", "content": "Hello!"}],\n)\nprint(res.choices[0].message.content)`;
  }
  if (options.lang === "node") {
    return `import OpenAI from "openai";\n\nconst client = new OpenAI({\n  baseURL: "${origin}/v1",\n  apiKey: "${key}",\n});\n\nconst res = await client.chat.completions.create({\n  model: "${model}",\n  messages: [{ role: "user", content: "Hello!" }],\n});\nconsole.log(res.choices[0].message.content);`;
  }
  return `curl ${origin}/v1/chat/completions \\\n  -H "Content-Type: application/json" \\\n  -H "Authorization: Bearer ${key}" \\\n  -d '{"model":"${model}","messages":[{"role":"user","content":"Hello!"}]}'`;
}

export interface OverviewCopy {
  eyebrow: string;
  title: string;
  subtitle: string;
  syncedPrefix: string;
  notSynced: string;
  refreshLabel: string;
  refreshingLabel: string;
  uninitializedBanner: string;
  runSetupCta: string;
  unreachableTitle: string;
  unreachableFallback: string;
  retryLabel: string;
  gatewayStatusLabel: string;
  statusOnline: string;
  statusDegraded: string;
  statusOffline: string;
  statusChecking: string;
  sqliteLabel: string;
  sqliteHealthy: string;
  sqliteError: string;
  sqliteUnknown: string;
  uptimeLabel: string;
  uptimeCaption: string;
  providersKeysLabel: string;
  providersUnit: string;
  keysUnit: string;
  manageProvidersLink: string;
  traffic24hLabel: string;
  requestsUnit: string;
  tokensUnit: string;
  viewUsageLink: string;
  quickActionsHeading: string;
  quickActions: {
    providers: { title: string; subtitle: string };
    keys: { title: string; subtitle: string };
    models: { title: string; subtitle: string };
    chat: { title: string; subtitle: string };
  };
  modulesTitle: string;
  modulesFilterPlaceholder: string;
  modulesLoading: string;
  modulesEmpty: string;
  modulesNoMatch: string;
  postureTitle: string;
  versionLabel: string;
  commitLabel: string;
  buildTimeLabel: string;
  goRuntimeLabel: string;
  initializedLabel: string;
  loginPolicyLabel: string;
  loadingLabel: string;
  unavailableLabel: string;
  activeBadge: string;
  pendingBadge: string;
  enforcedBadge: string;
  loopbackOpenBadge: string;
  copyCommitLabel: string;
  copiedLabel: string;
  endpointLabel: string;
  copyEndpointLabel: string;
  openSettingsCta: string;
  recentRequestsTitle: string;
  recentRequestsViewAll: string;
  recentRequestsEmptyTitle: string;
  recentRequestsEmptyBody: string;
  recentRequestsEmptyCta: string;
  colModel: string;
  colTokens: string;
  colLatency: string;
  colStatus: string;
  colTime: string;
  providersPanelTitle: string;
  providersPanelManage: string;
  providersEmptyTitle: string;
  providersEmptyBody: string;
  providersEmptyCta: string;
  providerEnabled: string;
  providerDisabled: string;
  modelsCountLabel: (count: number) => string;
}

const OVERVIEW_COPY: Record<UILocale, OverviewCopy> = {
  en: {
    eyebrow: "System Telemetry",
    title: "Control Plane Overview",
    subtitle: "Live runtime health, operational activity, and mounted gateway modules.",
    syncedPrefix: "Synced",
    notSynced: "Not synced",
    refreshLabel: "Refresh",
    refreshingLabel: "Refreshing",
    uninitializedBanner:
      "Administrator password has not been initialized yet. Complete first-run configuration before exposing this instance.",
    runSetupCta: "Run Setup",
    unreachableTitle: "Control plane endpoint unreachable",
    unreachableFallback: "Failed to reach gateway backend",
    retryLabel: "Retry",
    gatewayStatusLabel: "Gateway & Store",
    statusOnline: "Online",
    statusDegraded: "Degraded",
    statusOffline: "Offline",
    statusChecking: "Checking",
    sqliteLabel: "SQLite",
    sqliteHealthy: "Healthy",
    sqliteError: "Error",
    sqliteUnknown: "Unknown",
    uptimeLabel: "Process Uptime",
    uptimeCaption: "Elapsed since kernel start",
    providersKeysLabel: "Providers & Keys",
    providersUnit: "providers",
    keysUnit: "keys",
    manageProvidersLink: "Manage",
    traffic24hLabel: "24h Traffic",
    requestsUnit: "req",
    tokensUnit: "tokens",
    viewUsageLink: "Analytics",
    quickActionsHeading: "Quick Actions",
    quickActions: {
      providers: {
        title: "Providers",
        subtitle: "Connect upstream LLM accounts & OAuth",
      },
      keys: {
        title: "API Keys",
        subtitle: "Issue scoped keys & rate limits",
      },
      models: {
        title: "Models & Routing",
        subtitle: "Configure failover chains & aliases",
      },
      chat: {
        title: "Basic Chat",
        subtitle: "Test live routing in the studio",
      },
    },
    modulesTitle: "Registered Kernel Modules",
    modulesFilterPlaceholder: "Filter modules...",
    modulesLoading: "Querying kernel module registry...",
    modulesEmpty: "No active modules reported.",
    modulesNoMatch: "No modules match your filter.",
    postureTitle: "Build & Security Posture",
    versionLabel: "Version",
    commitLabel: "Git Commit",
    buildTimeLabel: "Build Time",
    goRuntimeLabel: "Go Runtime",
    initializedLabel: "Initialized",
    loginPolicyLabel: "Login Policy",
    loadingLabel: "Loading...",
    unavailableLabel: "Unavailable",
    activeBadge: "Active",
    pendingBadge: "Pending",
    enforcedBadge: "Enforced",
    loopbackOpenBadge: "Loopback Open",
    copyCommitLabel: "Copy commit hash",
    copiedLabel: "Copied",
    endpointLabel: "Endpoint",
    copyEndpointLabel: "Copy base endpoint URL",
    openSettingsCta: "Open System Settings",
    recentRequestsTitle: "Recent Requests",
    recentRequestsViewAll: "View all",
    recentRequestsEmptyTitle: "No traffic recorded yet",
    recentRequestsEmptyBody:
      "Send an API request or test a route in Basic Chat to see live gateway activity here.",
    recentRequestsEmptyCta: "Open Basic Chat",
    colModel: "Model",
    colTokens: "Tokens",
    colLatency: "Latency",
    colStatus: "Status",
    colTime: "Time",
    providersPanelTitle: "Provider Connections",
    providersPanelManage: "Manage",
    providersEmptyTitle: "No providers connected",
    providersEmptyBody:
      "Add an upstream OpenAI, Anthropic, or Gemini provider to start routing requests.",
    providersEmptyCta: "Add Provider",
    providerEnabled: "Active",
    providerDisabled: "Paused",
    modelsCountLabel: (count) => `${count} ${count === 1 ? "model" : "models"}`,
  },
  id: {
    eyebrow: "Telemetri Sistem",
    title: "Ringkasan Control Plane",
    subtitle: "Kesehatan runtime langsung, aktivitas operasional, dan modul gateway yang aktif.",
    syncedPrefix: "Sinkron",
    notSynced: "Belum sinkron",
    refreshLabel: "Muat Ulang",
    refreshingLabel: "Memuat...",
    uninitializedBanner:
      "Kata sandi administrator belum diinisialisasi. Selesaikan pengaturan awal sebelum mengekspos layanan ini.",
    runSetupCta: "Buka Pengaturan Awal",
    unreachableTitle: "Endpoint control plane tidak terjangkau",
    unreachableFallback: "Gagal terhubung ke backend gateway",
    retryLabel: "Coba Lagi",
    gatewayStatusLabel: "Gateway & Penyimpanan",
    statusOnline: "Aktif",
    statusDegraded: "Terdegradasi",
    statusOffline: "Offline",
    statusChecking: "Memeriksa",
    sqliteLabel: "SQLite",
    sqliteHealthy: "Normal",
    sqliteError: "Bermasalah",
    sqliteUnknown: "Tidak diketahui",
    uptimeLabel: "Waktu Aktif",
    uptimeCaption: "Sejak kernel dijalankan",
    providersKeysLabel: "Provider & Kunci API",
    providersUnit: "provider",
    keysUnit: "kunci",
    manageProvidersLink: "Kelola",
    traffic24hLabel: "Trafik 24 Jam",
    requestsUnit: "req",
    tokensUnit: "token",
    viewUsageLink: "Analitik",
    quickActionsHeading: "Aksi Cepat",
    quickActions: {
      providers: {
        title: "Provider",
        subtitle: "Hubungkan akun LLM upstream & OAuth",
      },
      keys: {
        title: "Kunci API",
        subtitle: "Terbitkan kunci dengan batas kuota",
      },
      models: {
        title: "Model & Routing",
        subtitle: "Atur rantai failover & alias model",
      },
      chat: {
        title: "Obrolan Dasar",
        subtitle: "Uji routing langsung di studio",
      },
    },
    modulesTitle: "Modul Kernel Terdaftar",
    modulesFilterPlaceholder: "Filter modul...",
    modulesLoading: "Memuat daftar modul kernel...",
    modulesEmpty: "Belum ada modul aktif yang dilaporkan.",
    modulesNoMatch: "Tidak ada modul yang cocok dengan filter.",
    postureTitle: "Build & Postur Keamanan",
    versionLabel: "Versi",
    commitLabel: "Git Commit",
    buildTimeLabel: "Waktu Build",
    goRuntimeLabel: "Runtime Go",
    initializedLabel: "Inisialisasi",
    loginPolicyLabel: "Kebijakan Login",
    loadingLabel: "Memuat...",
    unavailableLabel: "Tidak tersedia",
    activeBadge: "Aktif",
    pendingBadge: "Menunggu",
    enforcedBadge: "Wajib Login",
    loopbackOpenBadge: "Loopback Terbuka",
    copyCommitLabel: "Salin hash commit",
    copiedLabel: "Tersalin",
    endpointLabel: "Endpoint",
    copyEndpointLabel: "Salin URL base endpoint",
    openSettingsCta: "Buka Pengaturan Sistem",
    recentRequestsTitle: "Permintaan Terbaru",
    recentRequestsViewAll: "Lihat semua",
    recentRequestsEmptyTitle: "Belum ada trafik tercatat",
    recentRequestsEmptyBody:
      "Kirim permintaan API atau uji rute di Obrolan Dasar untuk melihat aktivitas gateway di sini.",
    recentRequestsEmptyCta: "Buka Obrolan Dasar",
    colModel: "Model",
    colTokens: "Token",
    colLatency: "Latensi",
    colStatus: "Status",
    colTime: "Waktu",
    providersPanelTitle: "Koneksi Provider",
    providersPanelManage: "Kelola",
    providersEmptyTitle: "Belum ada provider terhubung",
    providersEmptyBody:
      "Tambahkan provider upstream OpenAI, Anthropic, atau Gemini untuk mulai merutekan permintaan.",
    providersEmptyCta: "Tambah Provider",
    providerEnabled: "Aktif",
    providerDisabled: "Nonaktif",
    modelsCountLabel: (count) => `${count} model`,
  },
};

export function getOverviewCopy(locale: UILocale): OverviewCopy {
  return OVERVIEW_COPY[locale];
}
