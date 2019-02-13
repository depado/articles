package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/Depado/articles/code/qor/v1/migrate"
	"github.com/Depado/articles/code/qor/v1/models"
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

	adm.AddResource(&models.Product{})
	usr := adm.AddResource(&models.AdminUser{}, &admin.Config{Menu: []string{"User Management"}})
	usr.IndexAttrs("-Password")
	usr.Meta(&admin.Meta{
		Name: "Password",
		Type: "password",
		Setter: func(resource interface{}, metaValue *resource.MetaValue, context *qor.Context) {
			values := metaValue.Value.([]string)
			if len(values) > 0 {
				if np := values[0]; np != "" {
					pwd, err := bcrypt.GenerateFromPassword([]byte(np), bcrypt.DefaultCost)
					if err != nil {
						context.DB.AddError(validations.NewError(usr, "Password", "Can't encrypt password")) // nolint: gosec,errcheck
						return
					}
					u := resource.(*models.AdminUser)
					u.Password = pwd
				}
			}
		},
	})

	r := gin.New()
	r.Any("/admin/*resources", gin.WrapH(mux))
	r.Run("127.0.0.1:8080")
}
