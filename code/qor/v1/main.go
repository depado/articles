package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
	"github.com/sirupsen/logrus"

	"github.com/Depado/articles/code/qor/v1/migrate"
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

	adm := admin.New(&admin.AdminConfig{SiteName: "Admin", DB: db})
	mux := http.NewServeMux()
	adm.MountTo("/admin", mux)

	r := gin.New()
	r.Any("/admin/*resources", gin.WrapH(mux))
	r.Run("127.0.0.1:8080")
}
