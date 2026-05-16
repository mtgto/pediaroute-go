# pediaroute

Source of https://pediaroute.com/ .

# How to build

You need

- Golang 1.26
- Make
- Node.js 24+
- pnpm

```console
$ cd web && pnpm install && pnpm build:dev
$ make
```

# How to generate application data

Before generating application data, you need to download wikipedia sql files.

Japanese

1. jawiki-YYYYMMDD-page.sql.gz
2. jawiki-YYYYMMDD-pagelinks.sql.gz
3. jawiki-YYYYMMDD-linktarget.sql.gz

English

1. enwiki-YYYYMMDD-page.sql.gz
2. enwiki-YYYYMMDD-pagelinks.sql.gz
3. enwiki-YYYYMMDD-linktarget.sql.gz

```console
$ ./build/gen \
-ip /path/to/jawiki-YYYYMMDD-page.sql.gz \
-il /path/to/jawiki-YYYYMMDD-pagelinks.sql.gz \
-ilt /path/to/jawiki-YYYYMMDD-linktarget.sql.gz \
-o .
```

You will find 4 files in current directory:

- config.json
- title.dat
- page.dat
- link.dat

For English:

```console
$ mkdir -p en
$ ./build/gen \
-lang en \
-ip /path/to/enwiki-YYYYMMDD-page.sql.gz \
-il /path/to/enwiki-YYYYMMDD-pagelinks.sql.gz \
-ilt /path/to/enwiki-YYYYMMDD-linktarget.sql.gz \
-o en/
```

For Japanese:

```console
$ mkdir -p ja
$ ./build/gen \
-lang ja \
-ip /path/to/jawiki-YYYYMMDD-page.sql.gz \
-il /path/to/jawiki-YYYYMMDD-pagelinks.sql.gz \
-ilt /path/to/jawiki-YYYYMMDD-linktarget.sql.gz \
-o ja/
```

# How to run with Docker

Build the app image:

```console
$ make image
```

Run with the generated `en/` and `ja/` directories mounted:

```console
$ docker run \
-v $(pwd)/ja:/data/ja \
-v $(pwd)/en:/data/en \
-e JA=/data/ja/config.json \
-e EN=/data/en/config.json \
-p 8080:8080 \
mtgto/pediaroute:latest
```

The app will be available at http://localhost:8080 .

# License

GPL v3
