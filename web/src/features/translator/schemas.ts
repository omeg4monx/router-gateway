import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type TranslatorProtocol = components["schemas"]["TranslatorProtocol"];
export type TranslatorMode = components["schemas"]["TranslatorMode"];
export type TranslatorProtocolsResponse = components["schemas"]["TranslatorProtocolsResponse"];
export type TranslatePayloadRequest = components["schemas"]["TranslatePayloadRequest"];
export type TranslatePayloadResponse = components["schemas"]["TranslatePayloadResponse"];
export type ConsoleLogEntry = components["schemas"]["ConsoleLogEntry"];

export const TRANSLATOR_PROTOCOLS: readonly TranslatorProtocol[] = [
  "openai-chat",
  "anthropic",
  "openai-responses",
  "gemini",
] as const;

export const TRANSLATOR_MODES: readonly TranslatorMode[] = ["request", "response"] as const;

export const SAMPLE_PAYLOADS: Record<TranslatorMode, Record<TranslatorProtocol, string>> = {
  request: {
    "openai-chat": JSON.stringify(
      {
        model: "claude-sonnet-4",
        messages: [
          { role: "system", content: "You are a concise engineering assistant." },
          { role: "user", content: "Summarize HTTP/2 multiplexing in one sentence." },
        ],
        temperature: 0.2,
      },
      null,
      2,
    ),
    anthropic: JSON.stringify(
      {
        model: "claude-sonnet-4-20250514",
        max_tokens: 256,
        system: "You are a concise engineering assistant.",
        messages: [{ role: "user", content: "Explain WAL mode in SQLite." }],
      },
      null,
      2,
    ),
    "openai-responses": JSON.stringify(
      {
        model: "gpt-4.1",
        instructions: "Be direct.",
        input: "List three Go concurrency primitives.",
      },
      null,
      2,
    ),
    gemini: JSON.stringify(
      {
        contents: [{ role: "user", parts: [{ text: "Compare TCP and QUIC briefly." }] }],
        generationConfig: { maxOutputTokens: 256, temperature: 0.3 },
      },
      null,
      2,
    ),
  },
  response: {
    "openai-chat": JSON.stringify(
      {
        id: "chatcmpl-demo",
        object: "chat.completion",
        model: "gpt-4o",
        choices: [
          {
            index: 0,
            message: {
              role: "assistant",
              content: "HTTP/2 multiplexes streams over one TCP connection.",
            },
            finish_reason: "stop",
          },
        ],
        usage: { prompt_tokens: 18, completion_tokens: 11, total_tokens: 29 },
      },
      null,
      2,
    ),
    anthropic: JSON.stringify(
      {
        id: "msg_demo",
        type: "message",
        role: "assistant",
        model: "claude-sonnet-4-20250514",
        content: [
          {
            type: "text",
            text: "SQLite WAL appends writes to a separate log for concurrent reads.",
          },
        ],
        stop_reason: "end_turn",
        usage: { input_tokens: 14, output_tokens: 13 },
      },
      null,
      2,
    ),
    "openai-responses": JSON.stringify(
      {
        id: "resp_demo",
        object: "response",
        model: "gpt-4.1",
        status: "completed",
        output: [
          {
            type: "message",
            role: "assistant",
            content: [{ type: "output_text", text: "Goroutines, channels, and sync.Mutex." }],
          },
        ],
        usage: { input_tokens: 10, output_tokens: 8, total_tokens: 18 },
      },
      null,
      2,
    ),
    gemini: JSON.stringify(
      {
        candidates: [
          {
            content: {
              role: "model",
              parts: [
                {
                  text: "QUIC runs over UDP with built-in TLS 1.3 and zero head-of-line blocking.",
                },
              ],
            },
            finishReason: "STOP",
          },
        ],
        usageMetadata: { promptTokenCount: 12, candidatesTokenCount: 16, totalTokenCount: 28 },
      },
      null,
      2,
    ),
  },
};

export function parseJsonObject(raw: string): Record<string, unknown> | null {
  try {
    const parsed: unknown = JSON.parse(raw);
    if (parsed !== null && typeof parsed === "object" && !Array.isArray(parsed)) {
      return parsed as Record<string, unknown>;
    }
    return null;
  } catch {
    return null;
  }
}

export const translateFormSchema = z.object({
  mode: z.enum(["request", "response"]),
  from_protocol: z.enum(["openai-chat", "anthropic", "openai-responses", "gemini"]),
  to_protocol: z.enum(["openai-chat", "anthropic", "openai-responses", "gemini"]),
  payloadText: z
    .string()
    .refine((v) => parseJsonObject(v) !== null, "Payload must be a valid JSON object"),
});

export type TranslateFormValues = z.infer<typeof translateFormSchema>;
