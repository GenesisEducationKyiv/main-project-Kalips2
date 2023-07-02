package service_mock

import (
	"btc-app/model"
	"btc-app/service"
	"github.com/stretchr/testify/mock"
)

type MockCryptoProvider struct {
	mock.Mock
}

func (m *MockCryptoProvider) GetRate(currFrom string, currTo string) (*model.Rate, error) {
	result := m.Called(currFrom, currTo)
	return result.Get(0).(*model.Rate), result.Error(1)
}

func (m *MockCryptoProvider) SetNext(chain service.CryptoChain) {
	m.Called(chain)
}
