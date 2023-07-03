package service

import (
	"btc-app/config"
	"btc-app/model"
	"btc-app/service/rate_chain"
	"btc-app/template/message"
	"github.com/pkg/errors"
)

type RateServiceImpl struct {
	rateProvider rate_chain.CryptoChain
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

func NewRateService(c config.CryptoConfig, rateProvider rate_chain.CryptoChain) *RateServiceImpl {
	return &RateServiceImpl{
		conf:         c,
		rateProvider: rateProvider,
	}
}
