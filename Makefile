all: build

VERSION ?= $(shell git describe --always --tags --long)
PORT ?= 5000
SOURCES = $(shell find . -type f -name '*.go')

build: ${SOURCES}
	go get -d -v .
	go build -ldflags "-X main.Version \"${VERSION}\"" main.go

run: build
	./main

watch: ${GOPATH}/bin/gin
	@gin --appPort ${PORT} --immediate --bin main run

${GOPATH}/bin/gin:
	@echo -e "\n\033[1mError: install 'gin' with 'go get -v github.com/codegangsta/gin' first\033[0m\n"
	@exit 1
