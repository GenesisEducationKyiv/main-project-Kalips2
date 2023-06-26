package main

import (
	"btc-app/config"
	"btc-app/server"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load env variables from .env file.", err)
	}

	var conf config.Config
	if err := conf.NewConfig(); err != nil {
		log.Fatal("Failed to initialize configuration.", err)
	}

	var curServer = server.NewServer(conf)
	curServer.NewHandlers(&conf)
}
