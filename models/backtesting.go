package models

import (
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

type AlgoBacktesting struct {
	Id              []int
	Buyvalue        []float64
	Buytime         []int64 // Seconds
	Sellvalue       []float64
	Selltime        []int64 // Seconds
	Return          []float64
	Buytimelength   []int64 // Seconds
	Selltimelength  []int64 // Seconds
	Tradetimelength []int64 // Seconds
}

type ResultBacktesting struct {
	Total_return float64
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

	return nil
}

func (a *AlgoBacktesting) InsertSell(line binance.Kline) error {
	value, err := strconv.ParseFloat(line.Close, 64)
	if err != nil {
		return fmt.Errorf("Error on backtesting insertSell: %v", err)
	}

	a.Sellvalue = append(a.Sellvalue, value)
	a.Selltime = append(a.Selltime, line.CloseTime/1000)

	// Calculating return considering Binance fees
	buyvalue := a.Buyvalue[len(a.Buyvalue)-1]
	a.Return = append(a.Return, ((value-(value*0.001))-(buyvalue+(buyvalue*0.001)))/buyvalue)

	// Calculating Buytimelength
	if len(a.Sellvalue) > 1 {
		selltimeBefore := a.Selltime[len(a.Selltime)-2]
		buytimelength := a.Buytime[len(a.Buytime)-1] - selltimeBefore
		a.Buytimelength = append(a.Buytimelength, buytimelength)
	} else {
		a.Buytimelength = append(a.Buytimelength, 0)
	}

	// Calculating Selltimelength: Selltime - Buytime
	a.Selltimelength = append(a.Selltimelength, a.Selltime[len(a.Selltime)-1]-a.Buytime[len(a.Buytime)-1])

	// Calculating Tradetimelength: Selltimelength + Buytimelength
	a.Tradetimelength = append(a.Tradetimelength, a.Selltimelength[len(a.Selltimelength)-1]+a.Buytimelength[len(a.Buytimelength)-1])

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
		return a.Id[len(a.Id)-1]
	case "buyvalue":
		return a.Buyvalue[len(a.Buyvalue)-1]
	case "buytime":
		return a.Buytime[len(a.Buytime)-1]
	case "selltime":
		return a.Selltime[len(a.Selltime)-1]
	case "sellvalue":
		return a.Sellvalue[len(a.Sellvalue)-1]
	default:
		return struct{}{}
	}
}
