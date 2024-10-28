package main

import (
	"github.com/ivelsantos/cryptor/front"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services"
	"log"
)

func main() {
	err := models.InitDB("algor.db")
	if err != nil {
		log.Fatal(err)
	}

	go services.Services()
	front.Front()
}
