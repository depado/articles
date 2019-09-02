title: Cheap Open-Source Stack
description: Because cheap doesn't mean it's not perfect for you
slug: cheap-open-source-stack
date: 2017-08-18 18:00:00
tags:
    - inspiration
    - stack
    - server
    - ci-cd
draft: true

# Title

## Subtitle

!!! note "Optional Title"
    This the is the *content* of my **note**

Other markdown content

!!! note "Note"
    This is a simple note.

!!! warning
	A warning without title

!!! danger "Dangerous Stuff Ahead"
    Danger admonition

!!! hint "Advice and Hints"
    Add advices and hints

!!! success "Success"
    This was a success.

Let's go back to non-admonition markdown now. 
**And see if that works properly.**

# H1

## H2

### H3

#### H4

##### H5

###### H6

## Goals and scope

For some time now I've been wanting to write a long article (or an article serie maybe) on what my stack for my personal projects look like. And also how to deploy your own.

This article will explain how to setup such a stack and get the better of the open-source world of today, get some free TLS certificate, and much, much more. If you've ever asked yourself one of these questions :

- How do I monitor my server(s) with beautiful graphs ?
- How do I monitor my apps ? 
- How can I get free TLS certificates so my sites are secure ?
- What's CI/CD ? How can I try one for free ?
- What if I don't want my code on public places ? Can I still do all that ?

Well then, this article is for you. For now it's a draft, mainly because I have no idea on how to structurate all these information in a logical way.

Most parts of this article will be independent from one another. 

## Prerequisites

I'll assume you can spare a few dollars for a domain name and a small server. And when I say "a few", it's like 3$/month.

I'll also assume you have some basic knowledge of how internet and computer works. I'm not going to explain in detail every part of the infrastructure (like, how DNS works, what's an IP, etc), but I'll try to provide appropriate resources if you want to dig further in the subject.

# Servers

## Bare-metal vs VMs

So I've been using [Scaleway](https://scaleway.com) for quite some time now. And people keep asking me why I prefer bare-metal. There's a simple reason : I don't like sharing. I don't like sharing my kernel, I don't like sharing my CPUs, I don't like sharing my disks. Alright, VMs are nice and isolated, but still. It's a question of preference. Also I had a pretty bad experience with a shared server, where the server would reboot for quite no reason without any warning. Long story short : I enjoy knowing I'm the only one using the hardware I'm renting.

So let's start by creating a small C1 instance on Scaleway. It's like, the cheapest bare-metal server, ever, even though it's basically a Raspberry Pi in the cloud. 

## Users
The first thing you'll want to do on your fresh installed server, is create a privileged user that isn't root. Root is a dangerous role. Like, for real, you don't want to do all your operations as root. 

## Security
## Firewall and fail2ban

# Monitoring

## Prometheus
## Grafana
## AlertManager
## Instrumenting an app

# CI and CD

## DroneCI