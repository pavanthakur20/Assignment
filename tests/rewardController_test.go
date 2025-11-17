package tests

import (
	"assignment/controllers"
	"assignment/initializers"
	"assignment/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRewardUserSuccess(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.POST("/api/reward", controllers.RewardUser)

	rewardReq := models.RewardRequest{
		ID:              "reward-123",
		UserID:          "user123",
		StockSymbol:     "RELIANCE",
		Quantity:        10.5,
		RewardTimestamp: time.Now(),
	}

	body, _ := json.Marshal(rewardReq)
	req, _ := http.NewRequest("POST", "/api/reward", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.RewardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Stock reward recorded successfully", response.Message)
	assert.NotNil(t, response.Reward)
	assert.NotNil(t, response.CompanyCharges)
	assert.Greater(t, response.INRValue, 0.0)

	var savedReward models.StockReward
	err = db.Where("id = ?", "reward-123").First(&savedReward).Error
	assert.NoError(t, err)
	assert.Equal(t, "user123", savedReward.UserID)
	assert.Equal(t, "RELIANCE", savedReward.StockSymbol)
	assert.Equal(t, 10.5, savedReward.Quantity)

	var ledgerEntries []models.LedgerEntry
	err = db.Where("reward_id = ?", "reward-123").Find(&ledgerEntries).Error
	assert.NoError(t, err)
	assert.Equal(t, 5, len(ledgerEntries))
}

func TestRewardUserInvalidRequest(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.POST("/api/reward", controllers.RewardUser)

	req, _ := http.NewRequest("POST", "/api/reward", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response models.RewardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}

func TestRewardUserDuplicateID(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.POST("/api/reward", controllers.RewardUser)

	existingReward := models.StockReward{
		ID:                 "reward-123",
		UserID:             "user123",
		StockSymbol:        "RELIANCE",
		Quantity:           10.0,
		RewardTimestamp:    time.Now(),
		StockPriceAtReward: 2500.0,
	}
	err := db.Create(&existingReward).Error
	assert.NoError(t, err)

	rewardReq := models.RewardRequest{
		ID:              "reward-123",
		UserID:          "user456",
		StockSymbol:     "TCS",
		Quantity:        5.0,
		RewardTimestamp: time.Now(),
	}

	body, _ := json.Marshal(rewardReq)
	req, _ := http.NewRequest("POST", "/api/reward", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response models.RewardResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Message, "already been processed")
}

func TestRewardUserCompanyChargesCalculation(t *testing.T) {
	db := setupTestDB(t)
	initializers.DB = db
	setupTestLogger()

	router := setupRouter()
	router.POST("/api/reward", controllers.RewardUser)

	rewardReq := models.RewardRequest{
		ID:              "reward-123",
		UserID:          "user123",
		StockSymbol:     "RELIANCE",
		Quantity:        10.0,
		RewardTimestamp: time.Now(),
	}

	body, _ := json.Marshal(rewardReq)
	req, _ := http.NewRequest("POST", "/api/reward", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.RewardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.CompanyCharges)

	charges := response.CompanyCharges
	assert.Greater(t, charges.StockCost, 0.0)
	assert.Greater(t, charges.Brokerage, 0.0)
	assert.Greater(t, charges.STT, 0.0)
	assert.Greater(t, charges.GST, 0.0)
	assert.Greater(t, charges.TotalCost, charges.StockCost)

	expectedBrokerage := charges.StockCost * 0.0003
	assert.InDelta(t, expectedBrokerage, charges.Brokerage, 1.0)

	expectedSTT := charges.StockCost * 0.001
	assert.InDelta(t, expectedSTT, charges.STT, 1.0)

	expectedGST := charges.Brokerage * 0.18
	assert.InDelta(t, expectedGST, charges.GST, 1.0)
}
