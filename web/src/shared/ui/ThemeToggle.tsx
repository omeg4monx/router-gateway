"use client";

import { Moon, Sun } from "lucide-react";
import { useEffect } from "react";
import { useUIStore } from "@/shared/store/uiStore";
import { Button } from "./button";

export function ThemeToggle() {
  const theme = useUIStore((state) => state.theme);
  const hydrateTheme = useUIStore((state) => state.hydrateTheme);
  const toggleTheme = useUIStore((state) => state.toggleTheme);

  useEffect(() => {
    hydrateTheme();
  }, [hydrateTheme]);

  return (
    <Button
      type="button"
      variant="secondary"
      size="sm"
      onClick={toggleTheme}
      aria-label={theme === "dark" ? "Switch to warm light theme" : "Switch to charcoal dark theme"}
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
    </Button>
  );
}
