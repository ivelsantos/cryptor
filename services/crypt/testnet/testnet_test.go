package testnet

import (
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

	brl, err := GetAccountQuote(account.ApiKey_test, account.SecretKey_test, "BRL")
	if err != nil {
		t.Error(err)
	}

	brl_float, err := strconv.ParseFloat(brl, 64)
	if err != nil {
		t.Error(err)
	}

	quoteOrder := brl_float / 10

	quoteOrderStr := strconv.FormatFloat(quoteOrder, 'f', -1, 64)

	err = Buy(account.ApiKey_test, account.SecretKey_test, "BTCBRL", quoteOrderStr)
	if err != nil {
		t.Error(err)
		return
	}

}
