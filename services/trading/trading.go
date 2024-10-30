package trading

import (
	"fmt"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt"
	"log"
)

func Trading() error {
	for {
		log.Println("PASSOU AQUI")
		err := crypt.InitCrypt()
		algos, err := models.GetAllAlgos()
		if err != nil {
			return fmt.Errorf("Failed to get algos: %v", err)
		}

		for _, algo := range algos {
			price, ok := crypt.GetCryptValue("@Price")
			if !ok {
				return fmt.Errorf("Failed to get price")
			}

			// Placing the values on the globalStore
			optAlgo := lang.GlobalStore("Algo", algo)
			optPrice := lang.GlobalStore("Price", price)
			optBase := lang.GlobalStore("Base", "BTC")
			optQuote := lang.GlobalStore("Quote", "BRL")

			_, err := lang.Parse("", []byte(algo.Buycode), optPrice, optBase, optQuote, optAlgo)
			if err != nil {
				log.Printf("%v: Parsing error: %v\n", algo.Name, err)
				continue
			}
		}
	}
}
