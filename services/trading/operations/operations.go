package operations

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/testnet"
	"github.com/ivelsantos/cryptor/services/crypt/values"
)

func Buy(algo models.Algor) error {
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

		asset, err := testnet.GetAccountQuote(account.ApiKey_test, account.SecretKey_test, algo.QuoteAsset)
		if err != nil {
			return err
		}
		asset_float, err := strconv.ParseFloat(asset, 64)
		if err != nil {
			return err
		}

		minNotional, err := testnet.GetMinNotional(account.ApiKey_test, account.SecretKey_test, algo.BaseAsset+algo.QuoteAsset)
		if err != nil {
			return err
		}

		minOrder := minNotional * 4

		if minOrder > asset_float {
			quoteOrder := minNotional * 2
			quoteOrderStr := strconv.FormatFloat(quoteOrder, 'f', -1, 64)

			_, err := testnet.SellQuote(account.ApiKey_test, account.SecretKey_test, algo.BaseAsset+algo.QuoteAsset, quoteOrderStr)
			if err != nil {
				return err
			}
		}

		minOrderStr := strconv.FormatFloat(minOrder, 'f', -1, 64)

		order, err := testnet.Buy(account.ApiKey_test, account.SecretKey_test, algo.BaseAsset+algo.QuoteAsset, minOrderStr)
		if err != nil {
			return err
		}

		tb := models.TestingBuy{Botid: algo.Id, Baseasset: algo.BaseAsset, Quoteasset: algo.QuoteAsset}
		tb.Orderid = int(order.OrderID)
		tb.Orderstatus = string(order.Status)

		cum, err := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
		if err != nil {
			return err
		}
		tb.Buyvalue = cum

		quant, err := strconv.ParseFloat(order.ExecutedQuantity, 64)
		if err != nil {
			return err
		}
		tb.Buyquantity = quant
		tb.Buytime = int(order.TransactTime)

		err = models.InsertTestingBuy(tb)
		if err != nil {
			return err
		}

		log.Printf("TESTING %v: Buy %s at price %v\n", algo.Name, algo.BaseAsset+algo.QuoteAsset, cum/quant)
		return nil
	case "waiting", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func Sell(algo models.Algor, price float64) error {
	switch algo.State {
	case "testing":
		// transactions, err := models.GetTesting(algo.Id)
		// if err != nil {
		// 	return err
		// }

		// for _, transaction := range transactions {
		// 	// current := int(time.Now().Unix())

		// 	ts := models.TestingSell{Entryid: transaction.Id}
		// 	// ts.Orderstatus = ...
		// 	// ts.Sellprice = ...
		// 	// ts.Sellquantity = ...
		// 	// ts.Selltime = ...

		// 	err = models.InsertTestingSell(ts)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	log.Printf("TESTING: Sell %s at price %v\n", transaction.Ticket, price)
		// }

		return nil
	case "waiting", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func StopLoss(algo models.Algor, stop float64) error {
	ticket := algo.BaseAsset + algo.QuoteAsset
	price, err := values.GetPrice(ticket)
	if err != nil {
		return err
	}

	switch algo.State {
	case "testing":
		transactions, err := models.GetTesting(algo.Id)
		if err != nil {
			return err
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			// if transaction.Orderstatus != "FILLED" {
			// 	continue
			// }
			buyprice := transaction.Buyvalue / transaction.Buyquantity
			sellPrice := buyprice - (stop * buyprice)

			if price <= sellPrice {

				quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
				order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
				if err != nil {
					return err
				}

				ts := models.TestingSell{Entryid: transaction.Id}
				ts.Orderstatus = string(order.Status)

				cum, err := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
				if err != nil {
					return err
				}
				ts.Sellvalue = cum
				ts.Selltime = int(order.TransactTime)

				err = models.InsertTestingSell(ts)
				if err != nil {
					return err
				}
				log.Printf("TESTING %v: StopLoss %s at price %v\n", algo.Name, transaction.Ticket, cum/transaction.Buyquantity)
			}

		}

		return nil
	case "waiting", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func TakeProfit(algo models.Algor, take float64) error {
	ticket := algo.BaseAsset + algo.QuoteAsset
	price, err := values.GetPrice(ticket)
	if err != nil {
		return err
	}

	switch algo.State {
	case "testing":
		transactions, err := models.GetTesting(algo.Id)
		if err != nil {
			return err
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			// if transaction.Orderstatus != "FILLED" {
			// 	continue
			// }
			buyprice := transaction.Buyvalue / transaction.Buyquantity
			sellPrice := buyprice + (take * buyprice)

			if price > sellPrice {

				quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
				order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
				if err != nil {
					return err
				}

				ts := models.TestingSell{Entryid: transaction.Id}
				ts.Orderstatus = string(order.Status)

				cum, err := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
				if err != nil {
					return err
				}
				ts.Sellvalue = cum
				ts.Selltime = int(order.TransactTime)

				err = models.InsertTestingSell(ts)
				if err != nil {
					return err
				}
				log.Printf("TESTING %v: TakeProfit %s at price %v\n", algo.Name, transaction.Ticket, cum/transaction.Buyquantity)
			}

		}
		return nil
	case "waiting", "live":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
