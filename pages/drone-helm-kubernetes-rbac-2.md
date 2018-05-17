title: CI/CD with Drone, Kubernetes and Helm - Part 2
description: "Building your CI/CD pipeline with Drone, Kubernetes and Helm. RBAC included."
banner: "/assets/kube-drone-helm/banner.png"
slug: ci-cd-with-drone-kubernetes-and-helm-2
tags: ["ci-cd", "drone", "helm", "kubernetes", "rbac"]
date: "2018-05-17 15:02:43"
draft: true

# Introduction

This is the second part of the article series. In [the first part](/post/ci-cd-with-drone-kubernetes-and-helm-1)
we saw how to start a Kubernetes cluster, how to deploy Tiller and use Helm with
it, and we deployed a Drone instance. 

In this article we'll see how to create a quality pipeline for a go project, how
to build and push a docker image to Google Cloud Registry from our CI according
to the different event that can be handled by Drone.

# Go Project

## TL;DR

You can find the code and the various files [here](https://github.com/Depado/articles/tree/master/code/dummy/).

## Simple Project

We're going to work on a dummy go project. So just create a new repository in
your VCS, clone it and create a new file named `main.go` in it:

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"alive": true})
	})

	if err := r.Run("127.0.0.1:8080"); err != nil {
		logrus.WithError(err).Fatal("Couldn't listen")
	}
}
```

This example is a simple server that listens on `127.0.0.1:8080` and has a
single route: `/health`, a health-check route, which will always respond 
`200 OK` with a small JSON message.

## Tooling and Docker

As mentioned in a [previous article](/post/checklist-for-go-projects), we're
going to use [dep](https://github.com/golang/dep) to handle our dependencies,
and we're going to have a multi-stage Dockerfile. 

So first, let's install `dep` and initialize it in our repository:

```
$ go get -u github.com/golang/dep/cmd/dep
$ dep init
  Locking in  (a5b47d3) for transitive dep gopkg.in/yaml.v2
  Locking in master (22d885f) for transitive dep github.com/gin-contrib/sse
  Locking in  (c88ee25) for transitive dep github.com/ugorji/go
  Locking in v8.18.1 (5f57d22) for transitive dep gopkg.in/go-playground/validator.v8
  Locking in master (1a580b3) for transitive dep golang.org/x/crypto
  Locking in master (7c87d13) for transitive dep golang.org/x/sys
  Locking in  (57fdcb9) for transitive dep github.com/mattn/go-isatty
  Using ^1.2.0 as constraint for direct dep github.com/gin-gonic/gin
  Locking in v1.2 (d459835) for direct dep github.com/gin-gonic/gin
  Locking in  (5a0f697) for transitive dep github.com/golang/protobuf
  Using ^1.0.5 as constraint for direct dep github.com/sirupsen/logrus
  Locking in v1.0.5 (c155da1) for direct dep github.com/sirupsen/logrus
```

We now have two more files in our repository: `Gopkg.lock` and `Gopkg.toml`. We
also have a `vendor/` directory. Great, now our dependencies are fixed. Let's
create the `Dockerfile`:

```dockerfile
# Build step
FROM golang:latest AS build

RUN mkdir -p $GOPATH/src/github.com/Depado/dummy
ADD . $GOPATH/src/github.com/Depado/dummy
WORKDIR $GOPATH/src/github.com/Depado/dummy
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 go build -o /dummy

# Final step
FROM alpine

RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*

COPY --from=build /dummy /home/
ENV TZ=Europe/Paris
WORKDIR /home
ENTRYPOINT ./dummy
EXPOSE 8080
```

We are also going to prevent Docker from copying the whole `vendor/` directory
by creating a `.dockerignore` at the root of our repository. And voilà. We're 
now ready to build our Docker image !

```
$ docker build .
...
Successfully built 9a763d6aa971 
$ docker images
REPOSITORY   TAG      IMAGE ID       CREATED             SIZE
<none>       <none>   9a763d6aa971   1 minutes ago       23.5MB
```

There we go, we now have a small docker image containing our go program, and
everything need to be deployed to a Kubernetes cluster.

# Drone Pipeline

## Basic Pipeline

As stated in the previous article of the series, Drone works the same way as 
Travis, which is: You create a `.drone.yml` file at the root of your repository.

So let's do that:

```yaml
workspace:
  base: /go
  path: src/github.com/Depado/dummy

pipeline:
  prerequisites:
    image: "golang:latest"
    commands: 
      - go version
      - go get -u github.com/golang/dep/cmd/dep
      - dep ensure -vendor-only
  
  build:
    image: "golang:latest"
    commands:
      - go build
```

This is the most basic pipeline you can create while you're using `dep`. The 
first step is, using the `golang:latest` Docker image, display the go version,
install `dep` and then install the dependencies. The second step of the pipeline 
is simply to build and check if our project builds. 

## Pushing to GCR

Now things are getting serious. In this section we are going to start using
Drone secrets. So you need to make sure that you 
[installed the `drone` CLI](http://docs.drone.io/cli-installation/), and that 
you [configured it correctly](http://docs.drone.io/cli-authentication/).

You can check if that worked properly by running:

```
$ drone info
User: you
Email: you@yourmail.com
```

# Helm Chart