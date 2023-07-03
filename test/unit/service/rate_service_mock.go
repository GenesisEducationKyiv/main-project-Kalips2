package service

import (
	"github.com/stretchr/testify/mock"
)

type MockRateService struct {
	mock.Mock
}

func (m *MockRateService) GetCurrentRate() (float64, error) {
	result := m.Called()
	return result.Get(0).(float64), result.Error(1)
}
