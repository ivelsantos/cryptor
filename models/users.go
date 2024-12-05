package models

import (
	"database/sql"
	"fmt"
)

type Account struct {
	Name           string
	ApiKey         string
	SecretKey      string
	ApiKey_test    string
	SecretKey_test string
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

		err := rows.Scan(&account.Name, &account.ApiKey, &account.SecretKey, &account.ApiKey_test, &account.SecretKey_test)
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

func GetAccountByName(name string) (Account, error) {
	query := `SELECT * FROM accounts WHERE name = ?`
	var account Account

	row := db.QueryRow(query, name)

	err := row.Scan(&account.Name, &account.ApiKey, &account.SecretKey, &account.ApiKey_test, &account.SecretKey_test)
	if err != nil {
		if err == sql.ErrNoRows {
			return account, fmt.Errorf("No account found with name: %s", name)
		}
		return account, fmt.Errorf("Failed to retrieve account: %v", err)
	}

	return account, nil
}

func InsertAccount(account Account) error {
	query := `
		INSERT INTO accounts (name, apikey, secretkey, apikey_test, secretkey_test)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.Name, account.ApiKey, account.SecretKey, account.ApiKey_test, account.SecretKey_test)
	if err != nil {
		return fmt.Errorf("Failed to execute statement: %v", err)
	}

	return nil
}

func DeleteUser(user string) error {
	query := `
		DELETE FROM accounts
		WHERE name = ?
		`

	result, err := db.Exec(query, user)
	if err != nil {
		return fmt.Errorf("Failed to delete algo: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Failed to retrieve number of rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("No rows found for user %s", user)
	}

	return nil
}
