package main

import (
	"log"
	"producer/config"
	"producer/pkg/presentation/server"
)

func main() {
	conf, err := config.NewConfig("config.json")
	if err != nil {
		log.Fatal("Failed to initialize configuration.", err)
	}

	curServer := server.NewServer(conf)
	curServer.SetupServer(conf)
}
