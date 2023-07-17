package service_mock

import (
	"github.com/stretchr/testify/mock"
	"producer/pkg/domain/model"
	"producer/pkg/infrastructure/provider"
)

type MockCryptoProvider struct {
	mock.Mock
}

func (m *MockCryptoProvider) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	result := m.Called(curPair)
	return result.Get(0).(*model.CurrencyRate), result.Error(1)
}

func (m *MockCryptoProvider) SetNext(next *provider.CryptoProvider) *provider.CryptoProvidersChainImpl {
	result := m.Called(next)
	return result.Get(0).(*provider.CryptoProvidersChainImpl)
}
