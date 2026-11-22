"use client";

import { InitForm } from "@/features/system/InitForm";
import { CaduceusLogo } from "@/shared/ui/CaduceusLogo";
import { EmberRule } from "@/shared/ui/EmberRule";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

export default function InitPage() {
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

          <ThemeToggle />
        </div>
        <EmberRule />
      </header>

      <main className="mx-auto flex w-full max-w-lg flex-1 flex-col justify-center px-4 py-10 sm:px-6">
        <InitForm />
      </main>
    </div>
  );
}
