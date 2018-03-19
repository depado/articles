title: Handle Dynamic Subdomains with Flask and Nginx
description: A short tutorial on how to achieve a per-user subdomain
slug: handle-dynamic-subdomains-flask-nginx
author: Depado
date: 2015-06-17 10:36:00
tags:
    - dev
    - python
    - nginx

# Part I : Problem !

First of all, why would you even want to do something like that ? Dynamic
subdomains I mean. Well for example that's what I used to create MarkdownBlog.
Each user that registers on here gets a subdomain of his slugged username. This
allows the users to feel more detached from the platform. Like "I got my own
.com blog !" and not something like `markdownblog.com/user`.

Admit it, it's pretty cool right ? So, how are we going to do that with Flask ?
The answer is pretty simple :
[Blueprints](http://flask.pocoo.org/docs/0.10/blueprints/) !

# Part II : The Flask side.

*Note : This tutorial will assume you have the basis of web application
development using Flask. If you don't I'd suggest you to read the
[quickstart page](http://flask.pocoo.org/docs/0.10/quickstart/).*

## Introduction and minimal application

Do you remember how to create a single-file application with Flask ? Let's start
fresh with a single-file app. I chose this way of showing blueprints dynamic
subdomains because it's a lot easier to understand. Once the example is complete
I'll give a small example architecture for a more modular way of using
blueprints (which is after all, the whole point of blueprints).

Let's start by creating a single python file (remember to use virtualenv and pip
for your projects) and name it `main.py`

```python
from flask import Flask
app = Flask(__name__)

@app.route('/')
def hello_world():
    return 'Hello World!'

if __name__ == '__main__':
    app.run()
```

Now if you run `python main.py` and head to `127.0.0.1:5000` you should see a
pretty "Hello World!" in your browser. That's the most basic web application you
can write with Flask, and it's on the front page of the Flask's quickstart
guide.

## Modifying the /etc/hosts file

**Warning : This operation is only useful for testing on your local/development
machine. Don't do that on your server, if configured properly, it already uses
bind ! In other words, DO NOT do that on any distant server.**

*Note : I assume you're running Linux. If you're running Windows or OSX (or
whatever really), you're on your own for that part. It's a mandatory step to
test blueprints on your local machine.*

Before we can go any further and in order to test the subdomain handling, we'll
modify the `/etc/hosts` file. The hosts file technique is pretty limited because
you need to declare each and everyone of the subdomains you're going to test.
Let's start by adding the following lines at the end of our `/etc/hosts` files :

```
127.0.0.1	flask.dev				localhost	silence
127.0.0.1	test.flask.dev			localhost	silence
127.0.0.1	othertest.flask.dev		localhost	silence
```

Now you can run your Flask server and go to `flask.dev:5000` in your browser.
Note that we added two other hosts, that are subdomains of our application.
Also note that `flask.dev` is just an alias for `localhost` or `127.0.0.1`.
Time to add the following line just under the declaration of our app :

```python
#...
app = Flask(__name__)
app.config['SERVER_NAME'] = 'flask.dev:5000'
#...
```

There ! Now we can start using blueprints. If we didn't do these steps, the
Flask dev server wouldn't understand the subdomains and wouldn't even bother to
use the blueprint we're going to define.

**Warning again : This variable (`SERVER_NAME`) must be set to the used domain,
so in production remember to put your domain instead of `flask.dev`**

## Using blueprints

```python
from flask import Flask
from flask import Blueprint

app = Flask(__name__)
app.config['SERVER_NAME'] = 'flask.dev:5000'

@app.route('/')
def hello_world():
    return 'Hello World!'

# Blueprint declaration
bp = Blueprint('subdomain', __name__, subdomain="<user>")

# Add a route to the blueprint
@bp.route("/")
def home(user):
    return 'Welcome to your subdomain, {}'.format(user)

# Register the blueprint into the application
app.register_blueprint(bp)

if __name__ == '__main__':
    app.run(debug=True)
```

Start the server and see what happens when you go to `flask.dev:5000` and then
`test.flask.dev:5000`. Same uri, different subdomain, different behaviour ! Now
you can start doing some useful stuff. I'm pretty sure you have tons of idea on
how to use the `user` variable you get. And of course, you're not forced to use
a dynamic subdomain, you can also declare your blueprint with just a string in
the `subdomain` parameter.

## Example architecture

```
.
├── app
│   ├── api
│   ├── forms
│   ├── models
│   ├── modules
│   │   └── blog
│   ├── static
│   │   ├── css
│   │   │   └── syntax
│   │   ├── fonts
│   │   ├── img
│   │   └── js
│   ├── templates
│   │   └── blog
│   ├── utils
│   └── views
├── database
└── env
```

As you can see, the only blueprint I'm using is located in the `modules/blog`
subdirectory. It's organized as a standard application, views, forms, etc...


# Part III : Nginx side.

*Note : Remember to change the `SERVER_NAME` configuration variable to match
your actual domain name in production*

Your application is ready, you know what to do with your dynamic subdomain.
There is still one problem left though. Nginx. What happens if someones goes to
`www.yourdomain.com` ? Your application will think `www` is a subdomain but it's
actually not. So what are you going to do about that ? Simple. Strip the `www`
part and redirect to the url without it but keeping the subdomain. Here is how I
did that for MarkdownBlog :

```nginx
server {
	listen 80;
	listen 443 ssl;

	ssl_certificate /usr/local/nginx/ssl/nginx.crt;
	ssl_certificate_key /usr/local/nginx/ssl/nginx.key;

	server_name ~^www\.(?<user>.+\.)?markdownblog\.com$;
	return 301 "$scheme://${user}markdownblog.com$request_uri";
}

server {
       listen 80;
       listen 443 ssl;

       ssl_certificate /usr/local/nginx/ssl/nginx.crt;
       ssl_certificate_key /usr/local/nginx/ssl/nginx.key;

       server_name ~^.+\.markdownblog\.com$ markdownblog.com;

       location / {
                proxy_set_header Host $http_host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_pass http://127.0.0.1:8085;
    }
}
```

I spent hours on that. The fact is that, the regex was correct but was never
executed. As said in the [Nginx Documentation about Wildcards Server Names](http://nginx.org/en/docs/http/server_names.html#wildcard_names "Nginx"), 
you can use the `.domain.com` to match every subdomain. This also includes the
`www.domain.com`, which I don't want, otherwise my application would think `www`
is a user (and a blog url) which is wrong. Fact is that the `.domain.com` syntax
isn't considered as a regex. Nginx executes regex tests if all the over url
checks failed. This behaviour made my regex pointless as it was not even
executed by nginx. That configuration is used to catch the `username` part in
the url and strip out the `www` part.
