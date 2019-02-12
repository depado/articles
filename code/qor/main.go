package main

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type product struct {
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name  string
	Price int
}

func main() {
	db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
	if err != nil {
		logrus.WithError(err).Fatal("Couldn't initialize database connection")
	}
	defer db.Close()
}
