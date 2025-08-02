package db

import (
    "fmt"
	"os"
    "log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?charset=utf8&parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_NAME"),
	)

	// Connect to the database
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
    
    // Migrate the database scheme
    err := DB.AutoMigrate(&Config{}, &Result{})
    if err := nil {
        log.Fatalf("failed to create or update tables according to the defined schemes: %v", err)
    }
}