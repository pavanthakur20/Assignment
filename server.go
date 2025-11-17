package main

import (
	"assignment/controllers"
	"assignment/initializers"
	"assignment/services"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
}

func main() {
	services.StartPriceUpdateScheduler()

	server := gin.Default()

	v1 := server.Group("/api/v1")
	{
		v1.POST("/reward", controllers.RewardUser)
		v1.GET("/today-stocks/:userId", controllers.GetTodayStocks)
		v1.GET("/historical-inr/:userId", controllers.GetHistoricalINR)
		v1.GET("/stats/:userId", controllers.GetUserStats)
		v1.GET("/portfolio/:userId", controllers.GetUserPortfolio)
	}

	port := initializers.GetEnv("PORT", "8080")

	fmt.Printf("Server starting on port %s\n", port)

	if err := server.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
