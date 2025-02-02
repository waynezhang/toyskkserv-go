OUTPUT_PATH=bin
BINARY=toyskkserv

VERSION=`git for-each-ref --sort=creatordate --format '%(refname)' refs/tags | tail -n 1 | sed 's/refs\/tags\/v\(.*\)/\1/g'`
BUILD_TIME=`date +%Y%m%d%H%M`

LDFLAGS=-ldflags "-X github.com/waynezhang/toyskkserv/internal/defs.Version=${VERSION} -X github.com/waynezhang/toyskkserv/internal/defs.Revision=${BUILD_TIME}"

all: build

build:
	@CGO_ENABLED=0 go build ${LDFLAGS} -ldflags "-w -s" -o ${OUTPUT_PATH}/${BINARY} main.go

dev:
	@CGO_ENABLED=0 go build ${LDFLAGS} -o tmp/main main.go

test:
	@go test ./...

coverage:
	@TMPFILE=$$(mktemp); \
		go test ./... -coverprofile=$$TMPFILE; \
		go tool cover -html $$TMPFILE

changelog:
	@TMP_FILE=$$(mktemp); \
	cat CHANGELOG.md > $$TMP_FILE; \
	./scripts/changelog > CHANGELOG.md; \
	echo "\n" >> CHANGELOG.md; \
	cat $$TMP_FILE >> CHANGELOG.md

.PHONY: install
install:
	@go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	@if [ -f ${OUTPUT_PATH}/${BINARY} ] ; then rm ${OUTPUT_PATH}/${BINARY} ; fi
