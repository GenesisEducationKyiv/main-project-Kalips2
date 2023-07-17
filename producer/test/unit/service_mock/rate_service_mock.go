package service_mock

import (
	"github.com/stretchr/testify/mock"
	"producer/pkg/domain/model"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetCurrencyRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	result := m.Called(curPair)
	return result.Get(0).(*model.CurrencyRate), result.Error(1)
}
