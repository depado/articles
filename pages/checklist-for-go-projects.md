title: Checklist for Go projects
description: "I always forget stuff."
banner: "/assets/go-checklist/banner.png"
author: ""
slug: checklist-for-go-projects
tags: 
- go
- dev
- foss
date: "2018-03-26 16:54:46"
draft: false

# Introduction

I start new Go projects regularly, either as part of my daily work or for fun
on my free time. And every time I forget something ! Every time I forget about 
the configurable logger, the Makefile which injects the build number, the 
Dockerfile or a way to configure the project using environment variables.

And that's bad. Because I often realize I forgot something right when I need it.

So here is an article. It will be quite short I think, but will list everything
you should do when starting a Go project.

# TL;DR

- [Choose a license !](https://choosealicense.com/)
- Use [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for conf and command management
- Add a way to configure the logger
- Inject build number and version in your compiled binaries
- Have a minimalist Dockerfile

# Checklist

## License

First of all, when you start a FOSS project, **choose a license** ! Don't leave
your repository with no `LICENSE`, otherwise, by default, your repository is
considered copyrighted. Which means : Literally nobody can use what you did.

As said on the excellent [choosealicense.com](https://choosealicense.com/) on 
[this page](https://choosealicense.com/no-permission/) :

> When you make a creative work (which includes code), the work is under 
> exclusive copyright by default. Unless you include a license that specifies 
> otherwise, nobody else can use, copy, distribute, or modify your work without
> being at risk of take-downs, shake-downs, or litigation. Once the work has 
> other contributors (each a copyright holder), “nobody” starts including you.

So, spread the love of FOSS, and [choose a license](https://choosealicense.com/) !

And once you're starting to have dependencies to your project, remember to check
their dependencies too ! Some tools do that nicely, such as 
[glice](https://github.com/ribice/glice) which will list all your dependencies
and their associated licenses.

## Flexible Configuration

![viper](/assets/go-checklist/viper.png)

I've been searching for Go's best practices about configuration management for a
while now. First I started by writing structs with yaml tags and loading a file.
I also had some flags defined in my `main.go` file. That did the trick.

Then I started working with [Kubernetes](https://kubernetes.io/) and although
it is simple to have a configuration file mounted in a volume, it's easier to
work with environment variables. 

So I wrote [conftags](https://github.com/Depado/conftags), which works fine
but I knew there was a better way to deal with that. So I had a look at 
[viper](https://github.com/spf13/viper).

```go
package conf

import (
	"github.com/spf13/pflag"
    "github.com/spf13/viper"
    "github.com/sirupsen/logrus"
)

func init() {
    // Flags
    pflag.String("server.host", "127.0.0.1", "host on which the server should listen")
    pflag.Int("server.port", 8080, "port on which the server should listen")
    pflag.Bool("server.debug", false, "debug mode for the server")
    
    if err := viper.BindPFlags(pflag.CommandLine); err != nil {
        logrus.WithError(err).Fatal("Couldn't bind flags")
    }
    
    // Environment variables
    viper.AutomaticEnv()
    
    // Configuration file
    viper.SetConfigName("conf")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/config/")
    if err := viper.ReadInConfig(); err != nil {
        logrus.WithError(err).Warning("Couldn't read configuration file")
    }
    
    // Defaults
    viper.SetDefault("server.host", "127.0.0.1")
    viper.SetDefault("server.port", 8080)
    
    // Parsing flags
    pflag.Parse()
}
```

This simple snippet will do so many things. First you will be able to customize
all your configuration with flags, environment variable and configuration file.
It also sets defaults and will generate a help text with the `-h,--help` flag.

If no configuration file is found, a warning is logged but we won't stop the
program. After all, there are defaults, flags and environment variables, the
user might want to setup the whole program without configuration file.

## Commands

![cobra](/assets/go-checklist/cobra.png)

This is actually not mandatory since it highly depends on how you program is
intended to work. But you might want to have multiple commands and subcommands
in your program. For example [smallblog](https://github.com/Depado/smallblog)
uses at least two commands for now, which are `serve` and `new`. To achieve
this I'm using the [cobra](https://github.com/spf13/cobra) library which
can be tightly integrated with viper. 

### Root command

Cobra allows us to split our program's behavior in multiple commands, and thus
allows us to have different flags and configuration options for the different
commands. Have a look at how smallblog achieves this : 

We first define our `Root` command 
[like this](https://github.com/Depado/smallblog/blob/master/cmd/root.go).
It contains what's called `PersistentFlags`, which are special flags that are
also passed to other commands and subcommands. Basically it's the global
configuration of your application. At the end of the `init()` function
we're doing this :

```go
// Flag binding
viper.BindPFlags(rootCmd.PersistentFlags())
```

We're telling viper to bind those persistent flags so we can access them through
viper in the rest of our program. Once everything is initialized, the function
we defined at the top of our `init()` function 
(`cobra.OnInitialize(initialize)`) will be executed.

In this function we're telling viper to parse the configuration file, setup
the environment variables, etc. We also setup the logger which we'll see in the
next section.

### Commands

We then proceed to define each individual command. For example, [here is the
code for the `serve` command in smallblog](https://github.com/Depado/smallblog/blob/master/cmd/serve.go).
Here we define the flags that are specific to the `serve` command. We also bind
those to viper. And now, when the `serve` command is executed, viper will parse
the root command flags, the serve command flags, and then read those in the 
environment and in the configuration file ! It's as simple as that !

## Logger setup

In every application you need to be able to setup your logger properly. In most
of my applications I'm using [logrus](https://github.com/sirupsen/logrus). Now
as you may have seen in the previous section, I usually define three flags
(three configuration values) which are : Level, Format and Line.

The level is the minimum level from which your logger should log. If you set
it to "warning" or "warn", then no log under this level will be displayed,
which means you won't get "info" and "debug" logs.

The format is actually self-explicit. Logrus supports two formats : text and
json. Both are useful, text is easy to read, and json is easy to parse.

And finally line. When this flag is active (or its configuration value is set
to true), then a hook is added to logrus which will also log as an extra field,
the name of the file and the line on which the log happened.

Here's the code I'm using :

```go
func setupLogger() {
	lvl := viper.GetString("log.level")
	l, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.WithField("level", lvl).Warn("Invalid log level, fallback to 'info'")
	} else {
		logrus.SetLevel(l)
	}
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if viper.GetBool("log.line") {
		logrus.AddHook(filename.NewHook())
	}
}
```

## Build Number and Version

When building a Go program you can inject variables at compile time. A good
thing to do, is to inject a build number and a version number. To store them
and to display them on startup, we're going to do something like this :

```go
package main

import (
	"github.com/sirupsen/logrus"
)

// Build variables injected at compile time
var (
	Build   = "unknown"
	Version = "unknown"
)

func main() {
	logrus.WithFields(logrus.Fields{
        "version": Version, 
        "build": Build,
    }).Info("Starting")
```

The build will be associated with our most recent commit sha1, and the version
will be an arbitrary value. We'll create a Makefile to handle both of these :

```makefile
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

clean:
	-rm $(BINARY)
```

That's it. Now if you `make`, your build number and version will be injected
in your binary. You can even create a new Cobra command to display those !

## Minimalist Dockerfile

```dockerfile
# Build Step
FROM golang:1.10 AS build

# Prerequisites and vendoring
RUN mkdir -p $GOPATH/src/path/to/your/project/
ADD . $GOPATH/src/path/to/your/project/
WORKDIR $GOPATH/src/path/to/your/project/
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only

# Build
ARG build
ARG version
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${version} -X main.Build=${build}" -o myprogram
RUN cp myprogram /

# Final Step
FROM alpine

# Base packages
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN rm -rf /var/cache/apk/*

# Copy binary from build step
COPY --from=build /myprogram /home/
# Define timezone
ENV TZ=Europe/Paris

# Define the ENTRYPOINT
WORKDIR /home
ENTRYPOINT ./myprogram

# Document that the service listens on port 8080.
EXPOSE 8080
```

Multiple things happen here. First we're using [multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/)
to separate the build of our program and the execution. So we're building
our binary in a step called `build` (as defined on this line : 
`FROM golang:1.10 AS build`). We parse the `build` and `version` arguments
as we defined those earlier, we'll get back to this later. We build the binary
and place it at the root of our container.

We then run the second step, the one that will build a small container with
our binary inside. Using alpine here guarantees that you have `ca-certificates`
and `tzdata`. 

We can now either build it on the command line... Or... Use the Makefile we 
defined earlier !

```makefile

docker:
	docker build -t "your-docker-repo/$(BINARY):$(VERSION)" \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .
```

And now all you have to do is to run `make docker` and that's it !

# Thanks

Thanks to [@ashleymcnamara](https://github.com/ashleymcnamara) for the amazing
[Gopher Artworks](https://github.com/ashleymcnamara/gophers).
