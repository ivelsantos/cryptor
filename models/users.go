package models

import (
	"database/sql"
	"fmt"
)

type Account struct {
	Name      string
	ApiKey    string
	SecretKey string
}

func GetAccounts() ([]Account, error) {
	query := `SELECT * FROM accounts`
	var accounts []Account

	rows, err := db.Query(query)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("Failed to retrieve accounts: %v", err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		var account Account

		err := rows.Scan(&account.Name, &account.ApiKey, &account.SecretKey)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func InsertAccount(account Account) error {
	query := `
		INSERT INTO accounts (name, apikey, secretkey)
		VALUES (?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Name, account.ApiKey, account.SecretKey)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}
