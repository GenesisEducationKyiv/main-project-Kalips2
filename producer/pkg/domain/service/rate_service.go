package service

import "producer/pkg/domain/model"

type RateService interface {
	GetCurrencyRate(curPair model.CurrencyPair) (*model.CurrencyRate, error)
}
