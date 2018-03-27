.PHONY: all clean

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

BINARY=myprogram
VERSION=0.1.0
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all:
	go build -o $(BINARY) $(LDFLAGS)

docker:
	docker build -t "your-docker-repo/$(BINARY):$(VERSION)" \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .

clean:
	-rm $(BINARY)