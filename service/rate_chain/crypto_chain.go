package rate_chain

import (
	"btc-app/config"
	"btc-app/model"
	"btc-app/template/message"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type CryptoChain interface {
	GetCurrencyRate(currFrom string, currTo string) (*model.Rate, error)
	SetNext(chain *CryptoChain)
}

func InitChainOfProviders(c config.CryptoConfig) CryptoChain {
	cryptoCompareProvider := NewCryptoCompareProvider(c.CryptoCompareProviderURL, c.CurrencyTo)
	coinMarketProvider := NewCoinMarketProvider(c.CoinMarketProviderURL, fmt.Sprintf("data.%s.quote.%s.price", c.CurrencyFrom, c.CurrencyTo))
	coinApiProvider := NewCoinApiProvider(c.CoinApiProviderURL, "rate")

	cryptoCompareProvider.SetNext(&coinMarketProvider)
	coinMarketProvider.SetNext(&coinApiProvider)
	return cryptoCompareProvider
}

func getCurrencyRate(prvUrl string, pathToRate string, currTo string, currFrom string) (*model.Rate, error) {
	var err error
	var resp *http.Response
	rate := &model.Rate{}

	url := fmt.Sprintf(prvUrl, currTo, currFrom)
	if resp, err = http.Get(url); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}

	if rate, err = getRateFromHttpResponse(resp, pathToRate); err != nil {
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return rate, err
}

func getRateFromHttpResponse(resp *http.Response, pathToRate string) (*model.Rate, error) {
	var err error

	rate := model.Rate{}
	body, err := getBytesFromResponse(resp)
	if err != nil {
		return &rate, err
	}

	price := gjson.GetBytes(body, pathToRate)
	if !price.Exists() {
		return &rate, errors.New("failed to get rate from response")
	}

	rate.SetValue(price.String())
	return &rate, err
}

func getBytesFromResponse(resp *http.Response) ([]byte, error) {
	var err error
	var body []byte

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read body of response")
	}

	return body, err
}
