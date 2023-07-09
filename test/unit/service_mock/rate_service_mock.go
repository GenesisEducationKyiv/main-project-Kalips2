package service_mock

import (
	domain2 "btc-app/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetRate(curPair domain2.CurrencyPair) (*domain2.CurrencyRate, error) {
	result := m.Called()
	return result.Get(0).(*domain2.CurrencyRate), result.Error(1)
}
