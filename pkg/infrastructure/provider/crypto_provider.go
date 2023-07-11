package provider

import (
	"btc-app/pkg/domain/model"
	"btc-app/template/message"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"strconv"
)

type (
	CryptoProvider struct {
		name       string
		URL        string
		pathToRate string
	}
)

func NewCryptoProvider(name string, URL string, pathToRate string) *CryptoProvider {
	return &CryptoProvider{name: name, URL: URL, pathToRate: pathToRate}
}

func (pr *CryptoProvider) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := pr.getRateByURL(pr.URL, pr.pathToRate, curPair)
	return rate, err
}

func (pr *CryptoProvider) getRateByURL(prvUrl string, pathToRate string, curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	var err error
	var resp *http.Response

	if resp, err = pr.sendGetRequestTo(fmt.Sprintf(prvUrl, curPair.GetQuote(), curPair.GetBase())); err != nil {
		log.Printf("Error during sending request to %s", pr.name)
		return nil, errors.Wrap(err, message.FailToGetRateMessage)
	}

	rate, err := pr.getRateFromJsonResponse(resp, pathToRate)
	if err != nil {
		log.Printf("Error during parsing response from %s", pr.name)
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

	log.Printf("Response from %s: %s", pr.name, body)
	return body, err
}
