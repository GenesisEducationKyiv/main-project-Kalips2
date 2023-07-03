package service

import (
	"btc-app/model"
	"github.com/stretchr/testify/mock"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetRate() (*model.Rate, error) {
	result := m.Called()
	return result.Get(0).(*model.Rate), result.Error(1)
}
