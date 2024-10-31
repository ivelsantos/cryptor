package crypt

import (
	"fmt"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
	"strconv"
)

func GetCloseMean(algo models.Algor, args []any) (float64, error) {
	if len(args) != 2 {
		return 0, fmt.Errorf("Usage: @Mean(window f64)")
	}

	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return 0, err
	}

	ticket := algo.BaseAsset + algo.QuoteAsset
	klines, err := functions.GetKlines(ticket, account.ApiKey, account.SecretKey)
	if err != nil {
		return 0, err
	}

	window := int(args[0].(float64))
	var val float64
	for i := len(klines); i > (len(klines) - window); i-- {
		v, err := strconv.ParseFloat(klines[i-1].Close, 64)
		if err != nil {
			return 0, err
		}
		val += v
	}

	mean := val / float64(window)
	return mean, nil
}
