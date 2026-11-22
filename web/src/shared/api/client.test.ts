import { describe, expect, test } from "bun:test";
import type { ClientOptions } from "openapi-fetch";
import { createApiClient, extractErrorMessage } from "./client";

type CustomFetch = NonNullable<ClientOptions["fetch"]>;

describe("extractErrorMessage", () => {
  test("reads structured Go ErrorEnvelope", () => {
    expect(
      extractErrorMessage({
        error: {
          code: "conflict",
          message: "Already initialized",
          request_id: "req-123",
        },
      }),
    ).toBe("Already initialized");
  });

  test("reads standard Error instances", () => {
    expect(extractErrorMessage(new Error("Network unreachable"))).toBe("Network unreachable");
  });

  test("reads non-empty string errors and plain message objects", () => {
    expect(extractErrorMessage("Gateway timeout")).toBe("Gateway timeout");
    expect(extractErrorMessage({ message: "Validation failed" })).toBe("Validation failed");
  });

  test("falls back cleanly for unknown or empty values", () => {
    expect(extractErrorMessage(null)).toBe("An unexpected error occurred");
    expect(extractErrorMessage(undefined)).toBe("An unexpected error occurred");
    expect(extractErrorMessage({ error: { code: "empty", message: "   " } })).toBe(
      "An unexpected error occurred",
    );
  });
});

describe("createApiClient", () => {
  test("parses /api/health response and enforces credentials include", async () => {
    let capturedCredentials: RequestCredentials | undefined;
    let capturedUrl = "";

    const customFetch: CustomFetch = async (input) => {
      capturedCredentials = input.credentials;
      capturedUrl = input.url;

      return new Response(
        JSON.stringify({
          status: "ok",
          db: "ok",
          modules: ["system"],
          uptimeSeconds: 5,
        }),
        {
          status: 200,
          headers: { "Content-Type": "application/json" },
        },
      );
    };

    const client = createApiClient({
      baseUrl: "http://127.0.0.1:20128",
      fetch: customFetch,
    });

    const { data, error } = await client.GET("/api/health");

    expect(error).toBeUndefined();
    expect(data).toBeDefined();
    expect(data?.status).toBe("ok");
    expect(data?.db).toBe("ok");
    expect(data?.modules).toContain("system");
    expect(data?.uptimeSeconds).toBe(5);
    expect(capturedCredentials).toBe("include");
    expect(capturedUrl).toBe("http://127.0.0.1:20128/api/health");
  });

  test("surfaces structured ErrorEnvelope on POST /api/init conflict", async () => {
    const customFetch: CustomFetch = async () =>
      new Response(
        JSON.stringify({
          error: {
            code: "conflict",
            message: "System is already initialized",
          },
        }),
        {
          status: 409,
          headers: { "Content-Type": "application/json" },
        },
      );

    const client = createApiClient({
      baseUrl: "http://127.0.0.1:20128",
      fetch: customFetch,
    });

    const { data, error } = await client.POST("/api/init", {
      body: { password: "correct-horse-battery" },
    });

    expect(data).toBeUndefined();
    expect(error).toBeDefined();
    expect(extractErrorMessage(error)).toBe("System is already initialized");
  });

  test("parses /api/version response", async () => {
    const customFetch: CustomFetch = async () =>
      new Response(
        JSON.stringify({
          version: "0.1.0",
          commit: "abc1234",
          buildTime: "2026-11-22T17:00:00Z",
          goVersion: "go1.26.2",
        }),
        {
          status: 200,
          headers: { "Content-Type": "application/json" },
        },
      );

    const client = createApiClient({
      baseUrl: "http://127.0.0.1:20128",
      fetch: customFetch,
    });

    const { data, error } = await client.GET("/api/version");

    expect(error).toBeUndefined();
    expect(data?.version).toBe("0.1.0");
    expect(data?.goVersion).toBe("go1.26.2");
  });
});
