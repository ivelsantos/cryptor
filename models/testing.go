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

// func InsertTestingBuy(botid int, orderid int, ticket string, orderstatus string, buyprice float64, buyquantity string, buytime int) error {
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

// func InsertTestingSell(entryid int, sellprice float64, sellquantity string, selltime int, orderstatus string) error {
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
