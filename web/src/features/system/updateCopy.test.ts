import { describe, expect, it } from "bun:test";
import type { SystemUpdateStatusResponse } from "@/shared/api/client";
import {
  getUnavailableReasonMessage,
  getUpdateCopy,
  getUpdateStateLabel,
  isUpdateInProgress,
  shouldShowUpdateButton,
} from "./updateCopy";

function makeStatus(
  overrides: Partial<SystemUpdateStatusResponse> = {},
): SystemUpdateStatusResponse {
  return {
    currentVersion: "v0.3.0",
    currentCommit: "1111111",
    branch: "main",
    latestCommit: "2222222222222222222222222222222222222222",
    latestShortCommit: "2222222",
    latestMessage: "feat: auto-update",
    latestPublishedAt: "2026-11-28T10:00:00Z",
    latestUrl: "https://github.com/omeg4monx/router-gateway/commit/2222222",
    updateAvailable: false,
    canApply: false,
    mode: "binary",
    unavailableReason: "",
    state: "idle",
    message: "",
    error: "",
    checkedAt: "2026-11-28T10:00:00Z",
    updatedAt: "2026-11-28T10:00:00Z",
    ...overrides,
  };
}

describe("system update visibility & state helpers", () => {
  it("hides update button when status is undefined or up to date and idle", () => {
    expect(shouldShowUpdateButton(undefined)).toBe(false);
    expect(shouldShowUpdateButton(makeStatus({ updateAvailable: false, state: "idle" }))).toBe(
      false,
    );
  });

  it("shows update button when update is available or update is in progress", () => {
    expect(shouldShowUpdateButton(makeStatus({ updateAvailable: true, state: "idle" }))).toBe(true);
    expect(
      shouldShowUpdateButton(makeStatus({ updateAvailable: false, state: "downloading" })),
    ).toBe(true);
    expect(
      shouldShowUpdateButton(makeStatus({ updateAvailable: false, state: "restarting" })),
    ).toBe(true);
  });

  it("identifies active update states for polling", () => {
    expect(isUpdateInProgress("idle")).toBe(false);
    expect(isUpdateInProgress("failed")).toBe(false);
    expect(isUpdateInProgress("checking")).toBe(true);
    expect(isUpdateInProgress("downloading")).toBe(true);
    expect(isUpdateInProgress("building")).toBe(true);
    expect(isUpdateInProgress("verifying")).toBe(true);
    expect(isUpdateInProgress("applying")).toBe(true);
    expect(isUpdateInProgress("restarting")).toBe(true);
  });

  it("returns localized bilingual labels for en and id", () => {
    const en = getUpdateCopy("en");
    const id = getUpdateCopy("id");
    expect(en.buttonLabel).toBe("Update");
    expect(id.buttonLabel).toBe("Perbarui");
    expect(getUpdateStateLabel("en", "restarting")).toContain("Restarting");
    expect(getUpdateStateLabel("id", "restarting")).toContain("Memulai ulang");
    expect(getUnavailableReasonMessage("en", "dirty_worktree")).toContain("uncommitted");
    expect(getUnavailableReasonMessage("id", "dirty_worktree")).toContain("belum di-commit");
    expect(getUnavailableReasonMessage("id", "release_building")).toContain("build");
    expect(getUnavailableReasonMessage("en", "wrong_branch")).toContain("main");
    expect(getUnavailableReasonMessage("id", "wrong_branch")).toContain("main");
    expect(en.restartTimeoutNotice.length).toBeGreaterThan(0);
    expect(id.reloadAction.length).toBeGreaterThan(0);
  });
});
