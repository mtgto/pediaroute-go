# pediaroute

Source of https://pediaroute.com/ .

# How to build

You need

- Golang 1.16+
- Make
- Node.js 14+
- yarnpkg

```console
$ go get

$ cd web
$ yarn install && yarn build && yarn asset
```

# How to generate application data

Before generating application data, you need to download wikipedia sql files.

Japanese

1. jawiki-YYYYMMDD-page.sql.gz
1. jawiki-YYYYMMDD-pagelinks.sql.gz

English

1. enwiki-YYYYMMDD-page.sql.gz
1. enwiki-YYYYMMDD-pagelinks.sql.gz

```console
$ ./build/gen \
-ip /path/to/jawiki-YYYYMMDD-page.sql.gz \
-il /path/to/jawiki-YYYYMMDD-pagelinks.sql.gz \
-o .
```

You will find 4 files in current directory:

- config.json
- title.dat
- page.dat
- link.dat

# License

GPL v3
