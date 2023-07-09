package service_mock

import (
	"btc-app/pkg/application"
	"btc-app/pkg/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockCryptoProvider struct {
	mock.Mock
}

func (m *MockCryptoProvider) GetRate(curPair model.CurrencyPair) (*model.CurrencyRate, error) {
	result := m.Called(curPair)
	return result.Get(0).(*model.CurrencyRate), result.Error(1)
}

func (m *MockCryptoProvider) SetNext(chain application.ProvidersChain) {
	m.Called(chain)
}
