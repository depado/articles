package resources

import (
	"github.com/qor/admin"

	"github.com/Depado/articles/code/qor/models"
)

// AddProduct adds the product model to the admin interface
func AddProduct(adm *admin.Admin) {
	adm.AddResource(&models.Product{})
}
