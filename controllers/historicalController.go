package controllers

import (
	"net/http"
	"time"

	"assignment/initializers"
	"assignment/models"
	"assignment/services"

	"github.com/gin-gonic/gin"
)

func GetHistoricalINR(c *gin.Context) {
	userID := c.Param("userId")

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var rewards []models.StockReward
	err := initializers.DB.Where("user_id = ? AND reward_timestamp < ?", userID, today).
		Find(&rewards).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch historical data",
			"details": err.Error(),
		})
		return
	}

	dailyValues := make(map[string]float64)

	for _, reward := range rewards {
		dateStr := reward.RewardTimestamp.Format("2006-01-02")
		currentPrice, err := services.GetCurrentStockPrice(reward.StockSymbol)
		if err != nil {
			continue
		}

		dailyValues[dateStr] += reward.Quantity * currentPrice
	}

	for date, value := range dailyValues {
		dailyValues[date] = float64(int(value*100)) / 100
	}

	response := models.HistoricalINRResponse{
		UserID:      userID,
		DailyValues: dailyValues,
	}

	if response.DailyValues == nil {
		response.DailyValues = make(map[string]float64)
	}

	c.JSON(http.StatusOK, response)
}
