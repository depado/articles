title: CI/CD with Drone, Kubernetes and Helm - Part 3
description: "Building your CI/CD pipeline with Drone, Kubernetes and Helm. RBAC included."
banner: "/assets/kube-drone-helm/banner.png"
slug: ci-cd-with-drone-kubernetes-and-helm-3
tags: ["ci-cd", "drone", "helm", "kubernetes", "rbac"]
date: "2018-05-29 18:10:00"
draft: true

# Helm Chart

## Creating the Chart

In the previous article we learned how to use an Helm Chart. In this section
we'll see how to create a basic chart that will simply create a deployment.

We won't bother about ingress, configmap and such because our goal here is
simply to use Helm from within our CI environment.

```
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