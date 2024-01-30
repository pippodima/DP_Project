package main

import (
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var store = sessions.NewCookieStore([]byte("SuperSecretKey"))

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsGameQueueHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("errore nell'upgrade in ws handler: ", err)
	}

	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	getUserFromUsername(username).Conn = conn

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	if len(gameQueue) >= *PlayerNumber {
		randomIntSlice = getRandomSlice(*QuestionPerRound)
		session.Values["inGame"] = true
		err := session.Save(r, w)
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
		for _, user := range activeUsers {
			if user.Conn != nil {
				err = user.Conn.WriteMessage(websocket.TextMessage, []byte("/quiz"))
				if err != nil {
					log.Println(err)
				}
			}
		}

		leaderboardQueue = make([]string, len(gameQueue))
		leaderboard = make([]string, len(gameQueue))
		copy(leaderboardQueue, gameQueue)
		copy(leaderboard, gameQueue)
		gameQueue = []string{}

	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

	}

}

func wsLeaderboardQueueHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("errore nell'upgrade in ws handler: ", err)
	}

	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	getUserFromUsername(username).Conn = conn

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	if len(leaderboardQueue) == 0 {
		time.Sleep(1 * time.Second)
		for _, user := range activeUsers {
			if user.Conn != nil {
				err = user.Conn.WriteMessage(websocket.TextMessage, []byte("/leaderboard"))
				if err != nil {
					log.Println(err)
				}
			}
		}

	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

	}

}
