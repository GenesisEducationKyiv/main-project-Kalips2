package functional

import (
	"btc-app/config"
	"btc-app/handler"
	"btc-app/service"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetRate(t *testing.T) {
	rateHandler := InitTestRateHandler()
	server := httptest.NewServer(rateHandler.GetCurrentRateHandler())
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(t, err, "Failed to send request")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code, expected %d but got %d", http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Failed to read response body")

	rate, err := strconv.ParseFloat(string(body), 64)
	assert.NoError(t, err, "failed to parse response body as float")

	assert.True(t, !math.IsNaN(rate) && !math.IsInf(rate, 0), "unexpected response body, expected a valid floating-point number")
}

func InitTestRateHandler() *handler.RateHandlerImpl {
	conf := createConfig()
	rateService := service.NewRateService(conf)
	return handler.NewRateHandler(conf, rateService)
}

func createConfig() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load env variables from .env file.", err)
	}

	var c config.Config
	if err := c.NewConfig(); err != nil {
		log.Fatal("Failed to initialize configuration.", err)
	}
	return &c
}
