import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/shared/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded border px-2 py-0.5 font-mono text-[10px] font-medium uppercase tracking-wider",
  {
    variants: {
      variant: {
        default:
          "border-[#E56A4A] bg-[var(--bg-surface-raised)] text-[var(--accent-text)]",
        secondary:
          "border-[var(--border-subtle)] bg-[var(--bg-canvas)] text-[var(--text-muted)]",
        outline:
          "border-[var(--border-subtle)] bg-[var(--bg-surface)] text-[var(--text-muted)]",
        success:
          "border-[var(--status-ok-border)] bg-[var(--status-ok-bg)] text-[var(--status-ok)]",
        warning:
          "border-[var(--status-warn-border)] bg-[var(--status-warn-bg)] text-[var(--status-warn)]",
        destructive:
          "border-[var(--status-error-border)] bg-[var(--status-error-bg)] text-[var(--status-error)]",
      },
    },
    defaultVariants: {
      variant: "secondary",
    },
  },
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <span className={cn(badgeVariants({ variant }), className)} {...props} />
  );
}

export { Badge, badgeVariants };
