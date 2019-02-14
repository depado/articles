package admin

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/validations"
	"golang.org/x/crypto/bcrypt"
	gormigrate "gopkg.in/gormigrate.v1"
)

// User defines how an admin user is represented in database
type User struct {
	gorm.Model
	Email     string `gorm:"not null;unique"`
	FirstName string
	LastName  string
	Password  []byte
	LastLogin *time.Time
}

// TableName allows to override the name of the table
func (u User) TableName() string {
	return "admin_users"
}

// DisplayName satisfies the interface for Qor Admin
func (u User) DisplayName() string {
	if u.FirstName != "" && u.LastName != "" {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	}
	return u.Email
}

// HashPassword is a simple utility function to hash the password sent via API
// before inserting it in database
func (u *User) HashPassword() error {
	pwd, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

// CheckPassword is a simple utility function to check the password given as raw
// against the user's hashed password
func (u User) CheckPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(raw)) == nil
}

// AdminUserMigration is the migration that creates our user model
var AdminUserMigration = &gormigrate.Migration{
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
			Email:    "you@yourcompany.com",
			Password: pwd,
		}
		return tx.Save(&usr).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.DropTable("admin_users").Error
	},
}

func addUser(adm *admin.Admin) {
	usr := adm.AddResource(&User{}, &admin.Config{Menu: []string{"User Management"}})
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
					u := resource.(*User)
					u.Password = pwd
				}
			}
		},
	})
}
