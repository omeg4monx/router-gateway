import type { SystemUpdateStatusResponse } from "@/shared/api/client";
import type { UILocale } from "@/shared/store/uiStore";

export type SystemUpdateState = SystemUpdateStatusResponse["state"];

export interface SystemUpdateCopy {
  buttonLabel: string;
  buttonAriaLabel: string;
  updatingLabel: string;
  modalTitle: string;
  modalDescription: string;
  branchLabel: string;
  currentCommitLabel: string;
  latestCommitLabel: string;
  modeLabel: string;
  modeBinary: string;
  modeSource: string;
  modeNone: string;
  latestMessageLabel: string;
  confirmAction: string;
  retryAction: string;
  recheckAction: string;
  recheckingLabel: string;
  cancelAction: string;
  closeAction: string;
  restartingNotice: string;
  restartTimeoutNotice: string;
  reloadAction: string;
  viewCommitLabel: string;
}

const UPDATE_COPY: Record<UILocale, SystemUpdateCopy> = {
  en: {
    buttonLabel: "Update",
    buttonAriaLabel: "Update router from main branch",
    updatingLabel: "Updating…",
    modalTitle: "Update Available",
    modalDescription:
      "A newer commit is available on the main branch. Updating will replace the running binary and restart the router automatically.",
    branchLabel: "Branch",
    currentCommitLabel: "Current",
    latestCommitLabel: "Latest",
    modeLabel: "Mode",
    modeBinary: "Standalone Binary",
    modeSource: "Local Git Rebuild",
    modeNone: "Unavailable",
    latestMessageLabel: "Latest Commit",
    confirmAction: "Update & Restart",
    retryAction: "Retry Update",
    recheckAction: "Check Again",
    recheckingLabel: "Checking…",
    cancelAction: "Cancel",
    closeAction: "Close",
    restartingNotice:
      "Waiting for the updated router to come back online. This page will reload automatically.",
    restartTimeoutNotice:
      "The router is taking longer than expected to come back online. Check the service logs or reload manually.",
    reloadAction: "Reload Page",
    viewCommitLabel: "View commit",
  },
  id: {
    buttonLabel: "Perbarui",
    buttonAriaLabel: "Perbarui router dari branch main",
    updatingLabel: "Memperbarui…",
    modalTitle: "Pembaruan Tersedia",
    modalDescription:
      "Ada commit terbaru di branch main. Proses pembaruan akan mengganti binary aktif dan memulai ulang router secara otomatis.",
    branchLabel: "Branch",
    currentCommitLabel: "Saat Ini",
    latestCommitLabel: "Terbaru",
    modeLabel: "Mode",
    modeBinary: "Binary Standalone",
    modeSource: "Rebuild Git Lokal",
    modeNone: "Tidak Tersedia",
    latestMessageLabel: "Commit Terbaru",
    confirmAction: "Perbarui & Restart",
    retryAction: "Coba Lagi",
    recheckAction: "Cek Ulang",
    recheckingLabel: "Mengecek…",
    cancelAction: "Batal",
    closeAction: "Tutup",
    restartingNotice:
      "Menunggu router versi baru kembali aktif. Halaman ini akan dimuat ulang otomatis.",
    restartTimeoutNotice:
      "Router membutuhkan waktu lebih lama dari biasanya untuk kembali aktif. Periksa log layanan atau muat ulang secara manual.",
    reloadAction: "Muat Ulang",
    viewCommitLabel: "Lihat commit",
  },
};

const STATE_LABELS: Record<UILocale, Record<SystemUpdateState, string>> = {
  en: {
    idle: "Ready",
    checking: "Checking repository and release metadata…",
    downloading: "Downloading update package…",
    building: "Building web UI and Go binary…",
    verifying: "Verifying SHA256 checksum and executable…",
    applying: "Replacing executable atomically…",
    restarting: "Restarting router…",
    failed: "Update failed",
  },
  id: {
    idle: "Siap",
    checking: "Memeriksa metadata repositori dan rilis…",
    downloading: "Mengunduh paket pembaruan…",
    building: "Membangun web UI dan binary Go…",
    verifying: "Memverifikasi checksum SHA256 dan binary…",
    applying: "Mengganti binary secara atomik…",
    restarting: "Memulai ulang router…",
    failed: "Pembaruan gagal",
  },
};

export function getUpdateCopy(locale: UILocale): SystemUpdateCopy {
  return UPDATE_COPY[locale];
}

export function isUpdateInProgress(state: SystemUpdateState | undefined): boolean {
  return (
    state === "checking" ||
    state === "downloading" ||
    state === "building" ||
    state === "verifying" ||
    state === "applying" ||
    state === "restarting"
  );
}

export function shouldShowUpdateButton(status: SystemUpdateStatusResponse | undefined): boolean {
  if (!status) {
    return false;
  }
  return status.updateAvailable || isUpdateInProgress(status.state);
}

export function getUpdateStateLabel(locale: UILocale, state: SystemUpdateState): string {
  return STATE_LABELS[locale][state];
}

export function getUnavailableReasonMessage(locale: UILocale, reason: string): string {
  if (locale === "id") {
    switch (reason) {
      case "dirty_worktree":
        return "Repositori lokal memiliki perubahan yang belum di-commit. Commit atau stash perubahan Anda terlebih dahulu sebelum update otomatis.";
      case "wrong_branch":
        return "Repositori lokal saat ini tidak berada di branch main. Pindah ke branch main terlebih dahulu sebelum menjalankan update otomatis.";
      case "release_building":
        return "Binary rilis untuk commit main terbaru masih dalam proses build di GitHub Actions. Coba cek ulang dalam beberapa menit.";
      case "not_writable":
        return "Direktori executable saat ini tidak memiliki izin tulis untuk mengganti binary.";
      case "release_unavailable":
        return "Metadata rilis GitHub tidak dapat diambil saat ini.";
      default:
        return "Pembaruan otomatis tidak dapat diterapkan pada lingkungan ini.";
    }
  }

  switch (reason) {
    case "dirty_worktree":
      return "The local git repository has uncommitted changes. Commit or stash your changes before running auto-update.";
    case "wrong_branch":
      return "The local git repository is not currently checked out on the main branch. Switch to main before running auto-update.";
    case "release_building":
      return "Prebuilt binaries for the latest main commit are still building in GitHub Actions. Check again in a few minutes.";
    case "not_writable":
      return "The current executable directory is not writable for an in-place binary swap.";
    case "release_unavailable":
      return "GitHub release metadata could not be reached right now.";
    default:
      return "Automatic update cannot be applied in the current environment.";
  }
}
