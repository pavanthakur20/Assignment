package models

import (
	"time"
)

type StockReward struct {
	ID                 string    `gorm:"type:varchar(255);primaryKey"`
	UserID             string    `gorm:"type:varchar(255);not null;index"`
	StockSymbol        string    `gorm:"type:varchar(50);not null"`
	Quantity           float64   `gorm:"type:numeric(18,6);not null"`
	RewardTimestamp    time.Time `gorm:"not null;index"`
	StockPriceAtReward float64   `gorm:"type:numeric(18,4);not null"`
	CreatedAt          time.Time `gorm:"autoCreateTime"`
}

func (StockReward) TableName() string {
	return "stock_rewards"
}

type RewardRequest struct {
	ID              string    `json:"id" binding:"required"`
	UserID          string    `json:"user_id" binding:"required"`
	StockSymbol     string    `json:"stock_symbol" binding:"required"`
	Quantity        float64   `json:"quantity" binding:"required,gt=0"`
	RewardTimestamp time.Time `json:"reward_timestamp" binding:"required"`
}

type RewardResponse struct {
	Success        bool
	Message        string
	Reward         *StockReward
	INRValue       float64
	CompanyCharges *CompanyCharges
}

type CompanyCharges struct {
	StockCost float64
	Brokerage float64
	STT       float64
	GST       float64
	TotalCost float64
}
