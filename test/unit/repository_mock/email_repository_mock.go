package repository_mock

import (
	"btc-app/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type MockEmailRepository struct {
	mock.Mock
}

func (m *MockEmailRepository) CheckEmailIsExist(email domain.Email) (bool, error) {
	result := m.Called(email)
	return result.Bool(0), result.Error(1)
}

func (m *MockEmailRepository) SaveEmail(email domain.Email) error {
	result := m.Called(email)
	return result.Error(0)
}

func (m *MockEmailRepository) GetEmailsFromStorage() ([]domain.Email, error) {
	result := m.Called()
	return result.Get(0).([]domain.Email), result.Error(1)
}
