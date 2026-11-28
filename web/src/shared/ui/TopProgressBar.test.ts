import { describe, expect, test } from "bun:test";
import {
  computeNextTrickleProgress,
  createTopProgressController,
  getTopProgressBarStyle,
  shouldFinishTopProgress,
  TOP_PROGRESS_INITIAL,
  TOP_PROGRESS_MAX_TRICKLE,
} from "./TopProgressBar";

describe("TopProgressBar state machine", () => {
  test("computeNextTrickleProgress advances monotonically and caps at TOP_PROGRESS_MAX_TRICKLE", () => {
    let progress = TOP_PROGRESS_INITIAL;
    for (let step = 0; step < 60; step += 1) {
      const next = computeNextTrickleProgress(progress);
      expect(next).toBeGreaterThanOrEqual(progress);
      expect(next).toBeLessThanOrEqual(TOP_PROGRESS_MAX_TRICKLE);
      progress = next;
    }
    expect(progress).toBeCloseTo(TOP_PROGRESS_MAX_TRICKLE, 0);
    expect(computeNextTrickleProgress(TOP_PROGRESS_MAX_TRICKLE)).toBe(TOP_PROGRESS_MAX_TRICKLE);
  });

  test("createTopProgressController tracks ref-counted holds and idempotent release", () => {
    const controller = createTopProgressController();
    const observed: number[] = [];
    const unsubscribe = controller.subscribe((count) => observed.push(count));

    expect(controller.getHoldCount()).toBe(0);

    const releaseFirst = controller.beginHold();
    expect(controller.getHoldCount()).toBe(1);

    const releaseSecond = controller.beginHold();
    expect(controller.getHoldCount()).toBe(2);

    // Calling releaseFirst twice must only decrement once.
    releaseFirst();
    releaseFirst();
    expect(controller.getHoldCount()).toBe(1);

    releaseSecond();
    expect(controller.getHoldCount()).toBe(0);

    unsubscribe();
    const releaseAfterUnsub = controller.beginHold();
    releaseAfterUnsub();

    expect(observed).toEqual([1, 2, 1, 0]);
  });

  test("shouldFinishTopProgress waits for nav lock, suspense holds, and initial pending queries unless timed out", () => {
    expect(
      shouldFinishTopProgress({
        navLocked: true,
        holdCount: 0,
        pendingQueries: 0,
        timedOut: false,
      }),
    ).toBe(false);
    expect(
      shouldFinishTopProgress({
        navLocked: false,
        holdCount: 1,
        pendingQueries: 0,
        timedOut: false,
      }),
    ).toBe(false);
    expect(
      shouldFinishTopProgress({
        navLocked: false,
        holdCount: 0,
        pendingQueries: 2,
        timedOut: false,
      }),
    ).toBe(false);
    expect(
      shouldFinishTopProgress({
        navLocked: false,
        holdCount: 0,
        pendingQueries: 0,
        timedOut: false,
      }),
    ).toBe(true);
    expect(
      shouldFinishTopProgress({ navLocked: true, holdCount: 1, pendingQueries: 3, timedOut: true }),
    ).toBe(true);
  });

  test("getTopProgressBarStyle keeps full width during fade-out and disables transition on idle reset", () => {
    expect(getTopProgressBarStyle("running", 42)).toEqual({
      transform: "scaleX(0.42)",
      opacity: 1,
    });
    expect(getTopProgressBarStyle("completing", 100)).toEqual({
      transform: "scaleX(1)",
      opacity: 1,
    });
    expect(getTopProgressBarStyle("fading", 100)).toEqual({
      transform: "scaleX(1)",
      opacity: 0,
    });
    expect(getTopProgressBarStyle("idle", 0)).toEqual({
      transform: "scaleX(0)",
      opacity: 0,
      transition: "none",
    });
  });
});
