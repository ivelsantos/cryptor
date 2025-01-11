package crypt

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
	"gonum.org/v1/gonum/stat"
)

func GetMaxValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	return slices.Max(values), nil
}

func GetMinValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	return slices.Min(values), nil
}

func GetRangeValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	min := slices.Min(values)
	max := slices.Max(values)

	return max - min, nil
}

func GetMeanValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	return stat.Mean(values, nil), nil
}

func GetMedianValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	result := stat.Quantile(0.5, stat.Empirical, values, nil)

	return result, nil
}

func GetVarValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	variance := stat.Variance(values, nil)

	return variance, nil
}

func GetStdValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	variance := stat.Variance(values, nil)

	stdDev := math.Sqrt(variance)

	return stdDev, nil
}

func GetEmaValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	// Parsing smoothing argument
	var smoothing float64 = 2
	num, ok := args["smoothing"]
	if ok {
		smoothing, err = strconv.ParseFloat(num, 64)
		return 0, fmt.Errorf("Invalid Smoothing Paremeter: %v", err)
	}
	alpha := smoothing / float64(len(values)+1)

	// Initial smoothed value set to the first value
	ema := values[0]

	for i := range values {
		if i == 0 {
			continue
		}
		ema = (values[i] * alpha) + (ema * (1 - alpha))
	}

	return ema, nil
}

func GetPrice(algo models.Algor, args map[string]string) (float64, error) {
	if algo.State == "backtesting" {
		num, _ := args["backindex"]
		n, err := strconv.Atoi(num)
		if err != nil {
			return 0, fmt.Errorf("Error on GetPrice: %v", err)
		}

		value, err := strconv.ParseFloat(models.Backtesting_Data[n].Close, 64)
		if err != nil {
			return 0.0, fmt.Errorf("Error parsing close value: %v", err)
		}

		return value, nil
	}

	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return 0, err
	}

	client := binance.NewClient(account.ApiKey, account.SecretKey)

	price, err := client.NewListPricesService().Symbol(algo.BaseAsset + algo.QuoteAsset).Do(context.Background())
	if err != nil {
		return 0, fmt.Errorf("Error on GetPrice: %v", err)
	}

	if len(price) != 1 {
		return 0, fmt.Errorf("Wrong number of ticker prices")
	}

	priceFloat, err := strconv.ParseFloat(price[0].Price, 64)
	if err != nil {
		return 0, fmt.Errorf("Error on GetPrice: %v", err)
	}

	return priceFloat, nil
}

func gettingKlines(algo models.Algor, args map[string]string) ([]float64, error) {
	var klines []binance.Kline
	var values []float64

	account, err := models.GetAccountByName(algo.Owner)
	if err != nil {
		return values, err
	}

	// Parsing window argument
	window, ok := args["window_size"]
	if !ok {
		return values, fmt.Errorf("window argument not set")
	}
	window_int, err := strconv.ParseInt(window, 0, 0)
	if err != nil {
		return values, err
	}
	if window_int == 0 {
		return values, nil
	}

	// Parsing the lag argument
	var lag int64 = 0
	num, ok := args["lag"]
	if ok {
		lag, err = strconv.ParseInt(num, 0, 0)
		if err != nil {
			return values, fmt.Errorf("Error on ParseInt on lag argument parsing: %v", err)
		}
	}

	ticket := algo.BaseAsset + algo.QuoteAsset
	if algo.State == "backtesting" {
		num, _ := args["backindex"]
		n, err := strconv.Atoi(num)
		if err != nil {
			return values, err
		}

		klines, err = functions.GetKlinesBacktesting(ticket, account.ApiKey, account.SecretKey, int(window_int), lag, n)
		if err != nil {
			return values, err
		}
	} else {
		klines, err = functions.GetKlines(ticket, account.ApiKey, account.SecretKey, int(window_int), lag)
		if err != nil {
			return values, err
		}
	}

	for _, kline := range klines {
		v, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			return values, err
		}
		values = append(values, v)
	}

	return values, nil
}

func timeMsToSeconds(mili int64) time.Time {
	seconds := mili / 1000
	nanoseconds := (mili % 1000) * 1_000_000

	return time.Unix(seconds, nanoseconds)
}
