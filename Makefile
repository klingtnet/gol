all: build

VERSION ?= "0.1.0"

SOURCES = $(shell find . -type f -name '*.go')

build: ${SOURCES}
	go build main.go

release: ${SOURCES}
	go build -ldflags "-X main.Environment \"production\" -X main.Version \"${VERSION}\"" main.go
