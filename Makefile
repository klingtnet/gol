all: build

VERSION ?= $(shell git describe --always --tags)
NAME = gol-${VERSION}
PORT ?= 5000
SOURCES = $(shell find . -type f -name '*.go')
SOURCE_DIRS = $(shell find . -type f -name '*.go' | xargs dirname | sort | uniq)
CONTAINER_NAME = 'gol-docker'

build: ${SOURCES} assets/main.css
	go get -d -v .
	go build -ldflags "-X main.Version \"${VERSION}\"" main.go

assets/main.css: assets/main.scss
	sassc -m assets/main.scss assets/main.css

docker: build
	docker build -t ${CONTAINER_NAME} .

test:
	go test -v ${SOURCE_DIRS}

release: build test
	mkdir ${NAME}
	cp -R assets ${NAME}/assets
	cp -R templates ${NAME}/templates
	cp main ${NAME}
	tar -caf ${NAME}-linux.tar.gz ${NAME}
	rm -rf ${NAME}

run: build
	./main

watch: ${GOPATH}/bin/gin
	@${GOPATH}/bin/gin --appPort ${PORT} --immediate --bin main run

${GOPATH}/bin/gin:
	@echo -e "\n\033[1mError: install 'gin' with 'go get -v github.com/heyLu/gin' first\033[0m\n"
	@exit 1
