# docker build --file build/package/data/Dockerfile --build-arg DATE=YYYYMMDD --tag mtgto/pediaroute-data:YYYYMMDD .

FROM golang:1.11-alpine3.8 as app
WORKDIR /go/src/github.com/mtgto/pediaroute-go
ENV GO111MODULE=on
COPY cmd ./cmd
COPY internal ./internal
COPY Makefile go.mod go.sum ./
RUN apk add --update git make gcc musl-dev
RUN make build/gen


FROM alpine:3.8 as data
WORKDIR /data
ARG DATE

COPY --from=app /go/src/github.com/mtgto/pediaroute-go/build/gen .

RUN apk --update add wget

RUN wget --no-verbose "https://dumps.wikimedia.org/jawiki/${DATE}/jawiki-${DATE}-page.sql.gz"
RUN wget --no-verbose "https://dumps.wikimedia.org/jawiki/${DATE}/jawiki-${DATE}-pagelinks.sql.gz"
RUN wget --no-verbose "https://dumps.wikimedia.org/enwiki/${DATE}/enwiki-${DATE}-page.sql.gz"
RUN wget --no-verbose "https://dumps.wikimedia.org/enwiki/${DATE}/enwiki-${DATE}-pagelinks.sql.gz"

RUN mkdir ja
RUN ./gen -lang ja -ip "jawiki-${DATE}-page.sql.gz" -il "jawiki-${DATE}-pagelinks.sql.gz" -o ja
RUN mkdir en
RUN ./gen -lang en -ip "enwiki-${DATE}-page.sql.gz" -il "enwiki-${DATE}-pagelinks.sql.gz" -o en


FROM tianon/true

WORKDIR /data
COPY --from=data /data/ja /data/en /data/
