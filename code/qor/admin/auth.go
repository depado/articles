package admin

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/sirupsen/logrus"
)

// Auth is a structure to handle authentication for QOR. It will satisify the
// qor.Auth interface.
type auth struct {
	db      *gorm.DB
	session sessionConfig
	paths   pathConfig
}

type sessionConfig struct {
	name  string
	key   string
	store cookie.Store
}

type pathConfig struct {
	login  string
	logout string
	admin  string
}

// GetLogin simply returns the login page
func (a *auth) GetLogin(c *gin.Context) {
	if sessions.Default(c).Get(a.session.key) != nil {
		c.Redirect(http.StatusSeeOther, a.paths.admin)
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

// PostLogin is the handler to check if the user can connect
func (a *auth) PostLogin(c *gin.Context) {
	session := sessions.Default(c)
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.Redirect(http.StatusSeeOther, a.paths.login)
		return
	}
	var u User
	if a.db.Where(&User{Email: email}).First(&u).RecordNotFound() {
		c.Redirect(http.StatusSeeOther, a.paths.login)
		return
	}
	if !u.CheckPassword(password) {
		c.Redirect(http.StatusSeeOther, a.paths.login)
		return
	}

	now := time.Now()
	u.LastLogin = &now
	a.db.Save(&u)

	session.Set(a.session.key, u.ID)
	err := session.Save()
	if err != nil {
		logrus.WithError(err).Warn("Couldn't save session")
		c.Redirect(http.StatusSeeOther, a.paths.login)
		return
	}
	c.Redirect(http.StatusSeeOther, a.paths.admin)
}

// GetLogout allows the user to disconnect
func (a *auth) GetLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(a.session.key)
	if err := session.Save(); err != nil {
		logrus.WithError(err).Warn("Couldn't save session")
	}
	c.Redirect(http.StatusSeeOther, a.paths.login)
}

// GetCurrentUser satisfies the Auth interface and returns the current user
func (a auth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	var userid uint

	s, err := a.session.store.Get(c.Request, a.session.name)
	if err != nil {
		return nil
	}
	if v, ok := s.Values[a.session.key]; ok {
		userid = v.(uint)
	} else {
		return nil
	}

	var user User
	if !a.db.First(&user, "id = ?", userid).RecordNotFound() {
		return &user
	}

	return nil
}

// LoginURL statisfies the Auth interface and returns the route used to log
// users in
func (a auth) LoginURL(c *admin.Context) string { // nolint: unparam
	return a.paths.login
}

// LogoutURL statisfies the Auth interface and returns the route used to logout
// a user
func (a auth) LogoutURL(c *admin.Context) string { // nolint: unparam
	return a.paths.logout
}
