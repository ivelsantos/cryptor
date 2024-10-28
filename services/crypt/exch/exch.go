package exch

import (
	"github.com/adshao/go-binance/v2"
)

func BuyTest(apikey, secretkey string) error {
	binance.UseTestnet = true
	client := binance.NewClient(apikey, secretkey)
	_ = client

	return nil
}
