package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sort"
)

var activeUsers []User
var gameQueue []string
var leaderboardQueue []string

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

func getUserListFromUsernames(usernames []string) []User {
	var users []User
	for _, username := range usernames {
		for _, user := range activeUsers {
			if user.Username == username {
				var tempUser = User{
					Username:        user.Username,
					GamePoints:      user.GamePoints,
					TotalPoints:     user.TotalPoints,
					GamesPlayed:     user.GamesPlayed,
					CurrentQuestion: user.CurrentQuestion,
					Conn:            user.Conn,
				}
				users = append(users, tempUser)
			}
		}
	}
	return users
}

func sortUsers(users []User) []User {
	sort.Slice(users, func(i, j int) bool {
		return users[i].GamePoints > users[j].GamePoints
	})
	return users
}
