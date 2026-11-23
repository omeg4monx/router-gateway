"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  Bell,
  Check,
  CheckCheck,
  Clock,
  Play,
  Plus,
  RefreshCw,
  Trash2,
} from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useBackgroundJobs,
  useClearNotificationsMutation,
  useCreateNotificationMutation,
  useDeleteNotificationMutation,
  useMarkAllNotificationsReadMutation,
  useMarkNotificationReadMutation,
  useNotificationsList,
  useNotificationsStream,
  useTriggerJobMutation,
} from "./hooks";
import {
  type CreateNotificationFormValues,
  createNotificationSchema,
  formatIntervalMs,
  NOTIFICATION_SEVERITIES,
} from "./schemas";

function notificationBadgeVariant(
  severity: string,
): "default" | "secondary" | "success" | "warning" | "destructive" {
  switch (severity) {
    case "success":
      return "success";
    case "warning":
      return "warning";
    case "error":
      return "destructive";
    default:
      return "secondary";
  }
}

function jobStatusVariant(
  status: string,
): "default" | "secondary" | "success" | "warning" | "destructive" {
  switch (status) {
    case "ok":
      return "success";
    case "running":
      return "warning";
    case "error":
    case "panic":
      return "destructive";
    default:
      return "secondary";
  }
}

export function NotificationsConsole() {
  const [unreadOnly, setUnreadOnly] = useState(false);
  const [severityFilter, setSeverityFilter] = useState<
    "" | "info" | "success" | "warning" | "error"
  >("");
  const [showCreateForm, setShowCreateForm] = useState(false);

  const streamConnected = useNotificationsStream();
  const notificationsQuery = useNotificationsList({
    unreadOnly,
    severity: severityFilter || undefined,
  });
  const jobsQuery = useBackgroundJobs();

  const createMutation = useCreateNotificationMutation();
  const markReadMutation = useMarkNotificationReadMutation();
  const markAllReadMutation = useMarkAllNotificationsReadMutation();
  const deleteMutation = useDeleteNotificationMutation();
  const clearMutation = useClearNotificationsMutation();
  const triggerJobMutation = useTriggerJobMutation();

  const form = useForm<CreateNotificationFormValues>({
    resolver: zodResolver(createNotificationSchema),
    defaultValues: {
      title: "",
      message: "",
      severity: "info",
      category: "system",
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    createMutation.mutate(values, {
      onSuccess: () => {
        form.reset();
        setShowCreateForm(false);
      },
    });
  });

  const handleRefresh = () => {
    void Promise.all([notificationsQuery.refetch(), jobsQuery.refetch()]);
  };

  const notifData = notificationsQuery.data;
  const jobsData = jobsQuery.data;
  const isRefreshing = notificationsQuery.isFetching || jobsQuery.isFetching;
  const inboxMutationError =
    markReadMutation.error ??
    markAllReadMutation.error ??
    deleteMutation.error ??
    clearMutation.error;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Operations &amp; Telemetry
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Notifications &amp; Background Jobs
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Real-time control-plane alerts over SSE alongside periodic maintenance job runners.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={streamConnected ? "success" : "secondary"}>
            {streamConnected ? "SSE Live" : "SSE Standby"}
          </Badge>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isRefreshing}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
            <span>{isRefreshing ? "Refreshing" : "Refresh"}</span>
          </Button>
          <Button type="button" size="sm" onClick={() => setShowCreateForm((prev) => !prev)}>
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Send Alert</span>
          </Button>
        </div>
      </section>

      {notificationsQuery.isError || jobsQuery.isError || inboxMutationError ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>Operations request failed</AlertTitle>
            <AlertDescription className="font-mono">
              {notificationsQuery.error?.message ??
                jobsQuery.error?.message ??
                inboxMutationError?.message}
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      {showCreateForm ? (
        <Card className="p-5">
          <h3 className="mb-4 text-sm font-semibold text-[var(--text-primary)]">
            Emit System Notification
          </h3>
          <form onSubmit={onSubmit} className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="notif-title">Title</Label>
              <Input
                id="notif-title"
                {...form.register("title")}
                placeholder="Operational notice"
              />
              {form.formState.errors.title ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.title.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="notif-category">Category</Label>
              <Input id="notif-category" {...form.register("category")} placeholder="system" />
              {form.formState.errors.category ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.category.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="notif-severity">Severity</Label>
              <select
                id="notif-severity"
                {...form.register("severity")}
                className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
              >
                {NOTIFICATION_SEVERITIES.map((sev) => (
                  <option key={sev} value={sev}>
                    {sev}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-1.5 md:col-span-2">
              <Label htmlFor="notif-message">Message</Label>
              <Input
                id="notif-message"
                {...form.register("message")}
                placeholder="Alert details broadcast to connected dashboard sessions"
              />
              {form.formState.errors.message ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.message.message}
                </p>
              ) : null}
            </div>
            {createMutation.isError ? (
              <p className="text-xs text-[var(--status-error)] md:col-span-2">
                {createMutation.error.message}
              </p>
            ) : null}
            <div className="flex justify-end gap-2 md:col-span-2">
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => setShowCreateForm(false)}
              >
                Cancel
              </Button>
              <Button type="submit" size="sm" disabled={createMutation.isPending}>
                {createMutation.isPending ? "Sending..." : "Broadcast"}
              </Button>
            </div>
          </form>
        </Card>
      ) : null}

      <Card>
        <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2.5">
            <Bell className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <div>
              <CardTitle>System Notifications</CardTitle>
              <p className="mt-0.5 font-mono text-xs text-[var(--text-muted)]">
                {notifData
                  ? `${notifData.unread_count} unread / ${notifData.total} total`
                  : "Loading inbox..."}
              </p>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <Button
              type="button"
              variant={unreadOnly ? "default" : "outline"}
              size="xs"
              onClick={() => setUnreadOnly((prev) => !prev)}
            >
              {unreadOnly ? "Unread Only" : "All Items"}
            </Button>
            <select
              value={severityFilter}
              onChange={(e) =>
                setSeverityFilter(e.target.value as "" | "info" | "success" | "warning" | "error")
              }
              aria-label="Filter notifications by severity"
              className="h-8 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 text-xs text-[var(--text-primary)]"
            >
              <option value="">All severities</option>
              {NOTIFICATION_SEVERITIES.map((sev) => (
                <option key={sev} value={sev}>
                  {sev}
                </option>
              ))}
            </select>
            <Button
              type="button"
              variant="outline"
              size="xs"
              onClick={() => markAllReadMutation.mutate()}
              disabled={markAllReadMutation.isPending || (notifData?.unread_count ?? 0) === 0}
            >
              <CheckCheck className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Mark All Read</span>
            </Button>
            <Button
              type="button"
              variant="destructive"
              size="xs"
              onClick={() => clearMutation.mutate()}
              disabled={clearMutation.isPending || (notifData?.total ?? 0) === 0}
            >
              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Clear</span>
            </Button>
          </div>
        </CardHeader>

        <div className="divide-y divide-[var(--border-subtle)]">
          {notifData && notifData.items.length > 0 ? (
            notifData.items.map((item) => (
              <div
                key={item.id}
                className="flex flex-col gap-3 px-5 py-4 sm:flex-row sm:items-start sm:justify-between"
              >
                <div className="min-w-0 space-y-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge variant={notificationBadgeVariant(item.severity)}>{item.severity}</Badge>
                    <span className="font-mono text-xs text-[var(--text-muted)]">
                      {item.category}
                    </span>
                    {!item.read ? (
                      <span className="font-mono text-[11px] uppercase tracking-wider text-[var(--accent-text)]">
                        UNREAD
                      </span>
                    ) : null}
                    <span className="font-mono text-xs tabular-nums text-[var(--text-muted)]">
                      {new Date(item.created_at).toLocaleString()}
                    </span>
                  </div>
                  <p className="text-sm font-medium text-[var(--text-primary)]">{item.title}</p>
                  <p className="text-sm text-[var(--text-secondary)]">{item.message}</p>
                </div>

                <div className="flex shrink-0 items-center gap-2">
                  {!item.read ? (
                    <Button
                      type="button"
                      variant="outline"
                      size="xs"
                      onClick={() => markReadMutation.mutate(item.id)}
                      disabled={markReadMutation.isPending}
                    >
                      <Check className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>Read</span>
                    </Button>
                  ) : null}
                  <Button
                    type="button"
                    variant="secondary"
                    size="xs"
                    onClick={() => deleteMutation.mutate(item.id)}
                    disabled={deleteMutation.isPending}
                    aria-label={`Dismiss notification ${item.title}`}
                  >
                    <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                  </Button>
                </div>
              </div>
            ))
          ) : (
            <div className="px-5 py-8 text-center text-sm text-[var(--text-muted)]">
              {notificationsQuery.isPending
                ? "Loading notifications..."
                : "No notifications to display."}
            </div>
          )}
        </div>
      </Card>

      <EmberRule />

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2.5">
            <Clock className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <div>
              <CardTitle>Background Jobs</CardTitle>
              <p className="mt-0.5 font-mono text-xs text-[var(--text-muted)]">
                GET /api/jobs &middot; POST /api/jobs/&#123;name&#125;/trigger
              </p>
            </div>
          </div>
        </CardHeader>

        {triggerJobMutation.isError ? (
          <div className="px-5 pb-3">
            <Alert variant="destructive">
              <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
              <AlertDescription className="font-mono">
                {triggerJobMutation.error.message}
              </AlertDescription>
            </Alert>
          </div>
        ) : null}

        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-y border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-xs uppercase tracking-wider text-[var(--text-muted)]">
              <tr>
                <th className="px-4 py-2.5">Job Name</th>
                <th className="px-4 py-2.5">Status</th>
                <th className="px-4 py-2.5">Interval</th>
                <th className="px-4 py-2.5">Runs / Errors</th>
                <th className="px-4 py-2.5">Last Duration</th>
                <th className="px-4 py-2.5">Last Finished</th>
                <th className="px-4 py-2.5 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--border-subtle)]">
              {jobsData && jobsData.jobs.length > 0 ? (
                jobsData.jobs.map((job) => (
                  <tr key={job.name} className="hover:bg-[var(--bg-surface-raised)]">
                    <td className="px-4 py-3">
                      <p className="font-mono text-sm font-medium text-[var(--text-primary)]">
                        {job.name}
                      </p>
                      {job.last_error ? (
                        <p className="mt-0.5 font-mono text-xs text-[var(--status-error)]">
                          {job.last_error}
                        </p>
                      ) : null}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3">
                      <Badge variant={jobStatusVariant(job.last_status)}>{job.last_status}</Badge>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-secondary)]">
                      {formatIntervalMs(job.interval_ms)}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-secondary)]">
                      {job.run_count} /{" "}
                      <span
                        className={job.error_count > 0 ? "text-[var(--status-error)]" : undefined}
                      >
                        {job.error_count}
                      </span>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-secondary)]">
                      {job.last_duration_ms}ms
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-muted)]">
                      {job.last_finished_at
                        ? new Date(job.last_finished_at).toLocaleString()
                        : "Never"}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-right">
                      <Button
                        type="button"
                        variant="outline"
                        size="xs"
                        disabled={job.running || triggerJobMutation.isPending}
                        onClick={() => triggerJobMutation.mutate(job.name)}
                      >
                        <Play className="h-3 w-3 text-[var(--accent-primary)]" aria-hidden="true" />
                        <span>Run Now</span>
                      </Button>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan={7}
                    className="px-4 py-8 text-center text-sm text-[var(--text-muted)]"
                  >
                    {jobsQuery.isPending
                      ? "Loading background jobs..."
                      : "No periodic background jobs registered."}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
