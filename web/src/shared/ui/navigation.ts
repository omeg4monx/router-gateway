import {
  Activity,
  ArrowRightLeft,
  Bell,
  Cpu,
  Database,
  Gauge,
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
  Sliders,
  SlidersHorizontal,
  Sparkles,
  Swords,
  Terminal,
  Workflow,
  Wrench,
  Zap,
} from "lucide-react";
import type { UILocale } from "@/shared/store/uiStore";

export type NavAccordionId = "features" | "system";

export interface NavItem {
  href: string;
  moduleId: string;
  icon: typeof Activity;
  labels: Record<UILocale, string>;
  keywords: readonly string[];
}

export interface NavAccordionGroup {
  id: NavAccordionId;
  icon: typeof Activity;
  labels: Record<UILocale, string>;
  defaultExpanded: boolean;
  items: readonly NavItem[];
}

export interface MainNavLayout {
  routing: readonly NavItem[];
  operations: readonly NavItem[];
  extensions: readonly NavAccordionGroup[];
  settings: NavItem | null;
}

export const ROUTING_ITEMS: readonly NavItem[] = [
  {
    href: "/dashboard/",
    moduleId: "system",
    icon: Activity,
    labels: { en: "Overview", id: "Ringkasan" },
    keywords: ["endpoint", "home", "status", "health", "monitor", "dashboard", "utama"],
  },
  {
    href: "/dashboard/providers/",
    moduleId: "providers",
    icon: Server,
    labels: { en: "Providers", id: "Provider" },
    keywords: ["openai", "anthropic", "gemini", "upstream", "oauth", "koneksi"],
  },
  {
    href: "/dashboard/combos/",
    moduleId: "combos",
    icon: Layers,
    labels: { en: "Combos", id: "Combo" },
    keywords: ["combos", "routing", "aliases", "fallback", "failover", "cooldown", "rute"],
  },
];

export const OPERATIONS_ITEMS: readonly NavItem[] = [
  {
    href: "/dashboard/keys/",
    moduleId: "keys",
    icon: KeyRound,
    labels: { en: "API Keys", id: "Kunci API" },
    keywords: ["tokens", "credentials", "quota", "rate limit", "kunci"],
  },
  {
    href: "/dashboard/usage/",
    moduleId: "usage",
    icon: Cpu,
    labels: { en: "Usage", id: "Penggunaan" },
    keywords: ["telemetry", "billing", "cost", "metrics", "tokens", "analytics", "biaya"],
  },
  {
    href: "/dashboard/quota/",
    moduleId: "providers",
    icon: Gauge,
    labels: { en: "Quota", id: "Kuota" },
    keywords: [
      "quota",
      "kuota",
      "limit",
      "sisa",
      "allowance",
      "balance",
      "reset",
      "tokenharbor",
      "provider",
    ],
  },
  {
    href: "/dashboard/monitor/",
    moduleId: "usage",
    icon: Terminal,
    labels: { en: "Live Monitor", id: "Monitor Langsung" },
    keywords: ["terminal", "realtime", "stream", "logs", "provider", "model", "aktifitas", "live"],
  },
  {
    href: "/dashboard/token-saver/",
    moduleId: "tokensaver",
    icon: Zap,
    labels: { en: "Token Saver", id: "Penghemat Token" },
    keywords: ["compress", "cache", "headroom", "hemat", "optimasi"],
  },
];

export const EXTENSION_GROUPS: readonly NavAccordionGroup[] = [
  {
    id: "features",
    icon: Sparkles,
    labels: { en: "Features", id: "Fitur" },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/arena/",
        moduleId: "arena",
        icon: Swords,
        labels: { en: "Compare Models", id: "Arena Model" },
        keywords: ["arena", "battle", "compare", "benchmark", "eval", "banding"],
      },
      {
        href: "/dashboard/chat/",
        moduleId: "arena",
        icon: MessageSquare,
        labels: { en: "Playground", id: "Playground" },
        keywords: ["playground", "prompt", "test", "messages", "chat", "obrolan"],
      },
      {
        href: "/dashboard/model-editor/",
        moduleId: "modeleditor",
        icon: Sliders,
        labels: { en: "Custom Models", id: "Editor Model" },
        keywords: ["overrides", "bypass", "temperature", "params", "tuning", "parameter"],
      },
      {
        href: "/dashboard/prompt-guard/",
        moduleId: "promptinjection",
        icon: ShieldAlert,
        labels: { en: "Prompt Guard", id: "Pelindung Prompt" },
        keywords: ["injection", "firewall", "safety", "filter", "keamanan"],
      },
      {
        href: "/dashboard/plugins/",
        moduleId: "plugins",
        icon: Puzzle,
        labels: { en: "Plugins", id: "Plugin" },
        keywords: ["extensions", "addons", "modules", "ekstensi"],
      },
      {
        href: "/dashboard/telegram/",
        moduleId: "telegram",
        icon: Send,
        labels: { en: "Telegram Bot", id: "Bot Telegram" },
        keywords: ["bot", "messaging", "webhook", "chat"],
      },
    ],
  },
  {
    id: "system",
    icon: Wrench,
    labels: { en: "System", id: "Sistem" },
    defaultExpanded: true,
    items: [
      {
        href: "/dashboard/media/",
        moduleId: "media",
        icon: Image,
        labels: { en: "Media Gateway", id: "Gateway Media" },
        keywords: ["images", "vision", "audio", "video", "web", "search", "gambar"],
      },
      {
        href: "/dashboard/translator/",
        moduleId: "translator",
        icon: ArrowRightLeft,
        labels: { en: "Translator", id: "Penerjemah" },
        keywords: ["protocol", "format", "convert", "console", "openai", "anthropic"],
      },
      {
        href: "/dashboard/proxy-pools/",
        moduleId: "proxypools",
        icon: Network,
        labels: { en: "Proxy Pools", id: "Pool Proxy" },
        keywords: ["http", "socks5", "egress", "ip", "jaringan"],
      },
      {
        href: "/dashboard/tunnel/",
        moduleId: "tunnel",
        icon: Radio,
        labels: { en: "Tunnel", id: "Tunnel" },
        keywords: ["cloudflare", "tailscale", "expose", "remote", "terowongan"],
      },
      {
        href: "/dashboard/mitm/",
        moduleId: "mitm",
        icon: Shield,
        labels: { en: "MITM Proxy", id: "Proxy MITM" },
        keywords: ["intercept", "tls", "ca", "certificate", "sertifikat"],
      },
      {
        href: "/dashboard/pxpipe/",
        moduleId: "pxpipe",
        icon: Workflow,
        labels: { en: "PXPIPE", id: "Pipeline PXPIPE" },
        keywords: ["pipeline", "middleware", "transform", "alur"],
      },
      {
        href: "/dashboard/security/",
        moduleId: "security",
        icon: Lock,
        labels: { en: "Security", id: "Keamanan" },
        keywords: ["audit", "logs", "policy", "compliance", "keamanan"],
      },
      {
        href: "/dashboard/notifications/",
        moduleId: "notifications",
        icon: Bell,
        labels: { en: "Notifications & Jobs", id: "Notifikasi & Tugas" },
        keywords: ["alerts", "cron", "background", "events", "peringatan"],
      },
      {
        href: "/dashboard/importer/",
        moduleId: "importer",
        icon: Database,
        labels: { en: "Legacy Importer", id: "Impor Data Lama" },
        keywords: ["migrate", "import", "sqlite", "migrasi"],
      },
    ],
  },
];

export const SETTINGS_NAV_ITEM: NavItem = {
  href: "/dashboard/settings/",
  moduleId: "settings",
  icon: SlidersHorizontal,
  labels: { en: "Settings", id: "Pengaturan" },
  keywords: ["config", "preferences", "password", "locale", "theme", "konfigurasi"],
};

export const ALL_NAV_ITEMS: readonly NavItem[] = [
  ...ROUTING_ITEMS,
  ...OPERATIONS_ITEMS,
  ...EXTENSION_GROUPS.flatMap((group) => group.items),
  SETTINGS_NAV_ITEM,
];

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
  sectionLabel: Record<UILocale, string>;
  accordionId?: NavAccordionId;
  item: NavItem;
}

export function findActiveNavContext(pathname: string): ActiveNavContext | null {
  const normalized = normalizeRoute(pathname);

  for (const item of ROUTING_ITEMS) {
    if (isNavItemActive(normalized, item.href)) {
      return { sectionLabel: { en: "Routing", id: "Routing" }, item };
    }
  }
  for (const item of OPERATIONS_ITEMS) {
    if (isNavItemActive(normalized, item.href)) {
      return { sectionLabel: { en: "Operations", id: "Operasional" }, item };
    }
  }
  for (const group of EXTENSION_GROUPS) {
    for (const item of group.items) {
      if (isNavItemActive(normalized, item.href)) {
        return { sectionLabel: group.labels, accordionId: group.id, item };
      }
    }
  }
  if (isNavItemActive(normalized, SETTINGS_NAV_ITEM.href)) {
    return { sectionLabel: { en: "System", id: "Sistem" }, item: SETTINGS_NAV_ITEM };
  }
  return null;
}

function itemMatchesQuery(item: NavItem, terms: readonly string[], extraHaystack = ""): boolean {
  if (terms.length === 0) {
    return true;
  }
  const haystack =
    `${item.labels.en} ${item.labels.id} ${item.href} ${item.keywords.join(" ")} ${extraHaystack}`.toLowerCase();
  return terms.every((term) => haystack.includes(term));
}

function isModuleRegistered(item: NavItem, registeredModules: ReadonlySet<string> | null): boolean {
  if (!registeredModules) {
    return true;
  }
  return registeredModules.has(item.moduleId);
}

export function resolveVisibleNavigation(options: {
  query: string;
  registeredModules?: readonly string[] | null;
}): MainNavLayout {
  const terms = options.query.trim().toLowerCase().split(/\s+/).filter(Boolean);

  const moduleSet =
    options.registeredModules && options.registeredModules.length > 0
      ? new Set(options.registeredModules)
      : null;

  const filterList = (items: readonly NavItem[], extra = "") =>
    items.filter(
      (item) => isModuleRegistered(item, moduleSet) && itemMatchesQuery(item, terms, extra),
    );

  const routing = filterList(ROUTING_ITEMS, "routing");
  const operations = filterList(OPERATIONS_ITEMS, "operations operasional");
  const extensions: NavAccordionGroup[] = [];

  for (const group of EXTENSION_GROUPS) {
    const matching = filterList(group.items, `${group.labels.en} ${group.labels.id} extensions`);
    if (matching.length > 0) {
      extensions.push({ ...group, items: matching });
    }
  }

  const settingsVisible =
    isModuleRegistered(SETTINGS_NAV_ITEM, moduleSet) &&
    itemMatchesQuery(SETTINGS_NAV_ITEM, terms, "system sistem settings pengaturan");

  return {
    routing,
    operations,
    extensions,
    settings: settingsVisible ? SETTINGS_NAV_ITEM : null,
  };
}

export function countVisibleNavItems(layout: MainNavLayout): number {
  return (
    layout.routing.length +
    layout.operations.length +
    layout.extensions.reduce((sum, group) => sum + group.items.length, 0) +
    (layout.settings ? 1 : 0)
  );
}

export interface ShellCopy {
  controlPlane: string;
  fallbackTitle: string;
  sectionRouting: string;
  sectionOperations: string;
  sectionExtensions: string;
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
  signOutConfirmTitle: string;
  signOutConfirmMessage: string;
  signOutCancelLabel: string;
  signOutConfirmLabel: string;
  gatewayOnline: string;
  gatewayDegraded: string;
  gatewayOffline: string;
  gatewayChecking: string;
}

const SHELL_COPY: Record<UILocale, ShellCopy> = {
  en: {
    controlPlane: "Control Plane",
    fallbackTitle: "Gateway Monitor",
    sectionRouting: "Routing",
    sectionOperations: "Operations",
    sectionExtensions: "Extensions",
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
    signOutConfirmTitle: "Sign out?",
    signOutConfirmMessage:
      "You will need to enter your administrator password to access the dashboard again.",
    signOutCancelLabel: "Cancel",
    signOutConfirmLabel: "Sign Out",
    gatewayOnline: "Online",
    gatewayDegraded: "Degraded",
    gatewayOffline: "Offline",
    gatewayChecking: "Checking",
  },
  id: {
    controlPlane: "Control Plane",
    fallbackTitle: "Monitor Gateway",
    sectionRouting: "Routing",
    sectionOperations: "Operasional",
    sectionExtensions: "Ekstensi",
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
    signOutConfirmTitle: "Keluar dari sesi?",
    signOutConfirmMessage:
      "Anda perlu memasukkan kata sandi administrator lagi untuk mengakses dashboard.",
    signOutCancelLabel: "Batal",
    signOutConfirmLabel: "Ya, Keluar",
    gatewayOnline: "Aktif",
    gatewayDegraded: "Terdegradasi",
    gatewayOffline: "Offline",
    gatewayChecking: "Memeriksa",
  },
};

export function getShellCopy(locale: UILocale): ShellCopy {
  return SHELL_COPY[locale];
}
