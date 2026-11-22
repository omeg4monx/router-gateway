# 9Router Embeddings Tests

Unit tests for the `/v1/embeddings` endpoint implementation.

## Setup

Install test dependencies from the `tests/` directory:

```bash
cd tests/ && npm install
```

## Running Tests

From the `tests/` directory:

```bash
npm test
```

Or run vitest directly with npx:

```bash
npx vitest run --reporter=verbose --config ./vitest.config.js
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
