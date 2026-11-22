"use client";

import { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Modal, Button } from "@/shared/components";
import { translate } from "@/i18n/runtime";

/**
 * Xiaomi MiMo Desktop Plan auth modal.
 *
 * Desktop Plan = Xiaomi account weekly quota (mimo-v2.6 family), connected by
 * importing local MiMo Desktop credentials (~/.local/share/mimocode/auth.json).
 *
 * Token Plan (cloud sk- API key) uses the standard "API Key" entry point.
 */

export default function XiaomiMimoAuthModal({ isOpen, onSuccess, onClose }) {
  const [phase, setPhase] = useState("detecting"); // detecting | found | not-found | importing
  const [detectResult, setDetectResult] = useState(null);
  const [existingConnection, setExistingConnection] = useState(null);
  const [error, setError] = useState(null);

  // Detect local credentials when the modal opens (non-blocking: never trap the spinner)
  useEffect(() => {
    if (!isOpen) return;

    setPhase("detecting");
    setError(null);
    setDetectResult(null);
    setExistingConnection(null);

    const runDetect = async () => {
      try {
        const res = await fetch("/api/oauth/xiaomi-mimo/auto-import", {
          signal: AbortSignal.timeout(5000),
        });
        const data = await res.json();

        if (data.found && data.apiKey) {
          setDetectResult(data);
          setPhase("found");

          // Non-blocking: flag an already-imported account
          fetch("/api/providers", { signal: AbortSignal.timeout(3000) })
            .then((r) => r.json())
            .then((provData) => {
              const foundConn = (provData.connections || []).find(
                (c) =>
                  c.provider === "xiaomi-mimo" &&
                  data.uid &&
                  (c.email === `${data.uid}@xiaomi` ||
                    c.providerSpecificData?.uid === data.uid ||
                    c.providerSpecificData?.mimoUserId === data.uid),
              );
              if (foundConn) setExistingConnection(foundConn);
            })
            .catch(() => {});
        } else {
          setPhase("not-found");
        }
      } catch {
        setPhase("not-found");
      }
    };

    runDetect();
  }, [isOpen]);

  // Import the auto-detected local desktop credentials
  const handleImportLocal = async () => {
    if (!detectResult?.apiKey) return;
    setPhase("importing");
    setError(null);

    try {
      const res = await fetch("/api/oauth/xiaomi-mimo/api-key", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          apiKey: detectResult.apiKey,
          uid: detectResult.uid,
          baseUrl: detectResult.baseUrl,
          mimoPassToken: detectResult.mimoPassToken || null,
          mimoUserId: detectResult.mimoUserId || null,
          mimoCUserId: detectResult.mimoCUserId || null,
        }),
      });
      const data = await res.json();
      if (!res.ok || !data.success) throw new Error(data.error || "Import failed");
      onSuccess?.(data.connection);
      onClose();
    } catch (err) {
      setPhase("found");
      setError(err.message);
    }
  };

  return (
    <Modal isOpen={isOpen} title={translate("Connect Xiaomi MiMo")} onClose={onClose} size="lg">
      <div className="flex flex-col gap-3.5">
        {/* Detecting */}
        {phase === "detecting" && (
          <div className="text-center py-6">
            <div className="size-12 mx-auto mb-3 rounded-full bg-primary/10 flex items-center justify-center">
              <span className="material-symbols-outlined text-2xl text-primary animate-spin">
                progress_activity
              </span>
            </div>
            <p className="text-sm text-text-muted">Reading local MiMo Desktop credentials...</p>
          </div>
        )}

        {/* Importing */}
        {phase === "importing" && (
          <div className="text-center py-6">
            <div className="size-12 mx-auto mb-3 rounded-full bg-primary/10 flex items-center justify-center">
              <span className="material-symbols-outlined text-2xl text-primary animate-spin">
                progress_activity
              </span>
            </div>
            <p className="text-sm font-medium">Connecting...</p>
          </div>
        )}

        {/* Found local desktop credentials */}
        {phase === "found" && detectResult && (
          <div className="flex flex-col gap-3">
            <div className="flex items-center gap-1.5 text-xs font-semibold text-text-muted px-0.5">
              <span className="material-symbols-outlined text-primary text-sm">desktop_windows</span>
              <span>Desktop Plan · Local credentials</span>
            </div>

            {existingConnection ? (
              <div className="bg-info-bg p-3 rounded-lg border border-info-border">
                <div className="flex gap-2.5 items-start">
                  <span className="material-symbols-outlined text-info text-lg mt-0.5">
                    check_circle
                  </span>
                  <div className="text-sm text-info">
                    <p className="font-medium">This account is already connected (no need to import again)</p>
                    <p className="text-xs mt-0.5 opacity-80">
                      UID: {detectResult.uid || "–"} · Status: {existingConnection.testStatus === "active" ? "Active" : "Untested"}
                    </p>
                  </div>
                </div>
              </div>
            ) : (
              <div className="bg-success-bg p-3 rounded-lg border border-success-border">
                <div className="flex gap-2.5 items-start">
                  <span className="material-symbols-outlined text-success text-lg mt-0.5">
                    check_circle
                  </span>
                  <div className="text-sm text-success">
                    <p className="font-medium">Xiaomi MiMo Desktop credentials found!</p>
                    <p className="text-xs mt-0.5 opacity-80">
                      UID: {detectResult.uid || "–"} · Source: {detectResult.source?.split(/[\\/]/).pop()}
                    </p>
                  </div>
                </div>
              </div>
            )}

            {error && (
              <div className="bg-danger-bg p-2.5 rounded-lg border border-danger-border">
                <p className="text-xs text-danger">{translate(error)}</p>
              </div>
            )}

            <div className="flex gap-2">
              <Button onClick={handleImportLocal} variant={existingConnection ? "outline" : "default"} fullWidth>
                {existingConnection ? "Re-sync local credentials" : "Connect with local credentials"}
              </Button>
              <Button onClick={onClose} variant="ghost" fullWidth>
                {existingConnection ? "Close" : "Cancel"}
              </Button>
            </div>
          </div>
        )}

        {/* No local desktop credentials */}
        {phase === "not-found" && (
          <div className="flex flex-col gap-3">
            <div className="bg-warning-bg p-3 rounded-lg border border-warning-border">
              <div className="flex gap-2.5 items-start">
                <span className="material-symbols-outlined text-warning text-lg mt-0.5">info</span>
                <div className="text-sm text-warning">
                  <p className="font-medium">No local Desktop credentials found</p>
                  <p className="text-xs mt-0.5 opacity-80">
                    Install and sign in to MiMo Desktop, then reopen this dialog — or use the "API Key" option for a Token Plan key.
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-1">
              <Button onClick={onClose} variant="ghost" fullWidth>
                Cancel
              </Button>
            </div>
          </div>
        )}
      </div>
    </Modal>
  );
}

XiaomiMimoAuthModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onSuccess: PropTypes.func,
  onClose: PropTypes.func.isRequired,
};
