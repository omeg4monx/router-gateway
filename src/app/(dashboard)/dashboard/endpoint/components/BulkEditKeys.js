"use client";

import { useState } from "react";
import PropTypes from "prop-types";
import { Button, Modal, Input } from "@/shared/components";
import { BULK_RESET_INTERVAL_OPTIONS } from "../endpointUtils";

const QUOTA_FIELDS = ["tokenLimit", "resetInterval", "rpmLimit", "tpmLimit"];
const ACCESS_FIELDS = ["allowedModels", "deniedModels", "allowedProviders", "ipWhitelist", "expiresAt", "isActive"];

/**
 * Bulk edit keys: write one or more fields to every selected key in one call.
 * Each field is opt-in (checkbox): fields left off are not touched, so the
 * form never rewrites what the user did not mean to change.
 */
export default function BulkEditKeys({ selectedCount, onApply, onClose }) {
  const [fields, setFields] = useState({
    tokenLimit: { on: false, value: "" },
    resetInterval: { on: false, value: "never" },
    rpmLimit: { on: false, value: "" },
    tpmLimit: { on: false, value: "" },
    allowedModels: { on: false, value: "" },
    deniedModels: { on: false, value: "" },
    allowedProviders: { on: false, value: "" },
    ipWhitelist: { on: false, value: "" },
    expiresAt: { on: false, value: "" },
    isActive: { on: false, value: "true" },
  });
  const [applying, setApplying] = useState(false);
  const [error, setError] = useState("");

  const toggle = (name) =>
    setFields((prev) => ({ ...prev, [name]: { ...prev[name], on: !prev[name].on } }));
  const setValue = (name, value) =>
    setFields((prev) => ({ ...prev, [name]: { ...prev[name], value } }));

  const handleApply = async () => {
    const patch = {};
    // Deny-list and provider allowlist live in their own table (apiKeyAcls), so they
    // travel as a nested `acl` and are written once per selected key by the caller.
    const acl = {};
    for (const [name, { on, value }] of Object.entries(fields)) {
      if (!on) continue;
      if (name === "isActive") {
        patch.isActive = value === "true";
      } else if (name === "tokenLimit" || name === "rpmLimit" || name === "tpmLimit") {
        const n = Number(value);
        if (!Number.isFinite(n) || n < 0) {
          setError(`${labelFor(name)} must be a number of 0 or more (0 = unlimited)`);
          return;
        }
        patch[name] = n;
      } else if (name === "allowedModels") {
        if (!String(value).trim()) {
          setError("Allowed models pattern cannot be empty when enabled");
          return;
        }
        patch.allowedModels = String(value).trim();
      } else if (name === "deniedModels" || name === "allowedProviders") {
        if (!String(value).trim()) {
          setError(`${labelFor(name)} cannot be empty when enabled (turn the checkbox off to leave it alone)`);
          return;
        }
        acl[name] = String(value).trim();
      } else {
        patch[name] = value;
      }
    }
    if (!Object.keys(patch).length && !Object.keys(acl).length) {
      setError("Enable at least one field to change");
      return;
    }
    if (Object.keys(acl).length) patch.acl = acl;
    setError("");
    setApplying(true);
    try {
      await onApply(patch);
      onClose();
    } catch (err) {
      setError(err.message || "Failed to apply bulk changes");
    } finally {
      setApplying(false);
    }
  };

  const enabledCount = Object.values(fields).filter((f) => f.on).length;

  const renderRow = (name) => (
    <div
      key={name}
      className="flex items-start gap-3 rounded-xl border border-border-subtle bg-surface-2/50 p-3 transition-colors"
    >
      <input
        type="checkbox"
        checked={fields[name].on}
        onChange={() => toggle(name)}
        className="mt-1 size-4 accent-[var(--color-primary)] shrink-0 cursor-pointer"
        aria-label={`Change ${labelFor(name)}`}
      />
      <div className="flex-1 min-w-0">
        <label className="text-xs font-medium text-text-main block mb-1.5">{labelFor(name)}</label>
        {name === "resetInterval" ? (
          <select
            value={fields[name].value}
            onChange={(e) => setValue(name, e.target.value)}
            disabled={!fields[name].on}
            className="field-control w-full h-9 px-3 text-sm disabled:opacity-40"
          >
            {BULK_RESET_INTERVAL_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        ) : name === "isActive" ? (
          <select
            value={fields[name].value}
            onChange={(e) => setValue(name, e.target.value)}
            disabled={!fields[name].on}
            className="field-control w-full h-9 px-3 text-sm disabled:opacity-40"
          >
            <option value="true">Active</option>
            <option value="false">Disabled</option>
          </select>
        ) : (
          <Input
            value={fields[name].value}
            onChange={(e) => setValue(name, e.target.value)}
            disabled={!fields[name].on}
            placeholder={placeholderFor(name)}
            type={name === "expiresAt" ? "date" : "text"}
          />
        )}
      </div>
    </div>
  );

  return (
    <Modal
      isOpen
      onClose={onClose}
      size="lg"
      title={`Bulk edit ${selectedCount} key${selectedCount === 1 ? "" : "s"}`}
      footer={
        <>
          <Button variant="ghost" onClick={onClose} disabled={applying}>
            Cancel
          </Button>
          <Button onClick={handleApply} loading={applying} disabled={enabledCount === 0}>
            Apply to {selectedCount} key{selectedCount === 1 ? "" : "s"}
          </Button>
        </>
      }
    >
      <div className="flex flex-col gap-4">
        <p className="text-xs text-text-muted">
          Check only the fields you want to overwrite. Unchecked fields remain untouched on every selected key.
        </p>

        {error && (
          <div
            role="alert"
            className="rounded-lg border border-danger-border bg-danger-bg px-3 py-2 text-xs text-danger-fg"
          >
            {error}
          </div>
        )}

        <div className="space-y-2">
          <p className="micro-label">Quota &amp; Rate Limits</p>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
            {QUOTA_FIELDS.map(renderRow)}
          </div>
        </div>

        <div className="space-y-2">
          <p className="micro-label">Access &amp; Lifecycle</p>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
            {ACCESS_FIELDS.map(renderRow)}
          </div>
        </div>
      </div>
    </Modal>
  );
}

function labelFor(name) {
  return {
    tokenLimit: "Token limit (0 = unlimited)",
    resetInterval: "Auto-reset interval",
    rpmLimit: "Requests per minute (0 = unlimited)",
    tpmLimit: "Tokens per minute (0 = unlimited)",
    allowedModels: "Allowed models pattern",
    deniedModels: "Denied models pattern",
    allowedProviders: "Allowed providers",
    ipWhitelist: "IP whitelist",
    expiresAt: "Expiry date",
    isActive: "Status",
  }[name] || name;
}

function placeholderFor(name) {
  return {
    tokenLimit: "e.g. 1000000",
    rpmLimit: "e.g. 60",
    tpmLimit: "e.g. 100000",
    allowedModels: "e.g. anthropic/*, openai/gpt-4o",
    deniedModels: "e.g. gpt-4*, *opus",
    allowedProviders: "e.g. opencode, anthropic, glm",
    ipWhitelist: "e.g. 203.0.113.5, 198.51.100.0/24",
    expiresAt: "leave blank for no expiry",
  }[name] || "";
}

BulkEditKeys.propTypes = {
  selectedCount: PropTypes.number.isRequired,
  onApply: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};
