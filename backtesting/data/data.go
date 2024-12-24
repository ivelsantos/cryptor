package data

import (
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
)

func GetData(algo models.Algor, window_size int) ([]binance.Kline, error) {
	var klines []binance.Kline

	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return klines, err
	}

	if window_size == 0 {
		return klines, fmt.Errorf("No window_size")
	}

	ticket := algo.BaseAsset + algo.QuoteAsset
	klines, err = functions.GetKlinesByStartEnd(ticket, account.ApiKey, account.SecretKey, window_size)
	if err != nil {
		return klines, err
	}

	return klines, nil
}
