"use client";

import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Modal, Button, Input, Select, Toggle, SegmentedControl } from "@/shared/components";
import { cn } from "@/shared/utils/cn";
import { RESET_INTERVAL_OPTIONS, validateCustomInterval } from "../endpointUtils";

const PRESET_RESET_VALUES = new Set(["never", "1h", "5h", "1d", "7d", "14d", "30d"]);

export const PERMISSION_OPTIONS = [
  { key: "manageApiKeys", label: "Create, edit and delete API keys", icon: "key", desc: "Lets this key manage other API keys" },
  { key: "manageModels", label: "Create, edit and delete models", icon: "auto_awesome", desc: "Custom models, aliases and combos" },
  { key: "manageProviders", label: "Create, edit and delete providers", icon: "dns", desc: "Provider connections and API keys" },
  { key: "manageTools", label: "Use CLI tools and token saver", icon: "terminal", desc: "Request transformers and the CLI install page" },
  { key: "manageAdvanced", label: "Open console log, translator and proxy pools", icon: "terminal", desc: "Debugging and network plumbing" },
  { key: "managePlugins", label: "Manage custom plugins", icon: "widgets", desc: "Plugins that change how requests are handled" },
  { key: "manageMediaProviders", label: "Manage media providers", icon: "perm_media", desc: "Image, audio and embedding connections" },
  { key: "viewUsage", label: "View usage", icon: "bar_chart", desc: "Usage numbers for this key only" },
];

export const EMPTY_PERMISSIONS = {
  manageApiKeys: false,
  manageModels: false,
  manageProviders: false,
  manageTools: false,
  manageAdvanced: false,
  managePlugins: false,
  manageMediaProviders: false,
  viewUsage: true,
};

const PERMISSIONS_LOCKED_REASON = "Locked. Sign in with the dashboard password to change permissions.";

function toDateTimeLocalValue(raw) {
  if (!raw || typeof raw !== "string") return "";
  if (/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}/.test(raw)) return raw.slice(0, 16);
  const d = new Date(raw);
  if (Number.isNaN(d.getTime())) return "";
  const pad = (n) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function parseModelsList(str) {
  if (!str || str.trim() === "*" || str.trim() === "") return [];
  return str
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
}

function PermissionsEditor({ value, onChange, allowed, locked, lockedReason }) {
  const effective = locked ? EMPTY_PERMISSIONS : value;
  return (
    <div className="flex flex-col gap-2">
      <div>
        <label className="text-sm font-medium text-text-main">Dashboard Permissions</label>
        <p className="text-xs text-text-muted mt-0.5">
          Controls what this key can access if used to sign in to the dashboard.
        </p>
      </div>
      {locked && (
        <div className="rounded-xl border border-border-subtle bg-surface-2 px-3 py-2 text-xs text-text-muted">
          {lockedReason}
        </div>
      )}
      <div className="flex flex-col gap-1.5">
        {PERMISSION_OPTIONS.map((opt) => {
          const unavailable = locked || !allowed?.[opt.key];
          const checked = Boolean(effective?.[opt.key]);
          const note = locked
            ? lockedReason
            : unavailable
              ? "Your key does not hold this permission"
              : opt.desc;
          return (
            <label
              key={opt.key}
              className={cn(
                "flex items-start gap-2.5 rounded-xl border border-border-subtle bg-surface-2/60 px-3 py-2.5 transition-colors",
                unavailable ? "opacity-50 cursor-not-allowed" : "cursor-pointer hover:border-primary/40"
              )}
            >
              <input
                type="checkbox"
                className="mt-0.5 size-4 accent-[var(--color-primary)] shrink-0"
                checked={checked}
                disabled={unavailable}
                onChange={(e) => onChange({ ...EMPTY_PERMISSIONS, ...value, [opt.key]: e.target.checked })}
              />
              <span className="flex min-w-0 flex-col gap-0.5">
                <span className="flex items-center gap-1.5 text-xs sm:text-sm font-medium text-text-main">
                  <span className="material-symbols-outlined text-[16px] text-primary">{opt.icon}</span>
                  {opt.label}
                </span>
                <span className="text-xs text-text-muted">{note}</span>
              </span>
            </label>
          );
        })}
      </div>
    </div>
  );
}

PermissionsEditor.propTypes = {
  value: PropTypes.object,
  onChange: PropTypes.func.isRequired,
  allowed: PropTypes.object,
  locked: PropTypes.bool,
  lockedReason: PropTypes.string,
};

export default function KeyFormModal({
  isOpen,
  mode = "create",
  initialKey = null,
  existingKeys = [],
  creatorPermissions,
  creatorTokenLimit = 0,
  creatorAllowedModels = "*",
  isApiKeyUser = false,
  onOpenModelPicker,
  selectedModelsValue = "*",
  onSelectedModelsChange,
  onSubmit,
  onClose,
}) {
  const [activeTab, setActiveTab] = useState("general");
  const [name, setName] = useState("");
  const [customKey, setCustomKey] = useState("");
  const [appendSecret, setAppendSecret] = useState(true);
  const [tokenLimit, setTokenLimit] = useState("");
  const [resetInterval, setResetInterval] = useState("never");
  const [customReset, setCustomReset] = useState("");
  const [rpmLimit, setRpmLimit] = useState("");
  const [tpmLimit, setTpmLimit] = useState("");
  const [ipWhitelist, setIpWhitelist] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [permissions, setPermissions] = useState({ ...EMPTY_PERMISSIONS });
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const permissionsLocked = isApiKeyUser;

  useEffect(() => {
    if (!isOpen) return;
    setActiveTab("general");
    setError("");
    setSubmitting(false);
    setCustomKey("");
    setAppendSecret(true);

    if (initialKey) {
      setName(initialKey.name || "");
      setTokenLimit(initialKey.tokenLimit ? String(initialKey.tokenLimit) : "");
      const resVal = initialKey.resetInterval || "never";
      if (PRESET_RESET_VALUES.has(resVal)) {
        setResetInterval(resVal);
        setCustomReset("");
      } else {
        setResetInterval("custom");
        setCustomReset(resVal);
      }
      setRpmLimit(initialKey.rpmLimit ? String(initialKey.rpmLimit) : "");
      setTpmLimit(initialKey.tpmLimit ? String(initialKey.tpmLimit) : "");
      setIpWhitelist(initialKey.ipWhitelist || "");
      setExpiresAt(toDateTimeLocalValue(initialKey.expiresAt));
      setPermissions(initialKey.permissions ? { ...EMPTY_PERMISSIONS, ...initialKey.permissions } : { ...EMPTY_PERMISSIONS });
    } else {
      setName("");
      setTokenLimit("");
      setResetInterval("never");
      setCustomReset("");
      setRpmLimit("");
      setTpmLimit("");
      setIpWhitelist("");
      setExpiresAt("");
      setPermissions({ ...EMPTY_PERMISSIONS });
    }
  }, [isOpen, mode, initialKey]);

  const removeModel = (modelVal) => {
    const next = parseModelsList(selectedModelsValue).filter((m) => m !== modelVal);
    onSelectedModelsChange(next.length === 0 ? "*" : next.join(", "));
  };

  const handleSubmit = async () => {
    const trimmedName = name.trim();
    if (!trimmedName) {
      setError("Key name is required.");
      setActiveTab("general");
      return;
    }

    const duplicate = existingKeys.some(
      (k) => k.name === trimmedName && (mode === "create" || k.id !== initialKey?.id)
    );
    if (duplicate) {
      setError(`A key named "${trimmedName}" already exists. Choose a different name.`);
      setActiveTab("general");
      return;
    }

    const limitNum = tokenLimit ? Number(tokenLimit) : 0;
    if (!Number.isFinite(limitNum) || limitNum < 0) {
      setError("Token limit must be 0 or a positive number.");
      setActiveTab("limits");
      return;
    }
    if (creatorTokenLimit > 0 && limitNum > creatorTokenLimit) {
      setError(`Token limit cannot exceed your maximum (${creatorTokenLimit.toLocaleString()}).`);
      setActiveTab("limits");
      return;
    }

    let finalReset = "never";
    if (limitNum > 0) {
      if (resetInterval === "custom") {
        const intervalErr = validateCustomInterval(customReset);
        if (intervalErr) {
          setError(intervalErr);
          setActiveTab("limits");
          return;
        }
        finalReset = customReset.trim();
      } else {
        finalReset = resetInterval;
      }
    }

    const rpmNum = rpmLimit ? Number(rpmLimit) : 0;
    const tpmNum = tpmLimit ? Number(tpmLimit) : 0;
    if (!Number.isFinite(rpmNum) || rpmNum < 0 || !Number.isFinite(tpmNum) || tpmNum < 0) {
      setError("RPM and TPM limits must be 0 or positive numbers.");
      setActiveTab("limits");
      return;
    }

    setError("");
    setSubmitting(true);
    try {
      const payload = {
        name: trimmedName,
        tokenLimit: limitNum,
        resetInterval: finalReset,
        allowedModels: (selectedModelsValue || "*").trim() || "*",
        rpmLimit: rpmNum,
        tpmLimit: tpmNum,
        ipWhitelist: ipWhitelist.trim(),
        expiresAt: expiresAt || null,
        permissions: permissionsLocked ? EMPTY_PERMISSIONS : permissions,
      };
      if (mode === "create") {
        payload.customKey = customKey.trim() || undefined;
        payload.appendSecret = appendSecret;
      }
      const result = await onSubmit(payload);
      if (result && !result.ok) {
        setError(result.error || (mode === "create" ? "Failed to create API key." : "Failed to update API key."));
      }
    } catch (err) {
      setError(err?.message || "Unexpected error saving API key.");
    } finally {
      setSubmitting(false);
    }
  };

  const hasLimitsConfigured =
    Number(tokenLimit) > 0 ||
    Number(rpmLimit) > 0 ||
    Number(tpmLimit) > 0 ||
    parseModelsList(selectedModelsValue).length > 0 ||
    Boolean(ipWhitelist.trim());

  const customKeySanitizedLen = customKey.trim().replace(/[^a-zA-Z0-9-]/g, "-").slice(0, 35).length;

  return (
    <Modal
      isOpen={isOpen}
      title={mode === "edit" ? `Edit API Key — ${initialKey?.name || ""}` : "Create API Key"}
      size="lg"
      onClose={() => {
        if (!submitting) onClose();
      }}
      footer={
        <>
          <Button variant="ghost" onClick={onClose} disabled={submitting}>
            Cancel
          </Button>
          <Button onClick={handleSubmit} loading={submitting} disabled={!name.trim()}>
            {mode === "edit" ? "Save Changes" : "Create Key"}
          </Button>
        </>
      }
    >
      <div className="flex flex-col gap-4">
        <SegmentedControl
          size="sm"
          value={activeTab}
          onChange={setActiveTab}
          className="w-full"
          options={[
            { value: "general", label: "General", icon: "badge" },
            {
              value: "limits",
              label: hasLimitsConfigured ? "Limits & Access •" : "Limits & Access",
              icon: "tune",
            },
            { value: "permissions", label: "Permissions", icon: "admin_panel_settings" },
          ]}
        />

        {error && (
          <div
            role="alert"
            className="flex items-start gap-2 rounded-xl border border-danger-border bg-danger-bg px-3.5 py-2.5 text-xs text-danger-fg"
          >
            <span className="material-symbols-outlined text-[16px] shrink-0 mt-0.5">error</span>
            <span>{error}</span>
          </div>
        )}

        {activeTab === "general" && (
          <div className="flex flex-col gap-4">
            <Input
              label="Key Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Production Key"
              hint="A unique label to identify this key across logs and usage reports."
            />

            {mode === "create" && (
              <>
                <Input
                  label="Custom API Key (Optional)"
                  value={customKey}
                  maxLength={appendSecret ? 19 : 35}
                  onChange={(e) => setCustomKey(e.target.value)}
                  placeholder="e.g. prod-team, dev-cli (leave empty to auto-generate)"
                  hint={
                    appendSecret
                      ? "Encrypt mode: up to 19 characters; a secure random suffix is appended (35 chars total)."
                      : "Non-encrypt mode: up to 35 characters; saved directly as entered."
                  }
                />
                {customKey.trim() && (
                  <div className="flex items-center justify-between gap-3 rounded-xl border border-border-subtle bg-surface-2/60 px-3.5 py-2.5 text-xs">
                    <div className="flex flex-col min-w-0">
                      <span className="font-medium text-text-main">
                        {appendSecret ? "Encrypt Mode (Recommended)" : "Non-Encrypt Mode"}
                      </span>
                      <span className="text-text-muted">
                        {appendSecret
                          ? "Appends a cryptographic suffix so the key reaches 35 characters."
                          : `Exact custom key length: ${customKeySanitizedLen}/35 characters.`}
                      </span>
                    </div>
                    <Toggle checked={appendSecret} onChange={setAppendSecret} />
                  </div>
                )}
              </>
            )}

            <div className="flex flex-col gap-1">
              <div className="flex items-end gap-2">
                <div className="flex-1 min-w-0">
                  <Input
                    label="Expiry Date (Optional)"
                    type="datetime-local"
                    value={expiresAt}
                    onChange={(e) => setExpiresAt(e.target.value)}
                    hint="Key stops authenticating requests after this date; leave empty for no expiry."
                  />
                </div>
                {expiresAt && (
                  <Button
                    size="sm"
                    variant="ghost"
                    className="mb-6"
                    onClick={() => setExpiresAt("")}
                  >
                    Clear
                  </Button>
                )}
              </div>
            </div>
          </div>
        )}

        {activeTab === "limits" && (
          <div className="flex flex-col gap-4">
            <Input
              label="Token Limit (0 = unlimited)"
              type="number"
              min="0"
              value={tokenLimit}
              onChange={(e) => setTokenLimit(e.target.value)}
              placeholder="e.g. 10000000"
              hint={
                creatorTokenLimit > 0
                  ? `Your maximum allowed limit is ${creatorTokenLimit.toLocaleString()} tokens.`
                  : "Leave empty or set 0 for unlimited token usage."
              }
            />

            {Number(tokenLimit) > 0 && (
              <Select
                label="Auto-Reset Interval"
                options={RESET_INTERVAL_OPTIONS}
                value={resetInterval}
                onChange={(e) => setResetInterval(e.target.value)}
              />
            )}

            {Number(tokenLimit) > 0 && resetInterval === "custom" && (
              <Input
                label="Custom Reset Interval"
                value={customReset}
                onChange={(e) => setCustomReset(e.target.value)}
                placeholder="e.g. 30m, 10h, 3d"
                hint="Use a positive number ending in m (minutes), h (hours), or d (days)."
              />
            )}

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <Input
                label="RPM Limit (0 = unlimited)"
                type="number"
                min="0"
                value={rpmLimit}
                onChange={(e) => setRpmLimit(e.target.value)}
                placeholder="0"
                hint="Max requests per minute"
              />
              <Input
                label="TPM Limit (0 = unlimited)"
                type="number"
                min="0"
                value={tpmLimit}
                onChange={(e) => setTpmLimit(e.target.value)}
                placeholder="0"
                hint="Max tokens per minute"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <div className="flex items-center justify-between">
                <label className="text-sm font-medium text-text-main">Allowed Models</label>
                <Button
                  size="sm"
                  variant="secondary"
                  icon="add"
                  onClick={onOpenModelPicker}
                >
                  Select Models
                </Button>
              </div>
              <Input
                value={selectedModelsValue}
                readOnly
                inputClassName="truncate font-mono"
                hint={
                  isApiKeyUser && creatorAllowedModels !== "*"
                    ? `Limited to your allowed models: ${creatorAllowedModels}`
                    : "Pick specific models or keep * to allow all models."
                }
              />
              {parseModelsList(selectedModelsValue).length > 0 && (
                <div className="flex flex-wrap gap-1.5 mt-1">
                  {parseModelsList(selectedModelsValue).map((m) => (
                    <span
                      key={m}
                      className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-primary/10 border border-primary/20 text-primary font-mono text-xs"
                    >
                      {m}
                      <button
                        type="button"
                        onClick={() => removeModel(m)}
                        aria-label={`Remove ${m}`}
                        className="hover:text-danger transition-colors cursor-pointer"
                      >
                        <span className="material-symbols-outlined text-[14px]">close</span>
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </div>

            <Input
              label="IP Whitelist (Optional)"
              value={ipWhitelist}
              onChange={(e) => setIpWhitelist(e.target.value)}
              placeholder="e.g. 203.0.113.5, 198.51.100.0/24"
              hint="Comma-separated IPs or CIDR ranges. Leave empty to allow any IP address."
            />
          </div>
        )}

        {activeTab === "permissions" && (
          <PermissionsEditor
            value={permissions}
            onChange={setPermissions}
            allowed={creatorPermissions}
            locked={permissionsLocked}
            lockedReason={PERMISSIONS_LOCKED_REASON}
          />
        )}
      </div>
    </Modal>
  );
}

KeyFormModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  mode: PropTypes.oneOf(["create", "edit"]),
  initialKey: PropTypes.object,
  existingKeys: PropTypes.array,
  creatorPermissions: PropTypes.object,
  creatorTokenLimit: PropTypes.number,
  creatorAllowedModels: PropTypes.string,
  isApiKeyUser: PropTypes.bool,
  onOpenModelPicker: PropTypes.func.isRequired,
  selectedModelsValue: PropTypes.string,
  onSelectedModelsChange: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};
