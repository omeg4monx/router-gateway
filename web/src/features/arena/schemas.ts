import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type ArenaModelOption = components["schemas"]["ArenaModelOption"];
export type ArenaModelListResponse = components["schemas"]["ArenaModelListResponse"];
export type ArenaContenderResult = components["schemas"]["ArenaContenderResult"];
export type ArenaBattle = components["schemas"]["ArenaBattle"];
export type ArenaBattleListResponse = components["schemas"]["ArenaBattleListResponse"];
export type RunArenaBattleRequest = components["schemas"]["RunArenaBattleRequest"];
export type VoteArenaBattleRequest = components["schemas"]["VoteArenaBattleRequest"];
export type ArenaLeaderboardEntry = components["schemas"]["ArenaLeaderboardEntry"];
export type ArenaLeaderboardResponse = components["schemas"]["ArenaLeaderboardResponse"];
export type ChatSession = components["schemas"]["ChatSession"];
export type ChatSessionListResponse = components["schemas"]["ChatSessionListResponse"];
export type CreateChatSessionRequest = components["schemas"]["CreateChatSessionRequest"];
export type UpdateChatSessionRequest = components["schemas"]["UpdateChatSessionRequest"];
export type ChatMessage = components["schemas"]["ChatMessage"];
export type ChatSessionDetailResponse = components["schemas"]["ChatSessionDetailResponse"];
export type SendChatMessageRequest = components["schemas"]["SendChatMessageRequest"];
export type SendChatMessageResponse = components["schemas"]["SendChatMessageResponse"];

export const arenaBattleFormSchema = z.object({
  prompt: z.string().trim().min(1, "Prompt is required"),
  system_prompt: z.string().trim().optional().default(""),
  models: z
    .array(z.string().trim().min(1))
    .min(2, "Select at least 2 models")
    .max(4, "Select at most 4 models"),
  temperature: z.coerce.number().min(0).max(2).optional(),
  max_tokens: z.coerce.number().int().min(0).max(128000).optional(),
});

export type ArenaBattleFormInput = z.input<typeof arenaBattleFormSchema>;
export type ArenaBattleFormOutput = z.output<typeof arenaBattleFormSchema>;

export const createChatSessionFormSchema = z.object({
  title: z.string().trim().max(120).optional().default(""),
  model: z.string().trim().min(1, "Model is required"),
  system_prompt: z.string().trim().optional().default(""),
});

export type CreateChatSessionFormInput = z.input<typeof createChatSessionFormSchema>;
export type CreateChatSessionFormOutput = z.output<typeof createChatSessionFormSchema>;

export const chatMessageFormSchema = z.object({
  content: z.string().trim().min(1, "Message cannot be empty"),
});

export type ChatMessageFormInput = z.input<typeof chatMessageFormSchema>;
export type ChatMessageFormOutput = z.output<typeof chatMessageFormSchema>;
