package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
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

		newUser(username)

		http.Redirect(w, r, "/profileDashboard", http.StatusSeeOther)
	}

	t, err := template.ParseFiles("templates/login.html")
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

	t, err := template.ParseFiles("templates/register.html")
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
