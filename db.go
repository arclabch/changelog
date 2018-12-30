// Copyright © 2018 ARClab, Lionel Riem - https://arclab.ch/
//
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (

	// Schema for the database
	dbSchema = `CREATE TABLE entries (
		timestamp	DATETIME NOT NULL DEFAULT (datetime('now','localtime')),
		user	TEXT NOT NULL,
		entry	TEXT NOT NULL,
		PRIMARY KEY(timestamp)
	);`

	// Insert Query
	dbInsert = "INSERT INTO entries (user, entry) VALUES (?, ?);"

	// Lookup Query Construction
	dbSelect = "SELECT timestamp, user, entry FROM entries "
	dbCount  = "SELECT count(*) FROM entries WHERE timestamp IN (SELECT timestamp FROM entries "
	dbLength = "SELECT max(length(user)) as len FROM entries WHERE timestamp IN (SELECT timestamp FROM entries "
	dbUser   = "WHERE user = ? "
	dbEnd    = "ORDER BY timestamp DESC LIMIT ?"
)

type entry struct {
	timestamp time.Time
	user      string
	message   string
}

var db *sql.DB
var dbOpen = false

// initDb opens the database
func initDb() error {

	var err error

	// Start by checking if the DB exists
	dbExists := checkDbExists()
	if !dbExists && !userIsRoot() {
		return errors.New("can't create a new database if not root")
	}

	// Open/create the database
	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("can't access the database")
	}

	// Check that we're all good
	err = db.Ping()
	if err != nil {
		return err
	}

	// At this point, the DB is open
	dbOpen = true

	// If new DB and root, create table and set chmod rights
	if !dbExists {
		_, err = db.Exec(dbSchema)

		err = setDbRights()
		if err != nil {
			return err
		}
	}

	return nil
}

// addEntry adds an entry.
func addEntry(user string, message string) error {

	// Prepare the insert query
	stmt, err := db.Prepare(dbInsert)
	if err != nil {
		return err
	}

	// Run the query
	_, err = stmt.Exec(user, message)
	if err != nil {
		return err
	}

	stmt.Close()
	return nil
}

// getEntries fetches entries from the database
func getEntries(limit int, user string) ([]entry, int, int, error) {

	var count int
	var length int
	var answer []entry
	var err error
	var rows *sql.Rows

	// Prepare the queries
	queryCount := dbCount
	queryLength := dbLength
	querySelect := dbSelect

	if user != "" {
		queryCount = queryCount + dbUser
		queryLength = queryLength + dbUser
		querySelect = querySelect + dbUser
	}

	queryCount = queryCount + dbEnd + ");"
	queryLength = queryLength + dbEnd + ");"
	querySelect = querySelect + dbEnd + ";"

	// Start by getting a count
	if user != "" {
		err = db.QueryRow(queryCount, user, limit).Scan(&count)
	} else {
		err = db.QueryRow(queryCount, limit).Scan(&count)
	}

	// Check the result
	if count == 0 {
		return answer, 0, 0, nil
	} else if err != nil {
		return answer, 0, 0, err
	}

	// Get the length
	if user != "" {
		err = db.QueryRow(queryLength, user, limit).Scan(&length)
	} else {
		err = db.QueryRow(queryLength, limit).Scan(&length)
	}

	// Check the result
	if err != nil {
		return answer, 0, 0, err
	}

	// Get the entries
	if user != "" {
		rows, err = db.Query(querySelect, user, limit)
	} else {
		rows, err = db.Query(querySelect, limit)
	}

	// Check the result
	if err != nil {
		return answer, 0, 0, err
	}

	// Fetch the results into answer
	for rows.Next() {
		item := entry{}
		err2 := rows.Scan(&item.timestamp, &item.user, &item.message)
		if err2 != nil {
			return answer, 0, 0, err
		}
		answer = append(answer, item)
	}

	// Close the query when finished.
	rows.Close()

	return answer, count, length, nil
}
