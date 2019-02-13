package admin

import (
	"path/filepath"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"
	"golang.org/x/crypto/bcrypt"

	"github.com/Depado/articles/code/qor/v1/models"
)

type Admin struct {
	db        *gorm.DB
	auth      auth
	adm       *admin.Admin
	adminpath string
	prefix    string
}

func New(db *gorm.DB, prefix, cookiesecret string) *Admin {
	adminpath := filepath.Join(prefix, "/admin")
	a := Admin{
		db:        db,
		prefix:    prefix,
		adminpath: adminpath,
		auth: auth{
			db: db,
			paths: pathConfig{
				admin:  adminpath,
				login:  filepath.Join(prefix, "/login"),
				logout: filepath.Join(prefix, "/logout"),
			},
			session: sessionConfig{
				key:   "userid",
				name:  "admsession",
				store: cookie.NewStore([]byte(cookiesecret)),
			},
		},
	}
	a.adm = admin.New(&admin.AdminConfig{SiteName: "My Admin Interface", DB: db, Auth: a.auth})

	a.adm.AddResource(&models.Product{})
	usr := a.adm.AddResource(&models.AdminUser{}, &admin.Config{Menu: []string{"User Management"}})
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

	return &a
}

func (a Admin) Bind(r *gin.Engine) {

}
