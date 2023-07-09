package application

import (
	"btc-app/config"
	"btc-app/pkg/domain/model"
	"btc-app/template/message"
	"github.com/pkg/errors"
)

type (
	RateServiceImpl struct {
		rateProvider ProvidersChain
		conf         config.CryptoConfig
	}

	ProvidersChain interface {
		GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error)
		SetNext(nextProvider ProvidersChain)
	}
)

func (service RateServiceImpl) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	var err error
	var rate *model.CurrencyRate

	if rate, err = service.rateProvider.GetRate(curPair); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return rate, err
}

func NewRateService(c config.CryptoConfig, rateProvider ProvidersChain) *RateServiceImpl {
	return &RateServiceImpl{
		conf:         c,
		rateProvider: rateProvider,
	}
}
