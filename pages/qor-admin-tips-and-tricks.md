title: "QOR Admin: Tips and Tricks"
description: |
    QOR Admin is a great admin interface creation tool but, just like gorm,
    the documentation sometimes lacks explanation and tutorials to help you get
    started
slug: qor-admin-tips-and-tricks
banner: "/assets/dialogflow/banner.png"
draft: true
date: 2019-02-11 17:30:00
tags: [go,dev,admin,qor]

# Introduction

## What is QOR and QOR Admin?

> QOR is a set of libraries written in Go that abstracts common features needed 
> for business applications, CMSs, and E-commerce systems.
>
> <cite>[QOR GitHub Repository](https://github.com/qor/qor)</cite>

QOR is basically everything we need to run a full e-commerce website. But what
interests us here is the admin interface. Let's say we're creating an API
and store our data (stuff like your users, products, etc...) in a database. You
just don't need QOR to do that, you can use whatever router you're used to work
with for example. One thing to remember though is that QOR is using 
[gorm](http://gorm.io/) to do most its processing, so if your project isn't 
using it, it may be difficult to implement QOR Admin in your project. 

QOR is a full framework, if you want to look at everything it can do, head to
the [getqor website](https://getqor.com/) and prepare to be amazed because it
can do so many awesome things:

![qor features](/assets/qor-admin/qor-features.png)

But let's get back to our admin interface.
Frameworks like [Django](https://www.djangoproject.com/) give you an admin
interface to manage your data. This is great because it allows you to display
data in a web interface, modify them or execute one-shot actions on some
records. For example, it's rather easy to create a CSV export of one of your
tables in the form a single button that any admin can click on, thus avoiding
the usual SQL query to export data, if you know what I mean.

So QOR Admin is a component of the QOR stack. And the great news is: You don't
need the whole QOR stack to make QOR Admin work! Plus, let's face it, QOR Admin
is really gorgeous with its material design theme (but that's subjective). The
following screenshot is from the [QOR Admin Demo](http://demo.getqor.com/admin):

![qor admin demo](/assets/qor-admin/qor-admin-demo.png)

## Motivations

Although QOR Admin is an amazing open-source lib and product, sometimes the
documentation lacks of a clear way to do things. This article's goal is to act
as a kind of enhanced documentation and tutorial.

# First Steps

## QOR Admin and Gorm

So QOR (in general and not just admin) is tightly coupled with 
[gorm](https://gorm.io), mostly because gorm is an amazing ORM for the Go
language, and also because [jinzhu](https://github.com/jinzhu) is an amazing
human being who created gorm and helped creating QOR. 

## QOR Admin and Gin

## Authentication



