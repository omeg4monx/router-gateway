import { useCallback, useEffect, useRef, useState } from "react";
import type { UILocale } from "@/shared/store/uiStore";

interface BeforeInstallPromptEvent extends Event {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed"; platform: string }>;
}

export interface PwaCopy {
  installLabel: string;
  installingLabel: string;
  installAriaLabel: string;
  offlineTitle: string;
  offlineDescription: string;
  updateTitle: string;
  updateDescription: string;
  updateAction: string;
  dismissAction: string;
}

export function getPwaCopy(locale: UILocale): PwaCopy {
  if (locale === "id") {
    return {
      installLabel: "Pasang",
      installingLabel: "Memasang...",
      installAriaLabel: "Pasang aplikasi OMEG4MON ROUTER",
      offlineTitle: "Sedang offline",
      offlineDescription:
        "Koneksi ke server router terputus. Data langsung dan panggilan API akan kembali aktif saat tersambung.",
      updateTitle: "Versi baru tersedia",
      updateDescription: "Muat ulang untuk memakai pembaruan antarmuka terbaru.",
      updateAction: "Muat ulang",
      dismissAction: "Nanti",
    };
  }
  return {
    installLabel: "Install",
    installingLabel: "Installing...",
    installAriaLabel: "Install OMEG4MON ROUTER app",
    offlineTitle: "You are offline",
    offlineDescription:
      "Connection to the router server is unavailable. Live telemetry and API calls will resume once reconnected.",
    updateTitle: "Update available",
    updateDescription: "Reload to switch to the latest interface version.",
    updateAction: "Reload",
    dismissAction: "Later",
  };
}

export function isStandaloneDisplayMode(): boolean {
  if (typeof window === "undefined") {
    return false;
  }
  const mediaStandalone =
    typeof window.matchMedia === "function" &&
    window.matchMedia("(display-mode: standalone)").matches;
  const navigatorStandalone =
    "standalone" in window.navigator &&
    Boolean((window.navigator as Navigator & { standalone?: boolean }).standalone);
  return mediaStandalone || navigatorStandalone;
}

export interface PwaStatus {
  isOffline: boolean;
  updateAvailable: boolean;
  canInstall: boolean;
  isInstalling: boolean;
  applyUpdate: () => void;
  dismissUpdate: () => void;
  promptInstall: () => Promise<void>;
}

export function usePwaStatus(): PwaStatus {
  const [isOffline, setIsOffline] = useState<boolean>(() =>
    typeof navigator !== "undefined" ? !navigator.onLine : false,
  );
  const [updateAvailable, setUpdateAvailable] = useState(false);
  const [deferredPrompt, setDeferredPrompt] = useState<BeforeInstallPromptEvent | null>(null);
  const [isInstalling, setIsInstalling] = useState(false);

  const waitingWorkerRef = useRef<ServiceWorker | null>(null);
  const reloadingRef = useRef(false);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    const handleOnline = () => setIsOffline(false);
    const handleOffline = () => setIsOffline(true);
    const handleBeforeInstallPrompt = (event: Event) => {
      event.preventDefault();
      if (!isStandaloneDisplayMode()) {
        setDeferredPrompt(event as BeforeInstallPromptEvent);
      }
    };
    const handleAppInstalled = () => {
      setDeferredPrompt(null);
      setIsInstalling(false);
    };

    window.addEventListener("online", handleOnline);
    window.addEventListener("offline", handleOffline);
    window.addEventListener("beforeinstallprompt", handleBeforeInstallPrompt);
    window.addEventListener("appinstalled", handleAppInstalled);

    return () => {
      window.removeEventListener("online", handleOnline);
      window.removeEventListener("offline", handleOffline);
      window.removeEventListener("beforeinstallprompt", handleBeforeInstallPrompt);
      window.removeEventListener("appinstalled", handleAppInstalled);
    };
  }, []);

  useEffect(() => {
    if (
      typeof window === "undefined" ||
      !import.meta.env.PROD ||
      !window.isSecureContext ||
      !("serviceWorker" in navigator)
    ) {
      return;
    }

    let cancelled = false;
    let activeRegistration: ServiceWorkerRegistration | null = null;
    let trackedInstalling: ServiceWorker | null = null;

    const markWaiting = (worker: ServiceWorker | null) => {
      if (!worker || cancelled) {
        return;
      }
      waitingWorkerRef.current = worker;
      setUpdateAvailable(true);
    };

    const handleInstallingStateChange = () => {
      if (
        trackedInstalling &&
        trackedInstalling.state === "installed" &&
        navigator.serviceWorker.controller
      ) {
        markWaiting(trackedInstalling);
      }
    };

    const handleUpdateFound = () => {
      if (!activeRegistration) {
        return;
      }
      if (trackedInstalling) {
        trackedInstalling.removeEventListener("statechange", handleInstallingStateChange);
      }
      trackedInstalling = activeRegistration.installing;
      if (trackedInstalling) {
        trackedInstalling.addEventListener("statechange", handleInstallingStateChange);
      }
    };

    const handleControllerChange = () => {
      if (reloadingRef.current) {
        return;
      }
      reloadingRef.current = true;
      window.location.reload();
    };

    navigator.serviceWorker.addEventListener("controllerchange", handleControllerChange);

    navigator.serviceWorker
      .register("/sw.js", { scope: "/" })
      .then((registration) => {
        if (cancelled) {
          return;
        }
        activeRegistration = registration;
        if (registration.waiting && navigator.serviceWorker.controller) {
          markWaiting(registration.waiting);
        }
        registration.addEventListener("updatefound", handleUpdateFound);
      })
      .catch(() => {
        // Non-fatal: app remains fully functional without a Service Worker.
      });

    return () => {
      cancelled = true;
      navigator.serviceWorker.removeEventListener("controllerchange", handleControllerChange);
      if (activeRegistration) {
        activeRegistration.removeEventListener("updatefound", handleUpdateFound);
      }
      if (trackedInstalling) {
        trackedInstalling.removeEventListener("statechange", handleInstallingStateChange);
      }
    };
  }, []);

  const applyUpdate = useCallback(() => {
    const waiting = waitingWorkerRef.current;
    if (!waiting) {
      setUpdateAvailable(false);
      return;
    }
    waiting.postMessage({ type: "SKIP_WAITING" });
  }, []);

  const dismissUpdate = useCallback(() => {
    setUpdateAvailable(false);
  }, []);

  const promptInstall = useCallback(async () => {
    if (!deferredPrompt || isInstalling) {
      return;
    }
    setIsInstalling(true);
    try {
      await deferredPrompt.prompt();
      const choice = await deferredPrompt.userChoice;
      if (choice.outcome === "accepted") {
        setDeferredPrompt(null);
      }
    } catch {
      // Ignore user cancellation or browser prompt errors.
    } finally {
      setIsInstalling(false);
    }
  }, [deferredPrompt, isInstalling]);

  return {
    isOffline,
    updateAvailable,
    canInstall: deferredPrompt !== null,
    isInstalling,
    applyUpdate,
    dismissUpdate,
    promptInstall,
  };
}
