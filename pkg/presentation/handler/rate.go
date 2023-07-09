package handler

import (
	"btc-app/config"
	"btc-app/pkg/application"
	"btc-app/pkg/domain"
	"btc-app/pkg/presentation/presenter"
	"net/http"
)

type RateHandlerImpl struct {
	conf        *config.Config
	rateService application.RateService
}

func (rateHr *RateHandlerImpl) GetCurrentRateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		curPair := domain.NewCurrencyPair(rateHr.conf.Crypto.CurrencyTo, rateHr.conf.Crypto.CurrencyFrom)
		if rate, err := rateHr.rateService.GetRate(*curPair); err != nil {
			presenter.PresentErrorByBadRequest(w, err)
		} else {
			presenter.PresentRate(w, rate)
		}
	}
}

func NewRateHandler(c *config.Config, rateService application.RateService) *RateHandlerImpl {
	return &RateHandlerImpl{
		conf:        c,
		rateService: rateService,
	}
}
