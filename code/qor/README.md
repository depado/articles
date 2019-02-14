# QOR Admin Example Code

This directory holds the code for the 
["QOR Admin Tutorial"](https://blog.depa.do/post/qor-admin-tutorial) article.

Most of the useful code is located in the 
[admin](https://github.com/Depado/articles/tree/master/code/qor/admin) package.

- Declaring a QOR Admin easily with a simple constructor
- Bind to an existing gin router and specify a prefix for all the admin related
  tasks (such as login, logout and admin interface itself)
- Initial migration and exported migration from the admin package to create the
  `admin_users` table
- Manage admin users with a special setter for their password and 
  [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt) encryption in the 
  database
- Authentication handled by [Gin](http://github.com/gin-gonic/gin) with 
  [gin-contrib/sessions](https://github.com/gin-contrib/sessions) and a cookie 
  backend
- Separating admin-related structs (like admin user) and business structs (the 
  [models](https://github.com/Depado/articles/tree/master/code/qor/models) 
  directory only has business structures)
- Simple login HTML template using 
  [spectre.css](https://picturepan2.github.io/spectre/index.html)
