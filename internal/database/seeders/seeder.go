package seeders

import (
	"bemyfaktur/internal/database/fakers"
	"bemyfaktur/internal/model"
	"fmt"
	"github.com/google/uuid"
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
			UUID:        StringToUUID("3ebb6fa1-39ee-4f67-b482-180179cad78c"),
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
			UUID:           StringToUUID("673a0ff4-a6ac-4e67-888f-38a546c13325"),
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
			UUID:           StringToUUID("171e9123-bd3c-416b-98a9-4e20a57e3903"),
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
			UUID:           StringToUUID("f1a8feaa-f60b-46f2-b851-935ec44cf29f"),
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
			UUID:           StringToUUID("7dbc9309-a45e-4af6-8ee7-a30e6d11ca40"),
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
			UUID:              StringToUUID("633438e8-0d79-47d9-bcd5-a841ec2ee642"),
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
			UUID:           StringToUUID("e8a4cc6b-3578-4093-88b4-e190a23fc606"),
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
			UUID:           StringToUUID("e8a4cc6b-3578-4093-88b4-e190a23fc606"),
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
			UUID:           StringToUUID("63ad750a-d7a7-4784-a87a-e1c4d6a16990"),
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
			UUID:           StringToUUID("7fb873ca-5662-411d-a4f8-f9a926c6d449"),
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
			UUID:           StringToUUID("ade8f8f6-5a4c-4f4a-8f36-d92637c855a0"),
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
			UUID:           StringToUUID("a388273d-b5a5-49f1-956a-928c1f30cd26"),
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
		{Seeder: fakers.FileServiceFakers(db, model.FileService{
			ID:        1,
			UUID:      StringToUUID("63a10096-ab8b-4b78-9d94-0ce007eba6fe"),
			CreatedAt: time.Now(),
			CreatedBy: "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			UpdatedBy: "0bdfa34b-cdc4-4c4b-b087-f6b6d7cf81d5",
			FileName:  "main.jpeg",
			UuidDoc:   StringToUUID("633438e8-0d79-47d9-bcd5-a841ec2ee642"),
			DocType:   "INV",
		})},
	}
}

// run on db seed with function to delete constaint first for smooth seeder
func DBSeed(db *gorm.DB) error {

	// Run seeders
	for _, seeder := range RegisterSeeders(db) {
		err := db.Debug().Create(seeder.Seeder).Error
		if err != nil {
			return err
		}
	}

	//change the squence to 10 after the data
	RunSequenceChange(db)

	return nil
}

// function to change all squence id for table to start at 10 to make room for seeder input
func RunSequenceChange(db *gorm.DB) {
	// List of sequences to update
	sequences := []string{
		"document_no_temps_id_seq",
		"file_services_id_seq",
		"invoices_id_seq",
		"invoice_lines_id_seq",
		"organizations_id_seq",
		"partners_id_seq",
		"payment_lines_id_seq",
		"payments_id_seq",
		"products_id_seq",
	}

	// Loop through the sequences and execute the setval query for each one
	for _, seq := range sequences {
		query := fmt.Sprintf("SELECT setval('%s', 100, true);", seq)
		db.Exec(query)
	}
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

func StringToUUID(id string) uuid.UUID {
	// Validate the input string format if necessary
	// ...

	uID, err := uuid.Parse(id)
	if err != nil {
		fmt.Println(err)
	}
	return uID
}
