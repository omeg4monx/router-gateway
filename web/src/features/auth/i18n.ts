import type { UILocale } from "@/shared/store/uiStore";

export interface LoginCopy {
  signInHeading: string;
  pendingLabel: string;
  unreachableTitle: string;
  retryLabel: string;
  retryingLabel: string;
  uninitializedTitle: string;
  uninitializedBody: string;
  uninitializedCta: string;
  redirectTitle: string;
  redirectCta: string;
  passwordLabel: string;
  passwordPlaceholder: string;
  showPasswordLabel: string;
  hidePasswordLabel: string;
  capsLockWarning: string;
  submitLabel: string;
  submittingLabel: string;
  errorTitle: string;
  technicalDetailLabel: string;
}

const LOGIN_COPY: Record<UILocale, LoginCopy> = {
  en: {
    signInHeading: "Sign in",
    pendingLabel: "Checking session...",
    unreachableTitle: "Unable to reach gateway",
    retryLabel: "Retry",
    retryingLabel: "Retrying...",
    uninitializedTitle: "Setup required",
    uninitializedBody: "Create an administrator password before signing in.",
    uninitializedCta: "Go to Setup",
    redirectTitle: "Redirecting to dashboard...",
    redirectCta: "Open Dashboard",
    passwordLabel: "Password",
    passwordPlaceholder: "Enter administrator password",
    showPasswordLabel: "Show password",
    hidePasswordLabel: "Hide password",
    capsLockWarning: "Caps Lock is on",
    submitLabel: "Sign In",
    submittingLabel: "Signing in...",
    errorTitle: "Sign in failed",
    technicalDetailLabel: "Detail",
  },
  id: {
    signInHeading: "Masuk",
    pendingLabel: "Memeriksa sesi...",
    unreachableTitle: "Tidak dapat terhubung ke gateway",
    retryLabel: "Coba Lagi",
    retryingLabel: "Menghubungkan ulang...",
    uninitializedTitle: "Pengaturan awal diperlukan",
    uninitializedBody: "Buat kata sandi administrator terlebih dahulu sebelum masuk.",
    uninitializedCta: "Buka Pengaturan Awal",
    redirectTitle: "Mengarahkan ke dashboard...",
    redirectCta: "Buka Dashboard",
    passwordLabel: "Kata Sandi",
    passwordPlaceholder: "Masukkan kata sandi administrator",
    showPasswordLabel: "Tampilkan kata sandi",
    hidePasswordLabel: "Sembunyikan kata sandi",
    capsLockWarning: "Caps Lock aktif",
    submitLabel: "Masuk",
    submittingLabel: "Memproses...",
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
