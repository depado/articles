package migrate

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"

	"github.com/Depado/articles/code/qor/admin"
)

// Start starts the migration process
func Start(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		uuidCheck,
		initial,
		admin.AdminUserMigration,
		productTags,
	})
	return m.Migrate()
}
