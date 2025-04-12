.PHONY: all clean build test

all:

clean:
	go clean
	rm -f cover.out cover.html rsrc_windows_*.syso

build:
	go generate
	env GOOS=linux go build
	env GOOS=windows go build

test:
	go test -v -coverprofile=cover.out
	go tool cover -func=cover.out
	go tool cover -html=./cover.out -o cover.html
