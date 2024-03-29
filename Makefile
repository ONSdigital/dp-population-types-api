BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -json -m all | nancy sleuth

.PHONY: lint
lint:

.PHONY: build
build:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-population-types-api

.PHONY: debug
debug:
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-population-types-api
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-population-types-api

.PHONY: run
debug-run:
	HUMAN_LOG=1 DEBUG=1 go run -race -tags 'debug' $(LDFLAGS) main.go

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:
	go test -cover -coverpkg=github.com/ONSdigital/dp-population-types-api/... -component -logging

.PHONY: test-feature
test-feature:
	go test -cover -coverpkg=github.com/ONSdigital/dp-population-types-api/... -component

.PHONY: run
run:
	HUMAN_LOG=1 go run -tags 'production' -ldflags "-X $(SERVICE_PATH).BuildTime=$(BUILD_TIME) -X $(SERVICE_PATH).GitCommit=$(GIT_COMMIT) -X $(SERVICE_PATH).Version=$(VERSION)" -race $(LDFLAGS) main.go
