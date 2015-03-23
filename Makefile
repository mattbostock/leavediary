VERSION := $(shell git describe --always)

export PATH := $(shell pwd)/Godeps/_workspace/bin:$(PATH)
export GOPATH := $(shell pwd)/Godeps/_workspace:$(GOPATH)

build:
	go build -ldflags "-X main.version $(VERSION)"

testdeps:
	which cover || go get golang.org/x/tools/cmd/cover

test:   testdeps
	go test -cover ./...

run:	build
	./timeoff
