import { useEffect, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";

const DEFAULT_DURATION_MS = 420;
const TICK_CLASS_MS = 380;

export function formatIntegerDigits(value: number): string {
  if (!Number.isFinite(value)) {
    return "0";
  }
  return Math.round(value).toLocaleString();
}

export function easeOutCubic(progress: number): number {
  const clamped = Math.min(1, Math.max(0, progress));
  return 1 - (1 - clamped) ** 3;
}

export function interpolateMetricValue(from: number, to: number, progress: number): number {
  return from + (to - from) * easeOutCubic(progress);
}

function prefersReducedMotion(): boolean {
  if (typeof window === "undefined" || typeof window.matchMedia !== "function") {
    return false;
  }
  return window.matchMedia("(prefers-reduced-motion: reduce)").matches;
}

export interface AnimatedNumberProps {
  value: number;
  format?: (value: number) => string;
  durationMs?: number;
  className?: string;
}

export function AnimatedNumber({
  value,
  format = formatIntegerDigits,
  durationMs = DEFAULT_DURATION_MS,
  className,
}: AnimatedNumberProps) {
  const safeTarget = Number.isFinite(value) ? value : 0;
  const [displayValue, setDisplayValue] = useState(safeTarget);
  const [isTicking, setIsTicking] = useState(false);

  const currentDisplayRef = useRef(safeTarget);
  const previousTargetRef = useRef(safeTarget);
  const hasMountedRef = useRef(false);

  useEffect(() => {
    if (!hasMountedRef.current) {
      hasMountedRef.current = true;
      currentDisplayRef.current = safeTarget;
      previousTargetRef.current = safeTarget;
      setDisplayValue(safeTarget);
      return;
    }

    if (previousTargetRef.current === safeTarget) {
      return;
    }
    previousTargetRef.current = safeTarget;

    const from = currentDisplayRef.current;
    if (from === safeTarget || durationMs <= 0 || prefersReducedMotion()) {
      currentDisplayRef.current = safeTarget;
      setDisplayValue(safeTarget);
      return;
    }

    setIsTicking(true);
    const tickTimer = window.setTimeout(() => {
      setIsTicking(false);
    }, TICK_CLASS_MS);

    let frameId = 0;
    const startTime = performance.now();

    const step = (now: number) => {
      const elapsed = now - startTime;
      const progress = Math.min(1, elapsed / durationMs);
      const next = interpolateMetricValue(from, safeTarget, progress);
      currentDisplayRef.current = next;
      setDisplayValue(next);

      if (progress < 1) {
        frameId = window.requestAnimationFrame(step);
      } else {
        currentDisplayRef.current = safeTarget;
        setDisplayValue(safeTarget);
      }
    };

    frameId = window.requestAnimationFrame(step);

    return () => {
      window.cancelAnimationFrame(frameId);
      window.clearTimeout(tickTimer);
    };
  }, [safeTarget, durationMs]);

  return (
    <span className={cn("inline-block tabular-nums", isTicking && "metric-value-tick", className)}>
      {format(displayValue)}
    </span>
  );
}
