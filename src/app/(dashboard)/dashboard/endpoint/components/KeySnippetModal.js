"use client";

import { useState } from "react";
import PropTypes from "prop-types";
import { Modal, Button, SegmentedControl } from "@/shared/components";
import { useCopyToClipboard } from "@/shared/hooks/useCopyToClipboard";
import { generateEndpointSnippet } from "../endpointSnippets";

export default function KeySnippetModal({
  apiKeyObj,
  activeOrigin,
  activeRouteLabel = "Local",
  defaultModel = "gpt-4o",
  onClose,
}) {
  const [format, setFormat] = useState("openai");
  const [lang, setLang] = useState("curl");
  const { copied, copy } = useCopyToClipboard();

  const code = apiKeyObj
    ? generateEndpointSnippet({
        format,
        lang,
        baseUrl: activeOrigin,
        apiKey: apiKeyObj.key,
        model: defaultModel,
      })
    : "";

  return (
    <Modal
      isOpen={Boolean(apiKeyObj)}
      title={`Code Snippet — ${apiKeyObj?.name || "API Key"}`}
      size="lg"
      onClose={onClose}
      footer={
        <>
          <Button variant="ghost" onClick={onClose}>
            Close
          </Button>
          <Button
            icon={copied === "key_snippet" ? "check" : "content_copy"}
            onClick={() => copy(code, "key_snippet")}
          >
            {copied === "key_snippet" ? "Copied!" : "Copy to Clipboard"}
          </Button>
        </>
      }
    >
      <div className="flex flex-col gap-3.5">
        <div className="flex flex-wrap items-center justify-between gap-2">
          <SegmentedControl
            size="sm"
            value={format}
            onChange={setFormat}
            options={[
              { value: "openai", label: "OpenAI (/v1)" },
              { value: "anthropic", label: "Anthropic (/v1/messages)" },
            ]}
          />
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg border border-border-subtle bg-surface-2 font-mono text-[11px] text-text-muted">
            <span className="material-symbols-outlined text-[14px] text-primary">hub</span>
            <span>Route: {activeRouteLabel}</span>
          </span>
        </div>

        <SegmentedControl
          size="sm"
          value={lang}
          onChange={setLang}
          className="w-full sm:w-auto"
          options={[
            { value: "curl", label: "cURL" },
            { value: "python", label: "Python" },
            { value: "node", label: "Node.js" },
            { value: "go", label: "Go" },
          ]}
        />

        <pre className="rounded-xl border border-border-subtle bg-bg-alt p-4 text-xs text-text-main font-mono overflow-x-auto max-h-72 whitespace-pre-wrap break-all leading-relaxed custom-scrollbar">
          {code}
        </pre>
      </div>
    </Modal>
  );
}

KeySnippetModal.propTypes = {
  apiKeyObj: PropTypes.object,
  activeOrigin: PropTypes.string,
  activeRouteLabel: PropTypes.string,
  defaultModel: PropTypes.string,
  onClose: PropTypes.func.isRequired,
};
