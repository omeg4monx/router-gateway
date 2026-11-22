import { normalizeOrigin } from "./endpointUtils";

/**
 * Generates ready-to-run client integration snippets for OpenAI-compatible
 * (/v1/chat/completions) and Anthropic-compatible (/v1/messages) endpoints.
 */
export function generateEndpointSnippet({
  format = "openai",
  lang = "curl",
  baseUrl = "",
  apiKey = "",
  model = "",
} = {}) {
  const origin = normalizeOrigin(baseUrl) || "http://localhost:20128";
  const key = (apiKey && String(apiKey).trim()) || "<YOUR_API_KEY>";
  const resolvedModel =
    (model && String(model).trim()) ||
    (format === "anthropic" ? "claude-sonnet-4-5" : "gpt-4o");

  if (format === "anthropic") {
    const url = `${origin}/v1/messages`;
    if (lang === "curl") {
      return `curl -X POST "${url}" \\
  -H "x-api-key: ${key}" \\
  -H "anthropic-version: 2023-06-01" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"${resolvedModel}","max_tokens":1024,"messages":[{"role":"user","content":"Hello"}]}'`;
    }
    if (lang === "python") {
      return `import anthropic

client = anthropic.Anthropic(
    base_url="${origin}",
    api_key="${key}",
)
msg = client.messages.create(
    model="${resolvedModel}",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Hello"}],
)
print(msg.content[0].text)`;
    }
    if (lang === "node") {
      return `import Anthropic from "@anthropic-ai/sdk";

const client = new Anthropic({
  baseURL: "${origin}",
  apiKey: "${key}",
});
const msg = await client.messages.create({
  model: "${resolvedModel}",
  max_tokens: 1024,
  messages: [{ role: "user", content: "Hello" }],
});
console.log(msg.content[0].text);`;
    }
    if (lang === "go") {
      return `package main

import (
\t"bytes"
\t"fmt"
\t"io"
\t"net/http"
)

func main() {
\tbody := []byte(\`{"model":"${resolvedModel}","max_tokens":1024,"messages":[{"role":"user","content":"Hello"}]}\`)
\treq, _ := http.NewRequest("POST", "${url}", bytes.NewBuffer(body))
\treq.Header.Set("x-api-key", "${key}")
\treq.Header.Set("anthropic-version", "2023-06-01")
\treq.Header.Set("Content-Type", "application/json")
\tresp, _ := http.DefaultClient.Do(req)
\tdefer resp.Body.Close()
\tb, _ := io.ReadAll(resp.Body)
\tfmt.Println(string(b))
}`;
    }
    return "";
  }

  const openaiBase = `${origin}/v1`;
  const url = `${openaiBase}/chat/completions`;

  if (lang === "curl") {
    return `curl -X POST "${url}" \\
  -H "Authorization: Bearer ${key}" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"${resolvedModel}","messages":[{"role":"user","content":"Hello"}]}'`;
  }
  if (lang === "python") {
    return `from openai import OpenAI

client = OpenAI(base_url="${openaiBase}", api_key="${key}")
res = client.chat.completions.create(
    model="${resolvedModel}",
    messages=[{"role": "user", "content": "Hello"}],
)
print(res.choices[0].message.content)`;
  }
  if (lang === "node") {
    return `import OpenAI from "openai";

const client = new OpenAI({ baseURL: "${openaiBase}", apiKey: "${key}" });
const res = await client.chat.completions.create({
  model: "${resolvedModel}",
  messages: [{ role: "user", content: "Hello" }],
});
console.log(res.choices[0].message.content);`;
  }
  if (lang === "go") {
    return `package main

import (
\t"bytes"
\t"fmt"
\t"io"
\t"net/http"
)

func main() {
\tbody := []byte(\`{"model":"${resolvedModel}","messages":[{"role":"user","content":"Hello"}]}\`)
\treq, _ := http.NewRequest("POST", "${url}", bytes.NewBuffer(body))
\treq.Header.Set("Authorization", "Bearer ${key}")
\treq.Header.Set("Content-Type", "application/json")
\tresp, _ := http.DefaultClient.Do(req)
\tdefer resp.Body.Close()
\tb, _ := io.ReadAll(resp.Body)
\tfmt.Println(string(b))
}`;
  }
  return "";
}
