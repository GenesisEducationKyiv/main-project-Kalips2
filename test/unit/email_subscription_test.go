package unit

import (
	"btc-app/config"
	"btc-app/handler"
	"btc-app/service"
	"btc-app/template/exception"
	"btc-app/test/unit/repository"
	serviceTest "btc-app/test/unit/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

type EmailSubscriptionTestInfo struct {
	testEmail    string
	emailRepo    *repository.MockEmailRepository
	emailService handler.EmailService
}

var emailInfo *EmailSubscriptionTestInfo

func TestMain(t *testing.M) {
	emailInfo = setUpTest()
}

func TestSubscribeEmailSuccess(t *testing.T) {
	testEmail, repo, emailService := getComponents(emailInfo)
	repo.On("CheckEmailIsExist", testEmail).Return(false, nil)
	repo.On("SaveEmailToStorage", testEmail).Return(nil)

	err := emailService.SubscribeEmail(testEmail)

	assert.NoError(t, err)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertCalled(t, "SaveEmailToStorage", testEmail)
}

func TestSubscribeEmailFailed(t *testing.T) {
	testEmail, repo, emailService := getComponents(emailInfo)
	repo.On("CheckEmailIsExist", testEmail).Return(true, nil)

	err := emailService.SubscribeEmail(testEmail)

	assert.Error(t, err, exception.ErrEmailIsAlreadySubscribed)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertNotCalled(t, "SaveEmailToStorage", testEmail)
}

func setUpTest() *EmailSubscriptionTestInfo {
	emailRepo := &repository.MockEmailRepository{}
	emailSender := &serviceTest.MockGoMailSender{}
	rateService := &serviceTest.MockRateService{}
	emailService := service.NewEmailService(&config.Config{}, rateService, emailRepo, emailSender)
	return &EmailSubscriptionTestInfo{
		testEmail:    "testEmail@gmail.com",
		emailRepo:    emailRepo,
		emailService: emailService,
	}
}

func getComponents(emailInfo *EmailSubscriptionTestInfo) (string, *repository.MockEmailRepository, handler.EmailService) {
	return emailInfo.testEmail, emailInfo.emailRepo, emailInfo.emailService
}
