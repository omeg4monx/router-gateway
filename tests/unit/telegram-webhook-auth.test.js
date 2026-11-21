import { describe, it, expect, vi, beforeEach } from "vitest";

const mocks = vi.hoisted(() => ({
  getConnectionConfig: vi.fn(),
  runCommand: vi.fn(),
  fetch: vi.fn(),
}));

vi.mock("@/lib/db/repos/connectionRepo", () => ({ getConnectionConfig: mocks.getConnectionConfig }));
vi.mock("../../src/app/api/connection/telegram/command/route.js", () => ({ POST: mocks.runCommand }));

const { POST } = await import("../../src/app/api/connection/telegram/webhook/route.js");

const SECRET = "s".repeat(64);

function update(headers = {}, fromId = 42) {
  return new Request("http://x/api/connection/telegram/webhook", {
    method: "POST",
    headers: { "content-type": "application/json", ...headers },
    body: JSON.stringify({ message: { chat: { id: fromId }, from: { id: fromId }, text: "/prune" } }),
  });
}

describe("telegram webhook auth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    globalThis.fetch = mocks.fetch.mockResolvedValue(new Response("{}"));
    mocks.runCommand.mockResolvedValue(Response.json({ text: "ok" }));
    mocks.getConnectionConfig.mockResolvedValue({
      telegram: { enabled: true, botToken: "123:abc", userId: "42", webhookSecret: SECRET },
    });
  });

  it("rejects forged updates without Telegram's secret header, even from the admin id", async () => {
    const res = await POST(update());
    expect(res.status).toBe(401);
    expect(mocks.runCommand).not.toHaveBeenCalled();
  });

  it("rejects updates when no webhook secret has been registered", async () => {
    mocks.getConnectionConfig.mockResolvedValue({
      telegram: { enabled: true, botToken: "123:abc", userId: "42" },
    });
    const res = await POST(update({ "x-telegram-bot-api-secret-token": "" }));
    expect(res.status).toBe(401);
    expect(mocks.runCommand).not.toHaveBeenCalled();
  });

  it("fails closed when no admin userId is configured", async () => {
    mocks.getConnectionConfig.mockResolvedValue({
      telegram: { enabled: true, botToken: "123:abc", userId: "", webhookSecret: SECRET },
    });
    const res = await POST(update({ "x-telegram-bot-api-secret-token": SECRET }));
    expect(res.status).toBe(403);
    expect(mocks.runCommand).not.toHaveBeenCalled();
  });

  it("ignores other Telegram users", async () => {
    await POST(update({ "x-telegram-bot-api-secret-token": SECRET }, 7));
    expect(mocks.runCommand).not.toHaveBeenCalled();
  });

  it("runs the command in-process for the admin with a valid secret", async () => {
    const res = await POST(update({ "x-telegram-bot-api-secret-token": SECRET }));
    expect(res.status).toBe(200);
    expect(mocks.runCommand).toHaveBeenCalledTimes(1);
    // No HTTP hop carrying the bot token: only the sendMessage reply goes out.
    expect(mocks.fetch).toHaveBeenCalledTimes(1);
    expect(String(mocks.fetch.mock.calls[0][0])).toContain("api.telegram.org");
  });
});
