package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var activeUsers []User
var queue []string

type User struct {
	Username        string
	GamePoints      int
	TotalPoints     int
	GamesPlayed     int
	CurrentQuestion int
	Conn            *websocket.Conn
}

func newActiveUser(username string) {
	var totPoints int
	var gamesPlayed int
	err = db.QueryRow("select totalPoints from users where username = ?", username).Scan(&totPoints)
	if err != nil {
		fmt.Println("error querying total points during user creation: ", err)
		return
	}
	err = db.QueryRow("select gamesPlayed from users where username = ?", username).Scan(&gamesPlayed)
	if err != nil {
		fmt.Println("error querying total points during user creation: ", err)
		return
	}
	activeUsers = append(activeUsers, User{
		Username:        username,
		GamePoints:      0,
		TotalPoints:     totPoints,
		GamesPlayed:     gamesPlayed,
		CurrentQuestion: 0,
		Conn:            nil,
	})
}

func removeActiveUser(username string) {
	for i, user := range activeUsers {
		if user.Username == username {
			activeUsers = append(activeUsers[:i], activeUsers[i+1:]...)
			return
		}
	}
}

func getUserFromUsername(username string) *User {
	for i, user := range activeUsers {
		if user.Username == username {
			return &activeUsers[i]
		}
	}
	return nil
}

func getCurrentQuestion(username string) int {
	for i, user := range activeUsers {
		if user.Username == username {
			return activeUsers[i].CurrentQuestion
		}
	}
	return -1
}
