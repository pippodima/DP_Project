package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Question struct {
	Text       string
	Options    []string
	CorrectIdx int
}

func getRandomSlice(length int) []int {
	rand.Seed(time.Now().UnixNano())

	randomSlice := make([]int, length)
	usedNumbers := make(map[int]bool)

	for i := 0; i < length; {
		num := rand.Intn(TotalQuestion) + 1

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

	_, err = db.Exec("update users set totalPoints = totalPoints + ? where username = ?", getUserFromUsername(username).GamePoints, username)
	if err != nil {
		fmt.Println("error in updating the total points: ", err)
	}

	_, err = db.Exec("update users set gamesPlayed = ? where username = ?", getUserFromUsername(username).GamesPlayed, username)
	if err != nil {
		fmt.Println("error in updating the total games played: ", err)
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
