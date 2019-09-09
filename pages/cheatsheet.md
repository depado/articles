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

<details class="code"><summary>Code</summary>

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello World")
    }
    ```
</details>

## Without Language

```
package main

import "fmt"

func main() {
    fmt.Println("Hello World")
}
```

<details class="code"><summary>Code</summary>

    ```
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello World")
    }
    ```
</details>

## Inline Code

This `text` should be displayed as `code`.

Also, let's `try with a lot of chars as well`.

<details class="code"><summary>Code</summary>

```md
This `text` should be displayed as `code`.

Also, let's `try with a lot of chars as well`.
```
</details>

# Admonitions

## With Title

!!! note "Title"
    This is a `note` with a title.

<details class="code"><summary>Code</summary>

```md
!!! note "Title"
    This is a `note` with a title.
```
</details>

## Without Title

!!! warning
    This is a `warning` admonition without a title.

<details class="code"><summary>Code</summary>

```md
!!! warning
    This is a `warning` admonition without a title.
```
</details>

## Extra Class

!!! danger wow "wow"
    This a `danger` admonition with the `wow` CSS class.

<details class="code"><summary>Code</summary>

```md
!!! danger wow "wow"
    This a `danger` admonition with the `wow` CSS class.
```
</details>

## Unknown Type

!!! unknown "Unknown Admonition Type"
    This is an `unknown` admonition, it will render just like a `note`/`info` 
    one.

<details class="code"><summary>Code</summary>

```md
!!! unknown "Unknown Admonition Type"
    This is an `unknown` admonition, it will render just like a `note`/`info` 
    one.
```
</details>

# Tables

## With Heading

| Column 1   | Column 2  |
|------------|-----------|
| First Row  | First Row |
| Third Row  | Third Row |

<details class="code"><summary>Code</summary>

```md
| Column 1   | Column 2  |
|------------|-----------|
| First Row  | First Row |
| Third Row  | Third Row |
```
</details>

## Without Heading

|            |            |
|------------|------------|
| First Row  | First Row  |
| Second Row | Second Row |
| Third Row  | Third Row  |

<details class="code"><summary>Code</summary>

```md
|            |            |
|------------|------------|
| First Row  | First Row  |
| Second Row | Second Row |
| Third Row  | Third Row  |
```
</details>

## Markdown Inside

| Column 1       | Column 2       |
|----------------|----------------|
| **First Row**  |  _First Row_   |
| `Second Row`   | ~~Second Row~~ |

<details class="code"><summary>Code</summary>

```md
| Column 1       | Column 2       |
|----------------|----------------|
| **First Row** |  _First Row_   |
| `Second Row`   | ~~Second Row~~ |
```
</details>

# Paragraphs

## Simple

Lorem ipsum dolor sit amet, **consectetur adipiscing** elit, sed do eiusmod 
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
quis nostrud _exercitation_ ullamco laboris nisi ut aliquip ex ea commodo 
consequat. Duis aute irure dolor in ~~reprehenderit~~ in voluptate velit esse 
cillum dolore eu fugiat nulla [pariatur](/). 

Excepteur sint occaecat cupidatat `non` proident, sunt in 
culpa qui officia deserunt mollit anim id est laborum.

<details class="code"><summary>Code</summary>

```md
Lorem ipsum dolor sit amet, **consectetur adipiscing** elit, sed do eiusmod 
tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, 
quis nostrud _exercitation_ ullamco laboris nisi ut aliquip ex ea commodo 
consequat. Duis aute irure dolor in ~~reprehenderit~~ in voluptate velit esse 
cillum dolore eu fugiat nulla [pariatur](/). 

Excepteur sint occaecat cupidatat `non` proident, sunt in 
culpa qui officia deserunt mollit anim id est laborum.
```
</details>

## Collapse

## Simple Collapse

<details><summary>Click to open</summary>

Markdown **content**

```go
fmt.Println("Can contain code")
```
</details>

<details class="code"><summary>Code</summary>

    <details>
    <summary>Click to open</summary>

    Markdown **content**

    ```go
    fmt.Println("Can contain code")
    ```
    </details>
</details>

## Already Opened

<details open=""><summary>Click to close</summary>

This is my markdown content.
</details>

<details class="code"><summary>Code</summary>

    <details open="">
    <summary>Click to close</summary>

    This is my markdown content.
    </details>
</details>

## With Admonition Class

<details class="warning" open=""><summary>Warning Details</summary>

A warning collapsible
</details>

<details class="code"><summary>Code</summary>

    <details class="warning" open="">
    <summary>Warning Details</summary>

    A warning collapsible
    </details>
</details>

## Code

<details class="code" open=""><summary>Code</summary>

```go
// This is a custom detail just for code showing
fmt.Println("CODE!")
```
</details>

<details class="code"><summary>Code</summary>

    <details class="code" open="">
    <summary>Code</summary>

    ```go
    // This is a custom detail just for code showing
    fmt.Println("CODE!")
    ```
    </details>
</details>

# Pictures

## Simple Picture

![test](/assets/avatar.jpg)

<details class="code"><summary>Code</summary>

```md
![test](/assets/avatar.jpg)
```
</details>

## Picture with hover title

![test](/assets/avatar.jpg "Doge")

<details class="code"><summary>Code</summary>

```md
![test](/assets/avatar.jpg "Doge")
```
</details>

## Picture with Caption

![test](/assets/avatar.jpg) *Picture Caption or Description*

<details class="code"><summary>Code</summary>

```md
![test](/assets/avatar.jpg) *Picture Caption or Description*
```
</details>

## Side-by-side Pictures

!!! warning "Unimplemented"
    This feature is currently not implemented

# Footnote

This is a footnote.[^1]
[^1]: The footnote text.

<details class="code"><summary>Code</summary>

    This is a footnote.[^1]
    [^1]: The footnote text.
</details>

# Lists

## Unordered

- Simple bullet lists, starting with `+`, `-` or `*`
- Add sublists by indenting
    - Either with a tab or multiple spaces
    - This will create a sublist

<details class="code"><summary>Code</summary>

    - Simple bullet lists, starting with `+`, `-` or `*`
    - Add sublists by indenting
        - Either with a tab or multiple spaces
        - This will create a sublist
</details>

## Ordered

1. Lorem ipsum dolor sit amet
2. Consectetur adipiscing elit
3. Integer molestie lorem at massa

<details class="code"><summary>Code</summary>

    1. Lorem ipsum dolor sit amet
    2. Consectetur adipiscing elit
    3. Integer molestie lorem at massa
</details>

## Definition List

Golang
:  AKA the Go Programming Language.

Markdown
:  Custom format to write stuff and translate it to HTML

<details class="code"><summary>Code</summary>

    Golang
    :  AKA the Go Programming Language.

    Markdown
    :  Custom format to write stuff and translate it to HTML
</details>

# Quotes

## Anonymous Quote

> This is a quote.
> It can span on multiple lines.
>
> Like that I guess.

<details class="code"><summary>Code</summary>

    > This is a quote.
    > It can span on multiple lines.
    >
    > Like that I guess.
</details>

## Citation Quote

> This quote might be famous, so make sure you credit whoever said or wrote 
> that. I mean, it's pretty important.
>
> <cite>Me</cite> 

<details class="code"><summary>Code</summary>

    > This quote might be famous, so make sure you credit whoever said or wrote 
    > that. I mean, it's pretty important.
    >
    > <cite>Me</cite> 
</details>