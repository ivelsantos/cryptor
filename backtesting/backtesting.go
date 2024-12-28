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

		_, err = lang.Parse("", []byte(algo.Buycode), optAlgo, optIndex)
		if err != nil {
			log.Fatalf("%v: Parsing error: %v\n", algo.Name, err)
		}

		//// LOGGING
		index := i + 1
		if (index % 1440) == 0 {
			log.Printf("\tDay %v\n", index/1440)
		}
	}

	days := len(models.Backtesting_Data) / 1440
	priceStart := models.Backtesting_Data[0].Close
	priceEnd := models.Backtesting_Data[len(models.Backtesting_Data)-1].Close
	metrics := models.Backtesting_Transactions.Metrics(days, priceStart, priceEnd)

	log.Printf("\n")
	log.Printf("\tNumber of trades: %v\n", len(models.Backtesting_Transactions.Id))
	log.Printf("\tAverage trade time: %v\n", metrics.Avg_trade_time/60)
	log.Printf("\tDaily return: %.4f\n", metrics.Daily_return)
	log.Printf("\tTicket daily return: %.4f\n", metrics.Ticket_Daily_return)
	log.Printf("\tSucess rate: %.4f\n", metrics.Sucess_rate)

	models.Backtesting_Data = []binance.Kline{}
	models.Backtesting_Prov_Data = []binance.Kline{}
	models.Backtesting_Transactions = models.AlgoBacktesting{}

	return nil
}
