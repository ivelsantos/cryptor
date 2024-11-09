package trading

import (
	"fmt"
	"log"
	"sync"

	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
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

		var wg sync.WaitGroup

		for _, algo := range algos {
			wg.Add(1)
			go func(algo models.Algor) {
				defer wg.Done()

				// Placing the values on the globalStore
				optAlgo := lang.GlobalStore("Algo", algo)

				_, err = lang.Parse("", []byte(algo.Buycode), optAlgo)
				if err != nil {
					log.Printf("%v: Parsing error: %v\n", algo.Name, err)
				}
			}(algo)

		}
		wg.Wait()
	}
}
