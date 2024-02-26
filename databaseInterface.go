package main

import "strings"

func DBgetHashPswFromUsername(username string) (string, error) {
	var hashPsw string
	err = db.QueryRow("select password from users where username = ?", username).Scan(&hashPsw)
	return hashPsw, err
}

func DBdoesUsernameExist(username string) (bool, error) {
	var count int
	err = db.QueryRow("select count(*) from users where username = ?", username).Scan(&count)
	if count != 0 {
		return true, err
	}
	return false, err

}

func DBcreateUser(username string, hashPassword []byte) error {
	_, err = db.Exec("insert into users (username, password, totalPoints, gamesPlayed) VALUES (?, ?, ?, ?)", username, hashPassword, 0, 0)
	return err
}

func DBgetQuestionFromId(id int) (Question, error) {
	question := Question{}
	var options string
	err = db.QueryRow("select text, options, correct_idx from questions where id = ?", id).Scan(&question.Text, &options, &question.CorrectIdx)
	question.Options = strings.Split(options, ",")
	return question, err
}
