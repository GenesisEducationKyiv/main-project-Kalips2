package application

import (
	"btc-app/config"
	"btc-app/pkg/domain"
	"btc-app/template/message"
	"github.com/pkg/errors"
)

type (
	RateServiceImpl struct {
		rateProvider ProvidersChain
		conf         config.CryptoConfig
	}

	ProvidersChain interface {
		GetRate(curPair domain.CurrencyPair) (*domain.CurrencyRate, error)
		SetNext(nextProvider ProvidersChain)
	}
)

func (service RateServiceImpl) GetRate(curPair domain.CurrencyPair) (*domain.CurrencyRate, error) {
	var err error
	var rate *domain.CurrencyRate

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
