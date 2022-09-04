# syntax = docker/dockerfile:experimental

# DOCKER_BUILDKIT=1 docker build --file build/package/data/Dockerfile --build-arg DATE=YYYYMMDD --tag mtgto/pediaroute-data:YYYYMMDD .

FROM golang:1.19-bullseye as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
COPY Makefile go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN --mount=type=cache,target=/root/.cache/go-build make build/gen

FROM golang:1.19-bullseye as ja
WORKDIR /data
ARG DATE
ARG BASE_URL=https://wikimedia.bringyour.com/

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN wget --no-verbose --continue "${BASE_URL}/jawiki/${DATE}/jawiki-${DATE}-page.sql.gz"
RUN wget --no-verbose --continue "${BASE_URL}/jawiki/${DATE}/jawiki-${DATE}-pagelinks.sql.gz"

RUN mkdir ja
RUN ./gen -lang ja -ip "jawiki-${DATE}-page.sql.gz" -il "jawiki-${DATE}-pagelinks.sql.gz" -o ja


FROM golang:1.19-bullseye as en
WORKDIR /data
ARG DATE
ARG BASE_URL=https://wikimedia.bringyour.com/

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN wget --no-verbose --continue "${BASE_URL}/enwiki/${DATE}/enwiki-${DATE}-page.sql.gz"
RUN wget --no-verbose --continue "${BASE_URL}/enwiki/${DATE}/enwiki-${DATE}-pagelinks.sql.gz"

RUN mkdir en
RUN ./gen -lang en -ip "enwiki-${DATE}-page.sql.gz" -il "enwiki-${DATE}-pagelinks.sql.gz" -o en


FROM gcr.io/distroless/base-debian10

WORKDIR /data
COPY --from=ja /data/ja /data/ja
COPY --from=en /data/en /data/en
