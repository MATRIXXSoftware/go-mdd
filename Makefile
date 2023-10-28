GO = $(shell which go 2>/dev/null)

APP 			:= go-mdd
APP_VERSION 	:= v0.2.0

ifeq ($(GO),)
$(warning "go is not in your system PATH")
endif

.PHONY: all clean build test test-race

all: clean build

build:
	$(GO) build ./...
clean:
	$(GO) clean -testcache
test:
	$(GO) test ./... -v
test-benchmark:
	$(GO) test ./... -v -bench .
test-race:
	$(GO) test ./... -v -race
test-cover:
	$(GO) test ./... -v -cover -coverprofile=coverage.out 
	$(GO) tool cover -html=coverage.out