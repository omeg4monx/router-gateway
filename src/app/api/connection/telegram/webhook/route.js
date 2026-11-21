import crypto from "node:crypto";
import { NextResponse } from "next/server";
import { getConnectionConfig } from "@/lib/db/repos/connectionRepo";
import { POST as runCommand } from "../command/route.js";

export const dynamic = "force-dynamic";

function safeEqual(a, b) {
  const x = Buffer.from(String(a || ""));
  const y = Buffer.from(String(b || ""));
  return x.length > 0 && x.length === y.length && crypto.timingSafeEqual(x, y);
}

/**
 * Telegram Webhook Handler
 * Accepts only updates carrying the secret registered via setWebhook, from the
 * configured userId, then runs the command in-process and replies.
 */
export async function POST(request) {
  try {
    const config = await getConnectionConfig();
    const tg = config.telegram || {};

    if (!tg.enabled || !tg.botToken) {
      return NextResponse.json({ ok: false, error: "Telegram bot not enabled" });
    }

    // Proves the update came from Telegram. Webhooks registered before this check
    // have no secret and must be re-registered from the dashboard.
    if (!safeEqual(request.headers.get("x-telegram-bot-api-secret-token"), tg.webhookSecret)) {
      return NextResponse.json({ ok: false, error: "Unauthorized" }, { status: 401 });
    }

    // Fail closed: without a configured admin id any Telegram user could drive the bot.
    const configuredId = String(tg.userId || "").trim();
    if (!configuredId) {
      return NextResponse.json({ ok: false, error: "Telegram userId not configured" }, { status: 403 });
    }

    const update = await request.json().catch(() => null);
    if (!update || !update.message) {
      return NextResponse.json({ ok: true, ignored: true });
    }

    const { message } = update;
    const chatId = message.chat?.id;
    const fromId = message.from?.id;
    const text = (message.text || "").trim();

    if (!chatId || !text) {
      return NextResponse.json({ ok: true });
    }

    // Only the configured admin; chat id covers a configured private/group chat.
    if (String(fromId) !== configuredId && String(chatId) !== configuredId) {
      await fetch(`https://api.telegram.org/bot${tg.botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          chat_id: chatId,
          text: "⛔ *Akses Ditolak*\nAnda tidak memiliki otorisasi untuk mengontrol X Router ini.",
          parse_mode: "Markdown",
        }),
      }).catch(() => {});
      return NextResponse.json({ ok: true, authorized: false });
    }

    const cmdRes = await runCommand(
      new Request("http://internal/api/connection/telegram/command", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ command: text }),
      })
    );

    let replyText = "⚠️ Gagal memproses perintah.";
    if (cmdRes.ok) {
      const resData = await cmdRes.json();
      replyText = resData.text || "Perintah dijalankan.";
    }

    // Send response back to Telegram
    await fetch(`https://api.telegram.org/bot${tg.botToken}/sendMessage`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        chat_id: chatId,
        text: replyText,
        parse_mode: "Markdown",
      }),
    });

    return NextResponse.json({ ok: true });
  } catch (err) {
    console.error("[Telegram Webhook] Error handling update:", err);
    return NextResponse.json({ ok: false, error: err.message }, { status: 500 });
  }
}
