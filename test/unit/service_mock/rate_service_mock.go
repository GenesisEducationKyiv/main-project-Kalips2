package service_mock

import (
	"btc-app/pkg/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	result := m.Called(curPair)
	return result.Get(0).(*model.CurrencyRate), result.Error(1)
}
