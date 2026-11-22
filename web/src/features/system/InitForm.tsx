"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";
import { AlertCircle, CheckCircle2, Eye, EyeOff, Lock } from "lucide-react";
import { EmberRule } from "@/shared/ui/EmberRule";
import { useInitializeSystem, useInitStatus } from "./hooks";
import {
  MAX_PASSWORD_BYTES,
  MIN_PASSWORD_RUNES,
  validateInitPassword,
} from "./password";

export function InitForm() {
  const router = useRouter();
  const initStatusQuery = useInitStatus();
  const initializeMutation = useInitializeSystem();

  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [clientError, setClientError] = useState<string | null>(null);

  const validation = validateInitPassword(password, confirmPassword);

  const clearErrors = () => {
    if (clientError) {
      setClientError(null);
    }
    if (initializeMutation.error) {
      initializeMutation.reset();
    }
  };

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setClientError(null);

    if (validation.error) {
      setClientError(validation.error);
      return;
    }

    initializeMutation.mutate(
      { password },
      {
        onSuccess: () => {
          router.replace("/dashboard/");
        },
      },
    );
  };

  if (initStatusQuery.isPending) {
    return (
      <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-6">
        <p className="font-mono text-xs uppercase tracking-widest text-[var(--text-muted)]">
          Checking initialization state
        </p>
        <p className="mt-2 text-sm text-[var(--text-secondary)]">
          Querying <span className="font-mono">GET /api/init</span>...
        </p>
      </div>
    );
  }

  if (initStatusQuery.data?.initialized) {
    return (
      <div className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-6">
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
              The administrator password has already been provisioned for this
              instance.
            </p>
          </div>
        </div>

        <EmberRule className="my-5" />

        <div className="flex flex-wrap items-center gap-3">
          <Link
            href="/dashboard/"
            className="inline-flex h-10 items-center justify-center rounded-md bg-[#E56A4A] px-4 text-sm font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)]"
          >
            Open Dashboard
          </Link>
          <Link
            href="/login/"
            className="inline-flex h-10 items-center justify-center rounded-md border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] px-4 text-sm font-medium text-[var(--text-primary)] transition-colors hover:border-[#E56A4A]"
          >
            Go to Sign In
          </Link>
        </div>
      </div>
    );
  }

  const displayedError = clientError ?? initializeMutation.error?.message ?? null;

  return (
    <form
      onSubmit={handleSubmit}
      noValidate
      className="rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface)] p-6 sm:p-8"
    >
      <div className="flex items-center gap-2 text-xs font-medium uppercase tracking-widest text-[var(--accent-text)]">
        <Lock className="h-3.5 w-3.5" aria-hidden="true" />
        <span>First-Run Provisioning</span>
      </div>
      <h2 className="mt-2 text-xl font-semibold tracking-tight text-[var(--text-primary)]">
        Set Administrator Password
      </h2>
      <p className="mt-1.5 text-sm text-[var(--text-secondary)]">
        This credential secures control plane operations and is hashed locally in
        SQLite via bcrypt.
      </p>

      <EmberRule className="my-5" />

      {initStatusQuery.isError ? (
        <div
          role="alert"
          className="mb-5 flex items-start gap-2.5 rounded-md border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-3.5 py-3 text-sm text-[var(--text-primary)]"
        >
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div>
            <p className="font-medium">Unable to verify current init posture</p>
            <p className="mt-0.5 font-mono text-xs text-[var(--text-secondary)]">
              {initStatusQuery.error.message}
            </p>
          </div>
        </div>
      ) : null}

      {displayedError ? (
        <div
          role="alert"
          className="mb-5 flex items-start gap-2.5 rounded-md border border-[var(--status-error-border)] bg-[var(--status-error-bg)] px-3.5 py-3 text-sm text-[var(--text-primary)]"
        >
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <span>{displayedError}</span>
        </div>
      ) : null}

      <div className="space-y-4">
        <div>
          <label
            htmlFor="init-password"
            className="block text-xs font-medium uppercase tracking-wider text-[var(--text-secondary)]"
          >
            Administrator Password
          </label>
          <div className="relative mt-1.5">
            <input
              id="init-password"
              name="password"
              type={showPassword ? "text" : "password"}
              autoComplete="new-password"
              required
              value={password}
              onChange={(event) => {
                setPassword(event.target.value);
                clearErrors();
              }}
              placeholder={`Minimum ${MIN_PASSWORD_RUNES} characters`}
              className="h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 pr-10 font-mono text-sm text-[var(--text-primary)] placeholder:font-sans placeholder:text-[var(--text-muted)] transition-colors focus:border-[#E56A4A]"
            />
            <button
              type="button"
              onClick={() => setShowPassword((prev) => !prev)}
              aria-label={showPassword ? "Hide password" : "Show password"}
              className="absolute right-2 top-1/2 inline-flex h-7 w-7 -translate-y-1/2 items-center justify-center rounded text-[var(--text-muted)] transition-colors hover:text-[var(--text-primary)]"
            >
              {showPassword ? (
                <EyeOff className="h-4 w-4" aria-hidden="true" />
              ) : (
                <Eye className="h-4 w-4" aria-hidden="true" />
              )}
            </button>
          </div>
        </div>

        <div>
          <label
            htmlFor="init-password-confirm"
            className="block text-xs font-medium uppercase tracking-wider text-[var(--text-secondary)]"
          >
            Confirm Password
          </label>
          <input
            id="init-password-confirm"
            name="confirmPassword"
            type={showPassword ? "text" : "password"}
            autoComplete="new-password"
            required
            value={confirmPassword}
            onChange={(event) => {
              setConfirmPassword(event.target.value);
              clearErrors();
            }}
            placeholder="Re-enter password"
            className="mt-1.5 h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 font-mono text-sm text-[var(--text-primary)] placeholder:font-sans placeholder:text-[var(--text-muted)] transition-colors focus:border-[#E56A4A]"
          />
        </div>

        {/* Clear inline validation checklist */}
        <ul className="space-y-1.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3.5 py-3 text-xs">
          <li
            className={`flex items-center justify-between font-mono ${
              validation.meetsMinLength
                ? "text-[var(--status-ok)]"
                : "text-[var(--text-muted)]"
            }`}
          >
            <span>At least {MIN_PASSWORD_RUNES} characters</span>
            <span className="tabular-nums">
              {validation.meetsMinLength
                ? "OK"
                : `${validation.runeCount}/${MIN_PASSWORD_RUNES}`}
            </span>
          </li>
          <li
            className={`flex items-center justify-between font-mono ${
              validation.withinMaxBytes
                ? "text-[var(--status-ok)]"
                : "text-[var(--status-error)]"
            }`}
          >
            <span>Maximum {MAX_PASSWORD_BYTES} UTF-8 bytes</span>
            <span className="tabular-nums">
              {validation.byteLength}/{MAX_PASSWORD_BYTES}
            </span>
          </li>
          <li
            className={`flex items-center justify-between font-mono ${
              validation.passwordsMatch
                ? "text-[var(--status-ok)]"
                : "text-[var(--text-muted)]"
            }`}
          >
            <span>Confirmation matches</span>
            <span>{validation.passwordsMatch ? "OK" : "PENDING"}</span>
          </li>
        </ul>
      </div>

      <div className="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <button
          type="submit"
          disabled={initializeMutation.isPending}
          className="inline-flex h-10 items-center justify-center rounded-md bg-[#E56A4A] px-5 text-sm font-semibold text-[#0a0a0a] transition-colors hover:bg-[var(--accent-primary-hover)] disabled:opacity-60"
        >
          {initializeMutation.isPending
            ? "Initializing..."
            : "Initialize Router"}
        </button>

        <Link
          href="/dashboard/"
          className="text-center text-xs font-medium text-[var(--text-secondary)] transition-colors hover:text-[var(--text-primary)]"
        >
          Return to Dashboard
        </Link>
      </div>
    </form>
  );
}
