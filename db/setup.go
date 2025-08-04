package db

import (
    "fmt"
	"os"
    "log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(WarningLog *log.Logger, InfoLog *log.Logger) error {
	InfoLog.Println("Initializing Database...")
    dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?charset=utf8&parseTime=True",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_URL"),
		os.Getenv("DB_NAME"),
	)

	// Connect to the database
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
    defer func() {
        if cerr = DB.Close(); cerr != nil {
            WarningLog.Printf("failed to close database: %v \n", cerr)
        }
    }()
    
    // Migrate the database scheme
    InfoLog.Println("Migrating Database Scheme...")
    err := DB.AutoMigrate(&Config{}, &Result{})
    if err := nil {
        return fmt.Errorf("failed to create or update tables according to the defined schemes: %w", err)
    }
    
    InfoLog.Println("Database ready!")
    return nil
}
