package provider

import (
	"fmt"
	"producer/config"
	"producer/pkg/application"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
)

type (
	CryptoProvidersChainImpl struct {
		cryptoProvider application.RateProvider
		next           CryptoProvidersChain
	}

	CryptoProvidersChain interface {
		GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error)
		SetNext(next *CryptoProvider) *CryptoProvidersChainImpl
	}
)

func (chain *CryptoProvidersChainImpl) SetNext(next *CryptoProvider) *CryptoProvidersChainImpl {
	chain.next = NewEmptyChain(next)
	return chain
}

func (chain *CryptoProvidersChainImpl) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	var err error
	rate, err := chain.cryptoProvider.GetRate(curPair)
	if err != nil && chain.next != nil {
		rate, err = chain.next.GetRate(curPair)
	}
	return rate, err
}

func NewEmptyChain(cryptoProvider *CryptoProvider) *CryptoProvidersChainImpl {
	return &CryptoProvidersChainImpl{cryptoProvider: cryptoProvider, next: nil}
}

func SetupChainOfProviders(c config.CryptoConfig, logger service.Logger) CryptoProvidersChain {
	CryptoCompareProvider := NewCryptoProvider("Crypto Compare", c.CryptoCompareProviderURL, c.CurrencyTo, logger)
	CoinMarketProvider := NewCryptoProvider("Coin Market", c.CoinMarketProviderURL, fmt.Sprintf("data.%s.quote.%s.price", c.CurrencyFrom, c.CurrencyTo), logger)
	CoinApiProvider := NewCryptoProvider("Coin Api", c.CoinApiProviderURL, "rate", logger)

	chain := NewEmptyChain(CryptoCompareProvider).SetNext(CoinMarketProvider).SetNext(CoinApiProvider)
	return chain
}
