.PHONY: all clean deps docker test watch

VERSION ?= $(shell git describe --always --tags)
NAME := gol-${VERSION}
PORT ?= 5000
SOURCES := $(shell find . -type f -name '*.go')
SOURCE_DIRS := $(shell find . -type f -name '*.go' | xargs dirname | sort | uniq)
CONTAINER_NAME := 'gol-docker'
GOPATH := $(PWD)/.go
PREFIX := '/usr/share'

all: gol

deps:
	go get -v -d ./...
	/bin/bash -c "go list -f '{{ join .Imports \"\n\" }}' ./... | grep -v '^_' | sort | uniq | xargs go get -v"

gol: ${SOURCES} assets/main.css
	go get -d -v .
	go build -o $@ -ldflags "-X gol.Version=\"${VERSION}\"" gol.go

assets/main.css: assets/main.scss
	bin/sassc -m assets/main.scss assets/main.css

docker: gol
	docker build -t ${CONTAINER_NAME} .

test:
	go test -v ${SOURCE_DIRS}

release: gol test
	mkdir ${NAME}
	cp -R assets ${NAME}/assets
	cp -R templates ${NAME}/templates
	cp gol ${NAME}
	tar -caf ${NAME}-linux.tar.gz ${NAME}
	rm -rf ${NAME}

watch: ${GOPATH}/bin/gin
	@${GOPATH}/bin/gin --appPort ${PORT} --immediate --bin main run

${GOPATH}/bin/gin:
	@echo -e "\n\033[1mError: install 'gin' with 'go get -v github.com/heyLu/gin' first\033[0m\n"
	@exit 1

install: gol
	install -Dm 755 gol 		/usr/bin/gol
	install -dm 755 templates 	$(PREFIX)/gol/templates
	cp -r 			templates 	$(PREFIX)/gol/templates
	install -dm 755 assets 		$(PREFIX)/gol/assets
	cp -r			assets		$(PREFIX)/gol/assets

clean:
	rm -f gol
