package main

import (
	// "github.com/ivelsantos/cryptor/front"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services"
	"github.com/ivelsantos/cryptor/tui"
	"log"
)

func main() {
	err := models.InitDB("algor.db")
	if err != nil {
		log.Fatal(err)
	}

	go services.Services()

	if err = tui.Tui(); err != nil {
		log.Fatal(err)
	}
	// front.Front()
}
