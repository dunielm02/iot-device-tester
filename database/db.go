package database

import (
	"iotTester/models"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToSQLite() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(os.Getenv("DB_NAME")), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&models.Log{})
	if err != nil {
		return err
	}

	return nil
}
