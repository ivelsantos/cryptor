package crypt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
)

func GetMaxValue(algo models.Algor, args map[string]string) (float64, error) {
	klines, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(klines) < 1 {
		return 0, nil
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

func GetMinValue(algo models.Algor, args map[string]string) (float64, error) {
	klines, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(klines) < 1 {
		return 0, nil
	}

	var minvalue float64
	minvalue, err = strconv.ParseFloat(klines[0].Close, 64)
	if err != nil {
		return 0, err
	}
	for _, kline := range klines {
		v, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			return 0, err
		}
		if v < minvalue {
			minvalue = v
		}
	}

	return minvalue, nil
}

func GetMeanValue(algo models.Algor, args map[string]string) (float64, error) {
	klines, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(klines) < 1 {
		return 0, nil
	}

	var sum float64 = 0
	for _, kline := range klines {
		v, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			return 0, err
		}
		sum += v
	}

	meanValue := sum / float64(len(klines))
	return meanValue, nil
}

func gettingKlines(algo models.Algor, args map[string]string) ([]binance.Kline, error) {
	var klines []binance.Kline

	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return klines, err
	}

	// Parsing window argument
	window, ok := args["window_size"]
	if !ok {
		return klines, fmt.Errorf("window argument not set")
	}
	window_int, err := strconv.ParseInt(window, 0, 0)
	if err != nil {
		return klines, err
	}
	if window_int == 0 {
		return klines, nil
	}

	// Parsing the lag argument
	var lag uint64 = 0
	num, ok := args["lag"]
	if ok {
		lag, err = strconv.ParseUint(num, 0, 0)
	}

	ticket := algo.BaseAsset + algo.QuoteAsset
	klines, err = functions.GetKlines(ticket, account.ApiKey, account.SecretKey, int(window_int), lag)
	if err != nil {
		return klines, err
	}

	return klines, nil
}

func timeMsToSeconds(mili int64) time.Time {
	seconds := mili / 1000
	nanoseconds := (mili % 1000) * 1_000_000

	return time.Unix(seconds, nanoseconds)
}
