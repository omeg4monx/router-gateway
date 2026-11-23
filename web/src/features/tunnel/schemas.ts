import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type TunnelProvider = components["schemas"]["TunnelProvider"];
export type TailscaleMode = components["schemas"]["TailscaleMode"];
export type TunnelState = components["schemas"]["TunnelState"];
export type TunnelStatus = components["schemas"]["TunnelStatus"];
export type UpdateTunnelConfigRequest = components["schemas"]["UpdateTunnelConfigRequest"];
export type StartTunnelRequest = components["schemas"]["StartTunnelRequest"];
export type TunnelEvent = components["schemas"]["TunnelEvent"];
export type TunnelEventListResponse = components["schemas"]["TunnelEventListResponse"];

export const TUNNEL_PROVIDERS = ["cloudflare", "tailscale"] as const;
export const TAILSCALE_MODES = ["funnel", "serve"] as const;

function isLoopbackURL(raw: string): boolean {
  try {
    const u = new URL(raw);
    if (u.protocol !== "http:" && u.protocol !== "https:") {
      return false;
    }
    const host = u.hostname.toLowerCase();
    return host === "localhost" || host === "127.0.0.1" || host === "::1" || host === "[::1]";
  } catch {
    return false;
  }
}

export const tunnelConfigSchema = z.object({
  provider: z.enum(TUNNEL_PROVIDERS),
  local_url: z
    .string()
    .trim()
    .refine(isLoopbackURL, "Local URL must target a loopback host (127.0.0.1, localhost, or ::1)"),
  tailscale_mode: z.enum(TAILSCALE_MODES),
  auto_start: z.boolean().default(false),
});

export type TunnelConfigFormInput = z.input<typeof tunnelConfigSchema>;
export type TunnelConfigFormOutput = z.output<typeof tunnelConfigSchema>;
