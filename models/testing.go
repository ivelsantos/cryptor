package models

import (
	"database/sql"
	"fmt"
)

type AlgoTesting struct {
	Id        int
	Botid     int
	Ticket    string
	Buyvalue  float64
	Buytime   int64
	Sellvalue float64
	Selltime  int64
}

type TestingBuy struct {
	Botid      int
	Baseasset  string
	Quoteasset string
	Buyvalue   float64
	Buytime    int64
}

type TestingSell struct {
	Entryid   int
	Sellvalue float64
	Selltime  int64
}

type AlgoStats struct {
	Botid             int
	TotalReturn       float64
	AvgReturnPerTrade float64
	AvgReturnPerMonth float64
	SucessRate        float64
	MaxDrawdown       float64
	AvgTradeTime      int
}

func InsertTestingBuy(tb TestingBuy) error {
	query := `
		INSERT INTO testing (botid, ticket, buyvalue, buytime)
		VALUES (?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tb.Botid, tb.Baseasset+tb.Quoteasset, tb.Buyvalue, tb.Buytime)
	for checkSqlBusy(err) {
		_, err = stmt.Exec(tb.Botid, tb.Baseasset+tb.Quoteasset, tb.Buyvalue, tb.Buytime)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func InsertTestingSell(ts TestingSell) error {
	query := `
		UPDATE testing
		SET sellvalue = ?,
			selltime = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ts.Sellvalue, ts.Selltime, ts.Entryid)
	for checkSqlBusy(err) {
		_, err = stmt.Exec(ts.Sellvalue, ts.Selltime, ts.Entryid)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetTestingSell(botid int) ([]AlgoTesting, error) {
	query := `
		SELECT id, botid, ticket, buyvalue, buytime FROM testing
		WHERE sellvalue IS NULL
		AND botid = ?
	`
	var algos []AlgoTesting

	rows, err := db.Query(query, botid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo AlgoTesting

		err := rows.Scan(&algo.Id, &algo.Botid, &algo.Ticket, &algo.Buyvalue, &algo.Buytime)
		if err != nil {
			return nil, err
		}

		algos = append(algos, algo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return algos, nil
}

func GetTestingBuy(botid int) ([]AlgoTesting, error) {
	query := `
		SELECT id, botid, ticket, buyvalue, buytime FROM testing
		WHERE sellvalue IS NULL
		AND botid = ?
	`
	var algos []AlgoTesting

	rows, err := db.Query(query, botid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo AlgoTesting

		err := rows.Scan(&algo.Id, &algo.Botid, &algo.Ticket, &algo.Buyvalue, &algo.Buytime)
		if err != nil {
			return nil, err
		}

		algos = append(algos, algo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return algos, nil
}

func GetUniqueAlgoTesting() ([]int, error) {
	query := `
	SELECT DISTINCT botid
	FROM testing
	`
	var botids []int

	rows, err := db.Query(query)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var botid int

		err := rows.Scan(&botid)
		if err != nil {
			return nil, err
		}

		botids = append(botids, botid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return botids, nil
}

func eraseTesting() error {
	query := `
		DELETE FROM testing
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	for checkSqlBusy(err) {
		_, err = stmt.Exec()
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func InsertTestingCalcTable() error {
	query := `
		INSERT INTO testing_calc (botid, return, buytime, selltime, buytimelength, selltimelength, tradetimelength)
SELECT *
  FROM (
           SELECT botid,
                  ( (sellvalue - (sellvalue * 0.001) ) - (buyvalue + (buyvalue * 0.001) ) ) / buyvalue AS return,
                  buytime,
                  selltime,
                  buytime - LAG(selltime) OVER (PARTITION BY botid ORDER BY id) AS buytimelength,
                  selltime - buytime AS selltimelength,
                  selltime - LAG(selltime) OVER (PARTITION BY botid ORDER BY id) AS tradetimelength
             FROM testing
            ORDER BY id
       )
  WHERE buytimelength IS NOT NULL AND 
       return IS NOT NULL;
	`
	_, err := db.Exec(query)
	for checkSqlBusy(err) {
		_, err = db.Exec(query)
	}
	if err != nil {
		return fmt.Errorf("Failed to insert into table: %v", err)
	}

	err = eraseTesting()

	return err
}

func GetAllAlgoStats() ([]AlgoStats, error) {
	query := `
	SELECT * FROM algo_stats
	`
	var algos []AlgoStats

	rows, err := db.Query(query)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo AlgoStats

		err := rows.Scan(&algo.Botid, &algo.TotalReturn, &algo.AvgReturnPerTrade, &algo.AvgReturnPerMonth, &algo.SucessRate, &algo.MaxDrawdown, &algo.AvgTradeTime)
		if err != nil {
			return nil, err
		}

		algos = append(algos, algo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return algos, nil
}

func GetStatsById(stats []AlgoStats, botid int) AlgoStats {
	for _, stat := range stats {
		if stat.Botid == botid {
			return stat
		}
	}

	return AlgoStats{AvgReturnPerMonth: 0}
}
