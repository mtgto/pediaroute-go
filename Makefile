GO ?= go
LDFLAGS := -w -s

GEN := build/gen
WEB := build/web
BENCH := build/bench

.PHONY: all
all: $(GEN) $(WEB) $(BENCH)

.PHONY: test
test: test-web

$(GEN): cmd/gen/main.go $(wildcard internal/app/core/*.go) $(wildcard internal/app/gen/*.go)
	CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o $@ $<

$(WEB): cmd/web/main.go $(wildcard internal/app/core/*.go) $(wildcard internal/app/web/*.go)
	CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o $@ $<

$(BENCH): cmd/bench/main.go $(wildcard internal/app/core/*.go)
	CGO_ENABLED=0 $(GO) build -ldflags "$(LDFLAGS)" -o $@ $<

.PHONY: test-web
test-web:
	$(GO) test -v github.com/mtgto/pediaroute-go/internal/app/web

.PHONY: image
image:
	docker build --file build/package/app/Dockerfile --tag mtgto/pediaroute:latest .
