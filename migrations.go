package main

import (
	"database/sql"
	"time"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	s := &Storage{db: db}
	if err := s.applyMigrations(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Storage) applyMigrations() error {
	// create new table on database
	var err error
	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS seen (url TEXT, date TEXT, summary TEXT)")

	if err != nil {
		return err
	}
	err = s.addSummaryColumnIfNotExists()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) addSummaryColumnIfNotExists() error {
	tx, err := s.db.Begin()
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

func (s *Storage) MarkAsSeen(url, summary string) error {
	today := time.Now().Format("2006-01-02")
	_, err := s.db.Exec("INSERT INTO seen(url, date, summary) values(?,?,?)", url, today, summary)
	return err
}

// IsSeen returns (seen, seen_today, summary)
func (s *Storage) IsSeen(link string) (bool, bool, string) {
	var urlStr, date, summary sql.NullString
	err := s.db.QueryRow("SELECT url, date, summary FROM seen WHERE url=?", link).Scan(&urlStr, &date, &summary)

	if err != nil {
		return false, false, ""
	}

	today := time.Now().Format("2006-01-02")
	isSeen := urlStr.Valid && date.String != today
	isSeenToday := urlStr.Valid && date.String == today

	return isSeen, isSeenToday, summary.String
}

func (s *Storage) Close() {
	s.db.Close()
}
