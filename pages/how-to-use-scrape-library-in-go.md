title: "Small Example : How to use the scrape library in Go"
description: A short tutorial on how to use the scrape library in Go.
slug: how-to-use-scrape-library-in-go
date: 2015-08-07 11:21:00
tags:
    - go
    - dev
    - scrape

# Introduction

When I needed to scrap a website before, it was always getting a bit complicated to do something efficient. In Python for example, if you want to do something that's really efficient (meaning you have multiple pages to scrap, not only one) you had to implement the threads mecanism. Now, threads in Python are pretty neat, but let's admit it, there is no comparison with what's going on in Go with the goroutines. Goroutines are simpler, they are more efficient and it's actually a lot easier to share data between them.

This small example doesn't use complex mecanisms in Go. My problem was pretty simple : I wanted to list all the art galleries in Paris. So I found [this website](http://galerie-art-paris.com/) which is pretty great but can you see the problem here ? Galleries are splitted on several pages. I could do that by hand but hey... Why would I do that ?

# The scrape library

Scrape is a library that was written by [yhat](https://github.com/yhat). The API is quite simple but still really powerful. Of course there is not as much features as in, let's say, [BeautifulSoup](http://www.crummy.com/software/BeautifulSoup/). You can find the scrape library on [GitHub](https://github.com/yhat/scrape). In the [README.md](https://github.com/yhat/scrape/blob/master/README.md) there is a small example on how to use it and you can also find the complete documentation on [GoDoc](https://godoc.org/github.com/yhat/scrape). Now the example given in the README.md file is pretty minimalistic, and in this article, I'll attempt to show how to create a more complete program.

# Tutorial

## Scraping the front page to gather the links

If you have a look at the website I gave the link earlier, you can see that galleries are splitted on several pages. I could, of course, write in the program the 20 links that are on that page. But we're not going to do that.

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	urlRoot = "http://www.galerie-art-paris.com/"
)

func gatherNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent, "class") == "menu"
	}
	return false
}

func main() {
	resp, err := http.Get(urlRoot)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	as := scrape.FindAll(root, gatherNodes)
}
```

The `gatherNodes` function is called a matcher. A matcher is a function that takes a pointer to an HTML node and returns `true` if the HTML node satisfies the matcher. Here, the matcher is satisfied if the element is an anchor (`atom.A` in HTML it would correspond to the <a ...></a> tags), has a parent, and the parent's class is "menu".  Otherwise it returns `false` and the node is ignored. Now, `scrape.FindAll(root, gatherNodes)` will browse the HTML tree (root, which corresponds to the parsed `resp.Body`) and return a list of all the nodes that satisfies the matcher, in other words, the links I want to process.

## Parsing the other links asynchronously

If you see where this is going you can already tell what I'm going to do next. Let's define a new function that will be executed as a goroutine and takes an URL as a parameter.

```go
func scrapGalleries(url string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && scrape.Attr(n, "class") == "galerie-art-titre"
	}
	for _, g := range scrape.FindAll(root, matcher) {
		fmt.Println(scrape.Text(g))
	}
}
```

As you can see, the matcher is defined inline because it's a really simple one. That function will just scrape a page and display the results it found, in this case, all the galleries name that are present on the page (which is defined in a span (`atom.Span`) and has the "galerie-art-titre" class). Let's add a few lines to the main function :

```go
func main() {
        // ...
	as := scrape.FindAll(root, gatherNodes)
	for _, link := range as {
		go scrapGalleries(urlRoot + scrape.Attr(link, "href"))
	}
}
```

The list of HTML nodes is... Well just nodes. So if you want to get the actual url, you have to get what's in the `href` attribute, and append it to the urlRoot (in that case these are not absolute links, but it can depend on the website you're scrapping). Now if you execute this program as it is right now, nothing will happen because the program will automatically exit. It won't wait for the gouroutines to finish because that's not the default behaviour of gouroutines. So let's add a `sync.WaitGroup` and see what the full program looks like :

```go
package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	urlRoot = "http://www.galerie-art-paris.com/"
)

var wg sync.WaitGroup

func gatherNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent, "class") == "menu"
	}
	return false
}

func scrapGalleries(url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}
	matcher := func(n *html.Node) bool {
		return n.DataAtom == atom.Span && scrape.Attr(n, "class") == "galerie-art-titre"
	}
	for _, g := range scrape.FindAll(root, matcher) {
		fmt.Println(scrape.Text(g))
	}
}

func main() {
	resp, err := http.Get(urlRoot)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)
	}

	as := scrape.FindAll(root, gatherNodes)
	for _, link := range as {
		wg.Add(1)
		go scrapGalleries(urlRoot + scrape.Attr(link, "href"))
	}
	wg.Wait()
}
```

# Conclusion

Go is really (and I mean it, **really**) efficient when it comes to scrapping. That program scraps 21 pages in 140 ms. Of course it depends on the bandwidth you have and your CPU. But still, isn't this amazing ?
