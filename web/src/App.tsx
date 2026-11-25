import { lazy, Suspense } from "react";
import { BrowserRouter, Navigate, Outlet, Route, Routes, useLocation } from "react-router";
import { AuthGuard } from "@/features/auth/AuthGuard";
import { QueryProvider } from "@/shared/api/QueryProvider";
import { AppShell } from "@/shared/ui/AppShell";
import { Card, CardContent } from "@/shared/ui/card";
import { ErrorBoundary } from "@/shared/ui/ErrorBoundary";

const HomePage = lazy(() => import("@/pages/HomePage").then((m) => ({ default: m.HomePage })));
const InitPage = lazy(() => import("@/pages/InitPage").then((m) => ({ default: m.InitPage })));
const LoginPage = lazy(() => import("@/pages/LoginPage").then((m) => ({ default: m.LoginPage })));

const SystemOverview = lazy(() =>
  import("@/features/system/SystemOverview").then((m) => ({ default: m.SystemOverview })),
);
const ArenaConsole = lazy(() =>
  import("@/features/arena/ArenaConsole").then((m) => ({ default: m.ArenaConsole })),
);
const ChatConsole = lazy(() =>
  import("@/features/arena/ChatConsole").then((m) => ({ default: m.ChatConsole })),
);
const CombosConsole = lazy(() =>
  import("@/features/combos/CombosConsole").then((m) => ({ default: m.CombosConsole })),
);
const ImporterStudio = lazy(() =>
  import("@/features/importer/ImporterStudio").then((m) => ({ default: m.ImporterStudio })),
);
const KeysConsole = lazy(() =>
  import("@/features/keys/KeysConsole").then((m) => ({ default: m.KeysConsole })),
);
const LiveMonitorConsole = lazy(() =>
  import("@/features/monitor/LiveMonitorConsole").then((m) => ({ default: m.LiveMonitorConsole })),
);
const MediaConsole = lazy(() =>
  import("@/features/media/MediaConsole").then((m) => ({ default: m.MediaConsole })),
);
const MitmConsole = lazy(() =>
  import("@/features/mitm/MitmConsole").then((m) => ({ default: m.MitmConsole })),
);
const ModelEditorConsole = lazy(() =>
  import("@/features/modeleditor/ModelEditorConsole").then((m) => ({
    default: m.ModelEditorConsole,
  })),
);
const ModelCatalogConsole = lazy(() =>
  import("@/features/models/ModelCatalogConsole").then((m) => ({
    default: m.ModelCatalogConsole,
  })),
);
const NotificationsConsole = lazy(() =>
  import("@/features/notifications/NotificationsConsole").then((m) => ({
    default: m.NotificationsConsole,
  })),
);
const PluginsStudio = lazy(() =>
  import("@/features/plugins/PluginsStudio").then((m) => ({ default: m.PluginsStudio })),
);
const PromptGuardConsole = lazy(() =>
  import("@/features/promptinjection/PromptGuardConsole").then((m) => ({
    default: m.PromptGuardConsole,
  })),
);
const ProvidersConsole = lazy(() =>
  import("@/features/providers/ProvidersConsole").then((m) => ({ default: m.ProvidersConsole })),
);
const ProviderDetailView = lazy(() =>
  import("@/features/providers/ProviderDetailView").then((m) => ({
    default: m.ProviderDetailView,
  })),
);
const ProxyPoolsConsole = lazy(() =>
  import("@/features/proxypools/ProxyPoolsConsole").then((m) => ({
    default: m.ProxyPoolsConsole,
  })),
);
const PxpipeConsole = lazy(() =>
  import("@/features/pxpipe/PxpipeConsole").then((m) => ({ default: m.PxpipeConsole })),
);
const SecurityConsole = lazy(() =>
  import("@/features/security/SecurityConsole").then((m) => ({ default: m.SecurityConsole })),
);
const SettingsConsole = lazy(() =>
  import("@/features/settings/SettingsConsole").then((m) => ({ default: m.SettingsConsole })),
);
const TelegramConsole = lazy(() =>
  import("@/features/telegram/TelegramConsole").then((m) => ({ default: m.TelegramConsole })),
);
const TokenSaverConsole = lazy(() =>
  import("@/features/tokensaver/TokenSaverConsole").then((m) => ({
    default: m.TokenSaverConsole,
  })),
);
const TranslatorConsole = lazy(() =>
  import("@/features/translator/TranslatorConsole").then((m) => ({
    default: m.TranslatorConsole,
  })),
);
const TunnelConsole = lazy(() =>
  import("@/features/tunnel/TunnelConsole").then((m) => ({ default: m.TunnelConsole })),
);
const UsageConsole = lazy(() =>
  import("@/features/usage/UsageConsole").then((m) => ({ default: m.UsageConsole })),
);

function RouteFallback() {
  return (
    <Card>
      <CardContent className="py-8 text-sm text-[var(--text-muted)]">Loading view...</CardContent>
    </Card>
  );
}

function PageFallback() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-sm text-[var(--text-muted)]">
      Loading OMEG4MON ROUTER...
    </main>
  );
}

function DashboardLayout() {
  const { pathname } = useLocation();
  return (
    <AuthGuard>
      <AppShell>
        <ErrorBoundary resetKey={pathname} compact>
          <Suspense fallback={<RouteFallback />}>
            <Outlet />
          </Suspense>
        </ErrorBoundary>
      </AppShell>
    </AuthGuard>
  );
}

export const DASHBOARD_ROUTE_PATHS = [
  "arena",
  "chat",
  "combos",
  "importer",
  "keys",
  "media",
  "mitm",
  "monitor",
  "model-editor",
  "models",
  "notifications",
  "plugins",
  "prompt-guard",
  "providers",
  "providers/detail",
  "proxy-pools",
  "pxpipe",
  "security",
  "settings",
  "telegram",
  "token-saver",
  "translator",
  "tunnel",
  "usage",
] as const;

export function App() {
  return (
    <ErrorBoundary>
      <QueryProvider>
        <BrowserRouter>
          <Suspense fallback={<PageFallback />}>
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/init" element={<InitPage />} />
              <Route path="/login" element={<LoginPage />} />

              <Route path="/dashboard" element={<DashboardLayout />}>
                <Route index element={<SystemOverview />} />
                <Route path="arena" element={<ArenaConsole />} />
                <Route path="chat" element={<ChatConsole />} />
                <Route path="combos" element={<CombosConsole />} />
                <Route path="importer" element={<ImporterStudio />} />
                <Route path="keys" element={<KeysConsole />} />
                <Route path="media" element={<MediaConsole />} />
                <Route path="mitm" element={<MitmConsole />} />
                <Route path="monitor" element={<LiveMonitorConsole />} />
                <Route path="model-editor" element={<ModelEditorConsole />} />
                <Route path="models" element={<ModelCatalogConsole />} />
                <Route path="notifications" element={<NotificationsConsole />} />
                <Route path="plugins" element={<PluginsStudio />} />
                <Route path="prompt-guard" element={<PromptGuardConsole />} />
                <Route path="providers" element={<ProvidersConsole />} />
                <Route path="providers/detail" element={<ProviderDetailView />} />
                <Route path="proxy-pools" element={<ProxyPoolsConsole />} />
                <Route path="pxpipe" element={<PxpipeConsole />} />
                <Route path="security" element={<SecurityConsole />} />
                <Route path="settings" element={<SettingsConsole />} />
                <Route path="telegram" element={<TelegramConsole />} />
                <Route path="token-saver" element={<TokenSaverConsole />} />
                <Route path="translator" element={<TranslatorConsole />} />
                <Route path="tunnel" element={<TunnelConsole />} />
                <Route path="usage" element={<UsageConsole />} />
              </Route>

              <Route path="*" element={<Navigate to="/dashboard/" replace />} />
            </Routes>
          </Suspense>
        </BrowserRouter>
      </QueryProvider>
    </ErrorBoundary>
  );
}
