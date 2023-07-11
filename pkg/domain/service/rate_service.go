package service

import "btc-app/pkg/domain/model"

type RateService interface {
	GetCurrencyRate(curPair model.CurrencyPair) (*model.CurrencyRate, error)
}
