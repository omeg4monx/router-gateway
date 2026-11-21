import { describe, it, expect, vi } from "vitest";
import { VeniceExecutor, getExecutor, hasSpecializedExecutor } from "../../open-sse/executors/index.js";
import REGISTRY from "../../open-sse/providers/registry/index.js";
import { FREE_PROVIDERS } from "../../src/shared/constants/providers.js";
import { parseModel, resolveProviderAlias } from "../../open-sse/services/model.js";
import { getModelsByProviderId } from "../../open-sse/config/providerModels.js";

describe("Venice.ai keyless provider (xrouter 3b0f6c3 port)", () => {
  it("registry exposes a valid no-auth free entry with six models", () => {
    const reg = REGISTRY.find((p) => p.id === "venice-free");
    expect(reg).toBeDefined();
    expect(reg.category).toBe("free");
    expect(reg.noAuth).toBe(true);
    expect(reg.authType).toBe("none");
    expect(reg.transport.format).toBe("venice");
    expect(reg.transport.baseUrl).toBe("https://outerface.venice.ai/api/inference/chat");
    expect(reg.models).toHaveLength(6);
    const ids = reg.models.map((m) => m.id);
    expect(ids).toContain("deepseek-v4-1-flash-uncensored");
    expect(ids).toContain("qwen3-coder-480b-a35b-instruct-turbo");
    expect(ids).toContain("nvidia-nemotron-3-ultra-550b-a55b");
  });

  it("does not shadow the existing venice API-key provider", () => {
    const apiKeyReg = REGISTRY.find((p) => p.id === "venice");
    expect(apiKeyReg).toBeDefined();
    expect(apiKeyReg.noAuth).not.toBe(true);
  });

  it("is exposed as noAuth in FREE_PROVIDERS so auth.js skips credentials", () => {
    expect(FREE_PROVIDERS["venice-free"]).toBeDefined();
    expect(FREE_PROVIDERS["venice-free"].noAuth).toBe(true);
    expect(FREE_PROVIDERS["venice-free"].authType).toBe("none");
  });

  it("static model catalog resolves for the models route", () => {
    const models = getModelsByProviderId("venice-free");
    expect(models).toHaveLength(6);
    expect(models.map((m) => m.id)).toContain("deepseek-v4-1-flash-uncensored");
  });

  it("routes venice model prefixes to the keyless executor", () => {
    expect(resolveProviderAlias("venice-free")).toBe("venice-free");
    expect(resolveProviderAlias("venice-web")).toBe("venice-free");
    expect(resolveProviderAlias("venice-keyless")).toBe("venice-free");
    const parsed = parseModel("venice-free/qwen3-coder-480b-a35b-instruct-turbo");
    expect(parsed.provider).toBe("venice-free");
    expect(getExecutor(parsed.provider)).toBeInstanceOf(VeniceExecutor);
  });

  it("registers the keyless executor under venice-free aliases", () => {
    expect(hasSpecializedExecutor("venice-free")).toBe(true);
    expect(hasSpecializedExecutor("venice-web")).toBe(true);
    expect(getExecutor("venice-free")).toBeInstanceOf(VeniceExecutor);
    expect(getExecutor("venice-free").noAuth).toBe(true);
  });

  it("formats an OpenAI request into the Venice outerface payload", () => {
    const executor = new VeniceExecutor();
    const payload = executor.formatPayload("venice/qwen3-coder-480b-a35b-instruct-turbo", {
      messages: [
        { role: "system", content: "You are an AI assistant." },
        { role: "user", content: "What is 2+2?" },
        { role: "assistant", content: "4" },
        { role: "user", content: "And plus 1?" },
      ],
      temperature: 0.5,
      top_p: 0.8,
    });
    expect(payload.modelId).toBe("qwen3-coder-480b-a35b-instruct-turbo");
    expect(payload.systemPrompt).toBe("You are an AI assistant.");
    expect(payload.temperature).toBe(0.5);
    expect(payload.topP).toBe(0.8);
    expect(payload.simpleMode).toBe(false);
    expect(payload.prompt).toHaveLength(3);
    expect(payload.prompt.map((m) => m.role)).toEqual(["user", "assistant", "user"]);
    expect(typeof payload.requestId).toBe("string");
  });

  it("resolves a connection proxy when enabled", () => {
    const executor = new VeniceExecutor();
    expect(executor.resolveProxy({
      connectionProxyEnabled: true,
      connectionProxyUrl: "http://1.2.3.4:8080",
    })).toBe("http://1.2.3.4:8080");
    expect(executor.resolveProxy({
      connectionProxyEnabled: false,
      connectionProxyUrl: "http://1.2.3.4:8080",
    })).toBeNull();
  });

  it("caches minted attestation tokens per proxy without re-minting", async () => {
    const executor = new VeniceExecutor();
    const session = {
      fetch: vi.fn().mockResolvedValue({
        status: 200,
        json: async () => ({ tokens: [{ token: "tok-a" }, { token: "tok-b" }] }),
      }),
    };
    expect(await executor.getAttestationToken(session, "proxy-1")).toBe("tok-a");
    expect(await executor.getAttestationToken(session, "proxy-1")).toBe("tok-b");
    expect(session.fetch).toHaveBeenCalledTimes(1);
  });

  it("converts Venice NDJSON content chunks into an OpenAI SSE stream", async () => {
    const executor = new VeniceExecutor();
    executor.resolveProxy = () => null;
    executor.getAttestationToken = vi.fn().mockResolvedValue("attestation-token");

    const ndjson = [
      JSON.stringify({ kind: "metadata", requestId: "x" }),
      JSON.stringify({ kind: "content", content: "Hello" }),
      JSON.stringify({ kind: "content", content: " world" }),
    ].join("\n");

    const session = {
      fetch: vi.fn().mockImplementation(async (url) => {
        if (String(url).includes("/mint")) {
          return { status: 200, json: async () => ({ token: "attestation-token" }) };
        }
        return {
          ok: true,
          status: 200,
          body: new ReadableStream({
            start(controller) {
              controller.enqueue(new TextEncoder().encode(ndjson));
              controller.close();
            },
          }),
        };
      }),
    };

    const wreqModule = await import("wreq-js");
    const origCreate = wreqModule.default?.createSession || wreqModule.createSession;
    // The executor captured wreq at import time via createRequire; patch that object.
    const wreqReq = (await import("node:module")).createRequire(import.meta.url)("wreq-js");
    const orig = wreqReq.createSession;
    wreqReq.createSession = vi.fn().mockResolvedValue(session);
    try {
      const result = await executor.execute({
        model: "venice/deepseek-v4-1-flash-uncensored",
        body: { messages: [{ role: "user", content: "hi" }] },
        stream: true,
      });
      expect(result.response.status).toBe(200);
      expect(result.response.headers.get("Content-Type")).toContain("text/event-stream");
      const text = await result.response.text();
      expect(text).toContain('"content":"Hello"');
      expect(text).toContain('"content":" world"');
      expect(text).toContain('"finish_reason":"stop"');
      expect(text).toContain("data: [DONE]");
    } finally {
      wreqReq.createSession = orig;
    }
  });
});
