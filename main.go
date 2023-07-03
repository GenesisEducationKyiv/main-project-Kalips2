package main

import (
	"btc-app/config"
	"btc-app/server"
	"log"
)

func main() {
	conf, err := config.NewConfig("config.json")
	if err != nil {
		log.Fatal("Failed to initialize configuration.", err)
	}

	var curServer = server.NewServer(conf)
	curServer.SetupServer(conf)
}
