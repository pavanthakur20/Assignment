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

func TestGetTodayStocksSuccess(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/today/:userId/stocks", controllers.GetTodayStocks)

	now := time.Now()
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
			ID:                 "today-2",
			UserID:             "user123",
			StockSymbol:        "TCS",
			Quantity:           5.0,
			RewardTimestamp:    now.Add(2 * time.Hour),
			StockPriceAtReward: 3500.0,
		},
	}

	for _, reward := range rewards {
		err := db.Create(&reward).Error
		assert.NoError(t, err)
	}

	req, _ := http.NewRequest("GET", "/api/today/user123/stocks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user123", response["user_id"])
	assert.NotNil(t, response["date"])
	assert.Equal(t, float64(2), response["total_rewards"])
	assert.NotNil(t, response["rewards"])
}

func TestGetTodayStocksNoData(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/today/:userId/stocks", controllers.GetTodayStocks)

	req, _ := http.NewRequest("GET", "/api/today/user999/stocks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user999", response["user_id"])
	assert.Equal(t, float64(0), response["total_rewards"])
}

func TestGetTodayStocksExcludesHistoricalData(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/today/:userId/stocks", controllers.GetTodayStocks)

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

	req, _ := http.NewRequest("GET", "/api/today/user123/stocks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), response["total_rewards"])

	rewardsArray := response["rewards"].([]interface{})
	assert.Equal(t, 1, len(rewardsArray))
}
