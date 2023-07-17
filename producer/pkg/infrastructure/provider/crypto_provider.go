package provider

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
	"producer/template/message"
	"strconv"
)

type (
	CryptoProvider struct {
		name       string
		URL        string
		pathToRate string
		logger     service.Logger
	}
)

func NewCryptoProvider(name string, URL string, pathToRate string, logger service.Logger) *CryptoProvider {
	return &CryptoProvider{name: name, URL: URL, pathToRate: pathToRate, logger: logger}
}

func (pr *CryptoProvider) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := pr.getRateByURL(pr.URL, pr.pathToRate, curPair)
	return rate, err
}

func (pr *CryptoProvider) getRateByURL(prvUrl string, pathToRate string, curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	var err error
	var resp *http.Response

	url := fmt.Sprintf(prvUrl, curPair.GetQuote(), curPair.GetBase())
	_ = pr.logger.LogInfo(fmt.Sprintf("Request %s to %s was sended.", url, pr.name))
	if resp, err = pr.sendGetRequestTo(url); err != nil {
		_ = pr.logger.LogError(fmt.Sprintf("Fail to send %s to %s. Error: %s", prvUrl, pr.name, resp.Proto))
		return nil, errors.Wrap(err, message.FailToGetRateMessage)
	}

	rate, err := pr.getRateFromJsonResponse(resp, pathToRate)
	if err != nil {
		_ = pr.logger.LogError(fmt.Sprintf("Error during parsing response %s from %s", resp.Proto, pr.name))
		return nil, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return model.NewCurrencyRate(curPair, rate), err
}

func (pr *CryptoProvider) sendGetRequestTo(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)
	return
}

func (pr *CryptoProvider) getRateFromJsonResponse(resp *http.Response, pathToRate string) (float64, error) {
	var err error
	var rate float64

	body, err := pr.getBytesFromResponse(resp)
	if err != nil {
		return rate, err
	}

	price := gjson.GetBytes(body, pathToRate)
	if !price.Exists() {
		return rate, errors.New("failed to get rate from response")
	}

	rate, err = strconv.ParseFloat(price.String(), 64)
	return rate, err
}

func (pr *CryptoProvider) getBytesFromResponse(resp *http.Response) ([]byte, error) {
	var err error
	var body []byte

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read body of response")
	}

	_ = pr.logger.LogInfo(fmt.Sprintf("Response from %s: %s", pr.name, body))
	return body, err
}
