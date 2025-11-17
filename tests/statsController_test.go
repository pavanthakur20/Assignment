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

func TestGetUserStatsSuccess(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/stats/:userId", controllers.GetUserStats)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	rewards := []models.StockReward{
		{
			ID:                 "stat-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "stat-2",
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

	req, _ := http.NewRequest("GET", "/api/stats/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user123", response.UserID)
	assert.NotNil(t, response.TodayRewards)
	assert.Greater(t, response.CurrentPortfolioINR, 0.0)
	assert.Greater(t, response.TotalSharesRewarded, 0.0)
}

func TestGetUserStatsNoData(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/stats/:userId", controllers.GetUserStats)

	req, _ := http.NewRequest("GET", "/api/stats/user999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user999", response.UserID)
	assert.NotNil(t, response.TodayRewards)
	assert.Equal(t, 0, len(response.TodayRewards))
	assert.Equal(t, 0.0, response.CurrentPortfolioINR)
	assert.Equal(t, 0.0, response.TotalSharesRewarded)
}

func TestGetUserStatsTodayRewardsOnly(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/stats/:userId", controllers.GetUserStats)

	now := time.Now()

	rewards := []models.StockReward{
		{
			ID:                 "stat-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "stat-2",
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

	req, _ := http.NewRequest("GET", "/api/stats/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response.TodayRewards, "RELIANCE")
	assert.Contains(t, response.TodayRewards, "TCS")
	assert.Equal(t, 10.0, response.TodayRewards["RELIANCE"])
	assert.Equal(t, 5.0, response.TodayRewards["TCS"])
}

func TestGetUserStatsMixedRewards(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/stats/:userId", controllers.GetUserStats)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	rewards := []models.StockReward{
		{
			ID:                 "stat-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "stat-2",
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

	req, _ := http.NewRequest("GET", "/api/stats/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.StatsResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response.TodayRewards, "RELIANCE")
	assert.NotContains(t, response.TodayRewards, "TCS")

	assert.Greater(t, response.CurrentPortfolioINR, 0.0)

	assert.Equal(t, 15.0, response.TotalSharesRewarded)
}
