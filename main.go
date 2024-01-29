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
var randomIntSlice []int

func main() {
	db, err = sql.Open("sqlite3", "Data/database.db")
	if err != nil {
		fmt.Println("Error opening db: ", err)
	}
	defer db.Close()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", authenticationMiddleware(logoutHandler))
	http.HandleFunc("/profileDashboard", authenticationMiddleware(profileDashboardHandler))
	http.HandleFunc("/addGameQueue", authenticationMiddleware(addGameQueueHandler))
	http.HandleFunc("/gameQueue", authenticationMiddleware(gameQueueHandler))
	http.HandleFunc("/wsGameQueue", authenticationMiddleware(wsGameQueueHandler))
	http.HandleFunc("/quiz", authenticationMiddleware(quizHandler))
	http.HandleFunc("/submitAnswer", authenticationMiddleware(submitAnswerHandler))
	http.HandleFunc("/addLeaderboardQueue", authenticationMiddleware(addLeaderboardQueue))
	http.HandleFunc("/leaderboardQueue", authenticationMiddleware(leaderboardQueueHandler))
	http.HandleFunc("/wsLeaderboardQueue", authenticationMiddleware(wsLeaderboardQueueHandler))
	http.HandleFunc("/leaderboard", authenticationMiddleware(leaderboardHandler))

	fmt.Println("server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("error listenAndServe: ", err)
		return
	}

}
