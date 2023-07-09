package service

import "btc-app/pkg/domain/model"

type RateService interface {
	GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error)
}
