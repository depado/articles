package migrate

import (
	"github.com/jinzhu/gorm"
	gormigrate "gopkg.in/gormigrate.v1"
)

type Article struct {
	gorm.Model
	Title    string
	Slug     string `gorm:"unique_index"`
	Body     string
	Author   Author
	AuthorID uint
	Tags     []Tag `gorm:"many2many:article_tags;"`
}

type Tag struct {
	gorm.Model
	Name     string
	Articles []Article `gorm:"many2many:article_tags;"`
}

type Author struct {
	gorm.Model
	Name string

	Articles []Article
}

// Start starts the migration process
func Start(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "initial",
			Migrate: func(tx *gorm.DB) error {
				type Article struct {
					gorm.Model
					Title  string
					Author Author
					Tags   []Tag `gorm:"many2many:article_tags;"`
				}

				type Tag struct {
					gorm.Model
					Name     string
					Articles []Article `gorm:"many2many:article_tags;"`
				}
				return tx.CreateTable(&Article{}, &Tag{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("articles", "tags").Error
			},
		}, {
			ID: "add fields",
			Migrate: func(tx *gorm.DB) error {
				type Article struct {
					gorm.Model
					Title string
					Slug  string `gorm:"unique_index"`
					Body  string
					Tags  []Tag `gorm:"many2many:article_tags;"`
				}
				return tx.AutoMigrate(&Article{}).Error
			},
		},
		{
			ID: "fixtures",
			Migrate: func(tx *gorm.DB) error {
				a := Article{
					Title: "Hello World",
					Slug:  "hello-world",
					Body:  "This is an hello world article",
					Tags: []Tag{
						{Name: "dev"},
						{Name: "go"},
					},
				}
				return tx.Save(&a).Error
			},
		}, {
			ID: "add fk",
			Migrate: func(tx *gorm.DB) error {
				type Article struct {
					gorm.Model
					Title    string
					Slug     string `gorm:"unique_index"`
					Body     string
					Author   Author
					AuthorID uint
					Tags     []Tag `gorm:"many2many:article_tags;"`
				}

				type Author struct {
					gorm.Model
					Name     string
					Articles []Article
				}
				tx.CreateTable(&Author{})
				return tx.Model(&Article{}).AddForeignKey("author_id", "authors(id)", "RESTRICT", "RESTRICT").Error
			},
		},
	})
	return m.Migrate()
}
