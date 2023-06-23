package repository

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"os"
)

var (
	failToSaveEmailMessage = "Failed to subscribe email"
	failToGetEmailsMessage = "Failed to get emails from storage"
	permToOpenTheStorage   = 0644
)

type EmailRepositoryImpl struct {
	pathToStorage string
}

func (repo EmailRepositoryImpl) SaveEmailToStorage(email string) error {
	var err error
	var file *os.File
	defer func(file *os.File) {
		err = file.Close()
	}(file)

	file, err = setUpConnectionWithStorage(repo.pathToStorage)
	if err != nil {
		return errors.Wrap(err, failToSaveEmailMessage)
	}

	err = writeToStorage(email, file)
	if err != nil {
		return errors.Wrap(err, failToSaveEmailMessage)
	}
	return err
}

func (repo EmailRepositoryImpl) GetEmailsFromStorage() ([]string, error) {
	var file *os.File
	var err error
	var emails []string

	file, err = setUpConnectionWithStorage(repo.pathToStorage)
	if err != nil {
		return nil, errors.Wrap(err, failToGetEmailsMessage)
	}

	emails, err = readFromStorage(file)
	if err != nil {
		return nil, errors.Wrap(err, failToGetEmailsMessage)
	}
	return emails, err
}

func writeToStorage(email string, storage *os.File) error {
	writer := csv.NewWriter(storage)
	defer writer.Flush()
	err := writer.Write([]string{email})

	if err != nil {
		return errors.New("Failed to write into storage")
	}
	return err
}

func readFromStorage(storage *os.File) ([]string, error) {
	var data []string
	reader := csv.NewReader(storage)
	records, err := reader.ReadAll()

	if err != nil {
		return nil, errors.New("Failed to read from storage")
	}

	for _, record := range records {
		data = append(data, record[0])
	}
	return data, err
}

func setUpConnectionWithStorage(pathToStorage string) (*os.File, error) {
	storage, err := os.OpenFile(pathToStorage, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.FileMode(permToOpenTheStorage))
	if err != nil {
		return nil, errors.New("Failed to set up connection with storage")
	}
	return storage, err
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

func NewEmailRepository(pathToStorage string) EmailRepositoryImpl {
	return EmailRepositoryImpl{
		pathToStorage: pathToStorage,
	}
}
