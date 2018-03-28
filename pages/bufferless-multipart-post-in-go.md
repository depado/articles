title: Buffer-less Multipart POST in Golang
description: It's not as easy as it looks.
slug: bufferless-multipart-post-in-go
date: 2016-01-09 23:51:00
tags:
    - dev
    - go

# Introduction

For the client of my [goploader](https://up.depado.eu/) I started by using a simple POST method. Posting raw data was effective but there was a small problem. I couldn't name the file when serving it, so you ended up downloading things that were named `aefa3d32-c222-437e-4d6b-5181bca2d3d1` without even knowing the type of the file you're downloading. Of course when the content type can be determined, it's not really a problem but it is still inconvenient for the users. Around that time I had the idea of using multipart. My first idea was to have two fields, `file` and `name` which the server could understand.  

Then I realised that a multipart file upload would contain the name of the file anyway. I kept the `name` field in case the data source isn't properly a file but would come from `os.Stdin` for example. Also it allows to set a name that is different than the file name.  

I had a rough time understanding what was going on. A simple `http.Post` was pretty easy to do when you write raw data in it. A multipart post is somewhat more complicated and I ended up loading the whole file in ram which is... Bad. Also, I am using the [github.com/cheggaaa/pb](https://github.com/cheggaaa/pb) progress-bar and it didn't make any sense to monitor the speed in which the file is read from disk to memory. (“Wow my connection is blazing fast, 350Mo/s !”)

# Enters `io.Pipe()`

“Pipe creates a synchronous in-memory pipe. It can be used to connect code expecting an io.Reader with code expecting an io.Writer. Reads on one end are matched with writes on the other, copying data directly between the two; **there is no internal buffering**. It is safe to call Read and Write in parallel with each other or with Close. Close will complete once pending I/O is done. Parallel calls to Read, and parallel calls to Write, are also safe: the individual calls will be gated sequentially. ” - [Godoc about io.Pipe](https://golang.org/pkg/io/#Pipe)

`io.Pipe()` looks like exactly what we need as we are going to use a `multipart.Writer` to write the content of our file as the request body. But `http.Post()` takes an `io.Reader` as argument, not an `io.Writer`. The basic approach would be to write down the entire body in a byte buffer and then pass the said buffer to the request. What if we simply read the content while its written ? That's the role of `io.Pipe()`.

```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

const service = "https://url.of.your.service"

func main() {
	var err error
	var f *os.File
	var fi os.FileInfo
	var bar *pb.ProgressBar

	if f, err = os.Open("test.txt"); err != nil {
		log.Fatal(err)
	}
	if fi, err = f.Stat(); err != nil {
		log.Fatal(err)
	}
	bar = pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.Start()
}
```

Here we start by declaring a few variables and initialize them. We open a file (`test.txt`), store its information in an `os.FileInfo` so we can get the size when we initialize the bar. That program doesn't do much, nothing complicated here. Let's head to the multipart part.

```go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/cheggaaa/pb"
)

const service = "https://url.of.your.service"

func main() {
	var err error
	var f *os.File
	var fi os.FileInfo
	var bar *pb.ProgressBar

	if f, err = os.Open("test.txt"); err != nil {
		log.Fatal(err)
	}
	if fi, err = f.Stat(); err != nil {
		log.Fatal(err)
	}
	bar = pb.New64(fi.Size()).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.Start()

	r, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	go func() {
		var part io.Writer
		defer w.Close()
		defer f.Close()

		if part, err = mpw.CreateFormFile("file", fi.Name()); err != nil {
			log.Fatal(err)
		}
		part = io.MultiWriter(part, bar)
		if _, err = io.Copy(part, f); err != nil {
			log.Fatal(err)
		}
		if err = mpw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Post(service, mpw.FormDataContentType(), r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(ret))
}
```

First of all we start by creating our pipe, and our `multipart.Writer` which will write on the “write” end of the pipe. The next thing we do is start a goroutine. It will first create the `file` field, attributing to it the name of our file using the `os.FileInfo` we gathered earlier. The role of this goroutine will be to write the content of our file into a reader that will be read at the same time by our `http.Post` so that there is no buffering. As we also want to update the progress bar during this process, we make `part` a multiple writer (it will write both to `part` and `bar`). We then copy the content of our file right into our `part` and don't forget to close the multipart writer at the end, otherwise the server won't understand.  

The rest of the program is pretty classic, we read the response of the server and print it to stdout.  

Hope this helps !
