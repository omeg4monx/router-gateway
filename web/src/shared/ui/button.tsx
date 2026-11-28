import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";
import { cn } from "@/shared/lib/utils";
import { Spinner } from "./spinner";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors disabled:pointer-events-none disabled:opacity-60 [&_svg]:pointer-events-none [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default:
          "bg-[var(--accent-primary)] font-semibold text-[var(--accent-primary-ink)] hover:bg-[var(--accent-primary-hover)]",
        secondary:
          "border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-[var(--text-secondary)] hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]",
        outline:
          "border border-[var(--border-strong)] bg-[var(--bg-surface-raised)] text-[var(--text-primary)] hover:border-[var(--accent-primary)]",
        ghost:
          "text-[var(--text-muted)] hover:bg-[var(--bg-surface-raised)] hover:text-[var(--text-primary)]",
        destructive:
          "border border-[var(--status-error-border)] bg-[var(--bg-surface)] text-[var(--text-primary)] hover:border-[var(--accent-primary)]",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 rounded-md px-3.5 text-xs",
        xs: "h-8 rounded px-3 text-xs",
        icon: "h-9 w-9",
        "icon-sm": "h-7 w-7 rounded",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
  loading?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    { className, variant, size, asChild = false, loading = false, disabled, children, ...props },
    ref,
  ) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        disabled={asChild ? disabled : disabled || loading}
        aria-busy={loading || undefined}
        {...props}
      >
        {asChild ? (
          children
        ) : (
          <>
            {loading ? <Spinner size="xs" decorative /> : null}
            {children}
          </>
        )}
      </Comp>
    );
  },
);
Button.displayName = "Button";

export { Button, buttonVariants };
