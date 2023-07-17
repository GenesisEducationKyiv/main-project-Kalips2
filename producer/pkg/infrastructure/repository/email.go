package repository

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"producer/config"
	"producer/pkg/domain/model"
	"producer/pkg/domain/service"
	cerror "producer/template/cerror"
	"producer/template/message"
	"strconv"
)

type EmailRepositoryImpl struct {
	pathToStorage        string
	permToOpenTheStorage string
	logger               service.Logger
}

func (repo *EmailRepositoryImpl) SaveEmail(email model.Email) error {
	err := WriteEmailToStorage(email, repo.pathToStorage, repo.permToOpenTheStorage)
	if err != nil {
		err := errors.Wrap(err, fmt.Sprintf(message.FailToAddEmailToStorageMessage, email.GetAddress()))
		_ = repo.logger.LogError(err.Error())
		return err
	}
	return err
}

func (repo *EmailRepositoryImpl) GetEmailsFromStorage() ([]model.Email, error) {
	emails, err := ReadEmailsFromStorage(repo.pathToStorage)
	if err != nil {
		err := errors.Wrap(err, message.FailToGetEmailsMessage)
		_ = repo.logger.LogError(err.Error())
		return nil, err
	}
	return emails, err
}

func (repo *EmailRepositoryImpl) CheckEmailIsExist(email model.Email) (bool, error) {
	var err error
	emails, err := repo.GetEmailsFromStorage()
	if err != nil {
		err := errors.Wrap(err, fmt.Sprintf(message.FailToCheckExistenceOfEmailMessage, email.GetAddress()))
		_ = repo.logger.LogError(err.Error())
		return false, err
	}

	for _, existingEmail := range emails {
		if existingEmail.GetAddress() == email.GetAddress() {
			return true, err
		}
	}
	return false, err
}

func WriteEmailToStorage(email model.Email, pathToStorage string, permToFile string) error {
	var err error
	records, err := ReadEmailsFromStorage(pathToStorage)
	if err != nil {
		return cerror.ErrWriteToStorage
	}
	records = append(records, email)

	data, err := json.Marshal(records)
	if err != nil {
		return cerror.ErrWriteToStorage
	}

	permission, _ := strconv.ParseInt(permToFile, 0, 32)
	err = os.WriteFile(pathToStorage, data, os.FileMode(permission))
	if err != nil {
		return cerror.ErrWriteToStorage
	}
	return err
}

func ReadEmailsFromStorage(pathToStorage string) ([]model.Email, error) {
	data, err := os.ReadFile(pathToStorage)
	if err != nil {
		return nil, cerror.ErrReadFromStorage
	}

	var records []model.Email
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, cerror.ErrJsonWithIncorrectFormat
	}

	return records, nil
}

func NewEmailRepository(c config.DatabaseConfig, logger service.Logger) *EmailRepositoryImpl {
	return &EmailRepositoryImpl{
		pathToStorage:        c.PathToStorage,
		permToOpenTheStorage: c.PermissionToStorage,
		logger:               logger,
	}
}
