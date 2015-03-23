VERSION := $(shell git describe --always)

build:
	go build -ldflags "-X main.version $(VERSION)"

test:
	go test -cover ./...

run:	build
	./timeoff
