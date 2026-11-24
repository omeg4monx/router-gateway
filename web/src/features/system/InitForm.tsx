"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { AlertCircle, CheckCircle2, Eye, EyeOff, Lock } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { useInitializeSystem, useInitStatus } from "./hooks";
import {
  type InitPasswordValues,
  initPasswordSchema,
  MAX_PASSWORD_BYTES,
  MIN_PASSWORD_RUNES,
  validateInitPassword,
} from "./password";

export function InitForm() {
  const navigate = useNavigate();
  const initStatusQuery = useInitStatus();
  const initializeMutation = useInitializeSystem();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<InitPasswordValues>({
    resolver: zodResolver(initPasswordSchema),
    defaultValues: {
      password: "",
      confirmPassword: "",
    },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const passwordValue = watch("password");
  const confirmPasswordValue = watch("confirmPassword");
  const validation = validateInitPassword(passwordValue, confirmPasswordValue);

  const clearMutationError = () => {
    if (initializeMutation.error) {
      initializeMutation.reset();
    }
  };

  const onValidSubmit = (values: InitPasswordValues) => {
    initializeMutation.mutate(
      { password: values.password },
      {
        onSuccess: () => {
          void navigate("/dashboard/", { replace: true });
        },
      },
    );
  };

  const passwordField = register("password", {
    onChange: clearMutationError,
  });
  const confirmPasswordField = register("confirmPassword", {
    onChange: clearMutationError,
  });

  if (initStatusQuery.isPending) {
    return (
      <Card className="p-6">
        <p className="font-mono text-xs uppercase tracking-widest text-[var(--text-muted)]">
          Checking initialization state
        </p>
        <p className="mt-2 text-sm text-[var(--text-secondary)]">
          Querying <span className="font-mono">GET /api/init</span>...
        </p>
      </Card>
    );
  }

  if (initStatusQuery.data?.initialized) {
    return (
      <Card className="p-6">
        <div className="flex items-start gap-3">
          <CheckCircle2
            className="mt-0.5 h-5 w-5 shrink-0 text-[var(--status-ok)]"
            aria-hidden="true"
          />
          <div>
            <h2 className="text-lg font-semibold text-[var(--text-primary)]">
              Router Already Initialized
            </h2>
            <p className="mt-1 text-sm text-[var(--text-secondary)]">
              The administrator password has already been provisioned for this instance.
            </p>
          </div>
        </div>

        <EmberRule className="my-5" />

        <div className="flex flex-wrap items-center gap-3">
          <Button asChild>
            <Link to="/dashboard/">Open Dashboard</Link>
          </Button>
          <Button asChild variant="outline">
            <Link to="/login/">Go to Sign In</Link>
          </Button>
        </div>
      </Card>
    );
  }

  const formSummaryError =
    errors.password?.message ??
    errors.confirmPassword?.message ??
    initializeMutation.error?.message ??
    null;

  return (
    <Card className="p-6 sm:p-8">
      <form onSubmit={handleSubmit(onValidSubmit)} noValidate>
        <div className="flex items-center gap-2 text-xs font-medium uppercase tracking-widest text-[var(--accent-text)]">
          <Lock className="h-3.5 w-3.5" aria-hidden="true" />
          <span>First-Run Provisioning</span>
        </div>
        <h2 className="mt-2 text-xl font-semibold tracking-tight text-[var(--text-primary)]">
          Set Administrator Password
        </h2>
        <p className="mt-1.5 text-sm text-[var(--text-secondary)]">
          This credential secures control plane operations and is hashed locally in SQLite via
          bcrypt.
        </p>

        <EmberRule className="my-5" />

        {initStatusQuery.isError ? (
          <Alert variant="destructive" className="mb-5">
            <AlertCircle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
              aria-hidden="true"
            />
            <div>
              <AlertTitle>Unable to verify current init posture</AlertTitle>
              <AlertDescription className="font-mono">
                {initStatusQuery.error.message}
              </AlertDescription>
            </div>
          </Alert>
        ) : null}

        {formSummaryError ? (
          <Alert variant="destructive" className="mb-5">
            <AlertCircle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
              aria-hidden="true"
            />
            <span>{formSummaryError}</span>
          </Alert>
        ) : null}

        <div className="space-y-4">
          <div>
            <Label htmlFor="init-password">Administrator Password</Label>
            <div className="relative mt-1.5">
              <Input
                id="init-password"
                type={showPassword ? "text" : "password"}
                autoComplete="new-password"
                aria-invalid={Boolean(errors.password)}
                placeholder={`Minimum ${MIN_PASSWORD_RUNES} characters`}
                className="pr-10"
                {...passwordField}
              />
              <Button
                type="button"
                variant="ghost"
                size="icon-sm"
                onClick={() => setShowPassword((prev) => !prev)}
                aria-label={showPassword ? "Hide password" : "Show password"}
                className="absolute right-1.5 top-1/2 -translate-y-1/2"
              >
                {showPassword ? (
                  <EyeOff className="h-4 w-4" aria-hidden="true" />
                ) : (
                  <Eye className="h-4 w-4" aria-hidden="true" />
                )}
              </Button>
            </div>
          </div>

          <div>
            <Label htmlFor="init-password-confirm">Confirm Password</Label>
            <Input
              id="init-password-confirm"
              type={showPassword ? "text" : "password"}
              autoComplete="new-password"
              aria-invalid={Boolean(errors.confirmPassword)}
              placeholder="Re-enter password"
              className="mt-1.5"
              {...confirmPasswordField}
            />
          </div>

          <ul className="space-y-1.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3.5 py-3 text-xs">
            <li
              className={`flex items-center justify-between font-mono ${
                validation.meetsMinLength ? "text-[var(--status-ok)]" : "text-[var(--text-muted)]"
              }`}
            >
              <span>At least {MIN_PASSWORD_RUNES} characters</span>
              <span className="tabular-nums">
                {validation.meetsMinLength ? "OK" : `${validation.runeCount}/${MIN_PASSWORD_RUNES}`}
              </span>
            </li>
            <li
              className={`flex items-center justify-between font-mono ${
                validation.withinMaxBytes ? "text-[var(--status-ok)]" : "text-[var(--status-error)]"
              }`}
            >
              <span>Maximum {MAX_PASSWORD_BYTES} UTF-8 bytes</span>
              <span className="tabular-nums">
                {validation.byteLength}/{MAX_PASSWORD_BYTES}
              </span>
            </li>
            <li
              className={`flex items-center justify-between font-mono ${
                validation.passwordsMatch ? "text-[var(--status-ok)]" : "text-[var(--text-muted)]"
              }`}
            >
              <span>Confirmation matches</span>
              <span>{validation.passwordsMatch ? "OK" : "PENDING"}</span>
            </li>
          </ul>
        </div>

        <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <Button type="submit" disabled={initializeMutation.isPending} className="px-5">
            {initializeMutation.isPending ? "Initializing..." : "Initialize Router"}
          </Button>

          <Link
            to="/dashboard/"
            className="text-center text-xs font-medium text-[var(--text-secondary)] transition-colors hover:text-[var(--text-primary)]"
          >
            Return to Dashboard
          </Link>
        </div>
      </form>
    </Card>
  );
}
