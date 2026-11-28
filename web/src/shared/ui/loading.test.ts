import { describe, expect, test } from "bun:test";
import { createElement } from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { Button } from "./button";
import { PageLoadingScreen, RouteLoadingSkeleton, Skeleton } from "./skeleton";
import { Spinner } from "./spinner";

describe("Shared loading UI primitives", () => {
  test("Spinner renders accessible status when label is given and decorative icon otherwise", () => {
    const decorativeHtml = renderToStaticMarkup(
      createElement(Spinner, { size: "sm", decorative: true }),
    );
    expect(decorativeHtml).toContain('aria-hidden="true"');
    expect(decorativeHtml).toContain("animate-spin");

    const statusHtml = renderToStaticMarkup(
      createElement(Spinner, { size: "md", label: "Loading providers" }),
    );
    expect(statusHtml).toContain('role="status"');
    expect(statusHtml).toContain("Loading providers");
  });

  test("Skeleton renders shimmer surface hidden from assistive tech", () => {
    const html = renderToStaticMarkup(createElement(Skeleton, { className: "h-6 w-32" }));
    expect(html).toContain("loading-shimmer");
    expect(html).toContain('aria-hidden="true"');
    expect(html).toContain("h-6 w-32");
  });

  test("RouteLoadingSkeleton and PageLoadingScreen expose busy status semantics", () => {
    const routeHtml = renderToStaticMarkup(createElement(RouteLoadingSkeleton));
    expect(routeHtml).toContain('role="status"');
    expect(routeHtml).toContain('aria-busy="true"');
    expect(routeHtml).toContain("loading-shimmer");

    const pageHtml = renderToStaticMarkup(createElement(PageLoadingScreen));
    expect(pageHtml).toContain('role="status"');
    expect(pageHtml).toContain('aria-busy="true"');
    expect(pageHtml).toContain("OMEG4MON ROUTER");
  });

  test("Button with loading prop disables control, sets aria-busy, and renders inline spinner", () => {
    const loadingHtml = renderToStaticMarkup(createElement(Button, { loading: true }, "Saving..."));
    expect(loadingHtml).toContain('disabled=""');
    expect(loadingHtml).toContain('aria-busy="true"');
    expect(loadingHtml).toContain("animate-spin");
    expect(loadingHtml).toContain("Saving...");

    const idleHtml = renderToStaticMarkup(createElement(Button, null, "Save"));
    expect(idleHtml).not.toContain('aria-busy="true"');
    expect(idleHtml).not.toContain("animate-spin");
  });
});
