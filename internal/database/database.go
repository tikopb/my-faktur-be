package database

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetDb() *gorm.DB {
	dbAddress := getDbAdress()

	// var logMode logger.LogLevel
	// if os.Getenv("environment_site") == "production" {
	// 	logMode = logger.Silent
	// } else {
	// 	logMode = logger.Info // Set the desired log level for non-production environments
	// }

	db, err := gorm.Open(postgres.Open(dbAddress), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("Failed to connect into database")
	}

	seedDB(db)

	return db
}

func getDbAdress() string {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic("config environment not found!")
	}

	config := DBConfig{
		Host:     viper.GetString("db_host"),
		Port:     viper.GetString("db_port"),
		User:     viper.GetString("db_user"),
		Password: viper.GetString("db_password"),
		DBName:   viper.GetString("db_dbname"),
		SSLMode:  viper.GetString("db_sslmode"),
	}

	dbAddress := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
	return dbAddress
}
