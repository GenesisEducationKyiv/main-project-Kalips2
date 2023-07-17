package email

import (
	"github.com/stretchr/testify/assert"
	"producer/config"
	"producer/pkg/application"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
	cerror "producer/template/cerror"
	"producer/test/unit/repository_mock"
	serviceTest "producer/test/unit/service_mock"
	"testing"
)

type EmailSubscriptionTestInfo struct {
	testEmail    model.Email
	emailRepo    *repository_mock.MockEmailRepository
	emailService service.EmailService
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

	assert.Error(t, err, cerror.ErrEmailIsAlreadySubscribed)
	repo.AssertCalled(t, "CheckEmailIsExist", testEmail)
	repo.AssertNotCalled(t, "SaveEmail", testEmail)
}

func setUpEmailTest() *EmailSubscriptionTestInfo {
	emailRepo := &repository_mock.MockEmailRepository{}
	emailSender := &serviceTest.MockGoMailSender{}
	rateService := &serviceTest.MockRateService{}
	logger := &serviceTest.MockLogger{}
	emailService := application.NewEmailService(config.CryptoConfig{}, rateService, emailRepo, emailSender, logger)
	return &EmailSubscriptionTestInfo{
		testEmail:    model.NewEmail("testEmail@gmail.com"),
		emailRepo:    emailRepo,
		emailService: emailService,
	}
}

func getComponents(emailInfo *EmailSubscriptionTestInfo) (model.Email, *repository_mock.MockEmailRepository, service.EmailService) {
	return emailInfo.testEmail, emailInfo.emailRepo, emailInfo.emailService
}
