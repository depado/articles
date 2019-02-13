package admin

import (
	"github.com/Depado/articles/code/qor/v1/models"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
)

// Auth is a structure to handle authentication for QOR. It will satisify the
// qor.Auth interface.
type Auth struct {
	db    *gorm.DB
	store cookie.Store

	login   string
	logout  string
	session string
	key     string
}

// GetCurrentUser satisfies the Auth interface and returns the current user
func (a Auth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	var userid uint
	s, err := a.store.Get(c.Request, a.session)
	if err != nil {
		return nil
	}
	if v, ok := s.Values[a.key]; ok {
		userid = v.(uint)
	} else {
		return nil
	}

	var user models.AdminUser
	if !a.db.First(&user, "id = ?", userid).RecordNotFound() {
		return &user
	}

	return nil
}

// LoginURL statisfies the Auth interface and returns the route used to log
// users in
func (a Auth) LoginURL(c *admin.Context) string { // nolint: unparam
	return "/login"
}

// LogoutURL statisfies the Auth interface and returns the route used to logout
// a user
func (a Auth) LogoutURL(c *admin.Context) string { // nolint: unparam
	return "/logout"
}
