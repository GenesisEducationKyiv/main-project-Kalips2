package repository

import (
	"btc-app/model"
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

func (repo *EmailRepositoryImpl) SaveEmail(email model.Email) error {
	err := WriteEmailToStorage(email, repo.pathToStorage, repo.permToOpenTheStorage)
	if err != nil {
		return errors.Wrap(err, message.FailToAddEmailToStorageMessage)
	}
	return err
}

func (repo *EmailRepositoryImpl) GetEmailsFromStorage() ([]model.Email, error) {
	emails, err := ReadEmailsFromStorage(repo.pathToStorage)
	if err != nil {
		return nil, errors.Wrap(err, message.FailToGetEmailsMessage)
	}
	return emails, err
}

func (repo *EmailRepositoryImpl) CheckEmailIsExist(email model.Email) (bool, error) {
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

func WriteEmailToStorage(email model.Email, pathToStorage string, permToFile int) error {
	var err error
	records, err := ReadEmailsFromStorage(pathToStorage)
	if err != nil {
		return exception.ErrWriteToStorage
	}
	records = append(records, email)

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

func ReadEmailsFromStorage(pathToStorage string) ([]model.Email, error) {
	data, err := os.ReadFile(pathToStorage)
	if err != nil {
		return nil, exception.ErrReadFromStorage
	}

	var records []model.Email
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
