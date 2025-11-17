package controllers

import (
	"assignment/initializers"
	"assignment/models"
	"assignment/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUserStats(c *gin.Context) {
	log := initializers.Log
	userID := c.Param("userId")

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var allRewards []models.StockReward
	err := initializers.DB.Where("user_id = ?", userID).
		Order("reward_timestamp DESC").
		Find(&allRewards).Error

	if err != nil {
		log.WithError(err).WithField("user_id", userID).Error("Failed to fetch user rewards")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch user rewards",
			"details": err.Error(),
		})
		return
	}

	if len(allRewards) == 0 {
		c.JSON(http.StatusOK, models.StatsResponse{
			UserID:              userID,
			TodayRewards:        make(map[string]float64),
			CurrentPortfolioINR: 0.0,
			TotalSharesRewarded: 0.0,
		})
		return
	}

	todayRewardsMap := make(map[string]float64)
	portfolioQuantities := make(map[string]float64)
	totalSharesRewarded := 0.0
	uniqueSymbols := make(map[string]bool)

	for _, reward := range allRewards {
		uniqueSymbols[reward.StockSymbol] = true
		totalSharesRewarded += reward.Quantity
		portfolioQuantities[reward.StockSymbol] += reward.Quantity
		if reward.RewardTimestamp.After(startOfDay) && reward.RewardTimestamp.Before(endOfDay) {
			todayRewardsMap[reward.StockSymbol] += reward.Quantity
		}
	}

	symbolsList := make([]string, 0, len(uniqueSymbols))
	for symbol := range uniqueSymbols {
		symbolsList = append(symbolsList, symbol)
	}

	prices, err := services.GetCurrentPrices(symbolsList)
	if err != nil {
		log.WithError(err).WithField("user_id", userID).Error("Failed to fetch current prices")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch current prices",
			"details": err.Error(),
		})
		return
	}

	totalPortfolioValue := 0.0
	for symbol, quantity := range portfolioQuantities {
		if price, ok := prices[symbol]; ok {
			totalPortfolioValue += quantity * price
		}
	}

	response := models.StatsResponse{
		UserID:              userID,
		TodayRewards:        todayRewardsMap,
		CurrentPortfolioINR: float64(int(totalPortfolioValue*100)) / 100,
		TotalSharesRewarded: float64(int(totalSharesRewarded*1000000)) / 1000000,
	}

	c.JSON(http.StatusOK, response)
}
