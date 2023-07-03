package repository

import (
	"btc-app/model"
	"github.com/stretchr/testify/mock"
)

type MockEmailRepository struct {
	mock.Mock
}

func (m *MockEmailRepository) CheckEmailIsExist(email model.Email) (bool, error) {
	result := m.Called(email)
	return result.Bool(0), result.Error(1)
}

func (m *MockEmailRepository) SaveEmail(email model.Email) error {
	result := m.Called(email)
	return result.Error(0)
}

func (m *MockEmailRepository) GetEmailsFromStorage() ([]model.Email, error) {
	result := m.Called()
	return result.Get(0).([]model.Email), result.Error(1)
}
