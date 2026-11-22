import { describe, it, expect } from "vitest";
import fs from "node:fs";
import path from "node:path";

const ROOT = process.cwd();
const config = fs.readFileSync(
  path.join(ROOT, "src/shared/constants/config.js"),
  "utf8"
);

describe("update banner installs the fork, not upstream npm", () => {
  it("the published command points at the fork repository", () => {
    expect(config).toContain("github:omeg4monx/router-gateway");
    expect(config).toContain("installCmdLatest");
  });

  it("no longer hands out the upstream npm package name", () => {
    expect(config).not.toMatch(/npm i -g 9router(@latest)?[\"'\s]/);
  });
});