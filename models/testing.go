package models

import (
	"database/sql"
	"fmt"
)

type AlgoTesting struct {
	Id        int
	Botid     int
	Ticket    string
	Buyprice  float64
	Buytime   int
	Sellprice float64
	Selltime  int
}

func InsertTestingBuy(botid int, ticket string, buyprice float64, buytime int) error {
	query := `
		INSERT INTO testing (botid, ticket, buyprice, buytime)
		VALUES (?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(botid, ticket, buyprice, buytime)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func InsertTestingSell(entryid int, sellprice float64, selltime int) error {
	query := `
		UPDATE testing
		SET sellprice = ?,
			selltime = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(sellprice, selltime, entryid)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetTesting(botid int) ([]AlgoTesting, error) {
	query := `
		SELECT id, botid, ticket, buyprice, buytime FROM testing
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
