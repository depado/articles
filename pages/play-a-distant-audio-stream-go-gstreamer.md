title: Play a distant audio stream using Go and Gstreamer
description: Sometimes I'm facing weird problems.
slug: play-a-distant-audio-stream-go-gstreamer
author: Depado
date: 2015-07-15 21:21:00
tags:
    - go
    - dev
    - audio


# Requirements
First of all you need to have [GStreamer](http://gstreamer.freedesktop.org/) library installed on your computer since we're going to use a Go binding of it.
You can then install the [Go Bindings for GStreamer](https://github.com/ziutek/gst) using a simple go get :

```
go get github.com/ziutek/gst
```

# Play the stream !

```go
package main

import (
    "fmt"

    "github.com/ziutek/gst"
)

func main() {
    player := gst.ElementFactoryMake("playbin", "player")
    player.SetProperty("uri", "UrlToYourStreamHere.mp3")
    // Setting the state to gst.STATE_PLAYING starts playing the stream
    player.SetState(gst.STATE_PLAYING)
    fmt.Scanln()
    fmt.Println("Exiting")
}
```
Gst will handle the streaming itself, buffering and things like that are automated !

# Pause
If you want to pause the stream, just use `player.SetState(gst.STATE_PAUSED)`

That's it ! That was easy wasn't it ? :)
