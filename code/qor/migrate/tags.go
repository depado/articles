package migrate

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"gopkg.in/gormigrate.v1"
)

var productTags = &gormigrate.Migration{
	ID: "product_tags",
	Migrate: func(tx *gorm.DB) error {
		type product struct {
			ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
			CreatedAt time.Time
			UpdatedAt time.Time
			DeletedAt *time.Time
			Name      string
			Price     int
			Tags      pq.StringArray `gorm:"type:varchar(100)[]"`
		}
		return tx.AutoMigrate(&product{}).Error
	}, Rollback: func(tx *gorm.DB) error {
		type product struct {
			ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
			CreatedAt time.Time
			UpdatedAt time.Time
			DeletedAt *time.Time
			Name      string
			Price     int
			Tags      pq.StringArray `gorm:"type:varchar(100)[]"`
		}
		return tx.Model(&product{}).DropColumn("tags").Error
	},
}
