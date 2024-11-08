package trading

import (
	"fmt"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/values"
	"log"
)

func Trading() error {
	err := models.EraseTesting()
	if err != nil {
		return err
	}

	for {
		algos, err := models.GetAllAlgos()
		if err != nil {
			return fmt.Errorf("Failed to get algos: %v", err)
		}

		for _, algo := range algos {
			price, err := values.GetPrice(algo.BaseAsset + algo.QuoteAsset)
			if err != nil {
				return err
			}

			// Placing the values on the globalStore
			optAlgo := lang.GlobalStore("Algo", algo)
			optPrice := lang.GlobalStore("Price", price)

			_, err = lang.Parse("", []byte(algo.Buycode), optPrice, optAlgo)
			if err != nil {
				log.Printf("%v: Parsing error: %v\n", algo.Name, err)
				continue
			}
		}
	}
}
