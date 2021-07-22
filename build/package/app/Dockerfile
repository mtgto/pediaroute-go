
# docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .

FROM node:14-buster as asset
WORKDIR /node
COPY web/package.json web/yarn.lock /node/
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

FROM mtgto/pediaroute-data:latest
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
ENV JA=/data/ja/config.json EN=/data/en/config.json
ENTRYPOINT ["/app/web"]
