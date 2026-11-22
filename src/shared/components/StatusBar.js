"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { APP_CONFIG } from "@/shared/constants/config";
import LogoMark from "@/shared/components/LogoMark";

/* Status strip along the bottom of the shell — the instrument-panel tell that
   the upstream template does not have. Reports gateway health with latency,
   the endpoint host (click-to-copy), and the build, in mono type. */

function Dot({ state }) {
  const color =
    state === "ok" ? "bg-success" : state === "down" ? "bg-danger" : "bg-warning";
  return (
    <span className="relative flex size-2 shrink-0">
      {state === "ok" && (
        <span className="absolute inline-flex size-full animate-ping motion-reduce:animate-none rounded-full bg-success/60" />
      )}
      <span className={`relative inline-flex size-2 rounded-full ${color}`} />
    </span>
  );
}

export default function StatusBar() {
  const [health, setHealth] = useState("checking");
  const [latencyMs, setLatencyMs] = useState(null);
  const [checking, setChecking] = useState(true);
  const [host, setHost] = useState("");
  const [copied, setCopied] = useState(false);

  const mountedRef = useRef(true);
  const checkSeqRef = useRef(0);
  const copyTimerRef = useRef(null);

  const runHealthCheck = useCallback(async () => {
    const seq = ++checkSeqRef.current;
    setChecking(true);
    const start = performance.now();
    try {
      const res = await fetch("/api/health", { cache: "no-store" });
      const elapsed = Math.round(performance.now() - start);
      if (!mountedRef.current || seq !== checkSeqRef.current) return;
      if (res.ok) {
        setHealth("ok");
        setLatencyMs(elapsed);
      } else {
        setHealth("down");
        setLatencyMs(null);
      }
    } catch {
      if (!mountedRef.current || seq !== checkSeqRef.current) return;
      setHealth("down");
      setLatencyMs(null);
    } finally {
      if (mountedRef.current && seq === checkSeqRef.current) setChecking(false);
    }
  }, []);

  useEffect(() => {
    mountedRef.current = true;
    runHealthCheck();
    if (typeof window !== "undefined") setHost(window.location.host);
    const id = setInterval(runHealthCheck, 30000);
    return () => {
      mountedRef.current = false;
      clearInterval(id);
      if (copyTimerRef.current) {
        clearTimeout(copyTimerRef.current);
        copyTimerRef.current = null;
      }
    };
  }, [runHealthCheck]);

  const handleCopyHost = async () => {
    const target = host || "localhost";
    try {
      if (typeof navigator === "undefined" || !navigator.clipboard?.writeText) return;
      await navigator.clipboard.writeText(target);
      if (!mountedRef.current) return;
      setCopied(true);
      if (copyTimerRef.current) clearTimeout(copyTimerRef.current);
      copyTimerRef.current = setTimeout(() => {
        if (mountedRef.current) setCopied(false);
        copyTimerRef.current = null;
      }, 1800);
    } catch {
      // Clipboard API unavailable or blocked in non-HTTPS / restricted contexts
    }
  };

  const label =
    health === "ok"
      ? "Gateway online"
      : health === "down"
        ? "Gateway unreachable"
        : "Checking gateway";

  return (
    <div className="flex h-7 shrink-0 items-center gap-3 border-t border-border/80 bg-surface/80 px-3 font-mono text-[10.5px] text-text-subtle lg:px-5">
      <button
        type="button"
        onClick={runHealthCheck}
        title="Re-check gateway health"
        aria-label={`${label}${health === "ok" && latencyMs !== null ? `, ${latencyMs}ms` : ""}. Click to re-check gateway health`}
        className="group -ml-1.5 flex items-center gap-1.5 rounded-md px-1.5 py-0.5 transition-colors hover:bg-surface-2 hover:text-text-main cursor-pointer"
      >
        {checking ? (
          <span
            className="material-symbols-outlined animate-spin text-[12px] leading-none text-primary"
            aria-hidden="true"
          >
            progress_activity
          </span>
        ) : (
          <Dot state={health} />
        )}
        <span className={health === "down" ? "text-danger-fg" : "text-text-muted group-hover:text-text-main"}>
          {label}
        </span>
        {health === "ok" && latencyMs !== null ? (
          <span className="tabular-nums text-text-subtle">· {latencyMs}ms</span>
        ) : null}
      </button>

      <button
        type="button"
        onClick={handleCopyHost}
        title={`Copy endpoint host (${host || "localhost"})`}
        aria-label={copied ? `Endpoint host ${host || "localhost"} copied` : `Copy endpoint host ${host || "localhost"}`}
        className="hidden items-center gap-1.5 rounded-md px-1.5 py-0.5 transition-colors hover:bg-surface-2 hover:text-text-main cursor-pointer sm:flex"
      >
        <span
          className={`material-symbols-outlined text-[12px] leading-none ${
            copied ? "text-success" : "text-text-subtle"
          }`}
        >
          {copied ? "check" : "api"}
        </span>
        <span className={`truncate ${copied ? "text-success-fg" : ""}`}>
          {copied ? "Copied" : host || "localhost"}
        </span>
      </button>

      <span className="ml-auto hidden items-center gap-1.5 md:flex">
        <span className="material-symbols-outlined text-[12px] leading-none text-text-subtle">bolt</span>
        <span>OpenAI · Anthropic compatible</span>
      </span>

      <span className="flex items-center gap-1.5">
        <LogoMark className="size-3.5" />
        <span className="tabular-nums">{APP_CONFIG.name} v{APP_CONFIG.version}</span>
      </span>
    </div>
  );
}
