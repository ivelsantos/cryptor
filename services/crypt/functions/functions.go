package functions

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
)

func GetKlines(symbol, apiKey, secretKey string, window int, lag uint64) ([]binance.Kline, error) {
	var klineData []binance.Kline

	client := binance.NewClient(apiKey, secretKey)

	lines, err := client.NewKlinesService().Symbol(symbol).
		Interval("1m").Limit(window + int(lag)).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return klineData, fmt.Errorf("Error fetching data: %v", err)
	}

	for i := 0; i < len(lines)-int(lag); i++ {
		klineData = append(klineData, *lines[i])
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
