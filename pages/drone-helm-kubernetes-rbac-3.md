title: CI/CD with Drone, Kubernetes and Helm - Part 3
description: "Building your CI/CD pipeline with Drone, Kubernetes and Helm. RBAC included."
banner: "/assets/kube-drone-helm/banner.png"
slug: ci-cd-with-drone-kubernetes-and-helm-3
tags: ["ci-cd", "drone", "helm", "kubernetes", "rbac"]
date: "2018-05-29 18:10:00"
draft: true

# Introduction

This is the third an final part of this article series. In [the first part](/post/ci-cd-with-drone-kubernetes-and-helm-1)
we learned how to:

- Start a Kubernetes cluster using [GKE](https://cloud.google.com/kubernetes-engine/)
- Deploy Tiller and use Helm
- Deploy a Drone instance using its Helm Chart
- Enable HTTPS on our Drone instance using [cert-manager](https://github.com/jetstack/cert-manager) 

In the [second part](/post/ci-cd-with-drone-kubernetes-and-helm-2) we created
our first Drone pipeline for an example project, in which we ran a linter, 
either [gometalinter](https://github.com/alecthomas/gometalinter) or 
[golangci-lint](https://github.com/golangci/golangci-lint), built the Docker
image and push it to [GCR](https://cloud.google.com/container-registry/) with
appropriate tags according to the events of our VCS (push or tag).

In this last article, we'll see how to create an Helm Chart, and how we can
automate the upgrade/installation procedure directly from within our Drone 
pipeline.

# Helm Chart

## Creating the Chart

Helm provides us with a nice set of helpers. So let's go in our 
[dummy repo](https://github.com/Depado/dummy) and `create` our chart.

```
$ mkir helm
$ cd helm/
$ helm create dummy
Creating dummy
```

This will create a new directory `dummy` where you are. This directory
will contain two directories and some files:

- `charts/` A directory containing any charts upon which this chart depends
- `templates/` A directory of templates that, when combined with values, will generate 
  valid Kubernetes manifest files
- `Charts.yaml` A YAML file containing information about the chart
- `values.yaml` The default configuration values for this chart

For more information, check [the documentation](https://github.com/kubernetes/helm/blob/master/docs/charts.md#the-chart-file-structure)
about the chart file structure.

Here, we're going to modify both the `values.yaml` files to use sane defaults
for our chart, and more importantly `templates/` to add and modify the rendered
k8s manifests.