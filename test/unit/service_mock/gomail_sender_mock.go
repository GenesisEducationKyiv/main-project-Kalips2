package service_mock

import (
	"btc-app/pkg/domain"
	"github.com/stretchr/testify/mock"
)

type MockGoMailSender struct {
	mock.Mock
}

func (m *MockGoMailSender) SendMessageTo(message *domain.EmailMessage, recipients []domain.Email) error {
	result := m.Called(message, recipients)
	return result.Error(0)
}
