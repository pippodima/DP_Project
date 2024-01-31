package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type Question struct {
	Text       string
	Options    []string
	CorrectIdx int
}

func getRandomSlice(length int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	randomSlice := make([]int, length)
	usedNumbers := make(map[int]bool)

	for i := 0; i < length; {
		num := r.Intn(TotalQuestion) + 1

		if !usedNumbers[num] {
			randomSlice[i] = num
			usedNumbers[num] = true
			i++
		}
	}

	return randomSlice
}

func addCorrectPoint(username string) {
	for i, user := range activeUsers {
		if user.Username == username {
			activeUsers[i].GamePoints += 3
			return
		}
	}
}

func addTimePoint(username string, timeStart int) {
	for i, user := range activeUsers {
		if user.Username == username {
			currenTime := time.Now()
			increment := 5 - (int(currenTime.Unix()) - timeStart)
			if increment > 0 {
				activeUsers[i].GamePoints += increment
			}
			return
		}
	}
}

func save(username string) {

	getUserFromUsername(username).GamesPlayed++
	getUserFromUsername(username).TotalPoints += getUserFromUsername(username).GamePoints

	_, err = db.Exec("update users set totalPoints = totalPoints + ? where username = ?", getUserFromUsername(username).GamePoints, username)
	if err != nil {
		log.Println("error in updating the total points: ", err)
	}

	_, err = db.Exec("update users set gamesPlayed = ? where username = ?", getUserFromUsername(username).GamesPlayed, username)
	if err != nil {
		log.Println("error in updating the total games played: ", err)
	}

	getUserFromUsername(username).GamePoints = 0
	getUserFromUsername(username).CurrentQuestion = 0

}

func remove(list []string, username string) []string {
	for i, v := range list {
		if v == username {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func checkFlags() {
	help := flag.Bool("help", false, "Show help message")
	verbose := flag.Bool("verbose", false, "Shows (if there are) log error messages")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if !*verbose {
		log.SetOutput(io.Discard)
	}

}
