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

func TestGetUserPortfolioSuccess(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/portfolio/:userId", controllers.GetUserPortfolio)

	now := time.Now()
	rewards := []models.StockReward{
		{
			ID:                 "port-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.5,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "port-2",
			UserID:             "user123",
			StockSymbol:        "TCS",
			Quantity:           5.25,
			RewardTimestamp:    now,
			StockPriceAtReward: 3500.0,
		},
		{
			ID:                 "port-3",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           2.5,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
	}

	for _, reward := range rewards {
		err := db.Create(&reward).Error
		assert.NoError(t, err)
	}

	req, _ := http.NewRequest("GET", "/api/portfolio/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PortfolioResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user123", response.UserID)
	assert.NotNil(t, response.Holdings)
	assert.Greater(t, response.TotalValue, 0.0)
	assert.NotZero(t, response.LastUpdated)

	var relianceHolding *models.UserStockHolding
	for _, holding := range response.Holdings {
		if holding.StockSymbol == "RELIANCE" {
			relianceHolding = &holding
			break
		}
	}
	assert.NotNil(t, relianceHolding)
	assert.Equal(t, 13.0, relianceHolding.TotalQuantity)
}

func TestGetUserPortfolioNoData(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/portfolio/:userId", controllers.GetUserPortfolio)

	req, _ := http.NewRequest("GET", "/api/portfolio/user999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PortfolioResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "user999", response.UserID)
	assert.NotNil(t, response.Holdings)
	assert.Equal(t, 0, len(response.Holdings))
	assert.Equal(t, 0.0, response.TotalValue)
}

func TestGetUserPortfolioMultipleStocks(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.GET("/api/portfolio/:userId", controllers.GetUserPortfolio)

	now := time.Now()
	rewards := []models.StockReward{
		{
			ID:                 "port-1",
			UserID:             "user123",
			StockSymbol:        "RELIANCE",
			Quantity:           10.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 2500.0,
		},
		{
			ID:                 "port-2",
			UserID:             "user123",
			StockSymbol:        "TCS",
			Quantity:           5.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 3500.0,
		},
		{
			ID:                 "port-3",
			UserID:             "user123",
			StockSymbol:        "INFOSYS",
			Quantity:           8.0,
			RewardTimestamp:    now,
			StockPriceAtReward: 1600.0,
		},
	}

	for _, reward := range rewards {
		err := db.Create(&reward).Error
		assert.NoError(t, err)
	}

	req, _ := http.NewRequest("GET", "/api/portfolio/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PortfolioResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(response.Holdings))

	stockSymbols := make(map[string]bool)
	for _, holding := range response.Holdings {
		stockSymbols[holding.StockSymbol] = true
		assert.Greater(t, holding.CurrentValue, 0.0)
		assert.Greater(t, holding.TotalQuantity, 0.0)
		assert.Greater(t, holding.CurrentPrice, 0.0)
	}
	assert.True(t, stockSymbols["RELIANCE"])
	assert.True(t, stockSymbols["TCS"])
	assert.True(t, stockSymbols["INFOSYS"])
}
