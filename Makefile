VERSION := $(shell git describe --always)

build:
	go build -ldflags "-X main.version $(VERSION)"

testdeps:
	which cover || go get golang.org/x/tools/cmd/cover

test:   testdeps
	go test -cover ./...

run:	build
	./timeoff
