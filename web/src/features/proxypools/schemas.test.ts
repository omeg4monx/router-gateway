import { describe, expect, it } from "bun:test";
import {
  createProxyPoolSchema,
  normalizeProxyInput,
  parseCandidateLines,
  scanCandidatesSchema,
} from "./schemas";

describe("proxypools schemas", () => {
  it("accepts valid http, https, socks5, and socks5h proxy URLs", () => {
    for (const url of [
      "http://127.0.0.1:8080",
      "https://user:pass@proxy.corp.local:8443",
      "socks5://127.0.0.1:1080",
      "socks5h://proxy.internal:1080",
    ]) {
      const parsed = createProxyPoolSchema.safeParse({ name: "Edge Pool", url, enabled: true });
      expect(parsed.success).toBe(true);
    }
  });

  it("rejects empty names, unsupported schemes, and missing hosts", () => {
    expect(
      createProxyPoolSchema.safeParse({ name: "  ", url: "http://127.0.0.1:8080" }).success,
    ).toBe(false);
    expect(
      createProxyPoolSchema.safeParse({ name: "Bad", url: "ftp://127.0.0.1:21" }).success,
    ).toBe(false);
    expect(createProxyPoolSchema.safeParse({ name: "Bad", url: "http://:8080" }).success).toBe(
      false,
    );
  });

  it("parses multiline candidates and caps scan list at 50", () => {
    const lines = parseCandidateLines(
      "http://127.0.0.1:8080\n\n  socks5://127.0.0.1:1080  \nhttp://127.0.0.1:8080",
    );
    expect(lines).toEqual(["http://127.0.0.1:8080", "socks5://127.0.0.1:1080"]);

    const okScan = scanCandidatesSchema.safeParse({
      candidatesText: "http://127.0.0.1:8080\nsocks5://127.0.0.1:1080",
      autoAdd: true,
      namePrefix: "Scanned",
      timeoutMs: 2500,
    });
    expect(okScan.success).toBe(true);

    const tooMany = Array.from({ length: 51 }, (_, i) => `http://127.0.0.1:${8000 + i}`).join("\n");
    expect(
      scanCandidatesSchema.safeParse({ candidatesText: tooMany, autoAdd: false }).success,
    ).toBe(false);
  });

  it("normalizes Webshare/ProxyScrape shorthand host:port and host:port:user:pass formats", () => {
    expect(normalizeProxyInput("198.23.239.134:6540:alice:secret")).toBe(
      "http://alice:secret@198.23.239.134:6540",
    );
    expect(normalizeProxyInput("198.23.239.134:6540")).toBe("http://198.23.239.134:6540");
    expect(normalizeProxyInput("socks5h://user:pass@10.0.0.1:1080")).toBe(
      "socks5h://user:pass@10.0.0.1:1080",
    );
    expect(
      parseCandidateLines(
        "198.23.239.134:6540:alice:secret\nhttp://alice:secret@198.23.239.134:6540",
      ),
    ).toEqual(["http://alice:secret@198.23.239.134:6540"]);
  });
});
