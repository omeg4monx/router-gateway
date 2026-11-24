"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  ArenaBattle,
  ArenaBattleListResponse,
  ArenaLeaderboardResponse,
  ArenaModelListResponse,
  ChatSession,
  ChatSessionDetailResponse,
  ChatSessionListResponse,
  CreateChatSessionRequest,
  RunArenaBattleRequest,
  SendChatMessageRequest,
  SendChatMessageResponse,
  UpdateChatSessionRequest,
} from "./schemas";

export const arenaQueryKeys = {
  all: ["arena"] as const,
  models: ["arena", "models"] as const,
  battles: ["arena", "battles"] as const,
  leaderboard: ["arena", "leaderboard"] as const,
  sessions: ["chat", "sessions"] as const,
  sessionDetail: (id: string) => ["chat", "sessions", id] as const,
};

export function useArenaModels() {
  return useQuery<ArenaModelListResponse, Error>({
    queryKey: arenaQueryKeys.models,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/arena/models");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useArenaBattles() {
  return useQuery<ArenaBattleListResponse, Error>({
    queryKey: arenaQueryKeys.battles,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/arena/battles");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useArenaLeaderboard() {
  return useQuery<ArenaLeaderboardResponse, Error>({
    queryKey: arenaQueryKeys.leaderboard,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/arena/leaderboard");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useRunArenaBattleMutation() {
  const queryClient = useQueryClient();
  return useMutation<ArenaBattle, Error, RunArenaBattleRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/arena/battles", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.battles });
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.leaderboard });
    },
  });
}

export function useVoteArenaBattleMutation() {
  const queryClient = useQueryClient();
  return useMutation<ArenaBattle, Error, { id: string; winner_model: string }>({
    mutationFn: async ({ id, winner_model }) => {
      const { data, error } = await apiClient.POST("/api/arena/battles/{id}/vote", {
        params: { path: { id } },
        body: { winner_model },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.battles });
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.leaderboard });
    },
  });
}

export function useDeleteArenaBattleMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/arena/battles/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.battles });
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.leaderboard });
    },
  });
}

export function useChatSessions() {
  return useQuery<ChatSessionListResponse, Error>({
    queryKey: arenaQueryKeys.sessions,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/chat/sessions");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useChatSessionDetail(id: string | null) {
  return useQuery<ChatSessionDetailResponse, Error>({
    queryKey: arenaQueryKeys.sessionDetail(id ?? ""),
    enabled: Boolean(id),
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/chat/sessions/{id}", {
        params: { path: { id: id ?? "" } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateChatSessionMutation() {
  const queryClient = useQueryClient();
  return useMutation<ChatSession, Error, CreateChatSessionRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/chat/sessions", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.sessions });
    },
  });
}

export function useUpdateChatSessionMutation() {
  const queryClient = useQueryClient();
  return useMutation<ChatSession, Error, { id: string; body: UpdateChatSessionRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/chat/sessions/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: (session) => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.sessions });
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.sessionDetail(session.id) });
    },
  });
}

export function useDeleteChatSessionMutation() {
  const queryClient = useQueryClient();
  return useMutation<void, Error, string>({
    mutationFn: async (id) => {
      const { error } = await apiClient.DELETE("/api/chat/sessions/{id}", {
        params: { path: { id } },
      });
      if (error) {
        throw new Error(extractErrorMessage(error));
      }
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.sessions });
    },
  });
}

export function useSendChatMessageMutation() {
  const queryClient = useQueryClient();
  return useMutation<SendChatMessageResponse, Error, { id: string; body: SendChatMessageRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.POST("/api/chat/sessions/{id}/messages", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: (res) => {
      void queryClient.invalidateQueries({ queryKey: arenaQueryKeys.sessions });
      void queryClient.invalidateQueries({
        queryKey: arenaQueryKeys.sessionDetail(res.session.id),
      });
    },
  });
}
