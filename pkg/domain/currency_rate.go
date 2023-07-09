package domain

import (
	"strconv"
)

type CurrencyRate struct {
	CurrencyPair
	rate float64
}

func NewCurrencyRate(p CurrencyPair, rate float64) *CurrencyRate {
	return &CurrencyRate{p, rate}
}

func (rate *CurrencyRate) ToString() string {
	return strconv.FormatFloat(rate.rate, 'f', 5, 64)
}
