"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { components } from "@/shared/api/schema";

export type AuthMeResponse = components["schemas"]["AuthMeResponse"];
export type LoginRequest = components["schemas"]["LoginRequest"];
export type LoginResponse = components["schemas"]["LoginResponse"];
export type LogoutResponse = components["schemas"]["LogoutResponse"];

export const authKeys = {
  session: ["auth", "me"] as const,
};

export function useAuthSession() {
  return useQuery<AuthMeResponse, Error>({
    queryKey: authKeys.session,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/auth/me");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    staleTime: 0,
  });
}

export function useLoginMutation() {
  const queryClient = useQueryClient();

  return useMutation<LoginResponse, Error, LoginRequest>({
    mutationFn: async (payload) => {
      const { data, error } = await apiClient.POST("/api/auth/login", {
        body: payload,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (response) => {
      queryClient.setQueryData<AuthMeResponse>(authKeys.session, (previous) =>
        previous
          ? {
              ...previous,
              authenticated: response.authenticated,
              subject: response.subject,
              initialized: true,
            }
          : previous,
      );
      await queryClient.invalidateQueries({ queryKey: authKeys.session });
    },
  });
}

export function useLogoutMutation() {
  const queryClient = useQueryClient();

  return useMutation<LogoutResponse, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.POST("/api/auth/logout");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: async (response) => {
      queryClient.setQueryData<AuthMeResponse>(authKeys.session, (previous) =>
        previous
          ? {
              ...previous,
              authenticated: response.authenticated,
              subject: "",
            }
          : previous,
      );
      await queryClient.invalidateQueries({ queryKey: authKeys.session });
    },
  });
}
