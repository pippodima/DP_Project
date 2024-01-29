package main

import (
	"fmt"
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
		fmt.Println("errore nell'upgrade in ws handler: ", err)
	}

	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	getUserFromUsername(username).Conn = conn

	defer conn.Close()

	if len(queue) >= PlayerNumber {
		randomIntSlice = getRandomSlice(QuestionPerRound)
		time.Sleep(1 * time.Second)
		for _, user := range activeUsers {
			if user.Conn != nil {
				err = user.Conn.WriteMessage(websocket.TextMessage, []byte("/quiz"))
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
