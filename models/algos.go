package models

import (
	"database/sql"
	"fmt"
)

type Algor struct {
	Id      int
	Owner   string
	Name    string
	Created int
	Buycode string
	State   string
}

func InsertAlgo(algor Algor) error {
	query := `
		INSERT INTO algos (owner, name, created, buycode, state)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(algor.Owner, algor.Name, algor.Created, algor.Buycode, algor.State)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func GetAlgos(owner string) ([]Algor, error) {
	query := `SELECT * FROM algos WHERE owner = ?`
	var algos []Algor

	rows, err := db.Query(query, owner)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo Algor

		err := rows.Scan(&algo.Id, &algo.Owner, &algo.Name, &algo.Created, &algo.Buycode, &algo.State)
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

func GetAlgosState(state string) ([]Algor, error) {
	query := `SELECT * FROM algos WHERE state = ?`
	var algos []Algor

	rows, err := db.Query(query, state)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve algos: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var algo Algor

		err := rows.Scan(&algo.Id, &algo.Owner, &algo.Name, &algo.Created, &algo.Buycode, &algo.State)
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

func DeleteAlgo(id int, owner string) error {
	query := `
		DELETE FROM algos
		WHERE id = ? AND owner = ?`

	result, err := db.Exec(query, id, owner)
	if err != nil {
		return fmt.Errorf("Failed to delete algo: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to retrieve number of rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No rows found for id %d and owner %s", id, owner)
	}

	return nil
}

func UpdateAlgoState(state string, id int) error {
	query := `
		UPDATE algos
		SET state = ?
		WHERE id = ?
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(state, id)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}
