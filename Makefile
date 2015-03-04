all: build

VERSION ?= $(shell git describe --always --tags --long)

SOURCES = $(shell find . -type f -name '*.go')

build: ${SOURCES}
	go build -ldflags "-X main.Version \"${VERSION}\"" main.go

release: ${SOURCES}
	go build -ldflags "-X main.Environment \"production\" -X main.Version \"${VERSION}\"" main.go
