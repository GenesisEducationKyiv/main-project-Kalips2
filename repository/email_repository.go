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

func SaveEmailToStorage(email string, pathToStorage string) error {
	err := writeToStorage(email, pathToStorage)
	if err != nil {
		return errors.Wrap(err, failToSubscribeEmailMessage)
	}
	return err
}

func GetEmailsFromStorage(pathToStorage string) ([]string, error) {
	emails, err := readFromStorage(pathToStorage)
	if err != nil {
		return nil, errors.Wrap(err, failToReadStorageMessage)
	}
	return emails, err
}

func writeToStorage(email string, pathToStorage string) error {
	var err error
	emails, err := readFromStorage(pathToStorage)
	if err != nil {
		return err
	}
	emails = append(emails, email)

	data, err := json.Marshal(emails)
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

	var emails []string
	err = json.Unmarshal(data, &emails)
	if err != nil {
		return nil, err
	}

	return emails, nil
}

func CheckEmailIsExist(email string, pathToStorage string) (bool, error) {
	var err error
	emails, err := GetEmailsFromStorage(pathToStorage)
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
