package operations

import (
	"fmt"
	"log"
	"time"

	"github.com/ivelsantos/cryptor/models"
)

func Buy(algo models.Algor, ticket string, price float64) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTesting(algo.Id)
		if err != nil {
			return err
		}
		if len(transactions) != 0 {
			return nil
		}

		account, err := models.GetAccountByName(algo.Owner)
		_ = account

		current := int(time.Now().Unix())

		err = models.InsertTestingBuy(algo.Id, ticket, price, current)
		if err != nil {
			return err
		}

		log.Printf("TESTING: Buy %s at price %v\n", ticket, price)
		return nil
	case "new", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func Sell(algo models.Algor, price float64) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTesting(algo.Id)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			current := int(time.Now().Unix())

			err = models.InsertTestingSell(transaction.Id, price, current)
			if err != nil {
				return err
			}
			log.Printf("TESTING: Sell %s at price %v\n", transaction.Ticket, price)
		}

		return nil
	case "new", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func StopLoss(algo models.Algor, stop float64, price float64) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTesting(algo.Id)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			sellPrice := transaction.Buyprice - (stop * transaction.Buyprice)
			if price <= sellPrice {
				current := int(time.Now().Unix())
				err = models.InsertTestingSell(transaction.Id, price, current)
				if err != nil {
					return err
				}
				log.Printf("TESTING: Sell %s at price %v\n", transaction.Ticket, price)
			}

		}

		return nil
	case "new", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func TakeProfit(algo models.Algor, take float64, price float64) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTesting(algo.Id)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			sellPrice := transaction.Buyprice + (take * transaction.Buyprice)
			if price >= sellPrice {
				current := int(time.Now().Unix())
				err = models.InsertTestingSell(transaction.Id, price, current)
				if err != nil {
					return err
				}
				log.Printf("TESTING: Sell %s at price %v\n", transaction.Ticket, price)
			}
		}

		return nil
	case "new", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}
