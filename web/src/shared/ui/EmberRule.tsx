interface EmberRuleProps {
  className?: string;
}

export function EmberRule({ className = "" }: EmberRuleProps) {
  return (
    <hr
      className={`h-px w-full border-0 bg-gradient-to-r from-transparent via-[var(--accent-primary)] to-transparent ${className}`.trim()}
    />
  );
}
