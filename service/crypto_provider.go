package service

import (
	"btc-app/config"
	"btc-app/model"
	"btc-app/template/message"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
)

type (
	CryptoProvider struct {
		name       string
		URL        string
		pathToRate string
		next       CryptoChain
	}

	CryptoChain interface {
		GetRate(currFrom string, currTo string) (*model.Rate, error)
		SetNext(chain CryptoChain)
	}
)

func (pr *CryptoProvider) SetNext(next CryptoChain) {
	pr.next = next
}

func (pr *CryptoProvider) GetRate(currFrom string, currTo string) (*model.Rate, error) {
	var err error
	rate, err := pr.getRateByURL(pr.URL, pr.pathToRate, currFrom, currTo)
	if err != nil && pr.next != nil {
		rate, err = pr.next.GetRate(currFrom, currTo)
	}
	return rate, err
}

func NewCryptoProvider(name string, providerURL string, pathToRate string) CryptoChain {
	return &CryptoProvider{
		name:       name,
		URL:        providerURL,
		pathToRate: pathToRate,
	}
}

func NewChainOfProviders(c config.CryptoConfig) CryptoChain {
	cryptoCompareProvider := NewCryptoProvider("Crypto Compare", c.CryptoCompareProviderURL, c.CurrencyTo)
	coinMarketProvider := NewCryptoProvider("Coin Market", c.CoinMarketProviderURL, fmt.Sprintf("data.%s.quote.%s.price", c.CurrencyFrom, c.CurrencyTo))
	coinApiProvider := NewCryptoProvider("Coin Api", c.CoinApiProviderURL, "rate")

	cryptoCompareProvider.SetNext(coinMarketProvider)
	coinMarketProvider.SetNext(coinApiProvider)
	return cryptoCompareProvider
}

func (pr *CryptoProvider) getRateByURL(prvUrl string, pathToRate string, currTo string, currFrom string) (*model.Rate, error) {
	var err error
	var resp *http.Response
	rate := &model.Rate{}

	if resp, err = pr.sendGetRequestTo(fmt.Sprintf(prvUrl, currTo, currFrom)); err != nil {
		log.Printf("Error during sending request to %s", pr.name)
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}

	if rate, err = pr.getRateFromHttpResponse(resp, pathToRate); err != nil {
		log.Printf("Error during parsing response from %s", pr.name)
		return rate, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return rate, err
}

func (pr *CryptoProvider) sendGetRequestTo(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)
	return
}

func (pr *CryptoProvider) getRateFromHttpResponse(resp *http.Response, pathToRate string) (*model.Rate, error) {
	var err error
	rate := model.Rate{}

	body, err := pr.getBytesFromResponse(resp)
	if err != nil {
		return &rate, err
	}

	price := gjson.GetBytes(body, pathToRate)
	if !price.Exists() {
		return &rate, errors.New("failed to get rate from response")
	}

	err = rate.SetValue(price.String())
	return &rate, err
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
