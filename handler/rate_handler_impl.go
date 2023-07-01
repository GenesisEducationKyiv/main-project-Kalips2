package handler

import (
	"btc-app/config"
	"btc-app/service"
	"fmt"
	"net/http"
	"strconv"
)

type RateHandlerImpl struct {
	conf        *config.Config
	rateService service.RateService
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

func NewRateHandler(c *config.Config) *RateHandlerImpl {
	return &RateHandlerImpl{
		rateService: service.NewRateService(c),
		conf:        c,
	}
}
