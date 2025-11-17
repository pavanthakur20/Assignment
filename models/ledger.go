package models

import (
	"time"
)

type LedgerEntry struct {
	ID           uint        `gorm:"primaryKey;autoIncrement"`
	RewardID     string      `gorm:"type:varchar(255);not null;index"`
	Reward       StockReward `gorm:"foreignKey:RewardID;references:ID;constraint:OnDelete:CASCADE"`
	AccountType  string      `gorm:"type:varchar(50);not null"`
	StockSymbol  *string     `gorm:"type:varchar(50)"`
	DebitAmount  float64     `gorm:"type:numeric(18,4);not null;default:0"`
	CreditAmount float64     `gorm:"type:numeric(18,4);not null;default:0"`
	Quantity     *float64    `gorm:"type:numeric(18,6)"`
	Description  string      `gorm:"type:text"`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
}

func (LedgerEntry) TableName() string {
	return "ledger_entries"
}

const (
	AccountTypeStockAsset   = "STOCK_ASSET"
	AccountTypeCashAccount  = "CASH_ACCOUNT"
	AccountTypeBrokerageExp = "BROKERAGE_EXPENSE"
	AccountTypeSTTExp       = "STT_EXPENSE"
	AccountTypeGSTExp       = "GST_EXPENSE"
)
