package seeders

import (
	"bemyfaktur/internal/database/fakers"
	"bemyfaktur/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Seeder struct {
	Seeder interface{}
}

func RegisterSeeders(db *gorm.DB) []Seeder {
	return []Seeder{
		//users
		{Seeder: fakers.UserFacker(db)},

		//organization
		{Seeder: fakers.OrganizationFakers(db, model.Organization{
			ID:          1,
			UUID:        StringParsingToUUID("11b98639-986d-4d55-857a-e0167e80a968"),
			CreatedAt:   time.Now(),
			UpdateAt:    time.Now(),
			CreatedBy:   "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:   "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			OrgCode:     "IK",
			Name:        "Ikea-Group",
			Description: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book.",
			IsActive:    true,
		})},

		//products
		{Seeder: fakers.ProductFaker(db, model.Product{
			ID:             1,
			Name:           "MARKUS",
			Value:          "C-1.1",
			Upc:            "502.611.51",
			Description:    "Kursi kantor ergonomis ini membuat Anda tetap nyaman dan fokus dengan fitur-fitur seperti ketegangan kemiringan yang dapat diatur secara manual, dan sandaran kepala/lengan untuk membantu mengendurkan otot-otot di leher dan punggung. Garansi 10 tahun.",
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			IsActive:       true,
			UUID:           StringParsingToUUID("c2a063aa-eba5-4dbe-b32a-14e1c1be6758"),
			OrganizationId: 1,
		})},
		{Seeder: fakers.ProductFaker(db, model.Product{
			ID:             2,
			Name:           "TROTTEN",
			Value:          "D-1.1",
			Upc:            "794.295.79",
			Description:    "Berganti posisi dari duduk ke berdiri sangat baik untuk Anda dan gagang putar memungkinkan Anda untuk menggerakkan lengan sekaligus menyesuaikan ketinggian. Menggerakkan tubuh membuat Anda bekerja lebih baik.",
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			IsActive:       true,
			UUID:           StringParsingToUUID("42c4fbaa-17e0-4cb7-b5c8-388fbe5ddf11"),
			OrganizationId: 1,
		})},
		{Seeder: fakers.ProductFaker(db, model.Product{
			ID:             3,
			Name:           "UPPSPEL",
			Value:          "T-1.1",
			Upc:            "704.998.40",
			Description:    "Garansi 10 tahun. Baca lebih lanjut mengenai syarat dan ketentuan di dalam brosur garansi. Unit ini dapat diletakkan di mana saja di dalam ruangan karena bagian belakang diberi sentuhan akhir. Roda memudahkan menggeser penyimpanan di bawah meja atau sekitar ruangan. Penyimpanan dapat dikunci untuk barang pribadi Anda.",
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			IsActive:       false,
			UUID:           StringParsingToUUID("0bb233e1-04ed-4d59-a322-ad4fa621d528"),
			OrganizationId: 1,
		})},

		//partners
		{Seeder: fakers.PartnerFaker(db, model.Partner{
			ID:             1,
			Name:           "IKEA-SUPP-ID",
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			DNAmount:       0,
			CNAmount:       0,
			Isactive:       true,
			Code:           "IK-BP-1",
			UUID:           StringParsingToUUID("c734ad43-15da-4f8b-9260-543a517dce9c"),
			OrganizationId: 1,
		})},
		{Seeder: fakers.PartnerFaker(db, model.Partner{
			ID:             2,
			Name:           "IKEA-SUPP-MY",
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			DNAmount:       0,
			CNAmount:       0,
			Isactive:       false,
			Code:           "IK-BP-2",
			UUID:           StringParsingToUUID("96285a1a-257a-4a29-bf73-68f76bf8643a"),
			OrganizationId: 1,
		})},

		//invoice
		{Seeder: fakers.InvoiceFaker(db, model.Invoice{
			ID:                1,
			CreatedBy:         "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:         "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			PartnerID:         1,
			GrandTotal:        150750,
			Discount:          0,
			BatchNo:           "IK-1.1",
			Status:            "CO",
			DocAction:         "CO",
			OustandingPayment: 150750,
			DocumentNo:        "INV-001-JAN-2025",
			IsPrecentage:      false,
			PayDate:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Replace with the appropriate field faker function for date
			UUID:              StringParsingToUUID("cef1eac1-b36b-4df0-84be-b40c229a996b"),
			OrganizationId:    1,
		})},

		//invoiceline
		{Seeder: fakers.InvoiceLineFaker(db, model.InvoiceLine{
			ID:             1,
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			Price:          5000,
			Discount:       100,
			Qty:            1,
			Amount:         4500,
			IsPrecentage:   false,
			ProductID:      1,
			InvoiceID:      1,
			UUID:           StringParsingToUUID("9596d041-8504-4c7b-af62-88e6ade36486"),
			OrganizationId: 1,
		})},
		{Seeder: fakers.InvoiceLineFaker(db, model.InvoiceLine{
			ID:             2,
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			Price:          150000,
			Discount:       2.5,
			Qty:            1,
			Amount:         146250,
			IsPrecentage:   true,
			ProductID:      2,
			InvoiceID:      1,
			UUID:           StringParsingToUUID("c2897274-60d1-4ef4-864e-9ba13abc22c9"),
			OrganizationId: 1,
		})},

		//payment
		{Seeder: fakers.PaymentFaker(db, model.Payment{
			ID:             1,
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			PartnerID:      1,
			GrandTotal:     145750,
			Discount:       0,
			BatchNo:        "PAY-DUMP-1.1",
			Status:         "DR",
			DocAction:      "DR",
			DocumentNo:     "PAY-001-DEC-2024",
			IsPrecentage:   false,
			UUID:           StringParsingToUUID("82de82e2-aa89-4098-a37f-b92f6d66e943"),
			OrganizationId: 1,
		})},
		{Seeder: fakers.PaymentFaker(db, model.Payment{
			ID:             2,
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			PartnerID:      1,
			GrandTotal:     5000,
			Discount:       0,
			BatchNo:        "PAY-DUMP-1.2",
			Status:         "DR",
			DocAction:      "DR",
			DocumentNo:     "PAY-002-DEC-2024",
			IsPrecentage:   false,
			UUID:           StringParsingToUUID("d7922a27-d012-4d12-adee-b47b0e8cf056"),
			OrganizationId: 1,
		})},

		//paymentline
		{Seeder: fakers.PaymentLineFaker(db, model.PaymentLine{
			ID:             1,
			PaymentID:      1,
			Price:          145750,
			Amount:         1,
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			InvoiceID:      1,
			Discount:       0,
			IsPrecentage:   false,
			UUID:           StringParsingToUUID("f924d2a4-5f2d-45ea-bf5c-927196529e6d"),
			OrganizationId: 1,
		})},
		{Seeder: fakers.PaymentLineFaker(db, model.PaymentLine{
			ID:             2,
			PaymentID:      2,
			Price:          5000,
			Amount:         1,
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			InvoiceID:      1,
			Discount:       0,
			IsPrecentage:   false,
			UUID:           StringParsingToUUID("a6fb839e-d14c-429a-a31e-fe70234acba9"),
			OrganizationId: 1,
		})},
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
		&model.FileService{},
		&model.Organization{},
	)
}

func CreateDb(db *gorm.DB) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("env of db_dbname not found")
	}

	value := viper.GetString("db_dbname")
	db.Exec("create database " + value + ";")
}

func StringParsingToUUID(value string) uuid.UUID {
	parsedUUID, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}

	return parsedUUID
}
