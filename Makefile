GO ?= go
LDFLAGS := -w -s

all: build-gen

build-gen:
	$(GO) build -o build/gen cmd/gen/main.go
