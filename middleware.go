package main

import (
	"fmt"
	"net/http"
)

func authenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "playerData")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			_, err = fmt.Fprintf(w, `
					<p>Session expired, please re-login</p>
						<form action="/login">
							<button type="submit">Login</button>
						</form>`)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func inGameMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "playerData")
		if inGame, ok := session.Values["inGame"].(bool); !ok || !inGame {
			_, err = fmt.Fprintf(w, `
					<p>You are not currently in a Game session, press play button in the profile</p>
						<form action="/profileDashboard">
							<button type="submit">Profile</button>
						</form>`)
			return
		}
		next.ServeHTTP(w, r)
	}
}
