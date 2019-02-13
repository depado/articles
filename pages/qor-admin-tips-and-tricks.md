title: "QOR Admin Tutorial"
description: |
    QOR Admin is a great admin interface creation tool but, just like gorm,
    the documentation sometimes lacks explanation and tutorials to help you get
    started
slug: qor-admin-tutorial
banner: "/assets/qor-admin/banner.png"
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

# Initial Setup

QOR (in general and not just admin) is tightly coupled with
[gorm](https://gorm.io), mostly because gorm is an amazing ORM for the Go
language, and also because [jinzhu](https://github.com/jinzhu) who created gorm
also helped creating QOR. The main issue here is that if you're not already
using gorm, you might have a hard time using it "only" for the admin interface.
For this reason we won't cover the case where you're not already using gorm for
your API or web service.

This section's goal is to provide a basic example of how to use QOR Admin and
Gorm. We'll be using this example throughout the article but if you already
have a use-case where you're connected to a database using gorm and already have
some tables in there with the associated structs you can skip this part
entirely.

## PostgreSQL

We'll first create our database and we will be using [PostgreSQL](https://www.postgresql.org/)
with the [uuid-ossp](https://www.postgresql.org/docs/10/uuid-ossp.html)
extension to generate UUID as primary key. This is mostly useless for this
example, but it's always nice to see another way of generating primary keys than
a simple `uint`.

```sql
postgres=# CREATE DATABASE qor_tutorial;
CREATE DATABASE
postgres=# CREATE USER qor WITH ENCRYPTED PASSWORD 'password';
CREATE ROLE
postgres=# GRANT ALL PRIVILEGE ON DATABASE qor_tutorial TO qor;
GRANT
postgres=# \c qor_tutorial
qor_tutorial=# CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
CREATE EXTENSION
```

## Base Program

Let's get started by creating a base program which we'll then use as a reference
for the rest of the article. In this snippet we'll do two things:

- Define a struct that is understandable by gorm
- Connect to the database we created earlier

```go
package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type product struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name  string
	Price int
}

func main() {
	var db *gorm.DB
	var err error
	if db, err = gorm.Open(
		"postgres",
		"user=qor dbname=qor_tutorial password=password sslmode=disable",
	); err != nil {
		logrus.WithError(err).Fatal("Couldn't initialize database connection")
	}
	defer db.Close()
}
```

This is a minimal example and this post isn't about gorm itself but about QOR
so we'll keep things simple. If you need more information on how to
use gorm, please refer to my [Gorm Gotchas](/post/gorm-gotchas) post. 

We'll add a few migrations which you can find in
[the code associated to this post](https://github.com/Depado/articles/blob/master/code/qor/v1/migrations/).
These migrations (`uuidCheck` and `initial`) are used to check if the
`uuid-ossp` extension exists in the database we're connecting to, and to create
our `products` table. Simply import the `migrate` package and run:

```go
if err = migrate.Start(db); err != nil {
	logrus.WithError(err).Fatal("Couldn't run migration")
}
```

So we have a database and a model. Now we just want to add QOR Admin. But we
want to use the [gin](https://github.com/gin-gonic/gin) router, so let's head to
[the QOR documentation](https://doc.getqor.com/admin/integration.html#integrate-with-gin)
and find out how we can use QOR Admin with it.

```go
adm := admin.New(&admin.AdminConfig{SiteName: "Admin", DB: db})
mux := http.NewServeMux()
adm.MountTo("/admin", mux)

r := gin.New()
r.Any("/admin/*resources", gin.WrapH(mux))
r.Run("127.0.0.1:8080")
```

## Adding Resource

Now let's add our `product` model to our new admin interface:

```go
// imports, db connection, etc.

adm := admin.New(&admin.AdminConfig{SiteName: "Admin", DB: db})
adm.AddResource(&product{})

// mux, gin, etc.
```

We can now head to `127.0.0.1:8080/admin` to see our admin interface in action!

![products](/assets/qor-admin/product.gif)

That's it. Now our product model is successfully registered with QOR, we can
interact with the data in our database, add/delete/edit products.

# Authentication

![lock](/assets/qor-admin/lock.gif)

Now we have a nice looking admin interface which integrates our model and allows
us to modify data in our database. But there's no authentication in front of it,
which is seriously problematic. Luckily, after a quick search on
[QOR's Documentation](https://doc.getqor.com/) we found a page about
authentication [here](https://doc.getqor.com/admin/authentication.html)!

We're told that we need to implement an interface. Well ok. How are we supposed
to do that? Well we need to understand how a simple auth system works, cookies
forms and other stuff like that. So let's get started.

## Admin User Model

The goal here is to create another table which will only contain our admin users
for our admin interface (people that are allowed to login and do things).
We don't need a really complex one, just enough to identify people connecting
on our admin interface, so let's go with something like that:

```go
// AdminUser defines how an admin user is represented in database
type AdminUser struct {
	gorm.Model
	Email     string `gorm:"not null;unique"`
	FirstName string
	LastName  string
	Password  []byte
	LastLogin *time.Time
}
```

Basic really. Email, name, password and last login operation. Security people
I see you, don't worry, the password won't be clear text. We're not barbarians.

So we're going to add a few methods:

```go
// DisplayName satisfies the interface for Qor Admin
func (u AdminUser) DisplayName() string {
	if u.FirstName != "" && u.LastName != "" {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	}
	return u.Email
}

// HashPassword is a simple utility function to hash the password sent via API
// before inserting it in database
func (u *AdminUser) HashPassword() error {
    pwd, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
    if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

// CheckPassword is a simple utility function to check the password given as raw
// against the user's hashed password
func (u AdminUser) CheckPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(raw)) == nil
}
```

That's right let's just use bcrypt. Now we need to create that table and add our
first user.

## Migration

If you read my previous post about gorm, you might know that I'm using a package
called [gormigrate](https://github.com/go-gormigrate/gormigrate) to handle
database migrations. This package is really useful because it allows to run
migrations and more importantly, to rollback if they fail. As I wrote in my
previous article, embedding your model inside the migration ensures the
migrations can be run in the right order even if you modify the main `AdminUser`
model. So let's create our table, and immediately add our first admin user
with its email address and a `changeme` password.

```go
import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	gormigrate "gopkg.in/gormigrate.v1"
)

var initAdmin = &gormigrate.Migration{
	ID: "init_admin",
	Migrate: func(tx *gorm.DB) error {
		var err error

		type adminUser struct {
			gorm.Model
			Email     string `gorm:"not null;unique"`
			FirstName string
			LastName  string
			Password  []byte
			LastLogin *time.Time
		}

		if err = tx.CreateTable(&adminUser{}).Error; err != nil {
			return err
		}
		var pwd []byte
		if pwd, err = bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost); err != nil {
			return err
		}
		usr := adminUser{
			Email:    "youremailaddress@yourcompany.com",
			Password: pwd,
		}
		return tx.Save(&usr).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.DropTable("admin_users").Error
	},
}
```

## Adding Admin User to Admin

If everything went well, you now have a new table called `admin_users` which
only contains one single record, the user you created in the migration. We're
going to add the `AdminUser` model to our admin interface **but** we need to
define the behavior intended for the `Password` field. QOR Admin is great but
if you don't tell it what to do with the data, it just puts it there. Meaning:
a clear text password if we modify our user in the admin interface and save it.

This is where QOR Admin can get tricky.

```go

```

# Deployment and Bindatafs

This part was tricky. As in, really tricky and it took me a lot of time to
actually understand what was going on.

# Thanks

- Gin-Gonic Framework Logo by
  [Javier Provecho](https://github.com/javierprovecho) is licensed under a
  [Creative Commons Attribution 4.0 International License](http://creativecommons.org/licenses/by/4.0/).
- Scientist Gopher by [marcusolsson](https://github.com/marcusolsson) from the [gophers repo](https://github.com/marcusolsson/gophers)
