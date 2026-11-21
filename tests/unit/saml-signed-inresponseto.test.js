import { describe, it, expect, vi } from "vitest";

const validated = vi.hoisted(() => ({ profile: null }));

vi.mock("@node-saml/node-saml", () => ({
  SAML: class {
    async validatePostResponseAsync() {
      return { profile: validated.profile };
    }
  },
}));

const { validateSamlResponse, getSamlBaseUrl } = await import("../../src/lib/auth/saml.js");

const settings = { samlCert: "dummy-cert", baseUrl: "https://router.example.com" };
const body = { SAMLResponse: Buffer.from('<Response InResponseTo="req-123"></Response>').toString("base64") };

describe("SAML InResponseTo from the signed document", () => {
  it("rejects when the raw XML matches but the signed response answers another request", async () => {
    validated.profile = { inResponseTo: "req-other", email: "a@example.com" };
    await expect(validateSamlResponse(null, body, "req-123", settings)).rejects.toThrow(/signed response/);
  });

  it("rejects when the signed response carries no InResponseTo", async () => {
    validated.profile = { email: "a@example.com" };
    await expect(validateSamlResponse(null, body, "req-123", settings)).rejects.toThrow(/signed response/);
  });

  it("accepts when the signed InResponseTo matches the stored request id", async () => {
    validated.profile = { inResponseTo: "req-123", email: "a@example.com" };
    await expect(validateSamlResponse(null, body, "req-123", settings)).resolves.toMatchObject({ email: "a@example.com" });
  });
});

describe("getSamlBaseUrl", () => {
  it("ignores forwarded-host headers when no base URL is configured", () => {
    const request = new Request("http://localhost:20128/api/auth/saml/start", {
      headers: { "x-forwarded-host": "evil.example", "x-forwarded-proto": "https" },
    });
    expect(getSamlBaseUrl(request, {})).toBe("http://localhost:20128");
  });
});
