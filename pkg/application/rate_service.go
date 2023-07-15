package application

import (
	"btc-app/config"
	"btc-app/pkg/domain/model"
	"btc-app/pkg/domain/service"
	"btc-app/template/message"
	"github.com/pkg/errors"
)

type (
	RateServiceImpl struct {
		rateProvider RateProvider
		conf         config.CryptoConfig
		logger       service.Logger
	}

	RateProvider interface {
		GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error)
	}
)

func (service RateServiceImpl) GetCurrencyRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	var err error
	var rate *model.CurrencyRate

	if rate, err = service.rateProvider.GetRate(curPair); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return rate, err
}

func NewRateService(c config.CryptoConfig, rateProvider RateProvider) *RateServiceImpl {
	return &RateServiceImpl{
		conf:         c,
		rateProvider: rateProvider,
	}
}
