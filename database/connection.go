package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" //
)

var db *sql.DB

func periodicPing() {
	for {
		time.Sleep(time.Minute * 5)

		err := db.Ping()
		if err != nil {
			log.Println(err)
		}
	}
}

//Connect sets up the database connection
func Connect() {

	//check for already open connection, close it
	if db != nil {
		db.Close()
	}

	//if NOT using sqlx:
	var err error
	db, err = sql.Open("mysql", "cs361_connorsa:0400@tcp(classmysql.engr.oregonstate.edu:3306)/cs361_connorsa")
	//if using sqlx:
	//db, err = sqlx.Open("mysql", "cs361_connorsa:0400@tcp(classmysql.engr.oregonstate.edu:3306)/cs361_connorsa")

	if err != nil {
		log.Fatal(err)
	}

	//check database connection, bail out if bad
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	go periodicPing()

}

//DestroyConnection will destroy the database connection
func DestroyConnection() {

	db.Close()

}
