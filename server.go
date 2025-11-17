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
	initializers.InitLogger()
	initializers.ConnectDB()
}

func main() {
	services.StartPriceUpdateScheduler()

	server := gin.Default()

	server.POST("/reward", controllers.RewardUser)
	server.GET("/today-stocks/:userId", controllers.GetTodayStocks)
	server.GET("/historical-inr/:userId", controllers.GetHistoricalINR)
	server.GET("/stats/:userId", controllers.GetUserStats)
	server.GET("/portfolio/:userId", controllers.GetUserPortfolio)

	port := initializers.GetEnv("PORT", "8080")

	fmt.Printf("Server starting on port %s\n", port)

	if err := server.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}
