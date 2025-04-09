.PHONY: all clean build

all: clean build

clean:
	go clean

build:
	go generate
	env GOOS=linux go build
	env GOOS=windows go build
