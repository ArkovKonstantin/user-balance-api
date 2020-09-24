package main

import (
	"log"
	"user-balance-api/models"
	"user-balance-api/provider"
)

var (
	config models.Config
)

func init() {

}

func main() {
	p := provider.New(&config.SQLDataBase)
	err := p.Open()

	if err != nil {
		log.Fatal(err)
	}

}
