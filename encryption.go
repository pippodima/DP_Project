package main

import (
	"golang.org/x/crypto/bcrypt"
)

func getHashedPsw(username string) (string, error) {
	var hashedPassword string
	hashedPassword, err = DBgetHashPswFromUsername(username)
	return hashedPassword, err
}

func compareHashPsw(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func createHashPsw(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
