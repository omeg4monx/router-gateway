import type { UILocale } from "@/shared/store/uiStore";
import type { AuthMeResponse } from "./hooks";

export interface LoginCopy {
  brandEyebrow: string;
  heroTitle: string;
  heroDescription: string;
  gatewayStatusLabel: string;
  gatewayOnline: string;
  gatewayDegraded: string;
  gatewayOffline: string;
  gatewayChecking: string;
  buildVersionLabel: string;
  buildVersionLoading: string;
  buildVersionUnavailable: string;
  runtimeLabel: string;
  databaseLabel: string;
  databaseOk: string;
  databaseError: string;
  databaseUnknown: string;
  uptimeLabel: string;
  authPostureLabel: string;
  securityNoteTitle: string;
  securityNoteBody: string;
  footerContract: string;
  cardEyebrow: string;
  signInHeading: string;
  signInSubtitle: string;
  pendingTitle: string;
  pendingSubtitle: string;
  unreachableTitle: string;
  unreachableHint: string;
  unreachableTips: readonly string[];
  retryLabel: string;
  retryingLabel: string;
  openDashboardLabel: string;
  uninitializedTitle: string;
  uninitializedBody: string;
  uninitializedCta: string;
  redirectAuthenticatedTitle: string;
  redirectLoopbackTitle: string;
  redirectSubtitle: string;
  redirectCta: string;
  remoteEnforcedTitle: string;
  remoteEnforcedBody: string;
  passwordLabel: string;
  passwordPlaceholder: string;
  showPasswordLabel: string;
  hidePasswordLabel: string;
  capsLockWarning: string;
  submitLabel: string;
  submittingLabel: string;
  enterHint: string;
  helpTitle: string;
  helpBody: string;
  helpRecovery: string;
  errorTitle: string;
  technicalDetailLabel: string;
}

const LOGIN_COPY: Record<UILocale, LoginCopy> = {
  en: {
    brandEyebrow: "Self-Hosted LLM Gateway",
    heroTitle: "Control Plane Access",
    heroDescription:
      "Unified OpenAI and Anthropic compatible routing kernel. Single binary runtime with embedded SQLite state and deterministic provider failover.",
    gatewayStatusLabel: "Gateway",
    gatewayOnline: "Online",
    gatewayDegraded: "Degraded",
    gatewayOffline: "Offline",
    gatewayChecking: "Checking",
    buildVersionLabel: "Build",
    buildVersionLoading: "Loading...",
    buildVersionUnavailable: "Unavailable",
    runtimeLabel: "Runtime Health",
    databaseLabel: "SQLite",
    databaseOk: "Healthy",
    databaseError: "Error",
    databaseUnknown: "Unknown",
    uptimeLabel: "Uptime",
    authPostureLabel: "Access Policy",
    securityNoteTitle: "Signed HttpOnly Session",
    securityNoteBody:
      "Credentials are verified locally against your gateway database and issued as a scoped, signed session cookie.",
    footerContract: "API CONTRACT · OPENAPI 3.1",
    cardEyebrow: "Operator Authentication",
    signInHeading: "Sign in to Gateway",
    signInSubtitle: "Enter your administrator password to manage routes, providers, and API keys.",
    pendingTitle: "Checking gateway session",
    pendingSubtitle: "Verifying current authentication policy and control plane readiness...",
    unreachableTitle: "Unable to reach control plane",
    unreachableHint: "The dashboard cannot connect to the gateway API right now.",
    unreachableTips: [
      "Confirm the OMEG4MON Router process is running on this host.",
      "Verify your reverse proxy or port binding forwards /api/* requests.",
    ],
    retryLabel: "Retry Connection",
    retryingLabel: "Retrying...",
    openDashboardLabel: "Open Dashboard",
    uninitializedTitle: "First-run setup required",
    uninitializedBody:
      "This gateway has not been initialized yet. Create your administrator password before signing in.",
    uninitializedCta: "Complete First-Run Setup",
    redirectAuthenticatedTitle: "Operator session active",
    redirectLoopbackTitle: "Local loopback access active",
    redirectSubtitle: "Redirecting you to the control plane dashboard...",
    redirectCta: "Continue to Dashboard",
    remoteEnforcedTitle: "Remote Bind Enforcement Active",
    remoteEnforcedBody:
      "Non-loopback control plane access requires a signed operator session cookie.",
    passwordLabel: "Administrator Password",
    passwordPlaceholder: "Enter administrator password",
    showPasswordLabel: "Show password",
    hidePasswordLabel: "Hide password",
    capsLockWarning: "Caps Lock is on",
    submitLabel: "Sign In",
    submittingLabel: "Signing in...",
    enterHint: "Press Enter to sign in",
    helpTitle: "Need help signing in?",
    helpBody:
      "Use the administrator password configured during first-run setup (/init) or updated in Settings.",
    helpRecovery:
      "If you lost access, update or reset the admin credential directly on the host machine running the router.",
    errorTitle: "Could not sign in",
    technicalDetailLabel: "Detail",
  },
  id: {
    brandEyebrow: "Gateway LLM Self-Hosted",
    heroTitle: "Akses Control Plane",
    heroDescription:
      "Kernel routing terpadu yang kompatibel dengan OpenAI dan Anthropic. Runtime single-binary dengan SQLite tertanam dan failover provider deterministik.",
    gatewayStatusLabel: "Gateway",
    gatewayOnline: "Aktif",
    gatewayDegraded: "Terdegradasi",
    gatewayOffline: "Offline",
    gatewayChecking: "Memeriksa",
    buildVersionLabel: "Versi",
    buildVersionLoading: "Memuat...",
    buildVersionUnavailable: "Tidak tersedia",
    runtimeLabel: "Kesehatan Runtime",
    databaseLabel: "SQLite",
    databaseOk: "Normal",
    databaseError: "Bermasalah",
    databaseUnknown: "Tidak diketahui",
    uptimeLabel: "Waktu Aktif",
    authPostureLabel: "Kebijakan Akses",
    securityNoteTitle: "Sesi HttpOnly Bertanda Tangan",
    securityNoteBody:
      "Kredensial diverifikasi secara lokal pada database gateway dan diterbitkan sebagai cookie sesi terenkripsi.",
    footerContract: "KONTRAK API · OPENAPI 3.1",
    cardEyebrow: "Autentikasi Operator",
    signInHeading: "Masuk ke Gateway",
    signInSubtitle:
      "Masukkan kata sandi administrator untuk mengelola rute, provider, dan kunci API.",
    pendingTitle: "Memeriksa sesi gateway",
    pendingSubtitle: "Memverifikasi kebijakan autentikasi dan kesiapan control plane...",
    unreachableTitle: "Tidak dapat terhubung ke control plane",
    unreachableHint: "Dashboard belum dapat menjangkau API gateway saat ini.",
    unreachableTips: [
      "Pastikan proses OMEG4MON Router sedang berjalan di host ini.",
      "Periksa reverse proxy atau port binding agar meneruskan permintaan /api/*.",
    ],
    retryLabel: "Coba Hubungkan Ulang",
    retryingLabel: "Menghubungkan ulang...",
    openDashboardLabel: "Buka Dashboard",
    uninitializedTitle: "Pengaturan awal diperlukan",
    uninitializedBody:
      "Gateway ini belum diinisialisasi. Buat kata sandi administrator terlebih dahulu sebelum masuk.",
    uninitializedCta: "Lanjutkan Pengaturan Awal",
    redirectAuthenticatedTitle: "Sesi operator aktif",
    redirectLoopbackTitle: "Akses loopback lokal aktif",
    redirectSubtitle: "Mengarahkan Anda ke dashboard control plane...",
    redirectCta: "Lanjut ke Dashboard",
    remoteEnforcedTitle: "Penegakan Akses Jarak Jauh Aktif",
    remoteEnforcedBody:
      "Akses control plane di luar loopback mewajibkan cookie sesi operator yang sah.",
    passwordLabel: "Kata Sandi Administrator",
    passwordPlaceholder: "Masukkan kata sandi administrator",
    showPasswordLabel: "Tampilkan kata sandi",
    hidePasswordLabel: "Sembunyikan kata sandi",
    capsLockWarning: "Caps Lock aktif",
    submitLabel: "Masuk",
    submittingLabel: "Memproses...",
    enterHint: "Tekan Enter untuk masuk",
    helpTitle: "Butuh bantuan untuk masuk?",
    helpBody:
      "Gunakan kata sandi administrator yang dibuat saat pengaturan awal (/init) atau diperbarui di menu Pengaturan.",
    helpRecovery:
      "Jika lupa kata sandi, perbarui atau reset kredensial admin langsung dari mesin host yang menjalankan router.",
    errorTitle: "Gagal masuk",
    technicalDetailLabel: "Detail",
  },
};

export function getLoginCopy(locale: UILocale): LoginCopy {
  return LOGIN_COPY[locale];
}

export interface FormattedAuthError {
  title: string;
  message: string;
  detail?: string;
}

export function formatAuthError(
  rawMessage: string | null | undefined,
  locale: UILocale,
): FormattedAuthError | null {
  const trimmed = rawMessage?.trim();
  if (!trimmed) {
    return null;
  }

  const copy = getLoginCopy(locale);
  const normalized = trimmed.toLowerCase();

  if (normalized === "administrator password is required.") {
    return {
      title: copy.errorTitle,
      message:
        locale === "id"
          ? "Kata sandi administrator wajib diisi."
          : "Administrator password is required.",
    };
  }

  if (normalized.includes("invalid admin password") || normalized.includes("invalid_credentials")) {
    return {
      title: copy.errorTitle,
      message:
        locale === "id"
          ? "Kata sandi administrator tidak cocok. Periksa kembali penulisan huruf besar/kecil lalu coba lagi."
          : "The administrator password does not match. Check capitalization and try again.",
      detail: trimmed,
    };
  }

  if (normalized.includes("not initialized") || normalized.includes("not_initialized")) {
    return {
      title: copy.errorTitle,
      message:
        locale === "id"
          ? "Gateway belum diinisialisasi. Selesaikan pengaturan awal terlebih dahulu."
          : "This gateway is not initialized yet. Complete first-run setup first.",
      detail: trimmed,
    };
  }

  if (
    normalized.includes("failed to fetch") ||
    normalized.includes("networkerror") ||
    normalized.includes("network error")
  ) {
    return {
      title: copy.errorTitle,
      message:
        locale === "id"
          ? "Koneksi ke gateway terputus. Pastikan layanan masih aktif lalu coba lagi."
          : "Connection to the gateway was interrupted. Confirm the service is running and retry.",
      detail: trimmed,
    };
  }

  return {
    title: copy.errorTitle,
    message: trimmed,
  };
}

export interface AuthPostureDescriptor {
  code: string;
  badgeVariant: "default" | "secondary" | "success" | "warning" | "destructive";
  summary: string;
}

export function describeAuthPosture(
  input: {
    isPending: boolean;
    isError: boolean;
    session?: AuthMeResponse;
  },
  locale: UILocale,
): AuthPostureDescriptor {
  if (input.isPending) {
    return {
      code: locale === "id" ? "MEMERIKSA" : "CHECKING",
      badgeVariant: "secondary",
      summary:
        locale === "id"
          ? "Memeriksa kebijakan autentikasi control plane..."
          : "Inspecting control plane authentication policy...",
    };
  }

  if (input.isError || !input.session) {
    return {
      code: locale === "id" ? "TERPUTUS" : "UNAVAILABLE",
      badgeVariant: "destructive",
      summary:
        locale === "id"
          ? "Kebijakan akses tidak dapat dibaca karena API tidak terjangkau."
          : "Access policy could not be read while the API is unreachable.",
    };
  }

  if (!input.session.initialized) {
    return {
      code: locale === "id" ? "BELUM INIT" : "UNINITIALIZED",
      badgeVariant: "warning",
      summary:
        locale === "id"
          ? "Kata sandi admin belum dibuat. Lanjutkan ke pengaturan awal."
          : "No admin password configured yet. Complete initial setup first.",
    };
  }

  if (input.session.remote_enforced) {
    return {
      code: "REMOTE ENFORCED",
      badgeVariant: "default",
      summary:
        locale === "id"
          ? "Bind non-loopback aktif. Login wajib untuk semua akses operator."
          : "Non-loopback bind active. Signed operator login is strictly required.",
    };
  }

  if (input.session.require_login) {
    return {
      code: locale === "id" ? "LOGIN WAJIB" : "ENFORCED",
      badgeVariant: "default",
      summary:
        locale === "id"
          ? "Autentikasi operator diaktifkan oleh kebijakan gateway."
          : "Operator sign-in is required by the active gateway policy.",
    };
  }

  return {
    code: "LOOPBACK OPEN",
    badgeVariant: "success",
    summary:
      locale === "id"
        ? "Akses lokal (loopback) dibuka tanpa login."
        : "Local loopback access is open without requiring sign-in.",
  };
}

export function formatUptime(seconds: number | null | undefined, locale: UILocale): string {
  if (seconds === null || seconds === undefined || !Number.isFinite(seconds)) {
    return "—";
  }

  const wholeSeconds = Math.max(0, Math.floor(seconds));
  if (wholeSeconds < 60) {
    return locale === "id" ? "< 1 mnt" : "< 1m";
  }

  const totalMinutes = Math.floor(wholeSeconds / 60);
  const days = Math.floor(totalMinutes / (60 * 24));
  const hours = Math.floor((totalMinutes % (60 * 24)) / 60);
  const minutes = totalMinutes % 60;

  if (days > 0) {
    return locale === "id" ? `${days}h ${hours}j` : `${days}d ${hours}h`;
  }
  if (hours > 0) {
    return locale === "id" ? `${hours}j ${minutes}m` : `${hours}h ${minutes}m`;
  }
  return locale === "id" ? `${minutes} mnt` : `${minutes}m`;
}
