package models

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

type AlgoBacktesting struct {
	Id        []int
	Buyvalue  []float64
	Buytime   []int64
	Sellvalue []float64
	Selltime  []int64
}

var Backtesting_Transactions AlgoBacktesting
var Backtesting_Data []binance.Kline
var Backtesting_Prov_Data []binance.Kline

func (a *AlgoBacktesting) CheckBought() bool {
	if len(a.Buytime) != len(a.Selltime) {
		return true
	}
	return false
}

func (a *AlgoBacktesting) CheckSold() bool {
	if len(a.Buytime) == len(a.Selltime) {
		return true
	}
	return false
}

func (a *AlgoBacktesting) InsertBuy(line binance.Kline) error {
	value, err := strconv.ParseFloat(line.Close, 64)
	if err != nil {
		return fmt.Errorf("Error on backtesting insertBuy: %v", err)
	}

	a.Id = append(a.Id, len(a.Id))
	a.Buyvalue = append(a.Buyvalue, value)
	a.Buytime = append(a.Buytime, line.CloseTime/1000)
	a.Sellvalue = append(a.Sellvalue, 0)
	a.Selltime = append(a.Selltime, 0)

	return nil
}

func (a *AlgoBacktesting) InsertSell(line binance.Kline) error {
	value, err := strconv.ParseFloat(line.Close, 64)
	if err != nil {
		return fmt.Errorf("Error on backtesting insertSell: %v", err)
	}

	a.Sellvalue = append(a.Sellvalue, value)
	a.Selltime = append(a.Selltime, line.CloseTime/1000)

	return nil
}

func (a *AlgoBacktesting) Stoploss(line binance.Kline, stop float64) error {
	value, err := strconv.ParseFloat(line.Close, 64)
	if err != nil {
		return fmt.Errorf("Error on backtesting insertSell: %v", err)
	}

	buyvalue := a.last("buyvalue").(float64)

	threshold := buyvalue - (stop * buyvalue)
	if value <= threshold {
		a.InsertSell(line)
	}

	return nil
}

func (a *AlgoBacktesting) Takeprofit(line binance.Kline, take float64) error {
	value, err := strconv.ParseFloat(line.Close, 64)
	if err != nil {
		return fmt.Errorf("Error on backtesting insertSell: %v", err)
	}

	buyvalue := a.last("buyvalue").(float64)

	threshold := buyvalue + (take * buyvalue)
	if value > threshold {
		a.InsertSell(line)
	}

	return nil
}

func (a *AlgoBacktesting) last(field string) any {
	switch field {
	case "id":
		return a.Id
	case "buyvalue":
		return a.Buyvalue
	case "buytime":
		return a.Buytime
	case "selltime":
		return a.Selltime
	case "sellvalue":
		return a.Sellvalue
	default:
		return struct{}{}
	}
}
