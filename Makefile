HAS_GOLANGCI	:= $(shell command -v golangci-lint)
VERSION      	?= $(shell git describe --tags 2> /dev/null || \
            		git describe --match=$(git rev-parse --short=8 HEAD) --always --dirty --abbrev=8)
LDFLAGS			:= "-w -s -X 'main.version=${VERSION}'"
GOOS			?=linux
GOARCH			?=amd64
CGO_ENABLED     ?=0

golangci:
ifndef HAS_GOLANGCI
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.9.0
endif
	golangci-lint run --timeout 10m0s

test:
	go test -v -race -count 1 ./...

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build \
		 -ldflags $(LDFLAGS) \
		 -o build/k8spodsmetrics-$(GOOS)-$(GOARCH) \
		 ./cmd/k8spodsmetrics

lint: golangci
	go vet ./...

format:
	go fmt ./...

.PHONY: test list build format
