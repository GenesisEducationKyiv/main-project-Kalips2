package integration

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"producer/config"
	"producer/pkg/application"
	"producer/pkg/domain/model"
	"producer/pkg/infrastructure/provider"
	"producer/pkg/infrastructure/repository"
	"producer/pkg/infrastructure/sender"
	cerror "producer/template/cerror"
	"producer/test/unit/service_mock"
	"strconv"
	"testing"
)

var (
	pathToStorage            = "subscriptions-test.csv"
	permToOpenTheTestStorage = "0644"
)

func TestSubscribeEmailSuccess(t *testing.T) {
	resetStorageFile()
	testEmail := model.NewEmail("test@example.com")
	emailService := InitTestEmailService()

	serviceError := emailService.SubscribeEmail(testEmail)

	records, _ := repository.ReadEmailsFromStorage(pathToStorage)
	assert.NoError(t, serviceError)
	assert.Equal(t, 1, countOfEmailsIn(testEmail, records))
}

func TestSubscribeEmailFailed(t *testing.T) {
	resetStorageFile()
	testEmail := model.NewEmail("loremipsum@gmail.com")
	emailService := InitTestEmailService()
	_ = repository.WriteEmailToStorage(testEmail, pathToStorage, permToOpenTheTestStorage)

	serviceError := emailService.SubscribeEmail(testEmail)

	records, _ := repository.ReadEmailsFromStorage(pathToStorage)
	assert.Error(t, serviceError, cerror.ErrEmailIsAlreadySubscribed)
	assert.Equal(t, 1, countOfEmailsIn(testEmail, records))
}

func countOfEmailsIn(element model.Email, in []model.Email) int {
	count := 0
	for _, record := range in {
		if record.GetAddress() == element.GetAddress() {
			count++
		}
	}
	return count
}

func resetStorageFile() {
	emptyJSONArray := []byte("[]")
	permission, _ := strconv.ParseInt(permToOpenTheTestStorage, 0, 32)
	if err := os.WriteFile(pathToStorage, emptyJSONArray, os.FileMode(permission)); err != nil {
		log.Fatal(err)
	}
}

func InitTestEmailService() *application.EmailServiceImpl {
	databaseConfig := config.DatabaseConfig{PathToStorage: "subscriptions-test.csv", PermissionToStorage: "0644"}
	cryptoConfig := config.CryptoConfig{}
	mailConfig := config.MailConfig{}

	logger := &service_mock.MockLogger{}
	emailRepository := repository.NewEmailRepository(databaseConfig, logger)
	emailSender := sender.NewEmailSender(mailConfig)
	rateService := application.NewRateService(cryptoConfig, provider.SetupChainOfProviders(cryptoConfig, logger), logger)
	return application.NewEmailService(cryptoConfig, rateService, emailRepository, emailSender, logger)
}
