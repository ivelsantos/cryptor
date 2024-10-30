package testnet

import (
	"log"
	"math"
	"strconv"
	"testing"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/models"
)

var coin string = "USDC"

func TestBuy(t *testing.T) {
	err := models.InitDB("../../../algor.db")
	if err != nil {
		t.Fatal(err)
	}

	algos, err := models.GetAllAlgos()
	if err != nil {
		t.Errorf("Failed to get algos: %v", err)
		return
	}

	account, err := models.GetAccountByName(algos[0].Owner)
	if err != nil {
		t.Error(err)
	}

	asset, err := GetAccountQuote(account.ApiKey_test, account.SecretKey_test, coin)
	if err != nil {
		t.Error(err)
	}

	asset_float, err := strconv.ParseFloat(asset, 64)
	if err != nil {
		t.Error(err)
	}

	quoteOrder := roundFloat(asset_float/5, 2)
	quoteOrderStr := strconv.FormatFloat(quoteOrder, 'f', -1, 64)

	order, err := Buy(account.ApiKey_test, account.SecretKey_test, "BTC"+coin, quoteOrderStr)
	if err != nil {
		log.Println("ERROR ON BUY")
		t.Error(err)
		return
	}

	quantity, _ := strconv.ParseFloat(order.ExecutedQuantity, 64)
	Cummulative, _ := strconv.ParseFloat(order.CummulativeQuoteQuantity, 64)
	price := Cummulative / quantity

	log.Printf("\nInitial account balance: %v\nQuantity: %v\nPrice: %v\nCummulative: %v\n\n", asset_float, order.ExecutedQuantity, price, Cummulative)

	sellorder, err := Sell(account.ApiKey_test, account.SecretKey_test, "BTC"+coin, order.ExecutedQuantity)
	if err != nil {
		log.Println("ERROR ON SELL")
		t.Error(err)
		return
	}

	log.Printf("\nSell cummulative: %v\nSell quantity: %v\n\n", sellorder.CummulativeQuoteQuantity, order.ExecutedQuantity)

}

func TestGetBalance(t *testing.T) {
	err := models.InitDB("../../../algor.db")
	if err != nil {
		t.Fatal(err)
	}

	algos, err := models.GetAllAlgos()
	if err != nil {
		t.Errorf("Failed to get algos: %v", err)
		return
	}

	account, err := models.GetAccountByName(algos[0].Owner)
	if err != nil {
		t.Error(err)
	}

	asset, err := GetAccountQuote(account.ApiKey_test, account.SecretKey_test, coin)
	if err != nil {
		t.Error(err)
	}

	log.Printf("\n%s Balance: %v\n\n", coin, asset)
}

func sumFills(fills []*binance.Fill) float64 {
	var sum float64

	for _, fill := range fills {
		commission, err := strconv.ParseFloat(fill.Commission, 64)
		if err != nil {
			return -1
		}
		sum += commission
	}

	return sum
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
