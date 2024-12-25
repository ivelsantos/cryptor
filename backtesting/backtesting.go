package backtesting

import (
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/backtesting/data"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
)

var Backtesting_Transactions models.AlgoBacktesting
var Backtesting_Data []binance.Kline

func BackTesting(algo models.Algor, window_size int) error {
	var err error

	Backtesting_Data, err = data.GetData(algo, window_size)
	if err != nil {
		return err
	}

	algo.State = "backtesting"

	for i := range Backtesting_Data {
		optAlgo := lang.GlobalStore("Algo", algo)
		optIndex := lang.GlobalStore("Back", i)
		optBackData := lang.GlobalStore("BackData", Backtesting_Data)
		optBackTransaction := lang.GlobalStore("BackTransactions", Backtesting_Transactions)

		_, err = lang.Parse("", []byte(algo.Buycode), optAlgo, optIndex, optBackData, optBackTransaction)
		if err != nil {
			log.Printf("%v: Parsing error: %v\n", algo.Name, err)
		}
	}

	Backtesting_Data = []binance.Kline{}
	Backtesting_Transactions = models.AlgoBacktesting{}

	return nil
}
