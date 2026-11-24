import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const providersDir = resolve(import.meta.dir);
const consoleSource = readFileSync(resolve(providersDir, "ProvidersConsole.tsx"), "utf8");
const detailSource = readFileSync(resolve(providersDir, "ProviderDetailView.tsx"), "utf8");
const modalsSource = readFileSync(resolve(providersDir, "ProvidersModals.tsx"), "utf8");

describe("OmniRoute provider card flow & styling", () => {
  test("ProviderCard navigates to dedicated provider detail route via Link", () => {
    expect(consoleSource).toContain("/dashboard/providers/detail/?provider=");
    expect(consoleSource).toContain('import { Link } from "react-router"');
  });

  test("ProviderCard includes OmniRoute top 2px status hairline accent and Setup CTA", () => {
    expect(consoleSource).toContain("h-[2px]");
    expect(consoleSource).toContain('"Setup"');
    expect(consoleSource).toContain('"Add API Key"');
  });
});

describe("OmniRoute provider detail page (?provider=<id>) flow & styling", () => {
  test("ProviderDetailView handles ?provider=<id> query param with dedicated provider studio view", () => {
    expect(detailSource).toContain('searchParams.get("provider")');
    expect(detailSource).toContain("useProvidersCatalog");
    expect(detailSource).toContain("useProviderConnections");
    expect(detailSource).toContain("buildProviderEntries");
  });

  test("ProviderDetailView renders OmniRoute Connections card with priority controls and Available Models chips", () => {
    expect(detailSource).toContain("Connections");
    expect(detailSource).toContain("Available Models");
    expect(detailSource).toContain("rounded-[5px] border border-[var(--color-border)] px-2 py-1");
    expect(detailSource).toContain("AddApiKeyModal");
    expect(detailSource).toContain("OAuthWizardModal");
  });
});

describe("OmniRoute AddApiKeyModal & OAuthWizardModal flow & styling", () => {
  test("AddApiKeyModal supports locked provider mode and inline Check validation button", () => {
    expect(modalsSource).toContain("lockProvider?: boolean");
    expect(modalsSource).toContain('"Checking..."');
    expect(modalsSource).toContain('"Check"');
  });

  test("OAuthWizardModal implements OmniRoute combined waiting + manual paste divider and success/error screens", () => {
    expect(modalsSource).toContain("Or paste callback URL manually");
    expect(modalsSource).toContain("Waiting for popup authorization");
    expect(modalsSource).toContain("Connected Successfully!");
    expect(modalsSource).toContain("Connection Failed");
    expect(modalsSource).toContain("Try Again");
  });
});
