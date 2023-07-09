package integration

import (
	"btc-app/config"
	"btc-app/pkg/application"
	"btc-app/pkg/domain/model"
	"btc-app/pkg/infrastructure/provider"
	"btc-app/pkg/infrastructure/repository"
	"btc-app/pkg/infrastructure/sender"
	"btc-app/template/cerror"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
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
	assert.Error(t, serviceError, error.ErrEmailIsAlreadySubscribed)
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

	emailRepository := repository.NewEmailRepository(databaseConfig)
	emailSender := sender.NewEmailSender(mailConfig)
	rateService := application.NewRateService(cryptoConfig, provider.NewChainOfProviders(cryptoConfig))
	return application.NewEmailService(cryptoConfig, rateService, emailRepository, emailSender)
}
