package service_mock

import (
	"btc-app/model"
	"github.com/stretchr/testify/mock"
)

type MockGoMailSender struct {
	mock.Mock
}

func (m *MockGoMailSender) SendMessageTo(message *model.Message, recipients []model.Email) error {
	result := m.Called(message, recipients)
	return result.Error(0)
}
