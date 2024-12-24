package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
)

func GetKlines(symbol, apiKey, secretKey string, window int, lag int64) ([]binance.Kline, error) {
	var klineData []binance.Kline

	client := binance.NewClient(apiKey, secretKey)

	now := time.Now().UnixMilli() - (lag * 1000)
	before := now - (int64(window) * 1000)

	lines, err := client.NewKlinesService().Symbol(symbol).
		Interval("1m").StartTime(before).EndTime(now).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return klineData, fmt.Errorf("Error fetching data: %v", err)
	}

	for i := 0; i < len(lines)-int(lag); i++ {
		klineData = append(klineData, *lines[i])
	}

	last := klineData[len(klineData)-1].CloseTime

	for last < now && len(lines) > 0 {
		lines, err = client.NewKlinesService().Symbol(symbol).
			Interval("1m").StartTime(last + 1).EndTime(now).Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return klineData, fmt.Errorf("Error fetching data: %v", err)
		}

		for i := 0; i < len(lines); i++ {
			klineData = append(klineData, *lines[i])
		}
		last = klineData[len(klineData)-1].CloseTime
	}

	return klineData, nil
}

func GetKlinesByStartEnd(symbol, apiKey, secretKey string, window_size int) ([]binance.Kline, error) {
	var klineData []binance.Kline

	client := binance.NewClient(apiKey, secretKey)

	now, before, err := getWindowTimes(window_size)
	if err != nil {
		return klineData, fmt.Errorf("Error in getWindowTimes")
	}

	lines, err := client.NewKlinesService().Symbol(symbol).
		Interval("1m").StartTime(before).EndTime(now).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return klineData, fmt.Errorf("Error fetching data: %v", err)
	}

	for i := 0; i < len(lines); i++ {
		klineData = append(klineData, *lines[i])
	}

	last := klineData[len(klineData)-1].CloseTime

	for last < now && len(lines) > 0 {
		lines, err = client.NewKlinesService().Symbol(symbol).
			Interval("1m").StartTime(last + 1).EndTime(now).Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return klineData, fmt.Errorf("Error fetching data: %v", err)
		}

		for i := 0; i < len(lines); i++ {
			klineData = append(klineData, *lines[i])
		}
		last = klineData[len(klineData)-1].CloseTime
	}

	return klineData, nil
}

func GetSymbols(apiKey, secretKey string) ([]binance.Symbol, error) {
	var symbols []binance.Symbol
	client := binance.NewClient(apiKey, secretKey)

	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return symbols, err
	}

	symbols = exchangeInfo.Symbols

	return symbols, nil
}

func getWindowTimes(window int) (int64, int64, error) {
	now := time.Now().UnixMilli()
	before := time.Now().AddDate(0, 0, (window * -1)).UnixMilli()

	return now, before, nil
}
