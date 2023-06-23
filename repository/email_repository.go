package repository

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

var (
	failToSubscribeEmailMessage = "Failed to subscribe email"
	failToReadStorageMessage    = "Failed to read from storage"
	permToOpenTheStorage        = 0644
)

type EmailRepository interface {
	SaveEmailToStorage(email string) error
	GetEmailsFromStorage() ([]string, error)
	CheckEmailIsExist(email string) (bool, error)
}

type EmailRepositoryImpl struct {
	pathToStorage string
}

func (repo *EmailRepositoryImpl) SaveEmailToStorage(email string) error {
	err := writeToStorage(email, repo.pathToStorage)
	if err != nil {
		return errors.Wrap(err, failToSubscribeEmailMessage)
	}
	return err
}

func (repo *EmailRepositoryImpl) GetEmailsFromStorage() ([]string, error) {
	emails, err := readFromStorage(repo.pathToStorage)
	if err != nil {
		return nil, errors.Wrap(err, failToReadStorageMessage)
	}
	return emails, err
}

func (repo *EmailRepositoryImpl) CheckEmailIsExist(email string) (bool, error) {
	var err error
	emails, err := repo.GetEmailsFromStorage()
	if err != nil {
		return false, errors.Wrap(err, "Failed to check the existence of email")
	}

	for _, existingEmail := range emails {
		if existingEmail == email {
			return true, err
		}
	}
	return false, err
}

func writeToStorage(record string, pathToStorage string) error {
	var err error
	records, err := readFromStorage(pathToStorage)
	if err != nil {
		return err
	}
	records = append(records, record)

	data, err := json.Marshal(records)
	if err != nil {
		return err
	}

	err = os.WriteFile(pathToStorage, data, os.FileMode(permToOpenTheStorage))
	if err != nil {
		return err
	}

	return err
}

func readFromStorage(pathToStorage string) ([]string, error) {
	data, err := os.ReadFile(pathToStorage)
	if err != nil {
		return nil, errors.New(failToReadStorageMessage)
	}

	var records []string
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func NewEmailRepository(pathToStorage string) *EmailRepositoryImpl {
	return &EmailRepositoryImpl{
		pathToStorage: pathToStorage,
	}
}
