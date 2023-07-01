package repository

import "github.com/stretchr/testify/mock"

type MockEmailRepository struct {
	mock.Mock
}

func (m *MockEmailRepository) CheckEmailIsExist(email string) (bool, error) {
	result := m.Called(email)
	return result.Bool(0), result.Error(1)
}

func (m *MockEmailRepository) SaveEmailToStorage(email string) error {
	result := m.Called(email)
	return result.Error(0)
}

func (m *MockEmailRepository) GetEmailsFromStorage() ([]string, error) {
	result := m.Called()
	return result.Get(0).([]string), result.Error(1)
}
