"use client";

import { useState } from "react";
import { cn } from "@/shared/lib/utils";
import { getProviderMonogram } from "./catalogUtils";

const PROVIDER_ICON_ALIASES: Record<string, string> = {
  "openai-responses": "oai-r",
  "openai-compatible": "custom-openai.svg",
  "anthropic-compatible": "custom-anthropic.svg",
  "gemini-oauth": "gemini-cli",
  "github-copilot": "copilot",
  "aws-quick": "aws-polly",
  freepi: "freepi.svg",
  "xiaomi-mimo": "xiaomi-mimo",
  "opencode-zen": "opencode",
  "perplexity-agent": "perplexity",
  "gitlab-duo": "gitlab",
  "vercel-ai-gateway": "vercel",
  "ollama-search": "ollama",
  "minimax-code": "minimax-cn",
  "minimax-code-global": "minimax",
  "deepseek-web": "deepseek",
  "gemini-web": "gemini",
  "kimi-web": "kimi",
  "venice-free": "venice",
  "venice-web": "venice",
  "venice-keyless": "venice",
  "venice-ai": "venice",
};

const PREFIX_ICON_ALIASES: Array<[string, string]> = [
  ["openai-compatible-", "custom-openai.svg"],
  ["anthropic-compatible-", "custom-anthropic.svg"],
  ["custom-embedding-", "selfhosted-embedding"],
];

const DARK_INVERT_PATTERN =
  /\/(openai|oai-r|oai-cc|github|copilot|cursor|xai|grok-cli|ollama|vercel|zed|openrouter|kimi|nous|freepi|custom-openai|custom-anthropic|custom-moonshot)\.(png|svg)$/i;

const failedIconSrcs = new Set<string>();

export function resolveProviderIconSrc(
  providerId: string,
  protocol?: string,
  isCustom?: boolean,
): string | null {
  const normalized = providerId.trim().toLowerCase();
  if (!normalized) {
    return null;
  }

  let mapped = PROVIDER_ICON_ALIASES[normalized];
  if (!mapped) {
    for (const [prefix, target] of PREFIX_ICON_ALIASES) {
      if (normalized.startsWith(prefix)) {
        mapped = target;
        break;
      }
    }
  }

  if (!mapped && isCustom) {
    if (normalized.includes("moonshot") || normalized.includes("kimi")) {
      mapped = "custom-moonshot.svg";
    } else if (protocol === "anthropic") {
      mapped = "custom-anthropic.svg";
    } else if (protocol === "gemini") {
      mapped = "gemini";
    } else {
      mapped = "custom-openai.svg";
    }
  }

  const target = mapped ?? normalized;
  const src = target.endsWith(".svg") ? `/providers/${target}` : `/providers/${target}.png`;
  return failedIconSrcs.has(src) ? null : src;
}

export function ProviderLogo({
  providerId,
  name,
  protocol,
  isCustom = false,
  size = 24,
  className,
}: {
  providerId: string;
  name: string;
  protocol?: string;
  isCustom?: boolean;
  size?: number;
  className?: string;
}) {
  const src = resolveProviderIconSrc(providerId, protocol, isCustom);
  const [erroredSrc, setErroredSrc] = useState<string | null>(null);

  if (!src || erroredSrc === src) {
    return (
      <span className={cn("font-mono text-xs font-bold leading-none", className)}>
        {getProviderMonogram(name, providerId)}
      </span>
    );
  }

  const shouldInvertDark = DARK_INVERT_PATTERN.test(src);

  return (
    <img
      src={src}
      alt={name}
      width={size}
      height={size}
      loading="eager"
      className={cn(
        "provider-logo-img object-contain select-none",
        shouldInvertDark && "provider-logo-invert",
        className,
      )}
      onError={() => {
        failedIconSrcs.add(src);
        setErroredSrc(src);
      }}
    />
  );
}
