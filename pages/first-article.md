title: First Article
description: The reasons I made SmallBlog
slug: first-article
author: Depado
date: 2016-05-06 11:22:00
tags:
    - inspiration
    - dev
    - go
    - markdown

# Quick Challenge
The main goal of this project is to show how easily you can develop a flat file
blog with markdown as the primary writing language. It's not perfect, it will
never be, some people are already doing great things based on that idea, like
[Hugo](https://gohugo.io/) for instance. Let's note though, that's **not** a
static website generator.

---

# Behavior
As stated earlier, SmallBlog isn't a static website generator. It's a dynamic
web server, meaning that when you start it, it will parse the content of the
`pages` directory. While the server is running you can add new files to this
directory, they will be automatically parsed and added to the site. Same goes if
you modify the content of an already parsed file, it will be parsed once more.
You can even change the metadata of the file, the server will keep its
consistency.  

# Design
I already made a markdown based blog engine :
[MarkDownBlog](http://markdownblog.com). I'm tired of its design, it was my
first real project, and I made a lot of bad technology choices. Like, why would
I syntax-highlight on the server-side when there are perfectly good JS libs that
do that automatically on the client-side ? That's one less problem to deal with.

That's why I'm using [prismjs](http://prismjs.com/index.html) to do that. Plus
it integrates perfectly with
[Blackfriday](https://github.com/russross/blackfriday), the markdown parser I'm
using. Blackfriday provides a parser for Github formated markdown. Which means
you can use the notations of Github to write your articles, include the
triple-backticks notation to write some code.

I also decided not to use any CSS framework this time. Keeping things simple and
sober. And because I suck at designing things, I'd like to thank
[bettermotherfuckingwebsite.com](http://bettermotherfuckingwebsite.com/) for the
inspiration.

# How it works
When starting, the server will scan for the `pages` directory. For each file
present in that directory it will parse the headers that are written in the yaml
format, as well as the markdown under the header (which are separated by a blank
line). A map containing the slug of the posts will be created as well as a slice
of articles sorted by date. The map will be used for a low complexity access to
the articles by their slugs, which will result in fast performances, something
like 150µs to display a post.

This method has its flaws. The more posts you have, the more RAM the server will
consume.

While the server is running, it will also monitor the `pages` dir for new files,
modification of exisiting files, or removal. This way, the server doesn't need
to be restarted when you want to publish or remove one of your posts. If you
detect that something went wrong, the state of your filesystem doesn't match
your blog for example, then restart the server. I make mistakes sometimes. And
make sure to open an issue on Github so this can be fixed.

---

# Everything is grey
Just to be sure to please everyone, everything is grey. Customize the CSS to
match your color preferences.

# Rendering examples

```go
func main() {
    log.Println("Hello World !")
}
```
You can `go build` to build this, and here's a list for your pretty eyes :

 - With a first item there. With maybe some _italic_ stuff ?
 - And a second one. This time with **bold** content.
 - And a third one with `inlined` content !
 - ~~Not putting a fourth one to show striked-through text.~~
 - http://bettermotherfuckingwebsite.com/

> Here is how a quote looks like. I hope you like it because I didn't put much
work on styling that.   
> <cite>-- Benjamin Franklin </cite>

![doge](http://fanaru.com/doge/image/18361-doge-follow-your-dreams.jpg)
![spinning](http://ljdchost.com/ilzb1nb.gif)
