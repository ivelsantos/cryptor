package crypt

import (
	"fmt"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
	"github.com/ivelsantos/cryptor/services/crypt/values"

	"github.com/adshao/go-binance/v2"
)

type Crypt struct {
	vals   map[string]float64
	klines []binance.Kline
}

var cp *Crypt

func InitCrypt() error {
	cp = &Crypt{}

	ticket := "BTCBRL"

	val, err := values.GetPrice(ticket)
	if err != nil {
		return err
	}
	cp.vals = make(map[string]float64)
	cp.vals["@Price"] = val

	accounts, err := models.GetAccounts()
	if err != nil {
		return err
	}

	if len(accounts) > 0 {
		apiKey := accounts[0].ApiKey
		secretKey := accounts[0].SecretKey

		cp.klines, err = functions.GetKlines(ticket, apiKey, secretKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetCryptValue(key string) (float64, bool) {
	val, ok := cp.vals[key]
	return val, ok
}

func GetFuncValue(funcName string, args any) (float64, error) {
	argsSlice := args.([]any)
	switch funcName {

	case "@Mean":
		val, err := cp.GetCloseMean(argsSlice)
		if err != nil {
			return 0, err
		}
		return val, nil

	default:
		return 0, fmt.Errorf("Function %s does not exists", "@"+funcName)
	}
}
