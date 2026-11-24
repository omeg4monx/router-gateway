import { AlertTriangle, RefreshCw } from "lucide-react";
import { Component, type ErrorInfo, type ReactNode } from "react";
import { Alert, AlertDescription, AlertTitle } from "./alert";
import { Button } from "./button";
import { CaduceusLogo } from "./CaduceusLogo";
import { Card } from "./card";
import { EmberRule } from "./EmberRule";

export function isChunkLoadError(error: unknown): boolean {
  if (!(error instanceof Error)) {
    return false;
  }
  const msg = error.message.toLowerCase();
  return (
    msg.includes("failed to fetch dynamically imported module") ||
    msg.includes("importing a module script failed") ||
    msg.includes("loading chunk") ||
    msg.includes("dynamically imported module")
  );
}

interface ErrorBoundaryProps {
  children: ReactNode;
  resetKey?: string;
  compact?: boolean;
}

interface ErrorBoundaryState {
  error: Error | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = { error: null };

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo): void {
    console.error("[ErrorBoundary]", error, info.componentStack);
  }

  componentDidUpdate(prevProps: ErrorBoundaryProps): void {
    if (this.state.error && prevProps.resetKey !== this.props.resetKey) {
      this.setState({ error: null });
    }
  }

  private handleReset = (): void => {
    this.setState({ error: null });
  };

  private handleReload = (): void => {
    window.location.reload();
  };

  render(): ReactNode {
    const { error } = this.state;
    if (!error) {
      return this.props.children;
    }

    const chunkError = isChunkLoadError(error);
    const title = chunkError ? "Application update required" : "Unexpected interface error";
    const description = chunkError
      ? "A lazy-loaded view chunk could not be fetched (the server binary may have been rebuilt). Reload to load the latest bundle."
      : error.message || "An unhandled error occurred while rendering this view.";

    const card = (
      <Card className="w-full max-w-lg rounded-xl border-[var(--border-strong)]/80 p-6 shadow-sm sm:p-8">
        <div className="flex items-center gap-3">
          <CaduceusLogo size={28} />
          <div>
            <p className="text-sm font-semibold tracking-wider text-[var(--text-primary)]">
              OMEG4MON ROUTER
            </p>
            <p className="font-mono text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
              UI Recovery Boundary
            </p>
          </div>
        </div>

        <EmberRule className="my-5" />

        <div className="space-y-4">
          <Alert variant="destructive" className="rounded-lg">
            <AlertTriangle
              className="mt-0.5 h-4 w-4 shrink-0 text-[var(--status-error)]"
              aria-hidden="true"
            />
            <div>
              <AlertTitle className="text-xs font-semibold">{title}</AlertTitle>
              <AlertDescription className="font-mono text-xs">{description}</AlertDescription>
            </div>
          </Alert>

          <div className="flex flex-wrap items-center gap-2.5">
            <Button type="button" size="sm" onClick={this.handleReload}>
              <RefreshCw className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Reload Page</span>
            </Button>
            {!chunkError && (
              <Button type="button" variant="outline" size="sm" onClick={this.handleReset}>
                Try Again
              </Button>
            )}
          </div>
        </div>
      </Card>
    );

    if (this.props.compact) {
      return <div className="flex items-center justify-center py-12">{card}</div>;
    }

    return (
      <main className="flex min-h-screen items-center justify-center bg-[var(--bg-canvas)] px-4 py-12 text-[var(--text-primary)]">
        {card}
      </main>
    );
  }
}
