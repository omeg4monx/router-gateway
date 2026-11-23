import { describe, expect, test } from "bun:test";
import {
  telegramCommandFormSchema,
  telegramConfigFormSchema,
  telegramTestFormSchema,
} from "./schemas";

describe("telegram schemas", () => {
  test("accepts valid telegram config form", () => {
    const res = telegramConfigFormSchema.safeParse({
      enabled: true,
      bot_token: "123456:AAH-valid-token",
      chat_id: "10001",
      user_id: "777001",
      webhook_url: "https://router.example.com/api/telegram/webhook",
      notify_on_job_failure: true,
      notify_on_auth_change: true,
      notify_on_emit: false,
    });
    expect(res.success).toBe(true);
  });

  test("rejects non-https webhook_url when non-empty", () => {
    const res = telegramConfigFormSchema.safeParse({
      enabled: true,
      bot_token: "",
      chat_id: "10001",
      user_id: "777001",
      webhook_url: "http://insecure.example.com/webhook",
      notify_on_job_failure: true,
      notify_on_auth_change: false,
      notify_on_emit: false,
    });
    expect(res.success).toBe(false);
  });

  test("validates command simulator and test message forms", () => {
    expect(telegramCommandFormSchema.safeParse({ command: "/status" }).success).toBe(true);
    expect(telegramCommandFormSchema.safeParse({ command: "   " }).success).toBe(false);
    expect(telegramTestFormSchema.safeParse({ message: "Probe" }).success).toBe(true);
  });
});
