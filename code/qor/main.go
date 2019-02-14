package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"

	"github.com/Depado/articles/code/qor/admin"
	"github.com/Depado/articles/code/qor/migrate"
)

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
	if err = migrate.Start(db); err != nil {
		logrus.WithError(err).Fatal("Couldn't run migration")
	}

	r := gin.New()
	a := admin.New(db, "", "secret")
	a.Bind(r)
	r.Run("127.0.0.1:8080")
}
