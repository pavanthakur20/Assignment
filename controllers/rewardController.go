package controllers

import (
	"assignment/initializers"
	"assignment/models"
	"assignment/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RewardUser(c *gin.Context) {
	log := initializers.Log
	var req models.RewardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Invalid reward request")
		c.JSON(http.StatusBadRequest, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	var existingReward models.StockReward
	err := initializers.DB.Where("id = ?", req.ID).First(&existingReward).Error

	if err == nil {
		log.WithField("reward_id", req.ID).Warn("Duplicate reward ID")
		c.JSON(http.StatusConflict, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Reward with ID '%s' has already been processed", req.ID),
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		log.WithError(err).WithField("reward_id", req.ID).Error("Database error checking reward")
		c.JSON(http.StatusInternalServerError, models.RewardResponse{
			Success: false,
			Message: fmt.Sprintf("Database error: %v", err),
		})
		return
	}

	stockPrice, err := services.GetCurrentStockPrice(req.StockSymbol)
	if err != nil {
		log.WithError(err).WithField("stock_symbol", req.StockSymbol).Error("Failed to get stock price")
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
		log.WithError(err).WithField("reward_id", req.ID).Error("Failed to record reward")
		c.JSON(http.StatusInternalServerError, models.RewardResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	log.WithFields(logrus.Fields{
		"reward_id":  reward.ID,
		"user_id":    reward.UserID,
		"symbol":     reward.StockSymbol,
		"quantity":   reward.Quantity,
		"total_cost": charges.TotalCost,
	}).Info("Reward recorded successfully")

	c.JSON(http.StatusCreated, models.RewardResponse{
		Success:        true,
		Message:        "Stock reward recorded successfully",
		Reward:         reward,
		INRValue:       stockCost,
		CompanyCharges: charges,
	})
}
