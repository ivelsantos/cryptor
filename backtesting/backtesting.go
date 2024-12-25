package backtesting

import (
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/backtesting/data"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
)

func BackTesting(algo models.Algor, window_size int) error {
	var err error

	models.Backtesting_Data, err = data.GetData(algo, window_size)
	if err != nil {
		return err
	}

	algo.State = "backtesting"

	for i := range models.Backtesting_Data {
		optAlgo := lang.GlobalStore("Algo", algo)
		optIndex := lang.GlobalStore("Back", i)
		optBackData := lang.GlobalStore("BackData", models.Backtesting_Data)
		optBackTransaction := lang.GlobalStore("BackTransactions", models.Backtesting_Transactions)

		_, err = lang.Parse("", []byte(algo.Buycode), optAlgo, optIndex, optBackData, optBackTransaction)
		if err != nil {
			log.Printf("%v: Parsing error: %v\n", algo.Name, err)
		}
	}

	models.Backtesting_Data = []binance.Kline{}
	models.Backtesting_Transactions = models.AlgoBacktesting{}

	return nil
}
