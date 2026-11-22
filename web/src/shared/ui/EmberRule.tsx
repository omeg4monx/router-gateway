interface EmberRuleProps {
  className?: string;
}

export function EmberRule({ className = "" }: EmberRuleProps) {
  return (
    <hr
      className={`h-px w-full border-0 bg-gradient-to-r from-transparent via-[#E56A4A] to-transparent ${className}`.trim()}
    />
  );
}
