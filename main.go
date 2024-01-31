package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *sql.DB
var err error
var PlayerNumber = flag.Int("p", 3, "set the number of Players needed to start the Quiz")
var QuestionPerRound = flag.Int("q", 5, "set the number of Questions per round")
var TotalQuestion = 50
var randomIntSlice []int

func main() {
	checkFlags()
	db, err = sql.Open("sqlite3", "Data/database.db")
	if err != nil {
		log.Println("Error opening db: ", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", authenticationMiddleware(logoutHandler))
	http.HandleFunc("/profileDashboard", authenticationMiddleware(profileDashboardHandler))
	http.HandleFunc("/addGameQueue", authenticationMiddleware(addGameQueueHandler))
	http.HandleFunc("/gameQueue", authenticationMiddleware(gameQueueHandler))
	http.HandleFunc("/wsGameQueue", authenticationMiddleware(wsGameQueueHandler))
	http.HandleFunc("/quiz", inGameMiddleware(authenticationMiddleware(quizHandler)))
	http.HandleFunc("/submitAnswer", inGameMiddleware(authenticationMiddleware(submitAnswerHandler)))
	http.HandleFunc("/addLeaderboardQueue", inGameMiddleware(authenticationMiddleware(addLeaderboardQueue)))
	http.HandleFunc("/leaderboardQueue", inGameMiddleware(authenticationMiddleware(leaderboardQueueHandler)))
	http.HandleFunc("/wsLeaderboardQueue", inGameMiddleware(authenticationMiddleware(wsLeaderboardQueueHandler)))
	http.HandleFunc("/leaderboard", inGameMiddleware(authenticationMiddleware(leaderboardHandler)))
	http.HandleFunc("/save", inGameMiddleware(authenticationMiddleware(saveHandler)))

	fmt.Println("server starting at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("error listenAndServe: ", err)
		return
	}

}
