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

# License

GPL v3
