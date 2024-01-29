package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("Templates/welcome.html")
	if err != nil {
		fmt.Println("error parsing welcome.html: ", err)
		http.Error(w, "error parsing welcome page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("error executing welcome.html: ", err)
		http.Error(w, "error executing welcome page", http.StatusInternalServerError)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var hashedPassword string
		err = db.QueryRow("select password from users where username = ?", username).Scan(&hashedPassword)
		if err != nil {
			fmt.Println("error querying the hashedpassword during login: ", err)
			http.Error(w, "username not in the database", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			fmt.Println("error in comparing hashpsw and psw: ", err)
			http.Error(w, "wrong password", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "playerData")
		session.Values["username"] = username
		session.Values["authenticated"] = true
		session.Options = &sessions.Options{
			MaxAge: 300,
		}
		session.Save(r, w)

		newActiveUser(username)

		http.Redirect(w, r, "/profileDashboard", http.StatusSeeOther)
	}

	t, err := template.ParseFiles("Templates/login.html")
	if err != nil {
		fmt.Println("error parsing login.html: ", err)
		http.Error(w, "error parsing login page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("error executing login.html: ", err)
		http.Error(w, "error executing login page", http.StatusInternalServerError)
		return
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var countUsernames int
		err = db.QueryRow("select count(*) from users where username = ?", username).Scan(&countUsernames)
		if err != nil {
			fmt.Println("error in queryrow searching for username already used: ", err)
			http.Error(w, "error in queryrow searching for username already used", http.StatusInternalServerError)
			return
		}
		if countUsernames > 0 {
			http.Error(w, "Username already exists, please change username", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("error hashing register password: ", err)
			http.Error(w, "error hashing register password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("insert into users (username, password, totalPoints, gamesPlayed) VALUES (?, ?, ?, ?)", username, hashedPassword, 0, 0)
		if err != nil {
			fmt.Println("error inserting new users in db: ", err)
			http.Error(w, "error inserting new users in db", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	t, err := template.ParseFiles("Templates/register.html")
	if err != nil {
		fmt.Println("error parsing register.html: ", err)
		http.Error(w, "error parsing register page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("error executing register.html: ", err)
		http.Error(w, "error executing register page", http.StatusInternalServerError)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username, _ := session.Values["username"].(string)

	save(username)
	removeActiveUser(username)

	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func profileDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "playerData")
	if err != nil {
		fmt.Println(err)
	}

	t, err := template.ParseFiles("Templates/playerDashboard.html")
	if err != nil {
		fmt.Println("error parsing playerDashboard.html: ", err)
		http.Error(w, "error parsing playerDashboard page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, getUserFromUsername(session.Values["username"].(string)))
	if err != nil {
		fmt.Println("error executing playerDashboard.html: ", err)
		http.Error(w, "error executing playerDashboard page", http.StatusInternalServerError)
		return
	}

}

func addGameQueueHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)
	session.Values["inGame"] = true
	session.Save(r, w)
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

	t, err := template.ParseFiles("templates/queue.html")
	if err != nil {
		fmt.Println("error parsing queue.html: ", err)
		http.Error(w, "error parsing queue page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"QueueList": gameQueue,
	})
	if err != nil {
		fmt.Println("error executing queue.html: ", err)
		http.Error(w, "error executing queue page", http.StatusInternalServerError)
		return
	}
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "playerData")
	username := session.Values["username"].(string)

	if getCurrentQuestion(username) == -1 {
		http.Error(w, "Error getting the current question", http.StatusInternalServerError)
	}

	if getCurrentQuestion(username) >= QuestionPerRound {
		http.Redirect(w, r, "/addLeaderboardQueue", http.StatusSeeOther)
		return
	}

	currenTime := time.Now()
	session.Values["timeStart"] = int(currenTime.Unix())
	session.Save(r, w)

	question := Question{}
	var options string
	db.QueryRow("select text, options, correct_idx from questions where id = ?", randomIntSlice[getCurrentQuestion(username)]).Scan(&question.Text, &options, &question.CorrectIdx)
	question.Options = strings.Split(options, ",")
	getUserFromUsername(username).CurrentQuestion++

	t, err := template.ParseFiles("Templates/quiz.html")
	if err != nil {
		fmt.Println("error parsing quiz.html: ", err)
		http.Error(w, "error parsing quiz page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, question)
	if err != nil {
		fmt.Println("error executing quiz.html: ", err)
		http.Error(w, "error executing quiz page", http.StatusInternalServerError)
		return
	}
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
	leaderboardQueue = append(leaderboardQueue, username)

	for _, user := range activeUsers {
		if user.Conn != nil {
			err := user.Conn.WriteMessage(websocket.TextMessage, []byte(username))
			if err != nil {
				log.Println("Error writing to WebSocket:", err)
			}
		}
	}

	http.Redirect(w, r, "/leaderboardQueue", http.StatusSeeOther)
}

func leaderboardQueueHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/leaderboardQueue.html")
	if err != nil {
		fmt.Println("error parsing leaderboardQueue.html: ", err)
		http.Error(w, "error parsing leaderboardQueue page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"QueueList": leaderboardQueue,
	})
	if err != nil {
		fmt.Println("error executing leaderboardQueue.html: ", err)
		http.Error(w, "error executing leaderboardQueue page", http.StatusInternalServerError)
		return
	}
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/leaderboard.html")
	if err != nil {
		fmt.Println("error parsing leaderboard.html: ", err)
		http.Error(w, "error parsing leaderboard page", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, sortUsers(getUserListFromUsernames(leaderboardQueue)))
	if err != nil {
		fmt.Println("error executing leaderboard.html: ", err)
		http.Error(w, "error executing leaderboard page", http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "playerData")
	save(session.Values["username"].(string))

	//	TODO: resetQueues and better savePoints

}
