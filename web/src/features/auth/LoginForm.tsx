"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertCircle,
  ArrowRight,
  ChevronDown,
  CornerDownLeft,
  Eye,
  EyeOff,
  HelpCircle,
  Loader2,
  Lock,
  ShieldCheck,
} from "lucide-react";
import { type KeyboardEvent, useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";
import { cn } from "@/shared/lib/utils";
import { useUIStore } from "@/shared/store/uiStore";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { useLoginMutation } from "./hooks";
import { formatAuthError, getLoginCopy } from "./i18n";
import { type LoginFormValues, loginSchema } from "./schemas";

interface LoginFormProps {
  remoteEnforced?: boolean;
}

export function LoginForm({ remoteEnforced = false }: LoginFormProps) {
  const navigate = useNavigate();
  const locale = useUIStore((state) => state.locale);
  const copy = getLoginCopy(locale);
  const loginMutation = useLoginMutation();
  const [showPassword, setShowPassword] = useState(false);
  const [capsLockOn, setCapsLockOn] = useState(false);
  const [helpOpen, setHelpOpen] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      password: "",
    },
    mode: "onSubmit",
    reValidateMode: "onChange",
  });

  const clearMutationError = () => {
    if (loginMutation.error) {
      loginMutation.reset();
    }
  };

  const passwordField = register("password", {
    onChange: clearMutationError,
  });

  const handlePasswordKeyEvent = (event: KeyboardEvent<HTMLInputElement>) => {
    if (typeof event.getModifierState === "function") {
      setCapsLockOn(event.getModifierState("CapsLock"));
    }
  };

  const onValidSubmit = (values: LoginFormValues) => {
    loginMutation.mutate(
      { password: values.password },
      {
        onSuccess: () => {
          void navigate("/dashboard/", { replace: true });
        },
      },
    );
  };

  const rawError = errors.password?.message ?? loginMutation.error?.message ?? null;
  const formattedError = formatAuthError(rawError, locale);
  const hasError = Boolean(formattedError);

  return (
    <form onSubmit={handleSubmit(onValidSubmit)} noValidate className="space-y-5">
      {remoteEnforced ? (
        <Alert className="rounded-lg border-[var(--color-info-border)] bg-[var(--color-info-bg)]">
          <ShieldCheck
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--color-info-fg)]"
            aria-hidden="true"
          />
          <div>
            <AlertTitle className="text-xs font-semibold text-[var(--text-primary)]">
              {copy.remoteEnforcedTitle}
            </AlertTitle>
            <AlertDescription className="text-xs leading-relaxed text-[var(--text-secondary)]">
              {copy.remoteEnforcedBody}
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      {formattedError ? (
        <Alert variant="destructive" className="rounded-lg">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div className="min-w-0 space-y-1">
            <AlertTitle className="text-xs font-semibold">{formattedError.title}</AlertTitle>
            <AlertDescription className="text-xs leading-relaxed text-[var(--text-primary)]">
              {formattedError.message}
            </AlertDescription>
            {formattedError.detail ? (
              <p className="truncate font-mono text-[11px] text-[var(--text-muted)]">
                {copy.technicalDetailLabel}: {formattedError.detail}
              </p>
            ) : null}
          </div>
        </Alert>
      ) : null}

      <div className="space-y-2">
        <div className="flex items-center justify-between gap-2">
          <Label
            htmlFor="login-password"
            className="text-xs font-medium text-[var(--text-primary)]"
          >
            {copy.passwordLabel}
          </Label>
          {capsLockOn ? (
            <span
              role="status"
              className="inline-flex items-center gap-1 rounded border border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] px-2 py-0.5 font-mono text-[10px] font-medium text-[var(--status-warn)]"
            >
              {copy.capsLockWarning}
            </span>
          ) : null}
        </div>

        <div className="relative">
          <Lock
            className="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-[var(--text-muted)]"
            aria-hidden="true"
          />
          <Input
            id="login-password"
            type={showPassword ? "text" : "password"}
            autoComplete="current-password"
            autoFocus
            aria-invalid={hasError}
            aria-describedby={hasError ? "login-password-error-hint" : undefined}
            placeholder={copy.passwordPlaceholder}
            onKeyDown={handlePasswordKeyEvent}
            onKeyUp={handlePasswordKeyEvent}
            className={cn(
              "h-11 rounded-md pl-10 pr-11 text-sm",
              hasError && "border-[var(--status-error-border)] focus:border-[var(--status-error)]",
            )}
            {...passwordField}
          />
          <Button
            type="button"
            variant="ghost"
            size="icon-sm"
            onClick={() => setShowPassword((prev) => !prev)}
            aria-label={showPassword ? copy.hidePasswordLabel : copy.showPasswordLabel}
            aria-pressed={showPassword}
            className="absolute right-1.5 top-1/2 h-8 w-8 -translate-y-1/2 rounded-md"
          >
            {showPassword ? (
              <EyeOff className="h-4 w-4" aria-hidden="true" />
            ) : (
              <Eye className="h-4 w-4" aria-hidden="true" />
            )}
          </Button>
        </div>
        {hasError ? (
          <span id="login-password-error-hint" className="sr-only">
            {formattedError?.message}
          </span>
        ) : null}
      </div>

      <div className="space-y-3 pt-1">
        <Button
          type="submit"
          disabled={loginMutation.isPending}
          className="h-11 w-full rounded-md text-sm font-semibold shadow-xs"
        >
          {loginMutation.isPending ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
              <span>{copy.submittingLabel}</span>
            </>
          ) : (
            <>
              <span>{copy.submitLabel}</span>
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </>
          )}
        </Button>

        <div className="flex items-center justify-between gap-2 text-xs text-[var(--text-muted)]">
          <button
            type="button"
            onClick={() => setHelpOpen((prev) => !prev)}
            aria-expanded={helpOpen}
            aria-controls="login-help-panel"
            className="inline-flex items-center gap-1.5 rounded py-1 text-left transition-colors hover:text-[var(--text-primary)]"
          >
            <HelpCircle className="h-3.5 w-3.5 shrink-0" aria-hidden="true" />
            <span>{copy.helpTitle}</span>
            <ChevronDown
              className={cn(
                "h-3.5 w-3.5 shrink-0 transition-transform duration-150",
                helpOpen && "rotate-180",
              )}
              aria-hidden="true"
            />
          </button>

          <span className="hidden items-center gap-1 font-mono text-[11px] sm:inline-flex">
            <CornerDownLeft className="h-3 w-3" aria-hidden="true" />
            <span>{copy.enterHint}</span>
          </span>
        </div>

        {helpOpen ? (
          <div
            id="login-help-panel"
            className="space-y-1.5 rounded-lg border border-[var(--border-subtle)] bg-[var(--bg-canvas)]/80 p-3.5 text-xs leading-relaxed text-[var(--text-secondary)]"
          >
            <p>{copy.helpBody}</p>
            <p className="text-[var(--text-muted)]">{copy.helpRecovery}</p>
          </div>
        ) : null}
      </div>
    </form>
  );
}
