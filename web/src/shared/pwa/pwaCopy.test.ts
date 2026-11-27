import { describe, expect, it } from "bun:test";
import { getPwaCopy, isStandaloneDisplayMode } from "./usePwaStatus";

describe("PWA copy and standalone detection", () => {
  it("provides distinct bilingual EN and ID copy for install, offline, and update states", () => {
    const en = getPwaCopy("en");
    const id = getPwaCopy("id");

    expect(en.installLabel.length).toBeGreaterThan(0);
    expect(id.installLabel.length).toBeGreaterThan(0);
    expect(en.installLabel).not.toBe(id.installLabel);

    expect(en.offlineTitle).not.toBe(id.offlineTitle);
    expect(en.offlineDescription.length).toBeGreaterThan(0);
    expect(id.offlineDescription.length).toBeGreaterThan(0);

    expect(en.updateTitle).not.toBe(id.updateTitle);
    expect(en.updateAction).not.toBe(id.updateAction);
    expect(en.dismissAction).not.toBe(id.dismissAction);
  });

  it("returns false safely in non-browser environments", () => {
    expect(isStandaloneDisplayMode()).toBe(false);
  });
});
