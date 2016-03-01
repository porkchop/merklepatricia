.PHONY: all test get_deps

all: test install

install: get_deps
	go install github.com/porkchop/merklepatricia/cmd/...

test:
	go test -v github.com/porkchop/merklepatricia/...

get_deps:
	go get -d github.com/porkchop/merklepatricia/...
