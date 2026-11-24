import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const appShellSource = readFileSync(resolve(import.meta.dir, "AppShell.tsx"), "utf8");

describe("AppShell fixed desktop sidebar and layout contracts", () => {
  test("pins desktop aside rail with sticky top-2.5 self-start and viewport height so it never scrolls with page content", () => {
    const asideMatch = appShellSource.match(/<aside[\s\S]*?aria-label="Application navigation"/);
    expect(asideMatch).not.toBeNull();
    const asideBlock = asideMatch?.[0] ?? "";
    expect(asideBlock).toContain("sticky");
    expect(asideBlock).toContain("top-2.5");
    expect(asideBlock).toContain("self-start");
    expect(asideBlock).toContain("h-[calc(100dvh-20px)]");
    expect(asideBlock).toContain("overflow-hidden");
    expect(asideBlock).not.toContain("shell-rail relative");
  });

  test("constrains sidebar flex container with min-h-0 so nav scrolls internally while header and footer remain docked", () => {
    expect(appShellSource).toContain(
      'className="flex h-full min-h-0 flex-col text-[var(--color-text-main)]"',
    );
    expect(appShellSource).toContain(
      "custom-scrollbar min-h-0 flex-1 space-y-3 overflow-y-auto overflow-x-hidden",
    );
  });

  test("keeps top dashboard header sticky at top-0", () => {
    expect(appShellSource).toContain('<header className="sticky top-0 z-20');
  });

  test("requires confirmation dialog before triggering logout mutation", () => {
    expect(appShellSource).toContain("signOutConfirmOpen");
    expect(appShellSource).toContain('aria-labelledby="sign-out-confirm-title"');
    expect(appShellSource).toContain("onClick={openSignOutConfirm}");
    expect(appShellSource).toContain("onClick={handleConfirmSignOut}");
  });
});
