package structures

import "github.com/gorilla/websocket"

type Question struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	Options    string `json:"options"`
	CorrectIdx int    `json:"correct_idx"`
}

type User struct {
	Username        string
	GamePoints      int
	TotalPoints     int
	GamesPlayed     int
	CurrentQuestion int
	Conn            *websocket.Conn
}
