package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

var db *sql.DB
var err error
var PlayerNumber = 2
var TotalQuestion = 50
var QuestionPerRound = 3

func main() {
	db, err = sql.Open("sqlite3", "Data/database.db")
	if err != nil {
		fmt.Println("Error opening db: ", err)
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
		    totalPoints INTEGER,  
		    gamesPlayed INTEGER               
		)
	`)
	if err != nil {
		fmt.Println("error creating users table: ", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT,
			options TEXT,
			correct_idx INTEGER
		)
	`)
	if err != nil {
		fmt.Println("error creating questions table: ", err)
	}

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	fmt.Println("server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error listenAndServe: ", err)
		return
	}

}
