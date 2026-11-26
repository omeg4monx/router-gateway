import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  ChevronLeft,
  ChevronRight,
  Download,
  KeyRound,
  Lock,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Trash2,
  X,
} from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card } from "@/shared/ui/card";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClearSecurityAuditMutation,
  useCreateSecurityAuditMutation,
  useSecurityAuditEvents,
  useSecurityPosture,
} from "./hooks";
import {
  AUDIT_PAGE_LIMIT,
  AUDIT_SEVERITIES,
  type CreateSecurityAuditFormValues,
  computeAuditPagination,
  createSecurityAuditSchema,
  parseAuditMetadataJson,
} from "./schemas";

function severityBadgeVariant(
  severity: string,
): "default" | "secondary" | "warning" | "destructive" {
  switch (severity) {
    case "critical":
      return "destructive";
    case "warning":
      return "warning";
    default:
      return "secondary";
  }
}

function formatCompactTimestamp(iso: string): string {
  const date = new Date(iso);
  if (Number.isNaN(date.getTime())) {
    return iso;
  }
  return date.toLocaleString(undefined, {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  });
}

export function SecurityConsole() {
  const [categoryFilter, setCategoryFilter] = useState("");
  const [severityFilter, setSeverityFilter] = useState<"" | "info" | "warning" | "critical">("");
  const [searchInput, setSearchInput] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [offset, setOffset] = useState(0);
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [warningsExpanded, setWarningsExpanded] = useState(false);
  const [confirmClear, setConfirmClear] = useState(false);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      setDebouncedSearch(searchInput.trim());
      setOffset(0);
    }, 250);
    return () => window.clearTimeout(timer);
  }, [searchInput]);

  useEffect(() => {
    if (!createModalOpen) {
      return;
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setCreateModalOpen(false);
      }
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [createModalOpen]);

  const postureQuery = useSecurityPosture();
  const auditQuery = useSecurityAuditEvents({
    category: categoryFilter.trim() || undefined,
    severity: severityFilter || undefined,
    q: debouncedSearch || undefined,
    limit: AUDIT_PAGE_LIMIT,
    offset,
  });
  const createMutation = useCreateSecurityAuditMutation();
  const clearMutation = useClearSecurityAuditMutation();

  const form = useForm<CreateSecurityAuditFormValues>({
    resolver: zodResolver(createSecurityAuditSchema),
    defaultValues: {
      category: "security",
      action: "operator.note",
      severity: "info",
      resource_id: "",
      summary: "",
      metadata_text: "{}",
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    const metadata = parseAuditMetadataJson(values.metadata_text);
    createMutation.mutate(
      {
        category: values.category,
        action: values.action,
        severity: values.severity,
        resource_id: values.resource_id || undefined,
        summary: values.summary,
        metadata,
      },
      {
        onSuccess: () => {
          form.reset();
          setOffset(0);
          setCreateModalOpen(false);
        },
      },
    );
  });

  const handleRefresh = () => {
    void Promise.all([postureQuery.refetch(), auditQuery.refetch()]);
  };

  const handleResetFilters = () => {
    setSearchInput("");
    setDebouncedSearch("");
    setCategoryFilter("");
    setSeverityFilter("");
    setOffset(0);
  };

  const handleClearConfirm = () => {
    clearMutation.mutate(undefined, {
      onSuccess: () => {
        setConfirmClear(false);
        setOffset(0);
      },
    });
  };

  const posture = postureQuery.data;
  const auditData = auditQuery.data;
  const auditItems = auditData?.items ?? [];
  const auditTotal = auditData?.total ?? 0;
  const isRefreshing = postureQuery.isFetching || auditQuery.isFetching;
  const hasActiveFilters = Boolean(searchInput.trim() || categoryFilter.trim() || severityFilter);
  const warningCount = posture?.warnings.length ?? 0;

  const pagination = computeAuditPagination({
    total: auditTotal,
    pageSize: AUDIT_PAGE_LIMIT,
    offset,
    itemCount: auditItems.length,
  });

  return (
    <div className="space-y-3">
      {/* Page header */}
      <section className="flex flex-col gap-2.5 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0">
          <h2 className="text-xl font-semibold tracking-tight text-[var(--text-primary)]">
            Security Posture &amp; Audit
          </h2>
          <p className="text-xs text-[var(--text-secondary)]">
            Live listener posture, vault encryption source, and redacted control-plane audit log.
          </p>
        </div>

        <div className="grid grid-cols-3 gap-2 sm:flex sm:shrink-0 sm:items-center">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isRefreshing}
            className="h-8 justify-center px-3 text-xs"
          >
            <RefreshCw
              className={`h-3.5 w-3.5 text-[var(--accent-primary)] ${isRefreshing ? "animate-spin" : ""}`}
              aria-hidden="true"
            />
            <span>Refresh</span>
          </Button>
          <Button asChild variant="outline" size="sm" className="h-8 justify-center px-3 text-xs">
            <a href="/api/security/audit/export" download="security-audit.json">
              <Download className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Export</span>
            </a>
          </Button>
          <Button
            type="button"
            size="sm"
            onClick={() => setCreateModalOpen(true)}
            className="h-8 justify-center px-3 text-xs"
          >
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Record Event</span>
          </Button>
        </div>
      </section>

      {postureQuery.isError || auditQuery.isError || clearMutation.isError ? (
        <Alert variant="destructive" className="py-2.5">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div className="min-w-0">
            <AlertTitle className="text-xs">Security request failed</AlertTitle>
            <AlertDescription className="truncate font-mono text-xs">
              {postureQuery.error?.message ??
                auditQuery.error?.message ??
                clearMutation.error?.message}
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      {/* Proportional 4-card posture row */}
      <section
        aria-label="Security posture summary"
        className="grid grid-cols-2 gap-2.5 lg:grid-cols-4"
      >
        <Card className="flex items-center gap-3 px-3.5 py-2.5">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
            <ShieldCheck className="h-4 w-4 text-[var(--accent-primary)]" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-medium text-[var(--text-muted)]">Bind Scope</p>
            <p className="truncate font-mono text-sm font-semibold text-[var(--text-primary)]">
              {posture ? posture.bind_addr : "--"}
            </p>
            <p className="truncate text-[11px] text-[var(--text-muted)]">
              {posture ? (posture.is_loopback ? "Loopback only" : "Remote exposed") : "Checking..."}
            </p>
          </div>
        </Card>

        <Card className="flex items-center gap-3 px-3.5 py-2.5">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
            <Lock className="h-4 w-4 text-[var(--text-secondary)]" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-medium text-[var(--text-muted)]">Admin Login</p>
            <p
              className={`font-mono text-sm font-semibold ${
                posture?.require_login ? "text-[var(--status-ok)]" : "text-[var(--status-warn)]"
              }`}
            >
              {posture ? (posture.require_login ? "ENFORCED" : "DISABLED") : "--"}
            </p>
            <p className="truncate text-[11px] text-[var(--text-muted)]">
              {posture?.allow_insecure_remote ? "Insecure override" : "Remote guard active"}
            </p>
          </div>
        </Card>

        <Card className="flex items-center gap-3 px-3.5 py-2.5">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
            <KeyRound className="h-4 w-4 text-[var(--text-secondary)]" aria-hidden="true" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] font-medium text-[var(--text-muted)]">Vault Key</p>
            <p className="font-mono text-sm font-semibold uppercase text-[var(--text-primary)]">
              {posture ? posture.master_key_source : "--"}
            </p>
            <p className="truncate text-[11px] text-[var(--text-muted)]">AES-256-GCM</p>
          </div>
        </Card>

        <Card className="flex items-center gap-3 px-3.5 py-2.5">
          <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]">
            <AlertTriangle
              className={`h-4 w-4 ${
                warningCount > 0 ? "text-[var(--status-warn)]" : "text-[var(--text-secondary)]"
              }`}
              aria-hidden="true"
            />
          </div>
          <div className="min-w-0 flex-1">
            <p className="text-[11px] font-medium text-[var(--text-muted)]">Audit Records</p>
            <p className="font-mono text-sm font-semibold tabular-nums text-[var(--text-primary)]">
              {posture ? posture.audit_event_count.toLocaleString("en-US") : "--"}
            </p>
            {warningCount > 0 ? (
              <button
                type="button"
                onClick={() => setWarningsExpanded((prev) => !prev)}
                className="cursor-pointer font-mono text-[11px] text-[var(--status-warn)] underline-offset-2 hover:underline"
              >
                {warningCount} posture warning{warningCount > 1 ? "s" : ""}
              </button>
            ) : (
              <p className="text-[11px] text-[var(--text-muted)]">No warnings</p>
            )}
          </div>
        </Card>
      </section>

      {warningsExpanded && posture && posture.warnings.length > 0 ? (
        <Alert variant="warning" className="py-2">
          <AlertTriangle
            className="h-4 w-4 shrink-0 text-[var(--status-warn)]"
            aria-hidden="true"
          />
          <div className="min-w-0">
            <AlertTitle className="text-xs">Posture Recommendations</AlertTitle>
            <ul className="mt-1 list-disc space-y-0.5 pl-4 text-xs text-[var(--text-secondary)]">
              {posture.warnings.map((warning) => (
                <li key={warning}>{warning}</li>
              ))}
            </ul>
          </div>
        </Alert>
      ) : null}

      {/* Audit Trail Card — crisp natural height sized for 10 rows per page */}
      <Card className="overflow-hidden">
        {/* Filter toolbar */}
        <div className="flex flex-col gap-2.5 border-b border-[var(--border-subtle)] px-4 py-2.5 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              <h3 className="text-sm font-semibold text-[var(--text-primary)]">Audit Trail</h3>
              <Badge variant="secondary" className="font-mono text-[11px] tabular-nums">
                {auditQuery.isError
                  ? "Error"
                  : auditData
                    ? auditTotal.toLocaleString("en-US")
                    : "..."}
              </Badge>
            </div>

            {confirmClear ? (
              <div className="flex items-center gap-1.5 lg:hidden">
                <Button
                  type="button"
                  variant="destructive"
                  size="xs"
                  onClick={handleClearConfirm}
                  disabled={clearMutation.isPending}
                >
                  {clearMutation.isPending ? "Clearing..." : "Confirm"}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  size="xs"
                  onClick={() => setConfirmClear(false)}
                  disabled={clearMutation.isPending}
                >
                  Cancel
                </Button>
              </div>
            ) : (
              <Button
                type="button"
                variant="outline"
                size="xs"
                onClick={() => setConfirmClear(true)}
                disabled={clearMutation.isPending || auditTotal === 0}
                className="lg:hidden"
                title="Clear audit trail"
              >
                <Trash2 className="h-3.5 w-3.5 text-[var(--status-error)]" aria-hidden="true" />
                <span>Clear</span>
              </Button>
            )}
          </div>

          <div className="grid grid-cols-2 gap-2 sm:flex sm:flex-wrap sm:items-center">
            <div className="relative col-span-2 sm:col-span-1">
              <Search
                className="pointer-events-none absolute left-2.5 top-2 h-3.5 w-3.5 text-[var(--text-muted)]"
                aria-hidden="true"
              />
              <Input
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
                placeholder="Search summary or action..."
                className="h-8 w-full pl-8 text-xs sm:w-56"
                aria-label="Search audit trail"
              />
            </div>
            <Input
              value={categoryFilter}
              onChange={(e) => {
                setCategoryFilter(e.target.value);
                setOffset(0);
              }}
              placeholder="Category"
              className="h-8 w-full text-xs sm:w-32"
              aria-label="Filter by category"
            />
            <select
              value={severityFilter}
              onChange={(e) => {
                setSeverityFilter(e.target.value as "" | "info" | "warning" | "critical");
                setOffset(0);
              }}
              aria-label="Filter by severity"
              className="h-8 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2.5 text-xs text-[var(--text-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)] sm:w-32"
            >
              <option value="">All severities</option>
              {AUDIT_SEVERITIES.map((sev) => (
                <option key={sev} value={sev}>
                  {sev}
                </option>
              ))}
            </select>
            {hasActiveFilters ? (
              <Button
                type="button"
                variant="ghost"
                size="xs"
                onClick={handleResetFilters}
                className="col-span-2 h-8 justify-center sm:col-span-1"
              >
                <X className="h-3.5 w-3.5" aria-hidden="true" />
                <span>Reset</span>
              </Button>
            ) : null}
            <div className="hidden lg:flex lg:items-center lg:gap-1.5">
              {confirmClear ? (
                <>
                  <Button
                    type="button"
                    variant="destructive"
                    size="xs"
                    onClick={handleClearConfirm}
                    disabled={clearMutation.isPending}
                    className="h-8"
                  >
                    {clearMutation.isPending ? "Clearing..." : "Confirm Clear"}
                  </Button>
                  <Button
                    type="button"
                    variant="secondary"
                    size="xs"
                    onClick={() => setConfirmClear(false)}
                    disabled={clearMutation.isPending}
                    className="h-8"
                  >
                    Cancel
                  </Button>
                </>
              ) : (
                <Button
                  type="button"
                  variant="outline"
                  size="xs"
                  onClick={() => setConfirmClear(true)}
                  disabled={clearMutation.isPending || auditTotal === 0}
                  className="h-8"
                  title="Clear audit trail (records an audit.cleared marker)"
                >
                  <Trash2 className="h-3.5 w-3.5 text-[var(--status-error)]" aria-hidden="true" />
                  <span>Clear</span>
                </Button>
              )}
            </div>
          </div>
        </div>

        {/* Mobile-first card list (< md) */}
        <div className="divide-y divide-[var(--border-subtle)] md:hidden">
          {auditItems.length > 0 ? (
            auditItems.map((ev) => (
              <div key={ev.id} className="space-y-1 px-4 py-2.5">
                <div className="flex items-center justify-between gap-2">
                  <div className="flex min-w-0 items-center gap-1.5">
                    <Badge variant={severityBadgeVariant(ev.severity)}>{ev.severity}</Badge>
                    <span className="truncate font-mono text-xs text-[var(--accent-text)]">
                      {ev.category}
                    </span>
                    <span className="text-[var(--text-muted)]">/</span>
                    <span className="truncate font-mono text-xs text-[var(--text-primary)]">
                      {ev.action}
                    </span>
                  </div>
                  <span className="shrink-0 font-mono text-[11px] tabular-nums text-[var(--text-muted)]">
                    {formatCompactTimestamp(ev.created_at)}
                  </span>
                </div>
                <p className="line-clamp-2 text-xs text-[var(--text-primary)]">{ev.summary}</p>
                <div className="flex items-center gap-2 font-mono text-[11px] text-[var(--text-muted)]">
                  <span>{ev.actor}</span>
                  {ev.resource_id ? (
                    <>
                      <span>·</span>
                      <span className="truncate">{ev.resource_id}</span>
                    </>
                  ) : null}
                </div>
              </div>
            ))
          ) : (
            <p className="px-4 py-8 text-center text-sm text-[var(--text-muted)]">
              {auditQuery.isError
                ? "Unable to load security audit trail."
                : auditQuery.isPending
                  ? "Loading security audit trail..."
                  : "No security audit events match the current filter."}
            </p>
          )}
        </div>

        {/* Desktop & tablet table (>= md) */}
        <div className="hidden overflow-x-auto md:block">
          <table className="w-full text-left text-xs">
            <thead className="border-b border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
              <tr>
                <th className="whitespace-nowrap px-4 py-2">Time</th>
                <th className="whitespace-nowrap px-3 py-2">Severity</th>
                <th className="whitespace-nowrap px-3 py-2">Event</th>
                <th className="whitespace-nowrap px-3 py-2">Actor</th>
                <th className="px-4 py-2">Summary</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--border-subtle)]">
              {auditItems.length > 0 ? (
                auditItems.map((ev) => (
                  <tr key={ev.id} className="transition-colors hover:bg-[var(--bg-surface-raised)]">
                    <td
                      className="whitespace-nowrap px-4 py-2 font-mono text-[11px] tabular-nums text-[var(--text-muted)]"
                      title={new Date(ev.created_at).toLocaleString()}
                    >
                      {formatCompactTimestamp(ev.created_at)}
                    </td>
                    <td className="whitespace-nowrap px-3 py-2">
                      <Badge variant={severityBadgeVariant(ev.severity)}>{ev.severity}</Badge>
                    </td>
                    <td
                      className="whitespace-nowrap px-3 py-2 font-mono text-xs"
                      title={`${ev.category} / ${ev.action}`}
                    >
                      <span className="text-[var(--accent-text)]">{ev.category}</span>
                      <span className="text-[var(--text-muted)]"> / </span>
                      <span className="text-[var(--text-primary)]">{ev.action}</span>
                    </td>
                    <td className="whitespace-nowrap px-3 py-2 font-mono text-xs text-[var(--text-secondary)]">
                      {ev.actor}
                    </td>
                    <td
                      className="max-w-md truncate px-4 py-2 text-xs text-[var(--text-primary)] xl:max-w-2xl"
                      title={ev.resource_id ? `${ev.summary} (${ev.resource_id})` : ev.summary}
                    >
                      <span>{ev.summary}</span>
                      {ev.resource_id ? (
                        <span className="ml-2 font-mono text-[11px] text-[var(--text-muted)]">
                          {ev.resource_id}
                        </span>
                      ) : null}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan={5}
                    className="px-4 py-10 text-center text-sm text-[var(--text-muted)]"
                  >
                    {auditQuery.isError
                      ? "Unable to load security audit trail."
                      : auditQuery.isPending
                        ? "Loading security audit trail..."
                        : "No security audit events match the current filter."}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Pagination footer (10 rows / page) */}
        <div className="flex items-center justify-between gap-3 border-t border-[var(--border-subtle)] bg-[var(--bg-surface-raised)]/40 px-4 py-2.5">
          <span className="font-mono text-xs tabular-nums text-[var(--text-muted)]">
            Showing {pagination.showingFrom.toLocaleString("en-US")}–
            {pagination.showingTo.toLocaleString("en-US")} of {auditTotal.toLocaleString("en-US")}{" "}
            events
          </span>

          <div className="flex items-center gap-2">
            <span className="font-mono text-xs tabular-nums text-[var(--text-muted)]">
              Page {pagination.currentPage} of {pagination.totalPages}
            </span>
            <Button
              type="button"
              variant="secondary"
              size="xs"
              disabled={!pagination.canPrev || auditQuery.isFetching}
              onClick={() => setOffset((prev) => Math.max(0, prev - AUDIT_PAGE_LIMIT))}
              aria-label="Previous page"
              className="h-7 px-2.5"
            >
              <ChevronLeft className="h-3.5 w-3.5" aria-hidden="true" />
              <span className="hidden sm:inline">Prev</span>
            </Button>
            <Button
              type="button"
              variant="secondary"
              size="xs"
              disabled={!pagination.canNext || auditQuery.isFetching}
              onClick={() => setOffset((prev) => prev + AUDIT_PAGE_LIMIT)}
              aria-label="Next page"
              className="h-7 px-2.5"
            >
              <span className="hidden sm:inline">Next</span>
              <ChevronRight className="h-3.5 w-3.5" aria-hidden="true" />
            </Button>
          </div>
        </div>
      </Card>

      {/* Record Manual Event Modal */}
      {createModalOpen ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center p-4"
          role="dialog"
          aria-modal="true"
          aria-labelledby="record-audit-title"
        >
          <button
            type="button"
            className="fixed inset-0 bg-black/55"
            aria-label="Close modal"
            onClick={() => setCreateModalOpen(false)}
          />
          <Card className="relative z-10 w-full max-w-lg p-4 sm:p-5">
            <div className="mb-3 flex items-center justify-between gap-2">
              <div>
                <h3
                  id="record-audit-title"
                  className="text-sm font-semibold text-[var(--text-primary)]"
                >
                  Record Manual Audit Event
                </h3>
                <p className="font-mono text-[11px] text-[var(--text-muted)]">
                  Secrets are automatically redacted
                </p>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="h-7 w-7"
                onClick={() => setCreateModalOpen(false)}
                aria-label="Close"
              >
                <X className="h-4 w-4" aria-hidden="true" />
              </Button>
            </div>

            <form onSubmit={onSubmit} className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              <div className="space-y-1">
                <Label htmlFor="audit-category">Category</Label>
                <Input id="audit-category" {...form.register("category")} placeholder="security" />
                {form.formState.errors.category ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {form.formState.errors.category.message}
                  </p>
                ) : null}
              </div>
              <div className="space-y-1">
                <Label htmlFor="audit-action">Action</Label>
                <Input id="audit-action" {...form.register("action")} placeholder="operator.note" />
                {form.formState.errors.action ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {form.formState.errors.action.message}
                  </p>
                ) : null}
              </div>
              <div className="space-y-1">
                <Label htmlFor="audit-severity">Severity</Label>
                <select
                  id="audit-severity"
                  {...form.register("severity")}
                  className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[var(--accent-primary)]"
                >
                  {AUDIT_SEVERITIES.map((sev) => (
                    <option key={sev} value={sev}>
                      {sev}
                    </option>
                  ))}
                </select>
              </div>
              <div className="space-y-1">
                <Label htmlFor="audit-resource">Resource ID (optional)</Label>
                <Input
                  id="audit-resource"
                  {...form.register("resource_id")}
                  placeholder="key_..."
                />
              </div>
              <div className="space-y-1 sm:col-span-2">
                <Label htmlFor="audit-summary">Summary</Label>
                <Input
                  id="audit-summary"
                  {...form.register("summary")}
                  placeholder="Describe the security-relevant action or observation"
                />
                {form.formState.errors.summary ? (
                  <p className="text-xs text-[var(--status-error)]">
                    {form.formState.errors.summary.message}
                  </p>
                ) : null}
              </div>
              <details className="sm:col-span-2">
                <summary className="cursor-pointer text-xs font-medium text-[var(--text-muted)] hover:text-[var(--text-primary)]">
                  Metadata JSON (optional)
                </summary>
                <div className="mt-2 space-y-1">
                  <Input
                    id="audit-metadata"
                    className="font-mono text-xs"
                    {...form.register("metadata_text")}
                    aria-label="Metadata JSON"
                  />
                  {form.formState.errors.metadata_text ? (
                    <p className="text-xs text-[var(--status-error)]">
                      {form.formState.errors.metadata_text.message}
                    </p>
                  ) : null}
                </div>
              </details>
              {createMutation.isError ? (
                <p className="text-xs text-[var(--status-error)] sm:col-span-2">
                  {createMutation.error.message}
                </p>
              ) : null}
              <div className="flex justify-end gap-2 pt-1 sm:col-span-2">
                <Button
                  type="button"
                  variant="secondary"
                  size="sm"
                  onClick={() => setCreateModalOpen(false)}
                >
                  Cancel
                </Button>
                <Button type="submit" size="sm" disabled={createMutation.isPending}>
                  {createMutation.isPending ? "Saving..." : "Save Event"}
                </Button>
              </div>
            </form>
          </Card>
        </div>
      ) : null}
    </div>
  );
}
