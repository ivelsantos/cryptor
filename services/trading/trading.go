package trading

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/testnet"
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

				err := updateOrderStatusTesting(algo)
				if err != nil {
					log.Printf("%v: Error updating order status: %v\n", algo.Name, err)
				}

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

func updateOrderStatusTesting(algo models.Algor) error {
	transactions, err := models.GetTestingPending(algo.Id)
	if err != nil {
		return err
	}

	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return fmt.Errorf("models.GetAccountByName: %v", err)
	}

	for _, transaction := range transactions {
		status, err := testnet.GetOrderStatus(account.ApiKey_test, account.SecretKey_test, algo.BaseAsset+algo.QuoteAsset, transaction.Orderid)
		if err != nil {
			return fmt.Errorf("testnet.GetOrderStatus: %v", err)
		}

		if status != transaction.Orderstatus {
			count := 0
			err := models.UpdateOrderStatus(status, transaction.Id)
			for err != nil && count < 100 {
				err = models.UpdateOrderStatus(status, transaction.Id)
				time.Sleep(250 * time.Millisecond)
				count += 1
			}

			if err != nil {
				return fmt.Errorf("models.UpdateOrderStatus: %v", err)
			}
		}
	}

	return nil
}
