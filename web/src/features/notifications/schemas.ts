import { z } from "zod";
import type { components } from "@/shared/api/schema";

export type SystemNotification = components["schemas"]["SystemNotification"];
export type CreateNotificationRequest = components["schemas"]["CreateNotificationRequest"];
export type NotificationListResponse = components["schemas"]["NotificationListResponse"];
export type BackgroundJobStatus = components["schemas"]["BackgroundJobStatus"];
export type BackgroundJobListResponse = components["schemas"]["BackgroundJobListResponse"];
export type TriggerJobResponse = components["schemas"]["TriggerJobResponse"];

export const NOTIFICATION_SEVERITIES = ["info", "success", "warning", "error"] as const;

export const createNotificationSchema = z.object({
  title: z.string().trim().min(1, "Title is required").max(160),
  message: z.string().trim().min(1, "Message is required").max(1000),
  severity: z.enum(NOTIFICATION_SEVERITIES),
  category: z.string().trim().min(1, "Category is required").max(64),
});

export type CreateNotificationFormValues = z.infer<typeof createNotificationSchema>;

export function formatIntervalMs(ms: number): string {
  if (!Number.isFinite(ms) || ms <= 0) {
    return "--";
  }
  if (ms < 1_000) {
    return `${Math.round(ms)}ms`;
  }
  const seconds = Math.round(ms / 1_000);
  if (seconds < 60) {
    return `${seconds}s`;
  }
  const minutes = Math.round(seconds / 60);
  if (minutes < 60) {
    return `${minutes}m`;
  }
  const hours = Math.round(minutes / 60);
  return `${hours}h`;
}
