import { Suspense } from "react";
import { ProviderDetailView } from "@/features/providers/ProviderDetailView";
import { Card, CardContent } from "@/shared/ui/card";

export default function DashboardProviderDetailPage() {
  return (
    <Suspense
      fallback={
        <Card>
          <CardContent className="text-sm text-[var(--text-muted)]">
            Loading provider connection details...
          </CardContent>
        </Card>
      }
    >
      <ProviderDetailView />
    </Suspense>
  );
}
