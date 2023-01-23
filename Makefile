NAME    ?= health
VERSION ?= dev
COMMIT   = $(shell git rev-parse --short HEAD)
DATE     = $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

build:
	go build -o bin/$(NAME) -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" \
		cmd/health/*.go

test:
	go test -v ./pkg/...
