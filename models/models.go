package models

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(filename string) error {
	var err error

	db, err = sql.Open("sqlite", filename)
	if err != nil {
		return err
	}

	err = createAccountTable(db)
	if err != nil {
		return err
	}

	err = createAlgosTable(db)
	if err != nil {
		return err
	}

	err = createTestingTable(db)
	if err != nil {
		return err
	}

	err = createTestingFixedTable(db)
	if err != nil {
		return err
	}

	return db.Ping()
}

func createAccountTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		name TEXT PRIMARY KEY,
		apikey TEXT NOT NULL,
		secretkey TEXT NOT NULL,
		apikey_test TEXT NOT NULL,
		secretkey_test TEXT NOT NULL
	) WITHOUT ROWID`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create table: %v", err)
	}

	return nil
}

func createAlgosTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS algos (
		id INTEGER PRIMARY KEY,
		owner TEXT NOT NULL,
		name TEXT NOT NULL,
		created INTEGER NOT NULL,
		buycode TEXT NOT NULL,
		state TEXT NOT NULL,
		base_asset TEXT NOT NULL,
		quote_asset TEXT NOT NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create table: %v", err)
	}

	return nil
}

func createTestingTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS testing (
		id INTEGER PRIMARY KEY,
		botid INTEGER NOT NULL,
		orderid INTEGER NOT NULL,
		ticket TEXT NOT NULL,
		orderstatus TEXT NOT NULL,
		buyvalue REAL NOT NULL,
		buyquantity REAL NOT NULL,
		buytime INTEGER NOT NULL,
		sellvalue REAL,
		selltime INTEGER
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create table: %v", err)
	}

	return nil
}

func createTestingFixedTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS testing_fixed (
		id INTEGER PRIMARY KEY,
		botid INTEGER NOT NULL,
		return REAL NOT NULL,
		selltime INTEGER NOT NULL,
		buytimelength INTEGER NOT NULL,
		selltimelength INTEGER NOT NULL,
		tradetimelength INTEGER NOT NULL
		
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create table: %v", err)
	}

	return nil
}
