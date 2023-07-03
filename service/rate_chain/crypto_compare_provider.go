package rate_chain

import (
	"btc-app/model"
)

type CryptoCompareProvider struct {
	providerURL string
	pathToRate  string
	next        *CryptoChain
}

func NewCryptoCompareProvider(providerURL string, pathToRate string) CryptoChain {
	return &CryptoCompareProvider{
		providerURL: providerURL,
		pathToRate:  pathToRate,
	}
}

func (pr *CryptoCompareProvider) SetNext(next *CryptoChain) {
	pr.next = next
}

func (pr *CryptoCompareProvider) GetCurrencyRate(currFrom string, currTo string) (*model.Rate, error) {
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
