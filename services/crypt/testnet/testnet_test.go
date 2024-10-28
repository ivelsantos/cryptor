package testnet

import (
	"github.com/ivelsantos/cryptor/models"
	"testing"
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

	err = GetAccount(account.ApiKey_test, account.SecretKey_test)
	if err != nil {
		t.Error(err)
	}
}
