# syntax=docker/dockerfile:1

# docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .

FROM node:24-trixie AS asset
WORKDIR /node
RUN corepack enable
RUN --mount=type=bind,source=web/package.json,target=package.json \
    --mount=type=bind,source=web/pnpm-lock.yaml,target=pnpm-lock.yaml \
    --mount=type=cache,target=/root/.local/share/pnpm/store,sharing=locked \
    pnpm install --frozen-lockfile
COPY web .
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    pnpm build

FROM golang:1.26-trixie AS app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x
COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY --from=asset /node/dist ./cmd/web/assets
COPY internal ./internal
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=cache,target=/root/.cache/go-build \
    make build/web

FROM gcr.io/distroless/base-debian13
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
VOLUME ["/data"]
ENV JA=/data/ja/config.json EN=/data/en/config.json
ENTRYPOINT ["/app/web"]
