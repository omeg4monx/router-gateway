import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type MitmState = components["schemas"]["MitmState"];
export type MitmCAInfo = components["schemas"]["MitmCAInfo"];
export type MitmToolDescriptor = components["schemas"]["MitmToolDescriptor"];
export type MitmInstallGuide = components["schemas"]["MitmInstallGuide"];
export type MitmStatus = components["schemas"]["MitmStatus"];
export type UpdateMitmConfigRequest = components["schemas"]["UpdateMitmConfigRequest"];
export type GenerateMitmCARequest = components["schemas"]["GenerateMitmCARequest"];
export type MitmAlias = components["schemas"]["MitmAlias"];
export type MitmAliasListResponse = components["schemas"]["MitmAliasListResponse"];
export type CreateMitmAliasRequest = components["schemas"]["CreateMitmAliasRequest"];
export type MitmEvent = components["schemas"]["MitmEvent"];
export type MitmEventListResponse = components["schemas"]["MitmEventListResponse"];

export const MITM_TOOL_IDS = ["antigravity", "copilot", "kiro", "cursor"] as const;
export const MITM_ALIAS_TOOLS = ["*", ...MITM_TOOL_IDS] as const;

function isLoopbackHost(host: string): boolean {
  const clean = host
    .trim()
    .toLowerCase()
    .replace(/^\[|\]$/g, "");
  return clean === "localhost" || clean === "127.0.0.1" || clean === "::1";
}

function isLoopbackListenAddr(raw: string): boolean {
  const trimmed = raw.trim();
  let host: string;
  let portStr: string;
  if (trimmed.startsWith("[")) {
    const closeIdx = trimmed.indexOf("]:");
    if (closeIdx <= 1) {
      return false;
    }
    host = trimmed.slice(1, closeIdx);
    portStr = trimmed.slice(closeIdx + 2);
  } else {
    const firstColon = trimmed.indexOf(":");
    const lastColon = trimmed.lastIndexOf(":");
    if (firstColon <= 0 || firstColon !== lastColon || lastColon === trimmed.length - 1) {
      return false;
    }
    host = trimmed.slice(0, lastColon);
    portStr = trimmed.slice(lastColon + 1);
  }
  if (!/^\d+$/.test(portStr)) {
    return false;
  }
  const port = Number(portStr);
  if (!Number.isInteger(port) || port < 0 || port > 65535) {
    return false;
  }
  return isLoopbackHost(host);
}

function isLoopbackURL(raw: string): boolean {
  try {
    const u = new URL(raw.trim());
    if (u.protocol !== "http:" && u.protocol !== "https:") {
      return false;
    }
    return isLoopbackHost(u.hostname);
  } catch {
    return false;
  }
}

export const mitmConfigSchema = z.object({
  listen_addr: z
    .string()
    .trim()
    .refine(
      isLoopbackListenAddr,
      "Listen address must target a loopback host:port (127.0.0.1, localhost, or [::1])",
    ),
  target_gateway_url: z
    .string()
    .trim()
    .refine(
      isLoopbackURL,
      "Target gateway URL must use http/https on a loopback host (127.0.0.1, localhost, or ::1)",
    ),
  enabled_tools: z.array(z.enum(MITM_TOOL_IDS)),
  auto_start: z.boolean().default(false),
});

export type MitmConfigFormInput = z.input<typeof mitmConfigSchema>;
export type MitmConfigFormOutput = z.output<typeof mitmConfigSchema>;

export const createMitmAliasSchema = z.object({
  tool: z.enum(MITM_ALIAS_TOOLS),
  source_model: z.string().trim().min(1, "Source model is required"),
  target_model: z.string().trim().min(1, "Target model is required"),
  enabled: z.boolean().default(true),
});

export type CreateMitmAliasFormInput = z.input<typeof createMitmAliasSchema>;
export type CreateMitmAliasFormOutput = z.output<typeof createMitmAliasSchema>;
