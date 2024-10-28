package testnet

import (
	"context"
	"fmt"
	"log"

	"github.com/adshao/go-binance/v2"
)

func Buy(apikey, secretkey, symbol, orderquote string) error {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)
	_ = client

	order, err := client.NewCreateOrderService().Symbol(symbol).Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).TimeInForce(binance.TimeInForceTypeGTC).QuoteOrderQty(orderquote).Do(context.Background())
	if err != nil {
		return err
	}

	log.Println(order)

	return nil
}

func GetAccountQuote(apikey, secretkey, quote string) (string, error) {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "", err
	}

	brl, err := getCoinBalance("BRL", res.Balances)
	if err != nil {
		return "", err
	}

	// log.Println(brl)

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