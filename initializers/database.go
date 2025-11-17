package initializers

import (
	"assignment/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	fmt.Println("Connected to database successfully")

	err = runMigrations()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v\n", err)
	}
}

func runMigrations() error {
	err := DB.AutoMigrate(
		&models.StockReward{},
		&models.LedgerEntry{},
	)

	if err != nil {
		return fmt.Errorf("auto migration error: %v", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			fmt.Println("Database connection closed")
		}
	}
}
