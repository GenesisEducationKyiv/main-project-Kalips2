package service

import (
	"btc-app/config"
	"btc-app/model"
	"btc-app/template/message"
	"github.com/pkg/errors"
)

type RateServiceImpl struct {
	rateProvider CryptoChain
	conf         config.CryptoConfig
}

func (service RateServiceImpl) GetRate() (*model.Rate, error) {
	var err error
	var rate *model.Rate

	if rate, err = service.rateProvider.GetCurrencyRate(service.conf.CurrencyFrom, service.conf.CurrencyTo); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return rate, err
}

func NewRateService(c config.CryptoConfig, rateProvider CryptoChain) *RateServiceImpl {
	return &RateServiceImpl{
		conf:         c,
		rateProvider: rateProvider,
	}
}
