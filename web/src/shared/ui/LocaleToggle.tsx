import { Languages } from "lucide-react";
import { useEffect } from "react";
import { cn } from "@/shared/lib/utils";
import { type UILocale, useUIStore } from "@/shared/store/uiStore";

const LOCALE_OPTIONS: ReadonlyArray<{ code: UILocale; label: string; name: string }> = [
  { code: "en", label: "EN", name: "English" },
  { code: "id", label: "ID", name: "Bahasa Indonesia" },
];

export function LocaleToggle() {
  const locale = useUIStore((state) => state.locale);
  const hydrateLocale = useUIStore((state) => state.hydrateLocale);
  const setLocale = useUIStore((state) => state.setLocale);

  useEffect(() => {
    hydrateLocale();
  }, [hydrateLocale]);

  return (
    <fieldset
      aria-label={locale === "id" ? "Pilih bahasa tampilan" : "Select display language"}
      className="m-0 inline-flex h-9 min-w-0 items-center gap-0.5 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] p-0.5"
    >
      <span
        className="hidden items-center px-1.5 text-[var(--text-muted)] sm:flex"
        aria-hidden="true"
      >
        <Languages className="h-3.5 w-3.5" />
      </span>
      {LOCALE_OPTIONS.map((option) => {
        const active = locale === option.code;
        return (
          <button
            key={option.code}
            type="button"
            onClick={() => setLocale(option.code)}
            aria-pressed={active}
            aria-label={option.name}
            className={cn(
              "inline-flex h-7 min-w-7 items-center justify-center rounded px-1.5 font-mono text-[11px] font-semibold tracking-wider transition-colors sm:min-w-8 sm:px-2",
              active
                ? "bg-[var(--accent-primary)] text-[var(--accent-primary-ink)] shadow-xs"
                : "text-[var(--text-muted)] hover:text-[var(--text-primary)]",
            )}
          >
            {option.label}
          </button>
        );
      })}
    </fieldset>
  );
}
