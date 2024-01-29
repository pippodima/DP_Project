package main

import "github.com/gorilla/sessions"

var store = sessions.NewCookieStore([]byte("SuperSecretKey"))
