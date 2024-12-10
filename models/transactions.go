package models

import (
	"database/sql"
	"fmt"
)

type AlgoTransaction struct {
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

type TransactionBuy struct {
	Botid       int
	Orderid     int
	Baseasset   string
	Quoteasset  string
	Orderstatus string
	Buyvalue    float64
	Buyquantity float64
	Buytime     int
}

type TransactionSell struct {
	Entryid     int
	Orderstatus string
	Sellvalue   float64
	Selltime    int
	Orderid     int
}

type AlgoStatsLive struct {
	Botid             int
	TotalReturn       float64
	AvgReturnPerTrade float64
	AvgReturnPerMonth float64
	SucessRate        float64
	MaxDrawdown       float64
	AvgTradeTime      int
}

func InsertTransactionBuy(tb TransactionBuy) error {
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
	for checkSqlBusy(err) {
		_, err = stmt.Exec(tb.Botid, tb.Orderid, tb.Baseasset+tb.Quoteasset, tb.Orderstatus, tb.Buyvalue, tb.Buyquantity, tb.Buytime)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func InsertTransactionSell(ts TransactionSell) error {
	query := `
		UPDATE testing
		SET sellvalue = ?,
			selltime = ?,
			orderstatus = ?,
			orderid = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ts.Sellvalue, ts.Selltime, ts.Orderstatus, ts.Orderid, ts.Entryid)
	for checkSqlBusy(err) {
		_, err = stmt.Exec(ts.Sellvalue, ts.Selltime, ts.Orderstatus, ts.Orderid, ts.Entryid)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func UpdateOrderStatus(status string, id int) error {
	query := `
		UPDATE testing
		SET orderstatus = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, id)
	for checkSqlBusy(err) {
		_, err = stmt.Exec(status, id)
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetTransactionSell(botid int) ([]AlgoTransaction, error) {
	query := `
		SELECT id, botid, orderid, ticket, orderstatus, buyvalue, buyquantity, buytime FROM testing
		WHERE sellvalue IS NULL
		AND orderstatus IS 'FILLED'
		AND botid = ?
	`
	var algos []AlgoTransaction

	rows, err := db.Query(query, botid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo AlgoTransaction

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

func GetTransactionBuy(botid int) ([]AlgoTransaction, error) {
	query := `
		SELECT id, botid, orderid, ticket, orderstatus, buyvalue, buyquantity, buytime FROM testing
		WHERE (sellvalue IS NULL OR orderstatus IS NOT 'FILLED')
		AND botid = ?
	`
	var algos []AlgoTransaction

	rows, err := db.Query(query, botid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo AlgoTransaction

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

func GetTransactionPending(botid int) ([]AlgoTransaction, error) {
	query := `
		SELECT id, botid, orderid, ticket, orderstatus, buyvalue, buyquantity, buytime FROM testing
		WHERE orderstatus IS NOT 'FILLED'
		AND botid = ?
	`
	var algos []AlgoTransaction

	rows, err := db.Query(query, botid)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve testing algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo AlgoTransaction

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

func GetUniqueAlgoTransaction() ([]int, error) {
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

func EraseTransaction() error {
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
	for checkSqlBusy(err) {
		_, err = stmt.Exec()
	}
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetAllAlgoStatsLive() ([]AlgoStats, error) {
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

		err := rows.Scan(&algo.Botid, &algo.TotalReturn, &algo.AvgReturnPerTrade, &algo.AvgReturnPerDay, &algo.SucessRate, &algo.MaxDrawdown, &algo.AvgTradeTime)
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

func GetStatsByIdLive(stats []AlgoStats, botid int) AlgoStats {
	for _, stat := range stats {
		if stat.Botid == botid {
			return stat
		}
	}

	return AlgoStats{AvgReturnPerDay: 0}
}
