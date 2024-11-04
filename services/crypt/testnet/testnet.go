package testnet

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

func Buy(apikey, secretkey, symbol, orderquote string) (*binance.CreateOrderResponse, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	order, err := client.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).QuoteOrderQty(orderquote).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return order, nil
}

func Sell(apikey, secretkey, symbol, quantity string) (*binance.CreateOrderResponse, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	order, err := client.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).Quantity(quantity).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return order, nil
}

func SellQuote(apikey, secretkey, symbol, quoteOrder string) (*binance.CreateOrderResponse, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	order, err := client.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).QuoteOrderQty(quoteOrder).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return order, nil
}

func GetAccountQuote(apikey, secretkey, quote string) (string, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "", err
	}

	brl, err := getCoinBalance(quote, res.Balances)
	if err != nil {
		return "", err
	}

	return brl, nil
}

func getCoinBalance(coin string, balances []binance.Balance) (string, error) {
	if len(balances) < 1 {
		return "0", fmt.Errorf("Balances empty")
	}

	for _, balance := range balances {
		if balance.Asset == coin {
			return balance.Free, nil
		}
	}

	return "0", fmt.Errorf("Coin not found")
}

func GetMinNotional(apikey, secretkey, ticker string) (float64, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	res, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return 0, err
	}

	var symbol binance.Symbol

	for _, s := range res.Symbols {
		if s.Symbol == ticker {
			symbol = s
			break
		}
	}

	n, err := strconv.ParseFloat(symbol.NotionalFilter().MinNotional, 64)
	if err != nil {
		return n, err
	}

	return n, nil
}
