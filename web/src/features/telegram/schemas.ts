import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type TelegramConfig = components["schemas"]["TelegramConfig"];
export type UpdateTelegramConfigRequest = components["schemas"]["UpdateTelegramConfigRequest"];
export type TelegramActionResponse = components["schemas"]["TelegramActionResponse"];
export type TelegramTestRequest = components["schemas"]["TelegramTestRequest"];
export type TelegramTestResponse = components["schemas"]["TelegramTestResponse"];
export type TelegramCommandRequest = components["schemas"]["TelegramCommandRequest"];
export type TelegramCommandResponse = components["schemas"]["TelegramCommandResponse"];
export type TelegramEvent = components["schemas"]["TelegramEvent"];
export type TelegramEventListResponse = components["schemas"]["TelegramEventListResponse"];

export const telegramConfigFormSchema = z.object({
  enabled: z.boolean().default(false),
  bot_token: z.string().trim().optional().default(""),
  chat_id: z.string().trim().default(""),
  user_id: z.string().trim().default(""),
  webhook_url: z
    .string()
    .trim()
    .default("")
    .refine(
      (val) => val === "" || /^https:\/\/.+/i.test(val),
      "Webhook URL must start with https://",
    ),
  notify_on_job_failure: z.boolean().default(true),
  notify_on_auth_change: z.boolean().default(true),
  notify_on_emit: z.boolean().default(false),
});

export type TelegramConfigFormInput = z.input<typeof telegramConfigFormSchema>;
export type TelegramConfigFormOutput = z.output<typeof telegramConfigFormSchema>;

export const telegramCommandFormSchema = z.object({
  command: z.string().trim().min(1, "Command is required"),
});

export type TelegramCommandFormInput = z.input<typeof telegramCommandFormSchema>;
export type TelegramCommandFormOutput = z.output<typeof telegramCommandFormSchema>;

export const telegramTestFormSchema = z.object({
  message: z.string().trim().max(4000, "Message is too long").default(""),
});

export type TelegramTestFormInput = z.input<typeof telegramTestFormSchema>;
export type TelegramTestFormOutput = z.output<typeof telegramTestFormSchema>;
