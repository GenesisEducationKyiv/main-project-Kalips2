package unit

import (
	"btc-app/service"
	"btc-app/template/exception"
	"btc-app/tests/unit/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testEmail = "test@example.com"

func TestSubscribeEmailSuccess(t *testing.T) {
	repo := new(repository.MockEmailRepository)
	repo.On("CheckEmailIsExist", testEmail).Return(false, nil)
	repo.On("SaveEmailToStorage", testEmail).Return(nil)

	emailService := &service.EmailServiceImpl{}
	emailService.SetEmailRepository(repo)

	err := emailService.SubscribeEmail(testEmail)

	assert.NoError(t, err)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertCalled(t, "SaveEmailToStorage", testEmail)
}

func TestSubscribeEmailFailed(t *testing.T) {
	repo := new(repository.MockEmailRepository)
	repo.On("CheckEmailIsExist", testEmail).Return(true, nil)

	emailService := &service.EmailServiceImpl{}
	emailService.SetEmailRepository(repo)

	err := emailService.SubscribeEmail(testEmail)

	assert.Error(t, err, exception.EmailIsAlreadySubscribed)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertNotCalled(t, "SaveEmailToStorage", testEmail)
}
