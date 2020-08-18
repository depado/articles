title: "Small and secure Docker images for Go"
description: |
    How to get your backend in the smallest docker image possible 
slug: small-secure-docker-golang
date: 2020-08-14 20:30:00
tags: [go,dev,docker]


<details>
<summary>TL;DR</summary>

Use [upx](https://github.com/upx/upx) to compress your binaries to the extreme 
and [distroless](https://github.com/GoogleContainerTools/distroless) as your
final Docker image.
</details>

<details class="code">
<summary>TL;DR Dockerfile</summary>

```dockerfile
FROM node:14.5.0-alpine3.12 AS front_builder

ADD ./front /front
WORKDIR /front
RUN npm install && ./node_modules/.bin/quasar build

# Backend Build Step
FROM golang:1.15.0-alpine3.12 AS builder

# Prerequisites
RUN apk update && apk add --no-cache upx

# Dependencies
WORKDIR $GOPATH/src/github.com/Depado/vuemonit
COPY . .
RUN go mod download
RUN go mod verify
RUN go get github.com/rakyll/statik

# Copy frontend build
COPY --from=front_builder /front/dist $GOPATH/src/github.com/Depado/vuemonit/front/dist/

# Build
ARG build
ARG version
RUN mv front/dist/spa/index.html front/dist/spa/main.html
RUN statik -src=./front/dist/spa/ -f
RUN mv front/dist/spa/main.html front/dist/spa/index.html
RUN CGO_ENABLED=0 go build \
      -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" \
      -o /tmp/vuemonit
RUN upx --brute /tmp/vuemonit

# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/vuemonit /go/bin/vuemonit
VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/vuemonit", "--server.port", "8080", "--server.host", "0.0.0.0"]
```
</details>


# Context

I've been working on [vuemonit](https://github.com/Depado/vuemonit) for some
time now. This project aims to create an all-in-one app for simple observability
and monitoring while keeping things simple. Not a revolutionary idea but it's
fun to work on.

![vuemonit](/assets/small-secure-go-docker/vuemonit.png)
*Early version of vuemonit*

This project is composed of two parts:

- A Golang backend in charge of fetching data and storing it in a 
  [bbolt](https://github.com/etcd-io/bbolt) database using various optimization
  techniques:
    - Using [xid](https://github.com/rs/xid) to get sortable IDs containing the 
      date of creation in it, effectively avoiding data duplication.
    - Using [protobuf](https://github.com/golang/protobuf) to serialize the time
      series and keep everything small and efficient.
- A frontend that will consume the backend's API written in [Vue](https://vuejs.org/)
  using [quasar](https://quasar.dev/) framework.

Both these parts can be run independently as you'd expect, but the backend can
also embed the frontend within itself and serve it. Of course that means I'll
face configuration issues later on, but this will be addressed later. For now
I want to concentrate on making everything as small as possible.

For starters, let's say the frontend was optimized with tree shaking and other
size reducing techniques. It's currently 2.3MB. Once said frontend has been
embedded, the go binary is 23MB in size, which is quite big already. 

# Docker

## First Steps

Let's start with a simple multi-step Dockerfile. The first step's role is to 
build the frontend, the second one's role is to embed the frontend and build the
Go binary.

<details>
<summary>Makefile for this Dockerfile</summary>

The following Dockerfile will use `${version}` and `${build}` which are passed
as arguments (using `--build-arg`) to the Docker command. They are computed as
follow:

```makefile
export VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
export BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
BINARY=vuemonit
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD) -s -w"

.PHONY: build
build: ## Build
	go build $(LDFLAGS) -o $(BINARY) 

.PHONY: docker
docker: ## Build the docker image
	docker build -t $(BINARY):latest -t $(BINARY):$(BUILD) \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .
```

The version corresponds to the latest tag and defaults to `0.1.0` when no tag
was found in the repo. The build variable corresponds to the latest commit SHA1
and defaults to `undefined` if no commit was made.

These variables are then injected during compile time into the Go binary using 
`-ldflags="-X main.Version=${version} -X main.Build=${build}"`.

`main.Version` and `main.Build` are variables in our `main.go` file defined like
this:

```go
// Build number and versions injected at compile time, set yours
var (
	Version = "unknown"
	Build   = "unknown"
)

// Version command that will display the build number and version (if any)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run:   func(cmd *cobra.Command, args []string) { 
        fmt.Printf("Build: %s\nVersion: %s\n", Build, Version) 
    },
}
```

Using the above Makefile we can then simply run `make build`:

```sh
$ make build
go build -ldflags "-X main.Version=0.1.0 -X main.Build=1be06ebb91721f3c0ed613f17de4a58c840dc294 -s -w" -o vuemonit
$ ./vuemonit version
Build: 1be06ebb91721f3c0ed613f17de4a58c840dc294
Version: 0.1.0
```
</details>

```dockerfile
FROM node:14.5.0-alpine3.12 AS front_builder

ADD ./front /front
WORKDIR /front
RUN npm install && ./node_modules/.bin/quasar build

# Backend Build Step
FROM golang:1.15.0-alpine3.12 AS builder

# Dependencies
WORKDIR $GOPATH/src/github.com/Depado/vuemonit
COPY . .
RUN go mod download
RUN go mod verify
RUN go get github.com/rakyll/statik

# Copy frontend build
COPY --from=front_builder /front/dist $GOPATH/src/github.com/Depado/vuemonit/front/dist/

# Build
ARG build
ARG version
# Rename the main.html to index.html as to serve it properly with the backend
RUN mv front/dist/spa/index.html front/dist/spa/main.html
RUN statik -src=./front/dist/spa/ -f
RUN mv front/dist/spa/main.html front/dist/spa/index.html
RUN CGO_ENABLED=0 go build \ 
      -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" \ 
      -o /tmp/vuemonit
```

```
REPOSITORY   CREATED              SIZE
vuemonit     About a minute ago   1.06GB
```

As we can see, even when using something like Alpine as our builder, the docker
image is huge. You may notice that it's not a totally valid Dockerfile, there is
no entrypoint or things like that. That's because we're about to add another
step.

## Enter distroless

[distroless](https://github.com/GoogleContainerTools/distroless) is a minimal
and secure Docker image which contains the minimal runtime dependencies. That
means no shell, no package manager, etc. It does contains things like CA
certificates and timezone information though (the things you'd have to copy
from another image if using `scratch` for example).

> "Distroless" images contain only your application and its runtime 
> dependencies. They do not contain package managers, shells or any other 
> programs you would expect to find in a standard Linux distribution.
> 
> <cite>- distroless README</cite>

```dockerfile
FROM node:14.5.0-alpine3.12 AS front_builder

ADD ./front /front
WORKDIR /front
RUN npm install && ./node_modules/.bin/quasar build

# Backend Build Step
FROM golang:1.15.0-alpine3.12 AS builder

# Dependencies
WORKDIR $GOPATH/src/github.com/Depado/vuemonit
COPY . .
RUN go mod download
RUN go mod verify
RUN go get github.com/rakyll/statik

# Copy frontend build
COPY --from=front_builder /front/dist $GOPATH/src/github.com/Depado/vuemonit/front/dist/

# Build
ARG build
ARG version
RUN mv front/dist/spa/index.html front/dist/spa/main.html
RUN statik -src=./front/dist/spa/ -f
RUN mv front/dist/spa/main.html front/dist/spa/index.html
RUN CGO_ENABLED=0 go build \
      -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" \
      -o /tmp/vuemonit

# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/vuemonit /go/bin/vuemonit
VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/vuemonit", "--server.port", "8080", "--server.host", "0.0.0.0"]
```

```
REPOSITORY   CREATED         SIZE
vuemonit     5 seconds ago   18.6MB
```

That's a lot better! But what if we reduced the image size even more?

## Enter UPX

[upx](https://github.com/upx/upx) which stands for 
“The **U**ltimate **P**acker for e**X**ecutables” is a software that will
compress executable files. 

> UPX is an advanced executable file compressor. UPX will typically
> reduce the file size of programs and DLLs by around 50%-70%, thus
> reducing disk space, network load times, download times and
> other distribution and storage costs.
> 
> <cite>[UPX README](https://github.com/upx/upx)</cite>

As stated in the documentation, UPX will typically reduce the program sizes by
50%-70%, which is great. But this operation has a cost: it's time consuming. For
example, to pack my binary it took almost 7 minutes for a gain of 44% in size.
Which means, my binary went from ~23MB to ~9.9MB.

```
     File size         Ratio      Format      Name
--------------------   ------   -----------   -----------
23183975 ->  10334576   44.58%   linux/amd64   vuemonit
```

!!! warning "Tradeoff"
    Packing binaries is a complex task that can take a while to complete and 
    sometimes it isn't worth your time or CPU to do that, especially if the 
    gain is minimal. 

As to what exactly happens when your run a packer, I'll cite the excellent post
[“Packers, How They Work, Featuring UPX”](https://dzone.com/articles/packers-how-they-work-featuring-upx)
by [Christopher Lamb](https://dzone.com/users/2652714/chrislambistan.html):

> […] it compresses and compartmentalizes programs. It will take an executable, 
> compress it, and pack the compressed code into another section of the 
> executable. At runtime, it will uncompress the previously compressed code and
> execute it.
>
> <cite>["Packers, How They Work, Featuring UPX"](https://dzone.com/articles/packers-how-they-work-featuring-upx)</cite>

## Final Dockerfile

```dockerfile
FROM node:14.5.0-alpine3.12 AS front_builder

ADD ./front /front
WORKDIR /front
RUN npm install && ./node_modules/.bin/quasar build

# Backend Build Step
FROM golang:1.15.0-alpine3.12 AS builder

# Prerequisites, installs UPX
RUN apk update && apk add --no-cache upx

# Dependencies
WORKDIR $GOPATH/src/github.com/Depado/vuemonit
COPY . .
RUN go mod download
RUN go mod verify
RUN go get github.com/rakyll/statik

# Copy frontend build
COPY --from=front_builder /front/dist $GOPATH/src/github.com/Depado/vuemonit/front/dist/

# Build
ARG build
ARG version
RUN mv front/dist/spa/index.html front/dist/spa/main.html
RUN statik -src=./front/dist/spa/ -f
RUN mv front/dist/spa/main.html front/dist/spa/index.html
RUN CGO_ENABLED=0 go build \
      -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" \
      -o /tmp/vuemonit
RUN upx --brute /tmp/vuemonit

# Final Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/vuemonit /go/bin/vuemonit
VOLUME [ "/data" ]
WORKDIR /data
EXPOSE 8080
ENTRYPOINT ["/go/bin/vuemonit", "--server.port", "8080", "--server.host", "0.0.0.0"]
```

```
REPOSITORY   CREATED         SIZE
vuemonit     9 minutes ago   7.3MB
```

With the extra steps of adding `upx` to our build step, and executing it, we
reduced our image size by more than two times. 
