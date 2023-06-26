package handler

import (
	"btc-app/config"
	"fmt"
	"net/http"
	"strconv"
)

type RateHandlerImpl struct {
	conf        *config.Config
	rateService RateService
}

type RateService interface {
	GetCurrentRate() (float64, error)
}

func (rateHr *RateHandlerImpl) GetCurrentRateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rate, err := rateHr.rateService.GetCurrentRate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			rateStr := strconv.FormatFloat(rate, 'f', 5, 64)
			fmt.Fprint(w, rateStr)
		}
	}
}

func NewRateHandler(c *config.Config, rateService RateService) *RateHandlerImpl {
	return &RateHandlerImpl{
		rateService: rateService,
		conf:        c,
	}
}
