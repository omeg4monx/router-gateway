import * as React from "react";
import { cn } from "@/shared/lib/utils";

export type InputProps = React.InputHTMLAttributes<HTMLInputElement>;

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => (
    <input
      type={type}
      ref={ref}
      className={cn(
        "flex h-10 w-full rounded-md border border-[var(--border-strong)] bg-[var(--bg-input)] px-3 py-2 font-mono text-sm text-[var(--text-primary)] placeholder:font-sans placeholder:text-[var(--text-muted)] transition-colors focus:border-[var(--accent-primary)] disabled:cursor-not-allowed disabled:border-[var(--border-subtle)] disabled:bg-[var(--bg-canvas)] disabled:text-[var(--text-muted)]",
        className,
      )}
      {...props}
    />
  ),
);
Input.displayName = "Input";

export { Input };
