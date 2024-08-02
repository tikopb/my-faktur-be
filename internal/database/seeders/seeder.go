package seeders

import (
	"bemyfaktur/internal/database/fakers"
	"bemyfaktur/internal/model"
	"time"

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
		{Seeder: fakers.OrgFaker(db, model.Organization{
			ID:          1,
			CreatedAt:   time.Now(),
			UpdateAt:    time.Now(),
			CreatedBy:   "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy:   "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			OrgCode:     "IKE",
			Name:        "IKEA-ORG",
			Description: "Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.",
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
			OrganizationId: 1,
		})},

		//partners
		{Seeder: fakers.PartnerFaker(db, model.Partner{
			ID:             1,
			Name:           "IKEA-ID",
			CreatedBy:      "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			DNAmount:       0,
			CNAmount:       0,
			Isactive:       true,
			Code:           "BP-1",
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
			BatchNo:           "DUMP-1.1",
			Status:            "CO",
			DocAction:         "CO",
			OustandingPayment: 150750,
			DocumentNo:        "INV-001-DEC-2024",
			IsPrecentage:      false,
			PayDate:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), // Replace with the appropriate field faker function for date
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
			OrganizationId: 1,
		})},
	}
}

//func DBSeed(db *gorm.DB) error {
//	for _, seeder := range RegisterSeeders(db) {
//		err := db.Debug().Create(seeder.Seeder).Error
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func DBSeed(db *gorm.DB) error {

	// Disable foreign key constraints
	//db.Exec("ALTER TABLE users DROP CONSTRAINT fk_users_organization;")
	db.Exec("ALTER TABLE organizations DROP CONSTRAINT fk_organizations_user;")
	db.Exec("ALTER TABLE organizations DROP CONSTRAINT fk_organizations_user_updated;")

	// Run seeders
	for _, seeder := range RegisterSeeders(db) {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}

	// Re-enable foreign key constraints
	db.Exec("ALTER TABLE public.users ADD CONSTRAINT fk_users_organization FOREIGN KEY (org_id) REFERENCES public.organizations(id);")
	db.Exec("ALTER TABLE public.organizations ADD CONSTRAINT fk_organizations_user FOREIGN KEY (created_by) REFERENCES public.users(id);")
	db.Exec("ALTER TABLE public.organizations ADD CONSTRAINT fk_organizations_user_updated FOREIGN KEY (updated_by) REFERENCES public.users(id);")

	//change the scuence to 10 after the data
	RunSequenceChange(db)

	return nil
}

func RunSequenceChange(db *gorm.DB) {
	sql := `
		SELECT setval('document_no_temps_id_seq', 10, true);
		SELECT setval('file_services_id_seq', 10, true);
		SELECT setval('invoice_lines_id_seq', 10, true);
		SELECT setval('invoices_id_seq', 10, true);
		SELECT setval('organizations_id_seq', 10, true);
		SELECT setval('partners_id_seq', 10, true);
		SELECT setval('payment_lines_id_seq', 10, true);
		SELECT setval('payments_id_seq', 10, true);
		SELECT setval('products_id_seq', 10, true);
	`
	db.Exec(sql)
}

func MigrateDb(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.Partner{},
		&model.Invoice{},
		&model.InvoiceLine{},
		&model.Product{},
		&model.Payment{},
		&model.PaymentLine{},
		&model.DocumentNoTemp{},
		&model.FileService{},
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
