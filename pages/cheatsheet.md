title: "Smallblog Cheatsheet"
description: |
    A post that should display all the different capabilities of Smallblog and
    provide quick visual checks that everything works. Can also be used as a 
    cheatlist.
slug: cheatsheet
banner: ""
draft: true
date: 2019-09-03 10:00:00
tags: [markdown,render]

# Cheatsheet

This post displays Smallblog's capabilities and can be used as a cheatsheet.
For every example, the raw markdown is given first.

# Code Blocks

## Specify Language

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello World")
    }
    ```

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

## Without Language

    ```
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello World")
    }
    ```

```
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

## Inline Code

    This `text` should be displayed as `code`.

    Also, let's `try with a lot of chars as well`.

This `text` should be displayed as `code`.

Also, let's `try with a lot of chars as well`.

# Admonitions

## With Title

    !!! note "Title"
        This is a `note` with a title.

!!! note "Title"
    This is a `note` with a title.

## Without Title

    !!! warning
        This is a `warning` admonition without a title.

!!! warning
    This is a `warning` admonition without a title.

## Extra Class

    !!! danger wow "wow"
        This a `danger` admonition with the `wow` CSS class.

!!! danger wow "wow"
    This a `danger` admonition with the `wow` CSS class.

# Tables

## With Heading

    | Column 1   | Column 2  |
    |------------|-----------|
    | First Row  | First Row |
    | Third Row  | Third Row |

| Column 1   | Column 2  |
|------------|-----------|
| First Row  | First Row |
| Third Row  | Third Row |

## Without Heading

    |            |            |
    |------------|------------|
    | First Row  | First Row  |
    | Second Row | Second Row |
    | Third Row  | Third Row  |

|            |            |
|------------|------------|
| First Row  | First Row  |
| Second Row | Second Row |
| Third Row  | Third Row  |

## Markdown Inside

    | Column 1       | Column 2       |
    |----------------|----------------|
    | **First Row**  |  _First Row_   |
    | `Second Row`   | ~~Second Row~~ |

| Column 1       | Column 2       |
|----------------|----------------|
| **First Row**  |  _First Row_   |
| `Second Row`   | ~~Second Row~~ |

# Paragraphs

Lorem ipsum dolor sit amet, **consectetur adipiscing** elit, sed do eiusmod 
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
quis nostrud _exercitation_ ullamco laboris nisi ut aliquip ex ea commodo 
consequat. Duis aute irure dolor in ~~reprehenderit~~ in voluptate velit esse 
cillum dolore eu fugiat nulla [pariatur](/). 

Excepteur sint occaecat cupidatat `non` proident, sunt in 
culpa qui officia deserunt mollit anim id est laborum.

# Pictures

## Simple Picture

    ![test](/assets/avatar.jpg)

![test](/assets/avatar.jpg)

## Picture with hover title

    ![test](/assets/avatar.jpg "Doge")

![test](/assets/avatar.jpg "Doge")

## Picture with Caption

    ![test](/assets/avatar.jpg) *Picture Caption or Description*

![test](/assets/avatar.jpg) *Picture Caption or Description*

## Side-by-side Pictures

!!! warning "Unimplemented"
    This feature is currently not implemented

# Footnote

    This is a footnote.[^1]
    [^1]: The footnote text.

This is a footnote.[^1]
[^1]: The footnote text.

# Lists

## Unordered

    - Simple bullet lists, starting with `+`, `-` or `*`
    - Add sublists by indenting
        - Either with a tab or multiple spaces
        - This will create a sublist

- Simple bullet lists, starting with `+`, `-` or `*`
- Add sublists by indenting
    - Either with a tab or multiple spaces
    - This will create a sublist

## Ordered

    1. Lorem ipsum dolor sit amet
    2. Consectetur adipiscing elit
    3. Integer molestie lorem at massa

1. Lorem ipsum dolor sit amet
2. Consectetur adipiscing elit
3. Integer molestie lorem at massa

## Definition List

    Golang
    :  AKA the Go Programming Language.

    Markdown
    :  Custom format to write stuff and translate it to HTML

Golang
:  AKA the Go Programming Language.

Markdown
:  Custom format to write stuff and translate it to HTML

# Quotes

## Anonymous Quote

    > This is a quote.
    > It can span on multiple lines.
    >
    > Like that I guess.

> This is a quote.
> It can span on multiple lines.
>
> Like that I guess.

## Citation Quote

    > This quote might be famous, so make sure you credit whoever said or wrote 
    > that. I mean, it's pretty important.
    >
    > <cite>Me</cite> 

> This quote might be famous, so make sure you credit whoever said or wrote 
> that. I mean, it's pretty important.
>
> <cite>Me</cite> 