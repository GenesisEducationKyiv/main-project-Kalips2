package functional

import (
	"btc-app/config"
	"btc-app/handler"
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
	rateHandler := handler.NewRateHandler(createConfig())
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

func createConfig() *config.Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load env variables from .env file.", err)
	}

	var conf config.Config
	if err := conf.InitConfigFromEnv(); err != nil {
		log.Fatal("Failed to initialize configuration.", err)
	}
	return &conf
}
