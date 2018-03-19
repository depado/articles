package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"

	"github.com/Depado/articles/gorm/migrate"
)

func main() {
	db, err := gorm.Open("sqlite3", "gorm.db")
	if err != nil {
		logrus.WithError(err).Fatal("Couldn't open db")
	}
	defer db.Close()
	db.LogMode(true)

	migrate.Start(db)

	a := &migrate.Article{}
	db.Preload("Tags").First(a, 1)
	spew.Dump(a)

	// Or find all the articles that are tagged with "go" ?
	t := &migrate.Tag{}
	db.Preload("Articles").Preload("Articles.Tags").Where(&migrate.Tag{Name: "go"}).First(t)
	spew.Dump(t)
}
