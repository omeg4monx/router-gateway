interface EmberRuleProps {
  className?: string;
}

export function EmberRule({ className = "" }: EmberRuleProps) {
  return (
    <div
      role="separator"
      aria-hidden="true"
      className={`h-px w-full bg-gradient-to-r from-transparent via-[#E56A4A] to-transparent ${className}`.trim()}
    />
  );
}
