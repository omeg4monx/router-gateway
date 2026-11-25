import { describe, expect, test } from "bun:test";
import { ErrorBoundary, isChunkLoadError } from "./ErrorBoundary";

describe("ErrorBoundary", () => {
  test("getDerivedStateFromError captures thrown error", () => {
    const err = new Error("render boom");
    const next = ErrorBoundary.getDerivedStateFromError(err);
    expect(next).toEqual({ error: err });
  });

  test("isChunkLoadError detects Vite dynamic import and chunk failures", () => {
    expect(
      isChunkLoadError(
        new TypeError("Failed to fetch dynamically imported module: /assets/Foo.js"),
      ),
    ).toBe(true);
    expect(isChunkLoadError(new Error("Importing a module script failed."))).toBe(true);
    expect(isChunkLoadError(new Error("Loading chunk 42 failed"))).toBe(true);
    expect(isChunkLoadError(new Error("Ordinary runtime error"))).toBe(false);
    expect(isChunkLoadError(null)).toBe(false);
  });
});
