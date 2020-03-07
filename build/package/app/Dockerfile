FROM node:12-alpine as asset
WORKDIR /node
COPY web .
RUN yarn install && yarn build

FROM golang:1.14-alpine as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
ENV GO111MODULE=on
COPY --from=asset /node/dist ./dist
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile go.mod go.sum ./
RUN apk add --update-cache git make gcc musl-dev
RUN go get -u github.com/rakyll/statik
RUN statik -p web -src ./dist -dest internal/app
RUN make build/web

FROM alpine
WORKDIR /app
COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/web .
ENTRYPOINT /app/web
