package provider

import (
	"btc-app/pkg/application"
	"btc-app/pkg/domain/model"
	"time"
)

type (
	CachedCryptoProvider struct {
		cacheDuration time.Duration
		timeProvider  TimeProvider
		rateProvider  application.RateProvider
		lastRate      *model.CurrencyRate
		lastTime      time.Time
	}

	TimeProvider interface {
		Now() time.Time
	}
)

func (cache *CachedCryptoProvider) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	var err error
	rate := cache.lastRate
	timeDiff := cache.timeProvider.Now().Sub(cache.lastTime)

	if timeDiff >= cache.cacheDuration || rate == nil {
		rate, err = cache.rateProvider.GetRate(curPair)
		cache.lastRate = rate
		cache.lastTime = cache.timeProvider.Now()
	}
	return rate, err
}

func NewCachedCryptoProvider(cacheDuration time.Duration, rateProvider application.RateProvider, timeProvider TimeProvider) *CachedCryptoProvider {
	return &CachedCryptoProvider{
		cacheDuration: cacheDuration,
		rateProvider:  rateProvider,
		timeProvider:  timeProvider,
		lastTime:      timeProvider.Now().Add(-cacheDuration),
	}
}
