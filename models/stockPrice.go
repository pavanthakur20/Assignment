package models

import (
	"time"
)

type UserStockHolding struct {
	StockSymbol   string
	TotalQuantity float64
	CurrentPrice  float64
	CurrentValue  float64
}

type PortfolioResponse struct {
	UserID      string
	Holdings    []UserStockHolding
	TotalValue  float64
	LastUpdated time.Time
}

type StatsResponse struct {
	UserID              string
	TodayRewards        map[string]float64
	CurrentPortfolioINR float64
	TotalSharesRewarded float64
}

type HistoricalINRResponse struct {
	UserID      string
	DailyValues map[string]float64
}
