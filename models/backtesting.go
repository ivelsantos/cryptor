package models

import (
	"database/sql"
	"fmt"
	"math"
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
	Botid               int
	Daily_return        float64
	Ticket_Daily_return float64
	Sucess_rate         float64
	Avg_trade_time      int64
	Window              int
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

func (a *AlgoBacktesting) Metrics(days int, closeStart, closeEnd string) ResultBacktesting {
	result := ResultBacktesting{}
	n_trades := len(a.Return)
	if n_trades < 1 {
		return result
	}

	priceStart, _ := strconv.ParseFloat(closeStart, 64)
	priceEnd, _ := strconv.ParseFloat(closeEnd, 64)

	// Daily return
	tw := 1.0
	for _, value := range a.Return {
		tw = tw * (1 + value)
	}
	twr := tw - 1
	result.Daily_return = math.Pow(1+twr, 1/float64(days)) - 1

	// Ticket daily return
	ticketReturn := (priceEnd - priceStart) / priceStart
	result.Ticket_Daily_return = math.Pow(1+ticketReturn, 1/float64(days)) - 1

	// Sucess rate
	var sucessTrades float64
	for _, value := range a.Return {
		if value > 0 {
			sucessTrades += 1
		}
	}
	result.Sucess_rate = sucessTrades / float64(n_trades)

	// Average trade time
	var totalTradeTime int64
	for _, value := range a.Tradetimelength {
		totalTradeTime += value
	}

	result.Avg_trade_time = totalTradeTime / int64(n_trades)

	return result
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

func InsertBacktesting(algo Algor) error {
	query := `
		INSERT INTO backtesting (botid)
		VALUES (?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(algo.Id)
	for checkSqlBusy(err) {
		_, err = stmt.Exec(algo.Id)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func InsertMetricsBacktesting(mt ResultBacktesting) error {
	query := `
		UPDATE backtesting
		SET daily_return = ?,
			ticket_daily_return = ?
			sucess_rate = ?
			avg_trade_time = ?
			window = ?
		WHERE botid = ?
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(mt.Daily_return, mt.Ticket_Daily_return, mt.Sucess_rate, mt.Avg_trade_time, mt.Window, mt.Botid)
	for checkSqlBusy(err) {
		_, err = stmt.Exec(mt.Daily_return, mt.Ticket_Daily_return, mt.Sucess_rate, mt.Avg_trade_time, mt.Window, mt.Botid)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetBacktestingById(botid int) (ResultBacktesting, error) {
	query := `
	SELECT * FROM backtesting
	WHERE botid = ?
	`
	var mt ResultBacktesting

	row := db.QueryRow(query, botid)

	err := row.Scan(&mt.Botid, &mt.Daily_return, &mt.Ticket_Daily_return, &mt.Sucess_rate, &mt.Avg_trade_time, &mt.Window)
	if err != nil {
		if err == sql.ErrNoRows {
			// Because sqlite3 rowid autoincrement starts in 1 we can do that
			// Therefore, if mt.Botid == 0 there is not backtesting data on db for a given botid
			mt.Botid = 0
			return mt, nil
		}
		return mt, fmt.Errorf("Failed to retrieve backtesting: %v", err)
	}

	return mt, nil
}
func eraseBacktestingByBotid(botid int) error {
	query := `
		DELETE FROM backtesting
		WHERE botid = ?`

	_, err := db.Exec(query, botid)
	if err != nil {
		return fmt.Errorf("Failed to delete backtesting: %v", err)
	}

	return nil
}
