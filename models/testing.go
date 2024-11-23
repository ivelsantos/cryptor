package models

import (
	"database/sql"
	"fmt"
)

type AlgoTesting struct {
	Id          int
	Botid       int
	Orderid     int
	Ticket      string
	Orderstatus string
	Buyvalue    float64
	Buyquantity float64
	Buytime     int
	Sellvalue   float64
	Selltime    int
}

type TestingBuy struct {
	Botid       int
	Orderid     int
	Baseasset   string
	Quoteasset  string
	Orderstatus string
	Buyvalue    float64
	Buyquantity float64
	Buytime     int
}

type TestingSell struct {
	Entryid     int
	Orderstatus string
	Sellvalue   float64
	Selltime    int
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
		INSERT INTO testing (botid, orderid, ticket, orderstatus, buyvalue, buyquantity, buytime)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tb.Botid, tb.Orderid, tb.Baseasset+tb.Quoteasset, tb.Orderstatus, tb.Buyvalue, tb.Buyquantity, tb.Buytime)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func InsertTestingSell(ts TestingSell) error {
	query := `
		UPDATE testing
		SET sellvalue = ?,
			selltime = ?,
			orderstatus = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ts.Sellvalue, ts.Selltime, ts.Orderstatus, ts.Entryid)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetTesting(botid int) ([]AlgoTesting, error) {
	query := `
		SELECT id, botid, orderid, ticket, orderstatus, buyvalue, buyquantity, buytime FROM testing
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

		err := rows.Scan(&algo.Id, &algo.Botid, &algo.Orderid, &algo.Ticket, &algo.Orderstatus, &algo.Buyvalue, &algo.Buyquantity, &algo.Buytime)
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
		SELECT id, botid, orderid, ticket, orderstatus, buyvalue, buyquantity, buytime FROM testing
		WHERE sellvalue IS NULL
		OR orderstatus IS NOT 'FILLED'
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

		err := rows.Scan(&algo.Id, &algo.Botid, &algo.Orderid, &algo.Ticket, &algo.Orderstatus, &algo.Buyvalue, &algo.Buyquantity, &algo.Buytime)
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

func EraseTesting() error {
	query := `
		DELETE FROM testing
		WHERE sellvalue IS NULL
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
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
