package main

import (
	"log"
	"os"
	"user-balance-api/application"
	"user-balance-api/models"
	"user-balance-api/provider"
	"user-balance-api/repository"
)

var (
	config models.Config
)

func init() {
	env := os.Getenv("ENV")
	var path string
	if env == "" || env == "dev" {
		path = "config/config.dev.toml"
	} else if env == "prod" {
		path = "config/config.prod.toml"
	}
	err := config.LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	p := provider.New(&config.SQLDataBase)
	err := p.Open()

	if err != nil {
		log.Fatal(err)
	}

	rep := repository.NewAccountRepository(p)

	app := application.New(rep)
	app.Start()

}
