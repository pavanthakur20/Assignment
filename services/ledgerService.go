package services

import (
	"assignment/models"
	"fmt"

	"gorm.io/gorm"
)

func CalculateCompanyCharges(stockCost float64) *models.CompanyCharges {
	brokerage := stockCost * 0.0003
	stt := stockCost * 0.001
	gst := brokerage * 0.18

	totalCost := stockCost + brokerage + stt + gst

	return &models.CompanyCharges{
		StockCost: stockCost,
		Brokerage: float64(int(brokerage*100)) / 100,
		STT:       float64(int(stt*100)) / 100,
		GST:       float64(int(gst*100)) / 100,
		TotalCost: float64(int(totalCost*100)) / 100,
	}
}

func RecordLedgerEntriesGORM(tx *gorm.DB, reward *models.StockReward, charges *models.CompanyCharges) error {
	totalCashOutflow := charges.TotalCost

	entries := []models.LedgerEntry{
		{
			RewardID:     reward.ID,
			AccountType:  models.AccountTypeStockAsset,
			StockSymbol:  &reward.StockSymbol,
			DebitAmount:  charges.StockCost,
			CreditAmount: 0,
			Quantity:     &reward.Quantity,
			Description:  fmt.Sprintf("Stock acquired: %.6f shares of %s at ₹%.2f", reward.Quantity, reward.StockSymbol, reward.StockPriceAtReward),
		},
		{
			RewardID:     reward.ID,
			AccountType:  models.AccountTypeBrokerageExp,
			StockSymbol:  &reward.StockSymbol,
			DebitAmount:  charges.Brokerage,
			CreditAmount: 0,
			Description:  fmt.Sprintf("Brokerage expense for %s (0.03%%)", reward.StockSymbol),
		},
		{
			RewardID:     reward.ID,
			AccountType:  models.AccountTypeSTTExp,
			StockSymbol:  &reward.StockSymbol,
			DebitAmount:  charges.STT,
			CreditAmount: 0,
			Description:  fmt.Sprintf("Securities Transaction Tax for %s (0.1%%)", reward.StockSymbol),
		},
		{
			RewardID:     reward.ID,
			AccountType:  models.AccountTypeGSTExp,
			StockSymbol:  &reward.StockSymbol,
			DebitAmount:  charges.GST,
			CreditAmount: 0,
			Description:  fmt.Sprintf("GST on brokerage for %s (18%%)", reward.StockSymbol),
		},
		{
			RewardID:     reward.ID,
			AccountType:  models.AccountTypeCashAccount,
			StockSymbol:  &reward.StockSymbol,
			DebitAmount:  0,
			CreditAmount: totalCashOutflow,
			Description:  fmt.Sprintf("Cash paid for stock purchase and fees for %s", reward.StockSymbol),
		},
	}

	if err := tx.Create(&entries).Error; err != nil {
		return fmt.Errorf("failed to record ledger entry: %v", err)
	}

	return nil
}
