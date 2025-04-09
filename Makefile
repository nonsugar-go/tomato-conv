.PHONY: all clean build

all: clean build

clean:
	go clean
	rm -f rsrc_windows_*.syso

build:
	go generate
	env GOOS=linux go build
	env GOOS=windows go build
