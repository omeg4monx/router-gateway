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

  test("ProviderDetailView renders OmniRoute Connections card with priority controls and proportional Available Models grid", () => {
    expect(detailSource).toContain("Connections");
    expect(detailSource).toContain("Available Models");
    expect(detailSource).toContain(
      "grid grid-cols-1 items-start gap-2.5 sm:grid-cols-2 lg:grid-cols-3",
    );
    expect(detailSource).toContain("rounded-lg border border-[var(--color-border)]");
    expect(detailSource).toContain("AddApiKeyModal");
    expect(detailSource).toContain("OAuthWizardModal");
    expect(detailSource).toContain("EditConnectionModal");
  });

  test("ProviderDetailView edits a selected connection in a modal without navigating away", () => {
    expect(modalsSource).toContain("export function EditConnectionModal(");
    expect(detailSource).toContain("setEditingConnection(conn)");
    expect(detailSource).not.toContain("/dashboard/providers/detail/?id=");
  });

  test("EditConnectionModal provides dimmed placeholders and prevents browser autofill highlight on secret rotation fields", () => {
    expect(modalsSource).toContain("placeholder:text-[var(--color-text-muted)]/50");
    expect(modalsSource).toContain("[&:-webkit-autofill]:shadow-");
    expect(modalsSource).toContain('autoComplete="off"');
    expect(modalsSource).toContain("secretFieldsUnlocked");
    expect(modalsSource).toContain('placeholder="https://api.openai.com/v1"');
  });
});

describe("OmniRoute AddApiKeyModal & OAuthWizardModal flow & styling", () => {
  test("AddApiKeyModal supports locked provider mode, inline Check validation button, and opens Advanced Settings by default", () => {
    expect(modalsSource).toContain("lockProvider?: boolean");
    expect(modalsSource).toContain('"Checking..."');
    expect(modalsSource).toContain('"Check"');
    const addApiKeySection = modalsSource.slice(
      modalsSource.indexOf("export function AddApiKeyModal("),
      modalsSource.indexOf("export function OAuthWizardModal("),
    );
    expect(addApiKeySection).toContain("const [showAdvanced, setShowAdvanced] = useState(true);");
    expect(addApiKeySection).toContain("setShowAdvanced(true);");
  });

  test("OAuthWizardModal implements OmniRoute combined waiting + manual paste divider and success/error screens", () => {
    expect(modalsSource).toContain("Or paste callback URL manually");
    expect(modalsSource).toContain("Waiting for popup authorization");
    expect(modalsSource).toContain("Connected Successfully!");
    expect(modalsSource).toContain("Connection Failed");
    expect(modalsSource).toContain("Try Again");
  });
});
