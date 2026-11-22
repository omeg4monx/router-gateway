"use client";

import { useEffect, useState } from "react";
import { Moon, Sun } from "lucide-react";
import { InitForm } from "@/features/system/InitForm";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { EmberRule } from "@/shared/ui/EmberRule";

export default function InitPage() {
  const [theme, setTheme] = useState<"dark" | "light">("dark");

  useEffect(() => {
    const current = document.documentElement.getAttribute("data-theme");
    if (current === "light" || current === "dark") {
      setTheme(current);
    }
  }, []);

  const toggleTheme = () => {
    const next = theme === "dark" ? "light" : "dark";
    setTheme(next);
    document.documentElement.setAttribute("data-theme", next);
    try {
      window.localStorage.setItem("omeg4mon-theme", next);
    } catch {
      // Ignore storage failures.
    }
  };

  return (
    <div className="flex min-h-screen flex-col bg-[var(--bg-canvas)] text-[var(--text-primary)]">
      <header className="bg-[var(--bg-surface)]">
        <div className="mx-auto flex h-16 w-full max-w-4xl items-center justify-between px-4 sm:px-6">
          <div className="flex items-center gap-3">
            <CaduceusLogo size={26} />
            <div className="leading-tight">
              <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
                OMEG4MON ROUTER
              </p>
              <p className="font-mono text-[11px] text-[var(--text-muted)]">
                INITIAL CONFIGURATION
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={toggleTheme}
            className="inline-flex h-9 items-center gap-2 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-3 text-xs font-medium text-[var(--text-secondary)] transition-colors hover:border-[var(--border-strong)] hover:text-[var(--text-primary)]"
            aria-label={theme === "dark" ? "Switch to light theme" : "Switch to dark theme"}
          >
            {theme === "dark" ? (
              <>
                <Sun className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
                <span>Light</span>
              </>
            ) : (
              <>
                <Moon className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
                <span>Dark</span>
              </>
            )}
          </button>
        </div>
        <EmberRule />
      </header>

      <main className="mx-auto flex w-full max-w-lg flex-1 flex-col justify-center px-4 py-10 sm:px-6">
        <InitForm />
      </main>
    </div>
  );
}
