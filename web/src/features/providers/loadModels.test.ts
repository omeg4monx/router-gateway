import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { sortProviderModels } from "./modelIntelligence";

const providersDir = resolve(import.meta.dir);
const hooksSource = readFileSync(resolve(providersDir, "hooks.ts"), "utf8");
const detailSource = readFileSync(resolve(providersDir, "ProviderDetailView.tsx"), "utf8");

describe("OmniRoute Load Models feature contracts", () => {
  test("hooks.ts exposes useFetchConnectionModelsMutation calling /api/providers/connections/{id}/models", () => {
    expect(hooksSource).toContain("useFetchConnectionModelsMutation");
    expect(hooksSource).toContain("/api/providers/connections/");
    expect(hooksSource).toContain("/models");
  });

  test("ProviderDetailView renders Load Models button and Save to Connection action in Available Models card", () => {
    expect(detailSource).toContain("useFetchConnectionModelsMutation");
    expect(detailSource).toContain('"Load Models"');
    expect(detailSource).toContain('"Muat Model"');
    expect(detailSource).toContain("Save to Connection");
  });

  test("ProviderDetailView auto-loads live models for any provider with an active connection", () => {
    expect(detailSource).toContain("autoLoadedConnRef");
    expect(detailSource).toContain("fetchModelsMutate(activeConnectionForLoad.id");
  });

  test("ProviderDetailView renders detected models in a proportional responsive grid without stretching sibling cards", () => {
    expect(detailSource).toContain(
      "grid grid-cols-1 items-start gap-2.5 sm:grid-cols-2 lg:grid-cols-3",
    );
    expect(detailSource).toContain("flex min-h-[96px] min-w-0 flex-col justify-between");
    expect(detailSource).not.toContain("flex h-full min-w-0 flex-col justify-between");
  });

  test("ProviderDetailView supports per-model live testing with AbortController cleanup", () => {
    expect(detailSource).toContain("runSingleModelTest");
    expect(detailSource).toContain("testingModelIds");
    expect(detailSource).toContain("modelTestResults");
    expect(detailSource).toContain("handleTestModel");
    expect(detailSource).toContain("modelTestControllersRef");
  });

  test("sortProviderModels sorts models naturally in asc and desc order without mutating input", () => {
    const input = [
      { id: "gpt-4o" },
      { id: "claude-3-5-sonnet" },
      { id: "gpt-10" },
      { id: "gpt-5" },
    ];
    const snapshot = input.map((m) => m.id);

    expect(sortProviderModels(input, "name", "asc").map((m) => m.id)).toEqual([
      "claude-3-5-sonnet",
      "gpt-4o",
      "gpt-5",
      "gpt-10",
    ]);
    expect(sortProviderModels(input, "name", "desc").map((m) => m.id)).toEqual([
      "gpt-10",
      "gpt-5",
      "gpt-4o",
      "claude-3-5-sonnet",
    ]);
    expect(input.map((m) => m.id)).toEqual(snapshot);
  });

  test("ProviderDetailView includes name/intelligence sort toggle and tier badges for Available Models", () => {
    expect(detailSource).toContain("modelSortField");
    expect(detailSource).toContain("modelSortDirection");
    expect(detailSource).toContain("sortProviderModels");
    expect(detailSource).toContain("getModelIntelligenceProfile");
    expect(detailSource).toContain("MODEL_TIER_BADGE");
    expect(detailSource).toContain("ArrowDownAZ");
    expect(detailSource).toContain("ArrowUpZA");
  });

  test("ProviderDetailView includes per-model Unmask identity probe with smooth expand/collapse and limitations", () => {
    const modelsHooksSource = readFileSync(resolve(providersDir, "../models/hooks.ts"), "utf8");
    expect(modelsHooksSource).toContain("runModelIdentityProbe");
    expect(modelsHooksSource).toContain("/api/models/identify");
    expect(detailSource).toContain("runModelIdentityProbe");
    expect(detailSource).toContain("handleIdentifyModel");
    expect(detailSource).toContain("identifyingModelIds");
    expect(detailSource).toContain("modelIdentityResults");
    expect(detailSource).toContain("openIdentityModelIds");
    expect(detailSource).toContain("modelIdentityControllersRef");
    expect(detailSource).toContain("limitationNote");
    expect(detailSource).toContain("grid-rows-[1fr]");
    expect(detailSource).toContain("grid-rows-[0fr]");
    expect(detailSource).toContain("motion-reduce:transition-none");
  });
});
