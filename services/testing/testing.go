package testing

import (
	"fmt"
	"github.com/ivelsantos/cryptor/langSimulation"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt"
	"log"
)

func Testing() error {
	for {

		err := crypt.InitCrypt()
		algos, err := models.GetAlgosState("testing")
		if err != nil {
			return fmt.Errorf("Failed to get testing algor: %v", err)
		}

		for _, algo := range algos {
			price, ok := crypt.GetCryptValue("@Price")
			if !ok {
				return fmt.Errorf("Failed to get price")
			}

			// Placing the values on the globalStore
			optPrice := langSimulation.GlobalStore("Price", price)
			optBotid := langSimulation.GlobalStore("Botid", algo.Id)
			optTicket := langSimulation.GlobalStore("Ticket", "BTCBRL")
			optOwner := langSimulation.GlobalStore("Owner", algo.Owner)

			_, err := langSimulation.Parse("", []byte(algo.Buycode), optPrice, optBotid, optTicket, optOwner)
			if err != nil {
				log.Printf("%v: Parsing error: %v\n", algo.Name, err)
				continue
			}
		}
	}
}
