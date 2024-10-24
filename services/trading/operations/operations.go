package operations

import (
	"fmt"
	"log"
	"time"

	"github.com/ivelsantos/cryptor/models"
)

func Buy(botid int, ticket string, price float64, mode string) error {
	switch mode {
	case "testing":
		algos, err := models.GetTesting(botid)
		if err != nil {
			return err
		}
		if len(algos) != 0 {
			return nil
		}
		current := int(time.Now().Unix())

		err = models.InsertTestingBuy(botid, ticket, price, current)
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

func Sell(botid int, price float64) error {
	algos, err := models.GetTesting(botid)
	if err != nil {
		return err
	}

	for _, algo := range algos {
		current := int(time.Now().Unix())

		err = models.InsertTestingSell(algo.Id, price, current)
		if err != nil {
			return err
		}
		log.Printf("TESTING: Sell %s at price %v\n", algo.Ticket, price)
	}

	return nil
}

func StopLoss(stop float64, botid int, price float64) error {
	algos, err := models.GetTesting(botid)
	if err != nil {
		return err
	}

	for _, algo := range algos {
		sellPrice := algo.Buyprice - (stop * algo.Buyprice)
		if price <= sellPrice {
			current := int(time.Now().Unix())
			err = models.InsertTestingSell(algo.Id, price, current)
			if err != nil {
				return err
			}
			log.Printf("TESTING: Sell %s at price %v\n", algo.Ticket, price)
		}

	}

	return nil
}

func TakeProfit(take float64, botid int, price float64) error {
	algos, err := models.GetTesting(botid)
	if err != nil {
		return err
	}

	for _, algo := range algos {
		sellPrice := algo.Buyprice + (take * algo.Buyprice)
		if price >= sellPrice {
			current := int(time.Now().Unix())
			err = models.InsertTestingSell(algo.Id, price, current)
			if err != nil {
				return err
			}
			log.Printf("TESTING: Sell %s at price %v\n", algo.Ticket, price)
		}
	}

	return nil
}
