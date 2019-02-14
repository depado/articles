package migrate

import (
	"errors"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var uuidCheck = &gormigrate.Migration{
	ID: "precheck",
	Migrate: func(tx *gorm.DB) error {
		var (
			requiredext = "uuid-ossp"
			count       int
		)
		if err := tx.Table("pg_extension").Where("extname = ?", requiredext).Count(&count).Error; err != nil {
			return err
		}
		if count < 1 {
			return errors.New("extension uuid-ossp doesn't exist in the target database but is required")
		}
		return nil
	},
}
