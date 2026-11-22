import { NextResponse } from "next/server";
import { getMeta, setMeta } from "@/lib/db/helpers/metaStore";
import { getSessionContext } from "@/lib/auth/dashboardPermissions";

export const dynamic = "force-dynamic";

// GET /api/theme - fetch public/global theme preference set by instance owner
export async function GET() {
  try {
    const raw = await getMeta("global_theme_settings", null);
    if (!raw) {
      return NextResponse.json({
        ok: true,
        theme: "ember",
        accent: "orange",
        customColor: "#E56A4A",
      });
    }
    const parsed = JSON.parse(raw);
    return NextResponse.json({ ok: true, ...parsed });
  } catch (e) {
    return NextResponse.json({
      ok: false,
      theme: "ember",
      accent: "orange",
      customColor: "#E56A4A",
    });
  }
}

// POST /api/theme - update global theme preference (instance owner only).
// dashboardGuard already refuses anonymous and key-session writes; this is the
// route's own check so the rule survives a guard allow-list edit.
export async function POST(request) {
  try {
    const { session } = await getSessionContext();
    if (session?.role === "apikey") {
      return NextResponse.json({ ok: false, error: "Only the instance owner can change the site theme" }, { status: 403 });
    }
    const body = await request.json();
    const theme = body.theme || "ember";
    const accent = body.accent || "orange";
    const customColor = body.customColor || "#E56A4A";

    await setMeta("global_theme_settings", JSON.stringify({ theme, accent, customColor }));
    return NextResponse.json({ ok: true, theme, accent, customColor });
  } catch (e) {
    return NextResponse.json({ ok: false, error: e.message }, { status: 500 });
  }
}
