package testnet

import (
	"log"
	"strconv"
	"testing"

	"github.com/ivelsantos/cryptor/models"
)

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

	quote := "USDT"

	brl, err := GetAccountQuote(account.ApiKey_test, account.SecretKey_test, quote)
	if err != nil {
		t.Error(err)
	}

	brl_float, err := strconv.ParseFloat(brl, 64)
	if err != nil {
		t.Error(err)
	}

	quoteOrder := brl_float / 5

	quoteOrderStr := strconv.FormatFloat(quoteOrder, 'f', -1, 64)

	order, err := Buy(account.ApiKey_test, account.SecretKey_test, "BTC"+quote, quoteOrderStr)
	if err != nil {
		log.Println("ERROR ON BUY")
		t.Error(err)
		return
	}

	quantity, _ := strconv.ParseFloat(order.ExecutedQuantity, 64)
	price := quoteOrder / quantity

	log.Printf("\nInitial account balance: %v\nQuote Order: %v\nQuantity: %v\nPrice: %v\n", brl_float, quoteOrderStr, order.ExecutedQuantity, price)

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

	BRL, err := GetAccountQuote(account.ApiKey_test, account.SecretKey_test, "USDT")
	if err != nil {
		t.Error(err)
	}

	log.Printf("\nUSDT Balance: %v\n\n", BRL)
}
