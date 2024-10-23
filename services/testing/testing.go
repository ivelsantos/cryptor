package testing

import (
	"fmt"
	"github.com/ivelsantos/cryptor/lang"
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
			optPrice := lang.GlobalStore("Price", price)
			optBotid := lang.GlobalStore("Botid", algo.Id)
			optTicket := lang.GlobalStore("Ticket", "BTCBRL")
			optOwner := lang.GlobalStore("Owner", algo.Owner)

			_, err := lang.Parse("", []byte(algo.Buycode), optPrice, optBotid, optTicket, optOwner)
			if err != nil {
				log.Printf("%v: Parsing error: %v\n", algo.Name, err)
				continue
			}
		}
	}
}
