title: "Dialogflow Webhook, Golang and Protobuf"
description: "Using protobuf for our Dialogflow webhook. Also using go modules ♥"
banner: "/assets/dialogflow-protobuf/banner.png"
slug: dialogflow-webhook-golang-protobuf
tags: ["go", "dev", "protobuf", "dialogflow", "modules"]
date: "2018-09-11 11:54:00"
draft: true

# Introduction

In the [previous article about Dialogflow](/post/dialogflow-golang-webhook) we 
created a bot using a webhook written in Go. At the time of writing, there was
no Dialogflow SDK for Go, and I had to create a library to bootstrap the 
structures of a DialogFlow webhook call (which you can find 
[here](https://github.com/leboncoin/dialogflow-go-webhook)). Now with the `v2`
version of the Dialogflow API, there's a 
[Go SDK](https://github.com/GoogleCloudPlatform/google-cloud-go/tree/master/dialogflow/apiv2)
which allows you to control Dialogflow's behavior. The `v2` is now enabled by
default.

![banner](/assets/dialogflow-protobuf/warning.png)

But that's not the main point of interest of this release. When they released
the SDK, they also released the 
[Go code that was generated from the protocol buffer definition](https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2), 
and that's what we're going to use to properly handle an incoming webhook 
request.

## Dialogflow Webhook

Dialogflow is great to create conversational bots, but sometimes you want to do
more than just answer with predefined text. You might want to send an email, 
query an API and format its response to answer your user, or even store
information in a database. That's the goal of 
[Fulfillment](https://dialogflow.com/docs/fulfillment) in Dialogflow. When you
enable the fulfillment on a specific intent, you tell Dialogflow to send a query
with a payload to your backend, as shown in the 
[previous article](/post/dialogflow-golang-webhook#toc_4) and the schema below:

![webhook](/assets/dialogflow/global-flow.svg)

## Protocol Buffer (aka. Protobuf)

> Protocol buffers are Google's language-neutral, platform-neutral, extensible 
> mechanism for serializing structured data – think XML, but smaller, faster, 
> and simpler. You define how you want your data to be structured once, then you
> can use special generated source code to easily write and read your structured
> data to and from a variety of data streams and using a variety of languages.
> 
> <cite>[Official Protobuf Website](https://developers.google.com/protocol-buffers/)</cite>

Basically protobuf allows you to define how your data is structured once. 
Then you can generate the source code that will allow you to easily read and 
write in that format with the language you like the most or need to use.

# Handling Webhook

## Setup and initial code

We're going to start with a simple program that uses 
[gin](https://github.com/gin-gonic/gin) as the HTTP router, and 
[logrus](https://github.com/sirupsen/logrus) as the logger. So write that down
in a `main.go` file.

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func handleWebhook(c *gin.Context) {

}

func main() {
	var err error

	r := gin.Default()
	r.POST("/webhook", handleWebhook)

	if err = r.Run("127.0.0.1:8080"); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
```

This program only starts the router and listens on the `8080` port, or exits
with a proper error message if it fails to do so. It also accepts posts on the
`/webhook` endpoint but doesn't do anything with it.

Two options here. We'll first see how to use the new 
[go module](https://github.com/golang/go/wiki/Modules#go-111-modules) approach 
if you're already using [go1.11](https://blog.golang.org/go1.11). If you're not 
using modules or go1.11 yet, we'll also see how to use 
[dep](https://github.com/golang/dep) to manage our dependencies.


### Using go modules

The behavior can change if you're inside your `GOPATH` or outside. For this 
example we'll just create a new directory anywhere outside our `GOPATH` and put
the `main.go` file we created in the previous section inside it.

```
$ mkdir ~/dialogflowpb
$ cd ~/dialogflowpb
$ # Create that main.go file or copy it from elsewhere
$ go mod init github.com/Depado/articles/code/dialogflowpb
go: creating new go.mod: module github.com/Depado/articles/code/dialogflowpb
$ cat go.mod 
module github.com/Depado/articles/code/dialogflowpb
```

Here we initialized our package and gave it a name. It created a `mod.go` file
which just contains the package name and nothing else. Things get interesting
when we run the `go build` command though:

```
$ go build
go: finding github.com/gin-contrib/sse latest
go: finding github.com/ugorji/go/codec latest
go: finding github.com/golang/protobuf/proto latest
go: finding golang.org/x/sys/unix latest
go: finding golang.org/x/sys latest
go: finding golang.org/x/crypto/ssh/terminal latest
go: finding golang.org/x/crypto/ssh latest
go: finding golang.org/x/crypto latest

$ cat go.mod
module github.com/Depado/articles/code/dialogflowpb

require (
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/golang/protobuf v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/sirupsen/logrus v1.0.6
	github.com/ugorji/go/codec v0.0.0-20180831062425-e253f1f20942 // indirect
	golang.org/x/crypto v0.0.0-20180910181607-0e37d006457b // indirect
	golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/yaml.v2 v2.2.1 // indirect
)
```

And voilà. Not only did we installed the proper versions of the two libraries
we're using (namely gin and logrus) but also their transient dependencies.

<object type="image/svg+xml" data="/assets/dialogflow-protobuf/term.svg"></object>

### Using `dep`

We'll first init dep by running `dep init`.

```
$ dep init
...
(1/11) Wrote github.com/mattn/go-isatty@57fdcb988a5c543893cc61bce354a6e24ab70022
(2/11) Wrote gopkg.in/yaml.v2@a5b47d31c556af34a302ce5d659e6fea44d90de0
(3/11) Wrote github.com/gin-contrib/sse@master
(4/11) Wrote gopkg.in/go-playground/validator.v8@v8.18.1
(5/11) Wrote github.com/ugorji/go@c88ee250d0221a57af388746f5cf03768c21d6e2
(6/11) Wrote github.com/sirupsen/logrus@v1.0.6
(7/11) Wrote github.com/gin-gonic/gin@v1.3.0
(8/11) Wrote github.com/golang/protobuf@5a0f697c9ed9d68fef0116532c6e05cfeae00e55
(9/11) Wrote github.com/json-iterator/go@1.0.0
(10/11) Wrote golang.org/x/sys@master
(11/11) Wrote golang.org/x/crypto@master
```

Then we can add the following package:

```
$ dep ensure -add google.golang.org/genproto/googleapis/cloud/dialogflow/v2
```

## Handler

Now we're going to import the generated code from protobuf:

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)
```

We now have access to all the structures and especially the 
[WebhookRequest](https://godoc.org/google.golang.org/genproto/googleapis/cloud/dialogflow/v2#WebhookRequest)
one, which will allow us to properly unmarshal a webhook request. But now we 
have a problem. We have the Go code that corresponds to the protocol buffer, but 
Dialogflow will send JSON requests to our endpoint. Luckily, the protobuf team 
thought of that for us and created a package named 
[jsonpb](https://github.com/golang/protobuf/tree/master/jsonpb).

> Package jsonpb provides marshaling and unmarshalling between protocol buffers 
> and JSON. 
> It follows the specification at 
> https://developers.google.com/protocol-buffers/docs/proto3#json.
> This package produces a different output than the standard "encoding/json" 
> package, which does not operate correctly on protocol buffers.
>
> <cite>[Source](https://github.com/golang/protobuf/blob/master/jsonpb/jsonpb.go)</cite>

We can now modify our handler:

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func handleWebhook(c *gin.Context) {
	var err error

	wr := dialogflow.WebhookRequest{}
	if err = jsonpb.Unmarshal(c.Request.Body, &wr); err != nil {
		logrus.WithError(err).Error("Couldn't Unmarshal request to jsonpb")
		c.Status(http.StatusBadRequest)
		return
	}
}
```

If you used `dep` to handle your dependencies, run another `dep ensure`, if you
used the go modules approach, nothing to do but running `go build`.

We can now use the `wr` struct to retrieve the data Dialogflow sent us. For
example if you want to retrieve the output contexts:

```go 
contexts := wr.QueryResult.OutputContexts
// or
contexts = wr.GetQueryResult().GetOutputContexts()
```

# Conclusion

You can find the [complete code here](https://github.com/Depado/articles/tree/master/code/dialogflowpb).

Thanks to [@ashleymcnamara](https://github.com/ashleymcnamara) for the amazing
[Gopher Artworks](https://github.com/ashleymcnamara/gophers) !

