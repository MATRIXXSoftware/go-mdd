GO = $(shell which go 2>/dev/null)

ifeq ($(GO),)
$(warning "go is not in your system PATH")
endif

ROOT_MODULE=github.com/matrixxsoftware/go-mdd

.PHONY: all clean build test test-race

all: clean build

build:
	$(GO) build ./...
clean:
	$(GO) clean -testcache
test:
	$(GO) test ./... -v
test-race:
	$(GO) test ./... -v -race
test-cover:
	$(GO) test ./... -v -cover -coverprofile=coverage.out 
	$(GO) tool cover -html=coverage.out

test-bench:
	$(GO) test $(ROOT_MODULE)/cmdc -v -bench=. -benchtime=1s
test-profile:
	$(GO) test $(ROOT_MODULE)/cmdc -v -bench=. -benchmem -cpuprofile cpu.prof
	$(GO) tool pprof cpu.prof
test-memprofile:
	$(GO) test $(ROOT_MODULE)/cmdc -v -bench=. -benchmem -memprofile mem.prof
	$(GO) tool pprof mem.prof
