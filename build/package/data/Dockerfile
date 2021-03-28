# syntax = docker/dockerfile:experimental

# DOCKER_BUILDKIT=1 docker build --file build/package/data/Dockerfile --build-arg DATE=YYYYMMDD --tag mtgto/pediaroute-data:YYYYMMDD .

FROM golang:1.16-alpine as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile go.mod go.sum ./
RUN apk add --update-cache git make gcc musl-dev
RUN --mount=type=cache,target=/root/.cache/go-build make build/gen


FROM alpine as ja
WORKDIR /data
ARG DATE

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN apk --update add wget

RUN wget --no-verbose --continue "https://dumps.wikimedia.org/jawiki/${DATE}/jawiki-${DATE}-page.sql.gz"
RUN wget --no-verbose --continue "https://dumps.wikimedia.org/jawiki/${DATE}/jawiki-${DATE}-pagelinks.sql.gz"

RUN mkdir ja
RUN ./gen -lang ja -ip "jawiki-${DATE}-page.sql.gz" -il "jawiki-${DATE}-pagelinks.sql.gz" -o ja


FROM alpine as en
WORKDIR /data
ARG DATE

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN apk --update-cache add wget

RUN wget --no-verbose --continue "https://dumps.wikimedia.org/enwiki/${DATE}/enwiki-${DATE}-page.sql.gz"
RUN wget --no-verbose --continue "https://dumps.wikimedia.org/enwiki/${DATE}/enwiki-${DATE}-pagelinks.sql.gz"

RUN mkdir en
RUN ./gen -lang en -ip "enwiki-${DATE}-page.sql.gz" -il "enwiki-${DATE}-pagelinks.sql.gz" -o en


FROM tianon/true

WORKDIR /data
COPY --from=ja /data/ja /data/ja
COPY --from=en /data/en /data/en
