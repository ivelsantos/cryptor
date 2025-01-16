package operations

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/testingnet"
	"github.com/ivelsantos/cryptor/services/crypt/testnet"
)

func Buy(algo models.Algor, index any, args string) (bool, error) {
	// Parsing the arguments
	argsSlice := strings.Split(args, ",")
	argsSlice[len(argsSlice)-1] = strings.ReplaceAll(argsSlice[len(argsSlice)-1], ")", "")
	arguments := make(map[string]string)
	for _, arg := range argsSlice {
		sls := strings.Split(arg, "=")
		if len(sls) < 2 {
			continue
		}
		arguments[strings.Trim(sls[0], " ")] = strings.Trim(sls[1], " ")
	}
	// Up to this point, Buy() should have exactly one argument passed, quantity or percentage
	if len(arguments) != 1 {
		return false, fmt.Errorf("Buy() should have one of two arguments: quantity or percentage")
	}
	// Checking for arguments needed
	var arg string
	var valueStr string
	var ok bool
	if valueStr, ok = arguments["percentage"]; ok {
		arg = "percentage"
	} else if valueStr, ok = arguments["pct"]; ok {
		arg = "percentage"
	} else if valueStr, ok = arguments["quantity"]; !ok {
		if valueStr, ok = arguments["qty"]; !ok {
			return false, fmt.Errorf("Buy() should have one of two arguments: quantity or percentage")
		}
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return false, fmt.Errorf("Wrong format for %v value on Buy(): %v", arg, err)
	}

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
		transactions, err := models.GetTransactionBuy(algo.Id)
		if err != nil {
			return false, err
		}

		if len(transactions) != 0 {
			return false, nil
		}

		ticket := algo.BaseAsset + algo.QuoteAsset

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return false, err
		}

		// Getting quote balance available
		asset, err := testnet.GetAccountCoin(account.ApiKey_test, account.SecretKey_test, algo.QuoteAsset)
		if err != nil {
			return false, err
		}
		// Getting order value from argument
		var orderValue float64
		if arg == "percentage" {
			orderValue = asset * (value / 100)
		} else {
			orderValue = value
		}

		// Getting the minimal value allowed for an order
		minNotional, err := testnet.GetMinNotional(account.ApiKey_test, account.SecretKey_test, ticket)
		if err != nil {
			return false, err
		}

		// Verify if the order value is bigger enough
		if minNotional >= orderValue {
			// For now, if the order is not big enough, it will be ignored
			return false, nil
		}

		// Send the order
		orderValue = math.Floor(orderValue*100) / 100
		orderValueStr := strconv.FormatFloat(orderValue, 'f', -1, 64)
		order, err := testnet.Buy(account.ApiKey_test, account.SecretKey_test, ticket, orderValueStr)
		if err != nil {
			return false, err
		}

		//// Making sure the order is fulfilled
		time.Sleep(1 * time.Second) // wait for two seconds to make sure the transaction is done
		updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, ticket, int(order.OrderID))
		if err != nil {
			return false, fmt.Errorf("testnet.GetOrder: %v", err)
		}
		for string(updatedOrder.Status) != "FILLED" {
			time.Sleep(1 * time.Second)
			updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, ticket, int(order.OrderID))
			if err != nil {
				return false, fmt.Errorf("testnet.GetOrder: %v", err)
			}
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
		if err != nil {
			return false, fmt.Errorf("InsertTransactionBuy: %v", err)
		}
		return true, nil
	case "waiting":
		return false, nil
	case "verification":
		return false, nil
	case "backtesting":
		n := index.(int)

		ok := models.Backtesting_Transactions.CheckBought()
		if ok {
			return false, nil
		}

		err := models.Backtesting_Transactions.InsertBuy(models.Backtesting_Data[n])
		if err != nil {
			return false, err
		}
		return true, nil
	case "lang_test":
		return false, nil
	default:
		return false, fmt.Errorf("Unknown mode\n")
	}
}

func Sell(algo models.Algor, index any) error {
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

		bidPrice, _, err := testingnet.GetDepth(account.ApiKey, account.SecretKey, ticker)
		if err != nil {
			return err
		}

		for _, transaction := range transactions {
			ts := models.TestingSell{Entryid: transaction.Id, Sellvalue: bidPrice, Selltime: time.Now().Unix()}

			err = models.InsertTestingSell(ts)
			if err != nil {
				return fmt.Errorf("models.InsertTestingSell: %v", err)
			}
		}
		return nil
	case "live":
		transactions, err := models.GetTransactionSell(algo.Id)
		if err != nil {
			return fmt.Errorf("models.GetTesting: %v", err)
		}
		if len(transactions) != 1 {
			return nil
		}

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return fmt.Errorf("models.GetAccountByName: %v", err)
		}

		// Checking for available balance
		ok, err := checkBalance(account, algo.BaseAsset, transactions[0])
		if err != nil {
			return fmt.Errorf("models.EraseTransaction: %v", err)
		} else if !ok {
			return nil
		}

		quant := strconv.FormatFloat(transactions[0].Buyquantity, 'f', -1, 64)

		order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transactions[0].Ticket, quant)
		if err != nil {
			return fmt.Errorf("testnet.Sell: %v", err)
		}

		//// Making sure the order is fulfilled
		time.Sleep(1 * time.Second) // wait for two seconds to make sure the transaction is done

		updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transactions[0].Ticket, int(order.OrderID))
		if err != nil {
			return fmt.Errorf("testnet.GetOrder: %v", err)
		}

		for string(updatedOrder.Status) != "FILLED" {
			time.Sleep(1 * time.Second)
			updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transactions[0].Ticket, int(order.OrderID))
			if err != nil {
				return fmt.Errorf("testnet.GetOrder: %v", err)
			}
		}

		ts := models.TransactionSell{Entryid: transactions[0].Id, Orderid: transactions[0].Orderid}
		ts.Orderstatus = string(updatedOrder.Status)

		cum, err := strconv.ParseFloat(updatedOrder.CummulativeQuoteQuantity, 64)
		if err != nil {
			return fmt.Errorf("ParseFloat: %v", err)
		}
		ts.Sellvalue = cum
		ts.Selltime = int(order.TransactTime)

		err = models.InsertTransactionSell(ts)
		if err != nil {
			return fmt.Errorf("models.InsertTransactionSell: %v", err)
		}

		return nil
	case "verification", "waiting":
		return nil
	case "backtesting":
		n := index.(int)

		ok := models.Backtesting_Transactions.CheckSold()
		if ok {
			return nil
		}

		err := models.Backtesting_Transactions.InsertSell(models.Backtesting_Data[n])
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func StopLoss(algo models.Algor, stopPercentage float64, index any) error {
	// Converting percentage to proportion
	stop := stopPercentage / 100

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

		bidPrice, _, err := testingnet.GetDepth(account.ApiKey, account.SecretKey, ticker)
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
			}
		}
		return nil
	case "live":
		transactions, err := models.GetTransactionSell(algo.Id)
		if err != nil {
			return fmt.Errorf("models.GetTesting: %v", err)
		}
		if len(transactions) != 1 {
			return nil
		}
		transaction := transactions[0]

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return fmt.Errorf("models.GetAccountByName: %v", err)
		}

		// Checking for available balance
		ok, err := checkBalance(account, algo.BaseAsset, transaction)
		if err != nil {
			return fmt.Errorf("models.EraseTransaction: %v", err)
		} else if !ok {
			return nil
		}

		ticker := algo.BaseAsset + algo.QuoteAsset

		bidPrice, askPrice, err := testnet.GetDepth(account.ApiKey, account.SecretKey, ticker)
		if err != nil {
			return err
		}

		var price float64
		if askPrice > bidPrice {
			price = askPrice
		} else {
			price = bidPrice
		}

		buyprice := transaction.Buyvalue / transaction.Buyquantity
		sellPrice := buyprice - (stop * buyprice)

		if price <= sellPrice {
			quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
			order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
			if err != nil {
				return fmt.Errorf("testnet.Sell: %v", err)
			}

			//// Making sure the order is fulfilled
			time.Sleep(1 * time.Second) // wait for one seconds to make sure the transaction is done

			updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
			if err != nil {
				return fmt.Errorf("testnet.GetOrder: %v", err)
			}

			for string(updatedOrder.Status) != "FILLED" {
				time.Sleep(1 * time.Second)
				updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
				if err != nil {
					return fmt.Errorf("testnet.GetOrder: %v", err)
				}
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
			if err != nil {
				return fmt.Errorf("models.InsertTransactionSell: %v", err)
			}
		}
		return nil
	case "waiting":
		return nil
	case "verification":
		return nil
	case "backtesting":
		n := index.(int)

		ok := models.Backtesting_Transactions.CheckSold()
		if ok {
			return nil
		}

		err := models.Backtesting_Transactions.Stoploss(models.Backtesting_Data[n], stop)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}

}

func TakeProfit(algo models.Algor, takePercentage float64, index any) error {
	// Converting percentage to proportion
	take := takePercentage / 100

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

		bidPrice, _, err := testingnet.GetDepth(account.ApiKey, account.SecretKey, ticker)
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
			}
		}
		return nil
	case "live":
		transactions, err := models.GetTransactionSell(algo.Id)
		if err != nil {
			return err
		}
		if len(transactions) != 1 {
			return nil
		}

		transaction := transactions[0]

		account, err := models.GetAccountByName(algo.Owner)
		if err != nil {
			return err
		}

		// Checking for available balance
		ok, err := checkBalance(account, algo.BaseAsset, transaction)
		if err != nil {
			return fmt.Errorf("models.EraseTransaction: %v", err)
		} else if !ok {
			return nil
		}

		ticker := algo.BaseAsset + algo.QuoteAsset

		bidPrice, _, err := testnet.GetDepth(account.ApiKey, account.SecretKey, ticker)
		if err != nil {
			return err
		}

		buyprice := transaction.Buyvalue / transaction.Buyquantity
		sellPrice := buyprice + (take * buyprice)

		if bidPrice > sellPrice {
			quant := strconv.FormatFloat(transaction.Buyquantity, 'f', -1, 64)
			order, err := testnet.Sell(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, quant)
			if err != nil {
				return err
			}

			//// Making sure the order is fulfilled
			time.Sleep(1 * time.Second) // wait for two seconds to make sure the transaction is done

			updatedOrder, err := testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
			if err != nil {
				return fmt.Errorf("testnet.GetOrder: %v", err)
			}

			for string(updatedOrder.Status) != "FILLED" {
				time.Sleep(2 * time.Second)
				updatedOrder, err = testnet.GetOrder(account.ApiKey_test, account.SecretKey_test, transaction.Ticket, int(order.OrderID))
				if err != nil {
					return fmt.Errorf("testnet.GetOrder: %v", err)
				}
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
			if err != nil {
				return err
			}
		}
		return nil
	case "waiting":
		return nil
	case "verification":
		return nil
	case "backtesting":
		n := index.(int)

		ok := models.Backtesting_Transactions.CheckSold()
		if ok {
			return nil
		}

		err := models.Backtesting_Transactions.Takeprofit(models.Backtesting_Data[n], take)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("Unknown mode\n")
	}
}

func checkBalance(account models.Account, symbol string, transaction models.AlgoTransaction) (bool, error) {
	// Getting base balance available
	baseAsset, err := testnet.GetAccountCoin(account.ApiKey_test, account.SecretKey_test, symbol)
	if err != nil {
		return false, err
	}
	// If there is not enough base asset on the balance, the trade should be deleted
	if transaction.Buyquantity > baseAsset {
		err := models.EraseTransaction(transaction.Id)
		if err != nil {
			return false, fmt.Errorf("models.EraseTransaction: %v", err)
		}
		return false, nil
	}

	// If there is enough balance, return true
	return true, nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
