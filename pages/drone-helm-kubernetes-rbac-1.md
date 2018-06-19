title: "CI/CD with Drone, Kubernetes and Helm - Part 1"
description: "Building your CI/CD pipeline with Drone, Kubernetes and Helm. RBAC included."
banner: "/assets/kube-drone-helm/banner.png"
slug: ci-cd-with-drone-kubernetes-and-helm-1
tags: ["ci-cd", "drone", "helm", "kubernetes", "rbac"]
date: "2018-05-16 10:31:27"
draft: false

# Introduction

Continuous integration and delivery is hard. This is a fact everyone can agree
on. But now we have all this wonderful technology and the problems are mainly
"How do I plug this with that?" or "How do I make these two products work
together?"

Well, there's **never** a simple and universal answer to these questions. In
this article series we'll progressively build a complete pipeline for continuous
integration and delivery using three popular products, namely Kubernetes, Helm
and Drone.

This first article acts as an introduction to the various technology used
throughout the series. It is intended for beginners that have some
knowledge of Docker, how a container works and the basics of Kubernetes. You can
entirely skip it if you have a running k8s cluster and a running Drone instance.

## Steps

- Create a Kubernetes cluster with GKE
- Create a service account for Tiller
- Initialize Helm
- Add a repo to Helm
- Deploy Drone on the new k8s cluster
- Enable HTTPS on our new Drone instance

## Technology involved

### Drone

[Drone](https://drone.io/) is a Continuous Delivery platform built on Docker and
written in Go. Drone uses a simple YAML configuration file, a superset of
docker-compose, to define and execute Pipelines inside Docker containers.

It has the same approach as [Travis](https://travis-ci.org/), where you define
your pipeline as code in your repository. The cool feature is that every step in
your pipeline is executed in a Docker container. This may seem counter-intuitive
at first but it enables a great plugin system: Every plugin for Drone you might
use is a Docker image, which Drone will pull when needed. You have nothing to
install directly in Drone as you would do with Jenkins for example.

Another benefit of running inside Docker is that the
[installation procedure](http://docs.drone.io/installation/) for Drone is really
simple. But we're not going to install Drone on a bare-metal server or inside a
VM. More on that later in the tutorial.

### Kubernetes

> Kubernetes (commonly stylized as K8s) is an open-source
> container-orchestration system for automating deployment, scaling and
> management of containerized applications that was originally designed by
> Google and now maintained by the Cloud Native Computing Foundation. It aims to
> provide a "platform for automating deployment, scaling, and operations of
> application containers across clusters of hosts". It works with a range of
> container tools, including Docker.   
> <cite>[Wikipedia](https://en.wikipedia.org/wiki/Kubernetes) </cite>

Wikipedia summarizes k8s pretty well. Basically k8s abstracts the underlying
machines on which it runs and offers a platform where we can deploy our
applications. It is in charge of distributing our containers correctly on
different nodes so if one node shuts down or is disconnected from the network,
the application is still accessible while k8s works to repair the node or
provisions a new one for us.

I recommend at least reading [Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/)
for this tutorial.

### Helm

[Helm](https://helm.sh/) is the package manager for Kubernetes. It allows us to
create, maintain and deploy applications in a Kubernetes cluster.

Basically if you want to install something in your Kubernetes cluster you can
check if there's a Chart for it. For example we're going to use the Chart for
Drone to deploy it.

Helm allows you to deploy your application to different namespaces, change the
tag of your image and basically override every parameter you can put in your
Kubernetes deployment files when running it. This means you can use the same
chart to deploy your application in your staging environment and in production
simply by overriding some values on the command line or in a values file.

In this article we'll see how to use a preexisting chart. In the next one
we'll see how to create one from scratch.

## Disclaimers

In this tutorial, we'll use [Google Cloud Platform](https://cloud.google.com)
because it allows to create Kubernetes clusters easily and has a private
container registry which we'll use later.

# Kubernetes Cluster

<img src="/assets/kube-drone-helm/kube.png" style="max-height: 100px;" />

_You can skip this step if you already own a k8s cluster with a Kubernetes version above
1.8._

In this step we'll need the `gcloud` and `kubectl` CLI. Check out how to [install
the Google Cloud SDK](https://cloud.google.com/sdk/downloads) for your operating
system.

As said earlier, this tutorial isn't about creating and maintaining a Kubernetes
cluster. As such we're going to use [Google Kubernetes Engine](https://cloud.google.com/kubernetes-engine/)
to create our playground cluster. There are two options to create it: either
in the web interface offered by GCP, or directly using the `gcloud` command.
At the time of writing, the default version of k8s offered by Google is `1.8.8`,
but as long as you're above `1.8` you can pick whichever version you want.
_Even though there's no reason not to pick the highest stable version..._

The `1.8` choice is because in this version [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)
is activated by default and is the default authentication system.

To reduce the cost of your cluster you can modify the machine type, but try to
keep at least 3 nodes; this will allow zero-downtime migrations to different
machine types and upgrade k8s version if you ever want to keep this cluster
active and running.

To verify if your cluster is running, you can check the output of the following
command:

```
$ gcloud container clusters list
NAME       LOCATION        MASTER_VERSION  MASTER_IP    MACHINE_TYPE   NODE_VERSION  NUM_NODES  STATUS
mycluster  europe-west1-b  1.10.2-gke.1    <master ip>  custom-1-2048  1.10.2-gke.1  3          RUNNING
```

You should also get the `MASTER_IP`, `PROJECT`, and the `LOCATION` which I removed
from this snippet. From now on in the code snippets and command line examples,
`$LOCATION` will refer to your cluster's location, `$NAME` will refer to your
cluster's name, and `$PROJECT` will refer to your GCP project.

Once your cluster is running, you can then issue the following command to
retrieve the credentials to connect to your cluster:

```
$ gcloud container clusters get-credentials $NAME --zone $LOCATION --project $PROJECT
Fetching cluster endpoint and auth data.
kubeconfig entry generated for mycluster.
$ kubectl cluster-info
Kubernetes master is running at https://<master ip>
GLBCDefaultBackend is running at https://<master ip>/api/v1/namespaces/kube-system/services/default-http-backend/proxy
Heapster is running at https://<master ip>/api/v1/namespaces/kube-system/services/heapster/proxy
KubeDNS is running at https://<master ip>/api/v1/namespaces/kube-system/services/kube-dns/proxy
kubernetes-dashboard is running at https://<master ip>/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy
Metrics-server is running at https://<master ip>/api/v1/namespaces/kube-system/services/metrics-server/proxy
```

Now `kubectl` is configured to operate on your cluster. The last command will 
print out all the information you need to know about where your cluster is 
located.

# Helm and Tiller

<img src="/assets/kube-drone-helm/helm.png" style="max-height: 100px;" />

First of all we'll need the `helm` command. [See this page for installation
instructions](https://github.com/kubernetes/helm/blob/master/docs/install.md).

Helm is actually composed of two parts. Helm itself is the client, and Tiller
is the server. Tiller needs to be installed in our k8s cluster so Helm can
work with it, but first we're going to need a **service account** for Tiller.
Tiller must be able to interact with our k8s cluster, so it needs to
be able to create deployments, configmaps, secrets, and so on. Welcome to
**RBAC**.

So let's create a file named `tiller-rbac-config.yaml`


```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
```


In this yaml file we're declaring a [ServiceAccount](https://kubernetes.io/docs/admin/authorization/rbac/#service-account-permissions)
named tiller, and then we're declaring a [ClusterRoleBinding](https://kubernetes.io/docs/admin/authorization/rbac/#rolebinding-and-clusterrolebinding)
which associates the tiller service account to the cluster-admin authorization.

Now we can deploy tiller using the service account we just created like this:

```
$ helm init --service-account tiller
```

![tiller-service](/assets/kube-drone-helm/tiller-service.png)

Note that it's not necessarily good practice to deploy tiller this way. Using
RBAC, we can limit the actions Tiller can execute in our cluster and the
namespaces it can act on.
[See this documentation](https://github.com/kubernetes/helm/blob/master/docs/rbac.md)
to see how to use RBAC to restrict or modify the behavior of Tiller in your k8s
cluster.

This step is really important for the following parts of this series, as we'll
later use this service account to interact with k8s from Drone.

# Deploying Drone

<img src="/assets/kube-drone-helm/drone.png" style="max-height: 100px;" />

## Static IP

If you have a domain name and wish to associate a subdomain to your Drone
instance, you will have to create an external IP address in your Google Cloud
console. Give it a name and remember that name, we'll use it right after when
configuring the Drone chart.

Associate this static IP with your domain (and keep in mind DNS propagation
can take some time).

For the sake of this article, the external IP address name will be `drone-kube`
and the domain will be `drone.myhost.io`.

## Integration

First we need to setup Github integration for our Drone instance. Have a look
at [this documentation](http://docs.drone.io/install-for-github/) or if you're
using another version control system, check in the Drone documentation how to
create the proper integration. Currently, Drone supports the following VCS:

- [GitHub](http://docs.drone.io/install-for-github/)
- [GitLab](http://docs.drone.io/install-for-gitlab/)
- [Gitea](http://docs.drone.io/install-for-gitea/)
- [Gogs](http://docs.drone.io/install-for-gogs/)
- [Bitbucket Cloud](http://docs.drone.io/install-for-bitbucket-cloud/)
- [Bitbucket Server](http://docs.drone.io/install-for-bitbucket-server/)
- [Coding](http://docs.drone.io/install-for-coding/)

Keep in mind that if you're not using the Github integration, the changes in
the environment variables in the next section need to match.

## Chart and configuration

After a quick Google search, we can see there's a 
[Chart for Drone](https://github.com/kubernetes/charts/tree/master/stable/drone). 
We can have a look at the 
[configuration](https://github.com/kubernetes/charts/tree/master/stable/drone#configuration)
part for this Chart. We'll create a `values.yaml` file that will contain the
required information for our Drone instance to work properly.

```yaml
service:
  httpPort: 80
  nodePort: 32015
  type: NodePort
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: "gce"
    kubernetes.io/ingress.global-static-ip-name: "drone-kube"
    kubernetes.io/ingress.allow-http: "true"
  hosts:
    - drone.myhost.io
server:
  host: "http://drone.myhost.io"
  env:
    DRONE_PROVIDER: "github"
    DRONE_OPEN: "false"
    DRONE_GITHUB: "true"
    DRONE_ADMIN: "me"
    DRONE_GITHUB_CLIENT: "the github client secret you created earlier"
    DRONE_GITHUB_SECRET: "same thing with the secret"
```

Alright! We have our static IP associated with our domain. We have to put
the name of this reserved IP in the Ingress' annotations so it knows to which
IP it should bind. We're going to use a GCE load balancer, and since we don't
have a TLS certificate, we're going to tell Ingress that it's OK to accept
HTTP connections. (Please don't hit me, I promise we'll see how to enable TLS
later.)

We also declare all the variables used by Drone itself to communicate with our
VCS, in this case Github.

That's it. We're ready. Let's fire up Helm!

```
$ helm install --name mydrone -f values.yaml stable/drone
```

Given that your DNS record is now propagated, you should be able to access your
Drone instance using the `drone.myhost.io` URL!

# TLS

## Deploying cert-manager

In the past, we had [kube-lego](https://github.com/jetstack/kube-lego) which
is now deprecated in favor of [cert-manager](https://github.com/jetstack/cert-manager/).

[The documentation](http://cert-manager.readthedocs.io/en/latest/getting-started/2-installing.html)
states that installing cert-manager is as easy as running this command:

```
$ helm install --name cert-manager --namespace kube-system stable/cert-manager
```

## Creating an ACME Issuer

Cert-manager is composed of several components. It uses what's called [Custom Resource Definitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/)
and allows to use `kubectl` to control the certificates, issuers and so on.

An [Issuer](https://cert-manager.readthedocs.io/en/latest/reference/issuers.html) 
or [ClusterIssuer](https://cert-manager.readthedocs.io/en/latest/reference/clusterissuers.html) 
represents a certificate authority from which x509 certificates can be obtained.

The difference between an Issuer and a ClusterIssuer is that the Issuer can only
manage certificates in its own namespace and be called from within that 
namespace. The ClusterIssuer doesn't depend on a specific namespace.

We're going to create a Let'sEncrypt ClusterIssuer so we can issue a certificate
for our Drone instance and for our future deployments. Let's create a file named
`acme-issuer.yaml`:

```yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt
spec:
  acme:
    server: https://acme-v01.api.letsencrypt.org/directory
    email: your.email.address@gmail.com
    privateKeySecretRef:
      name: letsencrypt-production
    http01: {}
```

Here we're creating the ClusterIssuer with the HTTP challenge enabled. We're
only going to see this challenge in this article, refer to the 
[documentation](https://cert-manager.readthedocs.io/en/latest/) for more
information about challenges. 
**Remember to change the associated email address in your issuer !**

```
$ kubectl apply -f acme-issuer.yaml
```

We can also create a ClusterIssuer using Let'sEncrypt staging environment which
is more permissive with errors on requests. If you want to test out without
issuing true certificates, use this one instead. Create a new file 
`acme-staging-issuer.yaml`:

```yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: ClusterIssuer
metadata:
  name: letsencrypt-staging
spec:
  acme:
    server: https://acme-staging.api.letsencrypt.org/directory
    email: your.email.address@gmail.com
    privateKeySecretRef:
      name: letsencrypt-staging
    http01: {}
```

```
$ kubectl apply -f acme-staging-issuer.yaml
```

## Certificate

Now that we have our ClusterIssuer that is using the production of Let'sEncrypt,
we can create a manifest that will solve the ACME challenge for us. First we're
going to need the name of the ingress created by the Drone chart:

```
$ kubectl get ingress
NAME            HOSTS             ADDRESS       PORTS     AGE
mydrone-drone   drone.myhost.io   xx.xx.xx.xx   80        1h
```

Now that we have this information, let's create the `drone-cert.yaml` file:

```yaml
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: mydrone-drone
  namespace: default
spec:
  secretName: mydrone-drone-tls
  issuerRef:
    name: letsencrypt # This is where you put the name of your issuer
    kind: ClusterIssuer
  commonName: drone.myhost.io # Used for SAN
  dnsNames:
  - drone.myhost.io
  acme:
    config:
    - http01:
        ingress: mydrone-drone # The name of your ingress
      domains:
      - drone.myhost.io
```

There are many fields to explain here. Most of them are pretty explicit and can
be found [in the documentation](http://cert-manager.readthedocs.io/en/latest/tutorials/acme/http-validation.html)
about HTTP validation.

The important things here are:

- `spec.secretName`: The secret in which the certificate will be stored. Usually
  this will be prefixed with `-tls` so it doesn't get mixed up with other 
  secrets.
- `spec.issuerRef.name`: The named we defined earlier for our ClusterIssuer
- `spec.issuerRef.kind`: Specify that the issuer is a ClusterIssuer
- `spec.acme.config.http01.ingress`: The name of the ingress deployed with Drone

Now let's apply this:

```
$ kubectl apply -f drone-cert.yaml
$ kubectl get certificate
NAME            AGE
mydrone-drone   7m
$ kubectl describe certificate mydrone-drone
...
Events:
  Type     Reason                 Age              From                     Message
  ----     ------                 ----             ----                     -------
  Warning  ErrorCheckCertificate  33s              cert-manager-controller  Error checking existing TLS certificate: secret "mydrone-drone-tls" not found
  Normal   PrepareCertificate     33s              cert-manager-controller  Preparing certificate with issuer
  Normal   PresentChallenge       33s              cert-manager-controller  Presenting http-01 challenge for domain drone.myhost.io
  Normal   SelfCheck              32s              cert-manager-controller  Performing self-check for domain drone.myhost.io
  Normal   ObtainAuthorization    6s               cert-manager-controller  Obtained authorization for domain drone.myhost.io
  Normal   IssueCertificate       6s               cert-manager-controller  Issuing certificate...
  Normal   CertificateIssued      5s               cert-manager-controller  Certificate issued successfully
```

We need to wait for this last line to appear, the `CertificateIssued` event
before we can update our Ingress' values. This can take some time, be patient
as Google Cloud Load Balancers can take several minutes to update.

## Upgrade Drone's Values

Now that we have our secret containing the proper TLS certificate, we can go
back to our `values.yaml` file we used earlier to deploy Drone with its Chart
and add the TLS secret to the ingress section ! We're also going to disable
HTTP on our ingress (only HTTPS will be served), and modify our `server.host`
value to reflect this HTTPS change.

```yaml
service:
  httpPort: 80
  nodePort: 32015
  type: NodePort
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: "gce"
    kubernetes.io/ingress.global-static-ip-name: "drone-kube"
    kubernetes.io/ingress.allow-http: "false" # ← Let's disable HTTP and allow only HTTPS
  hosts:
    - drone.myhost.io
  # Add this ↓
  tls:
    - hosts:
      - drone.myhost.io
      secretName: mydrone-drone-tls
  # End
server:
  host: "https://drone.myhost.io" # ← Modify this too 
  env:
    DRONE_PROVIDER: "github"
    DRONE_OPEN: "false"
    DRONE_GITHUB: "true"
    DRONE_ADMIN: "me"
    DRONE_GITHUB_CLIENT: "the github client secret you created earlier"
    DRONE_GITHUB_SECRET: "same thing with the secret"
```

And we just have to upgrade our deployment:

```
$ helm upgrade mydrone -f values.yaml stable/drone
```

You're going to have to modify your Github application too. 

# Conclusion

In this article we saw how to deploy a Kubernetes cluster on GKE, how to create
a service account with the proper cluster role binding to deploy Tiller, how
to use helm and how to deploy a chart with the example of drone.

In the next article we'll see how to write a quality pipeline for a Go project as
well as how to push to Google Cloud Registry.

# Thanks

Thanks to [@shirley_leu](https://twitter.com/shirley_leu) for proofreading
this article and correcting my english mistakes !