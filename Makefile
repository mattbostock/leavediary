VERSION := $(shell git describe --always | tr -d '\n'; test -z "`git status --porcelain`" || echo '-dirty')

export PATH := $(shell pwd)/Godeps/_workspace/bin:$(PATH)
export GOPATH := $(shell pwd)/Godeps/_workspace:$(GOPATH)

build:
	go build -ldflags "-X main.version $(VERSION)"

testdeps:
	which cover || go get golang.org/x/tools/cmd/cover

test:   testdeps
	go test -cover ./...

run:	build
	./leavediary
