package models

import (
	"time"

	"github.com/lib/pq"

	uuid "github.com/satori/go.uuid"
)

// Product is our main product type
type Product struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name  string
	Price int
	Tags  pq.StringArray `gorm:"type:varchar(100)[]"`
}
