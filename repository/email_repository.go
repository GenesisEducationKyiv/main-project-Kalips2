package repository

import (
	"btc-app/template/exception"
	"btc-app/template/message"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

var permToOpenTheStorage = 0644

type EmailRepositoryImpl struct {
	pathToStorage        string
	permToOpenTheStorage int
}

func (repo *EmailRepositoryImpl) SaveEmailToStorage(email string) error {
	err := WriteToStorage(email, repo.pathToStorage, repo.permToOpenTheStorage)
	if err != nil {
		return errors.Wrap(err, message.FailToAddEmailToStorageMessage)
	}
	return err
}

func (repo *EmailRepositoryImpl) GetEmailsFromStorage() ([]string, error) {
	emails, err := ReadFromStorage(repo.pathToStorage)
	if err != nil {
		return nil, errors.Wrap(err, message.FailToGetEmailsMessage)
	}
	return emails, err
}

func (repo *EmailRepositoryImpl) CheckEmailIsExist(email string) (bool, error) {
	var err error
	emails, err := repo.GetEmailsFromStorage()
	if err != nil {
		return false, errors.Wrap(err, message.FailToCheckExistenceOfEmailMessage)
	}

	for _, existingEmail := range emails {
		if existingEmail == email {
			return true, err
		}
	}
	return false, err
}

func WriteToStorage(record string, pathToStorage string, permToFile int) error {
	var err error
	records, err := ReadFromStorage(pathToStorage)
	if err != nil {
		return exception.ErrWriteToStorage
	}
	records = append(records, record)

	data, err := json.Marshal(records)
	if err != nil {
		return exception.ErrWriteToStorage
	}

	err = os.WriteFile(pathToStorage, data, os.FileMode(permToFile))
	if err != nil {
		return exception.ErrWriteToStorage
	}
	return err
}

func ReadFromStorage(pathToStorage string) ([]string, error) {
	data, err := os.ReadFile(pathToStorage)
	if err != nil {
		return nil, exception.ErrReadFromStorage
	}

	var records []string
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, exception.ErrJsonWithIncorrectFormat
	}

	return records, nil
}

func NewEmailRepository(pathToStorage string) *EmailRepositoryImpl {
	return &EmailRepositoryImpl{
		pathToStorage:        pathToStorage,
		permToOpenTheStorage: permToOpenTheStorage,
	}
}
