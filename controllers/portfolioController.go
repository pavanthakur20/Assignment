package controllers

import (
	"net/http"
	"time"

	"assignment/initializers"
	"assignment/models"
	"assignment/services"

	"github.com/gin-gonic/gin"
)

func GetUserPortfolio(c *gin.Context) {
	userID := c.Param("userId")

	var rewards []models.StockReward
	err := initializers.DB.Where("user_id = ?", userID).
		Order("stock_symbol").
		Find(&rewards).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch portfolio",
			"details": err.Error(),
		})
		return
	}

	var symbols []string
	stockMap := make(map[string]float64)

	for _, reward := range rewards {
		if _, exists := stockMap[reward.StockSymbol]; !exists {
			symbols = append(symbols, reward.StockSymbol)
		}
		stockMap[reward.StockSymbol] += reward.Quantity
	}

	prices, err := services.GetCurrentPrices(symbols)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch current prices",
			"details": err.Error(),
		})
		return
	}

	var userHoldings []models.UserStockHolding
	totalValue := 0.0

	for _, symbol := range symbols {
		quantity := stockMap[symbol]
		price := prices[symbol]
		currentValue := quantity * price

		userHoldings = append(userHoldings, models.UserStockHolding{
			StockSymbol:   symbol,
			TotalQuantity: float64(int(quantity*1000000)) / 1000000,
			CurrentPrice:  float64(int(price*100)) / 100,
			CurrentValue:  float64(int(currentValue*100)) / 100,
		})

		totalValue += currentValue
	}

	if userHoldings == nil {
		userHoldings = []models.UserStockHolding{}
	}

	response := models.PortfolioResponse{
		UserID:      userID,
		Holdings:    userHoldings,
		TotalValue:  float64(int(totalValue*100)) / 100,
		LastUpdated: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
