package handler

import (
	"btc-app/config"
	"btc-app/model"
	"fmt"
	"net/http"
)

type RateHandlerImpl struct {
	conf        *config.Config
	rateService RateService
}

type RateService interface {
	GetRate() (model.Rate, error)
}

func (rateHr *RateHandlerImpl) GetCurrentRateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rate, err := rateHr.rateService.GetRate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			fmt.Fprint(w, rate.String())
		}
	}
}

func NewRateHandler(c *config.Config, rateService RateService) *RateHandlerImpl {
	return &RateHandlerImpl{
		rateService: rateService,
		conf:        c,
	}
}
