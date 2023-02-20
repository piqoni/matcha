package main

import "database/sql"

func applyMigrations(db *sql.DB) error {
	// create new table on database
	var err error
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS seen (url TEXT, date TEXT, summary TEXT)")

	if err != nil {
		return err
	}
	err = addSummaryColumnIfNotExists(db)
	if err != nil {
		return err
	}
	return nil
}

func addSummaryColumnIfNotExists(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if the 'summary' column already exists in the 'seen' table
	rows, err := tx.Query("PRAGMA table_info(seen)")
	if err != nil {
		return err
	}
	defer rows.Close()

	var columnExists bool
	for rows.Next() {
		var cid int
		var name string
		var dataType string
		var notNull int
		var dfltValue interface{}
		var pk int
		err := rows.Scan(&cid, &name, &dataType, &notNull, &dfltValue, &pk)
		if err != nil {
			return err
		}
		if name == "summary" {
			columnExists = true
			break
		}
	}

	// If the 'summary' column doesn't exist, add it to the table
	if !columnExists {
		_, err = tx.Exec("ALTER TABLE seen ADD COLUMN summary TEXT")
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
