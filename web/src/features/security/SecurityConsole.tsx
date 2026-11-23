"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  Download,
  KeyRound,
  Lock,
  Plus,
  RefreshCw,
  Search,
  ShieldCheck,
  Trash2,
} from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClearSecurityAuditMutation,
  useCreateSecurityAuditMutation,
  useSecurityAuditEvents,
  useSecurityPosture,
} from "./hooks";
import {
  AUDIT_SEVERITIES,
  type CreateSecurityAuditFormValues,
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

export function SecurityConsole() {
  const [categoryFilter, setCategoryFilter] = useState("");
  const [severityFilter, setSeverityFilter] = useState<"" | "info" | "warning" | "critical">("");
  const [searchQuery, setSearchQuery] = useState("");
  const [showCreateForm, setShowCreateForm] = useState(false);

  const postureQuery = useSecurityPosture();
  const auditQuery = useSecurityAuditEvents({
    category: categoryFilter || undefined,
    severity: severityFilter || undefined,
    q: searchQuery.trim() || undefined,
    limit: 100,
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
          setShowCreateForm(false);
        },
      },
    );
  });

  const handleRefresh = () => {
    void Promise.all([postureQuery.refetch(), auditQuery.refetch()]);
  };

  const posture = postureQuery.data;
  const auditData = auditQuery.data;
  const isRefreshing = postureQuery.isFetching || auditQuery.isFetching;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Security &amp; Compliance
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Security Posture &amp; Audit Trail
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Live listener posture, vault key provenance, and redacted control-plane audit log.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleRefresh}
            disabled={isRefreshing}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{isRefreshing ? "Refreshing" : "Refresh"}</span>
          </Button>
          <Button asChild variant="outline" size="sm">
            <a href="/api/security/audit/export" download="security-audit.json">
              <Download className="h-3.5 w-3.5" aria-hidden="true" />
              <span>Export JSON</span>
            </a>
          </Button>
          <Button type="button" size="sm" onClick={() => setShowCreateForm((prev) => !prev)}>
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Record Event</span>
          </Button>
        </div>
      </section>

      {postureQuery.isError ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>Failed to load security posture</AlertTitle>
            <AlertDescription className="font-mono">{postureQuery.error.message}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      <section
        aria-label="Security posture metrics"
        className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4"
      >
        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Bind Scope
            </span>
            <ShieldCheck className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
          </div>
          <p className="mt-3 font-mono text-lg font-semibold text-[var(--text-primary)]">
            {posture ? posture.bind_addr : "--"}
          </p>
          <p className="mt-1 text-xs text-[var(--text-muted)]">
            {posture
              ? posture.is_loopback
                ? "Loopback only"
                : "Remote interface exposed"
              : "Checking..."}
          </p>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Admin Login
            </span>
            <Lock className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
          </div>
          <p
            className={`mt-3 font-mono text-lg font-semibold ${
              posture?.require_login ? "text-[var(--status-ok)]" : "text-[var(--status-warn)]"
            }`}
          >
            {posture ? (posture.require_login ? "ENFORCED" : "DISABLED") : "--"}
          </p>
          <p className="mt-1 text-xs text-[var(--text-muted)]">
            {posture?.allow_insecure_remote
              ? "Insecure remote override active"
              : "Remote safety guard active"}
          </p>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Vault Master Key
            </span>
            <KeyRound className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
          </div>
          <p className="mt-3 font-mono text-lg font-semibold uppercase text-[var(--text-primary)]">
            {posture ? posture.master_key_source : "--"}
          </p>
          <p className="mt-1 text-xs text-[var(--text-muted)]">AES-256-GCM credential encryption</p>
        </Card>

        <Card className="p-4">
          <div className="flex items-center justify-between">
            <span className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
              Audit Records
            </span>
            <AlertTriangle className="h-4 w-4 text-[var(--text-muted)]" aria-hidden="true" />
          </div>
          <p className="mt-3 font-mono text-lg font-semibold tabular-nums text-[var(--text-primary)]">
            {posture ? posture.audit_event_count : 0}
          </p>
          <p className="mt-1 text-xs text-[var(--text-muted)]">
            {posture
              ? `${posture.warnings.length} posture warning(s)`
              : "Retained up to 90d / 5,000"}
          </p>
        </Card>
      </section>

      {posture && posture.warnings.length > 0 ? (
        <Alert variant="warning">
          <AlertTriangle
            className="h-4 w-4 shrink-0 text-[var(--status-warn)]"
            aria-hidden="true"
          />
          <div className="space-y-1">
            <AlertTitle>Security Posture Recommendations</AlertTitle>
            <AlertDescription>
              <ul className="list-disc space-y-1 pl-4 text-sm">
                {posture.warnings.map((warning) => (
                  <li key={warning}>{warning}</li>
                ))}
              </ul>
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      {showCreateForm ? (
        <Card className="p-5">
          <div className="mb-4 flex items-center justify-between">
            <h3 className="text-sm font-semibold text-[var(--text-primary)]">
              Record Manual Audit Event
            </h3>
            <span className="font-mono text-xs text-[var(--text-muted)]">
              Secrets automatically redacted
            </span>
          </div>
          <form onSubmit={onSubmit} className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="audit-category">Category</Label>
              <Input id="audit-category" {...form.register("category")} placeholder="security" />
              {form.formState.errors.category ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.category.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="audit-action">Action</Label>
              <Input id="audit-action" {...form.register("action")} placeholder="operator.note" />
              {form.formState.errors.action ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.action.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="audit-severity">Severity</Label>
              <select
                id="audit-severity"
                {...form.register("severity")}
                className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
              >
                {AUDIT_SEVERITIES.map((sev) => (
                  <option key={sev} value={sev}>
                    {sev}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="audit-resource">Resource ID (optional)</Label>
              <Input id="audit-resource" {...form.register("resource_id")} placeholder="key_..." />
            </div>
            <div className="space-y-1.5 md:col-span-2">
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
            <div className="space-y-1.5 md:col-span-2">
              <Label htmlFor="audit-metadata">Metadata JSON</Label>
              <Input
                id="audit-metadata"
                className="font-mono text-xs"
                {...form.register("metadata_text")}
              />
              {form.formState.errors.metadata_text ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.metadata_text.message}
                </p>
              ) : null}
            </div>
            {createMutation.isError ? (
              <p className="text-xs text-[var(--status-error)] md:col-span-2">
                {createMutation.error.message}
              </p>
            ) : null}
            <div className="flex justify-end gap-2 md:col-span-2">
              <Button
                type="button"
                variant="secondary"
                size="sm"
                onClick={() => setShowCreateForm(false)}
              >
                Cancel
              </Button>
              <Button type="submit" size="sm" disabled={createMutation.isPending}>
                {createMutation.isPending ? "Recording..." : "Save Event"}
              </Button>
            </div>
          </form>
        </Card>
      ) : null}

      <EmberRule />

      <Card>
        <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <CardTitle>Audit Trail</CardTitle>
            <p className="mt-0.5 font-mono text-xs text-[var(--text-muted)]">
              {auditData ? `${auditData.total} matching event(s)` : "Loading events..."}
            </p>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <div className="relative">
              <Search
                className="pointer-events-none absolute left-2.5 top-2.5 h-3.5 w-3.5 text-[var(--text-muted)]"
                aria-hidden="true"
              />
              <Input
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search summary, action..."
                className="h-8 w-48 pl-8 text-xs"
                aria-label="Search audit trail"
              />
            </div>
            <Input
              value={categoryFilter}
              onChange={(e) => setCategoryFilter(e.target.value)}
              placeholder="Category..."
              className="h-8 w-32 text-xs"
              aria-label="Filter by category"
            />
            <select
              value={severityFilter}
              onChange={(e) =>
                setSeverityFilter(e.target.value as "" | "info" | "warning" | "critical")
              }
              aria-label="Filter by severity"
              className="h-8 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 text-xs text-[var(--text-primary)]"
            >
              <option value="">All severities</option>
              {AUDIT_SEVERITIES.map((sev) => (
                <option key={sev} value={sev}>
                  {sev}
                </option>
              ))}
            </select>
            <Button
              type="button"
              variant="destructive"
              size="xs"
              onClick={() => clearMutation.mutate()}
              disabled={clearMutation.isPending}
              title="Clear audit trail (records an audit.cleared marker)"
            >
              <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
              <span>{clearMutation.isPending ? "Clearing..." : "Clear Trail"}</span>
            </Button>
          </div>
        </CardHeader>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-y border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-xs uppercase tracking-wider text-[var(--text-muted)]">
              <tr>
                <th className="px-4 py-2.5">Timestamp</th>
                <th className="px-4 py-2.5">Severity</th>
                <th className="px-4 py-2.5">Category / Action</th>
                <th className="px-4 py-2.5">Actor</th>
                <th className="px-4 py-2.5">Summary</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--border-subtle)]">
              {auditData && auditData.items.length > 0 ? (
                auditData.items.map((ev) => (
                  <tr key={ev.id} className="hover:bg-[var(--bg-surface-raised)]">
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-muted)]">
                      {new Date(ev.created_at).toLocaleString()}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3">
                      <Badge variant={severityBadgeVariant(ev.severity)}>{ev.severity}</Badge>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs">
                      <span className="text-[var(--accent-text)]">{ev.category}</span>
                      <span className="text-[var(--text-muted)]"> / </span>
                      <span className="text-[var(--text-primary)]">{ev.action}</span>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs text-[var(--text-secondary)]">
                      {ev.actor}
                    </td>
                    <td className="px-4 py-3 text-sm text-[var(--text-primary)]">{ev.summary}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan={5}
                    className="px-4 py-8 text-center text-sm text-[var(--text-muted)]"
                  >
                    {auditQuery.isPending
                      ? "Loading security audit trail..."
                      : "No security audit events match the current filter."}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
}
