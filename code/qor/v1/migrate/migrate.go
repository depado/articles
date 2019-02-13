package migrate

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

// Start starts the migration process
func Start(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		uuidCheck,
		initial,
	})
	return m.Migrate()
}
