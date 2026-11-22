"use client";

import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Badge, Button, Input, Modal, Select } from "@/shared/components";
import ProviderLogoField from "@/shared/components/ProviderLogoField";
import { useNotificationStore } from "@/store/notificationStore";

const VARIANT_CONFIG = {
  openai: {
    title: "Add OpenAI Compatible",
    type: "openai-compatible",
    defaultBaseUrl: "https://api.openai.com/v1",
    namePlaceholder: "OpenAI Compatible (Prod)",
    prefixPlaceholder: "oc-prod",
    baseUrlHint: "Use the base URL (ending in /v1) for your OpenAI-compatible API.",
    modelIdPlaceholder: "e.g. gpt-4o, claude-3-7-sonnet",
    errorLabel: "OpenAI Compatible",
    hasApiType: true,
  },
  anthropic: {
    title: "Add Anthropic Compatible",
    type: "anthropic-compatible",
    defaultBaseUrl: "https://api.anthropic.com/v1",
    namePlaceholder: "Anthropic Compatible (Prod)",
    prefixPlaceholder: "ac-prod",
    baseUrlHint:
      "Use the base URL (ending in /v1) for your Anthropic-compatible API. The system will append /messages.",
    modelIdPlaceholder: "e.g. claude-3-7-sonnet",
    errorLabel: "Anthropic Compatible",
    hasApiType: false,
  },
  moonshot: {
    title: "Add MoonshotAI Compatible",
    type: "openai-compatible",
    defaultBaseUrl: "https://api.moonshot.ai/v1",
    namePlaceholder: "MoonshotAI Compatible (Prod)",
    prefixPlaceholder: "mk-prod",
    baseUrlHint:
      "Use the base URL (ending in /v1) for your MoonshotAI (Kimi) OpenAI-compatible API. Docs: platform.kimi.ai",
    modelIdPlaceholder: "e.g. kimi-k2.6, moonshot-v1-8k",
    errorLabel: "MoonshotAI Compatible",
    hasApiType: true,
  },
};

const API_TYPE_OPTIONS = [
  { value: "chat", label: "Chat Completions" },
  { value: "responses", label: "Responses API" },
];

function createInitialFormData(variant) {
  const cfg = VARIANT_CONFIG[variant] || VARIANT_CONFIG.openai;
  return {
    name: "",
    prefix: "",
    ...(cfg.hasApiType ? { apiType: "chat" } : {}),
    baseUrl: cfg.defaultBaseUrl,
    logo: "",
  };
}

function AddCompatibleModal({ variant, isOpen, onClose, onCreated }) {
  const config = VARIANT_CONFIG[variant] || VARIANT_CONFIG.openai;

  const [formData, setFormData] = useState(() => createInitialFormData(variant));
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState("");
  const [checkKey, setCheckKey] = useState("");
  const [checkModelId, setCheckModelId] = useState("");
  const [validating, setValidating] = useState(false);
  const [validationResult, setValidationResult] = useState(null);

  // Reset form and validation cleanly each time the modal opens or switches variant,
  // without wiping a user-edited Base URL when toggling API Type.
  useEffect(() => {
    if (!isOpen) return;
    setFormData(createInitialFormData(variant));
    setCheckKey("");
    setCheckModelId("");
    setValidationResult(null);
    setSubmitError("");
  }, [isOpen, variant]);

  const canSubmit =
    Boolean(formData.name.trim()) &&
    Boolean(formData.prefix.trim()) &&
    Boolean(formData.baseUrl.trim()) &&
    !submitting;

  const handleSubmit = async () => {
    if (!canSubmit) return;
    setSubmitting(true);
    setSubmitError("");
    try {
      const res = await fetch("/api/provider-nodes", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: formData.name.trim(),
          prefix: formData.prefix.trim(),
          ...(config.hasApiType ? { apiType: formData.apiType } : {}),
          baseUrl: formData.baseUrl.trim(),
          type: config.type,
          brand: variant || undefined,
          logo: formData.logo,
        }),
      });
      const data = await res.json();
      if (res.ok) {
        useNotificationStore
          .getState()
          .success(`Added ${data.node?.name || formData.name.trim()}`);
        onCreated(data.node);
      } else {
        setSubmitError(
          data?.error || `Failed to create ${config.errorLabel} node`,
        );
      }
    } catch (error) {
      setSubmitError(error?.message || "Network error while creating node");
    } finally {
      setSubmitting(false);
    }
  };

  const handleValidate = async () => {
    setValidating(true);
    setValidationResult(null);
    try {
      const res = await fetch("/api/provider-nodes/validate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          baseUrl: formData.baseUrl.trim(),
          apiKey: checkKey,
          type: config.type,
          modelId: checkModelId.trim() || undefined,
        }),
      });
      const data = await res.json();
      setValidationResult(data);
    } catch {
      setValidationResult({ valid: false, error: "Network error" });
    } finally {
      setValidating(false);
    }
  };

  const renderValidationResult = () => {
    if (!validationResult) return null;
    const { valid, error, method } = validationResult;
    if (valid) {
      return (
        <div className="flex items-center gap-2">
          <Badge variant="success">Valid</Badge>
          {method === "chat" && (
            <span className="text-xs text-text-muted">(via inference test)</span>
          )}
        </div>
      );
    }
    return (
      <div className="flex flex-col gap-1">
        <Badge variant="error">Invalid</Badge>
        {error && <span className="text-xs text-danger break-words">{error}</span>}
      </div>
    );
  };

  const modalFooter = (
    <div className="flex items-center justify-end gap-2">
      <Button onClick={onClose} variant="ghost" disabled={submitting}>
        Cancel
      </Button>
      <Button onClick={handleSubmit} disabled={!canSubmit} loading={submitting}>
        {submitting ? "Creating..." : "Create Endpoint"}
      </Button>
    </div>
  );

  return (
    <Modal
      isOpen={isOpen}
      title={config.title}
      onClose={onClose}
      footer={modalFooter}
    >
      <div className="flex flex-col gap-4">
        <ProviderLogoField
          logo={formData.logo}
          onChange={(next) => setFormData((prev) => ({ ...prev, logo: next }))}
        />
        <Input
          label="Name"
          value={formData.name}
          onChange={(e) =>
            setFormData((prev) => ({ ...prev, name: e.target.value }))
          }
          placeholder={config.namePlaceholder}
          hint="A friendly label for this node"
        />
        <Input
          label="Prefix"
          value={formData.prefix}
          onChange={(e) =>
            setFormData((prev) => ({ ...prev, prefix: e.target.value }))
          }
          placeholder={config.prefixPlaceholder}
          hint="The provider prefix used in model IDs"
        />
        {config.hasApiType && (
          <Select
            label="API Type"
            options={API_TYPE_OPTIONS}
            value={formData.apiType}
            onChange={(e) =>
              setFormData((prev) => ({ ...prev, apiType: e.target.value }))
            }
          />
        )}
        <Input
          label="Base URL"
          value={formData.baseUrl}
          onChange={(e) =>
            setFormData((prev) => ({ ...prev, baseUrl: e.target.value }))
          }
          placeholder={config.defaultBaseUrl}
          hint={config.baseUrlHint}
        />

        <div className="rounded-xl border border-border-subtle bg-surface-2/50 p-3 flex flex-col gap-3">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-text-muted">
            Optional Connection Check
          </p>
          <Input
            label="API Key (for Check)"
            type="password"
            value={checkKey}
            onChange={(e) => setCheckKey(e.target.value)}
            placeholder="Not saved — used only to test reachability"
          />
          <Input
            label="Model ID (optional)"
            value={checkModelId}
            onChange={(e) => setCheckModelId(e.target.value)}
            placeholder={config.modelIdPlaceholder}
            hint="If provider lacks /models endpoint, enter a model ID to validate via inference instead."
          />
          <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
            <Button
              onClick={handleValidate}
              disabled={!checkKey || validating || !formData.baseUrl.trim()}
              variant="secondary"
              size="sm"
              className="w-full sm:w-auto"
            >
              {validating ? "Checking..." : "Check Endpoint"}
            </Button>
            {renderValidationResult()}
          </div>
        </div>

        {submitError && (
          <div
            role="alert"
            className="flex items-center gap-2 p-2.5 rounded-xl bg-danger/10 border border-danger/30 text-danger text-xs"
          >
            <span className="material-symbols-outlined text-[16px] shrink-0">
              error
            </span>
            <span>{submitError}</span>
          </div>
        )}
      </div>
    </Modal>
  );
}

AddCompatibleModal.propTypes = {
  variant: PropTypes.oneOf(["openai", "anthropic", "moonshot"]).isRequired,
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  onCreated: PropTypes.func.isRequired,
};

export default AddCompatibleModal;
