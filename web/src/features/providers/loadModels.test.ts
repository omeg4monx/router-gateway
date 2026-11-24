import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

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
});
