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
	repo, emailService := setUpSubscribeEmailTest()
	repo.On("CheckEmailIsExist", testEmail).Return(false, nil)
	repo.On("SaveEmailToStorage", testEmail).Return(nil)

	err := emailService.SubscribeEmail(testEmail)

	assert.NoError(t, err)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertCalled(t, "SaveEmailToStorage", testEmail)
}

func TestSubscribeEmailFailed(t *testing.T) {
	repo, emailService := setUpSubscribeEmailTest()
	repo.On("CheckEmailIsExist", testEmail).Return(true, nil)

	err := emailService.SubscribeEmail(testEmail)

	assert.Error(t, err, exception.ErrEmailIsAlreadySubscribed)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertNotCalled(t, "SaveEmailToStorage", testEmail)
}

func setUpSubscribeEmailTest() (*repository.MockEmailRepository, service.EmailService) {
	repo := &repository.MockEmailRepository{}
	emailService := &service.EmailServiceImpl{}
	emailService.SetEmailRepository(repo)
	return repo, emailService
}
