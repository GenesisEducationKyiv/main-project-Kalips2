package rate

import (
	"btc-app/config"
	"btc-app/handler"
	"btc-app/model"
	"btc-app/service"
	"btc-app/test/unit/service_mock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type RateTestInfo struct {
	rateProvider *service_mock.MockCryptoProvider
	cryptoConfig config.CryptoConfig
	rateService  handler.RateService
}

var rateInfo *RateTestInfo

func TestMain(t *testing.M) {
	rateInfo = setUpRateTest()
}

func TestGetRateSuccessful(t *testing.T) {
	expRate := model.Rate{Value: 999.876}
	rateProvider, cryptoConfig, rateService := getComponents(rateInfo)
	rateProvider.On("GetCurrencyRate", cryptoConfig.CurrencyFrom, cryptoConfig.CurrencyTo).Return(expRate, nil)

	rate, err := rateService.GetRate()

	assert.NoError(t, err)
	assert.Equal(t, expRate, rate)
	rateProvider.AssertNumberOfCalls(t, "GetCurrencyRate", 1)
}

func TestGetRateFailed(t *testing.T) {
	rateProvider, cryptoConfig, rateService := getComponents(rateInfo)
	expErr := errors.New("failed to get rate from response")
	rateProvider.On("GetCurrencyRate", cryptoConfig.CurrencyFrom, cryptoConfig.CurrencyTo).Return(nil, expErr)

	rate, err := rateService.GetRate()

	assert.Error(t, expErr, err)
	assert.Equal(t, nil, rate)
	rateProvider.AssertNumberOfCalls(t, "GetCurrencyRate", 1)
}

func getComponents(info *RateTestInfo) (*service_mock.MockCryptoProvider, config.CryptoConfig, handler.RateService) {
	return info.rateProvider, info.cryptoConfig, info.rateService
}

func setUpRateTest() *RateTestInfo {
	cryptoProvider := service_mock.MockCryptoProvider{}
	cryptoConfig := config.CryptoConfig{CurrencyFrom: "BTC", CurrencyTo: "UAH"}

	rateService := service.NewRateService(cryptoConfig, &cryptoProvider)
	return &RateTestInfo{rateProvider: &cryptoProvider, rateService: rateService, cryptoConfig: cryptoConfig}
}
