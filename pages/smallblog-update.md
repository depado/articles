title: "Smallblog Update & Capabilities"
description: |
    A list of new capabilities provided by Smallblog, and more importantly how
    they were implemented and how you can also implement them on your own.
slug: smallblog-update
banner: ""
draft: true
date: 2019-08-26 14:00:00
tags: [go,dev,markdown]

# Visual Changes

## Wider Code Blocks

When the size of the browser allows it, code blocks are now slightly larger in
order to limit the overflow and the display of scroll bars. There's no change
in mobile mode. 

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

## Auto Heading Link

Now every heading has an anchor that is only visible when hovered. Clicking said
anchor will scroll the page to the related heading, modifying the URL, just like
when we previously clicked in the Table of Content. This enables a new way of
sharing posts and articles and linking to specific parts without having to go
to the table of contents to do so.

## Refined Homepage

The homepage, listing all the posts has been reworked in order to be prettier.
This design also works better in mobile mode since there won't be any 
overlapping elements (such as the publication dates before). The new design also
allows to see draft posts quicker and faster due to color change. The 
description of the posts have also been added to make it easier to give your
readers a quick summary.

![new home](/assets/smallblog-update/newhome.png)

# Capabilities

## Admonitions

You can now use 
[admonitions](https://python-markdown.github.io/extensions/admonition/) to 
improve the readability of your posts. These blocks provide visual hints with
a clear color code. Unfortunately, this feature only allows to define a single
paragraph per admonition. I'll try to work on that, but the current state of
[Blackfriday](https://github.com/russross/blackfriday/tree/v2), the Markdown
renderer I'm currently extending, doesn't allow me to do that just yet. 
Nonetheless, you can now use admonitions like so:

```mardkown
!!! note "Note Title"
    Hello! I'm contained inside a simple note with a title!

!!! warning
    Hey. I'm contained in a warning without a title.

!!! warning "Wait, I found my title"
    Hey. I'm contained in a warning admonition **with** a title.

!!! danger "Danger Ahead"
    This is really dangerous. You shouldn't do that.

!!! info extra "Info with extra class"
    This admonition contains an extra CSS class used to modify the behavior of
    a single admonition.
```

!!! note "Note Title"
    Hello! I'm contained inside a simple note with a title!

!!! warning
    Hey. I'm contained in a warning without a title.

!!! warning "Wait, I found my title"
    Hey. I'm contained in a warning admonition **with** a title.

!!! danger "Danger Ahead"
    This is really dangerous. You shouldn't do that.

!!! info extra "Info with extra class"
    This admonition contains an extra CSS class used to modify the behavior of
    a single admonition.

The possible types of admonitions are: 

- `note`
- `info`
- `hint`/`tip`
- `question`
- `success`
- `warning`/`caution`
- `danger`/`error`

These are all defined directly using CSS, and you can add more if you'd like to.

The type of the admonition is always required, otherwise it won't be parsed. The
title is optional and you can add extra CSS classes in between the type and the
title.

## Caption Pictures

You can now add a caption to the pictures you insert in your posts.
This feature can be achieved by adding the caption in an `em` element 
(basically, surrounding the caption with `*`) right next or below the picture 
insertion.

```markdown
![alt](/path/to/image.jpg)
*Caption or Description*
```

```markdown
![alt](/path/to/image.jpg) *Caption or Description*
```

![test](/assets/avatar.jpg) *This is the caption*

## Stylish Details

Detail tags are useful when you want users to click on some content to reveal
it. Smallblog now have some styled details that follows the same design rules
as the admonitions. 

<details class="warning"><summary>Click to Open</summary>

You can **easily** add markdown in there since the `HTML` tags are not 
interpreted. You can even add codeblocks in there!

```go
fmt.Println("That is amazing.")
```
</details>

In addition to all the classes that match the admonitions ones, there's a
special class that allows for clean code display: `code`.

<details class="code"><summary>Code</summary>

    <details class="warning"><summary>Click to Open</summary>

    You can **easily** add markdown in there since the `HTML` tags are not 
    interpreted. You can even add codeblocks in there!

    ```go
    fmt.Println("That is amazing.")
    ```
    </details>
</details>

This `code` details block is intended to have one big chunk of code inside and
nothing else. You can also have a details block already open.

<details open=""><summary>Hey there</summary>

I'm already open you see. So people can see me first and hide me if I'm really 
boring them. 
</details>

Have a look at the [Smallblog Cheatsheet](/post/cheatsheet) for
more information.

# Technical Changes

In this section we'll see how those features were implemented and what changed
in the codebase. 

## From CSS to SASS

Smallblog was always a small project with very few CSS, the only external CSS
library used is [pure](https://purecss.io/)'s grid system. When developping the
admonition feature, I realized a lot of CSS was just copy-pasted, basically 
every type of admonition just changes one thing: the color and the icon. There
was also a lot of repeating CSS all over the place. I started using SASS not so
long ago but it's a great tool, so I decided to migrate to SASS. For example
this allowed to create a mixin to handle all the admonition types like so:

```
@mixin admonition($color, $icon) {
  $dark: rgba($color, 0.1);
  border-left: 0.2rem solid $color;
  > p.admonition-title {
    background-color: $dark;
    font-size: 0.75rem;
    &::before {
      font-family: "Font Awesome 5 Free";
      font-weight: 900;
      margin-right: 0.4rem;
      content: $icon;
      color: rgba($color, 1);
    }
  }
}

.admonition {
  &.note {
    @include admonition(map-get($colors, "note"), map-get($icons, "note"));
  }
  &.info {
    @include admonition(map-get($colors, "info"), map-get($icons, "note"));
  }
  &.question {
    @include admonition(
      map-get($colors, "question"),
      map-get($icons, "question")
    );
  }
  // More of those admonition types
}
```

You can see how admonitions were implemented using SASS, and even grab this file
to generate your own version
[in the Smallblog repository](https://github.com/Depado/smallblog/tree/master/assets/sass).

## Blackfriday Renderers

In the past I created a few extra renderers for Blackfriday. The first I created
was called [bfchroma](https://github.com/Depado/bfchroma) and added a simple
integration of the [chroma](https://github.com/alecthomas/chroma) library which
highlights code without requiring any JS (meaning it produces already colored
output HTML with themes support). Codeblocks can now be rendered according to
the specified language in the triple backtick code:

    ```go
    package main

    func main() {
      print("hello")
    }
    ```

```go
package main

func main() {
  print("hello")
}
```

bfchroma supports a whole bunch of options when initializing the renderer and it
must extend another renderer (usually the Blackfriday HTML renderer with all its
settings). This is how you'd usually setup bfchroma with Blackfriday:

```go
myrenderer := bfchroma.NewRenderer(
  bfchroma.WithoutAutodetect(),
	bfchroma.ChromaOptions(
	  html.WithLineNumbers(),
	),
	bfchroma.Extend(
	  bf.NewHTMLRenderer(bf.HTMLRendererParameters{
	    Flags: flags,
		}),
	),
)
```

Now, this is fine. But what if we add another renderer on top of that? The 
second renderer I worked on was called 
[bfadmonition](https://github.com/Depado/bfadmonition/) and its role was to 
handle another custom block, the admonition one mentionned above. Let's see how
our renderer looks like now:

```go
myrenderer := bfchroma.NewRenderer(
  bfchroma.WithoutAutodetect(),
  bfchroma.ChromaOptions(
    html.WithLineNumbers(),
  ),
  bfchroma.Extend(
    bfadmonition.NewRenderer(
      bfadmonition.Extend(
        bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: flags}),
      ),
    ),
  ),
)
```

This is starting to look a little messy. What if we add yet another renderer on
top of that? Namely, one that generates an anchor on every HTML title 
(`h1`, `h2`, etc…). You see where this is going, don't you? A lot of tabs, a lot
of comas. That's why I wrapped all these renderers in a single one that I called
[bfplus](https://github.com/Depado/bfplus) and that integrates all these
previous renderers like so:

```go
myrenderer := bfp.NewRenderer(
  bfp.WithAdmonition(), // Enables admonition support
  bfp.WithHeadingAnchors(), // Enables heading anchors support
  bfp.WithCodeHighlighting( // Enables chroma support
    bfp.WithoutAutodetect(),
    bfp.ChromaOptions(
      html.WithLineNumbers(),
    ),
  ),
  bfp.Extend(
    bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: flags}),
  ),
)
```

This library doesn't import the previously mentionned lib, as to only have one
`RenderNode` method and avoid getting too deep in the call stack. (I'm guessing
it's more efficient that way instead of going through 5 `RenderNode` function 
every time)