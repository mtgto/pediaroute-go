# syntax = docker/dockerfile:experimental

# DOCKER_BUILDKIT=1 docker build --file build/package/data/Dockerfile --build-arg DATE=YYYYMMDD --tag mtgto/pediaroute-data:YYYYMMDD .

FROM golang:1.26-trixie AS app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
COPY Makefile go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN --mount=type=cache,target=/root/.cache/go-build make build/gen

FROM golang:1.26-trixie AS ja
WORKDIR /data
ARG DATE
ARG BASE_URL=https://wikimedia.bringyour.com/

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN wget --no-verbose --continue "${BASE_URL}jawiki/${DATE}/jawiki-${DATE}-page.sql.gz"
RUN wget --no-verbose --continue "${BASE_URL}jawiki/${DATE}/jawiki-${DATE}-pagelinks.sql.gz"
RUN wget --no-verbose --continue "${BASE_URL}jawiki/${DATE}/jawiki-${DATE}-linktarget.sql.gz"

RUN mkdir ja
RUN ./gen -lang ja -ip "jawiki-${DATE}-page.sql.gz" -il "jawiki-${DATE}-pagelinks.sql.gz" -ilt "jawiki-${DATE}-linktarget.sql.gz" -o ja


FROM golang:1.26-trixie AS en
WORKDIR /data
ARG DATE
ARG BASE_URL=https://wikimedia.bringyour.com/

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN wget --no-verbose --continue "${BASE_URL}enwiki/${DATE}/enwiki-${DATE}-page.sql.gz"
RUN wget --no-verbose --continue "${BASE_URL}enwiki/${DATE}/enwiki-${DATE}-pagelinks.sql.gz"
RUN wget --no-verbose --continue "${BASE_URL}enwiki/${DATE}/enwiki-${DATE}-linktarget.sql.gz"

RUN mkdir en
RUN ./gen -lang en -ip "enwiki-${DATE}-page.sql.gz" -il "enwiki-${DATE}-pagelinks.sql.gz" -ilt "enwiki-${DATE}-linktarget.sql.gz" -o en


FROM gcr.io/distroless/base-debian13

WORKDIR /data
COPY --from=ja /data/ja /data/ja
COPY --from=en /data/en /data/en
