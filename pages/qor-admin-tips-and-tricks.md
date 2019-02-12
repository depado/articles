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

## What is QOR?

> QOR is a set of libraries written in Go that abstracts common features needed 
> for business applications, CMSs, and E-commerce systems.
>
> <cite>[QOR GitHub Repository](https://github.com/qor/qor)</cite>

QOR is basically everything we need to run a full e-commerce website. But what
interests us here is the admin interface. Let's say we're creating an API
and store our data (stuff like your users, products, etc...) in a database. You
just don't need QOR to do that, you can use whatever router you're used to work
with for example. But QOR is a full framework, if you want to look at everything
it can do, head to the [getqor website](https://getqor.com/) and prepare 
yourself to be amazed because it can do so many things:

![qor features](/assets/qor-admin/qor-features.png)

Although this article isn't sponsored by QOR or anything like that, I think it's
a great piece of technology and the people behind it are quite amazing. They
even have an [enterprise package](https://getqor.com/en/enterprise)!

## QOR Admin?

Frameworks like [Django](https://www.djangoproject.com/) give you an admin
interface to manage your data. This is great because it allows you to display
data in a web interface, modify them or execute one-shot actions on some
records. For example, it's rather easy to create a CSV export of one of your
tables in the form a single button that any admin can click on, thus avoiding
the usual SQL query to export data, if you know what I mean.

It also enables to modify data and more importantly **keep it consistent** by
writing your business rules as part of your admin interface. So we can prevent
someone from modifying one of our products and set the price to $0 (or 0€, or
whatever the currency, you get my point). Or prevent data loss. Or set specific
behavior for certain fields. The admin interface use case is then completely 
different of a dashboard that does only "read" operations to generate insights
like [Metabase](https://www.metabase.com/) (which is an amazing tool too!).

So QOR Admin is a component of the QOR stack. And the great news is: You don't
need the whole QOR stack to make QOR Admin work! Plus, let's face it, QOR Admin
is really gorgeous with its material design theme (but that's subjective). The
following screenshot is from the [QOR Admin Demo](http://demo.getqor.com/admin):

![qor admin demo](/assets/qor-admin/qor-admin-demo.png)

## Motivations

Although QOR Admin is an amazing open-source lib and product, sometimes the
documentation lacks of a clear way to do things. This article's goal is to act
as a kind of enhanced documentation and tutorial. We'll try to leverage the
annoying steps of setting up QOR following its best practices and create some 
kind of package that could be reused quickly without having to overthink things.

# First Steps

## QOR Admin and Gorm

So QOR (in general and not just admin) is tightly coupled with 
[gorm](https://gorm.io), mostly because gorm is an amazing ORM for the Go
language, and also because [jinzhu](https://github.com/jinzhu) is an amazing
human being who created gorm and helped creating QOR. 

## QOR Admin and Gin

## Authentication



