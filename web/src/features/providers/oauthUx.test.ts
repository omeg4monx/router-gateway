import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const providersDir = resolve(import.meta.dir);
const modalsSource = readFileSync(resolve(providersDir, "ProvidersModals.tsx"), "utf8");

describe("OAuthWizardModal UX-friendly contracts", () => {
  test("listens to popup auto-callback via postMessage, BroadcastChannel, and localStorage", () => {
    expect(modalsSource).toContain('event.data?.type === "oauth_callback"');
    expect(modalsSource).toContain('new BroadcastChannel("oauth_callback")');
    expect(modalsSource).toContain('event.key === "oauth_callback"');
  });

  test("provides one-click clipboard paste for manual callback URL", () => {
    expect(modalsSource).toContain("pasteCallbackFromClipboard");
    expect(modalsSource).toContain('"Tempel"');
    expect(modalsSource).toContain('"Paste"');
  });

  test("provides one-click Copy Code & Open action for Device Code flow", () => {
    expect(modalsSource).toContain("Copy Code & Open");
    expect(modalsSource).toContain("Salin Kode & Buka");
  });

  test("follows OmniRoute direct-open minimal UX without provider/label/stepper clutter when locked to provider", () => {
    expect(modalsSource).toContain("lockProvider?: boolean");
    expect(modalsSource).toContain("Opening sign-in window");
    expect(modalsSource).toContain("Membuka halaman login");
  });

  test("uses dynamic window.location.origin redirect_uri for PKCE and polls loopback status for any fixed_port provider", () => {
    expect(modalsSource).toContain("/api/providers/oauth/callback");
    expect(modalsSource).toContain("window.location.origin");
    expect(modalsSource).toContain("activeSession?.session.fixed_port");
    expect(modalsSource).toContain("fetchProviderOAuthStatus(loopbackProvider");
  });

  test("never renders manual OAuth Client ID or Client Secret input forms (OmniRoute/xrouter zero-config parity)", () => {
    expect(modalsSource).not.toContain("modal-oauth-client-id");
    expect(modalsSource).not.toContain("Masukkan OAuth Client ID");
    expect(modalsSource).not.toContain("Enter OAuth application Client ID");
    expect(modalsSource).not.toContain("Opsi OAuth Lanjutan");
  });

  test("auto-starts and auto-opens the Amazon Quick sign-in popup immediately on modal open", () => {
    expect(modalsSource).toContain("startAwsQuickFlow(true)");
    expect(modalsSource).not.toContain("startAwsQuickFlow(false)");
  });
});
