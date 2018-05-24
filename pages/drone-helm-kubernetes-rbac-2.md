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

## Simple Project

**TL;DR:** You can find the code and the various files [here](https://github.com/Depado/articles/tree/master/code/dummy/).

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

[Google Container Registry](https://cloud.google.com/container-registry/) is a
private Docker registry which we're going to use to host our tagged Docker
images for our future deployments. As it's private, we're going to need to
authenticate and we can achieve this by creating what's called a service
account. But this one won't be inside our k8s cluster, but will grant some
authorization.

So let's head to the [IAM Console](https://console.cloud.google.com/iam-admin/serviceaccounts)
and create a new service account. Name it as you like and select the 
"Storage Admin" role, check the box to generate a new JSON key and hit "Save".
You'll be offered to download the JSON key, download and save it on your
computer. 

![new-sa](/assets/kube-drone-helm/new-sa.png)

We're going to use our first Drone plugin, namely the 
[Google Container Registry](http://plugins.drone.io/drone-plugins/drone-gcr/).
There's some explanation on this page on how to use this plugin, and it's
written that this plugin is actually an extension of the 
[Docker](http://plugins.drone.io/drone-plugins/drone-docker/) plugin.

So let's add this step to our `.drone.yml` file:

```yaml
  gcr:
    image: plugins/gcr
    repo: project-id/dummy
    tags: latest
    secrets: [google_credentials]
    when:
      event: push
      branch: master
```

Obviously you need to replace the `project-id` with your own project ID. Here's
what Drone understands at this step:

> If a commit is pushed on the master branch, then build the Docker image, tag 
> it with `latest` and push it to GCR using the credentials defined in the 
> `google_credentials` secret.

We're missing something here. The `google_credentials` secret is yet to be
created. So in your terminal, we'll use the `drone` CLI to create a new
secret for our repository, and we'll limit this secret to be used only by
the `plugins/gcr` image:

```
$ drone secret add --image plugins/gcr --repository Depado/dummy \
  --name google_credentials --value @your_key.json
```

# Helm Chart

## Creating the Chart

In the previous article we learned how to use an Helm Chart. In this section
we'll see how to create a basic chart that will simply create a deployment.

We won't bother about ingress, configmap and such because our goal here is
simply to use Helm from within our CI environment.

```
$ helm create dummy
Creating dummy
```

This will create a new directory `dummy` where you are. This directory
will contain two directories and some files:

- `charts/` A directory containing any charts upon which this chart depends
- `templates/` A directory of templates that, when combined with values, will generate 
  valid Kubernetes manifest files
- `Charts.yaml` A YAML file containing information about the chart
- `values.yaml` The default configuration values for this chart

For more information, check [the documentation](https://github.com/kubernetes/helm/blob/master/docs/charts.md#the-chart-file-structure)
about the chart file structure.

Here, we're going to modify both the `values.yaml` files to use sane defaults
for our chart, and more importantly `templates/` to add and modify the rendered
k8s manifests.