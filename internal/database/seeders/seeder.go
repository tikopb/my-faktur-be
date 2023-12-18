package seeders

import (
	"bemyfaktur/internal/database/fakers"
	"bemyfaktur/internal/model"

	"gorm.io/gorm"
)

type Seeder struct {
	Seeder interface{}
}

func RegisterSeeders(db *gorm.DB) []Seeder {
	return []Seeder{
		{Seeder: fakers.PartnerFaker(db)},
		{Seeder: fakers.PartnerFaker(db)},
		{Seeder: fakers.ProductFaker(db)},
	}
}

func DBSeed(db *gorm.DB) error {
	for _, seeder := range RegisterSeeders(db) {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func MigrateDb(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(
		&model.User{},
		&model.Partner{},
		&model.Invoice{},
		&model.InvoiceLine{},
		&model.Product{},
		&model.Payment{},
		&model.PaymentLine{},
		&model.DocumentNoTemp{},
	)
}
