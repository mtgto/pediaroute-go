GO ?= go
LDFLAGS := -w -s

GEN := build/gen
WEB := build/web

.PHONY: all
all: $(GEN) $(WEB)

.PHONY: test
test: test-web

$(GEN): cmd/gen/main.go $(wildcard internal/app/core/*.go) $(wildcard internal/app/gen/*.go)
	$(GO) build -ldflags "$(LDFLAGS)" -o $@ $<

$(WEB): cmd/web/main.go $(wildcard internal/app/core/*.go) $(wildcard internal/app/web/*.go)
	$(GO) build -ldflags "$(LDFLAGS)" -o $@ $<

test-web:
	$(GO) test github.com/mtgto/pediaroute-go/internal/app/web

.PHONY: image
image:
	docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .
