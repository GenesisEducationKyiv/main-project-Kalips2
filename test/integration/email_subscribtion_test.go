package integration

import (
	"btc-app/config"
	"btc-app/repository"
	"btc-app/service"
	"btc-app/template/exception"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	pathToStorage            = "subscriptions-test.csv"
	permToOpenTheTestStorage = 0644
)

func TestSubscribeEmailSuccess(t *testing.T) {
	resetStorageFile()
	testEmail := "test@example.com"
	emailService := InitTestEmailService()

	serviceError := emailService.SubscribeEmail(testEmail)

	records, _ := repository.ReadFromStorage(pathToStorage)
	assert.NoError(t, serviceError)
	assert.Equal(t, 1, countOfElementIn(testEmail, records))
}

func TestSubscribeEmailFailed(t *testing.T) {
	resetStorageFile()
	testEmail := "loremipsum@gmail.com"
	emailService := InitTestEmailService()
	_ = repository.WriteToStorage(testEmail, pathToStorage, 0)

	serviceError := emailService.SubscribeEmail(testEmail)

	records, _ := repository.ReadFromStorage(pathToStorage)
	assert.Error(t, serviceError, exception.ErrEmailIsAlreadySubscribed)
	assert.Equal(t, 1, countOfElementIn(testEmail, records))
}

func countOfElementIn(element string, in []string) int {
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
	if err := os.WriteFile(pathToStorage, emptyJSONArray, os.FileMode(permToOpenTheTestStorage)); err != nil {
		log.Fatal(err)
	}
}

func InitTestEmailService() *service.EmailServiceImpl {
	c := &config.Config{EmailStoragePath: pathToStorage}
	emailRepository := repository.NewEmailRepository(c.EmailStoragePath)
	emailSender := service.NewEmailSender(c)
	rateService := service.NewRateService(c)
	return service.NewEmailService(c, rateService, emailRepository, emailSender)
}
