package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // go tools makes me comment this...
)

func initDB() {
	fmt.Println("Open DB")
	// Open PG
	db, err := sql.Open("postgres", "user=jahn dbname=pqgotest sslmode=disable")
	if err != nil {
		writeErr(err)
	}

	// Test ping
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed DB Ping")
		writeErr(err)
	} else {
		writeLog("Success DB Ping")
	}
}

/*
	What do I want stored in PG?
	- User's game profile, stats, etc.
	- ...

	Tables?
	- User/Player
	- Each game
	- ...
*/
