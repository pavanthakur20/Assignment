package tests

import (
	"assignment/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentStockPrice(t *testing.T) {
	symbols := []string{"RELIANCE", "TCS", "INFOSYS", "HDFC", "WIPRO"}

	for _, symbol := range symbols {
		price, err := services.GetCurrentStockPrice(symbol)
		assert.NoError(t, err)
		assert.Greater(t, price, 0.0)
	}
}

func TestGetCurrentStockPriceConsistency(t *testing.T) {
	symbol := "RELIANCE"

	price1, err1 := services.GetCurrentStockPrice(symbol)
	assert.NoError(t, err1)

	price2, err2 := services.GetCurrentStockPrice(symbol)
	assert.NoError(t, err2)

	assert.Equal(t, price1, price2, "Prices should be consistent within same session")
}

func TestGetCurrentStockPriceUnknownSymbol(t *testing.T) {
	price, err := services.GetCurrentStockPrice("UNKNOWN")
	assert.NoError(t, err)
	assert.Greater(t, price, 0.0)
	assert.LessOrEqual(t, price, 5000.0)
}

func TestGetCurrentStockPriceKnownSymbolRange(t *testing.T) {
	testCases := []struct {
		symbol string
		min    float64
		max    float64
	}{
		{"RELIANCE", 2200.0, 2800.0},
		{"TCS", 3200.0, 4000.0},
		{"INFOSYS", 1400.0, 1800.0},
		{"HDFC", 1500.0, 2000.0},
		{"WIPRO", 400.0, 600.0},
		{"ITC", 380.0, 480.0},
		{"BHARTI", 800.0, 1200.0},
		{"SBIN", 500.0, 700.0},
		{"HDFCBANK", 1400.0, 1700.0},
		{"ICICIBANK", 900.0, 1200.0},
	}

	for _, tc := range testCases {
		t.Run(tc.symbol, func(t *testing.T) {
			price, err := services.GetCurrentStockPrice(tc.symbol)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, price, tc.min, "Price should be >= min")
			assert.LessOrEqual(t, price, tc.max, "Price should be <= max")
		})
	}
}

func TestGetCurrentPricesMultipleSymbols(t *testing.T) {
	symbols := []string{"RELIANCE", "TCS", "INFOSYS"}

	prices, err := services.GetCurrentPrices(symbols)

	assert.NoError(t, err)
	assert.NotNil(t, prices)
	assert.Equal(t, 3, len(prices))

	for _, symbol := range symbols {
		price, exists := prices[symbol]
		assert.True(t, exists, "Price should exist for %s", symbol)
		assert.Greater(t, price, 0.0)
	}
}

func TestGetCurrentPricesEmptyList(t *testing.T) {
	symbols := []string{}

	prices, err := services.GetCurrentPrices(symbols)

	assert.NoError(t, err)
	assert.NotNil(t, prices)
	assert.Equal(t, 0, len(prices))
}

func TestGetCurrentPricesSingleSymbol(t *testing.T) {
	symbols := []string{"RELIANCE"}

	prices, err := services.GetCurrentPrices(symbols)

	assert.NoError(t, err)
	assert.NotNil(t, prices)
	assert.Equal(t, 1, len(prices))
	assert.Greater(t, prices["RELIANCE"], 0.0)
}

func TestGetCurrentPricesMixedKnownUnknown(t *testing.T) {
	symbols := []string{"RELIANCE", "UNKNOWN1", "TCS", "UNKNOWN2"}

	prices, err := services.GetCurrentPrices(symbols)

	assert.NoError(t, err)
	assert.NotNil(t, prices)
	assert.Equal(t, 4, len(prices))

	for _, symbol := range symbols {
		price, exists := prices[symbol]
		assert.True(t, exists)
		assert.Greater(t, price, 0.0)
	}
}

func TestUpdateStockPrices(t *testing.T) {
	oldPrices := make(map[string]float64)
	symbols := []string{"RELIANCE", "TCS", "INFOSYS"}

	for _, symbol := range symbols {
		price, _ := services.GetCurrentStockPrice(symbol)
		oldPrices[symbol] = price
	}

	err := services.UpdateStockPrices()
	assert.NoError(t, err)

	for _, symbol := range symbols {
		newPrice, _ := services.GetCurrentStockPrice(symbol)
		assert.Greater(t, newPrice, 0.0)
	}
}

func TestStockPricePrecision(t *testing.T) {
	price, err := services.GetCurrentStockPrice("RELIANCE")
	assert.NoError(t, err)

	priceStr := string(rune(int(price * 100)))
	assert.LessOrEqual(t, len(priceStr), 4, "Price should have at most 2 decimal places")
}
