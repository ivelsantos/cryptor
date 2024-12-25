package backtesting

import (
	"log"

	"github.com/adshao/go-binance/v2"
	"github.com/ivelsantos/cryptor/backtesting/data"
	"github.com/ivelsantos/cryptor/lang"
	"github.com/ivelsantos/cryptor/models"
)

type AlgoBacktesting struct {
	Id        []int
	Botid     []int
	Ticket    []string
	Buyvalue  []float64
	Buytime   []int64
	Sellvalue []float64
	Selltime  []int64
}

var Backtesting_Transactions AlgoBacktesting
var Backtesting_Data []binance.Kline

func BackTesting(algo models.Algor, window_size int) error {
	var err error

	Backtesting_Data, err = data.GetData(algo, window_size)
	if err != nil {
		return err
	}

	for _, line := range Backtesting_Data {
		optAlgo := lang.GlobalStore("Algo", algo)
		optBack := lang.GlobalStore("Back", line.CloseTime)
		optBackData := lang.GlobalStore("BackData", Backtesting_Data)

		_, err = lang.Parse("", []byte(algo.Buycode), optAlgo, optBack, optBackData)
		if err != nil {
			log.Printf("%v: Parsing error: %v\n", algo.Name, err)
		}
	}

	Backtesting_Data = []binance.Kline{}

	return nil
}
