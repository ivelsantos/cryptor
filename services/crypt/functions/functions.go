package functions

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
)

func GetKlines(symbol, apiKey, secretKey string) ([]binance.Kline, error) {
	var klineData []binance.Kline

	client := binance.NewClient(apiKey, secretKey)

	lines, err := client.NewKlinesService().Symbol(symbol).
		Interval("1m").Limit(1000).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return klineData, fmt.Errorf("Error fetching data: %v", err)
	}
	for _, k := range lines {
		klineData = append(klineData, *k)
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
