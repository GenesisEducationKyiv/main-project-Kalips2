package provider

import (
	"btc-app/config"
	"btc-app/pkg/application"
	"btc-app/pkg/domain"
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
		next       application.ProvidersChain
	}
)

func (pr *CryptoProvider) SetNext(next application.ProvidersChain) {
	pr.next = next
}

func (pr *CryptoProvider) GetRate(curPair domain.CurrencyPair) (*domain.CurrencyRate, error) {
	var err error
	rate, err := pr.getRateByURL(pr.URL, pr.pathToRate, curPair)
	if err != nil && pr.next != nil {
		rate, err = pr.next.GetRate(curPair)
	}
	return rate, err
}

func NewCryptoProvider(name string, providerURL string, pathToRate string) application.ProvidersChain {
	return &CryptoProvider{
		name:       name,
		URL:        providerURL,
		pathToRate: pathToRate,
	}
}

func NewChainOfProviders(c config.CryptoConfig) application.ProvidersChain {
	cryptoCompareProvider := NewCryptoProvider("Crypto Compare", c.CryptoCompareProviderURL, c.CurrencyTo)
	coinMarketProvider := NewCryptoProvider("Coin Market", c.CoinMarketProviderURL, fmt.Sprintf("data.%s.quote.%s.price", c.CurrencyFrom, c.CurrencyTo))
	coinApiProvider := NewCryptoProvider("Coin Api", c.CoinApiProviderURL, "rate")

	cryptoCompareProvider.SetNext(coinMarketProvider)
	coinMarketProvider.SetNext(coinApiProvider)
	return cryptoCompareProvider
}

func (pr *CryptoProvider) getRateByURL(prvUrl string, pathToRate string, curPair domain.CurrencyPair) (*domain.CurrencyRate, error) {
	var err error
	var resp *http.Response

	if resp, err = pr.sendGetRequestTo(fmt.Sprintf(prvUrl, curPair.GetQuote(), curPair.GetBase())); err != nil {
		log.Printf("Error during sending request to %s", pr.name)
		return nil, errors.Wrap(err, message.FailToGetRateMessage)
	}

	rate, err := pr.getRateFromHttpResponse(resp, pathToRate)
	if err != nil {
		log.Printf("Error during parsing response from %s", pr.name)
		return nil, errors.Wrap(err, message.FailToGetRateMessage)
	}
	return domain.NewCurrencyRate(curPair, rate), err
}

func (pr *CryptoProvider) sendGetRequestTo(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)
	return
}

func (pr *CryptoProvider) getRateFromHttpResponse(resp *http.Response, pathToRate string) (float64, error) {
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
