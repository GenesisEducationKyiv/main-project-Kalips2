package application

import (
	"github.com/pkg/errors"
	"producer/config"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
	"producer/template/message"
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

	_ = service.logger.LogInfo("GetRate started from RateService")
	if rate, err = service.rateProvider.GetRate(curPair); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	_ = service.logger.LogInfo("GetRate finished from RateService")
	return rate, err
}

func NewRateService(c config.CryptoConfig, rateProvider RateProvider, logger service.Logger) *RateServiceImpl {
	return &RateServiceImpl{
		conf:         c,
		rateProvider: rateProvider,
		logger:       logger,
	}
}
