import { describe, expect, test } from "bun:test";
import {
  mergeManualModels,
  normalizeManualModelInput,
  removeModelFromList,
} from "./customModelUtils";

describe("customModelUtils", () => {
  test("splits comma and newline lists, trims, deduplicates, and strips provider prefix", () => {
    const result = normalizeManualModelInput(
      " omeg4mon/gpt-5 , claude-sonnet-4\nmy-alias/deepseek-r1\nopenai/gpt-4o\ngpt-5 ",
      ["omeg4mon", "my-alias"],
    );
    expect(result).toEqual(["gpt-5", "claude-sonnet-4", "deepseek-r1", "openai/gpt-4o"]);
  });

  test("ignores blank entries and lone provider prefixes", () => {
    expect(normalizeManualModelInput(" , \n omeg4mon/ , ", ["omeg4mon"])).toEqual([]);
  });

  test("mergeManualModels preserves existing models, avoids duplicates, and reports counts", () => {
    const res = mergeManualModels(["gpt-4o", "claude-3-7-sonnet"], ["gpt-4o", "o3-pro", "o3-pro"]);
    expect(res).toEqual({
      merged: ["gpt-4o", "claude-3-7-sonnet", "o3-pro"],
      added: ["o3-pro"],
      duplicates: ["gpt-4o"],
    });
  });

  test("removeModelFromList removes exact matching model while keeping order", () => {
    expect(removeModelFromList(["a", "b", "c"], "b")).toEqual(["a", "c"]);
    expect(removeModelFromList(["a", "b"], "missing")).toEqual(["a", "b"]);
  });
});
