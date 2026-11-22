# syntax=docker/dockerfile:1

FROM oven/bun:1-alpine AS web-builder
WORKDIR /src/web
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY contracts /src/contracts
COPY web/ ./
ENV NEXT_OUTPUT_EXPORT=true
RUN bun run build

FROM golang:1.26-alpine AS go-builder
ARG VERSION=dev
ARG COMMIT=none
ARG BUILD_TIME=unknown
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY --from=web-builder /src/web/out ./cmd/omeg4mon-router/ui_dist
RUN mkdir -p /data && \
    CGO_ENABLED=0 go build \
      -tags embedui \
      -trimpath \
      -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
      -o /out/omeg4mon-router \
      ./cmd/omeg4mon-router

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=go-builder /out/omeg4mon-router /omeg4mon-router
COPY --from=go-builder --chown=65532:65532 /data /data
ENV OMEG4MON_HOST=0.0.0.0 \
    OMEG4MON_PORT=20128 \
    OMEG4MON_DATA_DIR=/data
VOLUME ["/data"]
EXPOSE 20128
ENTRYPOINT ["/omeg4mon-router"]
