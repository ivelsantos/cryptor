package models

import (
	"database/sql"
	"fmt"
)

type AlgoTesting struct {
	Id           int
	Botid        int
	Orderid      int
	Ticket       string
	Orderstatus  string
	Buyprice     float64
	Buyquantity  float64
	Buytime      int
	Sellprice    float64
	Sellquantity float64
	Selltime     int
}

type TestingBuy struct {
	Botid       int
	Orderid     int
	Baseasset   string
	Quoteasset  string
	Orderstatus string
	Buyprice    float64
	Buyquantity float64
	Buytime     int
}

type TestingSell struct {
	Entryid      int
	Orderstatus  string
	Sellprice    float64
	Sellquantity float64
	Selltime     int
}

// func InsertTestingBuy(botid int, orderid int, ticket string, orderstatus string, buyprice float64, buyquantity string, buytime int) error {
func InsertTestingBuy(tb TestingBuy) error {
	query := `
		INSERT INTO testing (botid, orderid, ticket, orderstatus, buyprice, buyquantity, buytime)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(tb.Botid, tb.Orderid, tb.Baseasset+tb.Quoteasset, tb.Orderstatus, tb.Buyprice, tb.Buyquantity, tb.Buytime)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

// func InsertTestingSell(entryid int, sellprice float64, sellquantity string, selltime int, orderstatus string) error {
func InsertTestingSell(ts TestingSell) error {
	query := `
		UPDATE testing
		SET sellprice = ?,
			selltime = ?,
			sellquantity = ?,
			orderstatus = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(ts.Sellprice, ts.Selltime, ts.Sellquantity, ts.Orderstatus, ts.Entryid)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetTesting(botid int) ([]AlgoTesting, error) {
	query := `
		SELECT id, botid, orderid, ticket, orderstatus, buyprice, buyquantity, buytime FROM testing
		WHERE sellprice IS NULL
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

		// err := rows.Scan(&algo.Id, &algo.Botid, &algo.Ticket, &algo.Buyprice, &algo.Buytime, &algo.Sellprice, &algo.Selltime)
		err := rows.Scan(&algo.Id, &algo.Botid, &algo.Ticket, &algo.Buyprice, &algo.Buytime)
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
