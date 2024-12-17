package crypt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
)

func GetMaxValue(algo models.Algor, args map[string]string) (float64, error) {
	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return 0, err
	}

	// Parsing window argument
	window, ok := args["window_size"]
	if !ok {
		return 0, fmt.Errorf("window argument not set")
	}
	window_int, err := strconv.ParseInt(window, 0, 0)
	if err != nil {
		return 0, err
	}

	// Parsing the lag argument
	var lag uint64 = 0
	num, ok := args["lag"]
	if ok {
		lag, err = strconv.ParseUint(num, 0, 0)
	}

	ticket := algo.BaseAsset + algo.QuoteAsset
	klines, err := functions.GetKlines(ticket, account.ApiKey, account.SecretKey, int(window_int), lag)
	if err != nil {
		return 0, err
	}

	var maxvalue float64 = 0
	for _, kline := range klines {
		v, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			return 0, err
		}
		if v > maxvalue {
			maxvalue = v
		}
	}

	return maxvalue, nil
}

func timeMsToSeconds(mili int64) time.Time {
	seconds := mili / 1000
	nanoseconds := (mili % 1000) * 1_000_000

	return time.Unix(seconds, nanoseconds)
}
