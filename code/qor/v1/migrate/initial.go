package migrate

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

var initial = &gormigrate.Migration{
	ID: "initial",
	Migrate: func(tx *gorm.DB) error {
		type product struct {
			ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
			CreatedAt time.Time
			UpdatedAt time.Time
			DeletedAt *time.Time

			Name  string
			Price int
		}
		return tx.CreateTable(&product{}).Error
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.DropTable("products").Error
	},
}
