import type { ReactNode } from "react";
import { AppShell } from "@/shared/ui/AppShell";

interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return <AppShell>{children}</AppShell>;
}
