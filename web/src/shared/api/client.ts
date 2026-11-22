import createClient, { type ClientOptions } from "openapi-fetch";
import type { components, paths } from "./schema";

export type ErrorDetail = components["schemas"]["ErrorDetail"];
export type ErrorEnvelope = components["schemas"]["ErrorEnvelope"];
export type HealthResponse = components["schemas"]["HealthResponse"];
export type VersionResponse = components["schemas"]["VersionResponse"];
export type InitStatusResponse = components["schemas"]["InitStatusResponse"];
export type InitRequest = components["schemas"]["InitRequest"];
export type InitCompleteResponse = components["schemas"]["InitCompleteResponse"];
export type ShutdownResponse = components["schemas"]["ShutdownResponse"];

const DEFAULT_ERROR_MESSAGE = "An unexpected error occurred";

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}

export function createApiClient(options: ClientOptions = {}) {
  return createClient<paths>({
    baseUrl: "",
    credentials: "include",
    ...options,
  });
}

export const apiClient = createApiClient();

export function extractErrorMessage(err: unknown): string {
  if (typeof err === "string") {
    const trimmed = err.trim();
    return trimmed.length > 0 ? trimmed : DEFAULT_ERROR_MESSAGE;
  }

  if (err instanceof Error) {
    const trimmed = err.message.trim();
    return trimmed.length > 0 ? trimmed : DEFAULT_ERROR_MESSAGE;
  }

  if (isRecord(err)) {
    const nested = err.error;
    if (isRecord(nested) && typeof nested.message === "string") {
      const trimmed = nested.message.trim();
      if (trimmed.length > 0) {
        return trimmed;
      }
    }

    if (typeof err.message === "string") {
      const trimmed = err.message.trim();
      if (trimmed.length > 0) {
        return trimmed;
      }
    }
  }

  return DEFAULT_ERROR_MESSAGE;
}
