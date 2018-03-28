title: Rendering the TOC with Blackfriday
description: Because that's a nice feature to have in your markdown renderer.
slug: render-toc-with-blackfriday
date: 2016-07-19 20:17:00
tags:
    - go
    - dev
    - markdown

# Introduction

Sometimes I wonder how to do things. And then I remember there is a
documentation. In that case I wanted to render the table of contents (TOC) using
[Blackfriday](https://github.com/russross/blackfriday) so that articles written
and served with smallblog can have a nicer structure.

Fact is, Blackfriday (the markdown renderer smallblog is using) can already do
such a task for you, but there is no actual example of using this feature unless
you start looking at the source code.

# Setting up your own renderer

Blackfriday has a lot of flags you can customize your render with. Let's create
a new simple package.

```go
package renderer

import . "github.com/russross/blackfriday"

var flags = 0 |
	HTML_USE_XHTML |
	HTML_USE_SMARTYPANTS |
	HTML_SMARTYPANTS_FRACTIONS |
	HTML_SMARTYPANTS_DASHES |
	HTML_SMARTYPANTS_LATEX_DASHES |
	HTML_TOC

var extensions = 0 |
	EXTENSION_NO_INTRA_EMPHASIS |
	EXTENSION_TABLES |
	EXTENSION_FENCED_CODE |
	EXTENSION_AUTOLINK |
	EXTENSION_STRIKETHROUGH |
	EXTENSION_SPACE_HEADERS |
	EXTENSION_HEADER_IDS |
	EXTENSION_BACKSLASH_LINE_BREAK |
	EXTENSION_DEFINITION_LISTS

func HTML(input []byte) []byte {
    r := HtmlRenderer(flags, "", "")
    return MarkdownOptions(input, r, Options{Extensions: extensions})
}
```

Here we're simply re-using the flags and extensions used for the
`blackfriday.MarkdownCommon` function, with the exception we added the
`HTML_TOC` flag. Also note that we used the dot import mecanism to simplify our
code.

# Use your renderer

You can then use import your package and do something like this :

```go
package main

import (
    "fmt"

    "github.com/user/project/renderer"
)

func main() {
    fmt.Println(renderer.HTML([]byte(`your markdown`)))
}
```

You will notice that now you have a nice table of contents matching your titles
on top of your rendered HTML.
