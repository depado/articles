title: Setting up Caddy Server on Debian
description: When there was a lack of documentation. I guess this is quite outdated now.
slug: setting-up-caddy-server-on-debian
date: 2015-12-07 09:25:00
tags:
    - dev
    - admin

# Introduction

[Caddy Server](https://caddyserver.com/) ~~looks like~~ is the next-gen web server. Here are some specs about Caddy that might be relevant :
 - HTTP/2 & HTTP
 - IPv6 & IPv4
 - Out of the Box [Let's Encrypt](https://letsencrypt.org/) support
 - Markdown
 - Websockets
 - Proxy and load balancer
 - FastCGI
 - Overly simple configuration files
 - ... And a lot more !

For a complete list of directives that Caddy supports you can head to the [official documentation](https://caddyserver.com/docs)

# Let's Encrypt Integration

As of version 0.8 of Caddy, it is now integrating Let's Encrypt. See [Caddy 0.8 Released with Let's Encrypt Integration](https://caddyserver.com/blog/caddy-0_8-released) on the Caddy Blog.

Why is that such a big deal ? What does it mean ? Let's Encrypt could be the subject of a whole blog post. In short terms it allows you to receive and use free SSL certificates and thus allows anyone to provide a secure layer to their website. For a long time, having a valid (signed everywhere) certificate was complicated and expensive. Meaning : If you wanted to provide a security layer to your website you had to pay. And not only you had to pay, but you were also compelled to prove that you're actually the owner of your domain by giving lot of information about yourself. Let's Encrypt breaks this wall. It brings security to any site owner, even those without the funds to pay for a valid certificate. It's the end of the x509 certificate era.

Now what does Caddy have to do with that ? Well Caddy... Automatically serves your sites with HTTPS by using Let's Encrypt. You don't have to do anything, you don't have to worry about the certificates. You don't even have to give out any personal information about yourself. It just abstracts the process of requesting a certificate and using it. It also abstracts the renewal of these certificates. Meaning that with a two line long configuration file like this :


```
depado.eu {
    proxy / localhost:8080
}
```

It will automatically request Let's Encrypt for a valid certificate and serve `depado.eu` with HTTPS by default. Isn't that just great ?
To get more information about that feature of caddy, head over to [the documentation about automatic https](https://caddyserver.com/docs/automatic-https).

# Installation

First of all, go to the [Caddy Server Download Page](https://caddyserver.com/download) and select the features you want, your architecture and operating system. Instead of click on the button, right click and copy the url.
Time to ssh into your server and start having fun. First of all, let's download caddy in a relevant location.

```
# mkdir /etc/caddy/
# wget "https://caddyserver.com/download/build?os=linux&arch=amd64&features=" /etc/caddy/caddy.tar.gz
# tar xvf /etc/caddy/caddy.tar.gz
```

You can give Caddy the rights to bind to low ports. To do so, here is the command you can execute :

```
# setcap cap_net_bind_service=+ep /etc/caddy/caddy
```

Let's create our first `Caddyfile` in `/etc/caddy/`. Edit `/etc/caddy/Caddyfile` and add something like that :

```
yourdomain.com {
	proxy / localhost:5000
}
```

Assuming something is running on port 5000, Caddy will then proxy every request for the `yourdomain.com` domain directly to the application running on that port. If you already have a website you want to bind to caddy, then head over to the [full documentation](https://caddyserver.com/docs) and see what directives are useful for you. Let's start caddy for the first time so that it can bind itself to Let'sEncrypt service.

```
# cd /etc/caddy/
# ./caddy
```

Caddy will then ask you for an email to give to the Let'sEncrypt service. If you don't wish to give that out, then don't, but keep in mind that you won't be able to recover your keys if you loose them. Our initial setup is done. Let's move on to the `supervisor` section.

# Supervisor configuration

In this guide I'll assume you have a functionning supervisor installation. It will allow us to execute caddy as a daemon.
First of all we'll edit the `/etc/supervisor/supervisord.conf` and add this line under the `[supervisord]` section :

```
minfds=4096
```

Why is that ? In a production environment, caddy will complain that the number of open file descriptor is too low. The reason is that supervisor's default value is too low (1024, instead of 4096 as recommended by caddy). Now let's add a new program to our supervisor configuration. Create the file `/etc/supervisord/conf.d/caddy.conf` :

```
[program:caddy]
directory=/etc/caddy/
command=/etc/caddy/caddy -conf="/etc/caddy/Caddyfile"
user=www-data
autostart=true
autorestart=true
stdout_logfile=/var/log/supervisor/caddy_stdout.log
stderr_logfile=/var/log/supervisor/caddy_stderr.log
```

You can customize the user to the one you want. As I said earlier, caddy now doesn't need root privileges to bind to low ports, so any user will do (prefer a user with few rights).
Caddy is now ready to be started by supervisor ! Simply add the program into supervisor and enjoy.

```
# supervisorctl reread
# supervisorctl add caddy
# supervisorctl start caddy
```
