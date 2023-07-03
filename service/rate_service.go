package service

import (
	"btc-app/config"
	"btc-app/model"
	"btc-app/template/message"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type RateServiceImpl struct {
	conf *config.Config
}

func (rateService RateServiceImpl) GetRate() (model.Rate, error) {
	var resp *http.Response
	var err error
	rate := model.Rate{}

	conf := rateService.conf
	url := fmt.Sprintf("%s?fsym=%s&tsyms=%s", conf.CryptoApiURL, conf.CurrencyFrom, conf.CurrencyTo)
	if resp, err = http.Get(url); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}

	if rate, err = getRateFromHttpResponse(resp); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return rate, err
}

func getRateFromHttpResponse(resp *http.Response) (model.Rate, error) {
	var err error
	var data map[string]float64
	var rate model.Rate

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return rate, errors.New("Failed to read body of response")
	}
	err = json.Unmarshal(body, &data)
	rate.Value = data["UAH"]
	return rate, err
}

func NewRateService(c *config.Config) *RateServiceImpl {
	return &RateServiceImpl{
		conf: c,
	}
}
