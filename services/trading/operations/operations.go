package operations

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/testingnet"
	"github.com/ivelsantos/cryptor/services/crypt/testnet"
)

func Buy(algo models.Algor) (bool, error) {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTestingBuy(algo.Id)
		if err != nil {
			return false, err
		}
		if len(transactions) != 0 {
			return false, nil
		}

		ticket := algo.BaseAsset + algo.QuoteAsset
		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return false, fmt.Errorf("GetAccountByName: %v", err)
		}
		_, askPrice, err := testingnet.GetDepth(account.ApiKey, account.SecretKey, ticket)
		if err != nil {
			return false, fmt.Errorf("GetDepth: %v", err)
		}

		tb := models.TestingBuy{Botid: algo.Id, Baseasset: algo.BaseAsset, Quoteasset: algo.QuoteAsset}
		tb.Buyvalue = askPrice
		tb.Buytime = time.Now().Unix()

		err = models.InsertTestingBuy(tb)
		if err != nil {
			return false, fmt.Errorf("InsertTestingBuy: %v", err)
		}
		return true, nil
	case "live":
		transactions, err := models.GetTestingBuy(algo.Id)
		if err != nil {
			return false, err
		}

		if len(transactions) != 0 {
			return false, nil
		}

		ticket := algo.BaseAsset + algo.QuoteAsset

		account, err := models.GetAccountByName(algo.Owner)

		asset, err := testnet.GetAccountQuote(account.ApiKey_test, account.SecretKey_test, algo.QuoteAsset)
		if err != nil {
			return false, err
		}
		asset_float, err := strconv.ParseFloat(asset, 64)
		if err != nil {
			return false, err
		}

		minNotional, err := testnet.GetMinNotional(account.ApiKey_test, account.SecretKey_test, ticket)
		if err != nil {
			return false, err
		}

		minOrder := minNotional * 4

		if minOrder > asset_float {
			quoteOrder := minOrder * 2
			quoteOrderStr := strconv.FormatFloat(quoteOrder, 'f', -1, 64)

			_, err := testnet.SellQuote(account.ApiKey_test, account.SecretKey_test, ticket, quoteOrderStr)
			if err != nil {
				return false, err
			}
		}

		minOrderStr := strconv.FormatFloat(minOrder, 'f', -1, 64)

		order, err := testnet.Buy(account.ApiKey_test, account.SecretKey_test, ticket, minOrderStr)
		if err != nil {
			return false, err
		}

		//// Making sure the order is fulfilled
		time.Sleep(2 * time.Second) // wait for two seconds to make sure the transaction is done

		updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, ticket, int(order.OrderID))
		if err != nil {
			return false, fmt.Errorf("testnet.GetOrder: %v", err)
		}

		count := 0
		for string(updatedOrder.Status) != "FILLED" && count < 30 {
			time.Sleep(2 * time.Second)
			updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, ticket, int(order.OrderID))
			if err != nil {
				return false, fmt.Errorf("testnet.GetOrder: %v", err)
			}
			count++
		}

		if string(updatedOrder.Status) != "FILLED" {
			log.Fatalf("Buy order %v not fulfilled\n", order.OrderID)
		}

		tb := models.TransactionBuy{Botid: algo.Id, Baseasset: algo.BaseAsset, Quoteasset: algo.QuoteAsset}
		tb.Orderid = int(order.OrderID)
		tb.Orderstatus = string(updatedOrder.Status)

		cum, err := strconv.ParseFloat(updatedOrder.CummulativeQuoteQuantity, 64)
		if err != nil {
			return false, err
		}
		tb.Buyvalue = cum

		quant, err := strconv.ParseFloat(updatedOrder.ExecutedQuantity, 64)
		if err != nil {
			return false, err
		}
		tb.Buyquantity = quant
		tb.Buytime = int(order.TransactTime)

		err = models.InsertTransactionBuy(tb)
		count = 0

		for err != nil && count < 100 {
			err = models.InsertTransactionBuy(tb)
			time.Sleep(250 * time.Millisecond)
			count += 1
		}

		if err != nil {
			return false, fmt.Errorf("InsertTestingBuy: %v", err)
		}
		return true, nil
	case "waiting":
		return false, nil
	case "verification":
		return false, nil
	default:
		return false, fmt.Errorf("Unknown mode\n")
	}
}

func Sell(algo models.Algor) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTransactionSell(algo.Id)
		if err != nil {
			return fmt.Errorf("models.GetTesting: %v", err)
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return fmt.Errorf("models.GetAccountByName: %v", err)
		}

		for _, transaction := range transactions {
			quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
			order, err := testingnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
			if err != nil {
				return fmt.Errorf("testnet.Sell: %v", err)
			}

			//// Making sure the order is fulfilled
			time.Sleep(2 * time.Second) // wait for two seconds to make sure the transaction is done

			updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
			if err != nil {
				return fmt.Errorf("testnet.GetOrder: %v", err)
			}

			count := 0
			for string(updatedOrder.Status) != "FILLED" && count < 30 {
				time.Sleep(2 * time.Second)
				updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
				if err != nil {
					return fmt.Errorf("testnet.GetOrder: %v", err)
				}
				count++
			}

			if string(updatedOrder.Status) != "FILLED" {
				log.Fatalf("Sell order %v not fulfilled\n", order.OrderID)
			}

			ts := models.TransactionSell{Entryid: transaction.Id, Orderid: transaction.Orderid}
			ts.Orderstatus = string(updatedOrder.Status)

			cum, err := strconv.ParseFloat(updatedOrder.CummulativeQuoteQuantity, 64)
			if err != nil {
				return err
			}
			ts.Sellvalue = cum
			ts.Selltime = int(order.TransactTime)

			err = models.InsertTransactionSell(ts)
			count = 0

			for err != nil && count < 100 {
				err = models.InsertTransactionSell(ts)
				count += 1
			}
			if err != nil {
				return fmt.Errorf("models.InsertTransactionSell: %v", err)
			}

			// res := (ts.Sellvalue - transaction.Buyvalue) / transaction.Buyvalue

			// log.Printf("TESTING %v Sell: \tMargin %v\tBuyvalue %v\tSellvalue %v\n", algo.Name, roundFloat(res, 5), transaction.Buyvalue, ts.Sellvalue)
		}
		return nil
	case "waiting", "live":
		return nil
	case "verification":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func StopLoss(algo models.Algor, stop float64) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTestingSell(algo.Id)
		if err != nil {
			return fmt.Errorf("models.GetTesting: %v", err)
		}
		if len(transactions) < 1 {
			return nil
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return fmt.Errorf("models.GetAccountByName: %v", err)
		}
		ticker := algo.BaseAsset + algo.QuoteAsset

		bidPrice, _, err := testingnet.GetDepth(account.ApiKey_test, account.SecretKey_test, ticker)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			buyprice := transaction.Buyvalue
			sellPrice := buyprice - (stop * buyprice)

			if bidPrice <= sellPrice {
				ts := models.TestingSell{Entryid: transaction.Id, Sellvalue: bidPrice, Selltime: time.Now().Unix()}

				err = models.InsertTestingSell(ts)
				if err != nil {
					return fmt.Errorf("models.InsertTestingSell: %v", err)
				}

				// res := (ts.Sellvalue - transaction.Buyvalue) / transaction.Buyvalue

				// log.Printf("TESTING %v Stop_loss: \tMargin %v\tBuyvalue %v\tSellvalue %v\n", algo.Name, roundFloat(res, 5), transaction.Buyvalue, ts.Sellvalue)
			}

		}
		return nil
	case "live":
		transactions, err := models.GetTransactionSell(algo.Id)
		if err != nil {
			return fmt.Errorf("models.GetTesting: %v", err)
		}
		if len(transactions) < 1 {
			return nil
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return fmt.Errorf("models.GetAccountByName: %v", err)
		}

		ticker := algo.BaseAsset + algo.QuoteAsset

		bidPrice, askPrice, err := testnet.GetDepth(account.ApiKey_test, account.SecretKey_test, ticker)
		if err != nil {
			return err
		}

		var price float64
		if askPrice > bidPrice {
			price = askPrice
		} else {
			price = bidPrice
		}

		for _, transaction := range transactions {
			buyprice := transaction.Buyvalue / transaction.Buyquantity
			sellPrice := buyprice - (stop * buyprice)

			if price <= sellPrice {
				quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
				order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
				if err != nil {
					return fmt.Errorf("testnet.Sell: %v", err)
				}

				//// Making sure the order is fulfilled
				time.Sleep(2 * time.Second) // wait for two seconds to make sure the transaction is done

				updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
				if err != nil {
					return fmt.Errorf("testnet.GetOrder: %v", err)
				}

				count := 0
				for string(updatedOrder.Status) != "FILLED" && count < 30 {
					time.Sleep(2 * time.Second)
					updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
					if err != nil {
						return fmt.Errorf("testnet.GetOrder: %v", err)
					}
					count++
				}

				if string(updatedOrder.Status) != "FILLED" {
					log.Fatalf("StopLoss order %v not fulfilled\n", order.OrderID)
				}

				ts := models.TransactionSell{Entryid: transaction.Id, Orderid: transaction.Orderid}
				ts.Orderstatus = string(updatedOrder.Status)

				cum, err := strconv.ParseFloat(updatedOrder.CummulativeQuoteQuantity, 64)
				if err != nil {
					return fmt.Errorf("ParseFloat: %v", err)
				}
				ts.Sellvalue = cum
				ts.Selltime = int(order.TransactTime)

				err = models.InsertTransactionSell(ts)
				count = 0
				for err != nil && count < 100 {
					err = models.InsertTransactionSell(ts)
					count += 1
				}
				if err != nil {
					return fmt.Errorf("models.InsertTransactionSell: %v", err)
				}

				// res := (ts.Sellvalue - transaction.Buyvalue) / transaction.Buyvalue

				// log.Printf("TESTING %v Stop_loss: \tMargin %v\tBuyvalue %v\tSellvalue %v\n", algo.Name, roundFloat(res, 5), transaction.Buyvalue, ts.Sellvalue)
			}

		}
		return nil
	case "waiting":
		return nil
	case "verification":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}

}

func TakeProfit(algo models.Algor, take float64) error {
	switch algo.State {
	case "testing":
		transactions, err := models.GetTestingSell(algo.Id)
		if err != nil {
			return fmt.Errorf("GetTestingSell: %v", err)
		}
		if len(transactions) < 1 {
			return nil
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return fmt.Errorf("GetAccountByName: %v", err)
		}

		ticker := algo.BaseAsset + algo.QuoteAsset

		bidPrice, _, err := testingnet.GetDepth(account.ApiKey_test, account.SecretKey_test, ticker)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			buyprice := transaction.Buyvalue
			sellPrice := buyprice + (take * buyprice)

			if bidPrice > sellPrice {
				ts := models.TestingSell{Entryid: transaction.Id, Sellvalue: bidPrice, Selltime: time.Now().Unix()}
				err = models.InsertTestingSell(ts)
				if err != nil {
					return fmt.Errorf("InsertTestingSell: %v", err)
				}

				// res := (ts.Sellvalue - transaction.Buyvalue) / transaction.Buyvalue

				// log.Printf("TESTING %v Take_profit: \tMargin %v\tBuyvalue %v\tSellvalue %v\n", algo.Name, roundFloat(res, 5), transaction.Buyvalue, ts.Sellvalue)
			}

		}
		return nil
	case "live":
		transactions, err := models.GetTransactionSell(algo.Id)
		if err != nil {
			return err
		}
		if len(transactions) < 1 {
			return nil
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return err
		}

		ticker := algo.BaseAsset + algo.QuoteAsset

		bidPrice, _, err := testnet.GetDepth(account.ApiKey_test, account.SecretKey_test, ticker)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			buyprice := transaction.Buyvalue / transaction.Buyquantity
			sellPrice := buyprice + (take * buyprice)

			if bidPrice > sellPrice {
				quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
				order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
				if err != nil {
					return err
				}

				//// Making sure the order is fulfilled
				time.Sleep(2 * time.Second) // wait for two seconds to make sure the transaction is done

				updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
				if err != nil {
					return fmt.Errorf("testnet.GetOrder: %v", err)
				}

				count := 0
				for string(updatedOrder.Status) != "FILLED" && count < 30 {
					time.Sleep(2 * time.Second)
					updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
					if err != nil {
						return fmt.Errorf("testnet.GetOrder: %v", err)
					}
					count++
				}

				if string(updatedOrder.Status) != "FILLED" {
					log.Fatalf("TakeProfit order %v not fulfilled\n", order.OrderID)
				}

				ts := models.TransactionSell{Entryid: transaction.Id, Orderid: transaction.Orderid}
				ts.Orderstatus = string(updatedOrder.Status)

				cum, err := strconv.ParseFloat(updatedOrder.CummulativeQuoteQuantity, 64)
				if err != nil {
					return err
				}
				ts.Sellvalue = cum
				ts.Selltime = int(order.TransactTime)

				err = models.InsertTransactionSell(ts)
				count = 0
				for err != nil && count < 100 {
					err = models.InsertTransactionSell(ts)
					count += 1
				}
				if err != nil {
					return err
				}

				// res := (ts.Sellvalue - transaction.Buyvalue) / transaction.Buyvalue

				// log.Printf("TESTING %v Take_profit: \tMargin %v\tBuyvalue %v\tSellvalue %v\n", algo.Name, roundFloat(res, 5), transaction.Buyvalue, ts.Sellvalue)
			}

		}
		return nil
	case "waiting":
		return nil
	case "verification":
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
