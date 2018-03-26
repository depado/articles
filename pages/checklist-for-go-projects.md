title: Checklist for Go projects
description: ""
banner: ""
author: ""
slug: checklist-for-go-projects
tags: 
- go
- dev
- foss
date: "2018-03-26 16:54:46"
draft: true

# Introduction

I start new Go projects regularly, either as part of my daily work or for fun
on my free time. And every time I forget something ! Every time I forget about 
the configurable logger, the Makefile which injects the build number, the 
Dockerfile or a way to configure the project using environment variables.

And that's bad. Because I often realize I forgot something right when I need it.

So here is an article. It will be quite short I think, but will list everything
you should do when starting a Go project.

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

## Flexible Configuration

