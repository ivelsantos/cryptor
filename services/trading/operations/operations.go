package operations

import (
	"fmt"
	"log"
	// "time"

	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/testnet"
)

func Buy(algo models.Algor, base_asset string, quote_asset string, price float64) error {
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
		brl := testnet.GetAccountQuote(account.ApiKey_test, account.SecretKey_test)

		tb := models.TestingBuy{Botid: algo.Id, Baseasset: base_asset, Quoteasset: quote_asset}
		// tb.Orderid = ...
		// tb.Orderstatus = ...
		// tb.Buyprice = ...
		// tb.Buyquantity = ...
		// tb.Buytime = ...

		err = models.InsertTestingBuy(tb)
		if err != nil {
			return err
		}

		log.Printf("TESTING: Buy %s at price %v\n", base_asset+quote_asset, price)
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
			// current := int(time.Now().Unix())

			ts := models.TestingSell{Entryid: transaction.Id}
			// ts.Orderstatus = ...
			// ts.Sellprice = ...
			// ts.Sellquantity = ...
			// ts.Selltime = ...

			err = models.InsertTestingSell(ts)
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
				// current := int(time.Now().Unix())

				ts := models.TestingSell{Entryid: transaction.Id}
				// ts.Orderstatus = ...
				// ts.Sellprice = ...
				// ts.Sellquantity = ...
				// ts.Selltime = ...
				err = models.InsertTestingSell(ts)
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
				// current := int(time.Now().Unix())

				ts := models.TestingSell{Entryid: transaction.Id}
				// ts.Orderstatus = ...
				// ts.Sellprice = ...
				// ts.Sellquantity = ...
				// ts.Selltime = ...

				err = models.InsertTestingSell(ts)
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
