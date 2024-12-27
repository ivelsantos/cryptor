package crypt

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/models"
	"github.com/ivelsantos/cryptor/services/crypt/functions"
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

	return calculateMean(values), nil
}

func GetMedianValue(algo models.Algor, args map[string]string) (float64, error) {
	values, err := gettingKlines(algo, args)
	if err != nil {
		return 0, err
	}
	if len(values) < 1 {
		return 0, nil
	}

	slices.Sort(values)

	n := len(values)
	var result float64
	if n%2 == 1 {
		result = values[int(n/2)]
	} else {
		i := n / 2
		result = (values[i-1] + values[i]) / 2
	}

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

	var sum float64
	for _, value := range values {
		sum += value
	}
	meanValue := sum / float64(len(values))

	var varianceSum float64
	for _, v := range values {
		diff := v - meanValue
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(values))

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

	var sum float64
	for _, value := range values {
		sum += value
	}
	meanValue := sum / float64(len(values))

	var varianceSum float64
	for _, v := range values {
		diff := v - meanValue
		varianceSum += diff * diff
	}
	variance := varianceSum / float64(len(values))

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

func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	if len(values) == 1 {
		return values[0]
	}

	var sum float64
	for _, value := range values {
		sum += value
	}
	meanValue := sum / float64(len(values))

	return meanValue
}

func timeMsToSeconds(mili int64) time.Time {
	seconds := mili / 1000
	nanoseconds := (mili % 1000) * 1_000_000

	return time.Unix(seconds, nanoseconds)
}
