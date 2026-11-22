const { readFileSync } = require("node:fs");
const { join } = require("node:path");

if (process.versions.bun || process.env.npm_config_user_agent?.startsWith("bun/")) {
  const expected = `v${readFileSync(join(__dirname, "../.nvmrc"), "utf8").trim()}`;
  if (process.versions.bun || process.version !== expected) {
    const actual = process.versions.bun ? `Bun ${process.versions.bun}` : `Node.js ${process.version}`;
    console.error(`Web scripts require Node.js ${expected} on PATH (found ${actual}).`);
    process.exit(1);
  }
}
