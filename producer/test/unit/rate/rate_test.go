package rate

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"producer/config"
	"producer/pkg/application"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
	"producer/test/unit/service_mock"
	"testing"
)

type RateTestInfo struct {
	rateProvider *service_mock.MockCryptoProvider
	cryptoConfig config.CryptoConfig
	rateService  service.RateService
}

var rateInfo *RateTestInfo

func TestMain(t *testing.M) {
	rateInfo = setUpRateTest()
}

func TestGetRateSuccessful(t *testing.T) {
	rateProvider, cryptoConfig, rateService := getComponents(rateInfo)
	curPair := model.NewCurrencyPair(cryptoConfig.CurrencyFrom, cryptoConfig.CurrencyTo)
	expRate := model.NewCurrencyRate(*curPair, 999.876)
	rateProvider.On("GetRate", curPair).Return(expRate, nil)

	rate, err := rateService.GetCurrencyRate(*curPair)

	assert.NoError(t, err)
	assert.Equal(t, expRate, rate)
	rateProvider.AssertNumberOfCalls(t, "GetRate", 1)
}

func TestGetRateFailed(t *testing.T) {
	rateProvider, cryptoConfig, rateService := getComponents(rateInfo)
	curPair := model.NewCurrencyPair(cryptoConfig.CurrencyFrom, cryptoConfig.CurrencyTo)
	expErr := errors.New("failed to get rate from response")
	rateProvider.On("GetRate", curPair).Return(nil, expErr)

	rate, err := rateService.GetCurrencyRate(*curPair)

	assert.Error(t, expErr, err)
	assert.Equal(t, nil, rate)
	rateProvider.AssertNumberOfCalls(t, "GetRate", 1)
}

func getComponents(info *RateTestInfo) (*service_mock.MockCryptoProvider, config.CryptoConfig, service.RateService) {
	return info.rateProvider, info.cryptoConfig, info.rateService
}

func setUpRateTest() *RateTestInfo {
	cryptoProvider := service_mock.MockCryptoProvider{}
	logger := &service_mock.MockLogger{}
	cryptoConfig := config.CryptoConfig{CurrencyFrom: "BTC", CurrencyTo: "UAH"}

	rateService := application.NewRateService(cryptoConfig, &cryptoProvider, logger)
	return &RateTestInfo{rateProvider: &cryptoProvider, rateService: rateService, cryptoConfig: cryptoConfig}
}
