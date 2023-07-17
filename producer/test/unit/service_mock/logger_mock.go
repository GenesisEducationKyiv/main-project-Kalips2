package service_mock

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"log"
)

type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) LogDebug(v ...any) error {
	log.Printf("Debug: " + fmt.Sprint(v...))
	return nil
}

func (l *MockLogger) LogError(v ...any) error {
	log.Printf("Error: " + fmt.Sprint(v...))
	return nil
}

func (l *MockLogger) LogInfo(v ...any) error {
	log.Printf("Info: " + fmt.Sprint(v...))
	return nil
}
