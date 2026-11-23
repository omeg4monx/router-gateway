"use client";

import { Moon, Sun } from "lucide-react";
import { useEffect } from "react";
import { useUIStore } from "@/shared/store/uiStore";
import { Button } from "./button";

export function ThemeToggle() {
  const theme = useUIStore((state) => state.theme);
  const locale = useUIStore((state) => state.locale);
  const hydrateTheme = useUIStore((state) => state.hydrateTheme);
  const toggleTheme = useUIStore((state) => state.toggleTheme);

  useEffect(() => {
    hydrateTheme();
  }, [hydrateTheme]);

  const isId = locale === "id";
  const lightLabel = isId ? "Terang" : "Light";
  const darkLabel = isId ? "Gelap" : "Dark";

  return (
    <Button
      type="button"
      variant="secondary"
      size="sm"
      onClick={toggleTheme}
      aria-label={
        theme === "dark"
          ? isId
            ? "Ganti ke tema terang Omegamon X-Antibody"
            : "Switch to Omegamon X-Antibody light theme"
          : isId
            ? "Ganti ke tema gelap Omegamon Zwart"
            : "Switch to Omegamon Zwart dark theme"
      }
    >
      {theme === "dark" ? (
        <>
          <Sun className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
          <span>{lightLabel}</span>
        </>
      ) : (
        <>
          <Moon className="h-3.5 w-3.5 text-[var(--accent-primary)]" aria-hidden="true" />
          <span>{darkLabel}</span>
        </>
      )}
    </Button>
  );
}
