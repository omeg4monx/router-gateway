# 9Router Embeddings Tests

Unit tests for the `/v1/embeddings` endpoint implementation.

## Setup

From the repository root, install the app and test dependencies using the pinned Bun version in `package.json` and the Node version in `.nvmrc`:

```bash
bun install --frozen-lockfile
(cd tests && bun install --frozen-lockfile)
```

## Running Tests

Run from the repository root:

```bash
RUN_REAL=0 bun run test
bun run test tests/unit/embeddingsCore.test.js
```

## Test Files

| File | What it tests |
|------|--------------|
| `unit/embeddingsCore.test.js` | `open-sse/handlers/embeddingsCore.js` — core logic: body builder, URL router, headers, handler flow |

## Coverage Summary (36 tests)

### `embeddingsCore.test.js` (36 tests)
- `buildEmbeddingsBody`: single string, array, encoding_format, default float
- `buildEmbeddingsUrl`: openai, openrouter, openai-compatible-*, unsupported providers
- `buildEmbeddingsHeaders`: per-provider header sets, fallback to accessToken
- `handleEmbeddingsCore` input validation: missing, wrong type, null, empty
- `handleEmbeddingsCore` success: response format, CORS, Content-Type, callbacks
- `handleEmbeddingsCore` errors: 400/429/500, network error, invalid JSON
- `handleEmbeddingsCore` token refresh: 401 retry, graceful fallback
