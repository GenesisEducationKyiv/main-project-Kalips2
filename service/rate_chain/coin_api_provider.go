package rate_chain

import (
	"btc-app/model"
)

type CoinApiProvider struct {
	providerURL string
	pathToRate  string
	next        *CryptoChain
}

func NewCoinApiProvider(providerURL string, pathToRate string) CryptoChain {
	return &CoinApiProvider{
		providerURL: providerURL,
		pathToRate:  pathToRate,
	}
}

func (pr *CoinApiProvider) SetNext(next *CryptoChain) {
	pr.next = next
}

func (pr *CoinApiProvider) GetCurrencyRate(currFrom string, currTo string) (*model.Rate, error) {
	var err error
	rate, err := getCurrencyRate(pr.providerURL, pr.pathToRate, currFrom, currTo)
	if err != nil {
		next := *pr.next
		if next == nil {
			return nil, err
		}

		rate, err = next.GetCurrencyRate(currFrom, currTo)
	}

	return rate, err
}
