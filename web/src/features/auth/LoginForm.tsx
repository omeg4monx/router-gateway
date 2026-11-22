"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { AlertCircle, ArrowRight, Eye, EyeOff, ShieldCheck } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Button } from "@/shared/ui/button";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import { useLoginMutation } from "./hooks";
import { type LoginFormValues, loginSchema } from "./schemas";

interface LoginFormProps {
  remoteEnforced?: boolean;
}

export function LoginForm({ remoteEnforced = false }: LoginFormProps) {
  const router = useRouter();
  const loginMutation = useLoginMutation();
  const [showPassword, setShowPassword] = useState(false);

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

  const onValidSubmit = (values: LoginFormValues) => {
    loginMutation.mutate(
      { password: values.password },
      {
        onSuccess: () => {
          router.replace("/dashboard/");
        },
      },
    );
  };

  const errorMessage = errors.password?.message ?? loginMutation.error?.message ?? null;

  return (
    <form onSubmit={handleSubmit(onValidSubmit)} noValidate className="space-y-5">
      {remoteEnforced ? (
        <Alert>
          <ShieldCheck className="mt-0.5 h-4 w-4 shrink-0 text-[#E56A4A]" aria-hidden="true" />
          <div>
            <AlertTitle>Remote Bind Enforcement Active</AlertTitle>
            <AlertDescription>
              Non-loopback control plane access requires a signed operator session cookie.
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      {errorMessage ? (
        <Alert variant="destructive">
          <AlertCircle
            className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
            aria-hidden="true"
          />
          <div>
            <AlertTitle>Authentication failed</AlertTitle>
            <AlertDescription className="font-mono">{errorMessage}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      <div>
        <Label htmlFor="login-password">Administrator Password</Label>
        <div className="relative mt-1.5">
          <Input
            id="login-password"
            type={showPassword ? "text" : "password"}
            autoComplete="current-password"
            autoFocus
            aria-invalid={Boolean(errors.password || loginMutation.isError)}
            placeholder="Enter administrator password"
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

      <Button type="submit" disabled={loginMutation.isPending} className="w-full">
        <span>{loginMutation.isPending ? "Signing in..." : "Sign In"}</span>
        {!loginMutation.isPending ? <ArrowRight className="h-4 w-4" aria-hidden="true" /> : null}
      </Button>
    </form>
  );
}
