package integration

import (
	"btc-app/config"
	"btc-app/model"
	"btc-app/repository"
	"btc-app/service"
	"btc-app/service/rate_chain"
	"btc-app/template/exception"
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
	testEmail := model.Email{Value: "test@example.com"}
	emailService := InitTestEmailService()

	serviceError := emailService.SubscribeEmail(testEmail.Value)

	records, _ := repository.ReadEmailsFromStorage(pathToStorage)
	assert.NoError(t, serviceError)
	assert.Equal(t, 1, countOfEmailsIn(testEmail, records))
}

func TestSubscribeEmailFailed(t *testing.T) {
	resetStorageFile()
	testEmail := model.Email{Value: "loremipsum@gmail.com"}
	emailService := InitTestEmailService()
	_ = repository.WriteEmailToStorage(testEmail, pathToStorage, permToOpenTheTestStorage)

	serviceError := emailService.SubscribeEmail(testEmail.Value)

	records, _ := repository.ReadEmailsFromStorage(pathToStorage)
	assert.Error(t, serviceError, exception.ErrEmailIsAlreadySubscribed)
	assert.Equal(t, 1, countOfEmailsIn(testEmail, records))
}

func countOfEmailsIn(element model.Email, in []model.Email) int {
	count := 0
	for _, record := range in {
		if record == element {
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

func InitTestEmailService() *service.EmailServiceImpl {
	databaseConfig := config.DatabaseConfig{PathToStorage: "subscriptions-test.csv", PermissionToStorage: "0644"}
	cryptoConfig := config.CryptoConfig{}
	mailConfig := config.MailConfig{}

	emailRepository := repository.NewEmailRepository(databaseConfig)
	emailSender := service.NewEmailSender(mailConfig)
	rateService := service.NewRateService(cryptoConfig, rate_chain.InitChainOfProviders(cryptoConfig))
	return service.NewEmailService(cryptoConfig, rateService, emailRepository, emailSender)
}
