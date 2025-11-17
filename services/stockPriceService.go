package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var stockPriceRanges = map[string]struct{ min, max float64 }{
	"RELIANCE":  {2200.0, 2800.0},
	"TCS":       {3200.0, 4000.0},
	"INFOSYS":   {1400.0, 1800.0},
	"HDFC":      {1500.0, 2000.0},
	"WIPRO":     {400.0, 600.0},
	"ITC":       {380.0, 480.0},
	"BHARTI":    {800.0, 1200.0},
	"SBIN":      {500.0, 700.0},
	"HDFCBANK":  {1400.0, 1700.0},
	"ICICIBANK": {900.0, 1200.0},
}

var (
	currentStockPrices = make(map[string]float64)
	pricesMutex        sync.RWMutex
)

func GetCurrentStockPrice(stockSymbol string) (float64, error) {
	pricesMutex.RLock()
	price, exists := currentStockPrices[stockSymbol]
	pricesMutex.RUnlock()

	if !exists {
		price = generateRandomPrice(stockSymbol)
		pricesMutex.Lock()
		currentStockPrices[stockSymbol] = price
		pricesMutex.Unlock()
	}

	return price, nil
}

func generateRandomPrice(stockSymbol string) float64 {
	priceRange, exists := stockPriceRanges[stockSymbol]
	if !exists {
		priceRange = struct{ min, max float64 }{100.0, 5000.0}
	}
	price := priceRange.min + rand.Float64()*(priceRange.max-priceRange.min)
	return float64(int(price*100)) / 100
}

func UpdateStockPrices() error {

	pricesMutex.Lock()
	defer pricesMutex.Unlock()

	for stockSymbol := range stockPriceRanges {
		price := generateRandomPrice(stockSymbol)
		currentStockPrices[stockSymbol] = price
	}

	fmt.Printf("Stock prices updated\n")
	return nil
}

func StartPriceUpdateScheduler() {
	go func() {
		if err := UpdateStockPrices(); err != nil {
			fmt.Printf("Error updating stock prices: %v\n", err)
		}
	}()

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			if err := UpdateStockPrices(); err != nil {
				fmt.Printf("Error updating stock prices: %v\n", err)
			}
		}
	}()

	fmt.Println("Stock price update scheduler started")
}

func GetCurrentPrices(stockSymbols []string) (map[string]float64, error) {
	prices := make(map[string]float64)

	for _, symbol := range stockSymbols {
		price, err := GetCurrentStockPrice(symbol)
		if err != nil {
			return nil, err
		}
		prices[symbol] = price
	}

	return prices, nil
}
