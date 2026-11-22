import type { ReactNode } from "react";
import { AuthGuard } from "@/features/auth/AuthGuard";
import { AppShell } from "@/shared/ui/AppShell";

interface DashboardLayoutProps {
  children: ReactNode;
}

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <AuthGuard>
      <AppShell>{children}</AppShell>
    </AuthGuard>
  );
}
