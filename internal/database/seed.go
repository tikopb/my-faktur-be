package database

import (
	"bemyfaktur/internal/model"

	"gorm.io/gorm"
)

func seedDB(db *gorm.DB) {
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
