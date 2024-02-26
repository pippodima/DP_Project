package main

import (
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "welcome.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string

		hashedPassword, err = getHashedPsw(username)
		if err != nil {
			log.Println("error querying the hashedpassword during login: ", err)
			http.Error(w, "username not in the database", http.StatusInternalServerError)
			return
		}

		err = compareHashPsw(password, hashedPassword)
		if err != nil {
			log.Println("error in comparing hashpsw and psw: ", err)
			http.Error(w, "wrong password", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "playerData")
		session.Values["username"] = username
		session.Values["authenticated"] = true
		session.Options = &sessions.Options{
			MaxAge: 300,
		}
		err := session.Save(r, w)
		if err != nil {
			log.Println(err)
			return
		}

		newActiveUser(username)

		http.Redirect(w, r, "/profileDashboard", http.StatusSeeOther)
	}

	renderTemplate(w, "login.html", nil)

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		exist, err := usernameAlreadyExist(username)
		if err != nil {
			log.Println("error in queryrow searching for username: ", err)
			http.Error(w, "error in queryrow searching for username", http.StatusInternalServerError)
			return
		}
		if exist {
			http.Error(w, "Username already exists, please change username", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := createHashPsw(password)
		if err != nil {
			log.Println("error hashing register password: ", err)
			http.Error(w, "error hashing register password", http.StatusInternalServerError)
			return
		}

		err = DBcreateUser(username, hashedPassword)
		if err != nil {
			log.Println("error inserting new users in db: ", err)
			http.Error(w, "error inserting new users in db", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "register.html", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username, _ := session.Values["username"].(string)

	removeActiveUser(username)

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)

		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func profileDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "playerData")
	if err != nil {
		log.Println(err)
	}

	renderTemplate(w, "playerDashboard.html", getUserFromUsername(session.Values["username"].(string)))

}

func addGameQueueHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	session.Values["inGame"] = true
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)

		return
	}
	gameQueue = append(gameQueue, username)

	for _, user := range activeUsers {
		if user.Conn != nil {
			err := user.Conn.WriteMessage(websocket.TextMessage, []byte(username))
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
			}
		}
	}

	http.Redirect(w, r, "/gameQueue", http.StatusSeeOther)

}

func gameQueueHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "queue.html", map[string]interface{}{
		"QueueList": gameQueue,
	})
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)

	if getCurrentQuestion(username) == -1 {
		http.Error(w, "Error getting the current question", http.StatusInternalServerError)
	}

	if getCurrentQuestion(username) >= *QuestionPerRound {
		http.Redirect(w, r, "/addLeaderboardQueue", http.StatusSeeOther)
		return
	}

	currenTime := time.Now()
	session.Values["timeStart"] = int(currenTime.Unix())
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
		return
	}

	question, err := DBgetQuestionFromId(randomIntSlice[getCurrentQuestion(username)])
	if err != nil {
		log.Println(err)
		http.Error(w, "error querying question", http.StatusInternalServerError)
		return
	}

	getUserFromUsername(username).CurrentQuestion++

	renderTemplate(w, "quiz.html", question)
}

func submitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.FormValue("selectedOption") == r.FormValue("correctIndex") {
			session, _ := store.Get(r, "playerData")
			username := session.Values["username"].(string)
			addCorrectPoint(username)
			addTimePoint(username, session.Values["timeStart"].(int))
		}
		http.Redirect(w, r, "/quiz", http.StatusSeeOther)
	}
}

func addLeaderboardQueue(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	leaderboardQueue = remove(leaderboardQueue, username)

	for _, user := range activeUsers {
		if user.Conn != nil {
			err := user.Conn.WriteMessage(websocket.TextMessage, []byte("/clear"))
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
			}
			for _, name := range leaderboardQueue {
				err := user.Conn.WriteMessage(websocket.TextMessage, []byte(name))
				if err != nil {
					log.Println("Error writing to WebSocket:", err)
				}
			}
		}
	}

	http.Redirect(w, r, "/leaderboardQueue", http.StatusSeeOther)
}

func leaderboardQueueHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "leaderboardQueue.html", map[string]interface{}{
		"QueueList": leaderboardQueue,
	})
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "leaderboard.html", sortUsers(getUserListFromUsernames(leaderboard)))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	session.Values["inGame"] = false
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
		return
	}
	save(username)
	remove(leaderboard, username)

	action := r.FormValue("action")
	switch action {
	case "logout":
		http.Redirect(w, r, "/logout", http.StatusSeeOther)
	case "playAgain":
		http.Redirect(w, r, "/addGameQueue", http.StatusSeeOther)
	case "profile":
		http.Redirect(w, r, "/profileDashboard", http.StatusSeeOther)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
	}
}
