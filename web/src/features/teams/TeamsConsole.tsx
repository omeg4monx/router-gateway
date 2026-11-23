"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import {
  AlertTriangle,
  KeyRound,
  Plus,
  RefreshCw,
  RotateCcw,
  Trash2,
  Users,
  X,
} from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useApiKeys } from "@/features/keys/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useAssignTeamKeyMutation,
  useCreateTeamMutation,
  useDeleteTeamMutation,
  useRemoveTeamKeyMutation,
  useResetTeamSpendMutation,
  useTeamsList,
  useUpdateTeamMutation,
} from "./hooks";
import { type CreateTeamFormValues, createTeamSchema } from "./schemas";

export function TeamsConsole() {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [selectedKeyByTeam, setSelectedKeyByTeam] = useState<Record<string, string>>({});

  const teamsQuery = useTeamsList();
  const keysQuery = useApiKeys();

  const createMutation = useCreateTeamMutation();
  const updateMutation = useUpdateTeamMutation();
  const deleteMutation = useDeleteTeamMutation();
  const assignKeyMutation = useAssignTeamKeyMutation();
  const removeKeyMutation = useRemoveTeamKeyMutation();
  const resetSpendMutation = useResetTeamSpendMutation();

  const form = useForm<CreateTeamFormValues>({
    resolver: zodResolver(createTeamSchema),
    defaultValues: {
      name: "",
      description: "",
      budget_usd: 0,
      rate_limit_rpm: 0,
      enabled: true,
    },
  });

  const onSubmit = form.handleSubmit((values) => {
    createMutation.mutate(values, {
      onSuccess: () => {
        form.reset();
        setShowCreateForm(false);
      },
    });
  });

  const mutationError =
    createMutation.error ??
    updateMutation.error ??
    deleteMutation.error ??
    assignKeyMutation.error ??
    removeKeyMutation.error ??
    resetSpendMutation.error;

  const teams = teamsQuery.data?.teams ?? [];
  const allKeys = keysQuery.data ?? [];
  const keyNameMap = new Map(allKeys.map((k) => [k.id, `${k.name} (${k.key_prefix})`]));

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Quota &amp; Budget Grouping
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Teams
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Group gateway API keys under shared USD spend caps and sliding-window RPM limits.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => void Promise.all([teamsQuery.refetch(), keysQuery.refetch()])}
            disabled={teamsQuery.isFetching}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{teamsQuery.isFetching ? "Refreshing" : "Refresh"}</span>
          </Button>
          <Button type="button" size="sm" onClick={() => setShowCreateForm((prev) => !prev)}>
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>New Team</span>
          </Button>
        </div>
      </section>

      {teamsQuery.isError || mutationError ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>Teams request failed</AlertTitle>
            <AlertDescription className="font-mono">
              {teamsQuery.error?.message ?? mutationError?.message}
            </AlertDescription>
          </div>
        </Alert>
      ) : null}

      {showCreateForm ? (
        <Card className="p-5">
          <h3 className="mb-4 text-sm font-semibold text-[var(--text-primary)]">
            Create API-Key Team
          </h3>
          <form onSubmit={onSubmit} className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="team-name">Team Name</Label>
              <Input id="team-name" {...form.register("name")} placeholder="Engineering" />
              {form.formState.errors.name ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.name.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="team-desc">Description</Label>
              <Input
                id="team-desc"
                {...form.register("description")}
                placeholder="Shared budget for backend services"
              />
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="team-budget">Shared USD Budget (0 = unlimited)</Label>
              <Input
                id="team-budget"
                type="number"
                step="0.01"
                min="0"
                {...form.register("budget_usd", { valueAsNumber: true })}
              />
              {form.formState.errors.budget_usd ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.budget_usd.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="team-rpm">Shared RPM Limit (0 = unlimited)</Label>
              <Input
                id="team-rpm"
                type="number"
                min="0"
                {...form.register("rate_limit_rpm", { valueAsNumber: true })}
              />
              {form.formState.errors.rate_limit_rpm ? (
                <p className="text-xs text-[var(--status-error)]">
                  {form.formState.errors.rate_limit_rpm.message}
                </p>
              ) : null}
            </div>
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
                {createMutation.isPending ? "Creating..." : "Create Team"}
              </Button>
            </div>
          </form>
        </Card>
      ) : null}

      <EmberRule />

      <div className="grid grid-cols-1 gap-4 lg:grid-cols-2">
        {teams.length > 0 ? (
          teams.map((team) => {
            const pct =
              team.budget_usd > 0
                ? Math.min(100, Math.round((team.spent_usd / team.budget_usd) * 100))
                : 0;
            const selectedKeyId = selectedKeyByTeam[team.id] ?? "";
            const unassignedKeys = allKeys.filter((k) => !team.key_ids.includes(k.id));

            return (
              <Card key={team.id} className="flex flex-col">
                <CardHeader className="flex flex-row items-start justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex items-center gap-2">
                      <Users className="h-4 w-4 shrink-0 text-[#E56A4A]" aria-hidden="true" />
                      <CardTitle className="truncate">{team.name}</CardTitle>
                      <Badge variant={team.enabled ? "success" : "destructive"}>
                        {team.enabled ? "Active" : "Disabled"}
                      </Badge>
                    </div>
                    {team.description ? (
                      <p className="mt-1 text-xs text-[var(--text-secondary)]">
                        {team.description}
                      </p>
                    ) : null}
                  </div>
                  <div className="flex shrink-0 items-center gap-1.5">
                    <Button
                      type="button"
                      variant="outline"
                      size="xs"
                      onClick={() =>
                        updateMutation.mutate({ id: team.id, body: { enabled: !team.enabled } })
                      }
                    >
                      {team.enabled ? "Disable" : "Enable"}
                    </Button>
                    <Button
                      type="button"
                      variant="destructive"
                      size="xs"
                      onClick={() => deleteMutation.mutate(team.id)}
                      aria-label={`Delete team ${team.name}`}
                    >
                      <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                    </Button>
                  </div>
                </CardHeader>

                <div className="space-y-4 px-5 pb-5">
                  <div className="grid grid-cols-2 gap-3 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] p-3">
                    <div>
                      <p className="text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                        Spend / Budget
                      </p>
                      <p className="mt-1 font-mono text-sm font-semibold tabular-nums text-[var(--text-primary)]">
                        ${team.spent_usd.toFixed(4)} /{" "}
                        {team.budget_usd > 0 ? `$${team.budget_usd.toFixed(2)}` : "Unlimited"}
                      </p>
                      {team.budget_usd > 0 ? (
                        <div className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-[var(--bg-surface-raised)]">
                          <div className="h-full bg-[#E56A4A]" style={{ width: `${pct}%` }} />
                        </div>
                      ) : null}
                    </div>
                    <div className="flex items-start justify-between">
                      <div>
                        <p className="text-[11px] uppercase tracking-wider text-[var(--text-muted)]">
                          Shared RPM
                        </p>
                        <p className="mt-1 font-mono text-sm font-semibold tabular-nums text-[var(--text-primary)]">
                          {team.rate_limit_rpm > 0 ? `${team.rate_limit_rpm} RPM` : "Unlimited"}
                        </p>
                      </div>
                      <Button
                        type="button"
                        variant="secondary"
                        size="xs"
                        onClick={() => resetSpendMutation.mutate(team.id)}
                        title="Reset accumulated spend to $0"
                      >
                        <RotateCcw className="h-3 w-3" aria-hidden="true" />
                        <span>Reset</span>
                      </Button>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <p className="text-xs font-medium uppercase tracking-wider text-[var(--text-muted)]">
                      Member API Keys ({team.key_ids.length})
                    </p>
                    {team.key_ids.length > 0 ? (
                      <div className="flex flex-wrap gap-1.5">
                        {team.key_ids.map((kid) => (
                          <span
                            key={kid}
                            className="inline-flex items-center gap-1.5 rounded border border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] px-2 py-1 font-mono text-xs text-[var(--text-primary)]"
                          >
                            <KeyRound className="h-3 w-3 text-[#E56A4A]" aria-hidden="true" />
                            <span>{keyNameMap.get(kid) ?? kid}</span>
                            <button
                              type="button"
                              onClick={() =>
                                removeKeyMutation.mutate({ teamId: team.id, keyId: kid })
                              }
                              aria-label={`Remove key ${kid} from ${team.name}`}
                              className="text-[var(--text-muted)] hover:text-[var(--status-error)]"
                            >
                              <X className="h-3 w-3" aria-hidden="true" />
                            </button>
                          </span>
                        ))}
                      </div>
                    ) : (
                      <p className="text-xs text-[var(--text-muted)]">No keys assigned yet.</p>
                    )}

                    <div className="flex items-center gap-2 pt-1">
                      <select
                        value={selectedKeyId}
                        onChange={(e) =>
                          setSelectedKeyByTeam((prev) => ({ ...prev, [team.id]: e.target.value }))
                        }
                        aria-label={`Select key to assign to ${team.name}`}
                        className="h-8 flex-1 rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-2 text-xs text-[var(--text-primary)]"
                      >
                        <option value="">Select API key to assign...</option>
                        {unassignedKeys.map((k) => (
                          <option key={k.id} value={k.id}>
                            {k.name} ({k.key_prefix})
                          </option>
                        ))}
                      </select>
                      <Button
                        type="button"
                        variant="outline"
                        size="xs"
                        disabled={!selectedKeyId || assignKeyMutation.isPending}
                        onClick={() => {
                          if (!selectedKeyId) return;
                          assignKeyMutation.mutate(
                            { teamId: team.id, keyId: selectedKeyId },
                            {
                              onSuccess: () =>
                                setSelectedKeyByTeam((prev) => ({ ...prev, [team.id]: "" })),
                            },
                          );
                        }}
                      >
                        Assign
                      </Button>
                    </div>
                  </div>
                </div>
              </Card>
            );
          })
        ) : (
          <Card className="p-8 text-center text-sm text-[var(--text-muted)] lg:col-span-2">
            {teamsQuery.isPending
              ? "Loading teams..."
              : "No teams configured yet. Create a team to group API keys under a shared budget or RPM cap."}
          </Card>
        )}
      </div>
    </div>
  );
}
