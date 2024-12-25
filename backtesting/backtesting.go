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

var Backtesting_Data []binance.Kline
var Backtesting_Transactions AlgoBacktesting

func BackTesting(algo models.Algor, window_size int) error {
	var err error

	Backtesting_Data, err = data.GetData(algo, window_size)
	if err != nil {
		return err
	}

	log.Printf("\tBacktesting with %v datapoints\n", len(Backtesting_Data))

	for _, line := range Backtesting_Data {
		optAlgo := lang.GlobalStore("Algo", algo)
		optBack := lang.GlobalStore("Back", line.CloseTime)

		_, err = lang.Parse("", []byte(algo.Buycode), optAlgo, optBack)
		if err != nil {
			log.Printf("%v: Parsing error: %v\n", algo.Name, err)
		}
	}

	return nil
}
