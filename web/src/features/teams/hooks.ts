import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type { CreateTeamRequest, Team, TeamListResponse, UpdateTeamRequest } from "./schemas";

export const teamsQueryKeys = {
  all: ["teams"] as const,
  list: ["teams", "list"] as const,
};

export function useTeamsList() {
  return useQuery<TeamListResponse, Error>({
    queryKey: teamsQueryKeys.list,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/teams");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateTeamMutation() {
  const queryClient = useQueryClient();
  return useMutation<Team, Error, CreateTeamRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/teams", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: teamsQueryKeys.all });
    },
  });
}

export function useUpdateTeamMutation() {
  const queryClient = useQueryClient();
  return useMutation<Team, Error, { id: string; body: UpdateTeamRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/teams/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: teamsQueryKeys.all });
    },
  });
}

export function useDeleteTeamMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ deleted: boolean }, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.DELETE("/api/teams/{id}", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: teamsQueryKeys.all });
    },
  });
}

export function useAssignTeamKeyMutation() {
  const queryClient = useQueryClient();
  return useMutation<Team, Error, { teamId: string; keyId: string }>({
    mutationFn: async ({ teamId, keyId }) => {
      const { data, error } = await apiClient.POST("/api/teams/{id}/keys", {
        params: { path: { id: teamId } },
        body: { key_id: keyId },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: teamsQueryKeys.all });
    },
  });
}

export function useRemoveTeamKeyMutation() {
  const queryClient = useQueryClient();
  return useMutation<Team, Error, { teamId: string; keyId: string }>({
    mutationFn: async ({ teamId, keyId }) => {
      const { data, error } = await apiClient.DELETE("/api/teams/{id}/keys/{keyId}", {
        params: { path: { id: teamId, keyId } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: teamsQueryKeys.all });
    },
  });
}

export function useResetTeamSpendMutation() {
  const queryClient = useQueryClient();
  return useMutation<Team, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.POST("/api/teams/{id}/reset-spend", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: teamsQueryKeys.all });
    },
  });
}
