package service_mock

import (
	"btc-app/pkg/domain/model"
	"github.com/stretchr/testify/mock"
)

type MockGoMailSender struct {
	mock.Mock
}

func (m *MockGoMailSender) SendMessageTo(message *model.EmailMessage, recipients []model.Email) error {
	result := m.Called(message, recipients)
	return result.Error(0)
}
