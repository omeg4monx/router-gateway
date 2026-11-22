import { describe, expect, test } from "bun:test";
import { cn } from "./utils";

describe("cn utility", () => {
  test("merges conditional classes and resolves conflicting Tailwind utilities", () => {
    expect(cn("px-2 py-1 text-sm", false && "hidden", "px-4")).toBe(
      "py-1 text-sm px-4",
    );
  });

  test("handles empty and undefined inputs cleanly", () => {
    expect(cn(undefined, null, "", "rounded-md")).toBe("rounded-md");
  });
});
