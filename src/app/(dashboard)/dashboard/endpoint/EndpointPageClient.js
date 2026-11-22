"use client";

// One auto-provision per page load: two concurrent loads (StrictMode remount)
// both seeing "no keys" would each create one, leaving a duplicate behind.
let defaultKeyProvisionInFlight = false;

import { useState, useEffect, useRef, useCallback, useMemo } from "react";
import PropTypes from "prop-types";
import { Button, Input, Modal, ConfirmModal, ModelSelectModal } from "@/shared/components";
import { useCopyToClipboard } from "@/shared/hooks/useCopyToClipboard";
import { cn } from "@/shared/utils/cn";
import useSettingsStore from "@/store/settingsStore";
import { useNotificationStore } from "@/store/notificationStore";
import {
  TUNNEL_BENEFITS,
  TUNNEL_PING_INTERVAL_MS,
  TUNNEL_PING_MAX_MS,
  STATUS_POLL_FAST_MS,
  REACHABLE_MISS_THRESHOLD,
  CLIENT_PING_FAST_MS,
} from "./endpointConstants";
import { clientPingUrl, clientPingAny } from "./endpointPing";
import {
  normalizeOrigin,
  filterApiKeys,
  countKeyBuckets,
  resolveDefaultSnippetModel,
} from "./endpointUtils";
import { generateEndpointSnippet } from "./endpointSnippets";
import GatewayConsole from "./components/GatewayConsole";
import ApiKeyCard from "./components/ApiKeyCard";
import KeyFormModal, { EMPTY_PERMISSIONS } from "./components/KeyFormModal";
import KeySnippetModal from "./components/KeySnippetModal";
import BulkEditKeys from "./components/BulkEditKeys";
import KeyShareModal from "./components/KeyShareModal";
import StatusAlert from "./components/StatusAlert";

const FILTER_TABS = [
  { id: "all", label: "All" },
  { id: "active", label: "Active" },
  { id: "disabled", label: "Disabled" },
  { id: "exhausted", label: "Exhausted" },
  { id: "expiring", label: "Expiring" },
];

function parseAllowedModelsList(str) {
  if (!str || str.trim() === "*" || str.trim() === "") return [];
  return str
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
}

export default function APIPageClient({ machineId }) {
  const [keys, setKeys] = useState([]);
  const [loading, setLoading] = useState(true);
  const [keySearch, setKeySearch] = useState("");
  const [statusFilter, setStatusFilter] = useState("all");
  const [activeRoute, setActiveRoute] = useState("local");

  // Live gateway health probe
  const [healthState, setHealthState] = useState("checking"); // "checking" | "online" | "offline"
  const [healthLatencyMs, setHealthLatencyMs] = useState(null);

  // Unified Create/Edit Key modal state
  const [formModal, setFormModal] = useState({ open: false, mode: "create", initialKey: null });
  const [formAllowedModels, setFormAllowedModels] = useState("*");
  const [activeProviders, setActiveProviders] = useState([]);
  const [modelAliases, setModelAliases] = useState({});
  const [showModelPicker, setShowModelPicker] = useState(false);

  // Secondary modals
  const [createdKey, setCreatedKey] = useState(null);
  const [confirmState, setConfirmState] = useState(null);
  const [snippetKeyObj, setSnippetKeyObj] = useState(null);
  const [shareKey, setShareKey] = useState(null);

  // Bulk edit state (admin-only)
  const [selectedKeyIds, setSelectedKeyIds] = useState(new Set());
  const [showBulkEdit, setShowBulkEdit] = useState(false);
  const selectAllRef = useRef(null);
  const securitySectionRef = useRef(null);

  // Security & auth state
  const [requireApiKey, setRequireApiKey] = useState(false);
  const [tunnelDashboardAccess, setTunnelDashboardAccess] = useState(false);
  const [authStatus, setAuthStatus] = useState(null);

  const isApiKeyUser = authStatus?.role === "apikey";
  const sessionApiKey = authStatus?.apiKey || null;
  const isOwnKey = useCallback(
    (key) => isApiKeyUser && Boolean(key) && key.key === sessionApiKey,
    [isApiKeyUser, sessionApiKey]
  );

  const creatorPermissions = authStatus?.permissions || {
    ...EMPTY_PERMISSIONS,
    manageApiKeys: true,
    manageModels: true,
    manageProviders: true,
    manageTools: true,
    manageAdvanced: true,
    managePlugins: true,
    manageMediaProviders: true,
  };
  const creatorTokenLimit = authStatus?.tokenLimit || 0;
  const creatorAllowedModels = authStatus?.allowedModels || "*";
  const scopedModelPatterns =
    isApiKeyUser && creatorAllowedModels !== "*"
      ? String(creatorAllowedModels)
          .split(",")
          .map((s) => s.trim().toLowerCase())
          .filter(Boolean)
      : null;

  // Cloudflare Tunnel state
  const [tunnelEnabled, setTunnelEnabled] = useState(false);
  const [tunnelReachable, setTunnelReachable] = useState(false);
  const [tunnelUrl, setTunnelUrl] = useState("");
  const [tunnelPublicUrl, setTunnelPublicUrl] = useState("");
  const [tunnelLoading, setTunnelLoading] = useState(false);
  const [tunnelProgress, setTunnelProgress] = useState("");
  const [tunnelStatus, setTunnelStatus] = useState(null);
  const [tunnelAuthGateWarning, setTunnelAuthGateWarning] = useState(false);
  const [showEnableTunnelModal, setShowEnableTunnelModal] = useState(false);
  const [showDisableTunnelModal, setShowDisableTunnelModal] = useState(false);

  // Tailscale state
  const [tsEnabled, setTsEnabled] = useState(false);
  const [tsPort, setTsPort] = useState(null);
  const [tsReachable, setTsReachable] = useState(false);
  const [tsUrl, setTsUrl] = useState("");
  const [tsLoading, setTsLoading] = useState(false);
  const [tsProgress, setTsProgress] = useState("");
  const [tsStatus, setTsStatus] = useState(null);
  const [tsAuthUrl, setTsAuthUrl] = useState("");
  const [tsAuthLabel, setTsAuthLabel] = useState("");
  const [tsInstalled, setTsInstalled] = useState(null);
  const [tsInstalling, setTsInstalling] = useState(false);
  const [tsInstallLog, setTsInstallLog] = useState([]);
  const [tsSudoPassword, setTsSudoPassword] = useState("");
  const [tsConnecting, setTsConnecting] = useState(false);
  const [showTsModal, setShowTsModal] = useState(false);
  const [showDisableTsModal, setShowDisableTsModal] = useState(false);
  const tsLogRef = useRef(null);

  // Debounce reachable=false: only flip to false after N consecutive misses.
  const tunnelMissRef = useRef(0);
  const tsMissRef = useRef(0);
  const tunnelClientReachableRef = useRef(false);
  const tsClientReachableRef = useRef(false);

  // API key visibility toggle state
  const [visibleKeys, setVisibleKeys] = useState(new Set());

  const { copied, copy } = useCopyToClipboard();

  const [localOrigin] = useState(() => {
    if (typeof window !== "undefined" && window.location?.origin) {
      return window.location.origin;
    }
    return "http://localhost:20128";
  });

  const probeGatewayHealth = useCallback(async () => {
    setHealthState("checking");
    const t0 = typeof performance !== "undefined" ? performance.now() : Date.now();
    try {
      const res = await fetch("/api/health", { cache: "no-store" });
      const elapsed = Math.max(
        1,
        Math.round((typeof performance !== "undefined" ? performance.now() : Date.now()) - t0)
      );
      if (res.ok) {
        setHealthState("online");
        setHealthLatencyMs(elapsed);
      } else {
        setHealthState("offline");
        setHealthLatencyMs(null);
      }
    } catch {
      setHealthState("offline");
      setHealthLatencyMs(null);
    }
  }, []);

  const updateReachable = useCallback((clientRef, missRef, setter) => {
    if (clientRef.current) {
      missRef.current = 0;
      setter(true);
    } else {
      missRef.current += 1;
      if (missRef.current >= REACHABLE_MISS_THRESHOLD) setter(false);
    }
  }, []);

  const syncTunnelStatus = useCallback(async () => {
    try {
      const statusRes = await fetch("/api/tunnel/status", { cache: "no-store" });
      if (!statusRes.ok) return;
      const data = await statusRes.json();
      const tRunning = data.tunnel?.running === true;
      const tEnabled = tRunning || data.tunnel?.enabled === true;
      const tUrl = data.tunnel?.tunnelUrl || "";
      if (tEnabled && !tunnelLoading) {
        setTunnelUrl(tUrl);
        setTunnelPublicUrl(data.tunnel?.publicUrl || "");
        setTunnelEnabled(true);
      } else if (!tunnelLoading && !tRunning && data.tunnel?.settingsEnabled !== true) {
        setTunnelEnabled(false);
      }
      updateReachable(tunnelClientReachableRef, tunnelMissRef, setTunnelReachable);

      const tsEn = data.tailscale?.settingsEnabled ?? data.tailscale?.running ?? data.tailscale?.enabled ?? false;
      setTsUrl(data.tailscale?.tunnelUrl || "");
      setTsEnabled(tsEn);
      setTsPort(data.tailscale?.port ?? data.tailscale?.localPort ?? null);
      updateReachable(tsClientReachableRef, tsMissRef, setTsReachable);
    } catch {
      /* ignore transient poll errors */
    }
  }, [tunnelLoading, updateReachable]);

  const loadSettings = useCallback(async () => {
    try {
      const [settingsData, statusRes] = await Promise.allSettled([
        useSettingsStore.getState().fetchSettings(),
        fetch("/api/tunnel/status", { cache: "no-store" }),
      ]);
      const sVal = settingsData.status === "fulfilled" ? settingsData.value : null;
      if (sVal) {
        setRequireApiKey(Boolean(sVal.requireApiKey));
        setTunnelDashboardAccess(Boolean(sVal.tunnelDashboardAccess));
      }
      if (statusRes.status === "fulfilled" && statusRes.value.ok) {
        const data = await statusRes.value.json();
        const tRunning = data.tunnel?.running === true;
        const tEnabled = tRunning || data.tunnel?.enabled === true;
        setTunnelUrl(data.tunnel?.tunnelUrl || "");
        setTunnelPublicUrl(data.tunnel?.publicUrl || "");
        setTunnelEnabled(tEnabled);
        updateReachable(tunnelClientReachableRef, tunnelMissRef, setTunnelReachable);

        const tsEn = data.tailscale?.settingsEnabled ?? data.tailscale?.running ?? data.tailscale?.enabled ?? false;
        setTsUrl(data.tailscale?.tunnelUrl || "");
        setTsEnabled(tsEn);
        setTsPort(data.tailscale?.port ?? data.tailscale?.localPort ?? null);
        updateReachable(tsClientReachableRef, tsMissRef, setTsReachable);
      }
    } catch (error) {
      console.error("Error loading endpoint settings:", error);
    }
  }, [updateReachable]);

  const fetchData = useCallback(async () => {
    try {
      const fetchKeys = async () => {
        const res = await fetch("/api/keys");
        if (!res.ok) return [];
        const data = await res.json();
        return data.keys || [];
      };

      const fetchProvidersAndAliases = async () => {
        try {
          const [providersRes, aliasesRes] = await Promise.all([
            fetch("/api/providers"),
            fetch("/api/models/alias"),
          ]);
          if (providersRes.ok) {
            const pData = await providersRes.json();
            setActiveProviders(pData.connections || []);
          }
          if (aliasesRes.ok) {
            const aData = await aliasesRes.json();
            setModelAliases(aData.aliases || {});
          }
        } catch (e) {
          console.error("Error fetching providers/aliases:", e);
        }
      };

      fetchProvidersAndAliases();

      let existing = await fetchKeys();
      if (existing.length === 0 && !defaultKeyProvisionInFlight) {
        defaultKeyProvisionInFlight = true;
        try {
          const createRes = await fetch("/api/keys", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name: "Default Keys" }),
          });
          if (createRes.ok) existing = await fetchKeys();
        } catch {
          /* fall through to empty render */
        } finally {
          defaultKeyProvisionInFlight = false;
        }
      }
      setKeys(existing);
    } catch (error) {
      console.error("Error fetching endpoint data:", error);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetch("/api/auth/status")
      .then((res) => res.json())
      .then((data) => setAuthStatus(data))
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (tsLogRef.current) tsLogRef.current.scrollTop = tsLogRef.current.scrollHeight;
  }, [tsInstallLog]);

  useEffect(() => {
    probeGatewayHealth();
    fetchData();
    loadSettings();
  }, [probeGatewayHealth, fetchData, loadSettings]);

  // Live real-time key usage refresh (polls keys every 8s when tab visible)
  useEffect(() => {
    const pollTimer = setInterval(() => {
      if (document.hidden) return;
      fetch("/api/keys")
        .then((res) => (res.ok ? res.json() : null))
        .then((data) => {
          if (Array.isArray(data?.keys)) {
            setKeys(data.keys);
          }
        })
        .catch(() => {});
    }, 8000);
    return () => clearInterval(pollTimer);
  }, []);

  // Adaptive tunnel status poll while degraded + visibility re-check
  useEffect(() => {
    const tunnelHealthy = !tunnelEnabled || tunnelReachable;
    const tsHealthy = !tsEnabled || tsReachable;
    const allHealthy = tunnelHealthy && tsHealthy;
    const onVisible = () => {
      if (!document.hidden) syncTunnelStatus();
    };
    document.addEventListener("visibilitychange", onVisible);
    if (allHealthy) return () => document.removeEventListener("visibilitychange", onVisible);
    const timer = setInterval(() => {
      if (!document.hidden) syncTunnelStatus();
    }, STATUS_POLL_FAST_MS);
    return () => {
      clearInterval(timer);
      document.removeEventListener("visibilitychange", onVisible);
    };
  }, [tunnelEnabled, tsEnabled, tunnelReachable, tsReachable, syncTunnelStatus]);

  // Browser-side periodic CORS ping for tunnel/Tailscale URLs
  useEffect(() => {
    const probeBoth = async () => {
      if (document.hidden) return;
      if (tunnelEnabled && (tunnelUrl || tunnelPublicUrl)) {
        const ok = await clientPingAny(tunnelPublicUrl, tunnelUrl);
        tunnelClientReachableRef.current = ok;
        if (ok) {
          tunnelMissRef.current = 0;
          setTunnelReachable(true);
        } else {
          tunnelMissRef.current += 1;
          if (tunnelMissRef.current >= REACHABLE_MISS_THRESHOLD) setTunnelReachable(false);
        }
      } else {
        tunnelClientReachableRef.current = false;
      }
      if (tsEnabled && tsUrl) {
        const ok = await clientPingUrl(tsUrl);
        tsClientReachableRef.current = ok;
        if (ok) {
          tsMissRef.current = 0;
          setTsReachable(true);
        } else {
          tsMissRef.current += 1;
          if (tsMissRef.current >= REACHABLE_MISS_THRESHOLD) setTsReachable(false);
        }
      } else {
        tsClientReachableRef.current = false;
      }
    };
    const anyEnabled = (tunnelEnabled && (tunnelUrl || tunnelPublicUrl)) || (tsEnabled && tsUrl);
    if (!anyEnabled) return;
    probeBoth();
    const tunnelHealthy = !tunnelEnabled || tunnelReachable;
    const tsHealthy = !tsEnabled || tsReachable;
    if (tunnelHealthy && tsHealthy) return;
    const id = setInterval(probeBoth, CLIENT_PING_FAST_MS);
    return () => clearInterval(id);
  }, [tunnelEnabled, tunnelUrl, tunnelPublicUrl, tsEnabled, tsUrl, tunnelReachable, tsReachable]);

  const handleTunnelDashboardAccess = async (value) => {
    try {
      const updated = await useSettingsStore.getState().patchSettings({ tunnelDashboardAccess: value });
      if (updated) {
        setTunnelDashboardAccess(value);
        useNotificationStore
          .getState()
          .success(value ? "Dashboard access over tunnel enabled." : "Dashboard access over tunnel disabled.");
      }
    } catch {
      useNotificationStore.getState().error("Failed to update tunnel dashboard access setting.");
    }
  };

  const handleRequireApiKey = async (value) => {
    try {
      const updated = await useSettingsStore.getState().patchSettings({ requireApiKey: value });
      if (updated) {
        setRequireApiKey(value);
        if (value) setTunnelAuthGateWarning(false);
        useNotificationStore
          .getState()
          .success(value ? "API key authentication is now required." : "API key authentication disabled.");
        return true;
      }
    } catch {
      useNotificationStore.getState().error("Failed to update API key authentication setting.");
    }
    return false;
  };

  const handleEnableTunnelClick = () => {
    if (!requireApiKey) {
      setTunnelAuthGateWarning(true);
      setTunnelStatus({
        type: "warning",
        message: "Security required: enable API Key Authentication before activating Cloudflare Tunnel.",
      });
      return;
    }
    setTunnelAuthGateWarning(false);
    setTunnelStatus(null);
    setShowEnableTunnelModal(true);
  };

  const handleEnableAuthAndTunnel = async () => {
    const ok = await handleRequireApiKey(true);
    if (ok) {
      setTunnelStatus(null);
      setTunnelAuthGateWarning(false);
      setShowEnableTunnelModal(true);
    }
  };

  // Cloudflare Tunnel lifecycle
  const pingTunnelHealth = async (...urls) => {
    setTunnelLoading(true);
    setTunnelProgress("Waiting for tunnel to become reachable…");
    const targets = urls.filter(Boolean).map((u) => `${u}/api/health`);
    const start = Date.now();
    while (Date.now() - start < TUNNEL_PING_MAX_MS) {
      await new Promise((r) => setTimeout(r, TUNNEL_PING_INTERVAL_MS));
      const ok = await Promise.any(
        targets.map(async (h) => {
          const p = await fetch(h, { mode: "cors", cache: "no-store" });
          if (p.ok) return true;
          throw new Error("not ready");
        })
      ).catch(() => false);
      if (ok) {
        tunnelClientReachableRef.current = true;
        setTunnelEnabled(true);
        setTunnelReachable(true);
        setTunnelLoading(false);
        setTunnelProgress("");
        setActiveRoute("cloudflare");
        useNotificationStore.getState().success("Cloudflare Tunnel is online.");
        return true;
      }
      if ((Date.now() - start) % 10000 < TUNNEL_PING_INTERVAL_MS) {
        try {
          const statusRes = await fetch("/api/tunnel/status");
          if (statusRes.ok) {
            const status = await statusRes.json();
            if (!status.tunnel?.enabled) {
              setTunnelStatus({ type: "error", message: "Tunnel process stopped unexpectedly." });
              setTunnelLoading(false);
              setTunnelProgress("");
              return false;
            }
          }
        } catch {
          /* ignore */
        }
      }
    }
    setTunnelStatus({ type: "error", message: "Tunnel created but not reachable yet. Please try again." });
    setTunnelLoading(false);
    setTunnelProgress("");
    return false;
  };

  const handleEnableTunnel = async () => {
    setShowEnableTunnelModal(false);
    setTunnelLoading(true);
    setTunnelStatus(null);
    setTunnelProgress("Creating tunnel…");

    let polling = true;
    const pollProgress = async () => {
      while (polling) {
        try {
          const r = await fetch("/api/tunnel/status");
          if (r.ok) {
            const s = await r.json();
            if (s.download?.downloading) {
              setTunnelProgress(`Downloading cloudflared… ${s.download.progress}%`);
            } else if (polling) {
              setTunnelProgress("Creating tunnel…");
            }
          }
        } catch {
          /* ignore */
        }
        await new Promise((r) => setTimeout(r, 1000));
      }
    };
    pollProgress();

    try {
      const res = await fetch("/api/tunnel/enable", { method: "POST" });
      polling = false;
      const data = await res.json();
      if (!res.ok) {
        setTunnelStatus({ type: "error", message: data.error || "Failed to enable tunnel." });
        return;
      }
      const url = data.tunnelUrl;
      if (!url) {
        setTunnelStatus({ type: "error", message: "No tunnel URL was returned by the server." });
        return;
      }
      setTunnelUrl(url);
      setTunnelPublicUrl(data.publicUrl || "");
      await pingTunnelHealth(data.publicUrl, url);
    } catch (error) {
      setTunnelStatus({ type: "error", message: error.message });
    } finally {
      polling = false;
      setTunnelLoading(false);
      setTunnelProgress("");
    }
  };

  const handleDisableTunnel = async () => {
    setTunnelLoading(true);
    setTunnelStatus(null);
    try {
      const res = await fetch("/api/tunnel/disable", { method: "POST" });
      const data = await res.json();
      if (res.ok) {
        setTunnelEnabled(false);
        setTunnelReachable(false);
        setTunnelUrl("");
        setTunnelPublicUrl("");
        setShowDisableTunnelModal(false);
        setTunnelStatus({ type: "success", message: "Cloudflare Tunnel disabled." });
        if (activeRoute === "cloudflare") setActiveRoute("local");
      } else {
        setTunnelStatus({ type: "error", message: data.error || "Failed to disable tunnel." });
      }
    } catch (error) {
      setTunnelStatus({ type: "error", message: error.message });
    } finally {
      setTunnelLoading(false);
    }
  };

  // Tailscale lifecycle
  const checkTailscaleInstalled = async () => {
    setTsInstalled(null);
    try {
      const res = await fetch("/api/tunnel/tailscale-check");
      if (res.ok) {
        const data = await res.json();
        setTsInstalled(data.installed);
        return data;
      }
    } catch {
      /* ignore */
    }
    setTsInstalled(false);
    return { installed: false };
  };

  const requestUserAuth = (url, label) => {
    setTsAuthUrl(url);
    setTsAuthLabel(label);
  };

  const clearUserAuth = () => {
    setTsAuthUrl("");
    setTsAuthLabel("");
  };

  const pingTsHealth = async (url) => {
    setTsProgress("Waiting for Tailscale Funnel to become reachable…");
    const healthUrl = `${url}/api/health`;
    const start = Date.now();
    while (Date.now() - start < TUNNEL_PING_MAX_MS) {
      await new Promise((r) => setTimeout(r, TUNNEL_PING_INTERVAL_MS));
      try {
        const ping = await fetch(healthUrl, { mode: "no-cors", cache: "no-store" });
        if (ping.ok || ping.type === "opaque") return true;
      } catch {
        /* not ready yet */
      }
    }
    return false;
  };

  const pollFunnelEnable = async (enableUrl) => {
    requestUserAuth(enableUrl, "Open Funnel Settings");
    setTsProgress('Click "Open Funnel Settings" to enable Funnel…');
    for (let i = 0; i < 40; i++) {
      await new Promise((r) => setTimeout(r, 3000));
      try {
        const res = await fetch("/api/tunnel/tailscale-enable", { method: "POST" });
        const data = await res.json();
        if (res.ok && data.success) {
          clearUserAuth();
          setTsUrl(data.tunnelUrl || "");
          const ok3 = await pingTsHealth(data.tunnelUrl);
          setTsEnabled(true);
          if (ok3) {
            tsClientReachableRef.current = true;
            setTsReachable(true);
            setActiveRoute("tailscale");
            useNotificationStore.getState().success("Tailscale Funnel is online.");
          }
          setTsStatus(ok3 ? null : { type: "warning", message: "Connected, waiting for Funnel reachability." });
          return;
        }
        if (data.funnelNotEnabled) continue;
        if (data.error) {
          clearUserAuth();
          setTsStatus({ type: "error", message: data.error });
          return;
        }
      } catch {
        /* retry */
      }
    }
    clearUserAuth();
    setTsStatus({ type: "error", message: "Timed out waiting for Tailscale Funnel to be enabled." });
  };

  const handleConnectTailscale = async () => {
    setShowTsModal(false);
    setTsConnecting(true);
    setTsLoading(true);
    setTsStatus(null);
    setTsProgress("Connecting…");
    clearUserAuth();
    try {
      const res = await fetch("/api/tunnel/tailscale-enable", { method: "POST" });
      const data = await res.json();

      if (res.ok && data.success) {
        setTsUrl(data.tunnelUrl || "");
        const reachable = await pingTsHealth(data.tunnelUrl);
        setTsEnabled(true);
        if (reachable) {
          tsClientReachableRef.current = true;
          setTsReachable(true);
          setActiveRoute("tailscale");
          useNotificationStore.getState().success("Tailscale Funnel is online.");
        }
        setTsStatus(reachable ? null : { type: "warning", message: "Connected, waiting for Funnel reachability." });
        return;
      }

      if (data.needsLogin && data.authUrl) {
        requestUserAuth(data.authUrl, "Open Login Page");
        setShowTsModal(true);
        setTsProgress('Login required — click "Open Login Page" to continue');
        for (let i = 0; i < 40; i++) {
          await new Promise((r) => setTimeout(r, 3000));
          try {
            const r2 = await fetch("/api/tunnel/tailscale-check");
            if (r2.ok) {
              const check = await r2.json();
              if (check.loggedIn) {
                clearUserAuth();
                setTsProgress("Starting funnel…");
                const res2 = await fetch("/api/tunnel/tailscale-enable", { method: "POST" });
                const data2 = await res2.json();
                if (res2.ok && data2.success) {
                  setTsUrl(data2.tunnelUrl || "");
                  const ok2 = await pingTsHealth(data2.tunnelUrl);
                  setTsEnabled(true);
                  if (ok2) {
                    tsClientReachableRef.current = true;
                    setTsReachable(true);
                    setActiveRoute("tailscale");
                  }
                  setTsStatus(ok2 ? null : { type: "warning", message: "Connected, waiting for Funnel reachability." });
                } else if (data2.funnelNotEnabled && data2.enableUrl) {
                  await pollFunnelEnable(data2.enableUrl);
                } else {
                  setTsStatus({ type: "error", message: data2.error || "Failed to start Tailscale Funnel." });
                }
                return;
              }
            }
          } catch {
            /* retry */
          }
        }
        clearUserAuth();
        setTsStatus({ type: "error", message: "Tailscale login timed out. Please try again." });
        return;
      }

      if (data.funnelNotEnabled && data.enableUrl) {
        await pollFunnelEnable(data.enableUrl);
        return;
      }

      setTsStatus({ type: "error", message: data.error || "Failed to connect Tailscale." });
    } catch (error) {
      setTsStatus({ type: "error", message: error.message });
    } finally {
      setTsLoading(false);
      setTsConnecting(false);
      setTsProgress("");
      clearUserAuth();
    }
  };

  const handleInstallTailscale = async () => {
    setTsInstalling(true);
    setTsStatus(null);
    setTsInstallLog([]);
    try {
      const res = await fetch("/api/tunnel/tailscale-install", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ sudoPassword: tsSudoPassword }),
      });
      setTsSudoPassword("");

      if (!res.ok) {
        const errJson = await res.json().catch(() => ({}));
        setTsStatus({ type: "error", message: errJson.error || `HTTP ${res.status}: Installation request failed` });
        setTsInstalling(false);
        return;
      }

      const reader = res.body.getReader();
      const decoder = new TextDecoder();
      let buffer = "";

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        buffer += decoder.decode(value, { stream: true });
        const parts = buffer.split("\n\n");
        buffer = parts.pop() || "";
        for (const part of parts) {
          const lines = part.split("\n");
          let event = "progress";
          let data = null;
          for (const line of lines) {
            if (line.startsWith("event: ")) event = line.slice(7).trim();
            if (line.startsWith("data: ")) {
              try {
                data = JSON.parse(line.slice(6));
              } catch {
                /* skip */
              }
            }
          }
          if (!data) continue;
          if (event === "progress") {
            setTsInstallLog((prev) => [...prev.slice(-50), data.message]);
          } else if (event === "done") {
            setTsInstalled(true);
            setTsInstalling(false);
            setShowTsModal(false);
            handleConnectTailscale();
            return;
          } else if (event === "error") {
            setTsStatus({ type: "error", message: data.error || "Tailscale installation failed." });
          }
        }
      }
    } catch (e) {
      setTsStatus({ type: "error", message: e.message });
    } finally {
      setTsInstalling(false);
    }
  };

  const handleDisableTailscale = async () => {
    setTsLoading(true);
    setTsStatus(null);
    try {
      const res = await fetch("/api/tunnel/tailscale-disable", { method: "POST" });
      const data = await res.json();
      if (res.ok) {
        setTsEnabled(false);
        setTsReachable(false);
        setTsUrl("");
        setShowDisableTsModal(false);
        setTsStatus({ type: "success", message: "Tailscale Funnel disabled." });
        if (activeRoute === "tailscale") setActiveRoute("local");
      } else {
        setTsStatus({ type: "error", message: data.error || "Failed to disable Tailscale." });
      }
    } catch (e) {
      setTsStatus({ type: "error", message: e.message });
    } finally {
      setTsLoading(false);
    }
  };

  const handleOpenTsModal = async () => {
    setTsStatus(null);
    setTsInstallLog([]);
    const data = await checkTailscaleInstalled();
    if (data?.installed && data?.hasCachedPassword) {
      handleConnectTailscale();
    } else {
      setShowTsModal(true);
    }
  };

  // Key form modal helpers
  const openCreateModal = () => {
    setFormAllowedModels("*");
    setFormModal({ open: true, mode: "create", initialKey: null });
  };

  const openEditModal = (key) => {
    setFormAllowedModels(key.allowedModels || "*");
    setFormModal({ open: true, mode: "edit", initialKey: key });
  };

  const openDuplicateModal = (sourceKey) => {
    const base = `${sourceKey.name} (copy)`;
    let candidate = base;
    let n = 2;
    while (keys.some((k) => k.name === candidate)) {
      candidate = `${base} ${n}`;
      n += 1;
    }
    setFormAllowedModels(sourceKey.allowedModels || "*");
    setFormModal({
      open: true,
      mode: "create",
      initialKey: { ...sourceKey, id: undefined, name: candidate },
    });
  };

  const closeFormModal = () => {
    setFormModal({ open: false, mode: "create", initialKey: null });
    setFormAllowedModels("*");
  };

  const handleSelectModelForPicker = (model) => {
    if (!model?.value) return;
    const list = parseAllowedModelsList(formAllowedModels);
    if (list.includes(model.value)) return;
    setFormAllowedModels([...list, model.value].join(", "));
  };

  const handleDeselectModelForPicker = (model) => {
    if (!model?.value) return;
    const next = parseAllowedModelsList(formAllowedModels).filter((m) => m !== model.value);
    setFormAllowedModels(next.length === 0 ? "*" : next.join(", "));
  };

  const handleFormSubmit = async (payload) => {
    if (formModal.mode === "create") {
      const res = await fetch("/api/keys", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) {
        return { ok: false, error: data?.error || "Failed to create API key." };
      }
      setCreatedKey({ key: data.key, name: data.name || payload.name });
      closeFormModal();
      await fetchData();
      useNotificationStore.getState().success(`API key "${payload.name}" created.`);
      return { ok: true };
    }

    const targetId = formModal.initialKey?.id;
    if (!targetId) return { ok: false, error: "Missing target key ID." };
    const res = await fetch(`/api/keys/${targetId}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      return { ok: false, error: data?.error || "Failed to update API key." };
    }
    closeFormModal();
    await fetchData();
    useNotificationStore.getState().success(`API key "${payload.name}" updated.`);
    return { ok: true };
  };

  const handleToggleKeyActive = async (key, isActive) => {
    const previous = key.isActive !== false;
    setKeys((prev) => prev.map((k) => (k.id === key.id ? { ...k, isActive } : k)));
    try {
      const res = await fetch(`/api/keys/${key.id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ isActive }),
      });
      if (!res.ok) {
        const errBody = await res.json().catch(() => ({}));
        setKeys((prev) => prev.map((k) => (k.id === key.id ? { ...k, isActive: previous } : k)));
        useNotificationStore.getState().error(errBody?.error || "Failed to update key state.");
        return;
      }
      useNotificationStore
        .getState()
        .success(`"${key.name}" ${isActive ? "activated" : "disabled"}.`);
    } catch {
      setKeys((prev) => prev.map((k) => (k.id === key.id ? { ...k, isActive: previous } : k)));
      useNotificationStore.getState().error("Failed to update key state.");
      return;
    }
    await fetchData();
  };

  const handleManualResetUsage = (key) => {
    setConfirmState({
      title: "Reset Token Usage",
      message: `Reset used tokens for "${key.name}" back to 0?`,
      onConfirm: async () => {
        setConfirmState(null);
        try {
          const res = await fetch(`/api/keys/${key.id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ usedTokens: 0, lastResetAt: new Date().toISOString() }),
          });
          if (res.ok) {
            await fetchData();
            useNotificationStore.getState().success(`Token usage reset for "${key.name}".`);
          } else {
            const errBody = await res.json().catch(() => ({}));
            useNotificationStore.getState().error(errBody?.error || "Failed to reset token usage.");
          }
        } catch {
          useNotificationStore.getState().error("Failed to reset token usage.");
        }
      },
    });
  };

  const handleDeleteKey = (id) => {
    const target = keys.find((k) => k.id === id);
    setConfirmState({
      title: "Delete API Key",
      message: `Delete "${target?.name || "this API key"}"? Clients using it will immediately lose access.`,
      onConfirm: async () => {
        setConfirmState(null);
        try {
          const res = await fetch(`/api/keys/${id}`, { method: "DELETE" });
          if (res.ok) {
            setKeys((prev) => prev.filter((k) => k.id !== id));
            setSelectedKeyIds((prev) => {
              const next = new Set(prev);
              next.delete(id);
              return next;
            });
            setVisibleKeys((prev) => {
              const next = new Set(prev);
              next.delete(id);
              return next;
            });
            useNotificationStore.getState().success("API key deleted.");
          } else {
            const errBody = await res.json().catch(() => null);
            useNotificationStore.getState().error(errBody?.error || "Failed to delete API key.");
          }
        } catch {
          useNotificationStore.getState().error("Failed to delete API key.");
        }
      },
    });
  };

  const toggleKeyVisibility = (keyId) => {
    setVisibleKeys((prev) => {
      const next = new Set(prev);
      if (next.has(keyId)) next.delete(keyId);
      else next.add(keyId);
      return next;
    });
  };

  const creatorLabelFor = (key) => {
    const creator = key?.createdByName || "";
    if (creator) return `Created by ${creator}`;
    if (key?.createdBy === "dashboard") return "Created by dashboard";
    return "";
  };

  // Filtered keys & counts
  const bucketCounts = useMemo(() => countKeyBuckets(keys), [keys]);
  const filteredKeys = useMemo(
    () => filterApiKeys(keys, { query: keySearch, statusFilter }),
    [keys, keySearch, statusFilter]
  );

  const canBulkManage = !isApiKeyUser;
  const selectableFilteredKeys = useMemo(
    () => (canBulkManage ? filteredKeys.filter((k) => !isOwnKey(k)) : []),
    [canBulkManage, filteredKeys, isOwnKey]
  );

  const selectedVisibleCount = useMemo(
    () => selectableFilteredKeys.filter((k) => selectedKeyIds.has(k.id)).length,
    [selectableFilteredKeys, selectedKeyIds]
  );

  const allVisibleSelected =
    selectableFilteredKeys.length > 0 && selectedVisibleCount === selectableFilteredKeys.length;
  const someVisibleSelected = selectedVisibleCount > 0 && !allVisibleSelected;

  useEffect(() => {
    if (selectAllRef.current) {
      selectAllRef.current.indeterminate = someVisibleSelected;
    }
  }, [someVisibleSelected]);

  const toggleKeySelection = (id) => {
    if (!canBulkManage) return;
    setSelectedKeyIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const toggleSelectAllKeys = () => {
    if (!canBulkManage) return;
    if (allVisibleSelected) {
      setSelectedKeyIds((prev) => {
        const next = new Set(prev);
        for (const k of selectableFilteredKeys) next.delete(k.id);
        return next;
      });
    } else {
      setSelectedKeyIds((prev) => {
        const next = new Set(prev);
        for (const k of selectableFilteredKeys) next.add(k.id);
        return next;
      });
    }
  };

  const handleBulkApply = async (patch) => {
    const ids = [...selectedKeyIds];
    if (ids.length === 0) return;
    const { acl, ...keyPatch } = patch || {};
    if (Object.keys(keyPatch).length) {
      const res = await fetch("/api/keys/bulk", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ ids, patch: keyPatch }),
      });
      if (!res.ok) {
        const errBody = await res.json().catch(() => ({}));
        throw new Error(errBody?.error || "Failed to apply bulk changes");
      }
      setKeys((prev) => prev.map((k) => (selectedKeyIds.has(k.id) ? { ...k, ...keyPatch } : k)));
    }
    if (acl && Object.keys(acl).length) {
      await Promise.all(
        ids.map((id) =>
          fetch(`/api/keys/${id}/acl`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(acl),
          }).catch(() => {})
        )
      );
    }
    await fetchData();
    setSelectedKeyIds(new Set());
    useNotificationStore.getState().success(`Updated ${ids.length} API key${ids.length === 1 ? "" : "s"}.`);
  };

  const handleBulkDelete = () => {
    const ids = [...selectedKeyIds];
    if (ids.length === 0) return;
    setConfirmState({
      title: `Delete ${ids.length} API Key${ids.length === 1 ? "" : "s"}`,
      message: `Delete ${ids.length} selected API key${ids.length === 1 ? "" : "s"}? This action cannot be undone.`,
      onConfirm: async () => {
        setConfirmState(null);
        try {
          const res = await fetch("/api/keys/bulk", {
            method: "DELETE",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ ids }),
          });
          if (res.ok) {
            setKeys((prev) => prev.filter((k) => !selectedKeyIds.has(k.id)));
            setSelectedKeyIds(new Set());
            useNotificationStore
              .getState()
              .success(`Deleted ${ids.length} API key${ids.length === 1 ? "" : "s"}.`);
          } else {
            const err = await res.json().catch(() => null);
            useNotificationStore.getState().error(err?.error || "Bulk delete failed.");
          }
        } catch {
          useNotificationStore.getState().error("Bulk delete failed.");
        }
      },
    });
  };

  // Effective origin for Created Key cURL & Per-Key Snippet modal
  const activeEffectiveOrigin = useMemo(() => {
    if (!isApiKeyUser && activeRoute === "cloudflare" && tunnelEnabled && tunnelReachable) {
      return normalizeOrigin(tunnelPublicUrl || tunnelUrl) || normalizeOrigin(localOrigin);
    }
    if (!isApiKeyUser && activeRoute === "tailscale" && tsEnabled && tsReachable) {
      return normalizeOrigin(tsUrl) || normalizeOrigin(localOrigin);
    }
    return normalizeOrigin(localOrigin);
  }, [isApiKeyUser, activeRoute, tunnelEnabled, tunnelReachable, tunnelPublicUrl, tunnelUrl, tsEnabled, tsReachable, tsUrl, localOrigin]);

  const activeRouteLabel =
    !isApiKeyUser && activeRoute === "cloudflare" && tunnelEnabled && tunnelReachable
      ? "Cloudflare Tunnel"
      : !isApiKeyUser && activeRoute === "tailscale" && tsEnabled && tsReachable
        ? "Tailscale Mesh"
        : "Local Host";

  const defaultModel = useMemo(
    () => resolveDefaultSnippetModel(modelAliases, activeProviders),
    [modelAliases, activeProviders]
  );

  if (loading) {
    return (
      <div className="flex flex-col gap-5 max-w-full" data-machine-id={machineId}>
        <div className="flex items-center justify-between gap-4">
          <div className="space-y-2">
            <div className="h-6 w-48 bg-surface-2 rounded-lg animate-pulse" />
            <div className="h-3.5 w-72 bg-surface-2 rounded-lg animate-pulse" />
          </div>
          <div className="h-9 w-28 bg-surface-2 rounded-lg animate-pulse" />
        </div>
        <div className="p-5 rounded-2xl border border-border-subtle bg-surface space-y-4">
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <div className="h-16 bg-surface-2 rounded-xl animate-pulse" />
            <div className="h-16 bg-surface-2 rounded-xl animate-pulse" />
            <div className="h-16 bg-surface-2 rounded-xl animate-pulse" />
          </div>
          <div className="h-11 bg-surface-2 rounded-xl animate-pulse" />
        </div>
        <div className="p-5 rounded-2xl border border-border-subtle bg-surface space-y-3">
          <div className="h-5 w-32 bg-surface-2 rounded-lg animate-pulse" />
          <div className="h-28 bg-surface-2 rounded-xl animate-pulse" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-5 pb-16 lg:pb-8" data-machine-id={machineId}>
      {/* Studio Page Header */}
      <header className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3 min-w-0">
          <div className="size-9 rounded-xl bg-primary/10 border border-primary/20 text-primary flex items-center justify-center shrink-0">
            <span className="material-symbols-outlined text-[20px]">api</span>
          </div>
          <div className="min-w-0">
            <h1 className="text-lg sm:text-xl font-semibold text-text-main tracking-tight truncate">
              Endpoint &amp; Keys
            </h1>
            <p className="text-xs text-text-muted truncate">
              Connect clients, manage remote exposure, and govern API keys.
            </p>
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2 shrink-0">
          {/* Live Health Probe Chip */}
          <button
            type="button"
            onClick={probeGatewayHealth}
            title="Click to re-check gateway health"
            className={cn(
              "inline-flex items-center gap-1.5 h-8 px-2.5 rounded-lg border font-mono text-[11px] transition-colors cursor-pointer",
              healthState === "online"
                ? "bg-success-bg border-success-border text-success-fg hover:brightness-95"
                : healthState === "offline"
                  ? "bg-danger-bg border-danger-border text-danger-fg"
                  : "bg-surface-2 border-border-subtle text-text-muted"
            )}
          >
            <span
              className={cn(
                "size-1.5 rounded-full shrink-0",
                healthState === "online"
                  ? "bg-success"
                  : healthState === "offline"
                    ? "bg-danger"
                    : "bg-text-subtle"
              )}
            />
            <span>
              {healthState === "online"
                ? `Online · ${healthLatencyMs}ms`
                : healthState === "offline"
                  ? "Unreachable"
                  : "Checking…"}
            </span>
          </button>

          {/* Auth Posture Chip */}
          <button
            type="button"
            onClick={() => {
              if (!isApiKeyUser && securitySectionRef.current) {
                securitySectionRef.current.scrollIntoView({ behavior: "smooth", block: "center" });
              }
            }}
            className={cn(
              "inline-flex items-center gap-1.5 h-8 px-2.5 rounded-lg border font-mono text-[11px] transition-colors",
              !isApiKeyUser && "cursor-pointer hover:brightness-95",
              requireApiKey
                ? "bg-primary/10 border-primary/25 text-primary"
                : "bg-warning-bg border-warning-border text-warning-fg"
            )}
            title={isApiKeyUser ? "Gateway authentication status" : "Jump to security policy"}
          >
            <span className="material-symbols-outlined text-[14px]">
              {requireApiKey ? "lock" : "lock_open"}
            </span>
            <span>{requireApiKey ? "Auth Required" : "Public Endpoint"}</span>
          </button>

          <Button size="md" icon="add" onClick={openCreateModal}>
            Create Key
          </Button>
        </div>
      </header>

      {/* Zona 1: Unified Gateway & Exposure Console */}
      <GatewayConsole
        isApiKeyUser={isApiKeyUser}
        localOrigin={localOrigin}
        activeRoute={activeRoute}
        onSelectRoute={setActiveRoute}
        tunnel={{
          enabled: tunnelEnabled,
          reachable: tunnelReachable,
          loading: tunnelLoading,
          progress: tunnelProgress,
          url: tunnelUrl,
          publicUrl: tunnelPublicUrl,
          status: tunnelStatus,
          authGateWarning: tunnelAuthGateWarning,
        }}
        tailscale={{
          enabled: tsEnabled,
          reachable: tsReachable,
          loading: tsLoading,
          connecting: tsConnecting,
          progress: tsProgress,
          url: tsUrl,
          port: tsPort,
          status: tsStatus,
          authUrl: tsAuthUrl,
          authLabel: tsAuthLabel,
        }}
        requireApiKey={requireApiKey}
        tunnelDashboardAccess={tunnelDashboardAccess}
        onToggleRequireApiKey={handleRequireApiKey}
        onToggleTunnelDashboardAccess={handleTunnelDashboardAccess}
        onEnableTunnelClick={handleEnableTunnelClick}
        onEnableAuthAndTunnel={handleEnableAuthAndTunnel}
        onDisableTunnelClick={() => setShowDisableTunnelModal(true)}
        onConnectTailscaleClick={handleOpenTsModal}
        onDisableTailscaleClick={() => setShowDisableTsModal(true)}
        onDismissTunnelStatus={() => {
          setTunnelStatus(null);
          setTunnelAuthGateWarning(false);
        }}
        onDismissTsStatus={() => setTsStatus(null)}
        keys={keys}
        defaultModel={defaultModel}
        securitySectionRef={securitySectionRef}
      />

      {/* Zona 2: API Keys Vault */}
      <section className="space-y-3.5">
        {/* Vault Header & Search */}
        <div className="flex flex-col gap-3">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2.5">
            <div className="flex items-center gap-2.5">
              {canBulkManage && selectableFilteredKeys.length > 0 && (
                <input
                  ref={selectAllRef}
                  type="checkbox"
                  checked={allVisibleSelected}
                  onChange={toggleSelectAllKeys}
                  className="size-4 accent-[var(--color-primary)] cursor-pointer shrink-0"
                  aria-label="Select all visible API keys"
                  title="Select all visible keys"
                />
              )}
              <span className="material-symbols-outlined text-primary text-[19px]">vpn_key</span>
              <h2 className="text-sm sm:text-base font-semibold text-text-main tracking-tight">
                API Keys
              </h2>
              <span className="px-2 py-0.5 rounded-md bg-surface-2 border border-border-subtle font-mono text-[11px] text-text-muted">
                {filteredKeys.length === keys.length
                  ? keys.length
                  : `${filteredKeys.length} / ${keys.length}`}
              </span>
            </div>

            <div className="relative w-full sm:w-64">
              <span className="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-text-subtle text-[16px] pointer-events-none">
                search
              </span>
              <input
                type="text"
                placeholder="Search by key name or prefix…"
                value={keySearch}
                onChange={(e) => setKeySearch(e.target.value)}
                aria-label="Search API keys"
                className="field-control w-full h-9 pl-8 pr-8 text-xs"
              />
              {keySearch && (
                <button
                  type="button"
                  onClick={() => setKeySearch("")}
                  aria-label="Clear search"
                  className="absolute right-2 top-1/2 -translate-y-1/2 size-5 rounded flex items-center justify-center text-text-subtle hover:text-text-main cursor-pointer"
                >
                  <span className="material-symbols-outlined text-[15px]">close</span>
                </button>
              )}
            </div>
          </div>

          {/* Quick Status Filter Pills */}
          <div className="flex flex-wrap items-center gap-1.5" role="tablist" aria-label="Filter API keys by status">
            {FILTER_TABS.map((tab) => {
              const count = bucketCounts[tab.id] || 0;
              const active = statusFilter === tab.id;
              return (
                <button
                  key={tab.id}
                  type="button"
                  role="tab"
                  aria-selected={active}
                  onClick={() => setStatusFilter(tab.id)}
                  className={cn(
                    "inline-flex items-center gap-1.5 h-7 px-2.5 rounded-lg border text-xs font-medium transition-colors cursor-pointer",
                    active
                      ? "bg-primary/12 border-primary/30 text-primary"
                      : "bg-surface border-border-subtle text-text-muted hover:bg-surface-2 hover:text-text-main"
                  )}
                >
                  <span>{tab.label}</span>
                  <span
                    className={cn(
                      "font-mono text-[10.5px] px-1 rounded",
                      active ? "bg-primary/15 text-primary" : "text-text-subtle"
                    )}
                  >
                    {count}
                  </span>
                </button>
              );
            })}
          </div>
        </div>

        {/* Admin-Only Contextual Bulk Action Bar */}
        {canBulkManage && selectedKeyIds.size > 0 && (
          <div className="flex flex-wrap items-center justify-between gap-2.5 rounded-xl border border-primary/30 bg-primary/8 px-3.5 py-2.5">
            <div className="flex items-center gap-2.5 text-xs">
              <span className="font-semibold text-primary">
                {selectedKeyIds.size} key{selectedKeyIds.size === 1 ? "" : "s"} selected
              </span>
              <button
                type="button"
                onClick={() => setSelectedKeyIds(new Set())}
                className="text-text-muted hover:text-text-main underline cursor-pointer"
              >
                Clear
              </button>
            </div>
            <div className="flex items-center gap-2">
              <Button size="sm" variant="secondary" icon="edit" onClick={() => setShowBulkEdit(true)}>
                Bulk Edit
              </Button>
              <Button size="sm" variant="danger" icon="delete" onClick={handleBulkDelete}>
                Delete Selected
              </Button>
            </div>
          </div>
        )}

        {/* Key Cards or Empty State */}
        {filteredKeys.length === 0 ? (
          <div className="text-center py-12 rounded-2xl border border-dashed border-border bg-surface p-6">
            <span className="material-symbols-outlined text-3xl text-text-subtle opacity-60 mb-2 block">
              vpn_key
            </span>
            <p className="text-sm font-medium text-text-main">
              {keys.length === 0 ? "No API keys yet" : "No keys match your current filter"}
            </p>
            <p className="text-xs text-text-muted mt-1 mb-4">
              {keys.length === 0
                ? "Create your first API key to authenticate IDE and SDK clients."
                : "Try adjusting your search query or switching back to All keys."}
            </p>
            {keys.length === 0 ? (
              <Button size="sm" icon="add" onClick={openCreateModal}>
                Create Key
              </Button>
            ) : (
              <Button
                size="sm"
                variant="secondary"
                onClick={() => {
                  setKeySearch("");
                  setStatusFilter("all");
                }}
              >
                Reset Filters
              </Button>
            )}
          </div>
        ) : (
          <div className="space-y-2.5">
            {filteredKeys.map((key) => (
              <ApiKeyCard
                key={key.id}
                apiKey={key}
                creatorLabel={creatorLabelFor(key)}
                canBulkSelect={canBulkManage}
                selected={selectedKeyIds.has(key.id)}
                onToggleSelect={toggleKeySelection}
                isOwnKey={isOwnKey(key)}
                revealed={visibleKeys.has(key.id)}
                onToggleReveal={toggleKeyVisibility}
                copied={copied === key.id}
                onCopy={copy}
                onToggleActive={handleToggleKeyActive}
                onOpenSnippet={setSnippetKeyObj}
                onEdit={openEditModal}
                onDuplicate={openDuplicateModal}
                onResetUsage={handleManualResetUsage}
                onShare={setShareKey}
                onDelete={handleDeleteKey}
              />
            ))}
          </div>
        )}
      </section>

      {/* Unified Create / Edit Key Modal */}
      <KeyFormModal
        isOpen={formModal.open}
        mode={formModal.mode}
        initialKey={formModal.initialKey}
        existingKeys={keys}
        creatorPermissions={creatorPermissions}
        creatorTokenLimit={creatorTokenLimit}
        creatorAllowedModels={creatorAllowedModels}
        isApiKeyUser={isApiKeyUser}
        onOpenModelPicker={() => setShowModelPicker(true)}
        selectedModelsValue={formAllowedModels}
        onSelectedModelsChange={setFormAllowedModels}
        onSubmit={handleFormSubmit}
        onClose={closeFormModal}
      />

      {/* Created Key Confirmation Modal */}
      <Modal
        isOpen={Boolean(createdKey)}
        title="API Key Ready"
        onClose={() => setCreatedKey(null)}
        footer={
          <Button onClick={() => setCreatedKey(null)}>Done</Button>
        }
      >
        <div className="flex flex-col gap-4">
          <div className="rounded-xl border border-success-border bg-success-bg p-3.5 text-xs text-success-fg">
            <p className="font-semibold mb-0.5">
              {createdKey?.name ? `"${createdKey.name}" is ready to use.` : "Your new API key is ready."}
            </p>
            <p>Copy the key below or copy a pre-filled cURL test command right away.</p>
          </div>

          <div className="flex gap-2">
            <Input
              value={createdKey?.key || ""}
              readOnly
              className="flex-1 min-w-0 font-mono text-sm"
              inputClassName="truncate"
            />
            <Button
              variant="secondary"
              icon={copied === "created_key" ? "check" : "content_copy"}
              className="shrink-0"
              onClick={() => copy(createdKey?.key || "", "created_key")}
            >
              {copied === "created_key" ? "Copied!" : "Copy Key"}
            </Button>
          </div>

          <Button
            variant="outline"
            icon={copied === "created_curl" ? "check" : "terminal"}
            fullWidth
            onClick={() =>
              copy(
                generateEndpointSnippet({
                  format: "openai",
                  lang: "curl",
                  baseUrl: activeEffectiveOrigin,
                  apiKey: createdKey?.key || "",
                  model: defaultModel,
                }),
                "created_curl"
              )
            }
          >
            {copied === "created_curl" ? "cURL Copied!" : "Copy cURL"}
          </Button>
        </div>
      </Modal>

      {/* Per-Key Code Snippet Modal */}
      <KeySnippetModal
        apiKeyObj={snippetKeyObj}
        activeOrigin={activeEffectiveOrigin}
        activeRouteLabel={activeRouteLabel}
        defaultModel={defaultModel}
        onClose={() => setSnippetKeyObj(null)}
      />

      {/* Model Select Modal for Create/Edit Key */}
      {showModelPicker && (
        <ModelSelectModal
          isOpen={showModelPicker}
          onClose={() => setShowModelPicker(false)}
          onSelect={handleSelectModelForPicker}
          onDeselect={handleDeselectModelForPicker}
          activeProviders={activeProviders}
          modelAliases={modelAliases}
          showStudioTargets={true}
          title="Select Allowed Models"
          addedModelValues={parseAllowedModelsList(formAllowedModels)}
          allowedModelPatterns={scopedModelPatterns}
          closeOnSelect={false}
        />
      )}

      {/* Enable Cloudflare Tunnel Modal */}
      <Modal
        isOpen={showEnableTunnelModal}
        title="Enable Cloudflare Tunnel"
        onClose={() => setShowEnableTunnelModal(false)}
        footer={
          <>
            <Button variant="ghost" onClick={() => setShowEnableTunnelModal(false)}>
              Cancel
            </Button>
            <Button icon="cloud_upload" onClick={handleEnableTunnel}>
              Start Tunnel
            </Button>
          </>
        }
      >
        <div className="flex flex-col gap-4">
          <div className="bg-surface-2 border border-border-subtle rounded-xl p-4">
            <div className="flex items-start gap-3">
              <span className="material-symbols-outlined text-primary">cloud_upload</span>
              <div>
                <p className="text-sm text-text-main font-medium mb-1">Cloudflare Tunnel</p>
                <p className="text-xs text-text-muted leading-relaxed">
                  Expose your local OMEG4MON ROUTER to the internet without port forwarding or a static IP, then use the URL in Cursor, Cline, and other tools from anywhere.
                </p>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-2.5">
            {TUNNEL_BENEFITS.map((benefit) => (
              <div
                key={benefit.title}
                className="flex flex-col items-center text-center p-3 rounded-xl border border-border-subtle bg-surface-2/50"
              >
                <span className="material-symbols-outlined text-xl text-primary mb-1">{benefit.icon}</span>
                <p className="text-xs font-semibold text-text-main">{benefit.title}</p>
                <p className="text-[11px] text-text-muted mt-0.5">{benefit.desc}</p>
              </div>
            ))}
          </div>

          <p className="text-xs text-text-muted">
            Requires outbound port 7844 (TCP/UDP) and typically takes 10–30 seconds to connect.
          </p>
        </div>
      </Modal>

      {/* Disable Cloudflare Tunnel Modal */}
      <Modal
        isOpen={showDisableTunnelModal}
        title="Disable Cloudflare Tunnel"
        onClose={() => !tunnelLoading && setShowDisableTunnelModal(false)}
        footer={
          <>
            <Button variant="ghost" onClick={() => setShowDisableTunnelModal(false)} disabled={tunnelLoading}>
              Cancel
            </Button>
            <Button variant="danger" onClick={handleDisableTunnel} loading={tunnelLoading}>
              Disable Tunnel
            </Button>
          </>
        }
      >
        <p className="text-sm text-text-muted">
          The Cloudflare tunnel will disconnect immediately and its public URL will stop responding.
        </p>
      </Modal>

      {/* Tailscale Modal */}
      <Modal
        isOpen={showTsModal}
        title="Tailscale Funnel"
        onClose={() => {
          if (!tsInstalling) {
            setShowTsModal(false);
            setTsSudoPassword("");
            setTsStatus(null);
          }
        }}
        footer={
          tsInstalled === false && !tsInstalling ? (
            <>
              <Button variant="ghost" onClick={() => setShowTsModal(false)}>
                Cancel
              </Button>
              <Button onClick={handleInstallTailscale}>Install Tailscale</Button>
            </>
          ) : tsInstalled === true && !tsInstalling ? (
            <>
              <Button variant="ghost" onClick={() => setShowTsModal(false)}>
                Cancel
              </Button>
              <Button onClick={handleConnectTailscale}>Connect</Button>
            </>
          ) : (
            <Button variant="ghost" onClick={() => setShowTsModal(false)} disabled={tsInstalling}>
              Close
            </Button>
          )
        }
      >
        <div className="flex flex-col gap-4">
          {tsInstalled === null && (
            <p className="text-sm text-text-muted flex items-center gap-2">
              <span className="material-symbols-outlined animate-spin text-sm">progress_activity</span>
              Checking Tailscale installation…
            </p>
          )}

          {tsInstalled === false && !tsInstalling && (
            <p className="text-sm text-text-muted">
              Tailscale is not installed on this machine. Install it now to enable Funnel access.
            </p>
          )}

          {tsInstalling && (
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2 text-sm text-text-muted">
                <span className="material-symbols-outlined animate-spin text-sm">progress_activity</span>
                Installing Tailscale…
              </div>
              {tsInstallLog.length > 0 && (
                <div
                  ref={tsLogRef}
                  className="bg-surface-2 rounded-xl border border-border-subtle p-3 max-h-40 overflow-y-auto font-mono text-xs text-text-muted"
                >
                  {tsInstallLog.map((line, i) => (
                    <div key={i}>{line}</div>
                  ))}
                </div>
              )}
            </div>
          )}

          {tsAuthUrl && (
            <div className="flex flex-col gap-3 p-3.5 rounded-xl border border-primary/35 bg-primary/8">
              <div className="flex items-center gap-2 text-sm font-semibold text-text-main">
                <span className="material-symbols-outlined text-primary text-[18px]">key</span>
                Tailscale Authentication Required
              </div>
              <p className="text-xs text-text-muted">
                Open the Tailscale login page in a new tab to authorize this router:
              </p>
              <div className="flex gap-2">
                <a
                  href={tsAuthUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex-1 inline-flex items-center justify-center gap-1.5 px-3 py-2 rounded-lg bg-primary text-on-primary text-xs font-semibold hover:bg-primary-hover transition-colors text-center"
                >
                  <span className="material-symbols-outlined text-[16px]">open_in_new</span>
                  <span>{tsAuthLabel || "Open Login Page"}</span>
                </a>
                <Button onClick={() => copy(tsAuthUrl, "modal_auth_url")} variant="secondary" size="sm">
                  {copied === "modal_auth_url" ? "Copied" : "Copy"}
                </Button>
              </div>
            </div>
          )}

          {tsInstalled === true && !tsInstalling && (
            <div className="flex items-center gap-2 text-sm text-success">
              <span className="material-symbols-outlined text-[16px]">check_circle</span>
              Tailscale is installed and ready to connect.
            </div>
          )}

          {tsStatus && <StatusAlert status={tsStatus} />}
        </div>
      </Modal>

      {/* Disable Tailscale Modal */}
      <Modal
        isOpen={showDisableTsModal}
        title="Disable Tailscale Funnel"
        onClose={() => !tsLoading && setShowDisableTsModal(false)}
        footer={
          <>
            <Button variant="ghost" onClick={() => setShowDisableTsModal(false)} disabled={tsLoading}>
              Cancel
            </Button>
            <Button variant="danger" onClick={handleDisableTailscale} loading={tsLoading}>
              Disable Tailscale
            </Button>
          </>
        }
      >
        <p className="text-sm text-text-muted">
          Tailscale Funnel will stop immediately and its tailnet URL will become unreachable.
        </p>
      </Modal>

      {/* Confirm Modal */}
      <ConfirmModal
        isOpen={Boolean(confirmState)}
        onClose={() => setConfirmState(null)}
        onConfirm={confirmState?.onConfirm}
        title={confirmState?.title || "Confirm"}
        message={confirmState?.message}
        variant="danger"
      />

      {/* Bulk Edit Modal */}
      {showBulkEdit && (
        <BulkEditKeys
          selectedCount={selectedKeyIds.size}
          onApply={handleBulkApply}
          onClose={() => setShowBulkEdit(false)}
        />
      )}

      {/* Usage-share links for one key */}
      {shareKey && (
        <KeyShareModal keyId={shareKey.id} keyName={shareKey.name} onClose={() => setShareKey(null)} />
      )}
    </div>
  );
}

APIPageClient.propTypes = {
  machineId: PropTypes.string.isRequired,
};
