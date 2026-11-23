"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { AlertTriangle, Check, Copy, Plus, RefreshCw, Ticket, Trash2 } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useTeamsList } from "@/features/teams/hooks";
import { Alert, AlertDescription, AlertTitle } from "@/shared/ui/alert";
import { Badge } from "@/shared/ui/badge";
import { Button } from "@/shared/ui/button";
import { Card, CardHeader, CardTitle } from "@/shared/ui/card";
import { EmberRule } from "@/shared/ui/EmberRule";
import { Input } from "@/shared/ui/input";
import { Label } from "@/shared/ui/label";
import {
  useClaimVoucherMutation,
  useCreateVoucherMutation,
  useDeleteVoucherMutation,
  useUpdateVoucherMutation,
  useVouchersList,
} from "./hooks";
import {
  type ClaimVoucherFormValues,
  type ClaimVoucherResponse,
  type CreateVoucherFormValues,
  claimVoucherSchema,
  createVoucherSchema,
} from "./schemas";

export function VouchersConsole() {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [claimedResult, setClaimedResult] = useState<ClaimVoucherResponse | null>(null);
  const [copiedSecret, setCopiedSecret] = useState(false);
  const [clipboardError, setClipboardError] = useState<string | null>(null);

  const vouchersQuery = useVouchersList();
  const teamsQuery = useTeamsList();
  const createMutation = useCreateVoucherMutation();
  const updateMutation = useUpdateVoucherMutation();
  const deleteMutation = useDeleteVoucherMutation();
  const claimMutation = useClaimVoucherMutation();

  const createForm = useForm<CreateVoucherFormValues>({
    resolver: zodResolver(createVoucherSchema),
    defaultValues: {
      code: "",
      name: "",
      description: "",
      quota_usd: 5,
      rate_limit_rpm: 30,
      allowed_models_text: "",
      team_id: "",
      max_redemptions: 1,
    },
  });

  const claimForm = useForm<ClaimVoucherFormValues>({
    resolver: zodResolver(claimVoucherSchema),
    defaultValues: {
      code: "",
      claimer_id: "",
      key_name: "",
    },
  });

  const onCreateSubmit = createForm.handleSubmit((values) => {
    const allowedModels = values.allowed_models_text
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean);
    createMutation.mutate(
      {
        code: values.code || undefined,
        name: values.name,
        description: values.description || undefined,
        quota_usd: values.quota_usd,
        rate_limit_rpm: values.rate_limit_rpm,
        allowed_models: allowedModels,
        team_id: values.team_id || undefined,
        max_redemptions: values.max_redemptions,
      },
      {
        onSuccess: () => {
          createForm.reset();
          setShowCreateForm(false);
        },
      },
    );
  });

  const onClaimSubmit = claimForm.handleSubmit((values) => {
    setClaimedResult(null);
    claimMutation.mutate(
      {
        code: values.code,
        claimer_id: values.claimer_id || undefined,
        key_name: values.key_name || undefined,
      },
      {
        onSuccess: (res) => {
          setClaimedResult(res);
          claimForm.reset();
        },
      },
    );
  });

  const handleCopySecret = async () => {
    if (!claimedResult?.secret) return;
    setClipboardError(null);
    try {
      await navigator.clipboard.writeText(claimedResult.secret);
      setCopiedSecret(true);
      setTimeout(() => setCopiedSecret(false), 1500);
    } catch {
      setClipboardError("Clipboard access denied; copy the secret manually.");
    }
  };

  const vouchers = vouchersQuery.data?.vouchers ?? [];
  const teams = teamsQuery.data?.teams ?? [];
  const teamNames = new Map(teams.map((t) => [t.id, t.name]));
  const errorMsg =
    vouchersQuery.error?.message ??
    createMutation.error?.message ??
    updateMutation.error?.message ??
    deleteMutation.error?.message ??
    claimMutation.error?.message ??
    clipboardError;

  return (
    <div className="space-y-6">
      <section className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="font-mono text-xs uppercase tracking-widest text-[var(--accent-text)]">
            Provisioning &amp; Distribution
          </p>
          <h2 className="mt-1 text-2xl font-semibold tracking-tight text-[var(--text-primary)]">
            Vouchers
          </h2>
          <p className="mt-1 text-sm text-[var(--text-secondary)]">
            Issue redeemable codes that atomically mint scoped gateway API keys with anti-abuse
            claim limits.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => void vouchersQuery.refetch()}
            disabled={vouchersQuery.isFetching}
          >
            <RefreshCw className="h-3.5 w-3.5 text-[#E56A4A]" aria-hidden="true" />
            <span>{vouchersQuery.isFetching ? "Refreshing" : "Refresh"}</span>
          </Button>
          <Button type="button" size="sm" onClick={() => setShowCreateForm((prev) => !prev)}>
            <Plus className="h-3.5 w-3.5" aria-hidden="true" />
            <span>Create Voucher</span>
          </Button>
        </div>
      </section>

      {errorMsg ? (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4 text-[var(--status-error)]" aria-hidden="true" />
          <div>
            <AlertTitle>Voucher operation failed</AlertTitle>
            <AlertDescription className="font-mono">{errorMsg}</AlertDescription>
          </div>
        </Alert>
      ) : null}

      {showCreateForm ? (
        <Card className="p-5">
          <h3 className="mb-4 text-sm font-semibold text-[var(--text-primary)]">
            New Redeemable Voucher
          </h3>
          <form
            onSubmit={onCreateSubmit}
            noValidate
            className="grid grid-cols-1 gap-4 md:grid-cols-2"
          >
            <div className="space-y-1.5">
              <Label htmlFor="voucher-name">Name</Label>
              <Input id="voucher-name" {...createForm.register("name")} placeholder="Beta Cohort" />
              {createForm.formState.errors.name ? (
                <p className="text-xs text-[var(--status-error)]">
                  {createForm.formState.errors.name.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="voucher-code">Custom Code (optional, auto-generated if blank)</Label>
              <Input
                id="voucher-code"
                className="font-mono uppercase"
                {...createForm.register("code")}
                placeholder="OMG-BETA-2026"
              />
              {createForm.formState.errors.code ? (
                <p className="text-xs text-[var(--status-error)]">
                  {createForm.formState.errors.code.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="voucher-quota">Key Quota USD</Label>
              <Input
                id="voucher-quota"
                type="number"
                step="0.01"
                min="0"
                {...createForm.register("quota_usd", {
                  setValueAs: (v: string) => (v === "" ? 0 : Number(v)),
                })}
              />
              {createForm.formState.errors.quota_usd ? (
                <p className="text-xs text-[var(--status-error)]">
                  {createForm.formState.errors.quota_usd.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="voucher-rpm">Key RPM Limit</Label>
              <Input
                id="voucher-rpm"
                type="number"
                min="0"
                {...createForm.register("rate_limit_rpm", {
                  setValueAs: (v: string) => (v === "" ? 0 : Number(v)),
                })}
              />
              {createForm.formState.errors.rate_limit_rpm ? (
                <p className="text-xs text-[var(--status-error)]">
                  {createForm.formState.errors.rate_limit_rpm.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="voucher-max">Max Redemptions</Label>
              <Input
                id="voucher-max"
                type="number"
                min="1"
                {...createForm.register("max_redemptions", {
                  setValueAs: (v: string) => (v === "" ? 1 : Number(v)),
                })}
              />
              {createForm.formState.errors.max_redemptions ? (
                <p className="text-xs text-[var(--status-error)]">
                  {createForm.formState.errors.max_redemptions.message}
                </p>
              ) : null}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="voucher-team">Attach to Team (optional)</Label>
              <select
                id="voucher-team"
                {...createForm.register("team_id")}
                className="flex h-9 w-full rounded-md border border-[var(--border-subtle)] bg-[var(--bg-canvas)] px-3 py-1 text-sm text-[var(--text-primary)]"
              >
                <option value="">No team attachment</option>
                {teams.map((t) => (
                  <option key={t.id} value={t.id}>
                    {t.name}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-1.5 md:col-span-2">
              <Label htmlFor="voucher-models">Allowed Models (comma-separated, blank = all)</Label>
              <Input
                id="voucher-models"
                {...createForm.register("allowed_models_text")}
                placeholder="gpt-4o, claude-3-7-sonnet"
              />
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
                {createMutation.isPending ? "Creating..." : "Create Voucher"}
              </Button>
            </div>
          </form>
        </Card>
      ) : null}

      <Card className="p-5">
        <div className="mb-3 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Ticket className="h-4 w-4 text-[#E56A4A]" aria-hidden="true" />
            <h3 className="text-sm font-semibold text-[var(--text-primary)]">
              Redeem Voucher (Test / Operator Claim)
            </h3>
          </div>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            POST /api/vouchers/claim
          </span>
        </div>
        <form onSubmit={onClaimSubmit} className="grid grid-cols-1 gap-3 sm:grid-cols-4">
          <div>
            <Input
              {...claimForm.register("code")}
              placeholder="Voucher code..."
              className="font-mono uppercase"
              aria-label="Voucher code to claim"
            />
            {claimForm.formState.errors.code ? (
              <p className="mt-1 text-xs text-[var(--status-error)]">
                {claimForm.formState.errors.code.message}
              </p>
            ) : null}
          </div>
          <Input
            {...claimForm.register("claimer_id")}
            placeholder="Claimer ID (optional)"
            aria-label="Claimer ID"
          />
          <Input
            {...claimForm.register("key_name")}
            placeholder="Custom key name (optional)"
            aria-label="Custom key name"
          />
          <Button type="submit" disabled={claimMutation.isPending}>
            {claimMutation.isPending ? "Claiming..." : "Claim & Mint Key"}
          </Button>
        </form>

        {claimedResult ? (
          <div className="mt-4 rounded-md border border-[#E56A4A] bg-[var(--bg-canvas)] p-3">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <p className="text-xs font-semibold text-[var(--text-primary)]">
                Provisioned API Key (shown once)
              </p>
              <Button type="button" variant="outline" size="xs" onClick={handleCopySecret}>
                {copiedSecret ? (
                  <Check className="h-3.5 w-3.5 text-[var(--status-ok)]" aria-hidden="true" />
                ) : (
                  <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                )}
                <span>{copiedSecret ? "Copied" : "Copy Secret"}</span>
              </Button>
            </div>
            <p className="mt-1.5 break-all font-mono text-xs text-[var(--accent-text)]">
              {claimedResult.secret}
            </p>
          </div>
        ) : null}
      </Card>

      <EmberRule />

      <Card>
        <CardHeader>
          <CardTitle>Configured Vouchers</CardTitle>
          <span className="font-mono text-xs text-[var(--text-muted)]">
            {vouchers.length} voucher(s)
          </span>
        </CardHeader>
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-y border-[var(--border-subtle)] bg-[var(--bg-surface-raised)] text-xs uppercase tracking-wider text-[var(--text-muted)]">
              <tr>
                <th className="px-4 py-2.5">Code / Name</th>
                <th className="px-4 py-2.5">Status</th>
                <th className="px-4 py-2.5">Redemptions</th>
                <th className="px-4 py-2.5">Quota / RPM</th>
                <th className="px-4 py-2.5">Team</th>
                <th className="px-4 py-2.5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-[var(--border-subtle)]">
              {vouchers.length > 0 ? (
                vouchers.map((v) => (
                  <tr key={v.id} className="hover:bg-[var(--bg-surface-raised)]">
                    <td className="px-4 py-3">
                      <p className="font-mono text-sm font-semibold text-[var(--text-primary)]">
                        {v.code}
                      </p>
                      <p className="text-xs text-[var(--text-secondary)]">{v.name}</p>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3">
                      <Badge
                        variant={
                          !v.enabled
                            ? "destructive"
                            : v.redeemed_count >= v.max_redemptions
                              ? "secondary"
                              : "success"
                        }
                      >
                        {!v.enabled
                          ? "Disabled"
                          : v.redeemed_count >= v.max_redemptions
                            ? "Exhausted"
                            : "Active"}
                      </Badge>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-primary)]">
                      {v.redeemed_count} / {v.max_redemptions}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs tabular-nums text-[var(--text-secondary)]">
                      ${v.quota_usd.toFixed(2)} &middot;{" "}
                      {v.rate_limit_rpm > 0 ? `${v.rate_limit_rpm} RPM` : "Unlimited"}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 font-mono text-xs text-[var(--text-secondary)]">
                      {v.team_id ? (teamNames.get(v.team_id) ?? v.team_id) : "None"}
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-right">
                      <div className="inline-flex items-center gap-1.5">
                        <Button
                          type="button"
                          variant="outline"
                          size="xs"
                          onClick={() =>
                            updateMutation.mutate({ id: v.id, body: { enabled: !v.enabled } })
                          }
                        >
                          {v.enabled ? "Disable" : "Enable"}
                        </Button>
                        <Button
                          type="button"
                          variant="destructive"
                          size="xs"
                          onClick={() => deleteMutation.mutate(v.id)}
                          aria-label={`Delete voucher ${v.code}`}
                        >
                          <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td
                    colSpan={6}
                    className="px-4 py-8 text-center text-sm text-[var(--text-muted)]"
                  >
                    {vouchersQuery.isPending
                      ? "Loading vouchers..."
                      : "No vouchers configured yet."}
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
