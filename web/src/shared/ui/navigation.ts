import {
  Activity,
  ArrowRightLeft,
  Bell,
  BookOpen,
  Cloud,
  Cpu,
  Database,
  Image,
  KeyRound,
  Layers,
  Lock,
  MessageSquare,
  Network,
  Puzzle,
  Radio,
  Send,
  Server,
  Shield,
  ShieldAlert,
  ShieldCheck,
  Sliders,
  SlidersHorizontal,
  Sparkles,
  Swords,
  Terminal,
  Ticket,
  Users,
  Workflow,
  Wrench,
  Zap,
} from "lucide-react";
import type { UILocale } from "@/shared/store/uiStore";

export type NavGroupId = "core" | "routing" | "network" | "access" | "workspace" | "system";

export interface NavItem {
  href: string;
  icon: typeof Activity;
  labels: Record<UILocale, string>;
  keywords: readonly string[];
}

export interface NavGroup {
  id: NavGroupId;
  icon: typeof Activity;
  labels: Record<UILocale, string>;
  defaultExpanded: boolean;
  items: readonly NavItem[];
}

export const NAV_GROUPS: readonly NavGroup[] = [
  {
    id: "core",
    icon: Activity,
    labels: {
      en: "Core & Monitoring",
      id: "Inti & Pemantauan",
    },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/",
        icon: Activity,
        labels: { en: "System Overview", id: "Ringkasan Sistem" },
        keywords: ["home", "status", "health", "monitor", "dashboard", "utama"],
      },
      {
        href: "/dashboard/usage/",
        icon: Cpu,
        labels: { en: "Usage & Telemetry", id: "Penggunaan & Telemetri" },
        keywords: ["billing", "cost", "metrics", "tokens", "analytics", "biaya", "statistik"],
      },
      {
        href: "/dashboard/notifications/",
        icon: Bell,
        labels: { en: "Notifications & Jobs", id: "Notifikasi & Tugas" },
        keywords: ["alerts", "cron", "background", "events", "peringatan"],
      },
    ],
  },
  {
    id: "routing",
    icon: Layers,
    labels: {
      en: "Routing & Providers",
      id: "Rute & Provider",
    },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/providers/",
        icon: Server,
        labels: { en: "Providers", id: "Provider" },
        keywords: ["openai", "anthropic", "gemini", "upstream", "oauth", "koneksi"],
      },
      {
        href: "/dashboard/models/",
        icon: Layers,
        labels: { en: "Models & Routing", id: "Model & Routing" },
        keywords: ["combos", "aliases", "fallback", "failover", "rute"],
      },
      {
        href: "/dashboard/model-editor/",
        icon: Sliders,
        labels: { en: "Model Editor", id: "Editor Model" },
        keywords: ["overrides", "temperature", "params", "tuning", "parameter"],
      },
      {
        href: "/dashboard/translator/",
        icon: ArrowRightLeft,
        labels: { en: "Translator", id: "Penerjemah Protokol" },
        keywords: ["protocol", "format", "convert", "openai", "anthropic", "konversi"],
      },
      {
        href: "/dashboard/token-saver/",
        icon: Zap,
        labels: { en: "Token Saver", id: "Penghemat Token" },
        keywords: ["compress", "cache", "headroom", "hemat", "optimasi"],
      },
    ],
  },
  {
    id: "network",
    icon: Network,
    labels: {
      en: "Network & Gateway",
      id: "Jaringan & Gateway",
    },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/proxy-pools/",
        icon: Network,
        labels: { en: "Proxy Pools", id: "Pool Proxy" },
        keywords: ["http", "socks5", "egress", "ip", "jaringan"],
      },
      {
        href: "/dashboard/media/",
        icon: Image,
        labels: { en: "Media Gateway", id: "Gateway Media" },
        keywords: ["images", "vision", "audio", "files", "gambar"],
      },
      {
        href: "/dashboard/tunnel/",
        icon: Radio,
        labels: { en: "Tunnel", id: "Tunnel" },
        keywords: ["cloudflare", "tailscale", "expose", "remote", "terowongan"],
      },
      {
        href: "/dashboard/mitm/",
        icon: Shield,
        labels: { en: "MITM Proxy", id: "Proxy MITM" },
        keywords: ["intercept", "tls", "ca", "certificate", "sertifikat"],
      },
      {
        href: "/dashboard/pxpipe/",
        icon: Workflow,
        labels: { en: "PXPIPE", id: "Pipeline PXPIPE" },
        keywords: ["pipeline", "middleware", "transform", "alur"],
      },
    ],
  },
  {
    id: "access",
    icon: ShieldCheck,
    labels: {
      en: "Access & Security",
      id: "Akses & Keamanan",
    },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/keys/",
        icon: KeyRound,
        labels: { en: "API Keys", id: "Kunci API" },
        keywords: ["tokens", "credentials", "quota", "rate limit", "kunci"],
      },
      {
        href: "/dashboard/teams/",
        icon: Users,
        labels: { en: "Teams", id: "Tim" },
        keywords: ["org", "members", "budget", "quota", "anggota"],
      },
      {
        href: "/dashboard/vouchers/",
        icon: Ticket,
        labels: { en: "Vouchers", id: "Voucher" },
        keywords: ["redeem", "credits", "codes", "kupon", "kredit"],
      },
      {
        href: "/dashboard/prompt-guard/",
        icon: ShieldAlert,
        labels: { en: "Prompt Guard", id: "Pelindung Prompt" },
        keywords: ["injection", "firewall", "safety", "filter", "keamanan"],
      },
      {
        href: "/dashboard/security/",
        icon: Lock,
        labels: { en: "Security", id: "Keamanan" },
        keywords: ["audit", "logs", "policy", "compliance", "keamanan"],
      },
    ],
  },
  {
    id: "workspace",
    icon: Sparkles,
    labels: {
      en: "Workspace & Studios",
      id: "Ruang Kerja & Studio",
    },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/chat/",
        icon: MessageSquare,
        labels: { en: "Basic Chat", id: "Obrolan Dasar" },
        keywords: ["playground", "prompt", "test", "messages", "chat"],
      },
      {
        href: "/dashboard/arena/",
        icon: Swords,
        labels: { en: "Model Arena", id: "Arena Model" },
        keywords: ["battle", "compare", "benchmark", "eval", "banding"],
      },
      {
        href: "/dashboard/telegram/",
        icon: Send,
        labels: { en: "Telegram Bot", id: "Bot Telegram" },
        keywords: ["bot", "messaging", "webhook", "chat"],
      },
      {
        href: "/dashboard/plugins/",
        icon: Puzzle,
        labels: { en: "Plugins", id: "Plugin" },
        keywords: ["extensions", "addons", "modules", "ekstensi"],
      },
      {
        href: "/dashboard/skills/",
        icon: BookOpen,
        labels: { en: "Skills", id: "Skill" },
        keywords: ["agents", "prompts", "playbooks", "panduan"],
      },
    ],
  },
  {
    id: "system",
    icon: Wrench,
    labels: {
      en: "System & Tools",
      id: "Sistem & Utilitas",
    },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/cloud-sync/",
        icon: Cloud,
        labels: { en: "Cloud Sync", id: "Sinkronisasi Cloud" },
        keywords: ["backup", "snapshot", "restore", "cadangan"],
      },
      {
        href: "/dashboard/importer/",
        icon: Database,
        labels: { en: "Legacy Importer", id: "Impor Data Lama" },
        keywords: ["migrate", "import", "sqlite", "migrasi"],
      },
      {
        href: "/dashboard/cli-tools/",
        icon: Terminal,
        labels: { en: "CLI Tools", id: "Alat CLI" },
        keywords: ["claude", "codex", "cursor", "setup", "terminal"],
      },
      {
        href: "/dashboard/settings/",
        icon: SlidersHorizontal,
        labels: { en: "Settings", id: "Pengaturan" },
        keywords: ["config", "preferences", "password", "locale", "theme", "konfigurasi"],
      },
    ],
  },
];

export const ALL_NAV_ITEMS: readonly NavItem[] = NAV_GROUPS.flatMap((group) => group.items);

export function normalizeRoute(path: string): string {
  const trimmed = path.trim();
  if (!trimmed || trimmed === "/") {
    return "/";
  }
  const withLeading = trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
  return withLeading.endsWith("/") ? withLeading : `${withLeading}/`;
}

export function isNavItemActive(normalizedPathname: string, itemHref: string): boolean {
  const normalizedItem = normalizeRoute(itemHref);
  if (normalizedItem === "/dashboard/") {
    return normalizedPathname === "/dashboard/";
  }
  return normalizedPathname === normalizedItem || normalizedPathname.startsWith(normalizedItem);
}

export interface ActiveNavContext {
  group: NavGroup;
  item: NavItem;
}

export function findActiveNavContext(pathname: string): ActiveNavContext | null {
  const normalized = normalizeRoute(pathname);
  for (const group of NAV_GROUPS) {
    for (const item of group.items) {
      if (isNavItemActive(normalized, item.href)) {
        return { group, item };
      }
    }
  }
  return null;
}

export function filterNavGroups(query: string, _locale: UILocale): NavGroup[] {
  const normalizedQuery = query.trim().toLowerCase();
  if (!normalizedQuery) {
    return [...NAV_GROUPS];
  }

  const terms = normalizedQuery.split(/\s+/).filter(Boolean);
  const filtered: NavGroup[] = [];

  for (const group of NAV_GROUPS) {
    const groupHaystack = `${group.labels.en} ${group.labels.id}`.toLowerCase();
    const matchingItems = group.items.filter((item) => {
      const itemHaystack =
        `${item.labels.en} ${item.labels.id} ${item.href} ${item.keywords.join(" ")} ${groupHaystack}`.toLowerCase();
      return terms.every((term) => itemHaystack.includes(term));
    });

    if (matchingItems.length > 0) {
      filtered.push({
        ...group,
        items: matchingItems,
      });
    }
  }

  return filtered;
}

export interface ShellCopy {
  controlPlane: string;
  fallbackTitle: string;
  searchPlaceholder: string;
  quickSearchLabel: string;
  clearSearchLabel: string;
  searchResultsLabel: (count: number) => string;
  emptySearchTitle: string;
  emptySearchSubtitle: string;
  resetSearchLabel: string;
  expandAllLabel: string;
  collapseAllLabel: string;
  categoriesLabel: (open: number, total: number) => string;
  expandSidebarLabel: string;
  collapseSidebarLabel: string;
  openMenuLabel: string;
  closeMenuLabel: string;
  dismissMenuLabel: string;
  signOutLabel: string;
  signingOutLabel: string;
  gatewayOnline: string;
  gatewayDegraded: string;
  gatewayOffline: string;
  gatewayChecking: string;
}

const SHELL_COPY: Record<UILocale, ShellCopy> = {
  en: {
    controlPlane: "Control Plane",
    fallbackTitle: "Gateway Monitor",
    searchPlaceholder: "Search menu...",
    quickSearchLabel: "Search",
    clearSearchLabel: "Clear search",
    searchResultsLabel: (count) => `${count} ${count === 1 ? "match" : "matches"}`,
    emptySearchTitle: "No matching menu",
    emptySearchSubtitle: "Try another keyword such as providers, keys, chat, or settings.",
    resetSearchLabel: "Reset search",
    expandAllLabel: "Expand all",
    collapseAllLabel: "Collapse all",
    categoriesLabel: (open, total) => `${open}/${total} open`,
    expandSidebarLabel: "Expand sidebar",
    collapseSidebarLabel: "Collapse sidebar",
    openMenuLabel: "Open navigation menu",
    closeMenuLabel: "Close navigation menu",
    dismissMenuLabel: "Dismiss navigation menu",
    signOutLabel: "Sign Out",
    signingOutLabel: "Signing out...",
    gatewayOnline: "Online",
    gatewayDegraded: "Degraded",
    gatewayOffline: "Offline",
    gatewayChecking: "Checking",
  },
  id: {
    controlPlane: "Control Plane",
    fallbackTitle: "Monitor Gateway",
    searchPlaceholder: "Cari menu...",
    quickSearchLabel: "Cari",
    clearSearchLabel: "Hapus pencarian",
    searchResultsLabel: (count) => `${count} hasil`,
    emptySearchTitle: "Menu tidak ditemukan",
    emptySearchSubtitle: "Coba kata kunci lain seperti provider, kunci, chat, atau pengaturan.",
    resetSearchLabel: "Atur ulang pencarian",
    expandAllLabel: "Buka semua",
    collapseAllLabel: "Tutup semua",
    categoriesLabel: (open, total) => `${open}/${total} terbuka`,
    expandSidebarLabel: "Perluas sidebar",
    collapseSidebarLabel: "Ringkas sidebar",
    openMenuLabel: "Buka menu navigasi",
    closeMenuLabel: "Tutup menu navigasi",
    dismissMenuLabel: "Tutup latar navigasi",
    signOutLabel: "Keluar",
    signingOutLabel: "Keluar...",
    gatewayOnline: "Aktif",
    gatewayDegraded: "Terdegradasi",
    gatewayOffline: "Offline",
    gatewayChecking: "Memeriksa",
  },
};

export function getShellCopy(locale: UILocale): ShellCopy {
  return SHELL_COPY[locale];
}
