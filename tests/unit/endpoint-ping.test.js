import { createServer } from "node:http";
import { once } from "node:events";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { clientPingAny, clientPingUrl } from "../../src/app/(dashboard)/dashboard/endpoint/endpointPing";

describe("gateway browser health probes", () => {
  let server;
  let origin;

  beforeAll(async () => {
    server = createServer((request, response) => {
      if (request.url === "/disconnected/api/health") {
        request.socket.destroy();
        return;
      }
      response.writeHead(request.url === "/online/api/health" ? 200 : 503).end();
    });
    server.listen(0, "127.0.0.1");
    await once(server, "listening");
    origin = `http://127.0.0.1:${server.address().port}`;
  });

  afterAll(async () => {
    if (server?.listening) await new Promise((resolve) => server.close(resolve));
  });

  it("reports a successful response as reachable", async () => {
    await expect(clientPingUrl(`${origin}/online`)).resolves.toBe(true);
  });

  it("reports an unsuccessful response or missing URL as unreachable", async () => {
    await expect(clientPingUrl(`${origin}/offline`)).resolves.toBe(false);
    await expect(clientPingUrl(undefined)).resolves.toBe(false);
  });

  it("reports a dropped connection as unreachable", async () => {
    await expect(clientPingUrl(`${origin}/disconnected`)).resolves.toBe(false);
    await expect(clientPingAny(`${origin}/disconnected`, `${origin}/offline`)).resolves.toBe(false);
  });

  it("reports reachability when any candidate succeeds", async () => {
    await expect(clientPingAny(null, `${origin}/offline`, `${origin}/online`)).resolves.toBe(true);
    await expect(clientPingAny(undefined, `${origin}/offline`)).resolves.toBe(false);
  });
});
