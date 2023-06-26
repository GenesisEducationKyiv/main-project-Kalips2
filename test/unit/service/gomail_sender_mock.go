package service

import (
	"github.com/go-gomail/gomail"
	"github.com/stretchr/testify/mock"
)

type MockGoMailSender struct {
	mock.Mock
}

func (m *MockGoMailSender) CreateMessage(emailFrom string, header string, body string) *gomail.Message {
	result := m.Called(emailFrom, header, body)
	return result.Get(0).(*gomail.Message)
}

func (m *MockGoMailSender) SendMessageTo(message *gomail.Message, recipients []string) error {
	result := m.Called(message, recipients)
	return result.Error(0)
}
