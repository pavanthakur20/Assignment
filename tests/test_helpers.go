package tests

import (
	"assignment/initializers"
	"assignment/models"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.StockReward{}, &models.LedgerEntry{})
	assert.NoError(t, err)

	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func setupTestLogger() {
	if initializers.Log == nil {
		initializers.Log = logrus.New()
		initializers.Log.SetOutput(os.Stdout)
		initializers.Log.SetLevel(logrus.WarnLevel)
	}
}
