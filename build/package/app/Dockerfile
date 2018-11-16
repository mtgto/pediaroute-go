FROM node:8-alpine as asset
WORKDIR /node
COPY web .
RUN yarn install && yarn build

FROM golang:1.11-alpine3.8 as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
ENV GO111MODULE=on
COPY --from=asset /node/dist ./dist
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile go.mod go.sum ./
RUN apk add --update git make gcc musl-dev && \
go get github.com/rakyll/statik && \
statik -p web -src ./dist -dest internal/app && \
GOOS=linux GOARCH=amd64 make all

FROM alpine:3.8
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
ENTRYPOINT /app/web
