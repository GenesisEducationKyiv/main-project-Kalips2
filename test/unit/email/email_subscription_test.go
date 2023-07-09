package email

import (
	"btc-app/config"
	"btc-app/pkg/application"
	"btc-app/template/cerror"
	"btc-app/test/unit/repository_mock"
	serviceTest "btc-app/test/unit/service_mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type EmailSubscriptionTestInfo struct {
	testEmail    string
	emailRepo    *repository_mock.MockEmailRepository
	emailService application.EmailService
}

var emailInfo *EmailSubscriptionTestInfo

func TestMain(t *testing.M) {
	emailInfo = setUpEmailTest()
}

func TestSubscribeEmailSuccess(t *testing.T) {
	testEmail, repo, emailService := getComponents(emailInfo)
	repo.On("CheckEmailIsExist", testEmail).Return(false, nil)
	repo.On("SaveEmail", testEmail).Return(nil)

	err := emailService.SubscribeEmail(testEmail)

	assert.NoError(t, err)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertCalled(t, "SaveEmail", testEmail)
}

func TestSubscribeEmailFailed(t *testing.T) {
	testEmail, repo, emailService := getComponents(emailInfo)
	repo.On("CheckEmailIsExist", testEmail).Return(true, nil)

	err := emailService.SubscribeEmail(testEmail)

	assert.Error(t, err, error.ErrEmailIsAlreadySubscribed)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertNotCalled(t, "SaveEmail", testEmail)
}

func setUpEmailTest() *EmailSubscriptionTestInfo {
	emailRepo := &repository_mock.MockEmailRepository{}
	emailSender := &serviceTest.MockGoMailSender{}
	rateService := &serviceTest.MockRateService{}
	emailService := application.NewEmailService(config.CryptoConfig{}, rateService, emailRepo, emailSender)
	return &EmailSubscriptionTestInfo{
		testEmail:    "testEmail@gmail.com",
		emailRepo:    emailRepo,
		emailService: emailService,
	}
}

func getComponents(emailInfo *EmailSubscriptionTestInfo) (string, *repository_mock.MockEmailRepository, application.EmailService) {
	return emailInfo.testEmail, emailInfo.emailRepo, emailInfo.emailService
}
