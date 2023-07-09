package service_mock

import (
	"btc-app/pkg/application"
	"btc-app/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type MockCryptoProvider struct {
	mock.Mock
}

func (m *MockCryptoProvider) GetRate(curPair domain.CurrencyPair) (*domain.CurrencyRate, error) {
	result := m.Called(curPair)
	return result.Get(0).(*domain.CurrencyRate), result.Error(1)
}

func (m *MockCryptoProvider) SetNext(chain application.ProvidersChain) {
	m.Called(chain)
}
