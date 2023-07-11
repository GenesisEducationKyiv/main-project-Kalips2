package handler

import (
	"btc-app/config"
	"btc-app/pkg/domain/model"
	"btc-app/pkg/domain/service"
	"btc-app/pkg/presentation/presenter"
	"net/http"
)

type RateHandlerImpl struct {
	conf        *config.Config
	rateService service.RateService
}

func (rateHr *RateHandlerImpl) GetCurrentRateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		curPair := model.NewCurrencyPair(rateHr.conf.Crypto.CurrencyTo, rateHr.conf.Crypto.CurrencyFrom)
		if rate, err := rateHr.rateService.GetCurrencyRate(*curPair); err != nil {
			presenter.PresentErrorByBadRequest(w, err)
		} else {
			presenter.PresentRate(w, rate)
		}
	}
}

func NewRateHandler(c *config.Config, rateService service.RateService) *RateHandlerImpl {
	return &RateHandlerImpl{
		conf:        c,
		rateService: rateService,
	}
}
