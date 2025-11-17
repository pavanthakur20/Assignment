package tests

import (
	"assignment/controllers"
	"assignment/initializers"
	"assignment/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHistoricalINRSuccess(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/historical/:userId/inr", controllers.GetHistoricalINR)

	yesterday := time.Now().AddDate(0, 0, -1)
	twoDaysAgo := time.Now().AddDate(0, 0, -2)

	rewards := []models.StockReward{
		{
			ID:                 "hist-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.0,
			RewardTimestamp:    yesterday,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "hist-2",
			UserID:             "user123",
			StockSymbol:        "TCS",
			Quantity:           5.0,
			RewardTimestamp:    twoDaysAgo,
			StockPriceAtReward: 3500.0,
		},
	}

	for _, reward := range rewards {
		err := db.Create(&reward).Error
		assert.NoError(t, err)
	}

	req, _ := http.NewRequest("GET", "/api/historical/user123/inr", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HistoricalINRResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user123", response.UserID)
	assert.NotNil(t, response.DailyValues)
}

func TestGetHistoricalINRNoData(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/historical/:userId/inr", controllers.GetHistoricalINR)

	req, _ := http.NewRequest("GET", "/api/historical/user999/inr", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HistoricalINRResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user999", response.UserID)
	assert.NotNil(t, response.DailyValues)
	assert.Equal(t, 0, len(response.DailyValues))
}

func TestGetHistoricalINRExcludesTodayRewards(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/historical/:userId/inr", controllers.GetHistoricalINR)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	rewards := []models.StockReward{
		{
			ID:                 "today-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "hist-1",
			UserID:             "user123",
			StockSymbol:        "TCS",
			Quantity:           5.0,
			RewardTimestamp:    yesterday,
			StockPriceAtReward: 3500.0,
		},
	}

	for _, reward := range rewards {
		err := db.Create(&reward).Error
		assert.NoError(t, err)
	}

	req, _ := http.NewRequest("GET", "/api/historical/user123/inr", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HistoricalINRResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	yesterdayStr := yesterday.Format("2006-01-02")
	todayStr := now.Format("2006-01-02")

	assert.Contains(t, response.DailyValues, yesterdayStr)
	assert.NotContains(t, response.DailyValues, todayStr)
}
