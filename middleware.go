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
