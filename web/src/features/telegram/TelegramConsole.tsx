"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  Bell,
  KeyRound,
  Link2,
  Play,
  RefreshCw,
  Send,
  Shield,
  Terminal,
  Trash2,
  Unlink,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClearTelegramEventsMutation,
  useDeleteTelegramWebhookMutation,
  useRegisterTelegramWebhookMutation,
  useRotateTelegramSecretMutation,
  useSendTelegramTestMutation,
  useSimulateTelegramCommandMutation,
  useTelegramConfig,
  useTelegramEvents,
  useUpdateTelegramConfigMutation,
} from "./hooks";
import {
  type TelegramCommandFormInput,
  type TelegramCommandFormOutput,
  type TelegramConfigFormInput,
  type TelegramConfigFormOutput,
  type TelegramTestFormInput,
  type TelegramTestFormOutput,
  telegramCommandFormSchema,
  telegramConfigFormSchema,
  telegramTestFormSchema,
} from "./schemas";

const QUICK_COMMANDS = ["/status", "/providers", "/models", "/jobs", "/ping", "/help"];

export function TelegramConsole() {
  const configQuery = useTelegramConfig();
  const eventsQuery = useTelegramEvents();

  const updateConfig = useUpdateTelegramConfigMutation();
  const rotateSecret = useRotateTelegramSecretMutation();
  const registerWebhook = useRegisterTelegramWebhookMutation();
  const deleteWebhook = useDeleteTelegramWebhookMutation();
  const sendTest = useSendTelegramTestMutation();
  const simulateCommand = useSimulateTelegramCommandMutation();
  const clearEvents = useClearTelegramEventsMutation();

  const [lastCommandReply, setLastCommandReply] = useState<{
    command: string;
    reply: string;
    executed_at: string;
  } | null>(null);
  const [statusNotice, setStatusNotice] = useState<string | null>(null);

  const configForm = useForm<TelegramConfigFormInput, unknown, TelegramConfigFormOutput>({
    resolver: zodResolver(telegramConfigFormSchema),
    defaultValues: {
      enabled: false,
      bot_token: "",
      chat_id: "",
      user_id: "",
      webhook_url: "",
      notify_on_job_failure: true,
      notify_on_auth_change: true,
      notify_on_emit: false,
    },
  });

  const commandForm = useForm<TelegramCommandFormInput, unknown, TelegramCommandFormOutput>({
    resolver: zodResolver(telegramCommandFormSchema),
    defaultValues: { command: "/status" },
  });

  const testForm = useForm<TelegramTestFormInput, unknown, TelegramTestFormOutput>({
    resolver: zodResolver(telegramTestFormSchema),
    defaultValues: { message: "OMEG4MON ROUTER Telegram probe" },
  });

  const cfg = configQuery.data;
  useEffect(() => {
    if (!cfg || configForm.formState.isDirty) {
      return;
    }
    configForm.reset({
      enabled: cfg.enabled,
      bot_token: "",
      chat_id: cfg.chat_id,
      user_id: cfg.user_id,
      webhook_url: cfg.webhook_url,
      notify_on_job_failure: cfg.notify_on_job_failure,
      notify_on_auth_change: cfg.notify_on_auth_change,
      notify_on_emit: cfg.notify_on_emit,
    });
  }, [cfg, configForm]);

  const activeError =
    configQuery.error ??
    eventsQuery.error ??
    updateConfig.error ??
    rotateSecret.error ??
    registerWebhook.error ??
    deleteWebhook.error ??
    sendTest.error ??
    simulateCommand.error ??
    clearEvents.error;

  const onSaveConfig = configForm.handleSubmit((values) => {
    setStatusNotice(null);
    updateConfig.mutate(
      {
        enabled: values.enabled,
        bot_token: values.bot_token ? values.bot_token : undefined,
        chat_id: values.chat_id,
        user_id: values.user_id,
        webhook_url: values.webhook_url,
        notify_on_job_failure: values.notify_on_job_failure,
        notify_on_auth_change: values.notify_on_auth_change,
        notify_on_emit: values.notify_on_emit,
      },
      {
        onSuccess: (updated) => {
          configForm.reset({
            enabled: updated.enabled,
            bot_token: "",
            chat_id: updated.chat_id,
            user_id: updated.user_id,
            webhook_url: updated.webhook_url,
            notify_on_job_failure: updated.notify_on_job_failure,
            notify_on_auth_change: updated.notify_on_auth_change,
            notify_on_emit: updated.notify_on_emit,
          });
          setStatusNotice("Telegram bot configuration saved to AES-256-GCM vault.");
        },
      },
    );
  });

  const onRunCommand = commandForm.handleSubmit((values) => {
    simulateCommand.mutate(
      { command: values.command },
      {
        onSuccess: (res) => {
          setLastCommandReply(res);
        },
      },
    );
  });

  const onSendTest = testForm.handleSubmit((values) => {
    setStatusNotice(null);
    sendTest.mutate(
      { message: values.message },
      {
        onSuccess: (res) => {
          setStatusNotice(`Test message delivered to chat ${res.chat_id}.`);
        },
      },
    );
  });

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Bot Control Plane
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Telegram Bot Controller
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Vault-encrypted bot credentials, fail-closed webhook authentication, eventbus alert
            forwarding, and in-process command simulator.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Badge variant={cfg?.enabled ? "default" : "secondary"}>
            {cfg?.enabled ? "BOT ENABLED" : "BOT DISABLED"}
          </Badge>
          <Badge variant={cfg?.bot_token_configured ? "default" : "secondary"}>
            {cfg?.bot_token_configured ? `TOKEN ${cfg.bot_token_mask}` : "NO TOKEN"}
          </Badge>
          <Badge variant={cfg?.webhook_secret_configured ? "default" : "secondary"}>
            {cfg?.webhook_secret_configured ? "SECRET READY" : "NO SECRET"}
          </Badge>
        </div>
      </section>

      <EmberRule />

      {activeError ? (
        <Alert variant="destructive">
          <AlertTitle>Telegram operation error</AlertTitle>
          <AlertDescription>{activeError.message}</AlertDescription>
        </Alert>
      ) : null}

      {statusNotice ? (
        <Alert>
          <AlertTitle>Status</AlertTitle>
          <AlertDescription>{statusNotice}</AlertDescription>
        </Alert>
      ) : null}

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-base">
              <Shield className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
              Credentials &amp; Webhook Security
            </CardTitle>
            <CardDescription>
              Public webhook endpoint: <span className="font-mono">/api/telegram/webhook</span>.
              Fails closed unless both secret token and admin user ID match.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={onSaveConfig} className="space-y-4">
              <label className="flex items-center justify-between rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 py-2.5 text-sm">
                <span className="font-medium text-[var(--text-primary)]">
                  Enable Telegram webhook &amp; alerts
                </span>
                <input
                  type="checkbox"
                  {...configForm.register("enabled")}
                  className="h-4 w-4 accent-[#E56A4A]"
                />
              </label>

              <div className="space-y-1.5">
                <Label htmlFor="tg-token">
                  Bot Token (encrypted in vault; leave blank to keep)
                </Label>
                <Input
                  id="tg-token"
                  type="password"
                  autoComplete="off"
                  placeholder={cfg?.bot_token_configured ? cfg.bot_token_mask : "123456:AAH..."}
                  {...configForm.register("bot_token")}
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="tg-chat-id">Target Chat ID</Label>
                  <Input
                    id="tg-chat-id"
                    placeholder="e.g. 123456789"
                    {...configForm.register("chat_id")}
                  />
                </div>
                <div className="space-y-1.5">
                  <Label htmlFor="tg-user-id">Authorized Admin User ID</Label>
                  <Input
                    id="tg-user-id"
                    placeholder="Required for webhook commands"
                    {...configForm.register("user_id")}
                  />
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="tg-webhook-url">Public HTTPS Webhook URL</Label>
                <Input
                  id="tg-webhook-url"
                  placeholder="https://router.example.com/api/telegram/webhook"
                  {...configForm.register("webhook_url")}
                />
                {configForm.formState.errors.webhook_url ? (
                  <p className="text-xs text-red-500">
                    {configForm.formState.errors.webhook_url.message}
                  </p>
                ) : null}
              </div>

              <div className="space-y-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3">
                <p className="flex items-center gap-1.5 text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <Bell className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
                  EventBus Alert Forwarding
                </p>
                <label className="flex items-center justify-between text-sm">
                  <span>Background job failures (jobs.failed)</span>
                  <input
                    type="checkbox"
                    {...configForm.register("notify_on_job_failure")}
                    className="h-4 w-4 accent-[#E56A4A]"
                  />
                </label>
                <label className="flex items-center justify-between text-sm">
                  <span>Auth policy changes (auth.policy_changed)</span>
                  <input
                    type="checkbox"
                    {...configForm.register("notify_on_auth_change")}
                    className="h-4 w-4 accent-[#E56A4A]"
                  />
                </label>
                <label className="flex items-center justify-between text-sm">
                  <span>General system notifications (notifications.emit)</span>
                  <input
                    type="checkbox"
                    {...configForm.register("notify_on_emit")}
                    className="h-4 w-4 accent-[#E56A4A]"
                  />
                </label>
              </div>

              <div className="flex flex-wrap items-center gap-2 pt-1">
                <Button type="submit" disabled={updateConfig.isPending}>
                  {updateConfig.isPending ? "Saving..." : "Save Configuration"}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  disabled={rotateSecret.isPending}
                  onClick={() =>
                    rotateSecret.mutate(undefined, {
                      onSuccess: () =>
                        setStatusNotice(
                          "Webhook secret rotated. Click Register Webhook to sync Telegram.",
                        ),
                    })
                  }
                >
                  <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Rotate Secret</span>
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  disabled={registerWebhook.isPending}
                  onClick={() =>
                    registerWebhook.mutate(undefined, {
                      onSuccess: (res) =>
                        setStatusNotice(`Webhook registered: ${res.webhook_url ?? "ok"}`),
                    })
                  }
                >
                  <Link2 className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Register Webhook</span>
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  disabled={deleteWebhook.isPending}
                  onClick={() =>
                    deleteWebhook.mutate(undefined, {
                      onSuccess: () => setStatusNotice("Webhook removed from Telegram Bot API."),
                    })
                  }
                >
                  <Unlink className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Delete Webhook</span>
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Terminal className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                In-Process Command Simulator
              </CardTitle>
              <CardDescription>
                Execute bot commands locally without an external Telegram webhook round-trip.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex flex-wrap gap-1.5">
                {QUICK_COMMANDS.map((cmd) => (
                  <Button
                    key={cmd}
                    type="button"
                    size="sm"
                    variant="secondary"
                    onClick={() => {
                      commandForm.setValue("command", cmd);
                      simulateCommand.mutate(
                        { command: cmd },
                        { onSuccess: (res) => setLastCommandReply(res) },
                      );
                    }}
                  >
                    <span className="font-mono text-xs">{cmd}</span>
                  </Button>
                ))}
              </div>

              <form onSubmit={onRunCommand} className="flex gap-2">
                <Input
                  placeholder="/status"
                  className="font-mono"
                  {...commandForm.register("command")}
                />
                <Button type="submit" disabled={simulateCommand.isPending}>
                  <Play className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>Run</span>
                </Button>
              </form>

              {lastCommandReply ? (
                <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-3">
                  <div className="mb-1.5 flex items-center justify-between font-mono text-xs text-[var(--text-muted)]">
                    <span>{lastCommandReply.command}</span>
                    <span>{lastCommandReply.executed_at}</span>
                  </div>
                  <pre className="whitespace-pre-wrap font-mono text-xs text-[var(--text-primary)]">
                    {lastCommandReply.reply}
                  </pre>
                </div>
              ) : null}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Send className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
                Outbound Delivery Probe
              </CardTitle>
              <CardDescription>
                Send a test notification via Telegram Bot API to the configured Chat ID.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={onSendTest} className="flex gap-2">
                <Input placeholder="Test message payload..." {...testForm.register("message")} />
                <Button type="submit" variant="outline" disabled={sendTest.isPending}>
                  {sendTest.isPending ? "Sending..." : "Send Test"}
                </Button>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle className="text-base">Delivery &amp; Command Audit Log</CardTitle>
            <CardDescription>
              Recent inbound webhook updates, simulator runs, and outbound alert deliveries.
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Button
              type="button"
              size="sm"
              variant="outline"
              onClick={() => void eventsQuery.refetch()}
            >
              <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Refresh</span>
            </Button>
            <Button
              type="button"
              size="sm"
              variant="secondary"
              disabled={clearEvents.isPending}
              onClick={() => clearEvents.mutate()}
            >
              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Clear</span>
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {eventsQuery.data?.events.length ? (
            <div className="divide-y divide-[var(--border-subtle)] rounded-md border border-[var(--border-subtle)]">
              {eventsQuery.data.events.map((ev) => (
                <div
                  key={ev.id}
                  className="flex flex-col gap-1 px-3 py-2.5 text-sm sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="space-y-0.5">
                    <div className="flex flex-wrap items-center gap-2">
                      <Badge
                        variant={
                          ev.status === "failed" || ev.status === "rejected"
                            ? "destructive"
                            : "secondary"
                        }
                      >
                        {ev.status}
                      </Badge>
                      <span className="font-mono text-xs text-[var(--text-muted)]">
                        {ev.direction}/{ev.kind}
                      </span>
                      {ev.command ? (
                        <span className="font-mono text-xs font-medium text-[var(--accent-text)]">
                          {ev.command}
                        </span>
                      ) : null}
                    </div>
                    <p className="text-xs text-[var(--text-secondary)]">{ev.error || ev.summary}</p>
                  </div>
                  <span className="shrink-0 font-mono text-[11px] text-[var(--text-muted)]">
                    {ev.created_at}
                  </span>
                </div>
              ))}
            </div>
          ) : (
            <p className="text-sm text-[var(--text-muted)]">No Telegram events recorded yet.</p>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
