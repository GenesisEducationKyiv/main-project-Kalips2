package config

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		Server      ServerConfig   `json:"server"`
		Crypto      CryptoConfig   `json:"crypto"`
		Database    DatabaseConfig `json:"database"`
		MailService MailConfig     `json:"emailService"`
	}

	ServerConfig struct {
		Port int `json:"port"`
	}

	CryptoConfig struct {
		CryptoCompareProviderURL string `json:"cryptoCompareProvider"`
		CoinMarketProviderURL    string `json:"coinMarketProvider"`
		CoinApiProviderURL       string `json:"coinApiProvider"`
		CurrencyFrom             string `json:"currencyFrom"`
		CurrencyTo               string `json:"currencyTo"`
	}

	DatabaseConfig struct {
		PathToStorage       string `json:"pathToStorage"`
		PermissionToStorage string `json:"permissionToStorage"`
	}

	MailConfig struct {
		MailHost     string `json:"mailServiceHost"`
		MailPort     int    `json:"mailPort"`
		MailSubject  string `json:"mailSubject"`
		MailFrom     string `json:"mailFrom"`
		MailPassword string `json:"mailPassword"`
	}
)

func NewConfig() (*Config, error) {
	conf := &Config{}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(data, conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}
