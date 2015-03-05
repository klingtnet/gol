all: build

VERSION ?= $(shell git describe --always --tags --long)

SOURCES = $(shell find . -type f -name '*.go')

build: ${SOURCES}
	go build -ldflags "-X main.Version \"${VERSION}\"" main.go

release: ${SOURCES}
	go build -ldflags "-X main.Environment \"production\" -X main.Version \"${VERSION}\"" main.go

run: build
	./main

watch: ${GOPATH}/bin/gin
	@gin --appPort 5000 --immediate --bin main run

${GOPATH}/bin/gin:
	@echo -e "\n\033[1mError: install 'gin' with 'go get -v github.com/codegangsta/gin' first\033[0m\n"
	@exit 1
