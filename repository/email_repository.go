package repository

import (
	"btc-app/template/exception"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

var permToOpenTheStorage = 0644

type EmailRepository interface {
	SaveEmailToStorage(email string) error
	GetEmailsFromStorage() ([]string, error)
	CheckEmailIsExist(email string) (bool, error)
}

type EmailRepositoryImpl struct {
	pathToStorage string
}

func (repo *EmailRepositoryImpl) SaveEmailToStorage(email string) error {
	err := WriteToStorage(email, repo.pathToStorage)
	if err != nil {
		return errors.Wrap(err, exception.FailToAddEmailToStorageMessage)
	}
	return err
}

func (repo *EmailRepositoryImpl) GetEmailsFromStorage() ([]string, error) {
	emails, err := ReadFromStorage(repo.pathToStorage)
	if err != nil {
		return nil, errors.Wrap(err, exception.FailToGetEmailsMessage)
	}
	return emails, err
}

func (repo *EmailRepositoryImpl) CheckEmailIsExist(email string) (bool, error) {
	var err error
	emails, err := repo.GetEmailsFromStorage()
	if err != nil {
		return false, errors.Wrap(err, exception.FailToCheckExistenceOfEmailMessage)
	}

	for _, existingEmail := range emails {
		if existingEmail == email {
			return true, err
		}
	}
	return false, err
}

func WriteToStorage(record string, pathToStorage string) error {
	var err error
	records, err := ReadFromStorage(pathToStorage)
	if err != nil {
		return exception.WriteToStorage
	}
	records = append(records, record)

	data, err := json.Marshal(records)
	if err != nil {
		return exception.WriteToStorage
	}

	err = os.WriteFile(pathToStorage, data, os.FileMode(permToOpenTheStorage))
	if err != nil {
		return exception.WriteToStorage
	}
	return err
}

func ReadFromStorage(pathToStorage string) ([]string, error) {
	data, err := os.ReadFile(pathToStorage)
	if err != nil {
		return nil, exception.ReadFromStorage
	}

	var records []string
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, exception.JsonWithIncorrectFormat
	}

	return records, nil
}

func NewEmailRepository(pathToStorage string) *EmailRepositoryImpl {
	return &EmailRepositoryImpl{
		pathToStorage: pathToStorage,
	}
}
