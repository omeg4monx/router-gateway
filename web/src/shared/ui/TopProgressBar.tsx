import { useIsFetching } from "@tanstack/react-query";
import { useCallback, useEffect, useRef, useState } from "react";
import { useLocation } from "react-router";

export const TOP_PROGRESS_INITIAL = 18;
export const TOP_PROGRESS_MAX_TRICKLE = 90;

const TRICKLE_INTERVAL_MS = 140;
const MIN_NAV_VISIBLE_MS = 160;
const COMPLETE_HOLD_MS = 140;
const FADE_OUT_MS = 200;
const MAX_ACTIVE_MS = 8000;

export type TopProgressPhase = "idle" | "running" | "completing" | "fading";

export function computeNextTrickleProgress(current: number): number {
  if (current >= TOP_PROGRESS_MAX_TRICKLE) {
    return TOP_PROGRESS_MAX_TRICKLE;
  }
  const remaining = TOP_PROGRESS_MAX_TRICKLE - current;
  const step = Math.max(0.6, remaining * 0.16);
  return Math.min(TOP_PROGRESS_MAX_TRICKLE, current + step);
}

export function shouldFinishTopProgress(params: {
  navLocked: boolean;
  holdCount: number;
  pendingQueries: number;
  timedOut: boolean;
}): boolean {
  if (params.timedOut) {
    return true;
  }
  return !params.navLocked && params.holdCount === 0 && params.pendingQueries === 0;
}

export function getTopProgressBarStyle(
  phase: TopProgressPhase,
  progress: number,
): {
  transform: string;
  opacity: number;
  transition?: string;
} {
  if (phase === "idle") {
    return {
      transform: "scaleX(0)",
      opacity: 0,
      transition: "none",
    };
  }
  if (phase === "completing") {
    return {
      transform: "scaleX(1)",
      opacity: 1,
    };
  }
  if (phase === "fading") {
    return {
      transform: "scaleX(1)",
      opacity: 0,
    };
  }
  return {
    transform: `scaleX(${Math.max(0, Math.min(100, progress)) / 100})`,
    opacity: 1,
  };
}

export interface TopProgressController {
  beginHold: () => () => void;
  getHoldCount: () => number;
  subscribe: (listener: (holdCount: number) => void) => () => void;
}

export function createTopProgressController(): TopProgressController {
  let holdCount = 0;
  const listeners = new Set<(count: number) => void>();

  const notify = () => {
    for (const listener of listeners) {
      listener(holdCount);
    }
  };

  return {
    beginHold() {
      let released = false;
      holdCount += 1;
      notify();
      return () => {
        if (released) {
          return;
        }
        released = true;
        holdCount = Math.max(0, holdCount - 1);
        notify();
      };
    },
    getHoldCount() {
      return holdCount;
    },
    subscribe(listener) {
      listeners.add(listener);
      return () => {
        listeners.delete(listener);
      };
    },
  };
}

const routeProgressController = createTopProgressController();

export function useTopProgressHold(): void {
  useEffect(() => {
    return routeProgressController.beginHold();
  }, []);
}

export function TopProgressBar() {
  const { pathname, search } = useLocation();
  const routeKey = `${pathname}${search}`;
  const pendingInitialQueries = useIsFetching({
    predicate: (query) => query.state.status === "pending",
  });

  const [phase, setPhase] = useState<TopProgressPhase>("idle");
  const [progress, setProgress] = useState(0);

  const prevRouteKeyRef = useRef<string | null>(null);
  const phaseRef = useRef<TopProgressPhase>("idle");
  const navLockRef = useRef(false);
  const pendingQueriesRef = useRef(pendingInitialQueries);
  pendingQueriesRef.current = pendingInitialQueries;

  const trickleTimerRef = useRef<number | null>(null);
  const minDelayTimerRef = useRef<number | null>(null);
  const completeTimerRef = useRef<number | null>(null);
  const fadeTimerRef = useRef<number | null>(null);
  const safetyTimerRef = useRef<number | null>(null);

  const setPhaseTracked = useCallback((next: TopProgressPhase) => {
    phaseRef.current = next;
    setPhase(next);
  }, []);

  const clearTimers = useCallback(() => {
    if (trickleTimerRef.current !== null) {
      window.clearInterval(trickleTimerRef.current);
      trickleTimerRef.current = null;
    }
    if (minDelayTimerRef.current !== null) {
      window.clearTimeout(minDelayTimerRef.current);
      minDelayTimerRef.current = null;
    }
    if (completeTimerRef.current !== null) {
      window.clearTimeout(completeTimerRef.current);
      completeTimerRef.current = null;
    }
    if (fadeTimerRef.current !== null) {
      window.clearTimeout(fadeTimerRef.current);
      fadeTimerRef.current = null;
    }
    if (safetyTimerRef.current !== null) {
      window.clearTimeout(safetyTimerRef.current);
      safetyTimerRef.current = null;
    }
  }, []);

  const finishProgress = useCallback(() => {
    if (phaseRef.current !== "running") {
      return;
    }
    if (trickleTimerRef.current !== null) {
      window.clearInterval(trickleTimerRef.current);
      trickleTimerRef.current = null;
    }
    if (safetyTimerRef.current !== null) {
      window.clearTimeout(safetyTimerRef.current);
      safetyTimerRef.current = null;
    }
    if (completeTimerRef.current !== null) {
      window.clearTimeout(completeTimerRef.current);
    }
    if (fadeTimerRef.current !== null) {
      window.clearTimeout(fadeTimerRef.current);
    }

    setProgress(100);
    setPhaseTracked("completing");

    completeTimerRef.current = window.setTimeout(() => {
      completeTimerRef.current = null;
      setPhaseTracked("fading");
      fadeTimerRef.current = window.setTimeout(() => {
        fadeTimerRef.current = null;
        setPhaseTracked("idle");
        setProgress(0);
      }, FADE_OUT_MS);
    }, COMPLETE_HOLD_MS);
  }, [setPhaseTracked]);

  const maybeFinish = useCallback(
    (timedOut = false) => {
      if (
        phaseRef.current === "running" &&
        shouldFinishTopProgress({
          navLocked: navLockRef.current,
          holdCount: routeProgressController.getHoldCount(),
          pendingQueries: pendingQueriesRef.current,
          timedOut,
        })
      ) {
        finishProgress();
      }
    },
    [finishProgress],
  );

  const startProgress = useCallback(() => {
    if (completeTimerRef.current !== null) {
      window.clearTimeout(completeTimerRef.current);
      completeTimerRef.current = null;
    }
    if (fadeTimerRef.current !== null) {
      window.clearTimeout(fadeTimerRef.current);
      fadeTimerRef.current = null;
    }

    setPhaseTracked("running");
    setProgress((prev) =>
      prev > TOP_PROGRESS_INITIAL && prev < 100 ? prev : TOP_PROGRESS_INITIAL,
    );

    if (trickleTimerRef.current === null) {
      trickleTimerRef.current = window.setInterval(() => {
        setProgress((prev) => computeNextTrickleProgress(prev));
      }, TRICKLE_INTERVAL_MS);
    }

    if (safetyTimerRef.current !== null) {
      window.clearTimeout(safetyTimerRef.current);
    }
    safetyTimerRef.current = window.setTimeout(() => {
      safetyTimerRef.current = null;
      navLockRef.current = false;
      maybeFinish(true);
    }, MAX_ACTIVE_MS);
  }, [maybeFinish, setPhaseTracked]);

  useEffect(() => {
    const handleHoldChange = (holdCount: number) => {
      if (holdCount > 0) {
        startProgress();
      } else {
        maybeFinish();
      }
    };

    if (routeProgressController.getHoldCount() > 0) {
      startProgress();
    }

    const unsubscribe = routeProgressController.subscribe(handleHoldChange);
    return () => {
      unsubscribe();
      clearTimers();
    };
  }, [clearTimers, maybeFinish, startProgress]);

  useEffect(() => {
    if (pendingInitialQueries === 0) {
      maybeFinish();
    }
  }, [maybeFinish, pendingInitialQueries]);

  useEffect(() => {
    if (prevRouteKeyRef.current === null) {
      prevRouteKeyRef.current = routeKey;
      return;
    }
    if (prevRouteKeyRef.current === routeKey) {
      return;
    }
    prevRouteKeyRef.current = routeKey;

    if (minDelayTimerRef.current !== null) {
      window.clearTimeout(minDelayTimerRef.current);
    }
    navLockRef.current = true;
    startProgress();

    minDelayTimerRef.current = window.setTimeout(() => {
      navLockRef.current = false;
      minDelayTimerRef.current = null;
      maybeFinish();
    }, MIN_NAV_VISIBLE_MS);
  }, [maybeFinish, routeKey, startProgress]);

  const isVisible = phase !== "idle";

  return (
    <div
      role="progressbar"
      aria-label="Page loading progress"
      aria-valuemin={0}
      aria-valuemax={100}
      aria-valuenow={Math.round(progress)}
      aria-hidden={!isVisible}
      className="pointer-events-none fixed inset-x-0 top-0 z-50 h-[2.5px] overflow-hidden"
    >
      <div
        className="top-progress-bar h-full w-full"
        style={getTopProgressBarStyle(phase, progress)}
      />
    </div>
  );
}
