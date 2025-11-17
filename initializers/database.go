package initializers

import (
	"assignment/models"
	"fmt"
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
		Log.Fatal("DB_URL environment variable is not set")
	}

	dbLogLevel := logger.Warn
	DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(dbLogLevel),
	})

	if err != nil {
		Log.WithError(err).Fatal("Unable to connect to database")
	}

	Log.Info("Connected to database successfully")

	err = runMigrations()
	if err != nil {
		Log.WithError(err).Fatal("Failed to run migrations")
	}
}

func runMigrations() error {
	Log.Info("Running database migrations...")
	err := DB.AutoMigrate(
		&models.StockReward{},
		&models.LedgerEntry{},
	)

	if err != nil {
		return fmt.Errorf("auto migration error: %v", err)
	}

	Log.Info("Database migrations completed successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			Log.Info("Database connection closed")
		}
	}
}
