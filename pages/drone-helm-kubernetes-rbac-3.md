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

So your repository structure should look like this:

```
.
├── Dockerfile
├── Gopkg.lock
├── Gopkg.toml
├── helm
│   └── dummy
│       ├── charts
│       ├── Chart.yaml
│       ├── templates
│       │   ├── deployment.yaml
│       │   ├── _helpers.tpl
│       │   ├── ingress.yaml
│       │   ├── NOTES.txt
│       │   └── service.yaml
│       └── values.yaml
├── LICENSE
├── main.go
└── README.md
```

Here, we're going to modify both the `values.yaml` files to use sane defaults
for our chart, and more importantly `templates/` to add and modify the rendered
k8s manifests.

We can see that Helm created a pretty chart ensuring the best practices, with
some nice helpers. As you can see the `metadata` section is quite always the
same:

```yaml
metadata:
  name: {{ template "dummy.fullname" . }}
  labels:
    app: {{ template "dummy.name" . }}
    chart: {{ template "dummy.chart" . }}
    release: {{ .Release.Name }}
    heritage: {{ .Release.Service }}
```

This will ensure we can deploy our application multiple times without our 
resources colliding. 

## Values File

Let's open up the `dummy/values.yaml` file:

```yaml
replicaCount: 1
image:
  repository: nginx
  tag: stable
  pullPolicy: IfNotPresent
service:
  type: ClusterIP
  port: 80
ingress:
  enabled: false
  annotations: {}
  path: /
  hosts:
    - chart-example.local
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local
resources: {}
nodeSelector: {}
tolerations: []
affinity: {}
```

Those are the default values (as well as the accepted values) our Chart
currently understands. For now we're just going to modify the `image` section
to reflect the work we have done in the previous part with our image deployment
to GCR:

```yaml
image:
  repository: gcr.io/project-id/dummy
  tag: latest
  pullPolicy: Always
```

We are setting the `pullPolicy` to `Always` because our `latest` image can, and
will, change a lot over time. These are the default values, we'll be able to
tweak those values for specific deployments. More on that later in the article.

## Deployment Manifest

Remember when we created the [dummy](https://github.com/Depado/dummy) project
that had a single endpoint `/health` that answers a `200 OK` in the 
[previous part](/post/ci-cd-with-drone-kubernetes-and-helm-2) ? Well this
endpoint is going to come handy here. It is what's called a 
[liveness probe](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/).

We are going to use this endpoint as our readiness probe too. Liveness probes
are used by Kubernetes to ensure your container is still running and has the
expected behavior. If our liveness probe were to answer anything else than a 
`200 OK` status, Kubernetes would consider the program crashed and would fire up
a new pod before evicting this one. The readiness probe, on the other hand 
determines if the pod is ready to accept incoming connections. While this probe
doesn't serve a suitable answer, Kubernetes won't route any traffic to the pod.

In our case, this application is really dumb. We can use the `/health` route
for both the liveness probe and readiness probe. So we'll open up the 
`dummy/templates/deployment.yaml` file and edit this section:

```yaml
          livenessProbe:
            httpGet:
              path: /health # There
              port: http
          readinessProbe:
            httpGet:
              path: /health # And there
              port: http
```

And... Well that's it. Our deployment manifest is complete already since the
Chart created by Helm is flexible enough to allow us to define what we need in
our `values.yaml` file.

Let's execute this, and check that our deployment is correctly rendered. We're
going to run Helm in debug mode and dry-run mode so it prints out the rendered
manifests and doesn't apply our Chart for real. Also we're going to name our
release with the `-n staging` and we'll fake install it in the `staging`
namespace.

```
$ helm install --dry-run --debug -n staging --namespace staging dummy/ 
```

```yaml
# Source: dummy/templates/deployment.yaml
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: staging-dummy
  labels:
    app: dummy
    chart: dummy-0.1.0
    release: staging
    heritage: Tiller
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dummy
      release: staging
  template:
    metadata:
      labels:
        app: dummy
        release: staging
    spec:
      containers:
        - name: dummy
          image: "gcr.io/project-id/dummy:latest"
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /health
              port: http
          readinessProbe:
            httpGet:
              path: /health
              port: http
```

## Service Manifest

> A Kubernetes Service is an abstraction which defines a logical set of Pods and
> a policy by which to access them - sometimes called a micro-service. The set 
> of Pods targeted by a Service is (usually) determined by a Label Selector.

So a [Service](https://kubernetes.io/docs/concepts/services-networking/service/) 
in Kubernetes is a way to create a stable link to access dynamically created 
pods using selectors. Remember all the things in our `metadata.labels` section
in our manifests ? This is the way we're going to access our application!

So let's open our `dummy/templates/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ template "dummy.fullname" . }}
  labels:
    app: {{ template "dummy.name" . }}
    chart: {{ template "dummy.chart" . }}
    release: {{ .Release.Name }}
    heritage: {{ .Release.Service }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app: {{ template "dummy.name" . }}
    release: {{ .Release.Name }}
```

Something is off here. Our `targetPort` is wrong. Remember our Docker image and
our dummy Go program ? We listen and expose the `8080` port. No problem! We're
simply going to allow the `targetPort` value to be customized:

```yaml
  ports:
    - port: {{ .Values.service.port }}
      targetPort: {{ .Values.service.targetPort }}
      protocol: TCP
      name: http
```

Sounds better. Let's modify our `dummy/values.yaml` file:

```yaml
service:
  type: ClusterIP
  targetPort: 8080
  port: 80
```

And once more, let's run Helm in dry-run and check if everything matches:

```
$ helm install --dry-run --debug -n staging --namespace staging dummy/ 
```

```yaml
# Source: dummy/templates/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: staging-dummy
  labels:
    app: dummy
    chart: dummy-0.1.0
    release: staging
    heritage: Tiller
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
    app: dummy
    release: staging
```

## Ingress

Helm Charts are supposed to be independent from the platform Kubernetes is 
deployed to and the technologies used. So we need to let people decide whether 
or not to activate the Ingress and the annotations that are associated with it.

Enforcing the use of a [GCLB](https://cloud.google.com/load-balancing/) instead
of an nginx load balancer doesn't make sense in the default values. So we'll
introduce a new file, which will be specific to our own deployments. When you
use Helm, it provides several ways to override the values defined in 
`values.yaml`. First, you can provide your own values file. If Helm doesn't
find a key, it will fallback to the sane defaults we declared in our default
values file. 

So let's create our first "user-supplied" values, and let's name it 
`staging.yml`:

```yaml
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: "gce"
    kubernetes.io/ingress.global-static-ip-name: "dummy-staging"
  path: /*
  hosts:
    - staging.dummy.myshost.io
```

**Note**: We need to define the path to be `/*` and not just `/` because of how
GCLB works.

Here we're using the same techniques we saw in the 
[first part](/post/ci-cd-with-drone-kubernetes-and-helm-1#toc_9) to link a 
load balancer to a static IP. And then what happens if we run helm in dry-run
once more, but this time we give it our custom values file ?

```
$ helm install --dry-run --debug -n staging --namespace staging -f staging.yml dummy/
```

We now have an Ingress !

```yaml
# Source: dummy/templates/ingress.yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: staging-dummy
  labels:
    app: dummy
    chart: dummy-0.1.0
    release: staging
    heritage: Tiller
  annotations:
    kubernetes.io/ingress.class: gce
    kubernetes.io/ingress.global-static-ip-name: dummy-staging
    
spec:
  rules:
    - host: staging.dummy.myshost.io
      http:
        paths:
          - path: /*
            backend:
              serviceName: staging-dummy
              servicePort: http
```

# Pipeline

## Service Account

Before we jump in how to continuously deploy our staging application (and then
our prod) using Drone, we first need to retrieve the Tiller credentials we
created in the [first part of this series](/post/ci-cd-with-drone-kubernetes-and-helm-1#toc_8).

We are going to inject these credentials in Drone so it can use Helm within our
pipeline. So first we're going to retrieve the Tiller credentials:

```
$ kubectl -n kube-system get secrets | grep tiller
tiller-token-xxxx
$ kubectl get secret tiller-token-xxx -n kube-system -o yaml
apiVersion: v1
data:
  ca.crt: xxx
  namespace: xxx
  token: xxx
kind: Secret
metadata:
  annotations:
    kubernetes.io/service-account.name: tiller
    kubernetes.io/service-account.uid: xxxx-xxxx-xxxx-xxxx
  creationTimestamp: 2018-05-15T14:51:35Z
  name: tiller-token-xxxx
  namespace: kube-system
  resourceVersion: "860311"
  selfLink: /api/v1/namespaces/kube-system/secrets/tiller-token-xxx
  uid: xxxx-xxxx-xxx-xxx
type: kubernetes.io/service-account-token
```

We're going to need what's inside the `data.token`. And just a reminder, this
is base64 encoded data. And since we're kind with our Drone instance, we're
going to decode it for him:

```
echo "that very long token of yours" | base64 -w 0
```

Store this somewhere, we'll explain later where we're going to use it. Also,
let's retrieve the IP of your Kubernetes Master:

```
$ kubectl cluster-info
Kubernetes master is running at <your master IP>
...
To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

## The Drone-Helm plugin

We are going to use the [drone-helm plugin](https://github.com/ipedrazas/drone-helm)
to automatically execute our Helm command. This plugin expects two secrets: 
`api_server` and `kubernetes_token`.

So we're going to create these secrets in our Drone instance:

```
$ drone secret add --image quay.io/ipedrazas/drone-helm --repository repo/dummy \
  --name kubernetes_token --value <the token you base64 decoded earlier>
$ drone secret add --image quay.io/ipedrazas/drone-helm --repository repo/dummy \
  --name api_server --value <your master IP>
```

And now it's time to configure our pipeline. I'll include the GCR part from the
previous article as well as the drone-helm plugin usage:

```yaml
  gcr:
    image: plugins/gcr
    repo: project-id/dummy
    tags: latest
    secrets: [google_credentials]
    when:
      event: push
      branch: master

  helm_deploy_staging:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./helm/dummy
    release: "staging"
    wait: true
    recreate_pods: true
    service_account: tiller
    secrets: [api_server, kubernetes_token]
    values_files: ["helm/staging.yml"]
    namespace: staging
    when:
      event: push
      branch: master
```

This is pretty self-explanatory when you're reading the docs but I'll explain it
anyway:

When there's a push on the master branch, first we're going to build and push
our Docker image to GCR. Then we're going to execute the `drone-helm` plugin,
giving it the path to our chart relative to our repository (`helm/dummy`).
We name our release `staging` in the namespace `staging` and we're using the
`tiller` service account. We're also going to `wait` for all the resources to
be created or recreated before exiting. Also, since we're using the `latest`
image we specify we want to recreate the pods using the `recreate_pods` option.

That's it. Now every time we push to master, we're going to update our staging
environment, given that all the tests pass. 

## Production

If you've learned things in this article series, you'll now understand what
makes Helm so special. Let's create a new file, and name it `prod.yml` (still in
our `helm/` directory):

```yaml
image:
  tag: 1.0.0
  pullPolicy: IfNotPresent
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: "gce"
    kubernetes.io/ingress.global-static-ip-name: "dummy-prod"
  path: /*
  hosts:
    - prod.dummy.myshost.io
```

And that's it. Now let's add these few lines to our Drone pipeline:

```yaml

  tagged_gcr:
    image: plugins/gcr
    repo: project-id/dummy
    tags: 
      - "${DRONE_TAG##v}"
      - "${DRONE_COMMIT_SHA}"
      - latest
    secrets: [google_credentials]
    when:
      event: tag
      branch: master

  helm_deploy_prod:
    image: quay.io/ipedrazas/drone-helm
    skip_tls_verify: true
    chart: ./helm/dummy
    release: "prod"
    wait: true
    recreate_pods: false
    service_account: tiller
    secrets: [api_server, kubernetes_token]
    values_files: ["helm/prod.yml"]
    values: image.tag=${DRONE_TAG##v}
    namespace: prod
    when:
      event: tag
      branch: master
```

And that's it. You now have a complete CI/CD pipeline that goes right into 
production when you tag a new release on Github. It will build the Docker image,
tag it with the given tag, the git commit's sha1, and the latest tag. It will
then use helm to deploy said image (using the tag) to our cluster.


# TLS

## Certificate Manifest

TODO
