package tests

import (
	"assignment/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateCompanyCharges(t *testing.T) {
	stockCost := 10000.0

	charges := services.CalculateCompanyCharges(stockCost)

	assert.NotNil(t, charges)
	assert.Equal(t, stockCost, charges.StockCost)

	assert.Greater(t, charges.Brokerage, 0.0)
	assert.Greater(t, charges.STT, 0.0)
	assert.Greater(t, charges.GST, 0.0)
	assert.Greater(t, charges.TotalCost, stockCost)
}

func TestCalculateCompanyChargesZeroAmount(t *testing.T) {
	stockCost := 1000.0

	charges := services.CalculateCompanyCharges(stockCost)

	assert.NotNil(t, charges)
	assert.Equal(t, stockCost, charges.StockCost)
	assert.Greater(t, charges.Brokerage, 0.0)
	assert.Greater(t, charges.STT, 0.0)
	assert.Greater(t, charges.GST, 0.0)
	assert.Greater(t, charges.TotalCost, stockCost)
}

func TestCalculateCompanyChargesLargeAmount(t *testing.T) {
	stockCost := 1000000.0

	charges := services.CalculateCompanyCharges(stockCost)

	assert.NotNil(t, charges)
	assert.Equal(t, stockCost, charges.StockCost)

	assert.Equal(t, 300.0, charges.Brokerage)
	assert.Equal(t, 1000.0, charges.STT)
	assert.Equal(t, 54.0, charges.GST)
	assert.Equal(t, 1001354.0, charges.TotalCost)
}
