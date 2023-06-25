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
	c := &config.Config{EmailStoragePath: pathToStorage}
	emailService := service.NewEmailService(c)

	serviceError := emailService.SubscribeEmail(testEmail)

	records, _ := repository.ReadFromStorage(pathToStorage)
	assert.NoError(t, serviceError)
	assert.Equal(t, 1, countOfElementIn(testEmail, records))
	resetStorageFile()
}

func TestSubscribeEmailFailed(t *testing.T) {
	resetStorageFile()
	testEmail := "loremipsum@gmail.com"
	_ = repository.WriteToStorage(testEmail, pathToStorage)
	emailService := service.NewEmailService(&config.Config{EmailStoragePath: pathToStorage})

	serviceError := emailService.SubscribeEmail(testEmail)

	records, _ := repository.ReadFromStorage(pathToStorage)
	assert.Error(t, serviceError, exception.ErrEmailIsAlreadySubscribed)
	assert.Equal(t, 1, countOfElementIn(testEmail, records))
	resetStorageFile()
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
