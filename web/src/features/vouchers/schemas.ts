import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type Voucher = components["schemas"]["Voucher"];
export type VoucherListResponse = components["schemas"]["VoucherListResponse"];
export type CreateVoucherRequest = components["schemas"]["CreateVoucherRequest"];
export type UpdateVoucherRequest = components["schemas"]["UpdateVoucherRequest"];
export type VoucherRedemption = components["schemas"]["VoucherRedemption"];
export type VoucherRedemptionListResponse = components["schemas"]["VoucherRedemptionListResponse"];
export type ClaimVoucherRequest = components["schemas"]["ClaimVoucherRequest"];
export type ClaimVoucherResponse = components["schemas"]["ClaimVoucherResponse"];

export const createVoucherSchema = z.object({
  code: z.string().trim().max(64),
  name: z.string().trim().min(1, "Voucher name is required").max(128),
  description: z.string().trim().max(500),
  quota_usd: z.number().finite().min(0, "Quota must be >= 0"),
  rate_limit_rpm: z.number().int().min(0, "RPM must be >= 0"),
  allowed_models_text: z.string().trim(),
  team_id: z.string().trim(),
  max_redemptions: z.number().int().min(1, "Max redemptions must be >= 1"),
});

export type CreateVoucherFormValues = z.infer<typeof createVoucherSchema>;

export const claimVoucherSchema = z.object({
  code: z.string().trim().min(1, "Voucher code is required").max(64),
  claimer_id: z.string().trim().max(128),
  key_name: z.string().trim().max(120),
});

export type ClaimVoucherFormValues = z.infer<typeof claimVoucherSchema>;
