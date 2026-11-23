import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type Team = components["schemas"]["Team"];
export type TeamListResponse = components["schemas"]["TeamListResponse"];
export type CreateTeamRequest = components["schemas"]["CreateTeamRequest"];
export type UpdateTeamRequest = components["schemas"]["UpdateTeamRequest"];

export const createTeamSchema = z.object({
  name: z.string().trim().min(1, "Team name is required").max(128),
  description: z.string().trim().max(500),
  budget_usd: z.number().finite().min(0, "Budget must be >= 0"),
  rate_limit_rpm: z.number().int().min(0, "RPM limit must be >= 0"),
  enabled: z.boolean(),
});

export type CreateTeamFormValues = z.infer<typeof createTeamSchema>;
