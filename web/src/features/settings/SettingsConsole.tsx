import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertCircle,
  CheckCircle2,
  Eye,
  EyeOff,
  History,
  KeyRound,
  Lock,
  RefreshCw,
  RotateCcw,
  Save,
  ShieldCheck,
  SlidersHorizontal,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useAuthSession } from "@/features/auth/hooks";
import {
  countPasswordRunes,
  MAX_PASSWORD_BYTES,
  MIN_PASSWORD_RUNES,
  utf8ByteLength,
} from "@/features/system/password";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useChangePasswordMutation,
  usePatchSettingsMutation,
  useRollbackSettingMutation,
  useSettings,
  useSettingsHistory,
  useUpdateAuthPolicyMutation,
} from "./hooks";
import {
  buildSettingsPatchRequest,
  type ChangePasswordFormValues,
  changePasswordSchema,
  DEFAULT_SETTINGS_FORM_VALUES,
  type SettingItem,
  type SettingsFormValues,
  settingsFormSchema,
  settingsItemsToFormValues,
} from "./schemas";

function formatSettingValue(value: unknown): string {
  if (value === null || value === undefined) {
    return "default";
  }
  if (typeof value === "boolean") {
    return value ? "true" : "false";
  }
  return String(value);
}

function formatHistoryTimestamp(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  return date.toLocaleString([], {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function syncThemeToUIStore(themeValue: unknown): void {
  if (themeValue === "dark" || themeValue === "light") {
    useUIStore.getState().setTheme(themeValue);
    return;
  }
  if (themeValue === "system" && typeof window !== "undefined") {
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    useUIStore.getState().setTheme(prefersDark ? "dark" : "light");
  }
}

export function SettingsConsole() {
  const settingsQuery = useSettings();
  const historyQuery = useSettingsHistory(50);
  const authSessionQuery = useAuthSession();

  const patchSettingsMutation = usePatchSettingsMutation();
  const rollbackMutation = useRollbackSettingMutation();
  const changePasswordMutation = useChangePasswordMutation();
  const updatePolicyMutation = useUpdateAuthPolicyMutation();

  const [settingsFeedback, setSettingsFeedback] = useState<string | null>(null);
  const [rollbackFeedback, setRollbackFeedback] = useState<string | null>(null);
  const [rollbackBusyId, setRollbackBusyId] = useState<number | null>(null);
  const [passwordFeedback, setPasswordFeedback] = useState<string | null>(null);
  const [policyFeedback, setPolicyFeedback] = useState<string | null>(null);
  const [showCurrentPassword, setShowCurrentPassword] = useState(false);
  const [showNewPassword, setShowNewPassword] = useState(false);

  const settingsForm = useForm<SettingsFormValues>({
    resolver: zodResolver(settingsFormSchema),
    defaultValues: DEFAULT_SETTINGS_FORM_VALUES,
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const passwordForm = useForm<ChangePasswordFormValues>({
    resolver: zodResolver(changePasswordSchema),
    defaultValues: {
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const {
    reset: resetSettingsForm,
    formState: { isDirty: isSettingsDirty },
  } = settingsForm;

  const settingsItems: SettingItem[] = settingsQuery.data ?? [];
  const settingsByKey = new Map(settingsItems.map((item) => [item.key, item]));

  useEffect(() => {
    if (settingsQuery.data && !isSettingsDirty) {
      resetSettingsForm(settingsItemsToFormValues(settingsQuery.data));
    }
  }, [settingsQuery.data, isSettingsDirty, resetSettingsForm]);

  const handleRefreshAll = () => {
    void Promise.all([settingsQuery.refetch(), historyQuery.refetch(), authSessionQuery.refetch()]);
  };

  const isRefreshing =
    settingsQuery.isFetching || historyQuery.isFetching || authSessionQuery.isFetching;

  const onSettingsSubmit = (values: SettingsFormValues) => {
    setSettingsFeedback(null);
    patchSettingsMutation.reset();

    const patchPayload = buildSettingsPatchRequest(values, settingsItems);
    if (!patchPayload) {
      setSettingsFeedback("No setting changes to save.");
      return;
    }

    patchSettingsMutation.mutate(patchPayload, {
      onSuccess: (updatedItems) => {
        if ("ui.theme" in patchPayload.updates) {
          syncThemeToUIStore(patchPayload.updates["ui.theme"]);
        }
        resetSettingsForm(settingsItemsToFormValues(updatedItems));
        const changedCount = Object.keys(patchPayload.updates).length;
        setSettingsFeedback(
          `Saved ${changedCount} ${changedCount === 1 ? "setting" : "settings"} atomically.`,
        );
      },
    });
  };

  const handleRollback = (historyId: number, settingKey: string) => {
    if (rollbackBusyId !== null) {
      return;
    }
    setRollbackFeedback(null);
    rollbackMutation.reset();
    setRollbackBusyId(historyId);

    rollbackMutation.mutate(
      {
        history_id: historyId,
        reason: `Dashboard rollback of ${settingKey}`,
      },
      {
        onSuccess: (restoredItem) => {
          if (restoredItem.key === "ui.theme") {
            syncThemeToUIStore(restoredItem.value);
          }
          const nextItems = settingsItems.some((item) => item.key === restoredItem.key)
            ? settingsItems.map((item) => (item.key === restoredItem.key ? restoredItem : item))
            : [...settingsItems, restoredItem];
          resetSettingsForm(settingsItemsToFormValues(nextItems));
          setRollbackFeedback(
            `Rolled back ${restoredItem.key} to ${formatSettingValue(restoredItem.value)} (rev ${restoredItem.revision}).`,
          );
        },
        onSettled: () => setRollbackBusyId(null),
      },
    );
  };

  const onChangePasswordSubmit = (values: ChangePasswordFormValues) => {
    setPasswordFeedback(null);
    changePasswordMutation.reset();

    changePasswordMutation.mutate(
      {
        current_password: values.currentPassword,
        new_password: values.newPassword,
      },
      {
        onSuccess: () => {
          passwordForm.reset({
            currentPassword: "",
            newPassword: "",
            confirmPassword: "",
          });
          setPasswordFeedback(
            "Administrator password updated and signed session cookie refreshed.",
          );
        },
      },
    );
  };

  const handleToggleRequireLogin = () => {
    const session = authSessionQuery.data;
    if (!session) {
      return;
    }
    const nextValue = !session.require_login;
    if (session.remote_enforced && !nextValue) {
      return;
    }

    setPolicyFeedback(null);
    updatePolicyMutation.reset();

    updatePolicyMutation.mutate(
      { require_login: nextValue },
      {
        onSuccess: (updated) => {
          setPolicyFeedback(
            updated.require_login
              ? "Dashboard login enforcement is now enabled."
              : "Dashboard login enforcement is now disabled for loopback access.",
          );
        },
      },
    );
  };

  const watchedNewPassword = passwordForm.watch("newPassword");
  const newPasswordRunes = countPasswordRunes(watchedNewPassword ?? "");
  const newPasswordBytes = utf8ByteLength(watchedNewPassword ?? "");

  const authSession = authSessionQuery.data;
  const historyItems = historyQuery.data ?? [];

  const getSettingRev = (key: string): number => settingsByKey.get(key)?.revision ?? 0;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Control Plane Configuration
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Settings &amp; Security
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Typed runtime configuration with atomic revision history, rollback, and single-admin
            access policy controls.
          </p>
        </div>

        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleRefreshAll}
          disabled={isRefreshing}
        >
          <RefreshCw className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
          <span>{isRefreshing ? "Refreshing" : "Refresh"}</span>
        </Button>
      </section>

      <Card>
        <CardHeader>
          <div>
            <div className="flex items-center gap-2">
              <SlidersHorizontal
                className="h-4 w-4 text-[var(--accent-primary)]"
                aria-hidden="true"
              />
              <CardTitle>Typed Runtime Settings</CardTitle>
            </div>
            <CardDescription className="mt-0.5">
              Changes are validated server-side and recorded in the revision ledger.
            </CardDescription>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">PATCH /api/settings</span>
        </CardHeader>

        <CardContent>
          {settingsQuery.isError ? (
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load runtime settings</AlertTitle>
                <AlertDescription className="font-mono">
                  {settingsQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>
          ) : (
            <form
              onSubmit={settingsForm.handleSubmit(onSettingsSubmit)}
              noValidate
              className="space-y-6"
            >
              {patchSettingsMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Failed to save settings</AlertTitle>
                    <AlertDescription className="font-mono">
                      {patchSettingsMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {settingsFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Settings Status</AlertTitle>
                    <AlertDescription>{settingsFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div className="grid grid-cols-1 gap-5 md:grid-cols-2">
                <div className="space-y-1.5">
                  <div className="flex items-center justify-between gap-2">
                    <Label htmlFor="setting-ui-theme">UI Theme (ui.theme)</Label>
                    <Badge variant="outline" className="tabular-nums">
                      rev {getSettingRev("ui.theme")}
                    </Badge>
                  </div>
                  <select
                    id="setting-ui-theme"
                    className="flex h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-sm text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]"
                    {...settingsForm.register("ui_theme")}
                  >
                    <option value="dark">dark</option>
                    <option value="light">light</option>
                    <option value="system">system</option>
                  </select>
                  <p className="text-xs text-[var(--text-muted)]">
                    Dashboard color theme mode. Syncs immediately with the active UI theme on save.
                  </p>
                </div>

                <div className="space-y-1.5">
                  <div className="flex items-center justify-between gap-2">
                    <Label htmlFor="setting-ui-locale">UI Locale (ui.locale)</Label>
                    <Badge variant="outline" className="tabular-nums">
                      rev {getSettingRev("ui.locale")}
                    </Badge>
                  </div>
                  <select
                    id="setting-ui-locale"
                    className="flex h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-sm text-[var(--text-primary)] transition-colors focus:border-[var(--accent-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]"
                    {...settingsForm.register("ui_locale")}
                  >
                    <option value="en">en (English)</option>
                    <option value="id">id (Bahasa Indonesia)</option>
                  </select>
                  <p className="text-xs text-[var(--text-muted)]">
                    Dashboard display language locale.
                  </p>
                </div>

                <div className="space-y-1.5">
                  <div className="flex items-center justify-between gap-2">
                    <Label htmlFor="setting-payload-retention">
                      Payload Retention Days (1..90)
                    </Label>
                    <Badge variant="outline" className="tabular-nums">
                      rev {getSettingRev("usage.payload_retention_days")}
                    </Badge>
                  </div>
                  <Input
                    id="setting-payload-retention"
                    type="number"
                    min={1}
                    max={90}
                    step={1}
                    className="tabular-nums"
                    aria-invalid={Boolean(settingsForm.formState.errors.payload_retention_days)}
                    {...settingsForm.register("payload_retention_days", { valueAsNumber: true })}
                  />
                  {settingsForm.formState.errors.payload_retention_days ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {settingsForm.formState.errors.payload_retention_days.message}
                    </p>
                  ) : (
                    <p className="text-xs text-[var(--text-muted)]">
                      Days to retain captured request and response payloads.
                    </p>
                  )}
                </div>

                <div className="space-y-1.5">
                  <div className="flex items-center justify-between gap-2">
                    <Label htmlFor="setting-metadata-retention">
                      Metadata Retention Days (1..3650)
                    </Label>
                    <Badge variant="outline" className="tabular-nums">
                      rev {getSettingRev("usage.metadata_retention_days")}
                    </Badge>
                  </div>
                  <Input
                    id="setting-metadata-retention"
                    type="number"
                    min={1}
                    max={3650}
                    step={1}
                    className="tabular-nums"
                    aria-invalid={Boolean(settingsForm.formState.errors.metadata_retention_days)}
                    {...settingsForm.register("metadata_retention_days", { valueAsNumber: true })}
                  />
                  {settingsForm.formState.errors.metadata_retention_days ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {settingsForm.formState.errors.metadata_retention_days.message}
                    </p>
                  ) : (
                    <p className="text-xs text-[var(--text-muted)]">
                      Days to retain request usage metadata and cost aggregates.
                    </p>
                  )}
                </div>

                <div className="space-y-1.5">
                  <div className="flex items-center justify-between gap-2">
                    <Label htmlFor="setting-default-timeout">
                      Default Upstream Timeout Seconds (5..600)
                    </Label>
                    <Badge variant="outline" className="tabular-nums">
                      rev {getSettingRev("gateway.default_timeout_seconds")}
                    </Badge>
                  </div>
                  <Input
                    id="setting-default-timeout"
                    type="number"
                    min={5}
                    max={600}
                    step={1}
                    className="tabular-nums"
                    aria-invalid={Boolean(settingsForm.formState.errors.default_timeout_seconds)}
                    {...settingsForm.register("default_timeout_seconds", { valueAsNumber: true })}
                  />
                  {settingsForm.formState.errors.default_timeout_seconds ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {settingsForm.formState.errors.default_timeout_seconds.message}
                    </p>
                  ) : (
                    <p className="text-xs text-[var(--text-muted)]">
                      Default upstream request timeout when no per-route override is set.
                    </p>
                  )}
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="setting-change-reason">Audit Reason (optional)</Label>
                  <Input
                    id="setting-change-reason"
                    type="text"
                    placeholder="e.g. Increase timeout for long reasoning calls"
                    aria-invalid={Boolean(settingsForm.formState.errors.reason)}
                    {...settingsForm.register("reason")}
                  />
                  {settingsForm.formState.errors.reason ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {settingsForm.formState.errors.reason.message}
                    </p>
                  ) : (
                    <p className="text-xs text-[var(--text-muted)]">
                      Stored alongside each modified setting in the revision ledger.
                    </p>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                <label
                  htmlFor="setting-payload-logging"
                  className="flex cursor-pointer items-start justify-between gap-4 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4 transition-colors hover:border-[var(--border-strong)]"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium text-[var(--text-primary)]">
                        Payload Body Logging
                      </span>
                      <Badge variant="outline" className="tabular-nums">
                        rev {getSettingRev("usage.payload_logging_enabled")}
                      </Badge>
                    </div>
                    <p className="text-xs text-[var(--text-muted)]">
                      Capture request and response payload bodies for usage inspection.
                    </p>
                  </div>
                  <input
                    id="setting-payload-logging"
                    type="checkbox"
                    className="mt-1 h-4 w-4 accent-[var(--accent-primary)]"
                    {...settingsForm.register("payload_logging_enabled")}
                  />
                </label>

                <label
                  htmlFor="setting-chain-fallback"
                  className="flex cursor-pointer items-start justify-between gap-4 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4 transition-colors hover:border-[var(--border-strong)]"
                >
                  <div className="space-y-1">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium text-[var(--text-primary)]">
                        Response Chain Fallback
                      </span>
                      <Badge variant="outline" className="tabular-nums">
                        rev {getSettingRev("gateway.allow_response_chain_fallback")}
                      </Badge>
                    </div>
                    <p className="text-xs text-[var(--text-muted)]">
                      Allow falling back to alternative upstream targets when a response chain
                      target fails.
                    </p>
                  </div>
                  <input
                    id="setting-chain-fallback"
                    type="checkbox"
                    className="mt-1 h-4 w-4 accent-[var(--accent-primary)]"
                    {...settingsForm.register("allow_response_chain_fallback")}
                  />
                </label>
              </div>

              <div className="flex items-center justify-end gap-3">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => {
                    if (settingsQuery.data) {
                      settingsForm.reset(settingsItemsToFormValues(settingsQuery.data));
                      setSettingsFeedback(null);
                      patchSettingsMutation.reset();
                    }
                  }}
                >
                  Discard Changes
                </Button>
                <Button
                  type="submit"
                  size="sm"
                  disabled={patchSettingsMutation.isPending || settingsQuery.isPending}
                >
                  <Save className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>{patchSettingsMutation.isPending ? "Saving..." : "Save Settings"}</span>
                </Button>
              </div>
            </form>
          )}
        </CardContent>
      </Card>

      <EmberRule />

      <section className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <CardTitle>Login Enforcement Policy</CardTitle>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">PUT /api/auth/policy</span>
          </CardHeader>
          <CardContent className="space-y-4">
            {updatePolicyMutation.error ? (
              <Alert variant="destructive">
                <AlertCircle
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                  aria-hidden="true"
                />
                <div>
                  <AlertTitle>Failed to update login policy</AlertTitle>
                  <AlertDescription className="font-mono">
                    {updatePolicyMutation.error.message}
                  </AlertDescription>
                </div>
              </Alert>
            ) : null}

            {policyFeedback ? (
              <Alert>
                <CheckCircle2
                  className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                  aria-hidden="true"
                />
                <div>
                  <AlertTitle>Policy Updated</AlertTitle>
                  <AlertDescription>{policyFeedback}</AlertDescription>
                </div>
              </Alert>
            ) : null}

            <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <span className="text-sm font-medium text-[var(--text-primary)]">
                  Require Dashboard Login
                </span>
                <div className="flex items-center gap-2">
                  {authSession?.remote_enforced ? (
                    <Badge variant="warning">Remote-Bind Locked</Badge>
                  ) : null}
                  {authSessionQuery.isError ? (
                    <Badge variant="destructive">Unavailable</Badge>
                  ) : (
                    <Badge variant={authSession?.require_login ? "success" : "secondary"}>
                      {authSessionQuery.isPending
                        ? "Checking"
                        : authSession?.require_login
                          ? "Enforced"
                          : "Optional"}
                    </Badge>
                  )}
                </div>
              </div>

              {authSessionQuery.isError ? (
                <Alert
                  variant="destructive"
                  className="mt-3 flex-col justify-between gap-3 sm:flex-row sm:items-center"
                >
                  <div className="flex items-start gap-2.5">
                    <AlertCircle
                      className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                      aria-hidden="true"
                    />
                    <div>
                      <AlertTitle>Unable to load authentication policy</AlertTitle>
                      <AlertDescription className="font-mono">
                        {authSessionQuery.error.message}
                      </AlertDescription>
                    </div>
                  </div>
                  <Button
                    type="button"
                    variant="destructive"
                    size="xs"
                    onClick={() => void authSessionQuery.refetch()}
                    disabled={authSessionQuery.isFetching}
                    className="shrink-0"
                  >
                    {authSessionQuery.isFetching ? "Retrying..." : "Retry"}
                  </Button>
                </Alert>
              ) : (
                <>
                  <p className="mt-2 text-xs text-[var(--text-muted)]">
                    Controls whether loopback requests to administrative APIs and dashboard routes
                    require a signed operator session cookie.
                  </p>

                  {authSession?.remote_enforced ? (
                    <Alert variant="warning" className="mt-3">
                      <Lock
                        className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-warn)]"
                        aria-hidden="true"
                      />
                      <div>
                        <AlertTitle>Non-Loopback Bind Lock Active</AlertTitle>
                        <AlertDescription>
                          This instance is bound to a non-loopback address. Disabling login
                          enforcement is locked unless{" "}
                          <span className="font-mono">OMEG4MON_ALLOW_INSECURE_REMOTE</span> is
                          explicitly enabled at startup.
                        </AlertDescription>
                      </div>
                    </Alert>
                  ) : null}

                  <div className="mt-4 flex items-center justify-between gap-3">
                    <span className="font-mono text-xs text-[var(--text-muted)]">
                      Subject:{" "}
                      {authSessionQuery.isPending
                        ? "checking..."
                        : authSession?.subject || "anonymous"}
                    </span>
                    <Button
                      type="button"
                      variant={authSession?.require_login ? "secondary" : "default"}
                      size="sm"
                      onClick={handleToggleRequireLogin}
                      disabled={
                        authSessionQuery.isPending ||
                        updatePolicyMutation.isPending ||
                        !authSession ||
                        (authSession.remote_enforced && authSession.require_login)
                      }
                    >
                      <ShieldCheck className="h-3.5 w-3.5" aria-hidden="true" />
                      <span>
                        {updatePolicyMutation.isPending
                          ? "Updating..."
                          : authSession?.require_login
                            ? "Disable Login Requirement"
                            : "Enable Login Requirement"}
                      </span>
                    </Button>
                  </div>
                </>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <KeyRound className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
              <CardTitle>Change Admin Password</CardTitle>
            </div>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              PUT /api/auth/password
            </span>
          </CardHeader>
          <CardContent>
            <form
              onSubmit={passwordForm.handleSubmit(onChangePasswordSubmit)}
              noValidate
              className="space-y-4"
            >
              {changePasswordMutation.error ? (
                <Alert variant="destructive">
                  <AlertCircle
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Password update failed</AlertTitle>
                    <AlertDescription className="font-mono">
                      {changePasswordMutation.error.message}
                    </AlertDescription>
                  </div>
                </Alert>
              ) : null}

              {passwordFeedback ? (
                <Alert>
                  <CheckCircle2
                    className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                    aria-hidden="true"
                  />
                  <div>
                    <AlertTitle>Password Updated</AlertTitle>
                    <AlertDescription>{passwordFeedback}</AlertDescription>
                  </div>
                </Alert>
              ) : null}

              <div>
                <Label htmlFor="current-admin-password">Current Password</Label>
                <div className="relative mt-1.5">
                  <Input
                    id="current-admin-password"
                    type={showCurrentPassword ? "text" : "password"}
                    autoComplete="current-password"
                    placeholder="Enter current password"
                    aria-invalid={Boolean(passwordForm.formState.errors.currentPassword)}
                    className="pr-10"
                    {...passwordForm.register("currentPassword")}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-sm"
                    onClick={() => setShowCurrentPassword((prev) => !prev)}
                    aria-label={
                      showCurrentPassword ? "Hide current password" : "Show current password"
                    }
                    className="absolute right-1.5 top-1/2 -translate-y-1/2"
                  >
                    {showCurrentPassword ? (
                      <EyeOff className="h-4 w-4" aria-hidden="true" />
                    ) : (
                      <Eye className="h-4 w-4" aria-hidden="true" />
                    )}
                  </Button>
                </div>
                {passwordForm.formState.errors.currentPassword ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {passwordForm.formState.errors.currentPassword.message}
                  </p>
                ) : null}
              </div>

              <div>
                <div className="flex items-center justify-between gap-2">
                  <Label htmlFor="new-admin-password">New Password</Label>
                  <span className="font-mono text-[11px] tabular-nums text-[var(--text-muted)]">
                    {newPasswordRunes} chars (min {MIN_PASSWORD_RUNES}) / {newPasswordBytes} of{" "}
                    {MAX_PASSWORD_BYTES} UTF-8 bytes
                  </span>
                </div>
                <div className="relative mt-1.5">
                  <Input
                    id="new-admin-password"
                    type={showNewPassword ? "text" : "password"}
                    autoComplete="new-password"
                    placeholder="8 to 72 UTF-8 bytes"
                    aria-invalid={Boolean(passwordForm.formState.errors.newPassword)}
                    className="pr-10"
                    {...passwordForm.register("newPassword")}
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon-sm"
                    onClick={() => setShowNewPassword((prev) => !prev)}
                    aria-label={showNewPassword ? "Hide new password" : "Show new password"}
                    className="absolute right-1.5 top-1/2 -translate-y-1/2"
                  >
                    {showNewPassword ? (
                      <EyeOff className="h-4 w-4" aria-hidden="true" />
                    ) : (
                      <Eye className="h-4 w-4" aria-hidden="true" />
                    )}
                  </Button>
                </div>
                {passwordForm.formState.errors.newPassword ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {passwordForm.formState.errors.newPassword.message}
                  </p>
                ) : null}
              </div>

              <div>
                <Label htmlFor="confirm-admin-password">Confirm New Password</Label>
                <Input
                  id="confirm-admin-password"
                  type={showNewPassword ? "text" : "password"}
                  autoComplete="new-password"
                  placeholder="Re-enter new password"
                  aria-invalid={Boolean(passwordForm.formState.errors.confirmPassword)}
                  className="mt-1.5"
                  {...passwordForm.register("confirmPassword")}
                />
                {passwordForm.formState.errors.confirmPassword ? (
                  <p className="mt-1 text-xs text-[var(--status-error)]">
                    {passwordForm.formState.errors.confirmPassword.message}
                  </p>
                ) : null}
              </div>

              <div className="flex justify-end">
                <Button type="submit" size="sm" disabled={changePasswordMutation.isPending}>
                  <KeyRound className="h-3.5 w-3.5" aria-hidden="true" />
                  <span>
                    {changePasswordMutation.isPending ? "Updating..." : "Update Admin Password"}
                  </span>
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      </section>

      <EmberRule />

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <History className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
            <CardTitle>Settings Revision History</CardTitle>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            GET /api/settings/history
          </span>
        </CardHeader>

        {rollbackMutation.error ? (
          <div className="px-5 pt-4">
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Rollback failed</AlertTitle>
                <AlertDescription className="font-mono">
                  {rollbackMutation.error.message}
                </AlertDescription>
              </div>
            </Alert>
          </div>
        ) : null}

        {rollbackFeedback ? (
          <div className="px-5 pt-4">
            <Alert>
              <CheckCircle2
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-ok)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Rollback Applied</AlertTitle>
                <AlertDescription>{rollbackFeedback}</AlertDescription>
              </div>
            </Alert>
          </div>
        ) : null}

        {historyQuery.isError ? (
          <CardContent>
            <Alert variant="destructive">
              <AlertCircle
                className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
                aria-hidden="true"
              />
              <div>
                <AlertTitle>Unable to load revision history</AlertTitle>
                <AlertDescription className="font-mono">
                  {historyQuery.error.message}
                </AlertDescription>
              </div>
            </Alert>
          </CardContent>
        ) : historyQuery.isPending ? (
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading setting revision history...
          </CardContent>
        ) : historyItems.length === 0 ? (
          <CardContent className="py-10 text-center">
            <p className="text-sm font-medium text-[var(--text-primary)]">
              No setting revisions recorded yet
            </p>
            <p className="mt-1 text-xs text-[var(--text-muted)]">
              Save any setting change above to start building an auditable rollback timeline.
            </p>
          </CardContent>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse text-left text-sm">
              <thead>
                <tr className="border-b border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[11px] font-medium uppercase tracking-wider text-[var(--text-muted)]">
                  <th scope="col" className="px-5 py-3">
                    Setting Key
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Transition
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Revision
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Actor &amp; Reason
                  </th>
                  <th scope="col" className="px-4 py-3">
                    Timestamp
                  </th>
                  <th scope="col" className="px-5 py-3 text-right">
                    Action
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[var(--border-subtle)]">
                {historyItems.map((entry) => {
                  const isRollingBack = rollbackBusyId === entry.id;
                  return (
                    <tr
                      key={entry.id}
                      className="align-middle transition-colors hover:bg-[var(--bg-surface-raised)]"
                    >
                      <td className="px-5 py-3 font-mono text-xs font-medium text-[var(--text-primary)]">
                        {entry.setting_key}
                      </td>
                      <td className="px-4 py-3 font-mono text-xs">
                        <span className="text-[var(--text-muted)]">
                          {formatSettingValue(entry.old_value)}
                        </span>
                        <span className="mx-1.5 text-[var(--accent-text)]">-&gt;</span>
                        <span className="text-[var(--text-primary)]">
                          {formatSettingValue(entry.new_value)}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <Badge variant="outline" className="tabular-nums">
                          rev {entry.revision}
                        </Badge>
                      </td>
                      <td className="px-4 py-3 text-xs">
                        <span className="font-mono text-[var(--text-secondary)]">
                          {entry.changed_by}
                        </span>
                        {entry.reason ? (
                          <p className="mt-0.5 text-[var(--text-muted)]">{entry.reason}</p>
                        ) : null}
                      </td>
                      <td className="px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-muted)]">
                        {formatHistoryTimestamp(entry.created_at)}
                      </td>
                      <td className="px-5 py-3 text-right">
                        <Button
                          type="button"
                          variant="outline"
                          size="xs"
                          onClick={() => handleRollback(entry.id, entry.setting_key)}
                          disabled={rollbackBusyId !== null}
                        >
                          <RotateCcw className="h-3.5 w-3.5" aria-hidden="true" />
                          <span>{isRollingBack ? "Rolling back..." : "Rollback"}</span>
                        </Button>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  );
}
