import { Loader2 } from "lucide-react";
import { cn } from "@/shared/lib/utils";

const SPINNER_SIZE_CLASS = {
  xs: "h-3.5 w-3.5",
  sm: "h-4 w-4",
  md: "h-5 w-5",
  lg: "h-6 w-6",
} as const;

export type SpinnerSize = keyof typeof SPINNER_SIZE_CLASS;

export interface SpinnerProps {
  size?: SpinnerSize;
  label?: string;
  decorative?: boolean;
  className?: string;
}

export function Spinner({ size = "sm", label, decorative = false, className }: SpinnerProps) {
  const icon = (
    <Loader2
      className={cn(
        "shrink-0 animate-spin motion-reduce:animate-none",
        SPINNER_SIZE_CLASS[size],
        className,
      )}
      aria-hidden="true"
    />
  );

  if (decorative || !label) {
    return icon;
  }

  return (
    <span role="status" aria-live="polite" className="inline-flex items-center gap-2">
      {icon}
      <span className="sr-only">{label}</span>
    </span>
  );
}
