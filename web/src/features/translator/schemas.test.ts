import { describe, expect, it } from "bun:test";
import { parseJsonObject, SAMPLE_PAYLOADS, translateFormSchema } from "./schemas";

describe("translator schemas", () => {
  it("validates request and response translation forms with JSON object payloads", () => {
    const ok = translateFormSchema.safeParse({
      mode: "request",
      from_protocol: "openai-chat",
      to_protocol: "anthropic",
      payloadText: SAMPLE_PAYLOADS.request["openai-chat"],
    });
    expect(ok.success).toBe(true);

    const okResp = translateFormSchema.safeParse({
      mode: "response",
      from_protocol: "anthropic",
      to_protocol: "openai-chat",
      payloadText: SAMPLE_PAYLOADS.response.anthropic,
    });
    expect(okResp.success).toBe(true);
  });

  it("rejects malformed JSON, non-object JSON, and unsupported protocols or modes", () => {
    expect(parseJsonObject("{bad json")).toBeNull();
    expect(parseJsonObject("[1, 2, 3]")).toBeNull();
    expect(parseJsonObject('"hello"')).toBeNull();

    const badPayload = translateFormSchema.safeParse({
      mode: "request",
      from_protocol: "openai-chat",
      to_protocol: "anthropic",
      payloadText: "[1,2,3]",
    });
    expect(badPayload.success).toBe(false);

    const badMode = translateFormSchema.safeParse({
      mode: "stream",
      from_protocol: "openai-chat",
      to_protocol: "anthropic",
      payloadText: "{}",
    });
    expect(badMode.success).toBe(false);
  });
});
