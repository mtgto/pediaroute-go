
# docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .

FROM node:24-trixie AS asset
WORKDIR /node
COPY web/package.json web/pnpm-lock.yaml /node/
RUN corepack enable && pnpm install --frozen-lockfile
COPY web .
RUN pnpm build

FROM golang:1.26-trixie AS app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
COPY Makefile go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY --from=asset /node/dist ./cmd/web/assets
COPY internal ./internal
RUN make build/web

FROM gcr.io/distroless/base-debian13
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
VOLUME ["/data"]
ENV JA=/data/ja/config.json EN=/data/en/config.json
ENTRYPOINT ["/app/web"]
