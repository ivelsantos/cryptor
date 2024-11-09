package trading

import (
	"fmt"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
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

			// Placing the values on the globalStore
			optAlgo := lang.GlobalStore("Algo", algo)

			_, err = lang.Parse("", []byte(algo.Buycode), optAlgo)
			if err != nil {
				log.Printf("%v: Parsing error: %v\n", algo.Name, err)
				continue
			}
		}
	}
}
