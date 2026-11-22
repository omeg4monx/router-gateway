import { describe, it, expect, vi, beforeEach, afterEach, afterAll } from "vitest";
import * as fsPromises from "fs/promises";
import * as childProcess from "child_process";
import { mkdirSync, rmSync } from "node:fs";
import { dirname, join } from "node:path";
import Database from "better-sqlite3";

// The route loads better-sqlite3 with a CommonJS `require`, which vi.mock does
// not intercept, so the tests read a real SQLite file under a throwaway home.
const { HOME } = vi.hoisted(() => ({ HOME: `${process.env.DATA_DIR}/cursor-auto-import-home` }));
const DARWIN_DB = join(HOME, "Library/Application Support/Cursor/User/globalStorage/state.vscdb");

vi.mock("next/server", () => ({
  NextResponse: {
    json: vi.fn((body, init) => ({
      status: init?.status || 200,
      body,
      json: async () => body,
    })),
  },
}));

vi.mock("os", () => ({
  default: { homedir: vi.fn(() => HOME) },
  homedir: vi.fn(() => HOME),
}));

// access() decides which candidate path the route picks.
vi.mock("fs/promises", () => ({
  access: vi.fn(),
  constants: { R_OK: 4 },
}));

// Strategy 2 shells out to `sqlite3` (and Linux probes `which cursor`); the
// route promisifies execFile, so a node-style callback mock is enough.
vi.mock("child_process", () => ({
  execFile: vi.fn(),
}));

function writeCursorDb(rows) {
  mkdirSync(dirname(DARWIN_DB), { recursive: true });
  const db = new Database(DARWIN_DB);
  db.exec("CREATE TABLE itemTable (key TEXT UNIQUE ON CONFLICT REPLACE, value BLOB)");
  const insert = db.prepare("INSERT INTO itemTable (key, value) VALUES (?, ?)");
  for (const [key, value] of Object.entries(rows)) insert.run(key, value);
  db.close();
}

function mockExecFile(stdoutFor = () => "") {
  vi.mocked(childProcess.execFile).mockImplementation((file, args, opts, cb) => {
    cb(null, { stdout: stdoutFor(file, args), stderr: "" });
  });
}

let GET;

describe("GET /api/oauth/cursor/auto-import", () => {
  const originalPlatform = process.platform;

  beforeEach(async () => {
    vi.clearAllMocks();
    rmSync(HOME, { recursive: true, force: true });
    mockExecFile();
    Object.defineProperty(process, "platform", { value: "darwin", writable: true });
    const mod = await import("../../src/app/api/oauth/cursor/auto-import/route.js");
    GET = mod.GET;
  });

  afterEach(() => {
    Object.defineProperty(process, "platform", { value: originalPlatform, writable: true });
  });

  afterAll(() => {
    rmSync(HOME, { recursive: true, force: true });
  });

  it("returns not-found listing every probed path when no cursor db is accessible", async () => {
    vi.mocked(fsPromises.access).mockRejectedValue(new Error("ENOENT"));

    const response = await GET();

    expect(response.body.found).toBe(false);
    expect(response.body.error).toContain("Cursor database not found. Checked locations:");
    expect(response.body.error).toContain(DARWIN_DB);
    expect(response.body.error).toContain("Cursor - Insiders");
  });

  it("extracts tokens using exact keys", async () => {
    vi.mocked(fsPromises.access).mockResolvedValue();
    writeCursorDb({
      "cursorAuth/accessToken": "test-token",
      "storage.serviceMachineId": "test-machine-id",
    });

    const response = await GET();

    expect(response.body).toEqual({ found: true, accessToken: "test-token", machineId: "test-machine-id" });
  });

  it("unwraps JSON-encoded string values", async () => {
    vi.mocked(fsPromises.access).mockResolvedValue();
    writeCursorDb({
      "cursorAuth/accessToken": '"json-token"',
      "storage.serviceMachineId": '"json-machine-id"',
    });

    const response = await GET();

    expect(response.body.found).toBe(true);
    expect(response.body.accessToken).toBe("json-token");
    expect(response.body.machineId).toBe("json-machine-id");
  });

  it("falls back to secondary keys in priority order", async () => {
    vi.mocked(fsPromises.access).mockResolvedValue();
    writeCursorDb({
      "cursorAuth/token": "alt-token",
      "telemetry.machineId": "alt-machine",
    });

    const response = await GET();

    expect(response.body).toEqual({ found: true, accessToken: "alt-token", machineId: "alt-machine" });
  });

  it("falls back to the sqlite3 CLI when the db cannot be opened in-process", async () => {
    // access() says readable, but no file exists, so better-sqlite3 throws.
    vi.mocked(fsPromises.access).mockResolvedValue();
    mockExecFile((file, args) => {
      if (file !== "sqlite3") return "";
      if (args[1].includes("cursorAuth/accessToken")) return '"cli-token"\n';
      if (args[1].includes("storage.serviceMachineId")) return "cli-machine\n";
      return "";
    });

    const response = await GET();

    expect(response.body).toEqual({ found: true, accessToken: "cli-token", machineId: "cli-machine" });
  });

  it("asks for manual paste when no strategy finds both tokens", async () => {
    vi.mocked(fsPromises.access).mockResolvedValue();
    writeCursorDb({ "cursorAuth/accessToken": "only-token" });

    const response = await GET();

    expect(response.body).toEqual({ found: false, windowsManual: true, dbPath: DARWIN_DB });
  });

  it("linux skips import when config exists but Cursor is not installed", async () => {
    Object.defineProperty(process, "platform", { value: "linux", writable: true });
    // First access() finds the db; the .desktop probe fails.
    vi.mocked(fsPromises.access)
      .mockResolvedValueOnce()
      .mockRejectedValue(new Error("ENOENT"));
    vi.mocked(childProcess.execFile).mockImplementation((file, args, opts, cb) => cb(new Error("not found")));

    const response = await GET();

    expect(response.body.found).toBe(false);
    expect(response.body.error).toContain("Cursor IDE does not appear to be installed");
  });
});
