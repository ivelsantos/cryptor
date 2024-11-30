package models

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(filename string) error {
	var err error

	db, err = sql.Open("sqlite", filename)
	if err != nil {
		return err
	}

	err = walMode(db)
	if err != nil {
		return err
	}

	err = busyTimeOut(db)
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

	err = createTransactionsTable(db)
	if err != nil {
		return err
	}

	err = createAlgoStatsView(db)
	if err != nil {
		return err
	}

	return db.Ping()
}

func walMode(db *sql.DB) error {
	_, err := db.Exec("PRAGMA journal_mode = wal")
	if err != nil {
		return fmt.Errorf("Failed to activate wal2 mode: %v", err)
	}

	return nil
}

func busyTimeOut(db *sql.DB) error {
	_, err := db.Exec("PRAGMA busy_timeout = 5000")
	if err != nil {
		return fmt.Errorf("Failed set busy_timeout: %v", err)
	}

	return nil
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
		ticket TEXT NOT NULL,
		buyvalue REAL NOT NULL,
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

func createTransactionsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS transactions (
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

func createAlgoStatsView(db *sql.DB) error {

	query := `
		CREATE VIEW IF NOT EXISTS algo_stats
		AS
		WITH temp_testing AS (
		SELECT *
		  FROM (
		           SELECT botid,
		                  ((sellvalue - (sellvalue * 0.001)) - (buyvalue + (buyvalue * 0.001))) / buyvalue AS return,
		                  selltime,
		                  buytime - LAG(selltime) OVER (PARTITION BY botid ORDER BY id) AS buytimelength,
		                  selltime - buytime AS selltimelength,
		                  selltime - LAG(selltime) OVER (PARTITION BY botid ORDER BY id) AS tradetimelength
		             FROM testing
		            ORDER BY id
		       )
		WHERE buytimelength IS NOT NULL
		AND return IS NOT NULL
		),

		bot_stats AS (
    SELECT 
        botid,
        SUM(return) AS total_return,
        AVG(return) AS average_return_per_trade,
        (SUM(return) / (SUM(tradetimelength) / 86400.0)) AS average_return_per_day,
        SUM(CASE WHEN return > 0 THEN 1 ELSE 0 END) * 100.0 / COUNT(*) AS success_rate,
        MAX(return) AS max_return,
        MIN(return) AS min_return,
        SUM(tradetimelength) / COUNT(*) AS average_trade_time
    FROM 
        temp_testing
    GROUP BY 
        botid
),
drawdown_calc AS (
    SELECT 
        botid,
        MAX(max_return - min_return) AS max_drawdown
    FROM 
        bot_stats
    GROUP BY 
        botid
)
SELECT 
    bs.botid,
    bs.total_return,
    bs.average_return_per_trade,
    bs.average_return_per_day,
    bs.success_rate,
    dc.max_drawdown,
    bs.average_trade_time
FROM 
    bot_stats AS bs
JOIN 
    drawdown_calc AS dc ON bs.botid = dc.botid;
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("Failed to create table: %v", err)
	}

	return nil
}

func checkSqlBusy(err error) bool {
	if err != nil && strings.Contains(err.Error(), "database is locked (5) (SQLITE_BUSY)") {
		return true
	}

	return false
}
