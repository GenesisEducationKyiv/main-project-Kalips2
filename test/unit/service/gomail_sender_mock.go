package service

import (
	"btc-app/model"
	"github.com/stretchr/testify/mock"
)

type MockGoMailSender struct {
	mock.Mock
}

func (m *MockGoMailSender) SendMessageTo(message *model.Message, recipients []string) error {
	result := m.Called(message, recipients)
	return result.Error(0)
}
