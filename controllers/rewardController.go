package controllers

import (
	"assignment/initializers"
	"assignment/models"
	"assignment/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RewardUser(c *gin.Context) {
	var req models.RewardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	var existingReward models.StockReward
	err := initializers.DB.Where("id = ?", req.ID).First(&existingReward).Error

	if err == nil {
		c.JSON(http.StatusConflict, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Reward with ID '%s' has already been processed", req.ID),
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Database error: %v", err),
		})
		return
	}

	stockPrice, err := services.GetCurrentStockPrice(req.StockSymbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get stock price: %v", err),
		})
		return
	}

	stockCost := stockPrice * req.Quantity
	charges := services.CalculateCompanyCharges(stockCost)

	reward := &models.StockReward{
		ID:                 req.ID,
		UserID:             req.UserID,
		StockSymbol:        req.StockSymbol,
		Quantity:           req.Quantity,
		RewardTimestamp:    req.RewardTimestamp,
		StockPriceAtReward: stockPrice,
	}

	err = initializers.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reward).Error; err != nil {
			return fmt.Errorf("failed to record reward: %v", err)
		}

		if err := services.RecordLedgerEntriesGORM(tx, reward, charges); err != nil {
			return fmt.Errorf("failed to record ledger entries: %v", err)
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RewardResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.RewardResponse{
		Success:        true,
		Message:        "Stock reward recorded successfully",
		Reward:         reward,
		INRValue:       stockCost,
		CompanyCharges: charges,
	})
}
