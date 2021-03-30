
# docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .

FROM node:14-buster-slim as asset
WORKDIR /node
COPY web/package.json web/yarn.lock .
RUN yarn install
COPY web .
RUN yarn build

FROM golang:1.16-buster as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
COPY Makefile go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/rakyll/statik
COPY --from=asset /node/dist ./dist
COPY cmd ./cmd
COPY internal ./internal
RUN statik -p web -src ./dist -dest internal/app
RUN make build/web

FROM gcr.io/distroless/base-debian10
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
COPY --from=mtgto/pediaroute-data:latest /data/ja /data/ja
COPY --from=mtgto/pediaroute-data:latest /data/en /data/en
ENTRYPOINT ["/app/web"]
