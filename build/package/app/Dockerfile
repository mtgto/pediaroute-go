
# docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .

FROM node:16-bullseye as asset
WORKDIR /node
COPY web/package.json web/yarn.lock /node/
RUN yarn install
COPY web .
RUN yarn build

FROM golang:1.19-bullseye as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
COPY Makefile go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY --from=asset /node/dist ./cmd/web/assets
COPY internal ./internal
RUN make build/web

FROM mtgto/pediaroute-data:latest
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
ENV JA=/data/ja/config.json EN=/data/en/config.json
ENTRYPOINT ["/app/web"]
