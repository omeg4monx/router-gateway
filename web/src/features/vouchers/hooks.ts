import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient, extractErrorMessage } from "@/shared/api/client";
import type {
  ClaimVoucherRequest,
  ClaimVoucherResponse,
  CreateVoucherRequest,
  UpdateVoucherRequest,
  Voucher,
  VoucherListResponse,
  VoucherRedemptionListResponse,
} from "./schemas";

export const vouchersQueryKeys = {
  all: ["vouchers"] as const,
  list: ["vouchers", "list"] as const,
  redemptions: (id: string) => ["vouchers", "redemptions", id] as const,
};

export function useVouchersList() {
  return useQuery<VoucherListResponse, Error>({
    queryKey: vouchersQueryKeys.list,
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/vouchers");
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useVoucherRedemptions(id: string | null) {
  return useQuery<VoucherRedemptionListResponse, Error>({
    queryKey: vouchersQueryKeys.redemptions(id ?? ""),
    enabled: Boolean(id),
    queryFn: async () => {
      const { data, error } = await apiClient.GET("/api/vouchers/{id}/redemptions", {
        params: { path: { id: id ?? "" } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
  });
}

export function useCreateVoucherMutation() {
  const queryClient = useQueryClient();
  return useMutation<Voucher, Error, CreateVoucherRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/vouchers", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: vouchersQueryKeys.all });
    },
  });
}

export function useUpdateVoucherMutation() {
  const queryClient = useQueryClient();
  return useMutation<Voucher, Error, { id: string; body: UpdateVoucherRequest }>({
    mutationFn: async ({ id, body }) => {
      const { data, error } = await apiClient.PATCH("/api/vouchers/{id}", {
        params: { path: { id } },
        body,
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: vouchersQueryKeys.all });
    },
  });
}

export function useDeleteVoucherMutation() {
  const queryClient = useQueryClient();
  return useMutation<{ deleted: boolean }, Error, string>({
    mutationFn: async (id) => {
      const { data, error } = await apiClient.DELETE("/api/vouchers/{id}", {
        params: { path: { id } },
      });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: vouchersQueryKeys.all });
    },
  });
}

export function useClaimVoucherMutation() {
  const queryClient = useQueryClient();
  return useMutation<ClaimVoucherResponse, Error, ClaimVoucherRequest>({
    mutationFn: async (body) => {
      const { data, error } = await apiClient.POST("/api/vouchers/claim", { body });
      if (data) {
        return data;
      }
      throw new Error(extractErrorMessage(error));
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: vouchersQueryKeys.all });
    },
  });
}
