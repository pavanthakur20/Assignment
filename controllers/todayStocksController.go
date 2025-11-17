package controllers

import (
	"assignment/initializers"
	"assignment/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTodayStocks(c *gin.Context) {
	log := initializers.Log
	userID := c.Param("userId")

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var rewards []models.StockReward
	err := initializers.DB.Where("user_id = ? AND reward_timestamp >= ? AND reward_timestamp < ?",
		userID, startOfDay, endOfDay).
		Order("reward_timestamp DESC").
		Find(&rewards).Error

	if err != nil {
		log.WithError(err).WithField("user_id", userID).Error("Failed to fetch today's stocks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch today's stocks",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"date":          startOfDay.Format("2006-01-02"),
		"total_rewards": len(rewards),
		"rewards":       rewards,
	})
}
