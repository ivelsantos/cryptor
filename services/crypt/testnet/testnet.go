package testnet

import (
	"context"
	"log"

	"github.com/adshao/go-binance/v2"
)

func Buy(apikey, secretkey string) error {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)
	_ = client

	return nil
}

func GetAccount(apikey, secretkey string) error {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return err
	}

	log.Println(res.Balances)

	return nil
}
