package services

import (
	"log"

	"github.com/ivelsantos/cryptor/services/trading"
)

func Services() {
	err := trading.Trading()
	if err != nil {
		log.Fatal(err)
	}
}
