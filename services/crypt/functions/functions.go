package functions

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/models"
)

func GetKlines(symbol, apiKey, secretKey string, window int, lag int64) ([]binance.Kline, error) {
	var klineData []binance.Kline

	client := binance.NewClient(apiKey, secretKey)

	now := time.Now().UnixMilli() - (lag * 1000 * 60)
	before := now - (int64(window) * 1000 * 60)

	lines, err := client.NewKlinesService().Symbol(symbol).
		Interval("1m").StartTime(before).EndTime(now).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return klineData, fmt.Errorf("Error fetching data: %v", err)
	}

	for i := 0; i < len(lines); i++ {
		klineData = append(klineData, *lines[i])
	}

	// Getting the remaining data
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

func GetKlinesBacktesting(symbol, apiKey, secretKey string, window int, lag int64, index int) ([]binance.Kline, error) {
	var klineData []binance.Kline

	first := models.Backtesting_Data[0].CloseTime

	if len(models.Backtesting_Prov_Data) < (window + int(lag)) {
		var err error
		end := int64(first) - 60000
		models.Backtesting_Prov_Data, err = getKlinesSupport(symbol, apiKey, secretKey, window+int(lag), end)
		if err != nil {
			return klineData, fmt.Errorf("Error on getKlinesSupport: %v", err)
		}
	}

	newIndex := index + len(models.Backtesting_Prov_Data) - int(lag)

	// Tests just to be sure
	if (newIndex - window) < 0 {
		return klineData, fmt.Errorf("Error on GetKlinesBacktesting: (newIndex - window) < 0")
	}
	diff := models.Backtesting_Data[0].CloseTime - models.Backtesting_Prov_Data[len(models.Backtesting_Prov_Data)-1].CloseTime
	if diff != 60000 {
		return klineData, fmt.Errorf("Error on GetKlinesBacktesting: (Data - Prov_Data) != 60000 ")
	}

	klineData = append(klineData, joinIndex(&models.Backtesting_Prov_Data, &models.Backtesting_Data, uint(newIndex-window), uint(newIndex))...)

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

	// Getting the remaining data
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

func getKlinesSupport(symbol, apiKey, secretKey string, window_adjusted int, now int64) ([]binance.Kline, error) {
	var klineData []binance.Kline

	client := binance.NewClient(apiKey, secretKey)

	before := now - (int64(window_adjusted) * 1000 * 60)

	lines, err := client.NewKlinesService().Symbol(symbol).
		Interval("1m").StartTime(before).EndTime(now).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return klineData, fmt.Errorf("Error fetching data: %v", err)
	}

	for i := 0; i < len(lines); i++ {
		klineData = append(klineData, *lines[i])
	}

	// Getting the remaining data
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

func joinIndex(s1 *[]binance.Kline, s2 *[]binance.Kline, start, end uint) []binance.Kline {
	var result []binance.Kline

	s1_n := len(*s1)

	if int(end) < s1_n {
		return (*s1)[start:end]
	} else if int(start) < s1_n {
		result = append(result, (*s1)[start:]...)
		result = append(result, (*s2)[:int(end)-s1_n]...)
	} else {
		return (*s2)[int(start)-s1_n : int(end)-s1_n]
	}

	return result
}
