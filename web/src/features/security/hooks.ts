import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import {
  AUDIT_PAGE_LIMIT,
  type ClearSecurityAuditResponse,
  type CreateSecurityAuditRequest,
  type SecurityAuditEvent,
  type SecurityAuditListResponse,
  type SecurityPostureReport,
} from "./schemas";

export const securityQueryKeys = {
  all: ["security"] as const,
  posture: ["security", "posture"] as const,
  audit: (filters: {
    category?: string;
    severity?: string;
    q?: string;
    limit?: number;
    offset?: number;
  }) => ["security", "audit", filters] as const,
};

export function useSecurityPosture() {
  return useQuery<SecurityPostureReport, Error>({
    queryKey: securityQueryKeys.posture,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/security/posture");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useSecurityAuditEvents(filters: {
  category?: string;
  severity?: "info" | "warning" | "critical";
  q?: string;
  limit?: number;
  offset?: number;
}) {
  return useQuery<SecurityAuditListResponse, Error>({
    queryKey: securityQueryKeys.audit(filters),
    placeholderData: (previousData) => previousData,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/security/audit", {
        params: {
          query: {
            category: filters.category || undefined,
            severity: filters.severity || undefined,
            q: filters.q || undefined,
            limit: filters.limit ?? AUDIT_PAGE_LIMIT,
            offset: filters.offset ?? 0,
          },
        },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateSecurityAuditMutation() {
  const queryClient = useQueryClient();
  return useMutation<SecurityAuditEvent, Error, CreateSecurityAuditRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/security/audit", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: securityQueryKeys.all });
    },
  });
}

export function useClearSecurityAuditMutation() {
  const queryClient = useQueryClient();
  return useMutation<ClearSecurityAuditResponse, Error, void>({
    mutationFn: async () => {
      const { data, error } = await apiClient.DELETE("/api/security/audit");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: securityQueryKeys.all });
    },
  });
}
