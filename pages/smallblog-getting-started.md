title: Smallblog Documentation
description: Because even though it's on Github, I'll write it here too.
slug: smallblog-doc
author: Depado
date: 2016-05-08 16:31:00
tags:
    - setup
    - markdown

# Install

You can either go to the [github release page](https://github.com/Depado/smallblog/releases) and download the latest release for your OS (currently only Linux 64Bit and Linux ARM are available) which consists of a `.tar.gz` archive. This archive contains everything you need to get started.

```
$ wget https://github.com/Depado/smallblog/releases/download/v1.0.1/smallblog-linux-amd64-1.0.1.tar.gz
$ tar xvf smallblog-linux-amd64-1.0.0.tar.gz
$ cd smallblog-linux-amd64
```
Once you're there, you must edit your `conf.yml`. (See the next section to do so). The `conf.yml` provided with this release contains some already filled values, including the `pages_dir` one.

Once you're done configuring Smallblog, head to the `pages` directory. There you'll find a starter page which you can either edit or remove if you want to (in which case you'll have to create a page from scratch). And that's it. You're all set. You just have to start the server.

```
$ ./smallblog
[SB] [INFO] [127.691µs] [first-article.md] [/post/first-article] First Article
[SB] [INFO] Generated 1 pages in 171.166µs
```

Your article will be listed on your front page, and will be available (as mentionned in the logs) at `host:port/post/first-article` or whatever you set your slug to (or let the server generate the slug for you from the title of your article).

# Configure

Put a `conf.yml` file next to your `smallblog` binary. Here are the options you
can customize

| Key         | Description                                                               | Default     |
| ----------- | ------------------------------------------------------------------------- | ----------- |
| host        | Interface on which the server should listen.                              | "127.0.0.1" |
| port        | Port on which the server should listen.                                   | 8080        |
| debug       | Activates the router's debug mode.                                        | false       |
| pages_dir   | Local or absolute path to the directory in which your articles are stored | "pages"     |
| title       | Blog title (front page)                                                   | ""          |
| description | Blog Description (front page)                                             | ""          |

# Write Posts

There is no naming convention for file names. You can name them whatever you
want, it won't chage the server's behaviour. A post (or page/article) file is
divided in two parts. The first part is yaml data. The second part is the actual
content of your article. The two parts are separated by a blank line.

Here is the list of yaml values you can fill

| Key         | Description                                                                           | Mandatory |
| ----------- | ------------------------------------------------------------------------------------- | --------- |
| title       | The title of your article.                                                            | **Yes**   |
| description | The description of your article (sub-title)                                           | No        |
| slug        | The link you want for your article. If left empty, will be generated from title.      | No        |
| author      | Author of the article                                                                 | No        |
| date        | The date of writing/publication of your article.                                      | **Yes**   |
| tags        | A list of tags you want to apply on the article (useless right now, but still pretty) | No        |

If any of the two mandatory values (`date` and `title`) are omitted, the parser will complain and simply ignore the file.

## Example Post

`pages/first-article`

```markdown
title: First Article
description: The reasons I made SmallBlog
slug: first-article
author: Depado
date: 2016-05-06 11:22:00
tags:
    - inspiration
    - dev

# Actual Markdown Content
Notice the blank line right after the `tags` list.
That's how you tell the parser that you are done with yaml format.
And that's a really long line you don't want to type every other day because it's excessively long.
```

This article will be parsed, and available at `example.com/post/first-article`.
It will also be listed at `example.com/`.

# Filesystem Monitoring

The directory you define in your `conf.yml` file is constantly watched by the
server. Which means several things :
 - If you create a new file, it will be parsed and added to your site.
   (Also if you `mv` a file inside the directory)
 - If you modify an exisiting file, it will be parsed and modified on your site
   if necessary (e.g if the slug changes).
 - If you delete an existing file, the article will be removed. (Also if you
   `mv` a file out of the directory)

All these changes are instant. Usually a file takes ~250µs to be parsed. When
you restart the server, all the files will be parsed again so they are stored in
RAM (which is really efficient unless you have 250Mo of markdown file).
