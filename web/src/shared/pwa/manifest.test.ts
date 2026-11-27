import { describe, expect, it } from "bun:test";
import { existsSync, readFileSync } from "node:fs";
import { resolve } from "node:path";
import { THEME_COLOR_BY_MODE } from "@/shared/store/uiStore";

const WEB_ROOT = resolve(import.meta.dirname, "../../..");
const PUBLIC_DIR = resolve(WEB_ROOT, "public");
const PNG_MAGIC = [0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a];

function hasPngMagicBytes(filePath: string): boolean {
  const buf = readFileSync(filePath);
  if (buf.length < PNG_MAGIC.length) {
    return false;
  }
  return PNG_MAGIC.every((byte, idx) => buf[idx] === byte);
}

interface WebAppManifestIcon {
  src: string;
  sizes: string;
  type: string;
  purpose?: string;
}

interface WebAppManifest {
  id?: string;
  name: string;
  short_name: string;
  start_url: string;
  scope: string;
  display: string;
  theme_color: string;
  background_color: string;
  icons: WebAppManifestIcon[];
}

describe("PWA manifest, icons, and HTML metadata", () => {
  it("provides a valid standalone Web App Manifest with existing PNG icons", () => {
    const manifestPath = resolve(PUBLIC_DIR, "manifest.webmanifest");
    expect(existsSync(manifestPath)).toBe(true);

    const manifest = JSON.parse(readFileSync(manifestPath, "utf-8")) as WebAppManifest;
    expect(manifest.name).toBe("OMEG4MON ROUTER");
    expect(manifest.short_name).toBe("OMEG4MON");
    expect(manifest.display).toBe("standalone");
    expect(manifest.start_url).toBe("/dashboard/");
    expect(manifest.scope).toBe("/");
    expect(manifest.theme_color).toBe("#050608");
    expect(manifest.background_color).toBe("#050608");

    const has192 = manifest.icons.some(
      (icon) => icon.sizes === "192x192" && icon.type === "image/png",
    );
    const has512 = manifest.icons.some(
      (icon) =>
        icon.sizes === "512x512" && icon.type === "image/png" && icon.purpose !== "maskable",
    );
    const hasMaskable512 = manifest.icons.some(
      (icon) =>
        icon.sizes === "512x512" &&
        icon.type === "image/png" &&
        (icon.purpose ?? "").includes("maskable"),
    );

    expect(has192).toBe(true);
    expect(has512).toBe(true);
    expect(hasMaskable512).toBe(true);

    for (const icon of manifest.icons) {
      const iconPath = resolve(PUBLIC_DIR, icon.src.replace(/^\//, ""));
      expect(existsSync(iconPath)).toBe(true);
      if (icon.type === "image/png") {
        expect(hasPngMagicBytes(iconPath)).toBe(true);
      }
    }
  });

  it("includes apple-touch-icon PNG and PWA links/meta in index.html", () => {
    const appleIconPath = resolve(PUBLIC_DIR, "icons/apple-touch-icon.png");
    expect(existsSync(appleIconPath)).toBe(true);
    expect(hasPngMagicBytes(appleIconPath)).toBe(true);

    const html = readFileSync(resolve(WEB_ROOT, "index.html"), "utf-8");
    expect(html).toContain('rel="manifest" href="/manifest.webmanifest"');
    expect(html).toContain('rel="apple-touch-icon"');
    expect(html).toContain('name="theme-color"');
    expect(html).toContain("viewport-fit=cover");
    expect(html).toContain('name="apple-mobile-web-app-capable"');
  });

  it("maps light and dark theme modes to Omegamon canvas theme-color hex values", () => {
    expect(THEME_COLOR_BY_MODE.light).toBe("#f5f6f8");
    expect(THEME_COLOR_BY_MODE.dark).toBe("#050608");
  });
});
