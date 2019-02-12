title: "Gorm Gotchas"
description: |
    Gorm is an amazing ORM for the Go programming language, but its
    documentation sometimes lack of explanation or tutorials.
slug: gorm-gotchas
date: 2017-12-20 12:00:00
tags:
    - go
    - dev
    - database

# Introduction

## What is Gorm ?

[Gorm](https://github.com/jinzhu/gorm) is an ORM for the Go programming 
language. Which means it abstracts the whole retrieval of data in a database and
mapping to a struct, which you would need to do manually without such a system. 
It uses struct tags to define the behavior of the fields in your structs, how 
they are represented in the database and how to access them, how multiple 
structs are related to each other, etc.

## Why this article ?

First of all I think gorm is fantastic, it just works perfectly. The only thing
I miss is that sometimes there is a lack of documentation/tutorial on how to do 
certain things without diving into gorm's code. Which requires that you 
understand how an ORM works.

That's why I'm writing this article, not only for me to remember those tricks
but also to help other people understanding how things work with gorm.

# Global Tips

## Gorm's version and dependencies

Gorm currently has only one release, with a tagged version. Which means that if
you're currently using [dep](https://github.com/golang/dep) (which you should
be using, really), it will download [the only tagged version](https://github.com/jinzhu/gorm/releases/tag/v1.0)
which is over one year old. You're going to miss stability improvement, some
features (like JsonB support as mentioned later). 

So I would advise either telling `dep` to use the `master` branch or pin a
specific commit. (Actually the commit pinning seems a lot safer).

```toml
[[constraint]]
  name = "github.com/jinzhu/gorm"
  branch = "master"
```

This is a known issue. (see [#1543](https://github.com/jinzhu/gorm/issues/1543))

## Use a migration system from the beginning

Having a migration system from the beginning will make your life easier and your
development faster. It will allow you to flawlessly deploy a new version without
having to handle SQL migrations by hand and is an absolute requirement.
[gormigrate](https://github.com/go-gormigrate/gormigrate) is a very nice tool
that allows to handle migrations of your database with gorm. 


## The AutoMigrate case

So you just started playing with gorm and you discovered that awesome function
called [AutoMigrate](https://godoc.org/github.com/jinzhu/gorm#DB.AutoMigrate).
This function is really useful. It automatically migrates your tables schema
to match the struct definitions in your code. So it's tempting to just execute
this when your program starts right ?

There are a few things to note here. First of all AutoMigrate does **not**
handle foreign keys at the time of writing. That means, those relations between
your structs won't be taken into account when gorm will create the tables.

Instead, use [CreateTable](https://godoc.org/github.com/jinzhu/gorm#DB.CreateTable)
which can handle this for you. Let's write an example on how to do that using
gormigrate :

```go
package migrate

import (
	"time"

	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

type Article struct {
	gorm.Model
	Title string
	Slug  string `gorm:"unique_index"`
	Body  string
	Tags  []Tag `gorm:"many2many:article_tags;"`
}

type Tag struct {
	gorm.Model
	Name     string
	Articles []Article `gorm:"many2many:article_tags;"`
}

// Start starts the migration process
func Start(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				return tx.CreateTable(&Article{}, &Tag{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable(&Article{}, &Tag{}).Error
			},
		},
	})
	return m.Migrate()
}

```

As you can see this code actually has two structs (`article` and `tag`) which
are linked together by a many-to-many relation. If you used `AutoMigrate` here,
this relationship would have been overlooked and simply not present in your
database schema. The related issue is [#450](https://github.com/jinzhu/gorm/issues/450)
and there is some work being done on that.

As a side note, this example code uses a many-to-many with back-reference. Which
means from a single article we can get all the associated tags, but also from
a single tag, get all the associated articles. (Perfect use case for tags, isn't
it ?) This behavior is documented [here](http://jinzhu.me/gorm/associations.html#many-to-many).

## Add a migration flag to your binary (or not)

Using gormigrate or any other migration system is good. But you should be able
to specify when the migration runs and when it doesn't even though it all
depends on how your programs are being shipped and executed. The migration flag
is always active in my deployments because I want my database to be always up to
date, but I want to be able to disable that using a simple flag or a config
value. 

# Gormigrate

## Include your models directly within your migrations

When using gormigrate, you want to ensure all your migrations can run in
sequence, that's why you need to declare the structs directly in the migration
and not use the ones that are currently used in the app. So you can have an 
history of the evolution of your structs.

```go
// Start starts the migration process
func Start(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				type Article struct {
					gorm.Model
					Title string
					Slug  string `gorm:"unique_index"`
					Body  string
					Tags  []Tag `gorm:"many2many:article_tags;"`
				}
				type Tag struct {
					gorm.Model
					Name     string
					Articles []Article `gorm:"many2many:article_tags;"`
				}
				return tx.CreateTable(&Article{}, &Tag{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("articles", "tags").Error
			},
		},
	})
	return m.Migrate()
}
```

## Handle multiple drivers in your migrations

If you're using gorm you might want to be able to use multiple database drivers. 
After all that's one of the main reason you might want an ORM.
It's important that your migration can run with multiple drivers, but for 
example when using SQLite, some things are just no possible. You can easily 
detect which driver is being used inside a transaction as shown 
[here](https://github.com/kleister/kleister-api/blob/master/store/data/migration.go#L70).

```go
{
	ID: "201609011303",
	Migrate: func(tx *gorm.DB) error {
		if tx.Dialect().GetName() == "sqlite3" {
			return nil
		}

		return tx.Table(
			"team_users",
		).AddForeignKey(
			"team_id",
			"teams(id)",
			"RESTRICT",
			"RESTRICT",
		).Error
	},
	Rollback: func(tx *gorm.DB) error {
		if tx.Dialect().GetName() == "sqlite3" {
			return nil
		}

		return gormigrate.ErrRollbackImpossible
	},
}
```

This migration for example will be completely ignored if the driver is `sqlite3`.
Thanks to [tboerger](https://github.com/tboerger) for the tip on this one !

# Querying tips

## Populate struct fields from relations

Your database is created, your tables are good looking and brand new. You have a
many-to-many table which you're proud of. But how do you use this without having
to make multiple queries with gorm ? 

This part is actually covered nicely in the documentation but not necessarily
where you wish it would be. This part of the documentation is 
[Preloading (Eager Loading)](http://jinzhu.me/gorm/crud.html#preloading-eager-loading).
(I honestly think there should be more examples in that documentation)

```go
// Get that article you're so eager to read (you know, the one with ID 1) and
// also load all the related tags
a := &Article{}
db.Preload("Tags").First(a, 1)
spew.Dump(a)

// Or find all the articles that are tagged with "go" ?
t := &Tag{}
db.Preload("Articles").Where(&Tag{Name: "go"}).First(t)
spew.Dump(t)
```

Wait a second. In the second query we get the articles linked to the `go` tag.
But we don't have the tags of these articles we just retrieved. That's kind of
a loopy thing. Fortunately gorm knows how to handle that. 

```go
// Find all the articles tagged with "go" and also load their Tags
t := &Tag{}
db.Preload("Articles").Preload("Articles.Tags").Where(&Tag{Name: "go"}).First(t)
```

## Clear relations and/or delete them

Gorm can get pretty obscure when you're working with relations. What's the
difference between `Preload()`, `Related()` and `Association()` for example ?
It's still not clear to me. Anyway, let's say you want to remove all the tags
from an article. How would you do that ? 

```go
// Let's say you already have your article
a := &Article{}
db.First(a, 1)

// Let's clear the associations
db.Model(a).Association("Tags").Clear()
```

That should do the trick right ? Well kind of. Now in your many to many table
linking your `articles` and `tags` table, you have some tags that are not
associated with anything, the foreign key is simply set to `NULL`.

In this case it's not bad. Maybe at some point you'll have articles linked to
those tags. In other cases, you might just want to completely delete the 
associated object and there's no easy way to do this. You're going to have to
delete those with a classic query like so :

```go
db.Unscoped().Where(Tag{ArticleID: a.ID}).Delete(&Tag{})
```

The `Unscoped()` method here is used to tell gorm to not just only soft-delete
(set a delete date), but to completely remove the records from the database.

**Warning** : This applies to a few specific case. In many-to-many relationships
it's often better to just leave those with `NULL` values.

# Struct tips

## Using gorm.Model gotchas

That's something most people will not care about, but if you need to marshal
your struct to JSON, using `gorm.Model` will create unwanted fields. That
problem is usually addressed by simply creating another struct that contains
only the data you want to expose using JSON tags. But you can also not embed
the `gorm.Model` and directly integrate those fields in your struct with `-` 
JSON tags.

```go
type Tag struct {
	ID        uint       `json:"-" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	Name     string

	Articles []Article `gorm:"many2many:article_tags;"`
}
```

This struct will behave exactly as if you used `gorm.Model`, but the four fields
that are in `gorm.Model` won't be marshalled to JSON.

## JSONB for PostgreSQL

So you want to be able to store JSON and enjoy all the possibilities that PG has
to offer. Since a few month, gorm knows how to handle that for you, even though
you'll need to handle some things yourself. The use case for this is typically
if you have a JSON that is too complex to flatten and store in the database
in a single row, or that JSON is too dynamic for example.

```go

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type MetaData struct {
	// A lot of fields with JSON bindings and stuff
}

type Article struct {
	gorm.Model
	Title string
	Slug  string `gorm:"unique_index"`
	Body  string

	MetaDataB postgres.Jsonb `gorm:"type:jsonb;"`
	MetaData  MetaData       `gorm:"-"`

	Tags  []Tag `gorm:"many2many:article_tags;"`
}
```

Here we're telling gorm not to save directly the `MetaData` object, but to store
the `MetaDataB` object which is of type `postgres.Jsonb`. Under the hood, this
`Jsonb` type is actually just a `json.RawMessage` (which is essentially, just
a `[]byte`). 

```go
func (a *Article) MarshalMetaData() error {
	var err error

	a.MetaDataB, err = json.Marshal(a.MetaData)
	
	return err
}

func (a *Article) UnmarshalMetaData() error {
	return json.Unmarshal(a.MetaDataB.RawMessage, a.MetaData)
}
```

Now you can marshal before saving your article, and unmarshal once you retrieve
it from the database !

## Using UUID as Primary Key

If you wish to use the UUID extension for postgres, this section is for you.
First we'll declare our model with a slight change, the ID won't be an `uint`
anymore. Se we can't use `gorm.Model` here:

```go
import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Email string
}
```

Now if you try to create this model using a migration (or `CreateTable`) you'll
be faced with the issue that the `uuid-ossp` doesn't exist in your database and
is required. So the simple solution is to create this extension on the database
you're using:

```sql
> CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;
```

Alright now you're good to go. But what if we had an initial database migration
that checked for that extension before creating tables and doing a lot of
other stuff?

Here's an example on how to do it:

```go
import (
	"errors"

	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

var precheck = &gormigrate.Migration{
	ID: "precheck",
	Migrate: func(tx *gorm.DB) error {
		var (
			requiredext = "uuid-ossp"
			count       int
		)
		if err := tx.Table("pg_extension").Where("extname = ?", requiredext).Count(&count).Error; err != nil {
			return err
		}
		if count < 1 {
			return errors.New("extension uuid-ossp doesn't exist in the target database but is required")
		}
		return nil
	},
}
```

Run this migration first, it will return an error this extension isn't found
thus preventing the following migrations to run.

# Troubleshooting

## Relations are not created using CreateTable

The examples in the [documentation](http://jinzhu.me/gorm/associations.html) are
great but it's easy to omit things when reading it.

```go
type User struct {
    gorm.Model
    CreditCard CreditCard
}

type CreditCard struct {
    gorm.Model
    UserID uint
    Number string
}
```

For example it's easy to omit that there is actually a `UserID uint` field.
And that field is required by gorm to understand that this is actually a
relation. Of course there could be other issues that could prevent gorm to
create the relationship between two models.