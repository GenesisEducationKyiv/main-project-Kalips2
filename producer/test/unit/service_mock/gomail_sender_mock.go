package service_mock

import (
	"github.com/stretchr/testify/mock"
	"producer/pkg/domain/model"
)

type MockGoMailSender struct {
	mock.Mock
}

func (m *MockGoMailSender) SendMessageTo(message *model.EmailMessage, recipients []model.Email) error {
	result := m.Called(message, recipients)
	return result.Error(0)
}
