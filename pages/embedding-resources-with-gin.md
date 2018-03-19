title: Embedding resources with rice.go in a Gin project
description: Because sometimes you want to package your resources within your binary.
slug: embedding-resources-with-gin
author: Depado
date: 2016-01-29 21:15:00
tags:
    - dev
    - gin
    - go

# Introduction

For the past few weeks I've been playing around with [gin](https://github.com/gin-gonic/gin) which pretty much covers all my needs when creating a web application. So, still about that [goploader](https://github.com/Depado/goploader) project of mine, I wanted to make the installation of the server part painless for people wanting to host the server themselves. What I had in mind was allowing people to download a single binary which would embed all the static assets (js, css, html templates, icons) and make the setup easy by first serving a form to automatically configure the server (which would generate a `conf.yml` file).  

Now I've also worked with [go.rice](https://github.com/GeertJohan/go.rice) which does a nice job at embedding resources in a binary file by simply generating go source files including all the assets in it. So, how can we make gin use those resources ? That will be covered in the first part of this tutorial. In the end I managed to do this, and I thought it was the end of the story. Except I didn't want to offer only a binary, but also an archive containing the static assets, which would then allow people to customize the look, information and content of the served web pages. Problem is : When the resources aren't embedded, go.rice doesn't check for relative paths, only for absolute paths. So when someone downloaded that archive, they would get an error telling them that the box wasn't found, although they had the right files in the right places.

# From r.LoadHTMLGlob() to r.SetHTMLTemplate()

Here we go. Ready to switch from on-disk templates and static files to embedded ones. Let's say this was your previous code :

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static/", "assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.Run(":8080")
}
```

Here we have the most classical project starter with gin. Loading all templates in `templates/`, serving static assets in `assets/` with the route `/static/`. Let's add some rice in this thing ! Adding the support for static files is pretty easy :

```go
// File : project/main.go
// Replace r.Static("/static/", "assets") with :
r.StaticFS("/static", rice.MustFindBox("assets").HTTPBox())
```

But things will get pretty complicated when it comes to templates. As you may know it, gin parses all the templates when it starts and doesn't dynamically load them when called. So you can't just give it an `HTTPBox()` like we did earlier for the static files. Instead we need to tell the engine to use some templates. No way around it than parsing them manually. We'll create a function called `InitAssetsTemplates` (it is exported because you may want your `main.go` file to remain clean and we'll put that inside an `utils` package) that will do that for us !

Let's first modify our `main.go` file :

```go
// File : project/main.go
package main

import (
	"log"
	"net/http"

	"github.com/Depado/articles/rice-gin/utils"
	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
)

func main() {
	var err error

	tbox, _ := rice.FindBox("templates")
	abox, _ := rice.FindBox("assets")

	r := gin.Default()
	if err = utils.InitAssetsTemplates(r, tbox, abox, "index.html"); err != nil {
		log.Fatal(err)
	}
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	r.Run(":8080")
}
```

Now we're talking. Stop writing your comment about how bad it is to not handle errors, and wait for the end of the article. Please. You'll see why it doesn't matter if an error is thrown or not at this point. For now we will handle only the templateBox (`tbox`) to load the templates into the engine. Let's look at what that function does, shall we ?

```go
// File : project/utils/router.go
package utils

import (
	"html/template"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
)

// InitAssetsTemplates initializes the router to use the rice boxes.
// r is our main router, tbox is our template rice box
// and names are the file names of the templates to load
func InitAssetsTemplates(r *gin.Engine, tbox, abox *rice.Box, names ...string) error {
	var err error
	var tmpl string
	var message *template.Template

	for _, x := range names {
		if tmpl, err = tbox.String(x); err != nil {
			return err
		}
		if message, err = template.New(x).Parse(tmpl); err != nil {
			return err
		}
		r.SetHTMLTemplate(message)
	}
	r.StaticFS("/static", abox.HTTPBox())
	return nil
}
```

Quite a lot of things to annotate here. First of all, the function declaration. We need to tell gin which files it needs to load in the engine, so we'll need to explicitly give the names of the templates we want to load. Then we will just cycle through the provided template names, load them, parse them and set them inside the engine with the name of the template being the key so that we can do `c.HTML(200, "index.html", gin.H{})` in our routes. Then we simply add the static route using our trusty assets box `abox`.

**“Yes but what happens if the boxes can't be loaded or found ?!”**  
Hey first of all, you calm down. Like right now. I told you we would come to that later. But if that's all you want to do (creating embedded binaries) you're good to go, just keep in mind that you **need** to generate the go source files using `rice embed-go` **before** compiling. Otherwise the boxes will never be found. And of course handle errors. (Ignoring them is only used in the next part).

# Not embedded ? No worries ! Fallback !

If you followed me right, what I wanted to do is that with the archive release, files should be loaded from disk. Now the thing is, rice.go doesn't do that. It will register the absolute path of the boxes, and try to find them no matter what at this exact location if it cannot be found in the binary. So let's handle the fallback !

```go
// File : project/utils/router.go
package utils

import (
	"html/template"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin"
)

// InitAssetsTemplates initializes the router to use the rice boxes.
// r is our main router, tbox is our template rice box, abox is our assets box
// and names are the file names of the templates to load
func InitAssetsTemplates(r *gin.Engine, tbox, abox *rice.Box, names ...string) error {
	var err error

	if tbox != nil {
		var tmpl string
		var message *template.Template
		for _, x := range names {
			if tmpl, err = tbox.String(x); err != nil {
				return err
			}
			if message, err = template.New(x).Parse(tmpl); err != nil {
				return err
			}
			r.SetHTMLTemplate(message)
		}
	} else {
		r.LoadHTMLGlob("templates/*")
	}

	if abox != nil {
		r.StaticFS("/static", abox.HTTPBox())
	} else {
		r.Static("/static", "assets")
	}
	return nil
}
```

Now that's why errors at that point didn't matter that much. We will check if the boxes are nil pointers, and if so, fallback to serving files and templates from disk. Embedded or not, your files will be served. Of course this will check if the `templates` directory exists, and if not it will panic in case the templates aren't embedded. Although there is no error thrown when the `assets` directory doesn't exist or can't be found.
